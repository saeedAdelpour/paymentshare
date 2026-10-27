# 08: Browse workspace and personal history

Status: ready-for-agent

**What to build:** Browse cursor-paginated entries and revision history across expenses, repayments, and cancellations.

**Blocked by:** 07: Edit and cancel entries with retained history.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 28, 29, 30.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. No UI is in scope. Follow the parent spec except for the explicitly approved edit-party clarification in ticket 07; the parent is preserved unchanged during ticket publication.

## Acceptance criteria

- [ ] Active members can read all workspace entries and filter to entries involving themselves. Do not expose other workspaces through filters, cursors, or identifiers.
- [ ] Include descriptions, occurrence dates, creation/edit times, actors, revisions, and cancellation state; retain voided history.
- [ ] Use a stable cursor order with a unique tie-breaker and document pagination behavior under concurrent writes; test equal timestamps and backdated entries.
- [ ] Current balances include all nonvoided entries irrespective of occurrence date. Do not implement historical balance-as-of-date queries.

