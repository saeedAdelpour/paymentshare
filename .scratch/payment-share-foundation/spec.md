# Payment Share project foundation

Status: ready-for-agent

The user approved GORM + QOR5 + Goose and the three-ticket foundation sequence on 2026-11-03. All existing MVP feature tickets follow completion of this foundation. Read the [MVP spec](../payment-share-mvp/spec.md), [glossary](../../GLOSSARY.md), and applicable [ADRs](../../docs/adr/) for the shared domain and accepted decisions.

## Problem Statement

Feature development needs a working, reproducible application foundation with identity, administration, controlled schema evolution, and quality checks. Developers also need to explore and execute the API before the member-facing application exists, with a reliable contract for that later client. Documentation describing absent code is not evidence that those capabilities work.

## Solution

Deliver the Go/Chi application with PostgreSQL persistence using GORM, versioned SQL migrations managed by Goose, and a QOR5 administration panel. Authenticate application users and staff through the same configurable OIDC provider. Persist administrative permissions in Payment Share and make the foundation observable through working HTTP and browser flows. Provide Swagger UI backed by an explicit OpenAPI contract and generated Go transport interfaces, including real OIDC authorization. Keep the contract aligned with API behavior throughout MVP delivery; build the React member-facing application after MVP ticket 10.

## User Stories

1. As a developer, I want documented local setup and start/check commands, so that a fresh checkout is reproducible.
2. As a developer, I want pinned formatting, linting, testing, and vulnerability tools, so that local and CI checks agree.
3. As a developer, I want target-version Up/Down migrations, so that I can switch application versions locally.
4. As an operator, I want explicit release migration and rollback/recovery procedures, so that deployment preserves data expectations.
5. As a user, I want to authenticate through an identity provider, so that my identity is stable across application restarts.
6. As a staff user, I want to sign into administration through that provider, so that staff identity is not a separate password system.
7. As a staff user, I want access only to permitted resources and actions, so that my responsibilities determine my access.
8. As a superuser, I want all administrative permissions without individual grants, so that I can administer every registered application resource.
9. As an operator, I want repeatable CLI superuser provisioning for a specific provider identity, so that initial administration is reproducible.
10. As a superuser, I want exclusive authority to manage staff status, administrative permissions, and superuser promotion, so that staff cannot escalate their own access.
11. As an authorized administrator, I want to deactivate users while retaining references and history, so that removing access does not destroy shared records.
12. As an authorized administrator, I want explicit reactivation to use current permissions, so that authentication alone cannot undo deactivation.
13. As a service administrator, I want a working panel for foundation tables and users that can grow with feature tables, so that application records are operable.

14. As a developer, I want Swagger UI during foundation, so that I can explore and execute API operations before the member-facing application exists.
15. As an API consumer, I want an explicit OpenAPI contract as the source of truth, so that requests and responses have a clear integration boundary.
16. As a developer, I want Go transport types and Chi server interfaces generated from the contract, so that handlers follow the declared API shape.
17. As a developer, I want to authorize Swagger UI through the configured OIDC provider using Authorization Code with PKCE, so that I can exercise authenticated operations without manually obtaining tokens.
18. As a user, I want Swagger requests subject to normal API authentication and permissions, so that the interactive documentation grants no additional access.
19. As an API consumer, I want every implemented application API operation to document inputs, success responses, relevant errors, authentication, and required headers, so that I can integrate without guessing.
20. As a developer, I want each backend slice to update its contract and Swagger documentation, so that API documentation grows with implemented behavior.
21. As an operator, I want Swagger UI enabled locally and disabled by default elsewhere with explicit configuration, so that its availability is intentional.
22. As a developer, I want local and CI checks to reject invalid contracts and stale generated Go code, so that inconsistent artifacts cannot silently pass checks.
23. As an API consumer, I want HTTP behavior tests to verify requests and responses against the contract, so that the running API matches its published interface.
24. As a developer, I want QOR5 browser routes and identity-provider endpoints kept outside the application API contract, so that the documented integration surface is clear.
25. As a developer, I want the React member-facing application to follow completed MVP backend validation, so that client work starts against the verified API.

## Implementation Decisions

- Use an explicit OpenAPI contract as the HTTP API source of truth, with oapi-codegen generating Go request/response types and Chi server interfaces and Swagger UI providing interactive API documentation. Keep generated transport types in the HTTP adapter. Deliver Swagger UI during foundation and update the contract alongside each MVP API ticket. Enable local development access; other deployments require explicit configuration and default to disabled. Swagger requests use normal API authentication and permissions. See [ADR 0007](../../docs/adr/0007-openapi-contract-first-http-api.md).
- Foundation 01 establishes the OpenAPI contract, Swagger UI, and pinned contract/code-generation checks. Foundation 02 adds Swagger UI authorization through the configured OIDC provider using Authorization Code with PKCE and API access tokens, separately from QOR5 sessions. Document every application API operation's inputs, success responses, relevant errors, authentication, and required headers; exclude QOR5 internal browser routes and identity-provider endpoints. Local checks and CI reject invalid contracts and stale generated Go code. HTTP behavior tests verify requests and responses against the contract, extending coverage with each API slice.
- Use Go, Chi, PostgreSQL, GORM, QOR5, and Goose. Pin compatible dependency and tool releases; verify the selected combination through integration behavior. Report an incompatibility with evidence before substituting a different stack.
- Keep domain rules independent of GORM, HTTP, OIDC, and QOR5. Application use cases coordinate authorization and transactions; adapters integrate external systems. Future API and administrative financial writes use the same application use cases.
- Use reviewed versioned SQL migrations. Verify Up → Down → Up with representative data and document data-preservation limits. Production rollback preserves previous-release compatibility where practical and requires explicit recovery for destructive changes. Execute migrations once as a release step; startup verifies schema compatibility.
- Use a reproducible local OIDC provider and a configurable issuer/audience. Validate API access tokens locally and reject ID tokens used as API credentials. Persist stable provider identity; ordinary registration grants no staff privileges.
- Staff browser login uses the same provider. Implement a browser-appropriate authorization flow and session handling; API bearer-token acceptance is not itself browser login. Document the chosen session configuration and enforce current application access on requests.
- Store staff status, superuser status, roles, and administrative permissions in Payment Share. Permissions distinguish read, create, update, and delete actions per application resource. Superusers need no individual resource grants; server enforcement must agree with the UI.
- Only superusers grant/revoke staff status, alter administrative permissions, or promote superusers initially. Provide repeatable, idempotent CLI provisioning tied to a specific provider identity.
- Deactivation blocks API and administration access, including existing tokens/sessions, and retains identity, memberships, financial references, and history. The external identity-provider account remains intact. Explicit reactivation restores access subject to current permissions.
- Register the application tables introduced by the foundation in administration; add feature tables when they exist. Protect sensitive operational fields from disclosure. Generic administrative writes must not bypass application use cases or mutate balances/revisions directly.
- Use gofumpt; golangci-lint v2 with govet, staticcheck, errcheck, ineffassign, and unused; govulncheck; and race-enabled Go tests. Pin versions and provide matching local/CI commands. Real PostgreSQL behavioral tests cover database consistency; the Go race detector covers memory races.

## Testing Decisions

The approved seams are the running application's HTTP interface (public API and administration browser/session requests), plus the real migration and superuser-provisioning commands for their command-only behavior. Use the same real PostgreSQL and controlled identity-provider setup across these checks; introduce no mock persistence seam. At synthesis, the checkout has setup/configuration documents and accepted ADRs but no application source or test suite to reuse. Foundation 01 establishes the harness, 02 adds identity/session behavior, and 03 extends it to the administration panel.

- Use the approved running HTTP/browser seam to fetch the served OpenAPI document and Swagger UI, exercise actual authorization through the controlled provider, and send authenticated API requests. Verify local availability and the default disabled behavior in other deployment configurations through observable HTTP behavior.
- Verify application API requests and responses against the contract, including success and relevant error status codes, response bodies, and required request headers. Extend the same real PostgreSQL/identity-provider harness per API slice; avoid tests of generated implementation internals.
- Run the real contract-validation and generation-check commands locally and in CI. Demonstrate rejection of an invalid contract and detection of stale generated output; pin generation/validation tooling and document reproducible commands.
- Exercise public HTTP behavior with real PostgreSQL and a controlled identity provider. Observe identity persistence, authorization failures, and access changes rather than private table layouts.
- Exercise browser login, session-backed administration, restricted staff, and superuser flows with browser-level or equivalent HTTP session tests, including server-side denials on direct requests.
- Verify migrations against a real database through target-version Up → Down → Up with representative fixtures and observable retained data where reversibility promises it. Document destructive limits rather than treating schema recreation as data recovery.
- Test repeatable superuser provisioning without duplicate users and without granting privileges to another provider identity.
- Verify deactivation with valid existing credentials/sessions, failed implicit reactivation at login, explicit reactivation, and permission changes taking effect on subsequent requests.
- Run the approved checks locally and in CI. Verify compatibility and behavior; no performance capacity claim follows from selecting an ORM or passing these foundation checks.

## Delivery Graph

1. [01: Runnable project and quality checks](issues/01-runnable-project-quality-checks.md) — no blockers; includes the contract, Swagger UI, generation, and contract checks.
2. [02: Authentication and user lifecycle](issues/02-authentication-user-lifecycle.md) — blocked by foundation 01; adds Swagger OIDC authorization and authenticated contract behavior checks.
3. [03: Administration and permissions](issues/03-administration-permissions.md) — blocked by foundation 02.

Completion of foundation 03 unblocks [MVP 01](../payment-share-mvp/issues/01-authenticate-workspaces.md); its existing descendants retain their dependencies.

## Out of Scope

- Workspace, expense, repayment, balance, and membership feature implementation; these remain in the existing MVP tickets.
- A member-facing frontend/mobile application, hosted identity-provider provisioning, and uninterrupted rolling deployment requirements.
- The member-facing application will use React and starts after [MVP ticket 10](../payment-share-mvp/issues/10-performance-workflows.md) is complete. Its detailed design is deferred until then; Swagger UI is delivered during foundation and QOR5 remains the staff administration interface.
- Exhaustive infrastructure design, proven production capacity, and the detailed financial lock protocol. The existing financial integrity and idempotency requirements remain binding for later feature work.

## Further Notes

The user confirmed the complete Swagger/OpenAPI plan on 2026-11-03, including the HTTP/browser and command testing seams. This specification incorporates that discussion into the existing three-ticket foundation rather than introducing a parallel delivery track.

The [research](../payment-share-mvp/implementation-research.md) records sources and alternatives. The [foundation discussion](../payment-share-mvp/foundation-plan.md) records the approved sequence and acceptance expectations. Browser/session settings and CLI ergonomics are implementation choices within this contract, not reasons to restart the requirements interview.
