package repository

import (
	"context"

	"erp_monolith/backend/internal/modules/system/user/domain"
	"erp_monolith/backend/internal/platform/database"

	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db database.DBTX
}

func NewUserRepository(db database.DBTX) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, company_id, branch_id, email, password_hash, name, role,
		       is_active, failed_login_attempts, locked_until,
		       created_at, created_by, updated_at, updated_by
		FROM usr_users
		WHERE email = $1
	`

	db := database.GetExecutor(ctx, r.db)
	rows, err := db.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return nil, err
	}

	return &user, nil
}
