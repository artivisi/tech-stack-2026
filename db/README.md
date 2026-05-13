# db

Canonical schema for the registration app.

`schema.sql` is the source of truth. Each stack's migration tool must produce this exact structure via its own format:

| Stack | Tool | Location |
|-------|------|----------|
| Spring Boot | Flyway | `spring-boot/src/main/resources/db/migration/V1__create_registration.sql` |
| Express | `node-pg-migrate` | `expressjs/migrations/<ts>_create-registration.sql` |
| Go | `golang-migrate/migrate` | `golang/migrations/000001_create_registration.up.sql` |

The migration bookkeeping tables differ (`flyway_schema_history`, `pgmigrations`, `schema_migrations`), but the `registration` table itself is identical across stacks.

Only one migrator should run against the same database at a time.
