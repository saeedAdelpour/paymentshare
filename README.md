# Payment Share API

A Go/Chi server with a GORM PostgreSQL connection pool, explicit Goose migrations,
OpenAPI-generated transport, and local Swagger UI. Authentication and QOR5
administration follow in foundation 02 and 03; business APIs follow those tickets.

Follow the [engineering standards](docs/engineering-standards.md),
[glossary](GLOSSARY.md), and [architecture decisions](docs/adr/) when contributing.

## Start development

Install Go **1.27.1**, Make, and optionally Docker with Compose. From a fresh
checkout, use the supplied local PostgreSQL setup:

```sh
cp .env.example .env
make db-up
make migrate
make dev
```

Already have PostgreSQL? Skip `make db-up` and set your own `DATABASE_URL` in
`.env` before running `make migrate` and `make dev`. The database and login must
already exist. **Migrations change the configured database**; verify the URL
before invoking them. The API never runs migrations or GORM AutoMigrate at startup.

`make dev` loads `.env`, enables Swagger UI by default (an explicit
`SWAGGER_ENABLED=false` wins), and starts on **127.0.0.1:8082**. Stop with Ctrl+C.
Go downloads missing dependencies automatically. After setup, starting development
only requires `make dev`. Run `make migrate` explicitly when the release changes.

Open:

- [Swagger UI](http://127.0.0.1:8082/docs/) to explore and execute the health API.
- [OpenAPI contract](http://127.0.0.1:8082/openapi.yaml).
- [Process health](http://127.0.0.1:8082/healthz).
- [Database readiness](http://127.0.0.1:8082/readyz).

The root `/` has no homepage and returns 404. Swagger uses local, pinned assets;
no CDN or public validator is contacted. OIDC authorization arrives in foundation 02.

`.env` uses shell syntax, is ignored by Git, and is loaded by `make dev` and
`make migrate`. These commands also accept an already-exported configuration if
`.env` is absent. `make run` uses the exported environment without loading `.env`
and leaves documentation disabled unless explicitly enabled.

| Variable | Meaning |
| --- | --- |
| `DATABASE_URL` | Required PostgreSQL URL or pgx keyword/value connection string |
| `HTTP_ADDR` | Defaults to `127.0.0.1:8082`; `:8082` permits container access |
| `SWAGGER_ENABLED` | Boolean; defaults to false in the application; `make dev` defaults to true |
| `TEST_DATABASE_URL` | Required for integration checks; role must be able to create schemas |

URL-encode special characters in URL usernames/passwords. Keep secrets in `.env`
or your deployment secret store. The Compose credentials and `sslmode=disable`
are for local development only; configure verified TLS for remote databases.
Compose binds PostgreSQL to loopback on port **55432**. It creates development
and test databases when initializing its named volume. `make db-down` preserves
that volume; changing initialization SQL does not modify an existing volume.

The API checks connectivity and exact schema compatibility before listening.
Startup failures exit nonzero with credential-safe messages; underlying causes
remain available programmatically. A PostgreSQL connection check has a five-second
timeout. Health responses disable caching and return JSON:

| Method and path | Behavior |
| --- | --- |
| `GET /healthz` | 200 `{"status":"ok"}`; independent of PostgreSQL |
| `GET /readyz` | 200 `{"status":"ok"}` when PostgreSQL and schema are ready; otherwise 503 `{"status":"unavailable"}` |

Readiness has a two-second timeout. SIGINT/SIGTERM drains in-flight requests for
up to ten seconds before closing connections. The pool permits ten open/five idle
connections, with a 30-minute lifetime and five-minute idle timeout. These are
initial limits, not measured capacity claims.

## Migrations and deployment

```sh
make migrate           # load local .env and migrate to this release's latest version
make migrate TO=0      # local reversal to version zero
make migrate TO=1      # local upgrade to the baseline
# Deployment: use exported DATABASE_URL rather than a local .env:
go run ./cmd/migrate -to 1
go run ./cmd/api
```

Goose **v3.29.0** executes embedded, reviewed SQL migrations transactionally and
serializes execution with a PostgreSQL session advisory lock. Commands have a
60-second work deadline plus at most five seconds for detached unlock cleanup.
Repeating a target is harmless. Unsupported targets and
schemas newer than the migration binary are rejected. Version zero is deliberately
incompatible with the current server, which expects **version 1**.

The first migration establishes Goose tracking only; there are no domain tables
yet. Its Down preserves the tracking table and all existing application data.
Round-trip tests retain a Unicode/numeric fixture across 1 → 0 → 1. This verifies
the baseline, not the reversibility of future domain migrations.

See [deployment and recovery](docs/deployment.md) for migration-before-start,
API/admin arrangement, previous-release compatibility, and production recovery.

## Contract and quality checks

[`api/openapi.yaml`](api/openapi.yaml) is the source of truth. Generated transport
types and Chi interfaces stay in `internal/httpapi`. Document application API
operations there as they are added; QOR5 internal browser routes and provider
endpoints are outside this contract. Documentation routes serve the contract and
UI but are not application API operations.

```sh
make deps                 # download application and separate tool-module dependencies
make tools                # build pinned tools into ignored bin/
make generate             # regenerate the Chi interface and Go transport types
make check-contract       # validate OpenAPI; SPEC=path validates a different document
make check-generated      # compare generated code without changing the checkout
make test-contract-tools  # prove invalid contracts and stale code are rejected
make fmt                  # apply pinned gofumpt formatting
make check                # formatting, contract, generation, build, lint, security
make test                 # race-enabled unit/HTTP tests
# Export TEST_DATABASE_URL from .env before database checks:
set -a; . ./.env; set +a
make test-integration     # real PostgreSQL command and HTTP behavior tests
make ci                   # same checks and full race suite used by CI
```

Tools are pinned with checksums in the separate [`tools/go.mod`](tools/go.mod)
and `tools/go.sum`: oapi-codegen v2.8.2, gofumpt v0.12.0, golangci-lint v2.15.0,
and govulncheck v1.8.0. The linter initially enables
**govet, staticcheck, errcheck, ineffassign, unused** for application and integration
code. `make check` fails on formatting, contract, stale output, analysis, or
vulnerability errors; unavailable vulnerability data fails visibly. It requires
network access for first-time tool downloads and vulnerability data.

The HTTP harness validates requests and responses against OpenAPI for both healthy
and unavailable readiness, alongside explicit behavior assertions. Database tests
execute the real migration binary, create random schemas, and remove only their
own schemas afterward. They check target-version Up/Down/Up, retained fixtures,
unmigrated/older/newer schema rejection, startup without automatic migrations,
readiness after a schema change, bounded cancellation/unlock cleanup through a
real PostgreSQL TCP proxy,
and graceful shutdown. The fault-injection test requires a local TCP database
that permits non-TLS connections (as provided by Compose/CI). Never point tests at
production. Missing `TEST_DATABASE_URL` is a failure, not a silent skip.

The [CI job](.github/workflows/quality.yml) supplies PostgreSQL 17.11 and runs
`make ci` on Go 1.27.1. Local Docker setup uses the same PostgreSQL version.
Browser assets and licensing are recorded in [third-party notices](THIRD_PARTY.md).
Verification evidence and remaining release gaps are in
[foundation verification](docs/foundation-verification.md).
