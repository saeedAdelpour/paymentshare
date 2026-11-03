# 02: Authentication and user lifecycle

Status: ready-for-agent

**What to build:** Authenticate through a local OIDC provider, persist stable users, provision the initial superuser, and enforce deactivation/reactivation before adding the administration UI.

**Blocked by:** [Foundation 01: Runnable project and quality checks](01-runnable-project-quality-checks.md).

**Parent:** [Payment Share project foundation](../spec.md)

**Spec user stories:** 5, 6, 9, 11, 12, 17, 18, 19, 20, 23, 24.

**MVP user stories:** 1, 44, 47, 52, 54, 55, 57 in the [MVP spec](../../payment-share-mvp/spec.md).

## Implementation constraints

Use the foundation's GORM/PostgreSQL persistence and Goose migrations. Keep OIDC and persistence outside the domain. Staff and ordinary users share the provider; staff status, superuser status, roles, and permissions belong to Payment Share. Expose the minimal server-side user-lifecycle operations needed to verify behavior; foundation 03 adds their administration screens. Preserve financial and actor references as later feature tables are introduced.

## Acceptance criteria

- [ ] Make Swagger UI's Authorize action work through the configured local OIDC provider using Authorization Code with PKCE and API access tokens, separately from the QOR5 session flow. Document reproducible client/redirect configuration and demonstrate authorized API requests with normal server-side access enforcement.
- [ ] Extend OpenAPI for every application API operation delivered here, including inputs, success responses, relevant errors, authentication, and required headers. HTTP behavior tests verify requests/responses against the contract, including rejected credentials and deactivated-user access; pass contract validation and generated-code freshness checks.
- [ ] Run a reproducible local identity provider with configurable issuer and audience. Validate API access tokens locally; reject expired, invalid, wrong-issuer, and wrong-audience tokens and ID tokens presented as API credentials.
- [ ] Persist stable identity using provider issuer/subject and verified email for later membership operations. Repeated and concurrent first authenticated requests create one identity; restart preserves it. Ordinary registration grants no staff or superuser privileges.
- [ ] Provide a browser login/session flow through the same provider suitable for the QOR5 integration. Document session behavior and enforce current application access on authenticated requests; exercise a real login/callback/session journey rather than replacing it with a test-only authentication bypass.
- [ ] Provide an idempotent CLI command to provision an initial superuser for a specific provider identity. Repeating it creates no duplicate user and grants no privileges to a different identity; document the prerequisite registration/provisioning workflow.
- [ ] Implement deactivation/reactivation through authorized user-lifecycle use cases. Deactivation blocks API and administration access with valid existing tokens/sessions, preserves the user record and references, and leaves the external provider account intact.
- [ ] Login or a first-request identity sync cannot silently reactivate a deactivated user. Explicit reactivation restores access subject to current permissions. Protect lifecycle operations with server-side administrative authorization.
- [ ] HTTP/session tests using real PostgreSQL and a controlled identity provider demonstrate persistence, rejected credentials, ordinary-user isolation from administration, superuser provisioning, and access before/after deactivation and reactivation.
- [ ] Include reversible migrations for introduced schema and run the approved quality checks.
