package dto

// CreateUserRequest DTO
type CreateUserRequest struct {
	CompanyID string `json:"company_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
}
