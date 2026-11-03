# GORM, QOR5, and Goose for the project foundation

Use GORM for PostgreSQL persistence, QOR5 for administration, and Goose for reviewed versioned SQL migrations. The user approved this implementation target on 2026-11-03 after comparing GORM, Ent, and Bun against the administration and migration requirements.

QOR5's GORM integration and replaceable mutation handlers favor this combination for a Django-like administration experience while keeping financial rules in application use cases. Goose provides explicit target-version Up/Down migrations independently of ORM schema synchronization. The foundation must validate OIDC integration, permissions/superuser behavior, custom handlers, and migration behavior; the selection itself establishes no performance result.
