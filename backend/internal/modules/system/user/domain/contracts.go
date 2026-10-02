package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserRepository is the outbound port for persisting and retrieving user entities
type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByCompanyID(ctx context.Context, companyID uuid.UUID) ([]User, error)
	Save(ctx context.Context, user *User) error
	UpdateLoginAttempts(ctx context.Context, id uuid.UUID, attempts int, lockedUntil *time.Time) error
}

// TokenService is the outbound port for issuing and verifying JWT tokens
type TokenService interface {
	GenerateToken(user *User) (token string, expiresAt time.Time, err error)
	ValidateToken(tokenStr string) (*Claims, error)
}

// Claims encapsulates identity data within the JWT
type Claims struct {
	UserID    uuid.UUID  `json:"uid"`
	CompanyID uuid.UUID  `json:"cid"`
	BranchID  *uuid.UUID `json:"bid,omitempty"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
}
