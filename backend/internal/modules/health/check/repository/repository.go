package repository

import (
	"context"

	"erp_monolith/backend/internal/modules/health/check/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresHealthChecker struct {
	pool *pgxpool.Pool
}

// NewPostgresHealthChecker creates a driven adapter to check PostgreSQL connectivity
func NewPostgresHealthChecker(pool *pgxpool.Pool) domain.HealthChecker {
	return &postgresHealthChecker{pool: pool}
}

func (r *postgresHealthChecker) Ping(ctx context.Context) error {
	if r.pool == nil {
		return nil
	}
	return r.pool.Ping(ctx)
}
