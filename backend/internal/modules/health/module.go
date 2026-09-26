package health

import (
	"fmt"

	healthDelivery "erp_monolith/backend/internal/modules/health/check/delivery/http"
	healthDomain "erp_monolith/backend/internal/modules/health/check/domain"
	healthRepo "erp_monolith/backend/internal/modules/health/check/repository"
	healthUseCase "erp_monolith/backend/internal/modules/health/check/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// RegisterModule initializes and wires all health check components to the Echo server
func RegisterModule(e *echo.Echo, pool *pgxpool.Pool) {
	var healthChecker healthDomain.HealthChecker
	if pool != nil {
		healthChecker = healthRepo.NewPostgresHealthChecker(pool)
	}

	uc := healthUseCase.NewHealthUseCase("1.0.0", healthChecker)
	handler := healthDelivery.NewHealthHandler(uc)
	handler.RegisterRoutes(e)

	fmt.Println("✅ Module [Health Check] route wired to /health")
}
