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
		if nowMilli < *user.LockedUntil {
			return nil, "", 0, "", 0, fmt.Errorf("%w: account is locked, try again later", apperrors.ErrForbidden)
		}

		txCtx, beginErr := u.txManager.Begin(ctx)
		if beginErr != nil {
			return nil, "", 0, "", 0, fmt.Errorf("begin tx: %w", beginErr)
		}

		uc := u
		rowsAffected, err := uc.userRepo.ResetFailedLogin(txCtx, user.ID.String(), nowMilli)
		if err != nil {
			resetErr := err
			_ = uc.txManager.CommitOrRollback(txCtx, resetErr)
			return nil, "", 0, "", 0, fmt.Errorf("reset failed login: %w", err)
		}

		if rowsAffected != 1 {
			err = fmt.Errorf("reset failed login: expected 1 row affected, got %d", rowsAffected)
			_ = uc.txManager.CommitOrRollback(txCtx, err)
			return nil, "", 0, "", 0, err
		}

		if commitErr := uc.txManager.CommitOrRollback(txCtx, nil); commitErr != nil {
			return nil, "", 0, "", 0, fmt.Errorf("commit reset failed login: %w", commitErr)
		}

		user.FailedLoginAttempts = 0
		user.LockedUntil = nil
	}

	return nil, "", 0, "", 0, nil
}
