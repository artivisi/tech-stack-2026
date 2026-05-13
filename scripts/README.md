# scripts/

Stack-agnostic verification scripts. These are the same ad-hoc curl/psql checks used to verify the Express implementation — saved here so they can be re-run against the Spring Boot and Go implementations to confirm cross-stack parity.

**These are not automated tests.** They share state across runs, have no per-case isolation, and don't replace JUnit / `node:test` / Go `testing`. They're acceptance scripts: same endpoints, same validation contract, same expected status codes.

## Files

| Script | Layer | What it does |
|--------|-------|--------------|
| `reset-db.sh` | infra | `docker compose down -v && up`; waits for healthy |
| `verify-http.sh` | HTTP | 20 scenarios (smoke + valid + invalid + duplicate) |
| `verify-db-constraints.sh` | DB | 6 direct-SQL inserts to exercise each CHECK constraint |

## Typical run

Reset the database between stacks so duplicate-email tests don't interfere across runs:

```bash
./scripts/reset-db.sh

# Pick one stack:
cd expressjs   && npm run migrate:up && npm run css:build && npm run dev &
# cd spring-boot && ./mvnw spring-boot:run &     (planned)
# cd golang      && go run ./cmd/server &        (planned)

# wait for server to be ready, then:
./scripts/verify-http.sh
./scripts/verify-db-constraints.sh
```

## Cross-stack contract

The scripts encode the contract every stack must satisfy:

| Endpoint | Method | Status on success | Status on invalid | Status on duplicate |
|----------|--------|-------------------|-------------------|---------------------|
| `/` | GET | 200 (HTML) | — | — |
| `/health` | GET | 200 (JSON `{status:"ok"}`) | — | — |
| `/css/app.css` | GET | 200 (text/css) | — | — |
| `/register` | POST | **302** → `/registrations` | **400** | **409** |
| `/registrations` | GET | 200 (HTML) | — | — |

If a stack returns 200-with-error-page instead of 302/400/409, the scripts will fail and the stack must be adjusted to match the contract — not the other way around.

## What the scripts don't assert

- Exact error message text (varies per stack default; defined in `db/validation.md` and enforced via UI review, not automated grep)
- HTML structure of returned pages
- CSS rendering / visual correctness
- Database row contents after `POST /register` (verify manually if needed)

## Configuration

| Env var | Default | Used by |
|---------|---------|---------|
| `BASE_URL` | `http://localhost:8080` | `verify-http.sh` |
| `CONTAINER` | `tech-stack-2026-db` | `verify-db-constraints.sh` |
