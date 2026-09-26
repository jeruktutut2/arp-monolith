package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"erp_monolith/backend/internal/modules/health"
	"erp_monolith/backend/internal/modules/system"
	"erp_monolith/backend/internal/platform/database"
	"erp_monolith/backend/internal/platform/eventbus"

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

	// Standard Middlewares
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 5. Wire Health Check Module
	health.RegisterModule(e, dbPool)

	// 6. Wire Business Modules (admin: 21_ADM, user: 19_USR, audit: 25_AUD)
	system.RegisterModule(e, dbPool)

	// 7. Graceful Server Start via Echo v5 StartConfig
	fmt.Printf("🌐 Server listening on http://localhost:%s\n", port)

	serverCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":" + port,
		GracefulTimeout: 10 * time.Second,
		HideBanner:      true,
		HidePort:        true,
	}

	if err := sc.Start(serverCtx, e); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("❌ Server error: %v\n", err)
	}

	fmt.Println("\n👋 ERP backend stopped cleanly.")
}
