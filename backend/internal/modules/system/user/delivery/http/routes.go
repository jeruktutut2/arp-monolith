package http

import (
	"github.com/labstack/echo/v5"
)

// RegisterRoutes mounts user and authentication routes to the Echo server
func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	// Authentication routes
	auth := e.Group("/api/v1/auth")
	auth.POST("/signin", h.SignIn)

	// User management routes
	users := e.Group("/api/v1/users")
	users.POST("/signin", h.SignIn) // convenient alias
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetUser)
	users.POST("", h.CreateUser)
}
