package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dsn string, logger *slog.Logger) (*pgxpool.Pool, error) {
	if dsn == "" {
		logger.Error("DSN is not set")
		return nil, fmt.Errorf("DSN is not set")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("Failed to parse DSN", "error", err)
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("Failed to create DB pool", "error", err)
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Check if 'books' table exists
	var exists bool
	err = pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'books'
		)
	`).Scan(&exists)

	if err != nil {
		logger.Error("Failed to check books table existence", "error", err)
		return nil, fmt.Errorf("failed to check books table: %w", err)
	}

	if !exists {
		logger.Warn("'books' table not found. Did you forget to run migrations?")
	}

	return pool, nil
}
