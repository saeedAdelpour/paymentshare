# Persistence, migrations, and administration research

Date: 2026-11-03

These findings informed the user-approved GORM + QOR5 + Goose implementation target and three-ticket foundation sequence on 2026-11-03. Alternatives and original assessments below are retained as research context. No workload benchmark or integration prototype was run.

## Requirements driving the choice

The requested administration panel should expose application tables and users, allow staff read/write/delete permissions, and provide a superuser with all administrative permissions. Migrations should support moving between application versions. Application quality, performance, and other non-functional requirements guide the choice.

The existing spec requires atomic financial writes, retained revisions, maintained net balances, request-key idempotency, and workspace access control. The panel's write/delete semantics must be reconciled with those requirements. The user subsequently approved administration in the first release and confirmed that superusers preserve financial invariants and history. Existing backend tickets now clarify that administration needs its own ticket breakdown.

## ORM comparison

| Option | Documented strengths relevant here | Assessment for this application |
| --- | --- | --- |
| GORM | Explicit transactions; configurable performance features; direct integration with QOR5's GORM2 data operator. | Leading proposal when a Go administration framework is a major requirement. Keep persistence models inside adapters and use explicit application transactions. |
| Ent | Schema-based generated CRUD; versioned migration integration with Atlas, including alternative migration file formats. | Strong alternative if generated APIs and schema definitions are prioritized. The QOR5 integration would require custom data operations rather than its stock GORM operator. |
| Bun | SQL-first Go ORM and migration support. | Strong alternative for developers prioritizing direct SQL control. A Django-like administration integration is additional work rather than a demonstrated native integration in this investigation. |

Sources: [GORM transactions](https://gorm.io/docs/transactions.html), [GORM performance](https://gorm.io/docs/performance.html), [Ent generated CRUD](https://entgo.io/docs/crud/), [Ent versioned migrations](https://entgo.io/docs/versioned-migrations/), [Bun guide](https://bun.uptrace.dev/guide/), [QOR5 GORM2 operator source](https://raw.githubusercontent.com/qor5/admin/main/presets/gorm2op/operator.go).

No source establishes which ORM is fastest for Payment Share. Selecting GORM for ecosystem fit is an engineering assessment, not a performance result. Verify the approved HTTP workload, query counts, indexes, lock contention, and transaction duration with the actual implementation.

## Administration candidates

### QOR5

QOR5 presets build list/create/edit/delete interfaces for registered Go models and expose an `http.Handler`. Registration is explicit; this is not proof of automatic discovery of every PostgreSQL table. [Presets documentation](https://docs.qor5.com/basics/presets-instant-crud.html).

Its built-in operator takes `*gorm.DB`; saves and deletes call GORM persistence operations. Reusing these defaults for financial data would bypass application use cases unless specifically integrated. [Operator source](https://raw.githubusercontent.com/qor5/admin/main/presets/gorm2op/operator.go).

Custom writes are supported: `EditingBuilder` exposes `SaveFunc`, `DeleteFunc`, their wrappers, and a validation function. This provides a concrete integration point for application use cases that preserve balances, revisions, actors, and transaction boundaries. [Editing implementation](https://raw.githubusercontent.com/qor5/admin/main/presets/editing.go).

Permissions distinguish list/get/create/update/delete and custom actions, with resource patterns and context conditions. An allow-all wildcard policy is documented. A superuser can therefore be represented by an administrative subject/policy; a special built-in bypass flag was not established. Deny-policy interactions and newly registered resources need verification. The Role module provides a UI for stored permission policies. [Permissions](https://docs.qor5.com/presets-guide/permissions.html), [Role UI](https://docs.qor5.com/presets-guide/role.html).

Login documentation covers passwords, TOTP, sessions, and OAuth through Goth. Turnkey compatibility with this application's configurable OIDC access-token contract was not established; staff browser login and session integration require validation. [Login documentation](https://docs.qor5.com/basics/login.html).

### GoAdmin and original QOR

GoAdmin advertises built-in RBAC. Its documented permissions use HTTP methods and routes, so per-resource/action permissions require configuring the corresponding routes. It remains an alternative, but model reuse and domain-command integration were not demonstrated here. [Official repository](https://github.com/GoAdminGroup/go-admin), [RBAC documentation](https://github.com/GoAdminGroup/docs/blob/master/zh/admin/rbac.md).

Original QOR registers GORM models and documents resource/field/action permissions, but its quick start imports the older `github.com/jinzhu/gorm` path. Prefer evaluating QOR5 for a new application rather than assuming compatibility with current GORM. [Original QOR quick start](https://doc.getqor.com/get_started.html), [authorization](https://doc.getqor.com/admin/authentication.html).

## Migration strategy

GORM AutoMigrate is schema synchronization, not the requested versioned rollback workflow. It explicitly avoids deleting unused columns. Use reviewed migration files with a separate manager rather than treating application startup synchronization as version history. [GORM migrations](https://gorm.io/docs/migration.html).

Goose supports SQL Up/Down sections, transactions by default, and target-version `up-to`/`down-to` commands. It is a leading proposal for explicit reversible migration files. Bun also has migrations, but its rollback operates on migration groups; that unit needs consideration when switching versions. Ent's Atlas integration can generate SQL in Goose or golang-migrate formats. [Goose SQL migrations](https://github.com/pressly/goose), [Goose commands](https://pressly.github.io/goose/documentation/cli-commands/), [Bun migrations](https://bun.uptrace.dev/guide/migrations.html), [Ent formats](https://entgo.io/docs/versioned-migrations/).

A Down migration can reverse schema operations without recovering information destroyed by an earlier change. For example, adding a dropped column back does not restore its values. This is an inherent data-preservation limit, not something ORM selection solves. Development branch switching can use separate databases or reversible target-version changes; production rollback requires its own compatibility and data-preservation policy.

## Proposed next validation

The selected combination is GORM + QOR5 + Goose, with integration verification required by the foundation acceptance criteria. Verify staff identity/session handling, permission changes and superuser behavior, model registration, custom save/delete routing, migration Up/Down/Up, and representative financial writes under concurrent requests. The user subsequently confirmed that superusers preserve domain invariants and retained history, financial deletion means voiding, and production rollback has a separate compatibility/recovery contract from local migration reversal. Staff now share the application OIDC provider; only superusers manage staff privileges, and user deletion means deactivation. Detailed integration behavior remains to be verified.

## Proposed foundation quality tooling — 2026-11-03

This is a small proposed baseline, not a claim that these tools establish production readiness or are universally best.

| Check | Proposed role |
| --- | --- |
| gofumpt | One formatter for authored Go code. It applies stricter formatting while remaining compatible with gofmt, so a separate gofmt gate is unnecessary. |
| golangci-lint v2 | Explicit initial analyzer set: `govet`, `staticcheck`, `errcheck`, `ineffassign`, and `unused`. Avoid enabling every available linter; add checks when they enforce an agreed rule or prevent a demonstrated mistake. |
| govulncheck | Scan known Go vulnerabilities and report call paths reaching vulnerable code. This complements static correctness checks. |
| `go test -race ./...` | Detect memory data races exercised by tests. Keep the approved real-PostgreSQL HTTP concurrency tests: a memory race detector does not establish database transaction correctness. |

Sources: [gofumpt](https://github.com/mvdan/gofumpt), [golangci-lint analyzers](https://golangci-lint.run/docs/linters/), [govulncheck tutorial](https://go.dev/doc/tutorial/govulncheck), [Go race detector](https://go.dev/doc/articles/race_detector).

The `govet` analyzer uses Go vet's passes, so run it through the lint gate rather than adding a duplicate mandatory standalone vet job. `go vet ./...` remains useful as a direct local command. In v2, formatters have a separate configuration section; use `golangci-lint fmt` for formatting and `golangci-lint run` for analysis. [Go vet](https://pkg.go.dev/cmd/vet), [v2 migration guide](https://golangci-lint.run/docs/product/migration-guide/).

Pin the Go toolchain and tool releases consistently in local setup and CI, and upgrade deliberately. Go supports tool dependencies through `tool` and versioned `require` directives; use an isolated tooling module where dependency separation is useful. For golangci-lint specifically, follow its recommendation to install a pinned release binary rather than adding it to the application's module with `go tool`: its installation guide documents dependency and build reproducibility concerns with source installation. [Go toolchains](https://go.dev/doc/toolchain), [Go tool dependencies](https://go.dev/ref/mod#tool-directive), [golangci-lint installation](https://golangci-lint.run/docs/welcome/install/local/), [CI version pinning](https://golangci-lint.run/docs/welcome/install/ci/).

Run formatting and analysis on each proposed change, and race-enabled behavioral tests with their PostgreSQL dependency in CI. Vulnerability checks also need a scheduled run because vulnerability information can change without a code change. Keep load tests separate from race-instrumented tests so instrumentation does not distort the application's performance measurements. These scheduling choices are engineering proposals; no tooling has been installed or configured by this research.

## Subsequent accepted foundation decisions

On 2026-11-03, the user accepted the proposed gofumpt/golangci-lint v2/govulncheck/race-enabled test baseline with pinned versions and matching local/CI commands. Staff share the application OIDC provider; staff status, superuser status, roles, and permissions are stored by Payment Share. A repeatable CLI command provisions the initial superuser for a specific provider identity. Framework integration and detailed browser/session configuration still need implementation validation.

The user subsequently approved GORM + QOR5 + Goose and the three-ticket sequence. The delivery contract is the [foundation spec](../payment-share-foundation/spec.md).
