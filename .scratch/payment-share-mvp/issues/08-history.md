# 08: Browse workspace and personal history

Status: ready-for-agent

**What to build:** Browse cursor-paginated entries and revision history across expenses, repayments, and cancellations.

**Blocked by:** 07: Edit and cancel entries with retained history.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 28, 29, 30, 45, 46, 47, 56, 58, 59, 60.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, QOR5, OIDC, and GORM/PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. Reuse the completed [project foundation](../../payment-share-foundation/spec.md). Register any application tables introduced by this slice in administration and route administrative mutations through the same application use cases; preserve financial invariants and history. Follow the current parent spec and the approved foundation stack and quality checks. Verify administrative behavior through browser/session or direct administrative HTTP requests with real PostgreSQL.

## Acceptance criteria

- [ ] Extend or audit the OpenAPI contract for this slice, documenting introduced/changed application API operations, inputs, success responses, relevant errors, authentication, and required headers. Verify HTTP requests/responses against the contract and pass contract validation and generated-Go freshness checks; keep Swagger UI usable with the foundation OIDC authorization flow.
- [ ] Active members can read all workspace entries and filter to entries involving themselves. Do not expose other workspaces through filters, cursors, or identifiers.
- [ ] Include descriptions, occurrence dates, creation/edit times, actors, revisions, and cancellation state; retain voided history.
- [ ] Use a stable cursor order with a unique tie-breaker and document pagination behavior under concurrent writes; test equal timestamps and backdated entries.
- [ ] Current balances include all nonvoided entries irrespective of occurrence date. Do not implement historical balance-as-of-date queries.
- [ ] Expose permitted administrative history and revision inspection, including voided records, with filtering/pagination behavior verified. Revisions remain immutable through administrative actions.
