# 06: Record external repayments

Status: ready-for-agent

**What to build:** Record money moved outside the app and immediately update balances.

**Blocked by:** 03: Record equal-split expenses and read balances.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 19, 20, 21, 22, 30, 39, 40, 42.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. No UI is in scope. Follow the parent spec except for the explicitly approved edit-party clarification in ticket 07; the parent is preserved unchanged during ticket publication.

## Acceptance criteria

- [ ] Either active party, owner, or authorized editor can record a positive repayment between distinct members of the same workspace, with description, occurrence date, and creator.
- [ ] No approval or actual money movement is performed. Permit partial, exact, excess, and unsolicited repayments regardless of current suggestions.
- [ ] Starting from A owing B $10, independently verify $5, $10, and $15 repayments leave $5 debt, zero, and B owing A $5.
- [ ] Lock currency when repayment is the first recorded entry. Preserve zero-sum balances and backdated occurrence semantics.
- [ ] Atomically commit repayment, initial revision, deltas, and request-key result; reuse the actor/workspace key namespace and lifetime from expenses.
- [ ] HTTP tests cover simultaneous same-key creation, distinct-key writes overlapping expense creation, conflicting payloads, rollback, and retry after response loss. No lost updates or duplicated effects.

