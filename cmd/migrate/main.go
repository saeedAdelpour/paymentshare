package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saeedAdelpour/paymentshare/internal/postgres"
	"github.com/saeedAdelpour/paymentshare/migrations"
)

func main() {
	target := flag.Int64("to", migrations.LatestVersion, "exact target version (0 reverses all local migrations)")
	flag.Parse()
	if flag.NArg() != 0 || *target < 0 || *target > migrations.LatestVersion {
		fmt.Fprintln(os.Stderr, "unsupported migration target or arguments")
		os.Exit(1)
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is required")
		os.Exit(1)
	}
	signalContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(signalContext, 60*time.Second)
	defer cancel()
	_, pool, err := postgres.Open(ctx, dsn)
	if err == nil {
		defer func() { _ = pool.Close() }()
		err = migrations.Apply(ctx, pool, *target)
	}
	if err != nil {
		// Raw PostgreSQL errors may contain credentials or sensitive SQL values.
		fmt.Fprintln(os.Stderr, "migration failed; check database availability, permissions, schema version, and database-side logs")
		os.Exit(1)
	}
	fmt.Printf("database migrated to version %d\n", *target)
}
