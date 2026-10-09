package repository

import (
	"context"

	"erp_monolith/backend/internal/modules/system/user/domain"
	"erp_monolith/backend/internal/platform/database"
)

type userRepository struct {
	db database.DBTX
}

func NewUserRepository(db database.DBTX) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
