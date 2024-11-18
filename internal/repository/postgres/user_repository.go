package postgres

import (
	"UserRESTfulApi/internal/domain"
	"UserRESTfulApi/internal/errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result := r.db.Create(user)
	if result.Error != nil {
		log.Printf("Failed to create user with email %s: %v", user.Email, result.Error)
		return errors.DatabaseError("create", result.Error)
	}

	return nil
}

// Get retrieves a user by ID
func (r *userRepository) Get(id uint) (*domain.User, error) {
	var user domain.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Printf("Failed to get user with id %d: %v", id, result.Error)
		return nil, errors.DatabaseError("get", result.Error)
	}

	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(user *domain.User) error {
	user.UpdatedAt = time.Now()

	result := r.db.Save(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Failed to update user with id %d: %v", user.ID, result.Error)
			return errors.NotFoundError("user", user.ID)
		}
		log.Printf("Failed to update user with id %d: %v", user.ID, result.Error)
		return errors.DatabaseError("update", result.Error)
	}

	return nil
}

// Delete deletes a user
func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.User{}, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Failed to delete user with id %d: %v", id, result.Error)
			return errors.NotFoundError("user", id)
		}
		log.Printf("Failed to delete user with id %d: %v", id, result.Error)
		return errors.DatabaseError("delete", result.Error)
	}

	return nil
}

// List retrieves users with pagination
func (r *userRepository) List(page, limit int) ([]*domain.User, error) {
	var users []*domain.User
	offset := (page - 1) * limit

	result := r.db.Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		log.Printf("Failed to list users: %v", result.Error)
		return nil, errors.DatabaseError("list", result.Error)
	}

	return users, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Printf("Failed to get user with email %s: %v", email, result.Error)
		return nil, errors.DatabaseError("get by email", result.Error)
	}

	return &user, nil
}
