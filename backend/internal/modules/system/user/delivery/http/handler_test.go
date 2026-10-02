package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	userDelivery "erp_monolith/backend/internal/modules/system/user/delivery/http"
	"erp_monolith/backend/internal/modules/system/user/domain"
	"erp_monolith/backend/internal/platform/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type mockUserUseCase struct {
	signInFunc     func(ctx context.Context, cmd domain.SignInCommand) (*domain.SignInResult, error)
	getUserFunc    func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	listUsersFunc  func(ctx context.Context, companyID uuid.UUID) ([]domain.User, error)
	createUserFunc func(ctx context.Context, cmd domain.CreateUserCommand) (*domain.User, error)
}

func (m *mockUserUseCase) SignIn(ctx context.Context, cmd domain.SignInCommand) (*domain.SignInResult, error) {
	if m.signInFunc != nil {
		return m.signInFunc(ctx, cmd)
	}
	return nil, nil
}

func (m *mockUserUseCase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserUseCase) ListUsers(ctx context.Context, companyID uuid.UUID) ([]domain.User, error) {
	if m.listUsersFunc != nil {
		return m.listUsersFunc(ctx, companyID)
	}
	return nil, nil
}

func (m *mockUserUseCase) CreateUser(ctx context.Context, cmd domain.CreateUserCommand) (*domain.User, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, cmd)
	}
	return nil, nil
}

func TestUserHandler_SignIn_Success(t *testing.T) {
	e := echo.New()
	userID := uuid.New()
	companyID := uuid.New()

	mockUC := &mockUserUseCase{
		signInFunc: func(ctx context.Context, cmd domain.SignInCommand) (*domain.SignInResult, error) {
			return &domain.SignInResult{
				Token:     "test.jwt.token",
				TokenType: "Bearer",
				ExpiresAt: time.Now().Add(24 * time.Hour),
				User: &domain.User{
					ID:        userID,
					CompanyID: companyID,
					Email:     cmd.Email,
					Name:      "Admin User",
					Role:      "admin",
					IsActive:  true,
					CreatedAt: time.Now(),
				},
			}, nil
		},
	}

	h := userDelivery.NewUserHandler(mockUC)

	reqBody := `{"email":"admin@example.com","password":"secretpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signin", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.SignIn(c)
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp response.StandardResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success true")
	}
}

func TestUserHandler_SignIn_InvalidCredentials(t *testing.T) {
	e := echo.New()

	mockUC := &mockUserUseCase{
		signInFunc: func(ctx context.Context, cmd domain.SignInCommand) (*domain.SignInResult, error) {
			return nil, domain.ErrInvalidCredentials
		},
	}

	h := userDelivery.NewUserHandler(mockUC)

	reqBody := `{"email":"admin@example.com","password":"wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signin", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.SignIn(c)
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized, got %d", rec.Code)
	}

	var resp response.StandardResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if resp.Code != "INVALID_CREDENTIALS" {
		t.Errorf("expected code INVALID_CREDENTIALS, got %s", resp.Code)
	}
}

func TestUserHandler_SignIn_AccountLocked(t *testing.T) {
	e := echo.New()

	mockUC := &mockUserUseCase{
		signInFunc: func(ctx context.Context, cmd domain.SignInCommand) (*domain.SignInResult, error) {
			return nil, domain.ErrUserLocked
		},
	}

	h := userDelivery.NewUserHandler(mockUC)

	reqBody := `{"email":"locked@example.com","password":"anypassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signin", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.SignIn(c)
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}

	if rec.Code != http.StatusLocked {
		t.Fatalf("expected status 423 Locked, got %d", rec.Code)
	}
}
