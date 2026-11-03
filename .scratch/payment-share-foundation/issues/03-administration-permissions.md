# 03: Administration and permissions

Status: ready-for-agent

**What to build:** Operate foundation users and application records through QOR5 with restricted staff permissions, unrestricted superuser permissions, and safe user lifecycle management.

**Blocked by:** [Foundation 02: Authentication and user lifecycle](02-authentication-user-lifecycle.md).

**Parent:** [Payment Share project foundation](../spec.md)

**Spec user stories:** 7, 8, 10, 11, 12, 13, 18, 19, 20, 23, 24.

**MVP user stories:** 45, 46, 47, 51, 52, 54, 56, 57 in the [MVP spec](../../payment-share-mvp/spec.md).

## Implementation constraints

Mount administration in the same initial application deployment and reuse the preceding identity/session and lifecycle use cases. Use QOR5 with GORM and Goose. Domain rules and authorization must hold for direct server requests as well as visible UI actions. Application feature tables are registered when later feature tickets introduce them; no workspace/financial feature implementation is required to finish this foundation.

## Acceptance criteria

- [ ] Update OpenAPI for any application API operations introduced or changed here, including inputs, success responses, relevant errors, authentication, and required headers. Keep QOR5 internal browser routes outside that contract. Verify application API requests/responses against the contract and pass contract validation and generated-code freshness checks; retain the separate administration/session behavior tests.
- [ ] Deliver a working browser administration panel through the local OIDC login/session flow. Ordinary users cannot access it; authenticated staff see only allowed resources/actions.
- [ ] Register foundation application tables and users for administrative browsing, protecting sensitive operational fields. Demonstrate list/detail and applicable create/update/delete operations with persisted data.
- [ ] Persist and manage staff roles/permissions with separate resource read, create, update, and delete actions. Check authorization on the server and apply changes on subsequent requests, including existing sessions.
- [ ] Superusers access every registered administrative resource without individual permission grants. Verify denial-policy interactions do not accidentally remove that authority; domain invariants and retained history still apply.
- [ ] Only superusers grant/revoke staff status, alter administrative permissions, and promote superusers. Restricted staff cannot escalate through direct user, role, permission, or bulk-action requests.
- [ ] Connect user deactivation/reactivation screens to foundation 02's use cases. Deactivation blocks existing credentials/sessions without deleting identity or history; explicit reactivation uses current permissions.
- [ ] Verify and document custom save/delete integration so later financial administrative operations invoke application use cases. Generic writes must not become a route to directly edit stored balances or revisions; financial deletion later means voiding.
- [ ] Browser/session and direct-request tests with real PostgreSQL prove ordinary-user denial, restricted staff, permission changes, superuser authority, privilege-management restrictions, and deactivation/reactivation behavior.
- [ ] Run approved formatting, linting, vulnerability, race-enabled test, and migration checks. Document working local setup and superuser provisioning; completion unblocks the existing MVP feature sequence.
