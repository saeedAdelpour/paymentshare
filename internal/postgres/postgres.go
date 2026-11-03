// Package postgres owns PostgreSQL connection setup. It never changes the schema.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open creates a bounded connection pool and verifies connectivity using ctx.
// The caller owns the returned SQL pool and must close it.
func Open(ctx context.Context, dsn string) (*gorm.DB, *sql.DB, error) {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("parse PostgreSQL configuration: %w", err)
	}
	config.ConnectTimeout = 5 * time.Second
	pool := stdlib.OpenDB(*config)
	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(30 * time.Minute)
	pool.SetConnMaxIdleTime(5 * time.Minute)
	if err := pool.PingContext(ctx); err != nil {
		_ = pool.Close()
		return nil, nil, fmt.Errorf("connect PostgreSQL: %w", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: pool}), &gorm.Config{
		DisableAutomaticPing: true,
		Logger:               logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		_ = pool.Close()
		return nil, nil, fmt.Errorf("initialize GORM: %w", err)
	}
	return db, pool, nil
}
