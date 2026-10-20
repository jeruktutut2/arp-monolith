package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"erp_monolith/backend/internal/modules/system/user/delivery/http/dto"
	"erp_monolith/backend/internal/modules/system/user/domain"
	"erp_monolith/backend/internal/platform/database"
	"erp_monolith/backend/internal/platform/sonyflake"
	"erp_monolith/backend/internal/platform/uuid"
	"erp_monolith/backend/internal/shared/apperrors"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	maxFailedAttempts = 5
	lockDuration      = 30 * time.Minute
)

type signInUseCase struct {
	userRepo             domain.UserRepository
	txManager            database.TxManager
	uuid                 uuid.Uuid
	sonyflake            sonyflake.Sonyflake
	accessTokenDuration  int64
	refreshTokenDuration int64
}

func NewSignInUseCase(userRepo domain.UserRepository, txManager database.TxManager, uuid uuid.Uuid, sonyflake sonyflake.Sonyflake, accessTokenDuration int64, refreshTokenDuration int64) domain.SignInUseCase {
	return &signInUseCase{
		userRepo:             userRepo,
		txManager:            txManager,
		uuid:                 uuid,
		sonyflake:            sonyflake,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
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

	nowMilli := now.UnixMilli()

	if user.LockedUntil != nil {
		if nowMilli < *user.LockedUntil {
			return nil, "", 0, "", 0, fmt.Errorf("%w: account is locked, try again later", apperrors.ErrForbidden)
		}

		txCtx, beginErr := u.txManager.Begin(ctx)
		if beginErr != nil {
			return nil, "", 0, "", 0, fmt.Errorf("begin tx: %w", beginErr)
		}

		rowsAffected, err := u.userRepo.ResetFailedLogin(txCtx, strconv.FormatInt(user.ID, 10), nowMilli)
		if err != nil {
			resetErr := err
			_ = u.txManager.CommitOrRollback(txCtx, resetErr)
			return nil, "", 0, "", 0, fmt.Errorf("reset failed login: %w", err)
		}

		if rowsAffected != 1 {
			err = fmt.Errorf("reset failed login: expected 1 row affected, got %d", rowsAffected)
			_ = u.txManager.CommitOrRollback(txCtx, err)
			return nil, "", 0, "", 0, err
		}

		if commitErr := u.txManager.CommitOrRollback(txCtx, nil); commitErr != nil {
			return nil, "", 0, "", 0, fmt.Errorf("commit reset failed login: %w", commitErr)
		}

		user.FailedLoginAttempts = 0
		user.LockedUntil = nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		newCount := user.FailedLoginAttempts + 1
		var lockedUntil *int64
		if newCount >= maxFailedAttempts {
			lockTime := nowMilli + lockDuration.Milliseconds()
			lockedUntil = &lockTime
		}

		txCtx, err := u.txManager.Begin(ctx)
		if err != nil {
			return nil, "", 0, "", 0, fmt.Errorf("begin tx: %w", err)
		}

		rowsAffected, err := u.userRepo.UpdateFailedLogin(txCtx, strconv.FormatInt(user.ID, 10), newCount, lockedUntil, nowMilli)
		if err != nil {
			_ = u.txManager.CommitOrRollback(txCtx, err)
			return nil, "", 0, "", 0, fmt.Errorf("update failed login: %w", err)
		}

		if rowsAffected != 1 {
			err = fmt.Errorf("update failed login: expected 1 row affected, got %d", rowsAffected)
			_ = u.txManager.CommitOrRollback(txCtx, err)
			return nil, "", 0, "", 0, err
		}

		if commitErr := u.txManager.CommitOrRollback(txCtx, nil); commitErr != nil {
			return nil, "", 0, "", 0, fmt.Errorf("commit update failed login: %w", commitErr)
		}

		return nil, "", 0, "", 0, fmt.Errorf("%w: invalid credentials", apperrors.ErrUnauthorized)
	}

	userSessionID, err := u.sonyflake.GenerateID()
	if err != nil {
		return nil, "", 0, "", 0, fmt.Errorf("generate session id: %w", err)
	}

	accessToken = u.uuid.GenerateAccessToken()
	accessTokenExpiredAt = nowMilli + u.accessTokenDuration
	refreshToken = u.uuid.GenerateRefreshToken()
	refreshTokenExpiredAt = nowMilli + u.refreshTokenDuration

	userSession := &domain.UserSession{
		ID:                    userSessionID,
		UserID:                user.ID,
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  accessTokenExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshTokenExpiredAt,
		IsRevoked:             false,
		CreatedAt:             nowMilli,
		CreatedBy:             strconv.FormatInt(user.ID, 10),
		UpdatedAt:             nil,
		UpdatedBy:             nil,
	}
	_ = userSession

	return nil, accessToken, accessTokenExpiredAt, refreshToken, refreshTokenExpiredAt, nil
}
