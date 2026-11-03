//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/saeedAdelpour/paymentshare/internal/server"
)

func TestMigrationTargetsPreserveFixturesAndGateStartup(t *testing.T) {
	dsn, pool := newDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Migration tests may inspect their fixtures directly: the command has no HTTP data API.
	if _, err := pool.ExecContext(ctx, `CREATE TABLE migration_fixture (
  id integer PRIMARY KEY, label text NOT NULL, amount numeric(20,2) NOT NULL
 ); INSERT INTO migration_fixture VALUES (1, 'خانه', 123456789.25)`); err != nil {
		t.Fatal(err)
	}
	assertRejected := func() {
		t.Helper()
		startup, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		err := server.Run(startup, server.Config{HTTPAddr: "127.0.0.1:0", DatabaseURL: dsn}, slog.New(slog.NewTextHandler(io.Discard, nil)))
		if err == nil {
			t.Fatal("server accepted an incompatible schema")
		}
	}
	assertRejected()
	var trackingTable *string
	if err := pool.QueryRowContext(ctx, "SELECT to_regclass('goose_db_version')::text").Scan(&trackingTable); err != nil {
		t.Fatal(err)
	}
	if trackingTable != nil {
		t.Fatal("startup initialized migration tracking instead of remaining read-only")
	}
	for _, target := range []string{"1", "1", "0", "0", "1"} {
		runMigration(t, dsn, target, true)
		var label, amount string
		if err := pool.QueryRowContext(ctx, "SELECT label, amount::text FROM migration_fixture WHERE id = 1").Scan(&label, &amount); err != nil {
			t.Fatal(err)
		}
		if label != "خانه" || amount != "123456789.25" {
			t.Fatal("migration changed retained fixture data")
		}
		var version string
		if err := pool.QueryRowContext(ctx, "SELECT version_id::text FROM goose_db_version ORDER BY id DESC LIMIT 1").Scan(&version); err != nil {
			t.Fatal(err)
		}
		if version != target {
			t.Fatalf("migration target %s produced version %s", target, version)
		}
		if target == "0" {
			assertRejected()
		}
	}
	runMigration(t, dsn, "2", false)
	runMigration(t, dsn, "-1", false)
	// A previous-release binary must reject a future schema instead of changing it.
	if _, err := pool.ExecContext(ctx, "INSERT INTO goose_db_version (version_id, is_applied) VALUES (2, true)"); err != nil {
		t.Fatal(err)
	}
	assertRejected()
	runMigration(t, dsn, "1", false)
}

func TestConcurrentMigrationCommandsApplyBaselineOnce(t *testing.T) {
	dsn, pool := newDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	results := make(chan error, 4)
	for range 4 {
		go func() {
			command := exec.CommandContext(ctx, migrationBinary, "-to", "1")
			command.Env = append(os.Environ(), "DATABASE_URL="+dsn)
			output, err := command.CombinedOutput()
			if err != nil {
				err = fmt.Errorf("migration command failed: %s", output)
			}
			results <- err
		}()
	}
	for range 4 {
		if err := <-results; err != nil {
			t.Error(err)
		}
	}
	var applied int
	if err := pool.QueryRowContext(ctx, "SELECT count(*) FROM goose_db_version WHERE version_id = 1").Scan(&applied); err != nil {
		t.Fatal(err)
	}
	if applied != 1 {
		t.Fatalf("baseline recorded %d times after concurrent commands", applied)
	}
}
