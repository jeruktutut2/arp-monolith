package domain

import (
	"context"

	"github.com/google/uuid"
)

// Inbound Port: Application services for user management
type UserUseCase interface {
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	ListUsers(ctx context.Context, companyID uuid.UUID) ([]User, error)
	CreateUser(ctx context.Context, cmd CreateUserCommand) (*User, error)
}

type CreateUserCommand struct {
	CompanyID uuid.UUID
	Email     string
	Name      string
	Role      string
}
