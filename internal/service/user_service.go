package service

import (
	"net/mail"

	"github.com/user-api/internal/domain"
	"github.com/user-api/internal/errors"
	"github.com/user-api/pkg/password"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(user *domain.User) error {
	// Validate email format
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.InvalidEmailError(user.Email)
	}

	// Check if email already exists
	existingUser, err := s.repo.GetByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.DuplicateEmailError(user.Email)
	}

	// Validate password
	if err := password.Validate(user.Password); err != nil {
		return errors.InvalidPasswordError(err.Error())
	}

	// Hash password
	hashedPassword, err := password.Hash(user.Password)
	if err != nil {
		return errors.InternalServerError(err)
	}
	user.Password = hashedPassword

	// Create user
	if err := s.repo.Create(user); err != nil {
		return errors.DatabaseError("create", err)
	}

	return nil
}

func (s *userService) GetUser(id uint) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.NotFoundError("user", id)
	}
	return user, nil
}

func (s *userService) UpdateUser(user *domain.User) error {
	// Check if user exists
	existingUser, err := s.repo.GetByID(user.ID)
	if err != nil {
		return errors.NotFoundError("user", user.ID)
	}

	// Validate email format if changed
	if user.Email != existingUser.Email {
		if _, err := mail.ParseAddress(user.Email); err != nil {
			return errors.InvalidEmailError(user.Email)
		}

		// Check if new email is already taken
		emailUser, err := s.repo.GetByEmail(user.Email)
		if err == nil && emailUser != nil && emailUser.ID != user.ID {
			return errors.DuplicateEmailError(user.Email)
		}
	}

	// Validate and hash password if changed
	if user.Password != "" {
		if err := password.Validate(user.Password); err != nil {
			return errors.InvalidPasswordError(err.Error())
		}

		hashedPassword, err := password.Hash(user.Password)
		if err != nil {
			return errors.InternalServerError(err)
		}
		user.Password = hashedPassword
	} else {
		user.Password = existingUser.Password
	}

	// Update user
	if err := s.repo.Update(user); err != nil {
		return errors.DatabaseError("update", err)
	}

	return nil
}

func (s *userService) DeleteUser(id uint) error {
	// Check if user exists
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.NotFoundError("user", id)
	}

	// Delete user
	if err := s.repo.Delete(id); err != nil {
		return errors.DatabaseError("delete", err)
	}

	return nil
}

func (s *userService) ListUsers(page, limit int) ([]*domain.User, error) {
	if page < 1 {
		return nil, errors.InvalidInputError("page", "must be greater than 0")
	}
	if limit < 1 {
		return nil, errors.InvalidInputError("limit", "must be greater than 0")
	}

	users, err := s.repo.List(page, limit)
	if err != nil {
		return nil, errors.DatabaseError("list", err)
	}

	return users, nil
}

func (s *userService) VerifyPassword(email, plainPassword string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.NotFoundError("user", email)
	}
	if user == nil {
		return nil, errors.NotFoundError("user", email)
	}

	if !password.Verify(plainPassword, user.Password) {
		return nil, errors.InvalidPasswordError("password does not match")
	}

	return user, nil
}
