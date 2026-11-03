# Administration preserves financial invariants

The first release includes an administration panel with staff permissions and a superuser who has all administrative permissions without individual grants. Superusers still preserve domain invariants and retained history: financial corrections and cancellation invoke application use cases, cancellation means voiding, and balances and revisions have no ordinary direct-edit operation.

Generic table writes would bypass maintained balances, revision history, and concurrency guarantees. Administrative authority therefore changes who may invoke an operation, while the operation retains its financial consistency contract. The user accepted this trade-off on 2026-11-03.
