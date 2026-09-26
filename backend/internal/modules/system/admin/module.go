package admin

import (
	"fmt"

	adminDelivery "erp_monolith/backend/internal/modules/system/admin/delivery/http"
	adminRepo "erp_monolith/backend/internal/modules/system/admin/repository"
	adminUseCase "erp_monolith/backend/internal/modules/system/admin/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// RegisterModule initializes and wires admin/multi-company components to the Echo router group
func RegisterModule(g *echo.Group, pool *pgxpool.Pool) {
	if pool == nil {
		return
	}

	companyRepo := adminRepo.NewPostgresCompanyRepo(pool)
	companyUC := adminUseCase.NewCompanyUseCase(companyRepo)
	companyHandler := adminDelivery.NewCompanyHandler(companyUC)
	companyHandler.RegisterRoutes(g)

	fmt.Println("✅ Submodule [system/admin] routes wired to /api/v1/companies")
}
