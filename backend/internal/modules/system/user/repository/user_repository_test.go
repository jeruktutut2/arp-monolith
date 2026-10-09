package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"erp_monolith/backend/internal/modules/system/user/domain"
	"erp_monolith/backend/internal/modules/system/user/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

func TestUserRepository_FindByEmail(t *testing.T) {
	columns := []string{
		"id", "company_id", "branch_id", "email", "password_hash", "name", "role",
		"is_active", "failed_login_attempts", "locked_until",
		"created_at", "created_by", "updated_at", "updated_by",
	}

	expectedQuery := `SELECT id, company_id, branch_id, email, password_hash, name, role,
		       is_active, failed_login_attempts, locked_until,
		       created_at, created_by, updated_at, updated_by
		FROM usr_users
		WHERE email = \$1`

	t.Run("Success - User ditemukan", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create pgxmock: %v", err)
		}
		defer mock.Close()

		repo := repository.NewUserRepository(mock)
		ctx := context.Background()

		userID := uuid.New()
		companyID := uuid.New()
		now := time.Now().Truncate(time.Second)
		email := "admin@erp.local"

		mock.ExpectQuery(expectedQuery).
			WithArgs(email).
			WillReturnRows(
				mock.NewRows(columns).AddRow(
					userID, companyID, nil, email, "hashed_secret", "Admin User", "admin",
					true, 0, nil, now, nil, now, nil,
				),
			)

		user, err := repo.FindByEmail(ctx, email)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if user == nil {
			t.Fatal("expected user, got nil")
		}
		if user.ID != userID {
			t.Errorf("expected ID %v, got %v", userID, user.ID)
		}
		if user.Email != email {
			t.Errorf("expected email %s, got %s", email, user.Email)
		}
		if user.Name != "Admin User" {
			t.Errorf("expected name 'Admin User', got %s", user.Name)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled mock expectations: %v", err)
		}
	})

	t.Run("Not Found - User tidak ditemukan (pgx.ErrNoRows)", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create pgxmock: %v", err)
		}
		defer mock.Close()

		repo := repository.NewUserRepository(mock)
		ctx := context.Background()
		email := "unknown@erp.local"

		mock.ExpectQuery(expectedQuery).
			WithArgs(email).
			WillReturnRows(mock.NewRows(columns))

		user, err := repo.FindByEmail(ctx, email)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			t.Errorf("expected pgx.ErrNoRows, got %v", err)
		}
		if user != nil {
			t.Errorf("expected nil user, got %+v", user)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled mock expectations: %v", err)
		}
	})

	t.Run("Error - Database query failure", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("failed to create pgxmock: %v", err)
		}
		defer mock.Close()

		repo := repository.NewUserRepository(mock)
		ctx := context.Background()
		email := "error@erp.local"
		dbErr := errors.New("connection reset by peer")

		mock.ExpectQuery(expectedQuery).
			WithArgs(email).
			WillReturnError(dbErr)

		user, err := repo.FindByEmail(ctx, email)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, dbErr) {
			t.Errorf("expected error %v, got %v", dbErr, err)
		}
		if user != nil {
			t.Errorf("expected nil user, got %+v", user)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled mock expectations: %v", err)
		}
	})
}
