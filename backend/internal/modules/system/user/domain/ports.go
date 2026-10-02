package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Inbound Port: Application service interface for user management and authentication
type UserUseCase interface {
	SignIn(ctx context.Context, cmd SignInCommand) (*SignInResult, error)
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	ListUsers(ctx context.Context, companyID uuid.UUID) ([]User, error)
	CreateUser(ctx context.Context, cmd CreateUserCommand) (*User, error)
}

// SignInCommand carries credentials to authenticate
type SignInCommand struct {
	Email    string
	Password string
}

// SignInResult contains the issued token and authenticated user details
type SignInResult struct {
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
	User      *User     `json:"user"`
}

// CreateUserCommand carries user creation parameters
type CreateUserCommand struct {
	CompanyID uuid.UUID
	BranchID  *uuid.UUID
	Email     string
	Password  string
	Name      string
	Role      string
}
