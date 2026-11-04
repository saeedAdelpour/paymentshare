# 01: Runnable project and quality checks

Status: ready-for-agent

**What to build:** Run the local Go/Chi application with PostgreSQL, GORM persistence integration, reversible Goose migrations, and reproducible quality checks before authentication and administration work.

**Blocked by:** None (can start immediately).

**Parent:** [Payment Share project foundation](../spec.md)

**Spec user stories:** 1, 2, 3, 4, 14, 15, 16, 19, 20, 21, 22, 23, 24.

**MVP user stories:** 44, 49, 50, 53, 57 in the [MVP spec](../../payment-share-mvp/spec.md).

## Implementation constraints

Use the approved Go/Chi + PostgreSQL + GORM + QOR5 + Goose stack. Establish application composition and configuration without coupling the domain to infrastructure. Reconcile existing configuration and README claims with the actual checkout; preserve existing user changes. This is the approved initialization prerequisite to the existing MVP feature tickets.

## Acceptance criteria

- [x] Establish the application OpenAPI contract as the source of truth, with pinned oapi-codegen generating Go transport types and Chi server interfaces inside the HTTP adapter. Serve working Swagger UI and its contract locally; other deployments require explicit configuration and default to disabled. Document all application API operations implemented in this slice, excluding QOR5 internal browser routes and identity-provider endpoints.
- [x] Provide matching local/CI checks rejecting invalid OpenAPI documents and stale generated Go code. Extend the HTTP test harness to verify requests and responses against the contract; document reproducible generation and check commands. Swagger authentication is completed in foundation 02.
- [x] A fresh checkout can run the application and local PostgreSQL using documented, reproducible commands. A minimal HTTP smoke/readiness check demonstrates a running application and database connectivity.
- [x] Integrate GORM with PostgreSQL and manage schema evolution through reviewed, versioned Goose SQL migrations. Application startup checks schema compatibility; migration execution is an explicit command/release step rather than runtime schema synchronization.
- [x] Support local migration to a specified version in either direction. Run Up → Down → Up against real PostgreSQL with representative fixtures; verify data retention where promised and document destructive limitations.
- [x] Document the initial API/admin deployment arrangement, migration-before-start sequence, and previous-release compatibility/recovery expectations without claiming recovery from Down files alone.
- [x] Pin gofumpt, golangci-lint v2, and govulncheck and compatible dependency versions. Configure govet, staticcheck, errcheck, ineffassign, and unused initially; run race-enabled Go tests and real PostgreSQL integration checks.
- [x] Provide matching local/CI quality commands and a reproducible CI job with its PostgreSQL dependency. Formatting and analysis fail on violations; security and test failures are visible rather than silently skipped.
- [x] Run the applicable checks and migration smoke tests; update setup/check documentation to match the delivered behavior.

## Comments

2026-11-03: Implementing incrementally at the user's request. The first increment
adds the Go/Chi server, GORM PostgreSQL pool, startup connectivity check, health
contract/generated transport, graceful shutdown, and matching local instructions.
The user provisions PostgreSQL; this increment performs no schema mutations.
The full ticket remains open, including schema compatibility, Goose, Swagger UI,
contract validation, pinned quality tools, and CI. Live database verification is
pending a configured database; integration checks require TEST_DATABASE_URL.


2026-11-03: Completed the remaining foundation slice. Added local pinned Swagger
assets with explicit enablement, source-contract validation and HTTP exchange
validation, pinned tooling in a separate module, an embedded Goose baseline and
explicit target-version migration command, read-only startup/readiness schema
checks, PostgreSQL Compose/CI setup, deployment/recovery documentation, and real
migration/HTTP tests. The user-configured database was not changed.

## Answer

All acceptance criteria are implemented. `make ci` passed locally on PostgreSQL
17.11 using the same pinned image as Compose/CI: formatting, contract/generation
checks (including deliberate invalid/stale fixtures), build, lint, vulnerability
scan, and full race-enabled tests. Real migration tests cover Up/Down/Up, repeated
and concurrent targets, retained fixtures, and schema compatibility failures.
A Chromium smoke check rendered Swagger UI and executed GET /healthz successfully.

See [verification evidence](../../../docs/foundation-verification.md) and
[deployment/recovery](../../../docs/deployment.md). Hosted GitHub Actions has not
been run in this session; authentication/admin, deployment capacity, backup restore,
and production recovery verification remain later work. The baseline tracks
schema version only and makes no claims about future domain migration reversibility.


Final verification after review: `make ci` passed again against the isolated
PostgreSQL 17.11 database, including the complete race-enabled suite. Standards
and Spec review findings were resolved: undeclared HTTP statuses are rejected,
and detached migration unlock cleanup is bounded and fault-injection tested.
