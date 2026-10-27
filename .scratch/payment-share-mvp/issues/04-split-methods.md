# 04: Support exact, percentage, and weighted splits

Status: ready-for-agent

**What to build:** Create expenses using exact amounts, percentages, and relative weights.

**Blocked by:** 03: Record equal-split expenses and read balances.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 12, 13.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. No UI is in scope. Follow the parent spec except for the explicitly approved edit-party clarification in ticket 07; the parent is preserved unchanged during ticket publication.

## Acceptance criteria

- [ ] Exact shares sum to total; percentages sum to 100%; relative weights determine proportional shares. Reject invalid or ambiguous split inputs without writes.
- [ ] Use exact monetary arithmetic, deterministic largest-remainder allocation, and member-ID tie-breaking; preserve the total and zero-sum workspace balances.
- [ ] Test unequal splits, tiny totals, fractional percentages, and ties through HTTP; supplementary domain tests may cover arithmetic edge cases.
- [ ] All new methods retain creation authorization, currency locking, atomic updates, and idempotency guarantees of equal splits.

