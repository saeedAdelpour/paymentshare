# 01: Create isolated workspaces

Status: ready-for-agent

**What to build:** Reuse foundation authentication to create, list, and rename isolated workspaces with a chosen currency and one owner.

**Blocked by:** [Foundation 03: Administration and permissions](../../payment-share-foundation/issues/03-administration-permissions.md).

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 2, 3, 4, 5, 38, 45, 46, 47, 56, 58, 59, 60.

## Implementation constraints

Authentication, the administration shell, migration management, and quality tooling are now prerequisites delivered by the foundation phase. Reuse and verify that foundation here; do not start workspace feature work before its acceptance checks pass.

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, QOR5, OIDC, and GORM/PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. Reuse the completed [project foundation](../../payment-share-foundation/spec.md). Register any application tables introduced by this slice in administration and route administrative mutations through the same application use cases; preserve financial invariants and history. Follow the current parent spec and the approved foundation stack and quality checks. Verify administrative behavior through browser/session or direct administrative HTTP requests with real PostgreSQL.

## Acceptance criteria

- [ ] Extend or audit the OpenAPI contract for this slice, documenting introduced/changed application API operations, inputs, success responses, relevant errors, authentication, and required headers. Verify HTTP requests/responses against the contract and pass contract validation and generated-Go freshness checks; keep Swagger UI usable with the foundation OIDC authorization flow.
- [ ] Reuse the completed foundation's running Go/Chi application, GORM/PostgreSQL persistence, OIDC identity, administration, Goose migrations, and pinned quality checks. Extend the existing harness rather than rebuilding authentication or project setup.
- [ ] Create, list, and rename workspaces using current authenticated identity. Unauthenticated/deactivated users and users outside a workspace cannot access its records.
- [ ] Persist workspace ownership and currency; a user can own multiple workspaces. Commit workspace creation and owner membership together, preserving exactly one owner.
- [ ] Only the current owner can rename a workspace through the member API. Check current ownership/membership with the mutation and test isolation, rejected writes, and persistence across restart.
- [ ] Add versioned migrations for workspace data and verify their local Up → Down → Up behavior with representative fixtures. Keep the foundation's release migration and rollback contract.
- [ ] Register workspace records in administration. Exercise permitted create/read/rename actions and denied staff actions; retain single-owner and currency rules and enforce application deactivation. Workspace deletion remains outside MVP scope.
