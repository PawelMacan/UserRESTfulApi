package errors

import (
	"errors"
	"fmt"
)

// Error types
var (
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInternalServer    = errors.New("internal server error")
	ErrDatabaseOperation = errors.New("database operation failed")
)

// ErrorType represents the type of error
type ErrorType string

const (
	NotFound          ErrorType = "NOT_FOUND"
	InvalidInput      ErrorType = "INVALID_INPUT"
	DuplicateEmail    ErrorType = "DUPLICATE_EMAIL"
	InvalidEmail      ErrorType = "INVALID_EMAIL"
	InvalidPassword   ErrorType = "INVALID_PASSWORD"
	Unauthorized      ErrorType = "UNAUTHORIZED"
	InternalServer    ErrorType = "INTERNAL_SERVER"
	DatabaseOperation ErrorType = "DATABASE_OPERATION"
)

// AppError represents a custom application error
type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Err     error    `json:"-"` // Original error (if any)
}

// Error returns the error message
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap returns the original error
func (e *AppError) Unwrap() error {
	return e.Err
}

// AppErrorf creates a new AppError with formatted message
func AppErrorf(errType ErrorType, format string, args ...interface{}) *AppError {
	return &AppError{
		Type:    errType,
		Message: fmt.Sprintf(format, args...),
		Err:     nil,
	}
}

// Error constructors for common cases
func NotFoundError(resource string, id interface{}) *AppError {
	return &AppError{
		Type:    NotFound,
		Message: fmt.Sprintf("%s with ID %v not found", resource, id),
		Err:     ErrNotFound,
	}
}

func InvalidInputError(field, reason string) *AppError {
	return &AppError{
		Type:    InvalidInput,
		Message: fmt.Sprintf("invalid %s: %s", field, reason),
		Err:     ErrInvalidInput,
	}
}

func DuplicateEmailError(email string) *AppError {
	return &AppError{
		Type:    DuplicateEmail,
		Message: fmt.Sprintf("email %s already exists", email),
		Err:     ErrDuplicateEmail,
	}
}

func InvalidEmailError(email string) *AppError {
	return &AppError{
		Type:    InvalidEmail,
		Message: fmt.Sprintf("invalid email format: %s", email),
		Err:     ErrInvalidEmail,
	}
}

func InvalidPasswordError(reason string) *AppError {
	return &AppError{
		Type:    InvalidPassword,
		Message: fmt.Sprintf("invalid password: %s", reason),
		Err:     ErrInvalidPassword,
	}
}

func UnauthorizedError(reason string) *AppError {
	return &AppError{
		Type:    Unauthorized,
		Message: fmt.Sprintf("unauthorized: %s", reason),
		Err:     ErrUnauthorized,
	}
}

func InternalServerError(err error) *AppError {
	return &AppError{
		Type:    InternalServer,
		Message: "internal server error",
		Err:     err,
	}
}

func DatabaseError(operation string, err error) *AppError {
	return &AppError{
		Type:    DatabaseOperation,
		Message: fmt.Sprintf("database %s operation failed", operation),
		Err:     err,
	}
}
