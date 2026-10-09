package repository

import (
	"context"
	"errors"
	"fmt"

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
	query := fmt.Sprintf(`SELECT %s FROM usr_users WHERE LOWER(email) = LOWER($1)`, userColumns)
	var u domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
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
