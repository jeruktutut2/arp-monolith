package usecase

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"erp_monolith/backend/internal/modules/system/user/domain"

	"github.com/google/uuid"
)

type userUseCase struct {
	repo         domain.UserRepository
	tokenService domain.TokenService
}

// NewUserUseCase creates an instance of UserUseCase
func NewUserUseCase(repo domain.UserRepository, tokenService domain.TokenService) domain.UserUseCase {
	return &userUseCase{
		repo:         repo,
		tokenService: tokenService,
	}
}

const (
	MaxLoginAttempts = 3
	LockoutDuration  = 30 * time.Minute
)

func (u *userUseCase) SignIn(ctx context.Context, cmd domain.SignInCommand) (*domain.SignInResult, error) {
	// 1. Validasi input
	email := strings.TrimSpace(cmd.Email)
	if email == "" {
		return nil, domain.ErrEmailRequired
	}

	if _, err := mail.ParseAddress(email); err != nil || !strings.Contains(email, ".") {
		return nil, domain.ErrInvalidEmail
	}

	if cmd.Password == "" {
		return nil, domain.ErrPasswordRequired
	}

	// 2. Cari pengguna berdasarkan email
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// Generic error untuk proteksi email enumeration
		return nil, domain.ErrInvalidCredentials
	}

	now := time.Now().UTC()

	// 3. Cek apakah akun sedang terkunci (lockout)
	if user.IsLocked(now) {
		return nil, domain.ErrUserLocked
	}

	// 4. Cek apakah akun aktif
	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	// 5. Validasi password hash
	if !user.ValidatePassword(cmd.Password) {
		user.RecordFailedLogin(now, MaxLoginAttempts, LockoutDuration)
		_ = u.repo.UpdateLoginAttempts(ctx, user.ID, user.FailedLoginAttempts, user.LockedUntil)

		if user.IsLocked(now) {
			return nil, domain.ErrUserLocked
		}
		return nil, domain.ErrInvalidCredentials
	}

	// 6. Reset percobaan login gagal jika berhasil
	if user.FailedLoginAttempts > 0 || user.LockedUntil != nil {
		user.ResetLoginAttempts()
		_ = u.repo.UpdateLoginAttempts(ctx, user.ID, 0, nil)
	}

	// 7. Terbitkan token otentikasi (JWT)
	token, expiresAt, err := u.tokenService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.SignInResult{
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

func (u *userUseCase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (u *userUseCase) ListUsers(ctx context.Context, companyID uuid.UUID) ([]domain.User, error) {
	return u.repo.FindByCompanyID(ctx, companyID)
}

func (u *userUseCase) CreateUser(ctx context.Context, cmd domain.CreateUserCommand) (*domain.User, error) {
	now := time.Now().UTC()
	user := &domain.User{
		ID:        uuid.New(),
		CompanyID: cmd.CompanyID,
		BranchID:  cmd.BranchID,
		Email:     strings.ToLower(strings.TrimSpace(cmd.Email)),
		Name:      cmd.Name,
		Role:      cmd.Role,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if cmd.Password != "" {
		if err := user.SetPassword(cmd.Password); err != nil {
			return nil, err
		}
	}

	if err := u.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
