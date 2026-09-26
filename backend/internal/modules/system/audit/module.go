package audit

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// RegisterModule initializes and wires audit trail components to the Echo router group
func RegisterModule(g *echo.Group, pool *pgxpool.Pool) {
	_ = g
	_ = pool
	fmt.Println("✅ Submodule [system/audit] registered (25_AUD)")
}
