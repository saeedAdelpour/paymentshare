# 03: Record equal-split expenses and read balances

Status: ready-for-agent

**What to build:** Create equal-split expenses and immediately read persisted net balances safely under retries and concurrent writes.

**Blocked by:** 02: Add members and delegate permissions.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 10, 11, 13, 14, 15, 16, 30, 39, 40, 42, 45, 46, 47, 48, 56, 58, 59, 60.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, QOR5, OIDC, and GORM/PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. Reuse the completed [project foundation](../../payment-share-foundation/spec.md). Register any application tables introduced by this slice in administration and route administrative mutations through the same application use cases; preserve financial invariants and history. Follow the current parent spec and the approved foundation stack and quality checks. Verify administrative behavior through browser/session or direct administrative HTTP requests with real PostgreSQL.

## Acceptance criteria

- [ ] Extend or audit the OpenAPI contract for this slice, documenting introduced/changed application API operations, inputs, success responses, relevant errors, authentication, and required headers. Verify HTTP requests/responses against the contract and pass contract validation and generated-Go freshness checks; keep Swagger UI usable with the foundation OIDC authorization flow.
- [ ] Record a positive amount in smallest currency units, description, occurrence date, creator, payer, and explicitly selected active participants. Allow backdating; include payer only if selected.
- [ ] Ordinary creators must be the payer; owners and authorized editors can record on behalf of another member. Validate participant membership and workspace currency.
- [ ] A pays $10 for B and C: A is owed $10, B and C owe $5 each. All workspace net balances sum to zero.
- [ ] Allocate fractional split remainders deterministically by largest remainder, then member ID. $20.23 split three ways becomes $6.75, $6.74, $6.74, never $20.
- [ ] Freeze currency after the first recorded expense. Maintain balances on writes and read them without scanning transaction history.
- [ ] Commit expense, initial revision, balance deltas, and idempotency result in one database transaction. Failed validation or persistence leaves no partial effects.
- [ ] Require a client `Idempotency-Key` scoped to actor and workspace, retained for workspace lifetime. Same operation and key returns the original result without another effect; different operation or payload with the same key conflicts; different keys permit identical expenses.
- [ ] Enforce duplicate-key uniqueness in PostgreSQL, not an in-memory check. Claim the key and commit its operation/target/payload fingerprint and replay result atomically with the financial mutation. A failed transaction commits neither effects nor a completed reservation. Concurrent requests with one key produce one entry and one balance effect; simulate a lost response followed by retry.
- [ ] Concurrent different keys affecting the same members preserve all successful deltas. Coordinate updates using database locks on balances scoped to workspace and member; acquire affected rows in a documented consistent order. Do not substitute a global user balance or a blanket workspace mutation lock. Use optimistic versions for edits/cancellation when those operations are introduced; a complex automatic retry framework is not required.
- [ ] Document and test the lock protocol coordinating authorization, participant membership, currency initialization, and balance writes, including concurrent ownership/permission changes. Acquire balance rows affected by old or new values in a consistent order; balance locks alone do not protect separate authorization state. Immediate reads after success reflect committed balances; retries recheck current access.
- [ ] Expose expense recording and balance inspection in administration with resource/action permissions. Verify administrative recording preserves creator/payer distinction, split validation, zero-sum balances, atomic effects, and idempotency; balances have no direct-edit operation.
