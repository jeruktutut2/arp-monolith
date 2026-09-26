package domain

import (
	"context"

	"github.com/google/uuid"
)

// Outbound Port: Repository interface for user storage
type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByCompanyID(ctx context.Context, companyID uuid.UUID) ([]User, error)
	Save(ctx context.Context, user *User) error
}
