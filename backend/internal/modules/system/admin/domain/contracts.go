package domain

import (
	"context"

	"github.com/google/uuid"
)

// Outbound Port: Repository interface required by system admin domain
type CompanyRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Company, error)
	FindAll(ctx context.Context) ([]Company, error)
	Save(ctx context.Context, company *Company) error
}
