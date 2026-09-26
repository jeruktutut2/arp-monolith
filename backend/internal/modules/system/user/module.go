package user

import (
	"fmt"

	userDelivery "erp_monolith/backend/internal/modules/system/user/delivery/http"
	userRepo "erp_monolith/backend/internal/modules/system/user/repository"
	userUseCase "erp_monolith/backend/internal/modules/system/user/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// RegisterModule initializes and wires user management components to the Echo server
func RegisterModule(e *echo.Echo, pool *pgxpool.Pool) {
	if pool == nil {
		return
	}

	userR := userRepo.NewPostgresUserRepo(pool)
	userUC := userUseCase.NewUserUseCase(userR)
	userHandler := userDelivery.NewUserHandler(userUC)
	userHandler.RegisterRoutes(e)

	fmt.Println("✅ Submodule [system/user] routes wired to /api/v1/users")
}
