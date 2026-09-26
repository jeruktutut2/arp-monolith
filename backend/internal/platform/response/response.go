package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// StandardResponse wraps all API responses
type StandardResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Pagination `json:"meta,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// Pagination metadata
type Pagination struct {
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
}

// Success returns a 200/201 HTTP JSON success response
func Success(c *echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, StandardResponse{
		Success: true,
		Code:    "OK",
		Message: message,
		Data:    data,
		TraceID: c.Response().Header().Get(echo.HeaderXRequestID),
	})
}

// Paginated returns a 200 HTTP JSON success response with pagination
func Paginated(c *echo.Context, message string, data interface{}, meta Pagination) error {
	return c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Code:    "OK",
		Message: message,
		Data:    data,
		Meta:    &meta,
		TraceID: c.Response().Header().Get(echo.HeaderXRequestID),
	})
}

// Error returns an error response
func Error(c *echo.Context, status int, code, message string) error {
	return c.JSON(status, StandardResponse{
		Success: false,
		Code:    code,
		Message: message,
		TraceID: c.Response().Header().Get(echo.HeaderXRequestID),
	})
}
