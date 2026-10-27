# 01: Authenticate and create workspaces

Status: ready-for-agent

**What to build:** Sign in and create, list, and rename isolated workspaces with a chosen currency and one owner.

**Blocked by:** None (can start immediately).

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 1, 2, 3, 4, 5, 38, 44.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. No UI is in scope. Follow the parent spec except for the explicitly approved edit-party clarification in ticket 07; the parent is preserved unchanged during ticket publication.

## Acceptance criteria

- [ ] Run the Go/Chi API with PostgreSQL and reproducible local setup, migrations, and HTTP integration tests.
- [ ] Use configurable OIDC issuer and audience, a local development provider, and locally validated API access tokens. Reject invalid, expired, wrong-issuer, and wrong-audience tokens; do not substitute an ID token for an API access token.
- [ ] Recognize registered users by stable provider identity and obtain verified email for later membership operations. Keep identity-provider details outside the domain model.
- [ ] Persist workspace ownership and currency; users can own multiple workspaces and cannot access another user's workspace.
- [ ] Owner can rename a workspace; enforce authorization on every request. Hosted identity-provider selection remains deployment work.

