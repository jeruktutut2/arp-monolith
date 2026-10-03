package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgresql defines the interface for database access, transaction lifecycle, and connection cleanup
type Postgresql interface {
	GetDB() *pgxpool.Pool
	Begin(ctx context.Context) (pgx.Tx, error)
	Commit(ctx context.Context, tx pgx.Tx) error
	Rollback(ctx context.Context, tx pgx.Tx) error
	CommitOrRollback(ctx context.Context, tx pgx.Tx, err *error) error
	Close()
}

type postgresql struct {
	pool *pgxpool.Pool
}

func (d *postgresql) GetDB() *pgxpool.Pool {
	return d.pool
}

func (d *postgresql) Begin(ctx context.Context) (pgx.Tx, error) {
	if d.pool == nil {
		return nil, errors.New("database pool is not initialized")
	}
	return d.pool.Begin(ctx)
}

func (d *postgresql) Commit(ctx context.Context, tx pgx.Tx) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}
	return tx.Commit(ctx)
}

func (d *postgresql) Rollback(ctx context.Context, tx pgx.Tx) error {
	if tx == nil {
		return nil
	}
	err := tx.Rollback(ctx)
	if errors.Is(err, pgx.ErrTxClosed) {
		return nil
	}
	return err
}

func (d *postgresql) CommitOrRollback(ctx context.Context, tx pgx.Tx, err *error) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}
	if err != nil && *err != nil {
		return d.Rollback(ctx, tx)
	}
	commitErr := d.Commit(ctx, tx)
	if commitErr != nil && err != nil && *err == nil {
		*err = commitErr
	}
	return commitErr
}

func (d *postgresql) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

// NewPostgresql initializes a pgxpool optimized for PgBouncer and returns the Postgresql interface
func NewPostgresql(
	url string,
	maxConns int32,
	minConns int32,
	maxConnLifetime time.Duration,
	maxConnIdleTime time.Duration,
) (Postgresql, error) {
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

	return &postgresql{pool: pool}, nil
}
