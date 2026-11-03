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

- [ ] Establish the application OpenAPI contract as the source of truth, with pinned oapi-codegen generating Go transport types and Chi server interfaces inside the HTTP adapter. Serve working Swagger UI and its contract locally; other deployments require explicit configuration and default to disabled. Document all application API operations implemented in this slice, excluding QOR5 internal browser routes and identity-provider endpoints.
- [ ] Provide matching local/CI checks rejecting invalid OpenAPI documents and stale generated Go code. Extend the HTTP test harness to verify requests and responses against the contract; document reproducible generation and check commands. Swagger authentication is completed in foundation 02.
- [ ] A fresh checkout can run the application and local PostgreSQL using documented, reproducible commands. A minimal HTTP smoke/readiness check demonstrates a running application and database connectivity.
- [ ] Integrate GORM with PostgreSQL and manage schema evolution through reviewed, versioned Goose SQL migrations. Application startup checks schema compatibility; migration execution is an explicit command/release step rather than runtime schema synchronization.
- [ ] Support local migration to a specified version in either direction. Run Up → Down → Up against real PostgreSQL with representative fixtures; verify data retention where promised and document destructive limitations.
- [ ] Document the initial API/admin deployment arrangement, migration-before-start sequence, and previous-release compatibility/recovery expectations without claiming recovery from Down files alone.
- [ ] Pin gofumpt, golangci-lint v2, and govulncheck and compatible dependency versions. Configure govet, staticcheck, errcheck, ineffassign, and unused initially; run race-enabled Go tests and real PostgreSQL integration checks.
- [ ] Provide matching local/CI quality commands and a reproducible CI job with its PostgreSQL dependency. Formatting and analysis fail on violations; security and test failures are visible rather than silently skipped.
- [ ] Run the applicable checks and migration smoke tests; update setup/check documentation to match the delivered behavior.
