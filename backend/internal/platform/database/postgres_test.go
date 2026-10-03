package database_test

import (
	"context"
	"errors"
	"testing"

	"erp_monolith/backend/internal/platform/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// mockTx implements pgx.Tx for testing transaction lifecycle methods
type mockTx struct {
	committed  bool
	rolledBack bool
	commitErr  error
	rbErr      error
}

func (m *mockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return m, nil
}

func (m *mockTx) Commit(ctx context.Context) error {
	m.committed = true
	return m.commitErr
}

func (m *mockTx) Rollback(ctx context.Context) error {
	m.rolledBack = true
	return m.rbErr
}

func (m *mockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func (m *mockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return nil
}

func (m *mockTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (m *mockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}

func (m *mockTx) Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
	return pgconn.CommandTag{}, nil
}

func (m *mockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *mockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}

func (m *mockTx) Conn() *pgx.Conn {
	return nil
}

func TestDatabase_CommitOrRollback_Success(t *testing.T) {
	db := database.NewPostgresql(nil)
	tx := &mockTx{}
	ctx := context.Background()

	var err error
	func() {
		defer db.CommitOrRollback(ctx, tx, &err)
		// Operation succeeded with no error
	}()

	if !tx.committed {
		t.Errorf("expected transaction to be committed")
	}
	if tx.rolledBack {
		t.Errorf("did not expect transaction to be rolled back")
	}
}

func TestDatabase_CommitOrRollback_OnError(t *testing.T) {
	db := database.NewPostgresql(nil)
	tx := &mockTx{}
	ctx := context.Background()

	var err error
	func() {
		defer db.CommitOrRollback(ctx, tx, &err)
		err = errors.New("business logic failed")
	}()

	if tx.committed {
		t.Errorf("did not expect transaction to be committed on error")
	}
	if !tx.rolledBack {
		t.Errorf("expected transaction to be rolled back on error")
	}
}

func TestDatabase_CommitOrRollback_OnPanic(t *testing.T) {
	db := database.NewPostgresql(nil)
	tx := &mockTx{}
	ctx := context.Background()

	var err error
	defer func() {
		p := recover()
		if p == nil {
			t.Fatalf("expected panic to be propagated")
		}
		if !tx.rolledBack {
			t.Errorf("expected transaction to be rolled back on panic")
		}
		if tx.committed {
			t.Errorf("did not expect transaction to be committed on panic")
		}
	}()

	func() {
		defer db.CommitOrRollback(ctx, tx, &err)
		panic("unexpected runtime panic")
	}()
}

func TestDatabase_GetDB_And_Close(t *testing.T) {
	db := database.NewPostgresql(nil)
	if db.GetDB() != nil {
		t.Errorf("expected nil pool for unit test")
	}
	// Calling close on nil pool should not panic
	db.Close()
}
