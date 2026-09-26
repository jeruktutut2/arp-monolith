package http

import (
	"github.com/labstack/echo/v5"
)

// RegisterRoutes mounts user routes directly to the Echo server
func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	users := e.Group("/api/v1/users")
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetUser)
	users.POST("", h.CreateUser)
}
