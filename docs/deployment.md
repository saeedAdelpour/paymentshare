# Foundation deployment and recovery

## Process arrangement

Deploy the `cmd/api` executable as one service behind an HTTPS reverse proxy.
The current service exposes only operational health endpoints and optional API
documentation. Foundation 03 will mount QOR5 administration in the same process,
sharing configuration and application use cases. There is no administration route
or separate admin service in this release. Do not assume staff authentication
exists before foundations 02 and 03 are implemented and verified.

Use a private PostgreSQL endpoint and a database role appropriate for runtime
queries. Use a separate migration role with schema privileges when deploying.
Set DATABASE_URL through a secret store, HTTP_ADDR for the listener, and leave
SWAGGER_ENABLED unset/false unless intentionally exposing docs. The development
helper enables docs; deploy the binary directly instead of using `make dev`.

## Release order

1. Review migration SQL and assess data size, locking, compatibility, and recovery.
2. Take a verified backup/PITR recovery point as appropriate for the change.
3. Run one release job using the new release's `cmd/migrate` binary against the
   intended database. It defaults to the release's latest version; `-to N` chooses
   an explicit local target. It serializes with other migration jobs using a
   PostgreSQL advisory lock, uses one transaction per SQL migration, and stops
   after a 60-second deadline. The database must support session advisory locks;
   use a direct connection rather than transaction-mode pooling for migrations.
4. Only after that job succeeds, start the API binary with runtime credentials.
   Startup reads Goose's version without creating its table. An absent,
   older, or newer version rejects startup. Do not run AutoMigrate or migrations
   from API replicas.
5. Route traffic after `/readyz` returns 200. Liveness `/healthz` deliberately
   ignores database outages; readiness checks connectivity and schema within
   two seconds and returns 503 when either is unavailable.

The baseline creates only Goose version tracking. No application tables, financial
records, or identity data are introduced here. Restrict schema-write privileges
to release tooling; version tracking alone does not detect manual schema drift.

## Previous-release compatibility

The current API requires exactly schema version 1. Version zero and unknown future
versions are rejected. The earlier connectivity-only server did not require a
schema, so adding this no-op baseline does not change its data access. This is the
only compatibility claim for these releases. Rolling upgrades and uninterrupted
mixed-version operation have not been demonstrated.

For future releases, review a compatibility matrix before deployment. Prefer
additive expand/contract changes that let the previous binary continue reading
and writing until replacement is complete. The current exact-version gate must
be deliberately revised and tested before claiming support for a version range.
A binary from before a new migration may intentionally refuse startup.

## Local reversal versus production recovery

For local switching, stop the server, run `make migrate TO=N`, then start a binary
compatible with N. Current supported targets are 0 and 1. The baseline's Down
preserves unrelated data and Goose tracking; the current server cannot run at 0.

Never infer production recovery from a Down file. Recreating a column/table
cannot restore destroyed information. Each future destructive change needs an
explicit plan: backup/PITR restore, replay of retained history, or a reviewed
forward repair. Reconcile all writes after the recovery point before returning
traffic. Assess downtime, external effects, and access to secrets separately.

On a failed migration, leave the old service in place if compatible; otherwise
keep traffic stopped while inspecting database-side diagnostics and the reviewed
SQL. Each baseline migration is transactional, so an interrupted SQL migration
rolls back as a unit. Do not edit an applied migration or manually mark it applied
to bypass an error. Application errors avoid raw driver text to prevent secret
leakage; use controlled operator diagnostics when additional detail is required.

Backup restoration, failover, production rollout, capacity, and recovery objectives
have not been exercised by this foundation. Define and verify them for the target
hosting environment before release.
