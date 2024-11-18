package errors

import "fmt"

type ErrorType string

const (
	NotFound          ErrorType = "NOT_FOUND"
	InvalidInput      ErrorType = "INVALID_INPUT"
	DuplicateEmail    ErrorType = "DUPLICATE_EMAIL"
	InvalidEmail      ErrorType = "INVALID_EMAIL"
	InvalidPassword   ErrorType = "INVALID_PASSWORD"
	DatabaseOperation ErrorType = "DATABASE_OPERATION"
	InternalServer    ErrorType = "INTERNAL_SERVER"
)

type AppError struct {
	Type    ErrorType
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

// NotFoundError creates a new not found error
func NotFoundError(resource string, id interface{}) error {
	return &AppError{
		Type:    NotFound,
		Message: fmt.Sprintf("%s with ID %v not found", resource, id),
	}
}

// InvalidInputError creates a new invalid input error
func InvalidInputError(field, reason string) error {
	return &AppError{
		Type:    InvalidInput,
		Message: fmt.Sprintf("Invalid input for %s: %s", field, reason),
	}
}

// DuplicateEmailError creates a new duplicate email error
func DuplicateEmailError(email string) error {
	return &AppError{
		Type:    DuplicateEmail,
		Message: fmt.Sprintf("Email %s is already registered", email),
	}
}

// InvalidEmailError creates a new invalid email error
func InvalidEmailError(email string) error {
	return &AppError{
		Type:    InvalidEmail,
		Message: fmt.Sprintf("Invalid email format: %s", email),
	}
}

// InvalidPasswordError creates a new invalid password error
func InvalidPasswordError(reason string) error {
	return &AppError{
		Type:    InvalidPassword,
		Message: fmt.Sprintf("Invalid password: %s", reason),
	}
}

// DatabaseError creates a new database operation error
func DatabaseError(operation string, err error) error {
	return &AppError{
		Type:    DatabaseOperation,
		Message: fmt.Sprintf("Database %s error: %v", operation, err),
	}
}

// InternalServerError creates a new internal server error
func InternalServerError(err error) error {
	return &AppError{
		Type:    InternalServer,
		Message: fmt.Sprintf("Internal server error: %v", err),
	}
}
