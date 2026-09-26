package dto

// CreateCompanyRequest DTO
type CreateCompanyRequest struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}
