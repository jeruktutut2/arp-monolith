package usecase

import (
	"context"
	"time"

	"erp_monolith/backend/internal/modules/system/user/domain"

	"github.com/google/uuid"
)

type userUseCase struct {
	repo domain.UserRepository
}

// NewUserUseCase creates an instance of UserUseCase
func NewUserUseCase(repo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *userUseCase) ListUsers(ctx context.Context, companyID uuid.UUID) ([]domain.User, error) {
	return u.repo.FindByCompanyID(ctx, companyID)
}

func (u *userUseCase) CreateUser(ctx context.Context, cmd domain.CreateUserCommand) (*domain.User, error) {
	now := time.Now().UTC()
	user := &domain.User{
		ID:        uuid.New(),
		CompanyID: cmd.CompanyID,
		Email:     cmd.Email,
		Name:      cmd.Name,
		Role:      cmd.Role,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
