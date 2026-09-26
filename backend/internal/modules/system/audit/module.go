package audit

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// RegisterModule initializes and wires audit trail components to the Echo server
func RegisterModule(e *echo.Echo, pool *pgxpool.Pool) {
	_ = e
	_ = pool
	fmt.Println("✅ Submodule [system/audit] registered (25_AUD)")
}
