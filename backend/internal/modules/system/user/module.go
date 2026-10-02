package user

import (
	"fmt"
	"time"

	userDelivery "erp_monolith/backend/internal/modules/system/user/delivery/http"
	userRepo "erp_monolith/backend/internal/modules/system/user/repository"
	userService "erp_monolith/backend/internal/modules/system/user/service"
	userUseCase "erp_monolith/backend/internal/modules/system/user/usecase"
	"erp_monolith/backend/internal/platform/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

// RegisterModule initializes and wires user management and auth components to Echo
func RegisterModule(e *echo.Echo, pool *pgxpool.Pool, cfg *config.Config) {
	if pool == nil {
		return
	}

	jwtSecret := "super-secret-erp-jwt-key-change-in-production"
	tokenDuration := 24 * time.Hour
	if cfg != nil {
		if cfg.Auth.JWTSecret != "" {
			jwtSecret = cfg.Auth.JWTSecret
		}
		if cfg.Auth.JWTExpirationHours > 0 {
			tokenDuration = time.Duration(cfg.Auth.JWTExpirationHours) * time.Hour
		}
	}

	tokenService := userService.NewJWTTokenService(jwtSecret, tokenDuration)
	userR := userRepo.NewPostgresUserRepo(pool)
	userUC := userUseCase.NewUserUseCase(userR, tokenService)
	userHandler := userDelivery.NewUserHandler(userUC)
	userHandler.RegisterRoutes(e)

	fmt.Println("✅ Submodule [system/user] routes wired to /api/v1/auth/signin and /api/v1/users")
}
