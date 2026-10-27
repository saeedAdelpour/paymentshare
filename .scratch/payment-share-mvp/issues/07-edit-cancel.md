# 07: Edit and cancel entries with retained history

Status: ready-for-agent

**What to build:** Directly correct expenses and repayments or cancel erroneous entries without losing history.

**Blocked by:** 04: Support exact, percentage, and weighted splits; 06: Record external repayments.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 23, 24, 25, 26, 27, 39, 40, 41, 42.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. No UI is in scope. Follow the parent spec except for the explicitly approved edit-party clarification in ticket 07; the parent is preserved unchanged during ticket publication.

## Acceptance criteria

- [ ] Active creators edit their own entries; owners and authorized editors edit others. Apply the same authorization to cancellation.
- [ ] Approved clarification: direct edits may change expense payer and participants and repayment sender/recipient as well as amount, split, description, and occurrence date. Revalidate resulting membership and recording permissions; preserve original creator.
- [ ] Preserve previous values and actors in revisions. Commit old-effect reversal, new effect, revision, version increment, and idempotency result atomically.
- [ ] Correct a fully repaid $10 share to $6: original creditor owes $4 back; keep the actual repayment unchanged.
- [ ] Cancellation marks an entry void and reverses its effect, retaining who and when. Repeated cancellation never reverses twice; cancelled entries cannot silently regain effects through edits.
- [ ] Require expected versions for mutations susceptible to stale writes. Concurrent edits from one version cannot both win; cover edit versus cancellation races.
- [ ] Use the shared lifetime request-key contract for edits and cancellation. Duplicate retries replay safely, changed payload conflicts, and current access is checked before replay.
- [ ] Expose revisions for authorized reads and test before/after values and balances through HTTP.

