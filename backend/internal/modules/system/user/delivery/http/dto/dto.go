package dto

import (
	"time"

	"erp_monolith/backend/internal/modules/system/user/domain"
)

// SignInRequest carries email and password for authentication
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignInResponse returns the authentication token and authenticated user profile
type SignInResponse struct {
	Token     string           `json:"token"`
	TokenType string           `json:"token_type"`
	ExpiresAt time.Time        `json:"expires_at"`
	User      UserResponseData `json:"user"`
}

// CreateUserRequest carries payload to register a user
type CreateUserRequest struct {
	CompanyID string  `json:"company_id"`
	BranchID  *string `json:"branch_id,omitempty"`
	Email     string  `json:"email"`
	Password  string  `json:"password,omitempty"`
	Name      string  `json:"name"`
	Role      string  `json:"role"`
}

// UserResponseData is the public representation of a user without sensitive credentials
type UserResponseData struct {
	ID        string  `json:"id"`
	CompanyID string  `json:"company_id"`
	BranchID  *string `json:"branch_id,omitempty"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	Role      string  `json:"role"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
}

// ToUserResponse converts domain.User to UserResponseData safely
func ToUserResponse(u *domain.User) UserResponseData {
	var branchIDStr *string
	if u.BranchID != nil {
		s := u.BranchID.String()
		branchIDStr = &s
	}

	return UserResponseData{
		ID:        u.ID.String(),
		CompanyID: u.CompanyID.String(),
		BranchID:  branchIDStr,
		Email:     u.Email,
		Name:      u.Name,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}
