package http

import (
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

func NewUserHandler(useCase domain.UserUseCase) *UserHandler {
	return &UserHandler{useCase: useCase}
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
	return response.Success(c, http.StatusOK, "Daftar pengguna berhasil diambil", users)
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

	return response.Success(c, http.StatusOK, "Data pengguna ditemukan", user)
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

	cmd := domain.CreateUserCommand{
		CompanyID: companyID,
		Email:     req.Email,
		Name:      req.Name,
		Role:      req.Role,
	}

	user, err := h.useCase.CreateUser(c.Request().Context(), cmd)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
	}

	return response.Success(c, http.StatusCreated, "Pengguna berhasil dibuat", user)
}
