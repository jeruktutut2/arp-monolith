package apperrors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource conflict or already exists")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("access forbidden")
	ErrValidation     = errors.New("validation failed")
	ErrInternalServer = errors.New("internal server error")
)

// AppError represents a typed application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
