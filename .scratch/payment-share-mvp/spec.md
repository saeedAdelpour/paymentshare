# Payment Share MVP backend

Status: ready-for-agent

Readiness: The backend behavioral baseline and HTTP API testing boundary with real PostgreSQL and a controlled local identity provider are agreed. The user approved GORM + QOR5 + Goose and the three-ticket foundation sequence. Complete the [foundation tickets](../payment-share-foundation/spec.md) before the existing feature sequence. The selected concurrency approach is optimistic entry versions with workspace-member balance locks and durable idempotency keys. Concrete aggregate interfaces and the lock ordering protocol are implementation work within the approved requirements.

## Problem Statement

People sharing expenses need to know what they owe across separate groups without manually reconstructing transaction history. Corrections, repayments, concurrent requests, and changes in membership must not leave balances incorrect or expose another workspace's records.

## Solution

Provide a backend API for workspaces with one currency each, shared expenses, externally made repayments, direct corrections with revision history, and cancellation of mistaken entries. Show current net balances and simplified repayment suggestions immediately after changes. Enforce workspace permissions and retain limited balance access for removed members. Include a Django-like administration panel for application tables and users, with staff permissions and superuser access that preserve financial integrity and retained history.

## User Stories

1. As a user, I want to authenticate through an identity provider, so that I can securely access my workspaces.
2. As a user, I want to create a workspace, so that I can track a group's expenses separately.
3. As a user, I want to belong to multiple workspaces, so that unrelated groups' balances stay separate.
4. As an owner, I want to choose a workspace currency, so that every amount has an unambiguous meaning.
5. As an owner, I want to rename my workspace, so that its name stays useful.
6. As an owner, I want to add a registered user by exact verified email, so that I can bring them into the group.
7. As an owner, I want to grant separate add-member, remove-member, and edit-others permissions, so that responsibilities can be delegated independently.
8. As an authorized member, I want to add users without receiving unrelated permissions, so that my access matches my responsibilities.
9. As an owner, I want to transfer ownership to an active member, so that the workspace can continue when I leave.
10. As a payer, I want to record an expense for selected members, so that their shares are reflected in balances.
11. As a payer, I want to include myself only when I share an expense, so that paying does not automatically make me a beneficiary.
12. As a member, I want equal, exact-amount, percentage, and relative-weight splits, so that I can describe different sharing arrangements.
13. As a member, I want rounding to preserve the expense total, so that no money appears or disappears.
14. As an owner or authorized editor, I want to record an expense on another member's behalf, so that bookkeeping can be delegated.
15. As a member, I want the creator and payer distinguished, so that I know who entered an expense and who paid it.
16. As a member, I want to see my net balance, so that I know my overall position in a workspace.
17. As a member, I want simplified repayment suggestions, so that intermediate debts do not require unnecessary transfers.
18. As a member, I want suggestions to be consistent for unchanged balances, so that repeated reads are predictable.
19. As a repayment party, I want to record an external repayment without waiting for confirmation, so that balances reflect money already moved.
20. As a member, I want to record partial, exact, and excess repayments, so that actual payments are represented accurately.
21. As a member, I want to record a repayment to someone other than my suggested recipient, so that the app reflects what actually happened.
22. As an owner or authorized editor, I want to record repayments for other members, so that I can maintain group records.
23. As an active creator, I want to correct my entries directly, so that mistakes can be fixed without recreating them.
24. As an owner or authorized editor, I want to correct other members' entries, so that shared records remain accurate.
25. As a member, I want earlier revisions and their actors retained, so that changes are explainable.
26. As a member, I want corrections after repayment to update balances, even when debt direction reverses, so that the current position is accurate.
27. As an authorized user, I want to cancel a mistaken or duplicated entry, so that its balance effect is reversed without losing history.
28. As an active member, I want to see workspace history and filter it to entries involving me, so that I can review relevant activity.
29. As a member, I want descriptions, occurrence dates, and paginated history, so that I can navigate older activity.
30. As a member, I want to record backdated activity, so that delayed bookkeeping remains possible.
31. As a member, I want to leave when my net balance is zero, so that I can stop participating without losing retained records.
32. As an authorized member, I want to remove a zero-balance member, so that access follows group participation.
33. As a removed member, I want read-only access to my own balance and suggestions involving me, so that I can settle changes affecting me.
34. As an active authorized user, I want to correct existing entries involving removed members, so that corrections do not require re-adding them.
35. As an active repayment party or authorized editor, I want to record repayments involving removed members, so that their remaining balances can be settled.
36. As a returning member, I want my existing history and balance reconnected, so that rejoining does not create a second identity in the workspace.
37. As an owner, I want rejoined members to start without delegated permissions, so that old privileges do not silently return.
38. As a user, I want workspace access enforced on every request, so that other groups' data stays private.
39. As an API client, I want retries to have only one effect, so that network failures do not duplicate transactions.
40. As an API client, I want conflicting reuse of request keys rejected, so that different operations are not mistaken for retries.
41. As an editor, I want stale edits rejected, so that I do not overwrite another person's correction unknowingly.
42. As a member, I want successful changes immediately reflected in balance reads, so that I can trust the displayed amounts.
43. As an operator, I want reproducible performance measurements, so that capacity claims are grounded in evidence.
44. As a developer, I want local authentication and database-backed tests, so that I can verify the API without a paid identity service.
45. As a service administrator, I want to browse application tables and users in an administration panel, so that I can operate Payment Share across workspaces.
46. As a staff user, I want administrative read, write, and delete permissions enforced, so that my access follows my assigned responsibilities.
47. As a superuser, I want all administrative permissions without individual grants, so that I can administer every application resource.
48. As a service administrator, I want financial corrections and cancellation to preserve balances and revision history, so that administrative operations maintain the same financial guarantees as member operations.
49. As a developer, I want tested migrations to a chosen schema version in either direction, so that I can switch application versions locally.
50. As an operator, I want an explicit production rollback and data recovery strategy, so that reverting an application release does not silently lose data.
51. As a superuser, I want exclusive authority to manage staff status, administrative permissions, and superuser promotion, so that administrative delegation is controlled.
52. As an authorized administrator, I want to deactivate and reactivate users while retaining their financial references and history, so that access can be managed without destroying shared records.
53. As a developer, I want pinned formatting, linting, vulnerability, and test commands shared by local development and CI, so that application quality checks are reproducible.
54. As a staff user, I want browser login through the same identity provider as application users, so that administrative access uses a shared identity.
55. As an operator, I want an idempotent CLI command provisioning a superuser for a specific provider identity, so that initial administration is reproducible.
56. As a service administrator, I want each new feature's records and permitted actions available in administration, so that the panel stays complete as the application grows.
57. As a developer, I want authentication, administration, migrations, and quality checks working before feature development, so that later tickets reuse a verified foundation.

58. As a developer, I want interactive Swagger documentation to grow with each API slice, so that I can exercise the backend before the member-facing application exists.
59. As an API consumer, I want complete OpenAPI inputs, success responses, relevant errors, authentication, and required headers, so that I can integrate against a clear contract.
60. As a developer, I want contract validation, generated-code freshness checks, and HTTP contract behavior tests, so that documentation and the running API stay aligned.
61. As a developer, I want to start the React member-facing application after MVP ticket 10, so that the backend contract and complete workflows are verified first.

## Implementation Decisions

### Foundation delivery order

- Complete project initialization before the existing feature tickets: authentication, administration panel with staff permissions and superuser support, migration management, and code formatting/linting checks. The approved [foundation spec and tickets](../payment-share-foundation/spec.md) define delivery; [foundation-plan.md](foundation-plan.md) preserves the discussion summary.
- Prove the foundation with a running local application, persistent users and staff permissions, browser administration, repeatable migration commands, and reproducible quality checks. As feature tables are introduced, register their administrative views and route mutations through the appropriate application use cases.
- Use gofumpt as the formatter and golangci-lint v2 with govet, staticcheck, errcheck, ineffassign, and unused initially. Run govulncheck and race-enabled Go tests alongside real PostgreSQL integration tests. Pin tool versions and use matching local/CI commands.
- Further infrastructure and non-functional requirements can be refined later; the user has not required an exhaustive infrastructure design before this foundation phase.

### Service and identity

- Maintain a contract-first OpenAPI API using oapi-codegen-generated Go transport types and Chi server interfaces, with Swagger UI established by foundation. Each API slice documents all introduced/changed operations, inputs, success responses, relevant errors, authentication, and required headers (including idempotency/version requirements where applicable). Local/CI checks reject invalid contracts and stale generated code; HTTP behavior tests verify requests/responses against the contract. Swagger authorization uses OIDC Authorization Code with PKCE and API access tokens. QOR5 internal browser routes and identity-provider endpoints are outside the application API contract. See [ADR 0007](../../docs/adr/0007-openapi-contract-first-http-api.md).
- Implement the member-facing application with React only after [ticket 10](issues/10-performance-workflows.md) is complete; defer detailed React design until then. Swagger UI grows with the backend throughout foundation and MVP delivery, and QOR5 remains the staff administration panel.
- Use the existing Go module and Chi router with PostgreSQL, initially as one backend API service.
- Use a configurable OpenID Connect provider for registration and login. Configure issuer and audience, validate access tokens locally, and check current workspace membership and permissions on each request.
- Provide a local development identity provider. Hosted provider selection is a deployment decision.
- Add already registered users immediately by exact verified email; provide no searchable directory or invitation acceptance flow.

### Domain design

- Use DDD within the existing single Payment Share bounded context, as already required by all implementation tickets. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions.
- Keep HTTP, administration UI, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants without loading all workspace history as one object graph. Concrete aggregate and module interfaces are implementation choices; the proposed Workspace/Expense/Repayment aggregate sketch was not separately approved and does not prescribe a mandatory structure.
- Use GORM, QOR5, and Goose, approved after the [implementation research](implementation-research.md). The foundation verifies staff identity integration, superuser behavior, custom mutation handlers, and migration behavior. Workload performance still requires measurement.

### Administration

- Include the administration panel in the first release. Register application tables and users for administrative browsing, with staff read/write/delete permissions and a superuser granted every administrative permission without individual grants.
- Administrative permissions are separate from workspace ownership and delegated member permissions. The panel supports authorized operations across workspaces.
- Superuser authority preserves domain invariants and retained history. Financial edits run through application use cases that atomically update entries, revisions, balances, and request-key results. Administrative cancellation of an expense or repayment means voiding rather than physical deletion.
- Balances and revisions are inspectable but have no ordinary direct-edit operation, including for superusers. Administrative mutation handlers must not bypass the financial use cases through generic table CRUD.
- Staff authenticate through the same OIDC provider as application users. Payment Share stores staff status, superuser status, roles, and administrative permissions. Ordinary registration grants no administrative privileges.
- Provide a repeatable, idempotent CLI command granting the initial superuser role to a specific provider identity.
- Only superusers may grant or revoke staff status, change administrative permissions, or promote other superusers initially. Enforce this on the server, including direct requests to role/permission resources; ordinary record-write permissions do not confer permission-management authority.
- User deletion in the initial administration panel means deactivation rather than physical deletion. Deactivation blocks access to Payment Share, including administration, while retaining memberships, financial records, and actor history; it does not delete or disable the external identity-provider account. Successful OIDC login must not implicitly reactivate a deactivated user. Reactivation restores access subject to current permissions.
- Detailed browser/session configuration and treatment of sensitive operational fields remain implementation design work.

### Schema migrations and rollback

- Use reviewed, versioned migrations supporting a target schema version in either direction for local development. Verify Up → Down → Up against representative data; document any data-preservation limitations.
- Local version switching must account for divergent migration histories. Revert with the departing version's migration files available, then apply the destination version's migrations from a compatible common state. Separate local databases remain an option when histories or data transformations make reuse unsuitable.
- Production rollback has a separate guarantee: preserve compatibility with the previous application release where practical, and define an explicit recovery strategy for destructive schema changes or data transformations. A Down migration alone does not establish data recoverability.
- Deploy the API and administration panel together initially, with PostgreSQL and an explicit migration step before the new application starts. Brief maintenance windows are acceptable for the initial release; startup checks schema compatibility.
- The migration runner and deployment execution process must be reproducible. Use Goose as the migration manager; ORM schema synchronization is not a substitute for versioned migration history.

### Workspace access

- A user may belong to multiple workspaces. Balances and repayment suggestions never cross workspace boundaries.
- Each workspace has exactly one owner and one currency. Currency becomes fixed after its first expense or repayment is recorded.
- Support workspace creation, renaming, membership management, delegated permissions, and ownership transfer.
- Only the owner grants permissions. Adding members, removing members, and editing others' entries are separate permissions.
- Ownership must transfer to another active member before the owner leaves; the former owner loses owner-only powers.
- Active members can read all workspace expenses and repayments and request a personal filtered view. Users who never joined have no workspace access.
- Removal or voluntary departure requires zero workspace net balance, not settlement of every original pairwise debt.
- Removed members retain read-only access only to their own balance and suggestions involving them, including counterparties. They cannot read transaction history, other balances, or write entries, even ones they created.
- Corrections to existing entries can affect removed members without restoring membership. New expenses require active participants.
- Rejoining reconnects the same history and balance, restores ordinary visibility, and starts with no delegated permissions.

### Expenses, repayments, and corrections

- An expense has one payer and selected participants with shares. Include the payer in the split only when explicitly selected.
- Ordinary members may create expenses only as payer. Owners and authorized editors may create on another member's behalf; retain creator separately from payer.
- Support equal, exact, percentage, and relative-weight splits. Exact shares sum to the total; percentages sum to 100%. Use largest fractional remainder to allocate leftover smallest units for calculated splits, with member ID as tie-breaker.
- Represent monetary amounts in the currency's smallest unit with positive expense and repayment totals. No currency conversion.
- Repayments record externally moved money and take effect immediately without counterparty approval. Permit amounts below, equal to, or above debt and recipients unrelated to current suggestions, including when no debt currently exists.
- Either active repayment party, the owner, or an authorized editor can record a repayment. Current members may record a repayment involving a removed member when they are the other party; owners and authorized editors may also record these.
- Creators who remain active may edit their own entries; owners and authorized editors may edit anyone's entries in the workspace.
- Normal corrections directly update the entry while preserving revisions and actors. Direct edits may change expense payer and participants and repayment sender/recipient, as well as amount, split, description, and occurrence date. Revalidate resulting membership and recording permissions and preserve the original creator. Corrections may occur after repayment and reverse debt direction; retain actual repayment records.
- Minimal cancellation voids an expense or repayment and reverses its balance effect atomically, preserving history, actor, and cancellation time. Use edit permissions; do not physically erase the record.
- Entries contain description and occurrence date; separately retain creation and edit times and actors. Permit backdating. Current balances include all nonvoided entries regardless of occurrence date.
- Use cursor pagination for history.

### Balances and consistency

- Maintain each member's net balance as writes occur instead of recomputing from history on reads. Preserve the invariant that workspace net balances sum to zero.
- Derive simplified debts from net balances, including removed members with nonzero balances. Preserve underlying expenses and repayments unchanged by simplification.
- Repeatedly match the largest remaining debtor and creditor, using a stable tie-breaker. Suggestions may change after any balance change; global minimization of payment count is not required.
- Commit each entry change, revision, balance update, and request-key record in one PostgreSQL transaction. Failed mutations must not leave partial balance or history changes.
- Successful writes guarantee subsequent reads reflect their changes, including repayment suggestions.
- Use optimistic concurrency checks on expense and repayment versions: an edit or cancellation conditioned on a stale version returns a conflict. Concurrent writes must not lose balance updates or allow repeated cancellation effects.
- Coordinate balance changes with database row locks scoped to workspace and member. Acquire affected balance rows in a consistent order, including members affected by either old or new entry values; a global user balance or blanket workspace mutation lock is not the selected strategy. Update all affected balances in the same transaction. Document and verify the full lock acquisition order and coordination with membership, permission, and currency changes while implementing the financial tickets; balance locks alone do not settle those races.
- Keep initial contention handling small. A complex automatic retry framework is not required; error and client retry semantics must preserve the idempotency contract.
- Creation, editing, cancellation, and repayment recording accept a client-provided `Idempotency-Key` scoped to requesting user and workspace. Retain them for the workspace lifetime.
- Enforce request-key uniqueness durably in PostgreSQL for actor and workspace. Claim the key, apply the financial effects, and store the operation fingerprint and replay result in the same database transaction. The fingerprint distinguishes operation, target, payload, and expected version where applicable. A rolled-back mutation leaves no committed effects or completed key reservation.
- Retry of the same operation has only one effect. Reusing a key for different data returns a conflict; identical expenses using different keys are distinct entries. Recheck current access before returning a stored response.

### Performance acceptance

- Target p95 balance-read latency below 200 ms and p95 write latency below 500 ms at 100 requests/second, with workspaces up to 100 members and 100,000 historical entries per tested workspace.
- Document hardware, database and service configuration, request mix, data distribution, duration, and errors with results. These are initial measurement targets, not demonstrated capacity or a hard product limit on workspace size.
- 10M daily active users is a future design target, not an MVP launch requirement.

## Testing Decisions

### Approved testing boundary

Use the running application's HTTP interface as the primary behavioral test boundary, with real application logic and PostgreSQL and a controlled local identity provider. Exercise public API requests and administration browser/session requests, including direct calls to administrative actions. Assert responses and subsequently observable balances, history, permissions, and revisions rather than private functions or table layouts. Exercise concurrent requests through the same boundary. Avoid mocked persistence for consistency tests. Use the real migration and provisioning commands for the behaviors only those entry points expose.

At synthesis, the checkout contains module configuration, setup documentation, specs, tickets, and accepted ADRs, but no application source or behavioral tests to reuse. README descriptions of absent implementation are not test evidence. The foundation establishes the approved HTTP/PostgreSQL test harness; feature tickets extend it.

Cover identity and access, workspace and membership management, expenses and splits, repayments, direct edits, cancellation, history, net balances, and simplified suggestions through this boundary. Add lower-level tests only where they materially improve mathematical edge-case coverage; they must not dictate internal structure.

Required behavioral scenarios:

- A pays $10 for B and C equally: B and C each owe $5; A has no share unless selected.
- B owes A $10 and A owes C $10: suggest B pays C $10, with A at zero and original history retained.
- From A owing B $10, independently record $5, $10, and $15 repayments: residual debt is $5, zero, and B owing A $5 respectively.
- Record a repayment outside current suggestions and when the sender owes nothing; verify updated balances and conservation of value.
- Cover all split methods, smallest-unit remainders, deterministic ties, invalid totals, and exact arithmetic.
- Correct a fully repaid $10 share to $6: the original creditor now owes $4 back; preserve repayment and revision history.
- Void an expense and a repayment; verify reversal, retained history, permissions, and protection against repeated effects.
- Concurrent edits from the same version: one succeeds and the other conflicts. Concurrent distinct mutations preserve all successful balance effects.
- Concurrent duplicate request keys produce one effect; changed payload with the same key conflicts; distinct keys preserve intentionally identical entries.
- Unauthorized mutations and failed writes leave balances and history unchanged. Immediate reads after success reflect committed changes.
- Exercise each delegated permission independently and verify owner transfer removes owner-only powers.
- Remove a zero-balance member, correct an old expense to give them a balance, and settle through an active party. Check removed-member access excludes history and other members' balances.
- Reject removal with nonzero net balance and new expense participation by removed members; verify rejoining retains history and resets delegated permissions.
- Verify requests and retries after access changes cannot disclose previously accessible transaction data.
- Verify currency locking, backdated entries, personal history filtering, and cursor pagination.
- Run the documented HTTP load test against realistic seeded history, including authentication and authorization checks. Report whether targets are met rather than inferring capacity from unit tests.
- Verify administration permissions and superuser access through observable administrative behavior. Administrative expense edits and cancellation must preserve balances, original creator, revisions, and concurrency guarantees; direct modification of balances or revisions is unavailable.
- Verify local migrations Up → Down → Up against representative data and document the data compatibility and recovery strategy for production release rollback.
- Exercise staff browser login/callback/session flows with the shared OIDC provider. Check permissions on subsequent requests and direct administrative actions, not just button visibility.
- Verify only superusers can manage staff privileges and permissions. Deactivation blocks API requests and existing admin sessions; provider login cannot silently reactivate a user. Explicit reactivation uses current permissions.
- Run superuser provisioning repeatedly for a specific provider identity and verify no duplicate users or grants to another identity.
- Fetch the served OpenAPI contract and Swagger UI through the existing HTTP/browser seam and exercise authorization through the controlled OIDC provider. Verify API requests/responses against the contract across success and relevant error cases, including idempotency and version requirements. Run contract validation and generated-Go freshness checks locally and in CI.
- Run the approved pinned formatter, lint analyzers, vulnerability checks, race-enabled tests, and PostgreSQL behavioral tests using matching local/CI commands. Measure performance independently of race instrumentation.

## Delivery and Story Coverage

`ready-for-agent` describes the ticket brief; work starts only when its blockers are complete. Foundation 01 is the first unblocked ticket. Foundation 03 gates MVP 01, and the existing MVP dependency graph then applies. Administrative acceptance is part of each relevant feature slice, not deferred to an unspecified future ticket.

| Ticket | Blocked by | User stories |
| --- | --- | --- |
| [Foundation 01: Runnable project and quality checks](../payment-share-foundation/issues/01-runnable-project-quality-checks.md) | None | 44, 49, 50, 53, 57 |
| [Foundation 02: Authentication and user lifecycle](../payment-share-foundation/issues/02-authentication-user-lifecycle.md) | Foundation 01 | 1, 44, 47, 52, 54, 55, 57 |
| [Foundation 03: Administration and permissions](../payment-share-foundation/issues/03-administration-permissions.md) | Foundation 02 | 45, 46, 47, 51, 52, 54, 56, 57 |
| [MVP 01: Create isolated workspaces](issues/01-authenticate-workspaces.md) | Foundation 03 | 2, 3, 4, 5, 38, 45, 46, 47, 56, 58, 59, 60 |
| [MVP 02: Add members and delegate permissions](issues/02-members-permissions.md) | MVP 01 | 6, 7, 8, 9, 38, 45, 46, 47, 56, 58, 59, 60 |
| [MVP 03: Record equal-split expenses and read balances](issues/03-equal-expenses-balances.md) | MVP 02 | 10, 11, 13, 14, 15, 16, 30, 39, 40, 42, 45, 46, 47, 48, 56, 58, 59, 60 |
| [MVP 04: Support exact, percentage, and weighted splits](issues/04-split-methods.md) | MVP 03 | 12, 13, 45, 46, 47, 48, 56, 58, 59, 60 |
| [MVP 05: Show simplified repayment suggestions](issues/05-simplified-debts.md) | MVP 03 | 17, 18, 42, 45, 46, 47, 56, 58, 59, 60 |
| [MVP 06: Record external repayments](issues/06-repayments.md) | MVP 03 | 19, 20, 21, 22, 30, 39, 40, 42, 45, 46, 47, 48, 56, 58, 59, 60 |
| [MVP 07: Edit and cancel entries with retained history](issues/07-edit-cancel.md) | MVP 04, MVP 06 | 23, 24, 25, 26, 27, 39, 40, 41, 42, 45, 46, 47, 48, 56, 58, 59, 60 |
| [MVP 08: Browse workspace and personal history](issues/08-history.md) | MVP 07 | 28, 29, 30, 45, 46, 47, 56, 58, 59, 60 |
| [MVP 09: Leave, remove, and rejoin workspaces](issues/09-membership-lifecycle.md) | MVP 05, MVP 08 | 31, 32, 33, 34, 35, 36, 37, 38, 45, 46, 47, 52, 56, 58, 59, 60 |
| [MVP 10: Verify MVP performance and complete workflows](issues/10-performance-workflows.md) | MVP 09 | 38, 43, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61 |

## Out of Scope

- A member-facing frontend or mobile application within this backend MVP. React application work starts after MVP ticket 10 is complete; detailed client design remains deferred. Swagger UI is included during foundation, and the administration panel is included in the first release.
- Moving real money, payment gateways, or counterparty approval of repayments.
- Currency conversion or currency per entry.
- Invitation acceptance, searchable user directories, workspace deletion, and archiving.
- Historical balance-as-of-date queries.
- Required void-and-replace correction workflows; these may be considered in a later version. Minimal cancellation remains in the MVP.
- Mathematically minimal repayment counts.
- Proven 10M DAU operation, distributed-service decomposition, and selection or provisioning of a hosted identity provider.

## Further Notes

- GORM + QOR5 + Goose, shared staff identity, application-owned administrative permissions, reversible local migrations, and the three-ticket foundation sequence are approved decisions. Integration acceptance tests establish whether the selected combination meets the contract.
- Direct editing with revision history and voiding remains the financial correction model, including administrative corrections and the approved changes to entry parties.
- Detailed endpoint names, response schemas, physical schema, aggregate interfaces, currency configuration, session settings, and the complete lock protocol are implementation choices within these requirements. They must be documented and verified in the relevant tickets; do not interpret them as approval to weaken consistency or restart the interview.
- Initial quality/performance targets remain measurement requirements. Broader infrastructure and operational requirements can be discussed later.
- Accepted rationale is recorded in the [ADRs](../../docs/adr/); alternatives are in the [implementation research](implementation-research.md).
