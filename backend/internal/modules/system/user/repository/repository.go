package repository

import (
	"context"
	"fmt"

	"erp_monolith/backend/internal/modules/system/user/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresUserRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepo creates a driven adapter repository for users
func NewPostgresUserRepo(pool *pgxpool.Pool) domain.UserRepository {
	return &postgresUserRepo{pool: pool}
}

func (r *postgresUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, company_id, email, name, role, is_active, created_at, updated_at FROM usr_users WHERE id = $1`
	var u domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.CompanyID, &u.Email, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	return &u, nil
}

func (r *postgresUserRepo) FindByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.User, error) {
	query := `SELECT id, company_id, email, name, role, is_active, created_at, updated_at FROM usr_users WHERE company_id = $1 ORDER BY name ASC`
	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.CompanyID, &u.Email, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *postgresUserRepo) Save(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO usr_users (id, company_id, email, name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			role = EXCLUDED.role,
			is_active = EXCLUDED.is_active,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.CompanyID, user.Email, user.Name, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}
