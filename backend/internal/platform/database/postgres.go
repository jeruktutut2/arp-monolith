package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool initializes a pgxpool optimized for PgBouncer transaction pooling
func NewPool(
	url string,
	maxConns int32,
	minConns int32,
	maxConnLifetime time.Duration,
	maxConnIdleTime time.Duration,
) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pg connection url: %w", err)
	}

	// Critical setting for PgBouncer transaction pooling compatibility:
	// PgBouncer in transaction mode does not support session-bound prepared statements.
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	if maxConns > 0 {
		poolConfig.MaxConns = maxConns
	} else {
		poolConfig.MaxConns = 30
	}

	if minConns > 0 {
		poolConfig.MinConns = minConns
	} else {
		poolConfig.MinConns = 5
	}

	if maxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = maxConnLifetime
	} else {
		poolConfig.MaxConnLifetime = time.Hour
	}

	if maxConnIdleTime > 0 {
		poolConfig.MaxConnIdleTime = maxConnIdleTime
	} else {
		poolConfig.MaxConnIdleTime = 30 * time.Minute
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	// Ping with timeout
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		// Log warning if database is not reachable yet during local dev/scaffolding
		fmt.Printf("⚠️ Warning: PostgreSQL/PgBouncer not immediately reachable: %v\n", err)
	}

	return pool, nil
}
