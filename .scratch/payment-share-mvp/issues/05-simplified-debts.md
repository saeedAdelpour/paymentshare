# 05: Show simplified repayment suggestions

Status: ready-for-agent

**What to build:** Read deterministic repayment suggestions from workspace net balances.

**Blocked by:** 03: Record equal-split expenses and read balances.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 17, 18, 42.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. No UI is in scope. Follow the parent spec except for the explicitly approved edit-party clarification in ticket 07; the parent is preserved unchanged during ticket publication.

## Acceptance criteria

- [ ] Match the largest remaining debtor and creditor repeatedly, using member ID for stable ties. Do not promise globally minimal payment count.
- [ ] B owes A $10 and A owes C $10 produces B paying C $10, leaving original entries unchanged.
- [ ] Generate suggestions from one consistent snapshot of stored balances, not transaction-history scans. Suggestions settle all balances and are stable for unchanged input.
- [ ] Expose workspace suggestions only to authorized active members at this stage; maintain immediate freshness after writes and workspace isolation.

