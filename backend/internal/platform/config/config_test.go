package config_test

import (
	"testing"
	"time"

	"erp_monolith/backend/internal/platform/config"
)

func TestConfigLoad(t *testing.T) {
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("SERVER_ENV", "production")
	t.Setenv("DB_MAX_CONNS", "25")
	t.Setenv("DB_MAX_CONN_LIFETIME", "45m")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error loading config, got %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Server.Port)
	}
	if !cfg.Server.IsProduction() {
		t.Errorf("expected production mode")
	}
	if cfg.Database.MaxConns != 25 {
		t.Errorf("expected 25 max conns, got %d", cfg.Database.MaxConns)
	}
	if cfg.Database.MaxConnLifetime != 45*time.Minute {
		t.Errorf("expected 45m lifetime, got %v", cfg.Database.MaxConnLifetime)
	}
}
