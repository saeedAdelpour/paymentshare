# Shared identity and application-owned staff permissions

Staff authenticate through the same OIDC provider as application users. Payment Share owns staff status, superuser status, roles, and administrative permissions; ordinary registration grants no administrative privileges. A repeatable CLI command grants the initial superuser role to a specific provider identity.

This keeps identity with one provider while allowing Payment Share to control administrative authorization independently of workspace ownership and identity-provider configuration. The user accepted this trade-off on 2026-11-03. Only superusers initially grant/revoke staff status, change administrative permissions, and promote other superusers. The administration framework's browser/session integration still needs validation.
