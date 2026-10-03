package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database defines the interface for database access, transaction lifecycle, and connection cleanup
type Database interface {
	GetDB() *pgxpool.Pool
	Begin(ctx context.Context) (pgx.Tx, error)
	Commit(ctx context.Context, tx pgx.Tx) error
	Rollback(ctx context.Context, tx pgx.Tx) error
	CommitOrRollback(ctx context.Context, tx pgx.Tx, err *error)
	Close()
}

// DB is an alias for Database
type DB = Database

type pgDatabase struct {
	pool *pgxpool.Pool
}

// NewDatabase wraps a pgxpool.Pool into the Database interface
func NewDatabase(pool *pgxpool.Pool) Database {
	return &pgDatabase{pool: pool}
}

func (d *pgDatabase) GetDB() *pgxpool.Pool {
	return d.pool
}

func (d *pgDatabase) Begin(ctx context.Context) (pgx.Tx, error) {
	if d.pool == nil {
		return nil, errors.New("database pool is not initialized")
	}
	return d.pool.Begin(ctx)
}

func (d *pgDatabase) Commit(ctx context.Context, tx pgx.Tx) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}
	return tx.Commit(ctx)
}

func (d *pgDatabase) Rollback(ctx context.Context, tx pgx.Tx) error {
	if tx == nil {
		return nil
	}
	err := tx.Rollback(ctx)
	if errors.Is(err, pgx.ErrTxClosed) {
		return nil
	}
	return err
}

func (d *pgDatabase) CommitOrRollback(ctx context.Context, tx pgx.Tx, err *error) {
	if tx == nil {
		return
	}
	if p := recover(); p != nil {
		_ = tx.Rollback(ctx)
		panic(p)
	} else if err != nil && *err != nil {
		_ = tx.Rollback(ctx)
	} else {
		if commitErr := tx.Commit(ctx); commitErr != nil && err != nil && *err == nil {
			*err = commitErr
		}
	}
}

func (d *pgDatabase) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

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
