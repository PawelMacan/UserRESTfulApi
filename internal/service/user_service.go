package service

import (
	"fmt"
	"net/mail"

	"github.com/user-api/internal/domain"
	"github.com/user-api/pkg/password"
	"github.com/user-api/pkg/validation"
)

type userService struct {
	repo domain.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(user *domain.User) error {
	// Validate user data
	if err := s.validateUser(user); err != nil {
		return err
	}

	// Validate password
	if err := validation.ValidatePassword(user.Password); err != nil {
		return fmt.Errorf("invalid password: %v", err)
	}

	// Hash the password before saving
	hashedPassword, err := password.Hash(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	user.Password = hashedPassword

	return s.repo.Create(user)
}

func (s *userService) GetUser(id uint) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *userService) UpdateUser(user *domain.User) error {
	// Check if user exists
	existing, err := s.repo.GetByID(user.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("user not found")
	}

	// If password is provided, validate and hash it
	if user.Password != "" {
		if err := validation.ValidatePassword(user.Password); err != nil {
			return fmt.Errorf("invalid password: %v", err)
		}
		hashedPassword, err := password.Hash(user.Password)
		if err != nil {
			return fmt.Errorf("error hashing password: %v", err)
		}
		user.Password = hashedPassword
	} else {
		// Keep existing password if not provided
		user.Password = existing.Password
	}

	// Validate user data
	if err := s.validateUser(user); err != nil {
		return err
	}

	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	// Check if user exists
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("user not found")
	}

	return s.repo.Delete(id)
}

func (s *userService) ListUsers(page, limit int) ([]*domain.User, error) {
	return s.repo.List(page, limit)
}

func (s *userService) validateUser(user *domain.User) error {
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if user.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Validate email format
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// VerifyPassword checks if the provided password matches the user's password
func (s *userService) VerifyPassword(email, plainPassword string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if !password.Verify(plainPassword, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}
