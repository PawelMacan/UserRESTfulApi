package service

import (
	"UserRESTfulApi/internal/domain"
	"UserRESTfulApi/internal/errors"
	"testing"
)

// Mock repository for testing
type mockUserRepository struct {
	users map[uint]*domain.User
	// Track function calls for testing
	getByIDCalled    bool
	getByEmailCalled bool
	createCalled     bool
	updateCalled     bool
	deleteCalled     bool
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[uint]*domain.User),
	}
}

func (m *mockUserRepository) GetByID(id uint) (*domain.User, error) {
	m.getByIDCalled = true
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.NotFoundError("user", id)
}

func (m *mockUserRepository) GetByEmail(email string) (*domain.User, error) {
	m.getByEmailCalled = true
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) Create(user *domain.User) error {
	m.createCalled = true
	if user.ID == 0 {
		user.ID = uint(len(m.users) + 1)
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Update(user *domain.User) error {
	m.updateCalled = true
	if _, exists := m.users[user.ID]; !exists {
		return errors.NotFoundError("user", user.ID)
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Delete(id uint) error {
	m.deleteCalled = true
	if _, exists := m.users[id]; !exists {
		return errors.NotFoundError("user", id)
	}
	delete(m.users, id)
	return nil
}

func (m *mockUserRepository) List(page, limit int) ([]*domain.User, error) {
	var users []*domain.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &domain.User{
				Email:    "test@example.com",
				Password: "Test123!@#",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: &domain.User{
				Email:    "invalid-email",
				Password: "Test123!@#",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			user: &domain.User{
				Email:    "test@example.com",
				Password: "weak",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			service := NewUserService(repo)

			err := service.CreateUser(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify repository was called
				if !repo.createCalled {
					t.Error("Create() was not called on repository")
				}
				// Verify password was hashed
				if tt.user.Password == "Test123!@#" {
					t.Error("Password was not hashed")
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	// Create initial user
	initialUser := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "Test123!@#",
	}

	tests := []struct {
		name    string
		user    *domain.User
		setup   func(*mockUserRepository)
		wantErr bool
	}{
		{
			name: "valid update",
			user: &domain.User{
				ID:       1,
				Email:    "new@example.com",
				Password: "NewTest123!@#",
			},
			setup: func(repo *mockUserRepository) {
				repo.Create(initialUser)
			},
			wantErr: false,
		},
		{
			name: "non-existent user",
			user: &domain.User{
				ID:       999,
				Email:    "new@example.com",
				Password: "NewTest123!@#",
			},
			setup:   func(repo *mockUserRepository) {},
			wantErr: true,
		},
		{
			name: "invalid email",
			user: &domain.User{
				ID:       1,
				Email:    "invalid-email",
				Password: "NewTest123!@#",
			},
			setup: func(repo *mockUserRepository) {
				repo.Create(initialUser)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			tt.setup(repo)
			service := NewUserService(repo)

			err := service.UpdateUser(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify repository was called
				if !repo.updateCalled {
					t.Error("Update() was not called on repository")
				}
				// Verify user was updated
				updated, _ := repo.GetByID(tt.user.ID)
				if updated.Email != tt.user.Email {
					t.Errorf("User email not updated, got = %v, want %v", updated.Email, tt.user.Email)
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	existingUser := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "Test123!@#",
	}

	tests := []struct {
		name    string
		userID  uint
		setup   func(*mockUserRepository)
		want    *domain.User
		wantErr bool
	}{
		{
			name:   "existing user",
			userID: 1,
			setup: func(repo *mockUserRepository) {
				repo.Create(existingUser)
			},
			want:    existingUser,
			wantErr: false,
		},
		{
			name:    "non-existent user",
			userID:  999,
			setup:   func(repo *mockUserRepository) {},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			tt.setup(repo)
			service := NewUserService(repo)

			got, err := service.GetUser(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.ID != tt.want.ID || got.Email != tt.want.Email {
					t.Errorf("GetUser() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
