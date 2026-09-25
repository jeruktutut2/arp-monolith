package usecase

import (
	"context"
	"erp_monolith/backend/internal/modules/system/domain"
	"time"

	"github.com/google/uuid"
)

type companyUseCase struct {
	repo domain.CompanyRepository
}

// NewCompanyUseCase creates an instance of CompanyUseCase
func NewCompanyUseCase(repo domain.CompanyRepository) domain.CompanyUseCase {
	return &companyUseCase{repo: repo}
}

func (u *companyUseCase) GetCompany(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *companyUseCase) ListCompanies(ctx context.Context) ([]domain.Company, error) {
	return u.repo.FindAll(ctx)
}

func (u *companyUseCase) CreateCompany(ctx context.Context, cmd domain.CreateCompanyCommand) (*domain.Company, error) {
	now := time.Now().UTC()
	company := &domain.Company{
		ID:        uuid.New(),
		Code:      cmd.Code,
		Name:      cmd.Name,
		Currency:  cmd.Currency,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.repo.Save(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}
