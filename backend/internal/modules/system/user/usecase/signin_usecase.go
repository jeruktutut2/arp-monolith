package usecase

import (
	"context"

	"erp_monolith/backend/internal/modules/system/user/delivery/http/dto"
	"erp_monolith/backend/internal/modules/system/user/domain"
)

type signInUseCase struct{}

func NewSignInUseCase() domain.SignInUseCase {
	return &signInUseCase{}
}

func (u *signInUseCase) SignIn(ctx context.Context, req dto.SignInRequest) (result *domain.SignInResult, accessToken string, accessTokenExpiredAt int64, refreshToken string, refreshTokenExpiredAt int64, err error) {
	return nil, "", 0, "", 0, nil
}
