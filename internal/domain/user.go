package domain

import "time"

// User represents the user entity
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password,omitempty" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserService defines the interface for user business logic
type UserService interface {
	Create(user *User) error
	Get(id uint) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List(page, limit int) ([]*User, error)
	GetByEmail(email string) (*User, error)
}

// UserRepository defines the interface for user data persistence
type UserRepository interface {
	Create(user *User) error
	Get(id uint) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List(page, limit int) ([]*User, error)
	GetByEmail(email string) (*User, error)
}
