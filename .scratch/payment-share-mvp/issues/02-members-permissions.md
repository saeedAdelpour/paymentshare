# 02: Add members and delegate permissions

Status: ready-for-agent

**What to build:** Add registered members, delegate independent permissions, and transfer ownership.

**Blocked by:** 01: Create isolated workspaces.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 6, 7, 8, 9, 38, 45, 46, 47, 56, 58, 59, 60.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, QOR5, OIDC, and GORM/PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. Reuse the completed [project foundation](../../payment-share-foundation/spec.md). Register any application tables introduced by this slice in administration and route administrative mutations through the same application use cases; preserve financial invariants and history. Follow the current parent spec and the approved foundation stack and quality checks. Verify administrative behavior through browser/session or direct administrative HTTP requests with real PostgreSQL.

## Acceptance criteria

- [ ] Extend or audit the OpenAPI contract for this slice, documenting introduced/changed application API operations, inputs, success responses, relevant errors, authentication, and required headers. Verify HTTP requests/responses against the contract and pass contract validation and generated-Go freshness checks; keep Swagger UI usable with the foundation OIDC authorization flow.
- [ ] Owner or a member with add permission adds an existing registered user immediately by exact verified email, without exposing a searchable directory.
- [ ] Only owners grant or revoke independent add-member, remove-member, and edit-others permissions. An add permission does not imply other powers.
- [ ] Active members see their workspace; outsiders remain excluded, including attempts using another workspace's identifiers.
- [ ] Transfer ownership to an active member atomically, preserving exactly one owner and removing the former owner's owner-only powers.
- [ ] Exercise each permission independently through HTTP. Revocation takes effect on subsequent requests, not only after token expiry.
- [ ] Expose membership and delegated-permission administration through authorized use cases. Verify staff administrative authority is distinct from workspace delegation and cannot grant staff/superuser privileges through workspace permissions.
