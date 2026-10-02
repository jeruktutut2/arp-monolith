package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"erp_monolith/backend/internal/modules/system/user/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresUserRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepo creates a driven adapter repository for users
func NewPostgresUserRepo(pool *pgxpool.Pool) domain.UserRepository {
	return &postgresUserRepo{pool: pool}
}

const userColumns = `
	id, company_id, branch_id, email, password_hash, name, role,
	is_active, failed_login_attempts, locked_until,
	created_at, created_by, updated_at, updated_by
`

func (r *postgresUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM usr_users WHERE id = $1`, userColumns)
	var u domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.CompanyID, &u.BranchID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
		&u.IsActive, &u.FailedLoginAttempts, &u.LockedUntil,
		&u.CreatedAt, &u.CreatedBy, &u.UpdatedAt, &u.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user by id: %w", err)
	}
	return &u, nil
}

func (r *postgresUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM usr_users WHERE LOWER(email) = LOWER($1)`, userColumns)
	var u domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.CompanyID, &u.BranchID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
		&u.IsActive, &u.FailedLoginAttempts, &u.LockedUntil,
		&u.CreatedAt, &u.CreatedBy, &u.UpdatedAt, &u.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}
	return &u, nil
}

func (r *postgresUserRepo) FindByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM usr_users WHERE company_id = $1 ORDER BY name ASC`, userColumns)
	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID, &u.CompanyID, &u.BranchID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
			&u.IsActive, &u.FailedLoginAttempts, &u.LockedUntil,
			&u.CreatedAt, &u.CreatedBy, &u.UpdatedAt, &u.UpdatedBy,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *postgresUserRepo) Save(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO usr_users (
			id, company_id, branch_id, email, password_hash, name, role,
			is_active, failed_login_attempts, locked_until,
			created_at, created_by, updated_at, updated_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (id) DO UPDATE SET
			company_id = EXCLUDED.company_id,
			branch_id = EXCLUDED.branch_id,
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash,
			name = EXCLUDED.name,
			role = EXCLUDED.role,
			is_active = EXCLUDED.is_active,
			failed_login_attempts = EXCLUDED.failed_login_attempts,
			locked_until = EXCLUDED.locked_until,
			updated_at = EXCLUDED.updated_at,
			updated_by = EXCLUDED.updated_by
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.CompanyID, user.BranchID, user.Email, user.PasswordHash, user.Name, user.Role,
		user.IsActive, user.FailedLoginAttempts, user.LockedUntil,
		user.CreatedAt, user.CreatedBy, user.UpdatedAt, user.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (r *postgresUserRepo) UpdateLoginAttempts(ctx context.Context, id uuid.UUID, attempts int, lockedUntil *time.Time) error {
	query := `
		UPDATE usr_users
		SET failed_login_attempts = $2, locked_until = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id, attempts, lockedUntil)
	if err != nil {
		return fmt.Errorf("failed to update login attempts: %w", err)
	}
	return nil
}
