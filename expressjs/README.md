# expressjs

Registration app — Express 5.1.0, Handlebars (`express-handlebars` v7), PostgreSQL via `pg`, migrations via `node-pg-migrate`, Tailwind CSS v4.

Color theme: **emerald**. Stacks share the same layout structure; only the accent color and stack badge differ.

## Prerequisites

- Node.js 22+
- PostgreSQL running (see root `docker-compose.yml`)
- Root `.env` populated (see root `.env.example`)

## Install

```
npm install
```

## Run

From repo root:

```
docker compose up -d db
```

From `expressjs/`:

```
npm run migrate:up
npm run css:build
npm run dev
```

Then open `http://localhost:${HTTP_PORT}` (8080 by default in `.env.example`).

## Scripts

| Script | Purpose |
|--------|---------|
| `npm run dev` | Start server with `node --watch` |
| `npm start` | Start server |
| `npm run css:build` | Build Tailwind CSS once (minified) |
| `npm run css:watch` | Build Tailwind CSS in watch mode |
| `npm run migrate:up` | Apply pending migrations |
| `npm run migrate:down` | Roll back last migration |
| `npm run migrate:create -- <name>` | Create new SQL migration |

## Layout

```
expressjs/
├── package.json
├── migrations/                 # node-pg-migrate SQL files
├── public/                     # static assets (Tailwind output goes here)
└── src/
    ├── server.js               # entry — reads HTTP_PORT, starts listener
    ├── app.js                  # express setup, view engine, middleware
    ├── db.js                   # pg pool (reads DATABASE_URL)
    ├── env.js                  # required-env helper (throws if missing)
    ├── routes/registration.js  # GET /, POST /register, GET /registrations
    ├── styles/input.css        # Tailwind entry
    └── views/                  # Handlebars templates
```

## Env vars

All read from repo-root `.env` via Node's `--env-file`. Missing vars cause startup to throw — no defaults.

| Var | Used by |
|-----|---------|
| `DATABASE_URL` | `pg` pool, `node-pg-migrate` |
| `HTTP_PORT` | `server.js` |
