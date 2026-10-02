package http

import (
	"errors"
	"net/http"

	"erp_monolith/backend/internal/modules/system/user/delivery/http/dto"
	"erp_monolith/backend/internal/modules/system/user/domain"
	"erp_monolith/backend/internal/platform/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	useCase domain.UserUseCase
}

// NewUserHandler creates a driving adapter handler for users and authentication
func NewUserHandler(useCase domain.UserUseCase) *UserHandler {
	return &UserHandler{useCase: useCase}
}

// SignIn handles POST /api/v1/auth/signin & /api/v1/users/signin
func (h *UserHandler) SignIn(c *echo.Context) error {
	var req dto.SignInRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Format JSON request tidak valid")
	}

	result, err := h.useCase.SignIn(c.Request().Context(), domain.SignInCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailRequired):
			return response.Error(c, http.StatusBadRequest, "EMAIL_REQUIRED", err.Error())
		case errors.Is(err, domain.ErrInvalidEmail):
			return response.Error(c, http.StatusBadRequest, "INVALID_EMAIL", err.Error())
		case errors.Is(err, domain.ErrPasswordRequired):
			return response.Error(c, http.StatusBadRequest, "PASSWORD_REQUIRED", err.Error())
		case errors.Is(err, domain.ErrInvalidCredentials):
			return response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
		case errors.Is(err, domain.ErrUserInactive):
			return response.Error(c, http.StatusForbidden, "USER_INACTIVE", err.Error())
		case errors.Is(err, domain.ErrUserLocked):
			return response.Error(c, http.StatusLocked, "USER_LOCKED", err.Error())
		default:
			return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan internal pada server")
		}
	}

	respData := dto.SignInResponse{
		Token:     result.Token,
		TokenType: result.TokenType,
		ExpiresAt: result.ExpiresAt,
		User:      dto.ToUserResponse(result.User),
	}

	return response.Success(c, http.StatusOK, "Login berhasil", respData)
}

func (h *UserHandler) ListUsers(c *echo.Context) error {
	companyIDStr := c.QueryParam("company_id")
	if companyIDStr == "" {
		return response.Error(c, http.StatusBadRequest, "MISSING_PARAM", "Parameter query company_id wajib diisi")
	}
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "INVALID_ID", "company_id tidak valid")
	}

	users, err := h.useCase.ListUsers(c.Request().Context(), companyID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	var dtos []dto.UserResponseData
	for i := range users {
		dtos = append(dtos, dto.ToUserResponse(&users[i]))
	}

	return response.Success(c, http.StatusOK, "Daftar pengguna berhasil diambil", dtos)
}

func (h *UserHandler) GetUser(c *echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID pengguna tidak valid")
	}

	user, err := h.useCase.GetUser(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "NOT_FOUND", "Pengguna tidak ditemukan")
	}

	return response.Success(c, http.StatusOK, "Data pengguna ditemukan", dto.ToUserResponse(user))
}

func (h *UserHandler) CreateUser(c *echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Format JSON tidak valid")
	}

	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "INVALID_COMPANY_ID", "company_id tidak valid")
	}

	var branchID *uuid.UUID
	if req.BranchID != nil && *req.BranchID != "" {
		parsed, err := uuid.Parse(*req.BranchID)
		if err == nil {
			branchID = &parsed
		}
	}

	cmd := domain.CreateUserCommand{
		CompanyID: companyID,
		BranchID:  branchID,
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		Role:      req.Role,
	}

	user, err := h.useCase.CreateUser(c.Request().Context(), cmd)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
	}

	return response.Success(c, http.StatusCreated, "Pengguna berhasil dibuat", dto.ToUserResponse(user))
}
