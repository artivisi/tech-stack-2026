# spring-boot

Registration app — Spring Boot 4.0.6, Java 25, Spring MVC, Thymeleaf (with `thymeleaf-layout-dialect`), `JdbcTemplate`, Flyway, Bean Validation, Tailwind CSS v4.

Color theme: **blue**.

## Prerequisites

- JDK 25+
- Node.js 22+ (for Tailwind build only)
- PostgreSQL running (`docker compose up -d db` from repo root)
- Root `.env` populated (see root `.env.example`)

## Build CSS

```
npm install
npm run css:build
```

For dev hot-reload: `npm run css:watch` in a separate terminal.

## Run

```
./run.sh
```

`run.sh` sources `../.env` and invokes `./mvnw spring-boot:run`. Flyway runs the migration on startup.

Server: `http://localhost:${HTTP_PORT}` (8080 by `.env.example`).

## Layout

```
spring-boot/
├── pom.xml
├── mvnw, .mvn/                       # Maven wrapper
├── package.json                      # Tailwind CLI only
├── run.sh                            # sources ../.env + ./mvnw spring-boot:run
└── src/main/
    ├── css/input.css                 # Tailwind entry
    ├── java/com/artivisi/techstack/registration/
    │   ├── RegistrationSpringBootApplication.java
    │   ├── config/DataSourceConfig.java       # DATABASE_URL -> JDBC
    │   ├── domain/Registration.java
    │   ├── repository/RegistrationRepository.java   # JdbcTemplate
    │   └── web/
    │       ├── RegistrationController.java
    │       └── RegistrationForm.java                 # @NotBlank/@Size/@Pattern
    └── resources/
        ├── application.properties
        ├── db/migration/V1__create_registration.sql  # Flyway
        ├── static/css/app.css                        # Tailwind output (gitignored)
        └── templates/
            ├── layouts/main.html
            ├── form.html
            ├── list.html
            └── error.html
```

## Env vars

Read from repo-root `.env` via `run.sh`. Missing required vars cause startup to throw — no defaults.

| Var | Used by |
|-----|---------|
| `DATABASE_URL` | `DataSourceConfig` (parses URI; constructs JDBC URL) |
| `HTTP_PORT` | `application.properties` `server.port` |

## Notes

- No automated tests included — test infrastructure is taught as a separate unit. The Initializr-generated test class was removed.
- Java side stays in camelCase (`fullName`). The HTML form field name `full_name` is bound via a `@ModelAttribute` factory method in the controller that takes `@RequestParam("full_name") String fullName` and constructs the record. The `@PostMapping` handler then uses `@Valid @ModelAttribute(binding=false)` so Spring doesn't re-bind from request (which would null the field since records are immutable and Spring would look for an `fullName` request param that doesn't exist).
- DB-level CHECK constraints in `V1__create_registration.sql` mirror app-level Bean Validation rules. See `db/validation.md` at repo root.
