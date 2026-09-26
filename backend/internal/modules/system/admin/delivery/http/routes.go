package http

import (
	"github.com/labstack/echo/v5"
)

// RegisterRoutes mounts admin company routes to the Echo v5 group
func (h *CompanyHandler) RegisterRoutes(g *echo.Group) {
	companies := g.Group("/companies")
	companies.GET("", h.ListCompanies)
	companies.GET("/:id", h.GetCompany)
	companies.POST("", h.CreateCompany)
}
