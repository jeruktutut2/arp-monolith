package http

import (
	"github.com/labstack/echo/v5"
)

// RegisterRoutes mounts user routes to the Echo v5 group
func (h *UserHandler) RegisterRoutes(g *echo.Group) {
	users := g.Group("/users")
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetUser)
	users.POST("", h.CreateUser)
}
