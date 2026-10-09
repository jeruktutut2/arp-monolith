package repository

import (
	"erp_monolith/backend/internal/modules/system/user/domain"
	"erp_monolith/backend/internal/platform/database"
)

type userRepository struct {
	db database.DBTX
}

func NewUserRepository(db database.DBTX) domain.UserRepository {
	return &userRepository{db: db}
}
