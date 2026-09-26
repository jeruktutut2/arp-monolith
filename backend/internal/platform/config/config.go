package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config aggregates all application configuration values
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port    string `env:"SERVER_PORT" envDefault:"8080"`
	Env     string `env:"SERVER_ENV" envDefault:"development"`
	Version string `env:"APP_VERSION" envDefault:"1.0.0"`
}

// IsProduction returns true if server runs in production environment
func (s ServerConfig) IsProduction() bool {
	return s.Env == "production"
}

// DatabaseConfig holds PostgreSQL & PgBouncer connection configuration
type DatabaseConfig struct {
	URL             string        `env:"DATABASE_URL" envDefault:"postgres://erp_user:erp_secret@localhost:6432/erp_db?sslmode=disable"`
	DirectURL       string        `env:"DIRECT_DATABASE_URL" envDefault:"postgres://erp_user:erp_secret@localhost:5432/erp_db?sslmode=disable"`
	MaxConns        int32         `env:"DB_MAX_CONNS" envDefault:"50"`
	MinConns        int32         `env:"DB_MIN_CONNS" envDefault:"10"`
	MaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME" envDefault:"1h"`
	MaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME" envDefault:"30m"`
}

// Load loads configuration from environment variables using caarlos0/env and godotenv
func Load() (*Config, error) {
	// Attempt to load .env if present (fails silently if file does not exist)
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
