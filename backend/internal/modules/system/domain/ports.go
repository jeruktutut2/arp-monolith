package domain

import (
	"context"

	"github.com/google/uuid"
)

// Inbound Port: Application services exposed to driving adapters (HTTP / CLI)
type CompanyUseCase interface {
	GetCompany(ctx context.Context, id uuid.UUID) (*Company, error)
	ListCompanies(ctx context.Context) ([]Company, error)
	CreateCompany(ctx context.Context, cmd CreateCompanyCommand) (*Company, error)
}

type CreateCompanyCommand struct {
	Code     string
	Name     string
	Currency string
}
