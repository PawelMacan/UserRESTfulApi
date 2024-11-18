package password

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong  = errors.New("password must not exceed 72 characters")
	ErrMissingLower     = errors.New("password must contain at least one lowercase letter")
	ErrMissingUpper     = errors.New("password must contain at least one uppercase letter")
	ErrMissingNumber    = errors.New("password must contain at least one number")
	ErrMissingSpecial   = errors.New("password must contain at least one special character")
)

// Validate checks if the password meets the required criteria
func Validate(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	if len(password) > 72 {
		return ErrPasswordTooLong
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrMissingUpper
	}
	if !hasLower {
		return ErrMissingLower
	}
	if !hasNumber {
		return ErrMissingNumber
	}
	if !hasSpecial {
		return ErrMissingSpecial
	}

	return nil
}

// Hash creates a bcrypt hash from a password string
func Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Verify checks if the provided password matches the hashed password
func Verify(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
