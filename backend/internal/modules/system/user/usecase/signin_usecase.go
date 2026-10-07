package usecase

import (
	"context"

	"erp_monolith/backend/internal/modules/system/user/delivery/http/dto"
)

type SignIn interface {
	SignIn(ctx context.Context, req dto.SignInRequest) (*dto.SignInResponse, error)
}
