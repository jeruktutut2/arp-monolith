package http

import (
	"github.com/labstack/echo/v5"
)

// RegisterRoutes mounts the /health route directly to the Echo instance
func (h *HealthHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/health", h.HealthCheck)
}

// RegisterGroupRoutes mounts the health route to a sub-route group (e.g. /api/v1/health)
func (h *HealthHandler) RegisterGroupRoutes(g *echo.Group) {
	g.GET("/health", h.HealthCheck)
}
