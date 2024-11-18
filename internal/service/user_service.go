package service

import (
	"UserRESTfulApi/internal/domain"
	"UserRESTfulApi/internal/errors"
	"net/mail"
	"strings"
	"unicode"
)

type userService struct {
	repo domain.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

// Create creates a new user
func (s *userService) Create(user *domain.User) error {
	if err := s.validateEmail(user.Email); err != nil {
		return err
	}

	if err := s.validatePassword(user.Password); err != nil {
		return err
	}

	if err := s.validateName(user.Name); err != nil {
		return err
	}

	existingUser, err := s.repo.GetByEmail(user.Email)
	if err != nil {
		return errors.InternalServerError(err)
	}
	if existingUser != nil {
		return errors.DuplicateEmailError(user.Email)
	}

	// TODO: Hash password before saving
	return s.repo.Create(user)
}

// Get retrieves a user by ID
func (s *userService) Get(id uint) (*domain.User, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NotFoundError("user", id)
	}
	return user, nil
}

// Update updates a user
func (s *userService) Update(user *domain.User) error {
	if err := s.validateEmail(user.Email); err != nil {
		return err
	}

	if user.Password != "" {
		if err := s.validatePassword(user.Password); err != nil {
			return err
		}
	}

	if err := s.validateName(user.Name); err != nil {
		return err
	}

	existingUser, err := s.repo.Get(user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.NotFoundError("user", user.ID)
	}

	// Check if email is being changed and if it's already taken
	if existingUser.Email != user.Email {
		emailUser, err := s.repo.GetByEmail(user.Email)
		if err != nil {
			return errors.InternalServerError(err)
		}
		if emailUser != nil {
			return errors.DuplicateEmailError(user.Email)
		}
	}

	// TODO: Hash password before saving if it's being updated
	return s.repo.Update(user)
}

// Delete deletes a user
func (s *userService) Delete(id uint) error {
	user, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.NotFoundError("user", id)
	}
	return s.repo.Delete(id)
}

// List lists users with pagination
func (s *userService) List(page, limit int) ([]*domain.User, error) {
	return s.repo.List(page, limit)
}

// GetByEmail retrieves a user by email
func (s *userService) GetByEmail(email string) (*domain.User, error) {
	if err := s.validateEmail(email); err != nil {
		return nil, err
	}
	return s.repo.GetByEmail(email)
}

// VerifyPassword verifies a user's password
func (s *userService) VerifyPassword(email, plainPassword string) (*domain.User, error) {
	user, err := s.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NotFoundError("user", email)
	}

	// TODO: Verify password
	return user, nil
}

// validateEmail validates email format
func (s *userService) validateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.InvalidEmailError("email cannot be empty")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.InvalidEmailError(email)
	}
	return nil
}

// validatePassword validates password strength
func (s *userService) validatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return errors.InvalidPasswordError("password cannot be empty")
	}

	if len(password) < 8 {
		return errors.InvalidPasswordError("password must be at least 8 characters long")
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

	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasNumber {
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return errors.InvalidPasswordError("password must contain at least one " + strings.Join(missing, ", "))
	}

	return nil
}

// validateName validates user name
func (s *userService) validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.InvalidInputError("name", "cannot be empty")
	}
	return nil
}
