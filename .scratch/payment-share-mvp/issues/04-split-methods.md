# 04: Support exact, percentage, and weighted splits

Status: ready-for-agent

**What to build:** Create expenses using exact amounts, percentages, and relative weights.

**Blocked by:** 03: Record equal-split expenses and read balances.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 12, 13, 45, 46, 47, 48, 56, 58, 59, 60.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, QOR5, OIDC, and GORM/PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. Reuse the completed [project foundation](../../payment-share-foundation/spec.md). Register any application tables introduced by this slice in administration and route administrative mutations through the same application use cases; preserve financial invariants and history. Follow the current parent spec and the approved foundation stack and quality checks. Verify administrative behavior through browser/session or direct administrative HTTP requests with real PostgreSQL.

## Acceptance criteria

- [ ] Extend or audit the OpenAPI contract for this slice, documenting introduced/changed application API operations, inputs, success responses, relevant errors, authentication, and required headers. Verify HTTP requests/responses against the contract and pass contract validation and generated-Go freshness checks; keep Swagger UI usable with the foundation OIDC authorization flow.
- [ ] Exact shares sum to total; percentages sum to 100%; relative weights determine proportional shares. Reject invalid or ambiguous split inputs without writes.
- [ ] Use exact monetary arithmetic, deterministic largest-remainder allocation, and member-ID tie-breaking; preserve the total and zero-sum workspace balances.
- [ ] Test unequal splits, tiny totals, fractional percentages, and ties through HTTP; supplementary domain tests may cover arithmetic edge cases.
- [ ] All new methods retain creation authorization, currency locking, atomic updates, and idempotency guarantees of equal splits.
- [ ] Exercise each split method through administration with the same validation and deterministic rounding as the API. Unauthorized or invalid administrative requests leave balances and history unchanged.
