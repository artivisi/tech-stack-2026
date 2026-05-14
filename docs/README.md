# Course material

Teaching module for STMIK class. The `tech-stack-2026` repo is the system under test (SUT) — three implementations of the same registration app (Spring Boot 4 / Express 5 / Go 1.26) that students can compare side-by-side.

This module covers two things:

1. **Web app architecture patterns** — rehash, mapped to the three stacks
2. **Test automation + CI/CD from scratch** with an 80% coverage gate

The test suite plan: Playwright functional tests, REST endpoint functional tests, integration tests, performance/capacity tests. All must be repeatable, portable, and independent.

## Outline

### Part A — Concepts

| # | Topic | Source notes |
|---|-------|--------------|
| 01 | [Web app architecture patterns](01-webapp-patterns.md) | `images/01-webapp-patterns.png` |
| 02 | [Automated test fundamentals](02-test-fundamentals.md) — mental model, lifecycle, good-test criteria | `images/02-test-mental-model.png`, `03-test-lifecycle.png`, `04-good-test.png` |
| 03 | [Testing tools landscape](03-testing-tools.md) | `images/05-testing-tools.png` |

### Part B — Test setup

| # | Topic |
|---|-------|
| 04 | [Test types in this project](04-test-types.md) — unit, integration, REST, Playwright UI, k6 perf |
| 05 | [Test infrastructure](05-test-infrastructure.md) — Testcontainers, seed data, CSV-parameterized validation |
| 06 | [CI/CD with 80% coverage gate](06-cicd-coverage.md) — GitHub Actions per stack |

### Part C — Per-stack walkthroughs (request → response)

Each walkthrough has a class diagram, sequence diagram, step-by-step trace with `file:line` references, and a "what's under the hood" section.

| # | Stack |
|---|-------|
| 07 | [Spring Boot — `DispatcherServlet` → `JdbcTemplate`](07-walkthrough-spring-boot.md) |
| 08 | [Express — middleware chain → `pg.Pool`](08-walkthrough-expressjs.md) |
| 09 | [Go — `ServeMux` → `database/sql`](09-walkthrough-golang.md) |

## How to read

Sections 01–03 are conceptual (the screenshots are the spine). Sections 04–06 are concrete: every claim references actual files in this repo. Test code shown in those sections is *example code for teaching*; it is **not** committed in the scaffold (test infrastructure is its own teaching unit — that's the whole point of this module).

## Pre-requisites

Students should already have:

- Run `./scripts/reset-db.sh`, started one of the three stacks, and seen `./scripts/verify-http.sh` pass 23/23.
- Read each stack's `README.md` and skimmed the source.
- Understood `db/schema.sql` and `db/validation.md` as the cross-stack contract.

If those steps don't make sense yet, finish the [root README](../README.md) first.
