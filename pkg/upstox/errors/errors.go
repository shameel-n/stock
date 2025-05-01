package errors

import "fmt"

// UpstoxError represents a custom error type for Upstox API
type UpstoxError struct {
	Code    int
	Message string
	Err     error
}

func (e *UpstoxError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("upstox error [%d]: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("upstox error [%d]: %s", e.Code, e.Message)
}

// New creates a new UpstoxError
func New(code int, message string, err error) *UpstoxError {
	return &UpstoxError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error codes
const (
	ErrInvalidRequest    = 400
	ErrUnauthorized      = 401
	ErrForbidden         = 403
	ErrNotFound          = 404
	ErrRateLimitExceeded = 429
	ErrInternalServer    = 500
)
