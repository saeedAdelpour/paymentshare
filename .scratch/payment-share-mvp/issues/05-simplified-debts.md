# 05: Show simplified repayment suggestions

Status: ready-for-agent

**What to build:** Read deterministic repayment suggestions from workspace net balances.

**Blocked by:** 03: Record equal-split expenses and read balances.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 17, 18, 42, 45, 46, 47, 56, 58, 59, 60.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, QOR5, OIDC, and GORM/PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. Reuse the completed [project foundation](../../payment-share-foundation/spec.md). Register any application tables introduced by this slice in administration and route administrative mutations through the same application use cases; preserve financial invariants and history. Follow the current parent spec and the approved foundation stack and quality checks. Verify administrative behavior through browser/session or direct administrative HTTP requests with real PostgreSQL.

## Acceptance criteria

- [ ] Extend or audit the OpenAPI contract for this slice, documenting introduced/changed application API operations, inputs, success responses, relevant errors, authentication, and required headers. Verify HTTP requests/responses against the contract and pass contract validation and generated-Go freshness checks; keep Swagger UI usable with the foundation OIDC authorization flow.
- [ ] Match the largest remaining debtor and creditor repeatedly, using member ID for stable ties. Do not promise globally minimal payment count.
- [ ] B owes A $10 and A owes C $10 produces B paying C $10, leaving original entries unchanged.
- [ ] Generate suggestions from one consistent snapshot of stored balances, not transaction-history scans. Suggestions settle all balances and are stable for unchanged input.
- [ ] Expose workspace suggestions only to authorized active members at this stage; maintain immediate freshness after writes and workspace isolation.
- [ ] Expose permitted balance/suggestion inspection in administration without direct balance editing. Check administrative resource permissions and deterministic suggestions from a consistent snapshot.
