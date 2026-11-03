# Payment Share API

Go/Chi backend with PostgreSQL persistence and OIDC access-token authentication.
This slice supports registered users and creating, listing, reading, and renaming
workspaces. Membership management and expense/repayment features are later tickets.

## Local setup

Install Go 1.27+, Docker with Compose, curl, and jq. From the repository root:

```sh
cp .env.example .env
set -a
. ./.env
set +a
make dev
```

PostgreSQL listens on `127.0.0.1:55432`; Keycloak listens on
`http://localhost:18080`. Wait for Keycloak's realm import to finish:

```sh
until curl -fsS "$OIDC_ISSUER/.well-known/openid-configuration" >/dev/null; do
  sleep 2
done
make run
```

The API listens on `127.0.0.1:8082`. Startup applies embedded, ordered migrations
in one transaction under a database advisory lock. `make migrate` applies them
without starting the API. Repeated runs preserve existing data. Compose uses
named volumes for both services; `docker compose down` preserves them.

The development realm has Alice (`alice` / `alice-local`) and Bob
(`bob` / `bob-local`), both with verified `@example.test` emails. The Keycloak
admin console is `http://localhost:18080/admin/` (`admin` / `local-admin`).
It supports registration/login; newly registered users need a verified email
before using this API. For local manual testing, mark their email verified in
the admin console. A deployed provider must perform real email verification.

The local CLI client permits a password grant to make backend testing convenient:

```sh
TOKENS=$(curl -fsS "$OIDC_ISSUER/protocol/openid-connect/token" \
  -d grant_type=password -d client_id=paymentshare-local \
  -d username=alice -d password=alice-local -d 'scope=openid email profile')
TOKEN=$(printf '%s' "$TOKENS" | jq -r .access_token)

curl -fsS http://localhost:8082/me -H "Authorization: Bearer $TOKEN"
curl -fsS http://localhost:8082/workspaces \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"Home","currency":"USD"}'
curl -fsS http://localhost:8082/workspaces -H "Authorization: Bearer $TOKEN"
# Use the workspace ID from the creation response:
curl -fsS -X PATCH http://localhost:8082/workspaces/WORKSPACE_ID \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"New home"}'
```

Use `access_token`, never `id_token`. Local credentials, password grants,
development-mode Keycloak, and HTTP are for local development only. Deployed
clients should use their provider's authorization-code flow with PKCE and HTTPS.
The supplied realm follows Keycloak's documented [container setup](https://www.keycloak.org/server/containers)
and [startup realm import](https://www.keycloak.org/server/importExport).
Realm imports skip existing realms; editing the JSON does not overwrite a realm
already persisted in the local Keycloak volume.

## Configuration and identity contract

| Variable | Meaning |
| --- | --- |
| `DATABASE_URL` | Required PostgreSQL connection string |
| `OIDC_ISSUER` | Required exact discovery issuer URL |
| `OIDC_AUDIENCE` | Required dedicated API audience, distinct from the login client |
| `SUPPORTED_CURRENCIES` | Comma-separated supported codes; defaults to `USD,EUR,GBP,IRR` |
| `HTTP_ADDR` | Defaults to `127.0.0.1:8082` |
| `TEST_DATABASE_URL` | Required for tests; role must be able to create schemas |

The provider must issue RS256-signed JWT **access tokens** carrying `iss`, `sub`,
`aud`, `exp`, `email`, and boolean `email_verified: true`. Supported purpose
markers are the protected JWT header `typ: at+jwt` (also `application/at+jwt`),
or Keycloak's protected `typ: JWT` with the signed payload `typ: Bearer`.
Other token profiles and opaque access tokens require an adapter change; merely
matching issuer and audience never makes an ID token acceptable. The verifier
uses discovery/JWKS, caches keys, validates signatures locally, and rejects
expired or not-yet-valid tokens. It does not introspect tokens on each request.
Provider logout/revocation therefore takes effect when an issued token expires;
workspace authorization is read from PostgreSQL on every request.

Registration in the API happens on the first authenticated request. Internal
user IDs are stable for `(issuer, subject)`, even if the verified email changes.
Emails are stored exactly as verified; they are not normalized or used to merge
identities. An email already held by another identity returns `409`. The provider
must restrict verification to addresses the person controls. Hosted provider
selection and production provisioning remain deployment work.

## HTTP contract

All endpoints require `Authorization: Bearer <access_token>` and return JSON.
Successful workspace responses contain `id`, `name`, `currency`, and `owner_id`.

| Method and path | Request / response |
| --- | --- |
| `GET /me` | `200` with internal `id` and verified `email` |
| `POST /workspaces` | `{ "name": "Home", "currency": "USD" }`; `201` plus `Location` |
| `GET /workspaces` | `200` array of accessible workspaces, ordered by ID; empty array when none |
| `GET /workspaces/{id}` | `200` workspace |
| `PATCH /workspaces/{id}` | `{ "name": "New home" }`; owner-only, `200` updated workspace |

Names are trimmed, must contain 1–120 Unicode characters, and cannot contain
control characters. Currency must exactly match a configured code. Rename
cannot change currency or ownership. Unknown fields, malformed JSON, invalid
IDs, and invalid workspace values return `400`; bodies are limited to 16 KiB.
Authentication failures return `401` with `WWW-Authenticate`; inaccessible or
missing workspaces return `404`. A future ordinary member attempting an owner
operation receives `403`. Errors are `{ "error": "message" }`.

Each workspace has one owner, who is also a member, and one currency. Workspace
and owner membership are committed atomically; PostgreSQL enforces the owner
membership relationship. Renames lock the workspace and caller's membership,
check the current owner, and commit the change together. Failed mutations leave
workspace state unchanged. Creation request-key semantics and currency changes
are not part of this ticket.

## Tests and checks

```sh
# After exporting .env as above; Keycloak is not required for automated tests.
docker compose up -d --wait postgres
make check
make test
# Focused HTTP scenario:
go test ./integration -run TestOnlyOwnerCanReadAndRenameWorkspace -count=1
```

Tests fail explicitly if `TEST_DATABASE_URL` is missing. Each test creates its
own random PostgreSQL schema, applies migrations, exercises the API over HTTP
with a controlled local RSA/JWKS identity provider, and drops its schema on
completion. They never truncate application tables. Assertions observe HTTP
responses, including persisted data after API restart, concurrent initial
registration, cross-user isolation, rejected writes, and invalid access tokens.

The single Payment Share bounded context is organized as:

- `internal/domain`: user and workspace concepts and workspace invariants.
- `internal/application`: use cases, authorization, repository contracts, transactions.
- `internal/httpapi`, `internal/identity`, `internal/postgres`: adapters.
- `internal/server`: composition; `cmd/api` and `cmd/migrate`: executables.
- `migrations`: embedded, versioned PostgreSQL schema.
