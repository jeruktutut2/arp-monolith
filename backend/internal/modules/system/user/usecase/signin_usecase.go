package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"erp_monolith/backend/internal/modules/system/user/delivery/http/dto"
	"erp_monolith/backend/internal/modules/system/user/domain"
	"erp_monolith/backend/internal/platform/database"
	"erp_monolith/backend/internal/shared/apperrors"

	"github.com/jackc/pgx/v5"
)

type signInUseCase struct {
	userRepo  domain.UserRepository
	txManager database.TxManager
}

func NewSignInUseCase(userRepo domain.UserRepository, txManager database.TxManager) domain.SignInUseCase {
	return &signInUseCase{
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (u *signInUseCase) SignIn(ctx context.Context, req dto.SignInRequest, now time.Time) (result *domain.SignInResult, accessToken string, accessTokenExpiredAt int64, refreshToken string, refreshTokenExpiredAt int64, err error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", 0, "", 0, fmt.Errorf("%w: invalid credentials", apperrors.ErrUnauthorized)
		}
		return nil, "", 0, "", 0, fmt.Errorf("find user: %w", err)
	}

	if user == nil {
		return nil, "", 0, "", 0, fmt.Errorf("%w: invalid credentials", apperrors.ErrUnauthorized)
	}

	if !user.IsActive {
		return nil, "", 0, "", 0, fmt.Errorf("%w: account is disabled", apperrors.ErrForbidden)
	}

	if user.LockedUntil != nil {
		nowMilli := now.UnixMilli()
		if nowMilli < user.LockedUntil.UnixMilli() {
			return nil, "", 0, "", 0, fmt.Errorf("%w: account is locked, try again later", apperrors.ErrForbidden)
		}
	}

	return nil, "", 0, "", 0, nil
}
