package system

import (
	"erp_monolith/backend/internal/modules/system/admin"
	"erp_monolith/backend/internal/modules/system/audit"
	"erp_monolith/backend/internal/modules/system/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// RegisterModule initializes and registers all system submodules (admin, user, audit)
func RegisterModule(e *echo.Echo, pool *pgxpool.Pool) {
	admin.RegisterModule(e, pool)
	user.RegisterModule(e, pool)
	audit.RegisterModule(e, pool)
}
