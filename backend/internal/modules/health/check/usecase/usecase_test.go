package usecase_test

import (
	"context"
	"errors"
	"testing"

	"erp_monolith/backend/internal/modules/health/check/usecase"
)

type mockChecker struct {
	err error
}

func (m *mockChecker) Ping(ctx context.Context) error {
	return m.err
}

func TestHealthUseCase_Check(t *testing.T) {
	t.Run("healthy when checker passes", func(t *testing.T) {
		uc := usecase.NewHealthUseCase("1.0.0", &mockChecker{err: nil})
		h := uc.Check(context.Background())

		if h.Status != "healthy" {
			t.Errorf("expected status healthy, got %s", h.Status)
		}
		if h.Version != "1.0.0" {
			t.Errorf("expected version 1.0.0, got %s", h.Version)
		}
	})

	t.Run("degraded when checker fails", func(t *testing.T) {
		uc := usecase.NewHealthUseCase("1.0.0", &mockChecker{err: errors.New("db error")})
		h := uc.Check(context.Background())

		if h.Status != "degraded" {
			t.Errorf("expected status degraded, got %s", h.Status)
		}
	})
}
