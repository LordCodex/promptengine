package errors

import "fmt"

// POSIX Exit Codes mapped to application outcomes
const (
	ExitSuccess            = 0
	ExitGeneralError       = 1
	ExitValidationWarning  = 2
	ExitHealthScoreBlock   = 3
)

// AppError is the standard error type for PromptEngine CLI
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError builds a custom CLI error
func NewAppError(code int, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}
