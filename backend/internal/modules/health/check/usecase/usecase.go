package usecase

import (
	"context"
	"time"

	"erp_monolith/backend/internal/modules/health/check/domain"
)

type healthUseCase struct {
	version string
	checker domain.HealthChecker
}

// NewHealthUseCase creates a new HealthUseCase implementation
func NewHealthUseCase(version string, checker domain.HealthChecker) domain.HealthUseCase {
	return &healthUseCase{
		version: version,
		checker: checker,
	}
}

func (u *healthUseCase) Check(ctx context.Context) domain.Health {
	status := "healthy"
	if u.checker != nil {
		if err := u.checker.Ping(ctx); err != nil {
			status = "degraded"
		}
	}

	return domain.Health{
		Status:    status,
		Timestamp: time.Now().UTC(),
		Version:   u.version,
	}
}
