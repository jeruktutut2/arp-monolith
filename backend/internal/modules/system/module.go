package system

import (
	"erp_monolith/backend/internal/modules/system/admin"
	"erp_monolith/backend/internal/modules/system/audit"
	"erp_monolith/backend/internal/modules/system/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// RegisterModule initializes and registers all system submodules (admin, user, audit)
func RegisterModule(g *echo.Group, pool *pgxpool.Pool) {
	admin.RegisterModule(g, pool)
	user.RegisterModule(g, pool)
	audit.RegisterModule(g, pool)
}
