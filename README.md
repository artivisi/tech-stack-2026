# tech-stack-2026

Single-table registration app implemented in three backend stacks for comparison. Same schema, same UI, same behavior — different runtimes.

## Stacks

| Stack | Language | Runtime | Web Framework |
|-------|----------|---------|---------------|
| `spring-boot/` | Java 21 | Spring Boot 4 | Spring MVC |
| `expressjs/`   | Node.js | Express 5.1.0 | Express |
| `golang/`      | Go 1.26 | net/http + chi/gin | stdlib templates |

Each implementation is self-contained and runs independently against the shared PostgreSQL instance.

## Frontend

- Server-rendered HTML, no SPA
- Tailwind CSS for styling (compiled, no CDN)
- Progressive enhancement only — forms work without JavaScript

## Database

PostgreSQL 18 Alpine, provisioned via `docker compose`.

```bash
docker compose up -d db
```

All three backends share the same schema and connect to the same database instance on `localhost:5432`.

## Domain

Registration form with a single table `registration`:

- `id` — UUID, primary key
- `email` — unique, validated
- `full_name`
- `phone`
- `created_at` — timestamp

Endpoints (identical across stacks):

| Method | Path             | Purpose                      |
|--------|------------------|------------------------------|
| GET    | `/`              | Show registration form       |
| POST   | `/register`      | Submit registration          |
| GET    | `/registrations` | List submitted registrations |

## Repository Layout

```
tech-stack-2026/
├── docker-compose.yml      # PostgreSQL 18 alpine
├── spring-boot/            # Spring Boot 4 implementation
├── expressjs/              # Express 5.1.0 implementation
├── golang/                 # Go 1.26 implementation
└── README.md
```

## Running

Start the database, then start any backend. Each backend's own README documents its build and run steps.

```bash
docker compose up -d db

# pick one:
cd spring-boot && ./mvnw spring-boot:run
cd expressjs   && npm start
cd golang      && go run ./cmd/server
```

All three serve on `http://localhost:8080` by default — run one at a time, or override the port per stack.

## Purpose

Side-by-side reference for evaluating the 2026 baseline of each stack: project layout, dependency footprint, build times, runtime characteristics, and ergonomics for the same trivial CRUD surface.

## License

See `LICENSE`.
