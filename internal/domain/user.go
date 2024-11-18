package domain

import "time"

// User represents the user entity
type User struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`  // "-" to exclude from JSON responses
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the user storage interface
type UserRepository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List(page, limit int) ([]*User, error)
}

// UserService defines the user business logic interface
type UserService interface {
	CreateUser(user *User) error
	GetUser(id uint) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
	ListUsers(page, limit int) ([]*User, error)
	VerifyPassword(email, password string) (*User, error)
}
