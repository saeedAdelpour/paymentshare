// Package migrations embeds reviewed SQL migrations. Execution is explicit.
package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

const LatestVersion int64 = 1

//go:embed *.sql
var files embed.FS

// Apply moves the database to an exact local target version. The session advisory
// lock serializes migration commands, while Goose transactions make each step atomic.
func Apply(ctx context.Context, pool *sql.DB, target int64) error {
	if target < 0 || target > LatestVersion {
		return errors.New("migration target is not supported by this release")
	}
	locker, err := lock.NewPostgresSessionLocker(lock.WithLockTimeout(1, 30), lock.WithUnlockTimeout(1, 5))
	if err != nil {
		return err
	}
	provider, err := goose.NewProvider(goose.DialectPostgres, pool, files, goose.WithSessionLocker(locker))
	if err != nil {
		return err
	}
	current, err := provider.GetDBVersion(ctx)
	if err != nil {
		return err
	}
	if current > LatestVersion {
		return errors.New("database is newer than this migration binary; use the matching release")
	}
	if current == target {
		return nil
	}
	if current > target {
		_, err = provider.DownTo(ctx, target)
	} else {
		_, err = provider.UpTo(ctx, target)
	}
	return err
}

// CheckSchema is strictly read-only: it never initializes Goose or changes a schema.
func CheckSchema(ctx context.Context, pool *sql.DB) error {
	var version int64
	err := pool.QueryRowContext(ctx, "SELECT version_id FROM goose_db_version ORDER BY id DESC LIMIT 1").Scan(&version)
	if err != nil {
		return fmt.Errorf("read schema version (run the migration command before startup): %w", err)
	}
	if version != LatestVersion {
		return fmt.Errorf("schema version %d is incompatible; require version %d", version, LatestVersion)
	}
	return nil
}
