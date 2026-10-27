# A payment share (like splitwise) app

# functional requirements

* users can add transaction containing other users
for example: user A pays $10 for B and C equally, so B and C owe $5 each to user A. The payer is included in the split only when explicitly selected.
* user can exists in multiple workspaces of payments and the permissions of seeing workspace must be considered
* users can record repayments made outside the app; the app does not transfer money
* show simplified debts using workspace net balances: if B owes A $10 and A owes C $10, show B owing C $10. Simplification changes repayment suggestions only; preserve original expenses and actual repayments
* choose repayment suggestions by matching the largest outstanding debtor with the largest creditor, using a stable tie-breaker. Suggestions may change when balances change; a mathematically minimal number of payments is not required
* repayments may be less than, equal to, or greater than the debt. If A owes B $10 and pays B $15, B then owes A $5 in a workspace containing only these debts
* ordinary members can record expenses only when they are the payer; owners and authorized editors can record expenses on behalf of other members. Keep creator and payer distinct
* either repayment party, the owner, or an authorized editor can record a repayment. It affects balances immediately without confirmation from the other party; preserve who recorded it
* repayments can be recorded between any two active workspace members, regardless of suggested recipients or whether the sender currently owes anything
* expense edits remain possible after repayment, including edits that reverse who owes whom. Preserve repayments and correction history, and update balances atomically with edits
* MVP corrections directly update the existing expense or repayment while preserving revision history. Do not require voiding and creating a replacement for normal corrections
* removal preserves entries and restricts access to a read-only view of the removed member's own balance and repayment suggestions involving them, including who to pay or receive money from. No transaction history, other members' balances, or write access. Their workspace net balance must be zero at removal, even if original expenses linked them to other members
* later corrections may change a removed member's balance without re-adding them. Include removed members with nonzero balances in simplified repayment suggestions
* current members can record repayments involving a removed member when they are the other repayment party; owners and authorized editors can also record these repayments
* removed members remain eligible for corrections to existing expenses and for repayments, but must be re-added to participate in new expenses
* MVP includes minimal cancellation (voiding) for mistaken or duplicate expenses and repayments, using the same permissions as editing. Atomically reverse the balance effect and preserve the entry and history, including who cancelled it and when. Cancellation is distinct from normal direct edits
* currency handled in each workspace, not per each transaction
* workspace currency becomes fixed when the first expense or repayment is recorded
* active workspace members can see all expenses and repayments in their workspace, with a personal view filtered to entries involving them. Removed members retain the limited access described above; users who never joined have no access
* each workspace has an owner, and it can grant access to other users, such as:
  * add another users
  * remove user
  * edit transaction
* creators can edit their own entries; owners and explicitly authorized members can edit anyone's entries in their workspace. Preserve history showing who changed what.
* adding members requires existing registered users. Adding members, removing members, and editing others' entries are separate permissions; only the owner can grant them
* add registered users immediately by exact verified email address; do not expose a searchable user directory. Invitation acceptance is outside the MVP
* re-adding a removed member reconnects their existing history and balance and restores normal workspace visibility. Delegated permissions start empty and must be explicitly granted again
* each workspace has exactly one owner, who must transfer ownership to another active member before leaving. Transfer removes the former owner's owner-only powers
* support workspace renaming and ownership transfer. Members can leave under the zero-balance rule; workspace deletion and archiving are outside the MVP
* support equal splits, exact amounts, percentages, and relative weights such as 2:1. Exact amounts must sum to the total and percentages to 100%. Distribute leftover smallest currency units deterministically for calculated splits so shares sum exactly to the total
* amounts use each currency's smallest unit; expense and repayment totals must be positive. No currency conversion. Allocate calculated split remainders by largest fractional remainder, breaking ties by member ID
* expenses and repayments contain a description and occurrence date. Record creation time, edit time, and actor separately; allow backdated entries. History uses cursor pagination. Balances include all current nonvoided entries regardless of occurrence date; historical balance-as-of-date queries are outside the first release

# non-functional requirements
* user must see the aggregation balance very fast, it's not good to calculate user balance on each request
* 10M DAU is a future design target. The first release must establish correct balances, safe concurrent writes, and measurable performance, including authentication.
* high consistency on adding transaction
* transactions must be prune to race conditions and idempotency
* reject stale edits with a conflict when another edit has changed the entry since the supplied version
* creating, editing, voiding, and recording repayments accept client-provided request keys. Retrying the same operation must not apply it twice; reusing a key with different data returns a conflict. Identical expenses submitted with different keys are separate expenses
* request keys are scoped to the requesting user and workspace and retained for the workspace's lifetime. Check current access permissions before returning a stored response to a retry
* users will edit transaction, so balance showing should be consistent
* successful writes guarantee subsequent balance reads reflect the change, including simplified repayment suggestions. Maintain net balances on writes; derive suggestions from net balances rather than transaction history
* initial load-test targets: 95% of balance reads under 200 ms and writes under 500 ms at 100 requests/second, workspaces up to 100 members, and 100,000 historical entries per tested workspace. Document the load-test conditions; these targets do not establish 10M DAU capacity

# core entities
user
transaction
workspace

# agreed delivery scope

The first release is an MVP backend API with automated tests. Further enhancements can be specified in later versions rather than expanding this release.
Domain terminology is defined in [GLOSSARY.md](GLOSSARY.md).

# api design

Use an external identity provider for registration and login, validating its access tokens locally. Check current workspace membership and permissions on each request.
Accept a configurable OpenID Connect provider with issuer and audience supplied through configuration. Include a local development provider so authentication and API tests can run without a paid service; hosted provider selection is a deployment decision.

# high level design

Build on Go and Chi with PostgreSQL, initially as one API service. Commit each entry change, balance update, and request-key record in the same database transaction.

# Possible later versions

* A void-and-replace correction workflow may be considered later. This is not a committed requirement and does not replace the MVP's direct updates or minimal cancellation operation.
* Invitation acceptance, workspace archiving or deletion, and historical balance-as-of-date queries are outside the MVP and can be specified separately if needed.
* 10M DAU remains a future design target, not the first release's capacity claim.
