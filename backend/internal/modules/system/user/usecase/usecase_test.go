package usecase_test

import (
	"context"
	"testing"

	"erp_monolith/backend/internal/modules/system/user/domain"
	"erp_monolith/backend/internal/modules/system/user/usecase"

	"github.com/google/uuid"
)

type mockUserRepo struct {
	users []domain.User
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) FindByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.User, error) {
	var res []domain.User
	for _, u := range m.users {
		if u.CompanyID == companyID {
			res = append(res, u)
		}
	}
	return res, nil
}

func (m *mockUserRepo) Save(ctx context.Context, user *domain.User) error {
	m.users = append(m.users, *user)
	return nil
}

func TestUserUseCase_CreateAndListUser(t *testing.T) {
	repo := &mockUserRepo{}
	uc := usecase.NewUserUseCase(repo)

	companyID := uuid.New()
	cmd := domain.CreateUserCommand{
		CompanyID: companyID,
		Email:     "admin@acme.com",
		Name:      "System Admin",
		Role:      "SuperAdmin",
	}

	user, err := uc.CreateUser(context.Background(), cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != cmd.Email || user.Name != cmd.Name {
		t.Errorf("user fields do not match command")
	}

	users, err := uc.ListUsers(context.Background(), companyID)
	if err != nil || len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}
