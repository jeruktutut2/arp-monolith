package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX abstracts common database execution methods, satisfied by *pgxpool.Pool and pgx.Tx.
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type txKey struct{}
type txRootKey struct{}

// TxManager defines the interface for managing database transactions.
type TxManager interface {
	Begin(ctx context.Context) (context.Context, error)
	CommitOrRollback(ctx context.Context, err error) error
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type txManager struct {
	db Postgresql
}

// NewTxManager creates a new transaction manager instance.
func NewTxManager(db Postgresql) TxManager {
	return &txManager{db: db}
}

// Begin starts a new database transaction and injects it into the returned context.
// If a transaction already exists in context, it marks the context as nested and reuses it.
func (tm *txManager) Begin(ctx context.Context) (context.Context, error) {
	if tx := TxFromContext(ctx); tx != nil {
		return context.WithValue(ctx, txRootKey{}, false), nil
	}

	tx, err := tm.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return ctx, fmt.Errorf("begin tx: %w", err)
	}

	ctxWithTx := ContextWithTx(ctx, tx)
	return context.WithValue(ctxWithTx, txRootKey{}, true), nil
}

// CommitOrRollback commits the transaction in ctx if err is nil, or rolls it back if err is not nil.
// If the context represents a nested transaction and err is nil, commit is deferred to the root.
func (tm *txManager) CommitOrRollback(ctx context.Context, err error) error {
	tx := TxFromContext(ctx)
	if tx == nil {
		return err
	}

	if isNested, ok := ctx.Value(txRootKey{}).(bool); ok && !isNested && err == nil {
		return nil
	}

	return tm.db.CommitOrRollback(ctx, tx, err)
}

// RunInTx executes the provided function within a transaction.
// If fn returns an error or panics, the transaction is rolled back.
func (tm *txManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	txCtx, err := tm.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tm.CommitOrRollback(txCtx, fmt.Errorf("panic: %v", p))
			err = fmt.Errorf("panic recovered: %v", p)
			return
		}
		err = tm.CommitOrRollback(txCtx, err)
	}()

	return fn(txCtx)
}

// TxFromContext extracts an active pgx.Tx from the context, or returns nil if none exists.
func TxFromContext(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}

// ContextWithTx injects a pgx.Tx into the given context.
func ContextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetExecutor returns the transaction from ctx if present; otherwise returns fallback DBTX.
func GetExecutor(ctx context.Context, fallback DBTX) DBTX {
	if tx := TxFromContext(ctx); tx != nil {
		return tx
	}
	return fallback
}
