package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"erp_monolith/backend/internal/modules/system/user/domain"
	"erp_monolith/backend/internal/modules/system/user/service"
	"erp_monolith/backend/internal/modules/system/user/usecase"

	"github.com/google/uuid"
)

type mockUserRepo struct {
	users map[uuid.UUID]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[uuid.UUID]*domain.User),
	}
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	cp := *u
	return &cp, nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, u := range m.users {
		if strings.EqualFold(u.Email, email) {
			cp := *u
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) FindByCompanyID(ctx context.Context, companyID uuid.UUID) ([]domain.User, error) {
	var res []domain.User
	for _, u := range m.users {
		if u.CompanyID == companyID {
			res = append(res, *u)
		}
	}
	return res, nil
}

func (m *mockUserRepo) Save(ctx context.Context, user *domain.User) error {
	cp := *user
	m.users[user.ID] = &cp
	return nil
}

func (m *mockUserRepo) UpdateLoginAttempts(ctx context.Context, id uuid.UUID, attempts int, lockedUntil *time.Time) error {
	u, ok := m.users[id]
	if !ok {
		return errors.New("user not found")
	}
	u.FailedLoginAttempts = attempts
	u.LockedUntil = lockedUntil
	return nil
}

func setupTest(t *testing.T) (domain.UserUseCase, *mockUserRepo, domain.TokenService) {
	t.Helper()
	repo := newMockUserRepo()
	tokenService := service.NewJWTTokenService("test-secret-key-12345", 2*time.Hour)
	uc := usecase.NewUserUseCase(repo, tokenService)
	return uc, repo, tokenService
}

func TestUserUseCase_CreateAndListUser(t *testing.T) {
	uc, _, _ := setupTest(t)

	companyID := uuid.New()
	cmd := domain.CreateUserCommand{
		CompanyID: companyID,
		Email:     "admin@acme.com",
		Password:  "password123",
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
	if !user.ValidatePassword("password123") {
		t.Errorf("password was not hashed properly")
	}

	users, err := uc.ListUsers(context.Background(), companyID)
	if err != nil || len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestUserUseCase_SignIn_Success(t *testing.T) {
	uc, repo, tokenService := setupTest(t)

	// Persist active user
	user := &domain.User{
		ID:        uuid.New(),
		CompanyID: uuid.New(),
		Email:     "user@enterprise.com",
		Name:      "John Doe",
		Role:      "admin",
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := user.SetPassword("SecurePass123!"); err != nil {
		t.Fatalf("failed to set password: %v", err)
	}
	_ = repo.Save(context.Background(), user)

	// Attempt sign in
	result, err := uc.SignIn(context.Background(), domain.SignInCommand{
		Email:    "User@Enterprise.com", // test case-insensitivity
		Password: "SecurePass123!",
	})
	if err != nil {
		t.Fatalf("expected successful sign in, got: %v", err)
	}

	if result.Token == "" {
		t.Errorf("expected non-empty JWT token")
	}
	if result.TokenType != "Bearer" {
		t.Errorf("expected token type Bearer, got %s", result.TokenType)
	}
	if result.User.Email != user.Email {
		t.Errorf("expected user email %s, got %s", user.Email, result.User.Email)
	}

	// Verify token validity
	claims, err := tokenService.ValidateToken(result.Token)
	if err != nil {
		t.Fatalf("expected valid token, got: %v", err)
	}
	if claims.UserID != user.ID || claims.Role != "admin" {
		t.Errorf("claims mismatch: %+v", claims)
	}
}

func TestUserUseCase_SignIn_ValidationErrors(t *testing.T) {
	uc, _, _ := setupTest(t)
	ctx := context.Background()

	// 1. Email empty
	_, err := uc.SignIn(ctx, domain.SignInCommand{Email: "", Password: "secret"})
	if !errors.Is(err, domain.ErrEmailRequired) {
		t.Errorf("expected ErrEmailRequired, got %v", err)
	}

	// 2. Email invalid format
	_, err = uc.SignIn(ctx, domain.SignInCommand{Email: "invalid-email", Password: "secret"})
	if !errors.Is(err, domain.ErrInvalidEmail) {
		t.Errorf("expected ErrInvalidEmail, got %v", err)
	}

	// 3. Password empty
	_, err = uc.SignIn(ctx, domain.SignInCommand{Email: "user@test.com", Password: ""})
	if !errors.Is(err, domain.ErrPasswordRequired) {
		t.Errorf("expected ErrPasswordRequired, got %v", err)
	}
}

func TestUserUseCase_SignIn_InvalidCredentials_UserNotFound(t *testing.T) {
	uc, _, _ := setupTest(t)

	_, err := uc.SignIn(context.Background(), domain.SignInCommand{
		Email:    "nonexistent@enterprise.com",
		Password: "password",
	})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUserUseCase_SignIn_WrongPassword_And_Lockout(t *testing.T) {
	uc, repo, _ := setupTest(t)
	ctx := context.Background()

	user := &domain.User{
		ID:        uuid.New(),
		CompanyID: uuid.New(),
		Email:     "target@enterprise.com",
		Name:      "Target User",
		Role:      "staff",
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = user.SetPassword("CorrectPassword123!")
	_ = repo.Save(ctx, user)

	// Attempt 1: Wrong password
	_, err := uc.SignIn(ctx, domain.SignInCommand{Email: user.Email, Password: "wrong-password-1"})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("attempt 1: expected ErrInvalidCredentials, got %v", err)
	}
	if repo.users[user.ID].FailedLoginAttempts != 1 {
		t.Errorf("expected 1 failed attempt, got %d", repo.users[user.ID].FailedLoginAttempts)
	}

	// Attempt 2: Wrong password
	_, err = uc.SignIn(ctx, domain.SignInCommand{Email: user.Email, Password: "wrong-password-2"})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("attempt 2: expected ErrInvalidCredentials, got %v", err)
	}
	if repo.users[user.ID].FailedLoginAttempts != 2 {
		t.Errorf("expected 2 failed attempts, got %d", repo.users[user.ID].FailedLoginAttempts)
	}

	// Attempt 3: Wrong password -> locks account
	_, err = uc.SignIn(ctx, domain.SignInCommand{Email: user.Email, Password: "wrong-password-3"})
	if !errors.Is(err, domain.ErrUserLocked) {
		t.Errorf("attempt 3: expected ErrUserLocked, got %v", err)
	}
	if repo.users[user.ID].FailedLoginAttempts != 3 || repo.users[user.ID].LockedUntil == nil {
		t.Errorf("expected user to be locked out after 3 attempts")
	}

	// Attempt 4: Even with correct password, user remains locked
	_, err = uc.SignIn(ctx, domain.SignInCommand{Email: user.Email, Password: "CorrectPassword123!"})
	if !errors.Is(err, domain.ErrUserLocked) {
		t.Errorf("attempt 4: expected ErrUserLocked, got %v", err)
	}
}

func TestUserUseCase_SignIn_InactiveUser(t *testing.T) {
	uc, repo, _ := setupTest(t)

	user := &domain.User{
		ID:        uuid.New(),
		CompanyID: uuid.New(),
		Email:     "inactive@enterprise.com",
		Name:      "Inactive User",
		Role:      "staff",
		IsActive:  false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = user.SetPassword("Password123!")
	_ = repo.Save(context.Background(), user)

	_, err := uc.SignIn(context.Background(), domain.SignInCommand{
		Email:    user.Email,
		Password: "Password123!",
	})
	if !errors.Is(err, domain.ErrUserInactive) {
		t.Errorf("expected ErrUserInactive, got %v", err)
	}
}

func TestUserUseCase_SignIn_ResetsFailedAttemptsOnSuccess(t *testing.T) {
	uc, repo, _ := setupTest(t)
	ctx := context.Background()

	user := &domain.User{
		ID:                  uuid.New(),
		CompanyID:           uuid.New(),
		Email:               "recovering@enterprise.com",
		Name:                "Recovering User",
		Role:                "staff",
		IsActive:            true,
		FailedLoginAttempts: 2,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}
	_ = user.SetPassword("CorrectPassword123!")
	_ = repo.Save(ctx, user)

	// Successful login resets failed attempts
	res, err := uc.SignIn(ctx, domain.SignInCommand{
		Email:    user.Email,
		Password: "CorrectPassword123!",
	})
	if err != nil {
		t.Fatalf("expected successful sign in, got %v", err)
	}
	if res.Token == "" {
		t.Errorf("expected token")
	}
	if repo.users[user.ID].FailedLoginAttempts != 0 {
		t.Errorf("expected failed attempts to reset to 0, got %d", repo.users[user.ID].FailedLoginAttempts)
	}
}
