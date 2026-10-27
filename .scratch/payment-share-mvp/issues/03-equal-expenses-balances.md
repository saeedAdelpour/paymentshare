# 03: Record equal-split expenses and read balances

Status: ready-for-agent

**What to build:** Create equal-split expenses and immediately read persisted net balances safely under retries and concurrent writes.

**Blocked by:** 02: Add members and delegate permissions.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 10, 11, 13, 14, 15, 16, 30, 39, 40, 42.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. No UI is in scope. Follow the parent spec except for the explicitly approved edit-party clarification in ticket 07; the parent is preserved unchanged during ticket publication.

## Acceptance criteria

- [ ] Record a positive amount in smallest currency units, description, occurrence date, creator, payer, and explicitly selected active participants. Allow backdating; include payer only if selected.
- [ ] Ordinary creators must be the payer; owners and authorized editors can record on behalf of another member. Validate participant membership and workspace currency.
- [ ] A pays $10 for B and C: A is owed $10, B and C owe $5 each. All workspace net balances sum to zero.
- [ ] Allocate fractional split remainders deterministically by largest remainder, then member ID. $20.23 split three ways becomes $6.75, $6.74, $6.74, never $20.
- [ ] Freeze currency after the first recorded expense. Maintain balances on writes and read them without scanning transaction history.
- [ ] Commit expense, initial revision, balance deltas, and idempotency result in one database transaction. Failed validation or persistence leaves no partial effects.
- [ ] Require a client request key scoped to actor and workspace, retained for workspace lifetime. Same operation and key returns the original result without another effect; different operation or payload with the same key conflicts; different keys permit identical expenses.
- [ ] Enforce duplicate-key uniqueness in PostgreSQL, not an in-memory check. Concurrent requests with one key produce one entry and one balance effect; simulate a lost response followed by retry.
- [ ] Concurrent different keys affecting the same members preserve all successful deltas. Use database locking, atomic updates, or equivalent serialization, and document the strategy.
- [ ] Coordinate authorization, participant membership, and currency validation with the write to prevent time-of-check races. Immediate reads after success reflect committed balances; retries recheck current access.

