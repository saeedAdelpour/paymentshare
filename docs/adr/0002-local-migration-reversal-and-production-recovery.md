# Separate local migration reversal from production recovery

Support tested target-version Up/Down migrations for local application version switching. Production rollback preserves compatibility with the previous application release where practical and requires an explicit recovery strategy for destructive schema changes or data transformations.

A schema inverse cannot recover information already destroyed, so the presence of a Down migration cannot promise production data recovery. The user accepted separate local version-switching and production recovery guarantees on 2026-11-03. The subsequent stack decision selects Goose; migrations run as an explicit release step before the API/admin application starts, and startup checks schema compatibility. Detailed migration commands and recovery procedures are verified during the foundation.
