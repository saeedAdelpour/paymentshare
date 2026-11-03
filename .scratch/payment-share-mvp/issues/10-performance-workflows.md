# 10: Verify MVP performance and complete workflows

Status: ready-for-agent

**What to build:** Verify complete MVP workflows and meet documented initial load-test targets.

**Blocked by:** 09: Leave, remove, and rejoin workspaces.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 38, 43, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, QOR5, OIDC, and GORM/PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. Reuse the completed [project foundation](../../payment-share-foundation/spec.md). Register any application tables introduced by this slice in administration and route administrative mutations through the same application use cases; preserve financial invariants and history. Follow the current parent spec and the approved foundation stack and quality checks. Verify administrative behavior through browser/session or direct administrative HTTP requests with real PostgreSQL.

## Acceptance criteria

- [ ] Extend or audit the OpenAPI contract for this slice, documenting introduced/changed application API operations, inputs, success responses, relevant errors, authentication, and required headers. Verify HTTP requests/responses against the contract and pass contract validation and generated-Go freshness checks; keep Swagger UI usable with the foundation OIDC authorization flow.
- [ ] Run authenticated HTTP load tests using real PostgreSQL, workspaces up to 100 members, and 100,000 historical entries per tested workspace at 100 requests/second.
- [ ] Measure p95 balance reads below 200 ms and writes below 500 ms; report hardware, configuration, request mix, duration, data distribution, errors, and actual achieved throughput. Fix failures and rerun before declaring acceptance.
- [ ] Exercise concurrent duplicate and distinct-key creations, edits, cancellations, and membership changes; verify final observable balances and retained history, including that every successful mutation is represented once.
- [ ] Verify complete expense, simplified repayment, overpayment, correction, removal, and rejoin journeys using the approved HTTP test seam.
- [ ] Run the full relevant test suite and document reproducible setup and commands. Make no 10M DAU capacity claim.
- [ ] Audit the complete spec and story coverage across these tickets; fix missing MVP behavior without introducing deferred features.
- [ ] Verify complete administration journeys with restricted staff and superusers alongside API journeys, including financial correction/voiding, permissions changed during an existing session, and deactivated-user denial. Confirm all implemented application tables have appropriate administrative views and protected mutations.

## Further Notes

Completion of this ticket gates the later React member-facing application. Audit the complete implemented application API against OpenAPI and confirm working Swagger authorization and contract checks before backend acceptance. React implementation and detailed client design remain outside this ticket.
