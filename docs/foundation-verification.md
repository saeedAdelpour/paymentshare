# Foundation 01 verification

Verified locally on 2026-11-03 using Linux amd64, Go
`go1.27.1-X:nodwarf5`, and PostgreSQL **17.11-alpine** from the same digest pinned
in Compose and CI. Verification used a separate Compose project and test database;
the user's configured application database was not migrated or otherwise changed.

## Evidence

- The documented Compose service started from a new named volume, initialized
  separate development/test databases, and passed its PostgreSQL health check.
- Focused real PostgreSQL tests passed for target-version 1 → 0 → 1, repeated
  targets (including zero), retained Unicode/numeric fixtures, unsupported target
  rejection, unmigrated/older/newer schema rejection, and read-only startup.
- Concurrent initial migration commands passed and recorded the baseline once.
- Real HTTP tests passed for healthy/unavailable readiness, schema reversal while
  running, independent liveness, and graceful shutdown. Test requests and responses
  are validated against the source OpenAPI document.
- HTTP documentation tests passed for default-disabled routes and explicitly
  enabled HTML, local JS/CSS assets, and the source contract.
- A Chromium browser smoke check loaded Swagger UI and used Try it out / Execute
  for GET /healthz, observing 200 with `{"status":"ok"}` and no JavaScript errors.
- The actual validation/generation commands rejected an invalid schema and stale
  generated Go in a temporary checkout. Valid source/generated output passed.
- Formatting, build, golangci-lint (govet/staticcheck/errcheck/ineffassign/unused),
  and pinned browser-asset checksum checks passed during implementation.
- govulncheck reported two reachable kin-openapi issues at v0.133.0. Upgrading
  to v0.144.0 removed the findings; the scan then reported no vulnerabilities.
  oapi-codegen v2.5.1 was incompatible with the patched parser, so the tool module
  was upgraded to v2.8.2 and generated code was regenerated.
- Full `make ci` passed: browser-asset integrity, formatting, OpenAPI validation,
  generation comparison, negative contract-tooling checks, build, lint (zero
  issues), vulnerability scan (no vulnerabilities), and the complete race-enabled
  suite with real PostgreSQL integration tests.

The first database checks used a cached PostgreSQL 17.6 image while Docker Hub
was unavailable. They were repeated successfully on the pinned 17.11 image after
its download succeeded. Only the final 17.11 result establishes parity with the
supplied local/CI database version.

## Scope and release gaps

This is a runnable foundation. Its SQL baseline creates Goose version tracking
only; there are no domain tables or business mutations to exercise yet. The
fixture round-trip proves baseline retention, not future migration reversibility.
Version tracking does not detect out-of-band manual schema changes.

The GitHub Actions workflow is supplied but has not been executed on GitHub in
this session. The shared `make ci` command provides the local verification path;
hosted runner behavior remains to be confirmed after pushing.

Authentication, OIDC Swagger authorization, staff permissions, and QOR5
administration belong to foundations 02/03. Capacity/load, backup restoration,
PITR, failover, target-environment TLS, real deployment/rolling compatibility,
and recovery objectives remain unverified. There is no production performance
or disaster-recovery guarantee from these checks. See [deployment](deployment.md).
