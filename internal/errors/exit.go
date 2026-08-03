package errors

import (
	"errors"
	"fmt"
)

const (
	ExitSuccess           = 0
	ExitGeneralError      = 1
	ExitValidationWarning = 2
	ExitHealthScoreBlock  = 3
	ExitConfiguration     = 4
)

type Category string

const (
	CategoryGeneral       Category = "general"
	CategoryConfiguration Category = "configuration"
	CategoryLifecycle     Category = "lifecycle"
	CategoryIO            Category = "io"
	CategoryCommand       Category = "command"
)

type AppError struct {
	Message        string
	Category       Category
	Code           int
	Err            error
	Severity       string
	Recommendation string
	DocRef         string
	AutoFixID      string
}

func NewAppError(code int, msg string, err error) *AppError {
	return &AppError{
		Code:     code,
		Message:  msg,
		Category: CategoryGeneral,
		Err:      err,
	}
}

func New(category Category, code int, msg string, err error) *AppError {
	return &AppError{
		Message:  msg,
		Category: category,
		Code:     code,
		Err:      err,
	}
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

func (e *AppError) ExitCode() int {
	if e == nil || e.Code == 0 {
		return ExitGeneralError
	}
	return e.Code
}

func (e *AppError) WithSeverity(sev string) *AppError {
	e.Severity = sev
	return e
}

func (e *AppError) WithCategory(cat string) *AppError {
	e.Category = Category(cat)
	return e
}

func (e *AppError) WithRecommendation(rec string) *AppError {
	e.Recommendation = rec
	return e
}

func (e *AppError) WithDocRef(ref string) *AppError {
	e.DocRef = ref
	return e
}

func (e *AppError) WithAutoFixID(id string) *AppError {
	e.AutoFixID = id
	return e
}

func ExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.ExitCode()
	}
	return ExitGeneralError
}
