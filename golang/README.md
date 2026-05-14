# golang

Registration app — Go 1.26, stdlib `net/http`, `html/template`, `database/sql` with `pgx/v5/stdlib` driver, `golang-migrate/migrate` with pgx/v5 driver, `go-playground/validator/v10`, Tailwind CSS v4.

Color theme: **cyan**.

## Prerequisites

- Go 1.26+
- Node.js 22+ (for Tailwind build only)
- PostgreSQL running (`docker compose up -d db` from repo root)
- Root `.env` populated

## Install Tailwind

```
npm install
```

## Build CSS

```
npm run css:build
```

For dev: `npm run css:watch` in a separate terminal.

## Run

```
./migrate.sh up   # apply migrations
./run.sh          # start server
```

Both scripts source `../.env` then invoke `go run`. Server listens on `http://localhost:${HTTP_PORT}`.

## Layout

```
golang/
├── go.mod, go.sum
├── package.json                              # Tailwind CLI only
├── input.css                                 # Tailwind entry (@source -> templates)
├── run.sh, migrate.sh                        # source ../.env + go run
├── cmd/
│   ├── server/main.go                        # binary entry, graceful shutdown
│   └── migrate/main.go                       # golang-migrate via pgx WithInstance
├── internal/
│   ├── config/                               # env helper, db connect
│   ├── domain/registration.go                # struct
│   ├── repository/registration.go            # database/sql ops, 23505 detection
│   └── web/
│       ├── handler.go                        # GET /, POST /register, GET /registrations, /health
│       ├── form.go                           # RegistrationForm + validator + error mapping
│       ├── middleware.go                     # RequestLogger
│       ├── templates.go                      # embed.FS + html/template loader
│       └── templates/
│           ├── layout.html                   # {{block "content"}}
│           ├── form.html
│           └── list.html
├── migrations/
│   ├── 000001_create_registration.up.sql
│   └── 000001_create_registration.down.sql
└── static/css/app.css                        # Tailwind output (gitignored)
```

## Env vars

Read from repo-root `.env` via `run.sh` / `migrate.sh`. Missing vars cause `log.Fatal` — no defaults.

| Var | Used by |
|-----|---------|
| `DATABASE_URL` | `pgx` (server), `golang-migrate` (migrate) |
| `HTTP_PORT` | `cmd/server` |

## Notes

- Both server (runtime queries) and migrate (`golang-migrate`) use the same `pgx/v5` driver via `database/sql`. No lib/pq.
- `golang-migrate` is invoked through `pgxmigrate.WithInstance(db, …)` so the migration tool reuses the same connection setup as the application code rather than re-parsing the URL.
- `html/template` layout pattern: `layout.html` defines `{{block "content" .}}{{end}}`, page templates redefine `{{define "content"}}…{{end}}`. Each page is parsed together with the layout into its own `*template.Template`.
- Validator uses `RegisterTagNameFunc` to expose the `form` struct-tag name (e.g. `email`, `fullName`) so error keys match the HTML field names.
- The DB column is `full_name` (SQL convention). `repository.Registration` reads it into `domain.Registration.FullName` (Go convention) via the column-by-column `Scan` call.
- No automated tests included — test infrastructure is taught as a separate unit.
