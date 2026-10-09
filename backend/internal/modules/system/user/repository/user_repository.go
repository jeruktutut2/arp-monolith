package repository

import (
	"erp_monolith/backend/internal/platform/database"
)

type userRepository struct {
	db database.DBTX
}
