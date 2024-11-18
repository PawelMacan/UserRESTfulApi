package validation

import (
	"fmt"
	"unicode"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72 // bcrypt's maximum input length
)

// PasswordError represents password validation errors
type PasswordError struct {
	TooShort      bool
	TooLong       bool
	NoUpper       bool
	NoLower       bool
	NoNumber      bool
	NoSpecial     bool
	ContainsSpace bool
}

func (e PasswordError) Error() string {
	if e.TooShort {
		return fmt.Sprintf("password must be at least %d characters long", MinPasswordLength)
	}
	if e.TooLong {
		return fmt.Sprintf("password must not exceed %d characters", MaxPasswordLength)
	}
	if e.NoUpper {
		return "password must contain at least one uppercase letter"
	}
	if e.NoLower {
		return "password must contain at least one lowercase letter"
	}
	if e.NoNumber {
		return "password must contain at least one number"
	}
	if e.NoSpecial {
		return "password must contain at least one special character"
	}
	if e.ContainsSpace {
		return "password must not contain spaces"
	}
	return ""
}

// ValidatePassword checks if a password meets the following criteria:
// - Between MinPasswordLength and MaxPasswordLength characters long
// - Contains at least one uppercase letter
// - Contains at least one lowercase letter
// - Contains at least one number
// - Contains at least one special character
// - Does not contain spaces
func ValidatePassword(password string) error {
	var err PasswordError

	if len(password) < MinPasswordLength {
		err.TooShort = true
		return err
	}

	if len(password) > MaxPasswordLength {
		err.TooLong = true
		return err
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
		case unicode.IsSpace(char):
			err.ContainsSpace = true
			return err
		}
	}

	if !hasUpper {
		err.NoUpper = true
		return err
	}
	if !hasLower {
		err.NoLower = true
		return err
	}
	if !hasNumber {
		err.NoNumber = true
		return err
	}
	if !hasSpecial {
		err.NoSpecial = true
		return err
	}

	return nil
}
