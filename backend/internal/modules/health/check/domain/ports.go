package domain

import "context"

// HealthUseCase defines the inbound port for health checking operations
type HealthUseCase interface {
	Check(ctx context.Context) Health
}
