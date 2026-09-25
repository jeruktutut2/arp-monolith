package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"erp_monolith/backend/internal/platform/database"
	"erp_monolith/backend/internal/platform/eventbus"
	sysDelivery "erp_monolith/backend/internal/modules/system/delivery/http"
	sysRepo "erp_monolith/backend/internal/modules/system/repository"
	sysUseCase "erp_monolith/backend/internal/modules/system/usecase"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	fmt.Println("🚀 Starting Enterprise ERP Monolith Backend...")

	// 1. Environment & Config
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default to PgBouncer port 6432
		dbURL = "postgres://erp_user:erp_secret@localhost:6432/erp_db?sslmode=disable"
	}

	// 2. Database Pool via PgBouncer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := database.NewPool(ctx, database.Config{
		URL: dbURL,
	})
	if err != nil {
		fmt.Printf("❌ Failed to initialize database pool: %v\n", err)
	} else {
		defer dbPool.Close()
		fmt.Println("✅ PgBouncer connection pool initialized (:6432)")
	}

	// 3. EventBus In-Memory Dispatcher
	bus := eventbus.NewInMemoryEventBus()
	_ = bus
	fmt.Println("✅ Decoupled EventBus initialized")

	// 4. Initialize Echo v5 Web Server
	e := echo.New()
	e.HideBanner = true

	// Standard Middlewares
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "1.0.0",
		})
	})

	// 5. Wire Hexagonal Modules
	apiV1 := e.Group("/api/v1")

	// System Module (21_ADM / Multi-Company)
	if dbPool != nil {
		companyRepo := sysRepo.NewPostgresCompanyRepo(dbPool)
		companyUC := sysUseCase.NewCompanyUseCase(companyRepo)
		companyHandler := sysDelivery.NewCompanyHandler(companyUC)
		companyHandler.RegisterRoutes(apiV1)
		fmt.Println("✅ Module [System / ADM] routes wired to /api/v1/companies")
	}

	// 6. Graceful Server Start
	go func() {
		addr := ":" + port
		fmt.Printf("🌐 Server listening on http://localhost%s\n", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server startup error: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("\n⏳ Shutting down ERP backend gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Forced shutdown: %v\n", err)
	}
	fmt.Println("👋 ERP backend stopped cleanly.")
}
