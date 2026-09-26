package domain

import "context"

// HealthChecker defines the outbound port for infrastructure pinging (e.g. database)
type HealthChecker interface {
	Ping(ctx context.Context) error
}
