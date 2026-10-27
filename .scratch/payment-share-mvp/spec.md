# Payment Share MVP backend

Status: ready-for-agent

Readiness: Content synthesized from the agreed requirements. The user confirmed the HTTP API testing boundary with real PostgreSQL and a controlled local identity provider.

## Problem Statement

People sharing expenses need to know what they owe across separate groups without manually reconstructing transaction history. Corrections, repayments, concurrent requests, and changes in membership must not leave balances incorrect or expose another workspace's records.

## Solution

Provide a backend API for workspaces with one currency each, shared expenses, externally made repayments, direct corrections with revision history, and cancellation of mistaken entries. Show current net balances and simplified repayment suggestions immediately after changes. Enforce workspace permissions and retain limited balance access for removed members.

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

## Implementation Decisions

### Service and identity

- Use the existing Go module and Chi router with PostgreSQL, initially as one backend API service.
- Use a configurable OpenID Connect provider for registration and login. Configure issuer and audience, validate access tokens locally, and check current workspace membership and permissions on each request.
- Provide a local development identity provider. Hosted provider selection is a deployment decision.
- Add already registered users immediately by exact verified email; provide no searchable directory or invitation acceptance flow.

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
- Normal corrections directly update the entry while preserving revisions and actors. Corrections may occur after repayment and reverse debt direction; retain actual repayment records.
- Minimal cancellation voids an expense or repayment and reverses its balance effect atomically, preserving history, actor, and cancellation time. Use edit permissions; do not physically erase the record.
- Entries contain description and occurrence date; separately retain creation and edit times and actors. Permit backdating. Current balances include all nonvoided entries regardless of occurrence date.
- Use cursor pagination for history.

### Balances and consistency

- Maintain each member's net balance as writes occur instead of recomputing from history on reads. Preserve the invariant that workspace net balances sum to zero.
- Derive simplified debts from net balances, including removed members with nonzero balances. Preserve underlying expenses and repayments unchanged by simplification.
- Repeatedly match the largest remaining debtor and creditor, using a stable tie-breaker. Suggestions may change after any balance change; global minimization of payment count is not required.
- Commit each entry change, revision, balance update, and request-key record in one PostgreSQL transaction. Failed mutations must not leave partial balance or history changes.
- Successful writes guarantee subsequent reads reflect their changes, including repayment suggestions.
- Reject stale edits with a conflict using the entry version supplied by the client. Concurrent writes must not lose balance updates or allow repeated cancellation effects.
- Creation, editing, cancellation, and repayment recording accept client-provided request keys scoped to requesting user and workspace. Retain them for the workspace lifetime.
- Retry of the same operation has only one effect. Reusing a key for different data returns a conflict; identical expenses using different keys are distinct entries. Recheck current access before returning a stored response.

### Performance acceptance

- Target p95 balance-read latency below 200 ms and p95 write latency below 500 ms at 100 requests/second, with workspaces up to 100 members and 100,000 historical entries per tested workspace.
- Document hardware, database and service configuration, request mix, data distribution, duration, and errors with results. These are initial measurement targets, not demonstrated capacity or a hard product limit on workspace size.
- 10M daily active users is a future design target, not an MVP launch requirement.

## Testing Decisions

### Approved testing boundary

Use the public HTTP API as the primary behavioral test boundary, running the real service logic and PostgreSQL with a controlled local identity provider. Assert responses and subsequently observable balances, history, permissions, and revisions rather than private functions or table layouts. Exercise concurrent requests through the same boundary. Avoid mocked persistence for consistency tests.

The repository currently contains module configuration and requirements documentation, with no application implementation, tests, or ADRs to reuse. This is a new testing boundary, not a replacement for existing tests.

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

## Out of Scope

- A frontend or mobile application.
- Moving real money, payment gateways, or counterparty approval of repayments.
- Currency conversion or currency per entry.
- Invitation acceptance, searchable user directories, workspace deletion, and archiving.
- Historical balance-as-of-date queries.
- Required void-and-replace correction workflows; these may be considered in a later version. Minimal cancellation remains in the MVP.
- Mathematically minimal repayment counts.
- Proven 10M DAU operation, distributed-service decomposition, and selection or provisioning of a hosted identity provider.

## Further Notes

- This spec synthesizes the approved requirements and later corrections; the latest decision is direct editing with revision history plus minimal cancellation.
- Detailed endpoint naming, response schemas, schema design, locking strategy, and supported currency configuration remain implementation design work within these behavioral constraints.
- The conversation did not settle whether direct edits can change an expense payer or repayment parties. Do not treat the earlier void-and-recreate proposal as approved; resolve that field-level contract before implementing those edits.
- The user approved the testing boundary. No implementation or ticket breakdown is included in this task.
