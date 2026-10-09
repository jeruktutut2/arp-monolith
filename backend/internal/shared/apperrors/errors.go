package apperrors

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource conflict or already exists")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("access forbidden")
	ErrValidation     = errors.New("validation failed")
	ErrInternalServer = errors.New("internal server error")
)
