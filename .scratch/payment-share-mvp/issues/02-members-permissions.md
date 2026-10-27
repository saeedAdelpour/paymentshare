# 02: Add members and delegate permissions

Status: ready-for-agent

**What to build:** Add registered members, delegate independent permissions, and transfer ownership.

**Blocked by:** 01: Authenticate and create workspaces.

**Parent:** [Payment Share MVP backend](../spec.md)

**Spec user stories:** 6, 7, 8, 9, 38.

## Implementation constraints

Use DDD within the existing single Payment Share bounded context. Domain entities and value objects own money, split, membership, and entry invariants; application use cases coordinate authorization and database transactions. Keep HTTP, OIDC, and PostgreSQL in adapters outside the domain. Choose aggregate boundaries around invariants, without loading all workspace history as one object graph. DDD does not require microservices, event sourcing, or asynchronous balances.

Implement a complete backend slice, including durable persistence, API behavior, and tests through HTTP with real PostgreSQL and the controlled local identity provider. No UI is in scope. Follow the parent spec except for the explicitly approved edit-party clarification in ticket 07; the parent is preserved unchanged during ticket publication.

## Acceptance criteria

- [ ] Owner or a member with add permission adds an existing registered user immediately by exact verified email, without exposing a searchable directory.
- [ ] Only owners grant or revoke independent add-member, remove-member, and edit-others permissions. An add permission does not imply other powers.
- [ ] Active members see their workspace; outsiders remain excluded, including attempts using another workspace's identifiers.
- [ ] Transfer ownership to an active member atomically, preserving exactly one owner and removing the former owner's owner-only powers.
- [ ] Exercise each permission independently through HTTP. Revocation takes effect on subsequent requests, not only after token expiry.

