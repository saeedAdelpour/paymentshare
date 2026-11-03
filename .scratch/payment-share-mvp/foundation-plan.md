# Project foundation before feature tickets

The user approved this delivery order, the GORM + QOR5 + Goose stack, and the three-ticket breakdown on 2026-11-03. The [foundation spec](../payment-share-foundation/spec.md) and its linked tickets are the delivery contract. Existing MVP feature tickets start after foundation 03 is complete.

## Required outcome

Run the local Go application with PostgreSQL, authentication, an administration panel, migration management, and code quality checks. Staff permissions and superuser behavior work before workspace/expense/repayment feature work starts. The panel grows to cover feature tables as they are implemented.

## Approved API documentation direction

On 2026-11-03, the user accepted contract-first OpenAPI with oapi-codegen-generated Go types and Chi server interfaces, displayed through Swagger UI during foundation. Keep transport types outside the domain and update the contract with each MVP API ticket. Swagger UI is enabled for local development; other deployments require explicit configuration and default to disabled. Its requests use normal API authentication and permissions. See [ADR 0007](../../docs/adr/0007-openapi-contract-first-http-api.md).

The member-facing application will use React after MVP ticket 10 is complete; detailed React design is deferred until then. QOR5 remains the separate staff administration panel.

The user also accepted Swagger UI authorization through OIDC Authorization Code with PKCE and API access tokens, documentation of every application API operation (inputs, success responses, relevant errors, authentication, and required headers), and local/CI checks rejecting invalid OpenAPI or stale generated Go code. HTTP behavior tests verify requests and responses against the contract. Exclude QOR5 internal browser routes and identity-provider endpoints. Foundation 01 delivers the contract, Swagger page, and checks; foundation 02 adds working authorization; later slices extend documentation and behavioral coverage. The user confirmed shared understanding of the complete plan on 2026-11-03, closing the design interview. The foundation and MVP tickets incorporate these accepted requirements.

## Approved implementation sequence

1. **Runnable application and checks.** Establish application composition and local PostgreSQL setup, reviewed versioned migrations, formatting/linting, and reproducible check commands. Demonstrate a clean database initialization and migration Up → Down → Up using representative seed data.
2. **Authenticated identity.** Run a local identity provider and exercise the agreed access-token contract through HTTP with real PostgreSQL. Persist stable user identity and verify rejected invalid tokens. Provide idempotent CLI superuser provisioning for a specific provider identity. Staff use the same OIDC provider as application users; Payment Share stores staff status, superuser status, roles, and administrative permissions. Implement browser login/session handling for the administration panel; detailed session configuration remains implementation work.
3. **Working administration.** Deliver user/staff views and permission management, reusing the idempotent CLI superuser provisioning delivered in step 2. Ordinary registration grants no staff privileges. Demonstrate an ordinary user's denied access, restricted staff access, permission changes, and superuser access without individual grants. Register new application resources as feature work proceeds.

GORM + QOR5 + Goose is the selected implementation stack. Verify staff identity/session integration, server-side permission enforcement, superuser policy behavior, and custom mutation handlers before treating it as proven. Financial mutations later use the same application use cases as the public API and preserve balances and revisions.

## Acceptance evidence

- A fresh local checkout has documented setup and start/check commands.
- Migrations move to a chosen version in either direction locally; data-preservation limits and production recovery expectations are documented.
- Local authentication and persistent user identity work through the real HTTP interface and PostgreSQL.
- A staff user can access only permitted resources/actions. An ordinary application user cannot access administration. Superuser access does not depend on individual resource grants. Only superusers can change staff status, administrative permissions, or superuser promotion. Deactivation blocks access and retains history; reactivation uses current permissions.
- The administration panel runs as part of the initial application deployment. Migration execution is an explicit release step before application startup; startup verifies schema compatibility.
- Run the approved gofumpt, golangci-lint v2, govulncheck, and race-enabled Go tests through pinned, reproducible tools and matching local/CI commands. Keep PostgreSQL behavioral integration tests for persistence and concurrency.
- The implementation includes behavior tests for identity, administrative access, and migrations; application quality claims rely on these checks rather than scaffolding alone.

## Implementation validation

- Implement and verify staff browser/session configuration and CLI ergonomics within the accepted shared OIDC, application-owned permission, and CLI provisioning contract.
- Integration behavior of the selected persistence/administration/migration combination must pass the foundation acceptance checks.
- The approved [three-ticket breakdown](../payment-share-foundation/spec.md#delivery-graph) is published. Ticket 01 is unblocked; 02 waits for 01; 03 waits for 02.
- Protect sensitive operational fields and document the session configuration during implementation. Core permission-management and deactivation policies are accepted below.

## Approved staff and user management policies

- Only superusers grant/revoke staff status, change administrative permissions, and promote other superusers initially. Restrict role and permission management endpoints as well as their UI controls.
- Initial administration offers user deactivation rather than physical deletion. Deactivated users cannot access the API or administration even with valid provider tokens or existing sessions. Keep memberships, financial records, and actor history.
- The external identity-provider account remains intact. Authentication must not silently reactivate a deactivated user. Explicit reactivation restores access subject to current permissions.
- Demonstrate these policies through behavior tests, including direct unauthorized permission-management requests and requests made after deactivation.

## Approved quality tooling

Use gofumpt as the single formatter; golangci-lint v2 with an explicit initial set of govet, staticcheck, errcheck, ineffassign, and unused; govulncheck; and race-enabled Go tests. Keep real PostgreSQL behavioral concurrency tests because the Go race detector does not detect database transaction races. Pin tool versions and use the same commands locally and in CI. The user accepted this baseline on 2026-11-03. The tooling sources and installation caveats are captured in [implementation-research.md](implementation-research.md).

## Financial concurrency direction for later feature work

The user prefers optimistic entry versions with database locks on workspace-member balances and durable request-key idempotency. This replaces the proposed blanket workspace mutation lock. Balance locks must be coordinated with membership, permissions, and currency changes before the financial tickets are implemented. Broader infrastructure and non-functional requirements can be revisited later, as requested by the user.
