# 05 — Test infrastructure

Three pieces of infrastructure that make the test types in [04-test-types.md](04-test-types.md) repeatable, independent, and portable:

1. **Testcontainers** — every test class spins up its own ephemeral Postgres
2. **Seed data** — known fixtures for "this email already exists"-style scenarios
3. **Parameterized CSV** — every validation rule's accept/reject cases live in a CSV that all three stacks consume

## 1. Testcontainers — ephemeral PostgreSQL per test class

**Problem it solves:** the test needs a real Postgres (mocks lie about CHECK constraints, unique violations, regex semantics). But it shouldn't require the developer to manually `docker compose up` first, and CI runners can't share a long-lived database with other jobs.

**How it works:** the test class declares "I need Postgres 18-alpine," and the Testcontainers library starts a fresh Docker container at `@BeforeAll`, exposes it on a random free port, and tears it down at `@AfterAll`. The connection string is injected into the application config.

### Spring Boot

`spring-boot/pom.xml` would add:

```xml
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>postgresql</artifactId>
    <scope>test</scope>
</dependency>
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>junit-jupiter</artifactId>
    <scope>test</scope>
</dependency>
```

Then base test class:

```java
@SpringBootTest
@Testcontainers
abstract class AbstractIT {
    @Container
    static PostgreSQLContainer<?> postgres =
        new PostgreSQLContainer<>("postgres:18-alpine")
            .withDatabaseName("registration")
            .withUsername("registration")
            .withPassword("test");

    @DynamicPropertySource
    static void register(DynamicPropertyRegistry r) {
        r.add("DATABASE_URL", () ->
            "postgres://" + postgres.getUsername() + ":" + postgres.getPassword()
            + "@" + postgres.getHost() + ":" + postgres.getFirstMappedPort()
            + "/" + postgres.getDatabaseName() + "?sslmode=disable");
        r.add("HTTP_PORT", () -> "0");  // random port for @SpringBootTest
    }
}
```

`@Container` lifecycle:
- `static` field → started once per class (`@BeforeAll`)
- non-static field → started before each test method (`@BeforeEach`)

For most cases you want `static` — starting Postgres takes 1–2 seconds; you don't want that per test.

Flyway runs the migration against the container automatically (Spring Boot's `spring.flyway.enabled=true` + the `DATABASE_URL` from `@DynamicPropertySource`).

### Express

```bash
npm i -D @testcontainers/postgresql
```

```javascript
import { PostgreSqlContainer } from '@testcontainers/postgresql';
import { test, before, after } from 'node:test';

let container;
before(async () => {
    container = await new PostgreSqlContainer('postgres:18-alpine')
        .withDatabase('registration')
        .withUsername('registration')
        .withPassword('test')
        .start();
    process.env.DATABASE_URL = container.getConnectionUri() + '?sslmode=disable';
    // run migrations
    await runMigrations();
});
after(async () => {
    await container?.stop();
});
```

### Go

```go
import (
    "testing"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupDB(t *testing.T) string {
    t.Helper()
    ctx := context.Background()
    container, err := postgres.Run(ctx,
        "postgres:18-alpine",
        postgres.WithDatabase("registration"),
        postgres.WithUsername("registration"),
        postgres.WithPassword("test"),
    )
    if err != nil { t.Fatal(err) }
    t.Cleanup(func() { container.Terminate(ctx) })

    dsn, _ := container.ConnectionString(ctx, "sslmode=disable")
    runMigrations(t, dsn)
    return dsn
}
```

`t.Cleanup` is the Go equivalent of `@AfterEach` — it runs after the test finishes regardless of success or failure.

### Why "fresh container per test class" is the right unit

| Granularity | Container starts | Speed | Isolation |
|-------------|------------------|-------|-----------|
| Per test method | ×N (one per test) | slow | best |
| Per test class | ×M (one per class) | balanced | good — clean DB per logical test group |
| Per JVM / process | ×1 (whole suite) | fast | weak — tests can pollute each other |

We pick per-class for the integration / REST tests. Within a class, use `@BeforeEach` to clean rows or wrap in a transaction that rolls back. Combine the two strategies for the best speed/isolation trade-off.

## 2. Seed data — "this email already exists"

Some tests require the DB to be in a *specific* pre-populated state. The classic case for our app: `POST /register` with an email that's already in the database — the response should be 409.

**Bad approach:** call `POST /register` to insert the fixture, then call `POST /register` again to test the duplicate. The first call exercises code we haven't necessarily tested in this test method; if it breaks, the duplicate test fails for the wrong reason.

**Good approach:** insert the fixture directly into the DB, bypassing the application.

### Option A — SQL seed file

`spring-boot/src/test/resources/seed/email-exists.sql`:

```sql
INSERT INTO registration (id, email, full_name, phone, created_at)
VALUES ('00000000-0000-0000-0000-000000000001',
        'exists@example.com', 'Existing User', '08123456789', '2025-01-01 00:00:00+00');
```

```java
@Test
@Sql(scripts = "/seed/email-exists.sql")
void duplicate_email_returns_409() {
    given().port(port)
        .formParam("email", "exists@example.com")
        .formParam("fullName", "Different Person")
        .formParam("phone", "08555555555")
    .when()
        .post("/register")
    .then().statusCode(409);
}
```

`@Sql` runs the file before the test method. Spring Boot Test handles ordering.

### Option B — programmatic fixture

For Express/Go, simpler: just call the repo directly.

```javascript
import { test } from 'node:test';

test('duplicate email returns 409', async () => {
    await db.query(
        `INSERT INTO registration (id, email, full_name, phone, created_at) VALUES ($1, $2, $3, $4, NOW())`,
        ['00000000-0000-0000-0000-000000000001', 'exists@example.com', 'Existing', '08123456789']
    );
    const response = await request(app).post('/register').type('form').send({
        email: 'exists@example.com', fullName: 'Different', phone: '08555555555',
    });
    assert.equal(response.status, 409);
});
```

### Fixed UUIDs in seed data — why

`00000000-0000-0000-0000-000000000001` instead of `uuid.New()`:
- Tests can reference the same row in multiple assertions
- A failing test message saying "expected id=...0001" is debuggable; "expected id=<random>" is not
- Cleanup logic can target known IDs without a `WHERE *` scan

## 3. Parameterized CSV — validation cases as data

**Problem it solves:** validation rules have dozens of edge cases (email no @, name with digits, phone too short, …). Writing one test method per case is repetitive and hides the *table of cases* in the noise.

**Better:** put the cases in a CSV, write one test method that iterates over the rows.

`shared-fixtures/validation-cases.csv`:

```csv
case_name,email,full_name,phone,expected_status,expected_field
empty_email,,Alice Doe,08123456789,400,email
email_no_at,invalid,Alice Doe,08123456789,400,email
email_no_dot,a@b,Alice Doe,08123456789,400,email
email_with_space,a b@c.com,Alice Doe,08123456789,400,email
name_too_short,a@b.c,X,08123456789,400,fullName
name_with_digits,a@b.c,Alice 3,08123456789,400,fullName
phone_too_short,a@b.c,Alice Doe,12345,400,phone
phone_with_letters,a@b.c,Alice Doe,abc12345,400,phone
valid_minimum,a@b.cd,Al,1234567,302,
valid_full,alice@example.com,Alice Doe,08123456789,302,
```

The same CSV drives tests in all three stacks. When the contract changes, edit the CSV in one place.

### JUnit 5 — `@ParameterizedTest` + `@CsvFileSource`

```java
@ParameterizedTest(name = "{0}")
@CsvFileSource(resources = "/validation-cases.csv", numLinesToSkip = 1)
void validation_table(String caseName, String email, String fullName, String phone,
                      int expectedStatus, String expectedField) {
    var response = given().port(port)
        .formParam("email", email != null ? email : "")
        .formParam("fullName", fullName != null ? fullName : "")
        .formParam("phone", phone != null ? phone : "")
    .when()
        .post("/register");

    response.then().statusCode(expectedStatus);
    if (expectedStatus == 400 && expectedField != null && !expectedField.isEmpty()) {
        response.then().body(containsString(expectedField));
    }
}
```

### Node — `node:test` parameterized via `parse-csv` or `csv-parse`

```javascript
import { readFileSync } from 'node:fs';
import { parse } from 'csv-parse/sync';
import { test } from 'node:test';

const cases = parse(readFileSync('../shared-fixtures/validation-cases.csv'),
    { columns: true, skip_empty_lines: true });

for (const c of cases) {
    test(`validation: ${c.case_name}`, async () => {
        const resp = await request(app).post('/register').type('form').send({
            email: c.email, fullName: c.full_name, phone: c.phone,
        });
        assert.equal(resp.status, Number(c.expected_status));
    });
}
```

### Go — table-driven from CSV

```go
func TestValidationTable(t *testing.T) {
    f, _ := os.Open("../../shared-fixtures/validation-cases.csv")
    defer f.Close()
    r := csv.NewReader(f)
    r.Read() // skip header
    for {
        row, err := r.Read()
        if err == io.EOF { break }
        if err != nil { t.Fatal(err) }
        t.Run(row[0], func(t *testing.T) {
            resp, _ := http.PostForm(serverURL+"/register", url.Values{
                "email": {row[1]}, "fullName": {row[2]}, "phone": {row[3]},
            })
            wantStatus, _ := strconv.Atoi(row[4])
            if resp.StatusCode != wantStatus {
                t.Errorf("want %d, got %d", wantStatus, resp.StatusCode)
            }
        })
    }
}
```

### Where to put the CSV

A shared `shared-fixtures/` directory at repo root, alongside `db/`:

```
tech-stack-2026/
├── db/
├── shared-fixtures/
│   └── validation-cases.csv
├── spring-boot/   (refers to ../shared-fixtures/)
├── expressjs/     (refers to ../shared-fixtures/)
└── golang/        (refers to ../shared-fixtures/)
```

Same fixture file across all stacks = the contract is data, not code. Editing the CSV automatically updates all three test suites.

## Putting it together — what a complete test class looks like

```
@SpringBootTest(WebEnvironment.RANDOM_PORT)
@Testcontainers
class RegistrationIT extends AbstractIT {

    @Test                                  // unit-ish: validation alone
    @CsvFileSource(...)
    void validation_table(...)

    @Test                                  // integration: repo + DB
    void insert_then_findAll_returns_in_DESC_order()

    @Test                                  // REST endpoint, happy path
    void post_register_valid_returns_302()

    @Test                                  // REST endpoint, seeded duplicate
    @Sql("/seed/email-exists.sql")
    void post_register_duplicate_returns_409()
}
```

One class, four test types, all three criteria (repeatable, independent, portable) satisfied by the combination of:
- Testcontainers ⇒ portable + isolated
- `@Sql` seed file or `@BeforeEach` cleanup ⇒ repeatable
- Each test method owns its own state ⇒ independent
