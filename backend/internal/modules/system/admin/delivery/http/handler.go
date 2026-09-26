package http

import (
	"net/http"

	"erp_monolith/backend/internal/modules/system/admin/delivery/http/dto"
	"erp_monolith/backend/internal/modules/system/admin/domain"
	"erp_monolith/backend/internal/platform/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type CompanyHandler struct {
	useCase domain.CompanyUseCase
}

func NewCompanyHandler(useCase domain.CompanyUseCase) *CompanyHandler {
	return &CompanyHandler{useCase: useCase}
}

func (h *CompanyHandler) ListCompanies(c *echo.Context) error {
	companies, err := h.useCase.ListCompanies(c.Request().Context())
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	return response.Success(c, http.StatusOK, "Daftar perusahaan berhasil diambil", companies)
}

func (h *CompanyHandler) GetCompany(c *echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "INVALID_ID", "ID perusahaan tidak valid")
	}

	company, err := h.useCase.GetCompany(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "NOT_FOUND", "Perusahaan tidak ditemukan")
	}

	return response.Success(c, http.StatusOK, "Data perusahaan ditemukan", company)
}

func (h *CompanyHandler) CreateCompany(c *echo.Context) error {
	var req dto.CreateCompanyRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Format JSON tidak valid")
	}

	cmd := domain.CreateCompanyCommand{
		Code:     req.Code,
		Name:     req.Name,
		Currency: req.Currency,
	}

	company, err := h.useCase.CreateCompany(c.Request().Context(), cmd)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
	}

	return response.Success(c, http.StatusCreated, "Perusahaan berhasil dibuat", company)
}
