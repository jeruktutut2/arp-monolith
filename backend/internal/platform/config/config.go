package config

import (
	"time"
)

// Config aggregates all application configuration values
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port    string
	Env     string // "development", "staging", "production"
	Version string
}

// IsProduction returns true if server runs in production environment
func (s ServerConfig) IsProduction() bool {
	return s.Env == "production"
}

// DatabaseConfig holds PostgreSQL & PgBouncer connection configuration
type DatabaseConfig struct {
	URL             string
	DirectURL       string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// CORSConfig holds CORS origin settings
type CORSConfig struct {
	AllowedOrigins []string
}

// Load loads configuration from environment variables (and optional .env file)
func Load() (*Config, error) {
	// Attempt loading from local .env files if present (does not overwrite existing environment variables)
	_ = LoadEnvFile(".env", "../.env")

	cfg := &Config{
		Server: ServerConfig{
			Port:    GetString("SERVER_PORT", "8080"),
			Env:     GetString("SERVER_ENV", "development"),
			Version: GetString("APP_VERSION", "1.0.0"),
		},
		Database: DatabaseConfig{
			URL:             GetString("DATABASE_URL", "postgres://erp_user:erp_secret@localhost:6432/erp_db?sslmode=disable"),
			DirectURL:       GetString("DIRECT_DATABASE_URL", "postgres://erp_user:erp_secret@localhost:5432/erp_db?sslmode=disable"),
			MaxConns:        GetInt32("DB_MAX_CONNS", 50),
			MinConns:        GetInt32("DB_MIN_CONNS", 10),
			MaxConnLifetime: GetDuration("DB_MAX_CONN_LIFETIME", time.Hour),
			MaxConnIdleTime: GetDuration("DB_MAX_CONN_IDLE_TIME", 30*time.Minute),
		},
		JWT: JWTConfig{
			Secret:          GetString("JWT_SECRET", "super-secret-erp-jwt-key-change-in-production"),
			ExpirationHours: GetInt("JWT_EXPIRATION_HOURS", 24),
		},
		CORS: CORSConfig{
			AllowedOrigins: GetStringSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173", "http://localhost:3000"}),
		},
	}

	return cfg, nil
}
