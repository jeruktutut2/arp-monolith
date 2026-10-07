package usecase

import (
	"erp_monolith/backend/internal/modules/system/user/domain"
)

type signInUseCase struct{}

func NewSignInUseCase() domain.SignInUseCase {
	return &signInUseCase{}
}
