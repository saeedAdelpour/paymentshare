# 10: Verify MVP performance and complete workflows

Status: ready-for-agent

**What to build:** Verify complete MVP workflows and meet documented initial load-test targets.

**Blocked by:** 09: Leave, remove, and rejoin workspaces.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 43.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. No UI is in scope. Follow the parent spec except for the explicitly approved edit-party clarification in ticket 07; the parent is preserved unchanged during ticket publication.

## Acceptance criteria

- [ ] Run authenticated HTTP load tests using real PostgreSQL, workspaces up to 100 members, and 100,000 historical entries per tested workspace at 100 requests/second.
- [ ] Measure p95 balance reads below 200 ms and writes below 500 ms; report hardware, configuration, request mix, duration, data distribution, errors, and actual achieved throughput. Fix failures and rerun before declaring acceptance.
- [ ] Exercise concurrent duplicate and distinct-key creations, edits, cancellations, and membership changes; verify final observable balances and retained history, including that every successful mutation is represented once.
- [ ] Verify complete expense, simplified repayment, overpayment, correction, removal, and rejoin journeys using the approved HTTP test seam.
- [ ] Run the full relevant test suite and document reproducible setup and commands. Make no 10M DAU capacity claim.
- [ ] Audit the complete spec and story coverage across these tickets; fix missing MVP behavior without introducing deferred features.

