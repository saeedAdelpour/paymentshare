# 06: Record external repayments

Status: ready-for-agent

**What to build:** Record money moved outside the app and immediately update balances.

**Blocked by:** 03: Record equal-split expenses and read balances.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 19, 20, 21, 22, 30, 39, 40, 42, 45, 46, 47, 48, 56, 58, 59, 60.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, QOR5, OIDC, and GORM/PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. Reuse the completed [project foundation](../../payment-share-foundation/spec.md). Register any application tables introduced by this slice in administration and route administrative mutations through the same application use cases; preserve financial invariants and history. Follow the current parent spec and the approved foundation stack and quality checks. Verify administrative behavior through browser/session or direct administrative HTTP requests with real PostgreSQL.

## Acceptance criteria

- [ ] Extend or audit the OpenAPI contract for this slice, documenting introduced/changed application API operations, inputs, success responses, relevant errors, authentication, and required headers. Verify HTTP requests/responses against the contract and pass contract validation and generated-Go freshness checks; keep Swagger UI usable with the foundation OIDC authorization flow.
- [ ] Either active party, owner, or authorized editor can record a positive repayment between distinct members of the same workspace, with description, occurrence date, and creator.
- [ ] No approval or actual money movement is performed. Permit partial, exact, excess, and unsolicited repayments regardless of current suggestions.
- [ ] Starting from A owing B $10, independently verify $5, $10, and $15 repayments leave $5 debt, zero, and B owing A $5.
- [ ] Lock currency when repayment is the first recorded entry. Preserve zero-sum balances and backdated occurrence semantics.
- [ ] Atomically commit repayment, initial revision, deltas, and request-key result; reuse the actor/workspace `Idempotency-Key` namespace, lifetime, fingerprint, and replay authorization contract from expenses. Lock affected workspace-member balances in the shared consistent order and preserve the documented membership/currency coordination protocol.
- [ ] HTTP tests cover simultaneous same-key creation, distinct-key writes overlapping expense creation, conflicting payloads, rollback, and retry after response loss. No lost updates or duplicated effects.
- [ ] Exercise administrative repayment recording through the same transactional use case and idempotency contract. Preserve actual parties and creator, zero-sum balances, and failed-write rollback; enforce staff permissions.
