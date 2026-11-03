# Payment Share API

The first foundation step is a Go/Chi HTTP server with a GORM PostgreSQL
connection pool. It checks database connectivity before listening, handles
SIGINT/SIGTERM with up to 10 seconds to drain requests, and closes the pool.
It does not yet implement authentication, administration, or business APIs.

Follow the [engineering standards](docs/engineering-standards.md) and
[architecture decisions](docs/adr/) when contributing.

## Run locally

Install Go 1.27 or newer. Create a PostgreSQL database and login yourself;
the server does not create databases, tables, or run migrations.

```sh
cp .env.example .env
# Edit .env with your database connection details, then export it:
set -a
. ./.env
set +a
make deps
make run
```

`.env` is ignored by Git and is not loaded automatically. Keep credentials there
or in your deployment's secret store. URL-encode special characters in URL
usernames/passwords. The pgx keyword/value connection format is also supported.

| Variable | Behavior |
| --- | --- |
| `DATABASE_URL` | Required PostgreSQL connection string |
| `HTTP_ADDR` | Defaults to `127.0.0.1:8082`; use `:8082` for container access |
| `TEST_DATABASE_URL` | Required by `make test-integration`; use a dedicated test database |

The example disables TLS only for local development. Configure TLS verification
(e.g. `sslmode=verify-full` and trusted certificates) for remote deployments.
A failed initial connection exits nonzero without printing the connection string
or driver errors, which can contain credentials. Check database-side logs for
connection failure details. Initial connection checks have a five-second timeout.
The pool allows ten open/five idle connections, with a 30-minute lifetime and
five-minute idle lifetime. These are initial limits, not measured capacity claims.

```sh
curl -i http://127.0.0.1:8082/healthz
curl -i http://127.0.0.1:8082/readyz
```

`GET /healthz` returns `200 {"status":"ok"}` without querying the database.
`GET /readyz` checks PostgreSQL with a two-second timeout and returns `200` with
`{"status":"ok"}` or `503` with `{"status":"unavailable"}`. Both are public,
return JSON, and disable caching. Readiness checks connectivity only, not schema
compatibility. No schema is required in this initial step.

## Contract and checks

[`api/openapi.yaml`](api/openapi.yaml) is the health API source of truth. The
Chi server interface and response types are generated inside the HTTP adapter
using pinned oapi-codegen v2.5.1. Commit generated code alongside the contract.

```sh
make generate        # regenerate the HTTP interface
make check-generated # fail if generated output is stale
make check           # formatting, generated output, build, vet
make test            # race-enabled unit/HTTP tests
make test-integration # real PostgreSQL health checks; requires TEST_DATABASE_URL
```

Integration tests never migrate, truncate, or change database records. They check
healthy readiness, unavailable readiness after closing the pool, and independent
process liveness. An explicitly invoked integration run fails if configuration is
missing. The normal unit suite does not claim to verify database connectivity.

## Remaining foundation work

Ticket [01](.scratch/payment-share-foundation/issues/01-runnable-project-quality-checks.md)
remains incomplete: Swagger UI, full contract validation and HTTP contract tests,
Goose migrations and startup schema compatibility, migration reversal/recovery
checks, pinned formatting/lint/security tools, and CI are subsequent steps.
Identity and QOR5 administration follow in foundation 02 and 03. There is no
schema change to roll back in this step. Production deployment and database
recovery procedures remain unverified; no production-readiness claim is made.
