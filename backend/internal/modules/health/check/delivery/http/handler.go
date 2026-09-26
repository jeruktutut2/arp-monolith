package http

import (
	"net/http"
	"time"

	"erp_monolith/backend/internal/modules/health/check/domain"

	"github.com/labstack/echo/v5"
)

type HealthHandler struct {
	useCase domain.HealthUseCase
}

func NewHealthHandler(useCase domain.HealthUseCase) *HealthHandler {
	return &HealthHandler{useCase: useCase}
}

func (h *HealthHandler) HealthCheck(c *echo.Context) error {
	hData := h.useCase.Check(c.Request().Context())
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    hData.Status,
		"timestamp": hData.Timestamp.Format(time.RFC3339),
		"version":   hData.Version,
	})
}
