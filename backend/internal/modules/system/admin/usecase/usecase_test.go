package usecase_test

import (
	"context"
	"testing"

	"erp_monolith/backend/internal/modules/system/admin/domain"
	"erp_monolith/backend/internal/modules/system/admin/usecase"

	"github.com/google/uuid"
)

type mockCompanyRepo struct {
	companies []domain.Company
}

func (m *mockCompanyRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	for _, c := range m.companies {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, nil
}

func (m *mockCompanyRepo) FindAll(ctx context.Context) ([]domain.Company, error) {
	return m.companies, nil
}

func (m *mockCompanyRepo) Save(ctx context.Context, company *domain.Company) error {
	m.companies = append(m.companies, *company)
	return nil
}

func TestCompanyUseCase_CreateCompany(t *testing.T) {
	repo := &mockCompanyRepo{}
	uc := usecase.NewCompanyUseCase(repo)

	cmd := domain.CreateCompanyCommand{
		Code:     "CMP-001",
		Name:     "Acme Corporation",
		Currency: "IDR",
	}

	company, err := uc.CreateCompany(context.Background(), cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if company.Code != cmd.Code || company.Name != cmd.Name {
		t.Errorf("company fields do not match command")
	}

	all, err := uc.ListCompanies(context.Background())
	if err != nil || len(all) != 1 {
		t.Errorf("expected 1 company in repo, got %d", len(all))
	}
}
