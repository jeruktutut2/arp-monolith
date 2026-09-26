package http

import (
	"github.com/labstack/echo/v5"
)

// RegisterRoutes mounts admin company routes directly to the Echo server
func (h *CompanyHandler) RegisterRoutes(e *echo.Echo) {
	companies := e.Group("/api/v1/companies")
	companies.GET("", h.ListCompanies)
	companies.GET("/:id", h.GetCompany)
	companies.POST("", h.CreateCompany)
}
