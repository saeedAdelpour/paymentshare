# 09: Leave, remove, and rejoin workspaces

Status: ready-for-agent

**What to build:** End and restore participation while retaining balances, records, and narrowly scoped removed-member access.

**Blocked by:** 05: Show simplified repayment suggestions; 08: Browse workspace and personal history.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 31, 32, 33, 34, 35, 36, 37, 38.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. No UI is in scope. Follow the parent spec except for the explicitly approved edit-party clarification in ticket 07; the parent is preserved unchanged during ticket publication.

## Acceptance criteria

- [ ] Allow voluntary departure and authorized removal only when net balance is zero; owner transfers ownership first. Check net balance atomically with removal against concurrent financial writes.
- [ ] Removed users can read only their own balance and suggestions involving them, including counterparties. Deny all history, other balances, and writes, even on their own old entries.
- [ ] Corrections may affect existing removed participants without rejoining; new expenses cannot include them. Include removed nonzero balances in suggestions.
- [ ] An active repayment party or owner/authorized editor can record repayments involving a removed member.
- [ ] Test remove at zero, correct an old entry to create debt, read restricted balance, and repay without rejoining. Net-zero intermediary members may leave despite original pairwise debts.
- [ ] Rejoin the same membership/history/balance and restore normal visibility with delegated permissions empty.
- [ ] Test old cursors, direct identifiers, and idempotency replay after removal or permission revocation. Coordinate membership and authorization changes with financial writes to prevent stale permission races.

