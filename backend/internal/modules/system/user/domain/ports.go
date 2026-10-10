package domain

import (
	"context"
	"time"

	"erp_monolith/backend/internal/modules/system/user/delivery/http/dto"

	"github.com/google/uuid"
)

type SignInUseCase interface {
	SignIn(ctx context.Context, req dto.SignInRequest, now time.Time) (result *SignInResult, accessToken string, accessTokenExpiredAt int64, refreshToken string, refreshTokenExpiredAt int64, err error)
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	ResetFailedLogin(ctx context.Context, userID string, now int64) (int64, error)
}

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

// SignInResult carries the authenticated user identity
type SignInResult struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Old SignInResult (commented out)
// type SignInResult struct {
// 	Token     string    `json:"token"`
// 	TokenType string    `json:"token_type"`
// 	ExpiresAt time.Time `json:"expires_at"`
// 	User      *User     `json:"user"`
// }

// CreateUserCommand carries user creation parameters
type CreateUserCommand struct {
	CompanyID uuid.UUID
	BranchID  *uuid.UUID
	Email     string
	Password  string
	Name      string
	Role      string
}
