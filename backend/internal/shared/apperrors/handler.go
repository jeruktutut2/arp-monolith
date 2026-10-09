package apperrors

import (
	"context"
	"errors"
	"net/http"

	"erp_monolith/backend/internal/platform/response"
)

func FromError(err error, traceID string) (int, interface{}) {
	if err == nil {
		return 0, nil
	}

	if errors.Is(err, context.Canceled) {
		return 0, nil
	}

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout, response.SetGatewayTimeout("request timeout", traceID)
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, response.SetNotFound(err.Error(), traceID)
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, response.SetUnauthorized(err.Error(), traceID)
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, response.SetForbidden(err.Error(), traceID)
	case errors.Is(err, ErrConflict):
		return http.StatusConflict, response.SetConflict(err.Error(), traceID)
	case errors.Is(err, ErrValidation), errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, response.SetBadRequest(err.Error(), traceID)
	default:
		// Technical / internal error: hide internal error details from client
		return http.StatusInternalServerError, response.SetInternalServerError(traceID)
	}
}
