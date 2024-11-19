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
	getCalled       bool
	getByEmailCalled bool
	createCalled     bool
	updateCalled     bool
	deleteCalled     bool
	listCalled       bool
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[uint]*domain.User),
	}
}

func (m *mockUserRepository) Get(id uint) (*domain.User, error) {
	m.getCalled = true
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
	m.listCalled = true
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func TestCreateUser(t *testing.T) {
	repo := newMockUserRepository()
	service := NewUserService(repo)

	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &domain.User{
				Email:    "test@example.com",
				Password: "Password123!",
				Name:    "Test User",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &domain.User{
				Email:    "test@example.com",
				Password: "Password123!",
				Name:    "Another User",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Create(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	repo := newMockUserRepository()
	service := NewUserService(repo)

	// Create initial user
	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "Password123!",
		Name:    "Test User",
	}
	repo.users[user.ID] = user

	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
	}{
		{
			name: "valid update",
			user: &domain.User{
				ID:       1,
				Email:    "updated@example.com",
				Password: "NewPassword123!",
				Name:    "Updated User",
			},
			wantErr: false,
		},
		{
			name: "non-existent user",
			user: &domain.User{
				ID:       999,
				Email:    "nonexistent@example.com",
				Password: "Password123!",
				Name:    "Non-existent User",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Update(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	repo := newMockUserRepository()
	service := NewUserService(repo)

	// Create test user
	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "Password123!",
		Name:    "Test User",
	}
	repo.users[user.ID] = user

	tests := []struct {
		name    string
		id      uint
		want    *domain.User
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      1,
			want:    user,
			wantErr: false,
		},
		{
			name:    "non-existent user",
			id:      999,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.Get(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.want.ID {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
