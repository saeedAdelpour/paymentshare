# 07: Edit and cancel entries with retained history

Status: ready-for-agent

**What to build:** Directly correct expenses and repayments or cancel erroneous entries without losing history.

**Blocked by:** 04: Support exact, percentage, and weighted splits; 06: Record external repayments.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 23, 24, 25, 26, 27, 39, 40, 41, 42, 45, 46, 47, 48, 56, 58, 59, 60.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, QOR5, OIDC, and GORM/PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. Reuse the completed [project foundation](../../payment-share-foundation/spec.md). Register any application tables introduced by this slice in administration and route administrative mutations through the same application use cases; preserve financial invariants and history. Follow the current parent spec and the approved foundation stack and quality checks. Verify administrative behavior through browser/session or direct administrative HTTP requests with real PostgreSQL.

## Acceptance criteria

- [ ] Extend or audit the OpenAPI contract for this slice, documenting introduced/changed application API operations, inputs, success responses, relevant errors, authentication, and required headers. Verify HTTP requests/responses against the contract and pass contract validation and generated-Go freshness checks; keep Swagger UI usable with the foundation OIDC authorization flow.
- [ ] Active creators edit their own entries; owners and authorized editors edit others. Apply the same authorization to cancellation.
- [ ] Approved clarification: direct edits may change expense payer and participants and repayment sender/recipient as well as amount, split, description, and occurrence date. Revalidate resulting membership and recording permissions; preserve original creator.
- [ ] Preserve previous values and actors in revisions. Commit old-effect reversal, new effect, revision, version increment, and idempotency result atomically.
- [ ] Correct a fully repaid $10 share to $6: original creditor owes $4 back; keep the actual repayment unchanged.
- [ ] Cancellation marks an entry void and reverses its effect, retaining who and when. Repeated cancellation never reverses twice; cancelled entries cannot silently regain effects through edits.
- [ ] Apply optimistic compare-and-update using the client's expected entry version for edits and cancellation; stale versions return 409 Conflict without effects. Concurrent edits from one version cannot both win; cover edit versus cancellation races.
- [ ] Use the shared lifetime `Idempotency-Key` contract for edits and cancellation, including operation, target, payload, and expected version in the fingerprint. Lock the union of old/new affected workspace-member balances in the consistent order. Duplicate retries replay safely, changed payload conflicts, and current access is checked before replay.
- [ ] Expose revisions for authorized reads and test before/after values and balances through HTTP.
- [ ] Exercise administrative edit/void actions with staff permissions and superuser authority. Verify API-versus-admin stale-edit conflicts, retained creator/revisions, atomic old/new balance effects, and idempotent retries. No physical deletion or direct editing of balances/revisions is available.
