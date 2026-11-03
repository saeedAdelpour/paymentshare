//go:build integration

package integration_test

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/saeedAdelpour/paymentshare/internal/postgres"
)

var migrationBinary string

func TestMain(m *testing.M) { os.Exit(runSuite(m)) }

func runSuite(m *testing.M) int {
	directory, err := os.MkdirTemp("", "paymentshare-integration-")
	if err != nil {
		fmt.Fprintln(os.Stderr, "create integration directory failed")
		return 1
	}
	defer func() { _ = os.RemoveAll(directory) }()
	migrationBinary = filepath.Join(directory, "migrate")
	build := exec.Command("go", "build", "-o", migrationBinary, "./cmd/migrate")
	build.Dir = ".."
	if output, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build migration command: %s\n", output)
		return 1
	}
	return m.Run()
}

// Each test owns one random schema. The configured database and its other schemas
// are never migrated, truncated, or dropped by this suite.
func newDatabase(t *testing.T) (string, *sql.DB) {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("TEST_DATABASE_URL is required; role must be able to create test schemas")
	}
	_, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatal("invalid TEST_DATABASE_URL (details suppressed)")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, admin, err := postgres.Open(ctx, dsn)
	if err != nil {
		t.Fatal("connect test PostgreSQL failed (details suppressed)")
	}
	schema := fmt.Sprintf("test_%x", rand.Text())
	if _, err := admin.ExecContext(ctx, "CREATE SCHEMA "+schema); err != nil {
		_ = admin.Close()
		t.Fatal("create isolated test schema failed")
	}
	t.Cleanup(func() {
		cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		if _, err := admin.ExecContext(cleanup, "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Error("drop isolated test schema failed")
		}
		_ = admin.Close()
	})
	isolated := dsn + " search_path=" + schema
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		parsed, err := url.Parse(dsn)
		if err != nil {
			t.Fatal("invalid test database URL")
		}
		query := parsed.Query()
		query.Set("search_path", schema)
		parsed.RawQuery = query.Encode()
		isolated = parsed.String()
	}
	_, pool, err := postgres.Open(ctx, isolated)
	if err != nil {
		t.Fatal("connect isolated test schema failed")
	}
	t.Cleanup(func() { _ = pool.Close() })
	var actual string
	if err := pool.QueryRowContext(ctx, "SELECT current_schema()").Scan(&actual); err != nil || actual != schema {
		t.Fatal("test database schema isolation failed")
	}
	return isolated, pool
}

func runMigration(t *testing.T, dsn string, target string, success bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, migrationBinary, "-to", target)
	command.Env = append(os.Environ(), "DATABASE_URL="+dsn)
	output, err := command.CombinedOutput()
	if (err == nil) != success {
		t.Fatalf("migration target %s: success=%v; output=%s", target, err == nil, output)
	}
}
