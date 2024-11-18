package integration

import (
	"UserRESTfulApi/internal"
	"UserRESTfulApi/internal/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	router *gin.Engine
	db     *gorm.DB
)

func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Wait for database to be ready
	fmt.Println("Waiting for database to be ready...")
	time.Sleep(2 * time.Second)

	// Setup test database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefault("DB_PORT", "5432"),
		getEnvOrDefault("DB_USER", "postgres"),
		getEnvOrDefault("DB_PASSWORD", "postgres"),
		getEnvOrDefault("DB_NAME", "test_db"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		fmt.Printf("Error migrating database: %v\n", err)
		os.Exit(1)
	}

	// Setup router
	router = internal.SetupRouter(db)

	// Run tests
	code := m.Run()

	// Clean up
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Error getting database instance: %v\n", err)
		os.Exit(1)
	}
	sqlDB.Close()

	os.Exit(code)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func clearDatabase() {
	db.Exec("DELETE FROM users")
}

func cleanupDatabase(t *testing.T) {
	err := db.Exec("TRUNCATE users CASCADE").Error
	if err != nil {
		t.Fatalf("Failed to cleanup database: %v", err)
	}
}

func setupTest(t *testing.T) {
	cleanupDatabase(t)
}

func createTestUser(t *testing.T) *domain.User {
	setupTest(t)

	user := &domain.User{
		Email:    "test@example.com",
		Password: "Test@123",
		Name:     "Test User",
	}

	w := httptest.NewRecorder()
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/users", bytes.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var createdUser domain.User
	err := json.NewDecoder(w.Body).Decode(&createdUser)
	assert.NoError(t, err)

	return &createdUser
}

func TestCreateUser(t *testing.T) {
	setupTest(t)

	user := &domain.User{
		Email:    "test@example.com",
		Password: "Test@123",
		Name:     "Test User",
	}

	w := httptest.NewRecorder()
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/users", bytes.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var createdUser domain.User
	err := json.NewDecoder(w.Body).Decode(&createdUser)
	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.NotEmpty(t, createdUser.CreatedAt)
	assert.NotEmpty(t, createdUser.UpdatedAt)
}

func TestGetUser(t *testing.T) {
	setupTest(t)
	user := createTestUser(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/users/%d", user.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var fetchedUser domain.User
	err := json.NewDecoder(w.Body).Decode(&fetchedUser)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, fetchedUser.ID)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.Name, fetchedUser.Name)
}

func TestUpdateUser(t *testing.T) {
	setupTest(t)
	user := createTestUser(t)

	updatedUser := &domain.User{
		Email: "updated@example.com",
		Name:  "Updated User",
	}

	w := httptest.NewRecorder()
	body, _ := json.Marshal(updatedUser)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/users/%d", user.ID), bytes.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var responseUser domain.User
	err := json.NewDecoder(w.Body).Decode(&responseUser)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, responseUser.ID)
	assert.Equal(t, updatedUser.Email, responseUser.Email)
	assert.Equal(t, updatedUser.Name, responseUser.Name)
}

func TestDeleteUser(t *testing.T) {
	setupTest(t)
	user := createTestUser(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/users/%d", user.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Verify user is deleted
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/users/%d", user.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestListUsers(t *testing.T) {
	setupTest(t)

	// Create multiple users
	for i := 0; i < 3; i++ {
		user := &domain.User{
			Email:    fmt.Sprintf("test%d@example.com", i),
			Password: "Test@123",
			Name:     fmt.Sprintf("Test User %d", i),
		}

		w := httptest.NewRecorder()
		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/api/users", bytes.NewReader(body))
		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
	}

	// Test listing users
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var users []*domain.User
	err := json.NewDecoder(w.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Len(t, users, 3)
}

func TestCreateUserValidation(t *testing.T) {
	setupTest(t)

	testCases := []struct {
		name     string
		user     domain.User
		wantCode int
	}{
		{
			name: "Invalid Email",
			user: domain.User{
				Email:    "invalid-email",
				Password: "Test@123",
				Name:     "Test User",
			},
			wantCode: 400,
		},
		{
			name: "Weak Password",
			user: domain.User{
				Email:    "test@example.com",
				Password: "weak",
				Name:     "Test User",
			},
			wantCode: 400,
		},
		{
			name: "Empty Name",
			user: domain.User{
				Email:    "test@example.com",
				Password: "Test@123",
				Name:     "",
			},
			wantCode: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.user)
			req := httptest.NewRequest("POST", "/api/users", bytes.NewReader(body))
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantCode, w.Code)

			var response map[string]string
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Contains(t, response, "error")
		})
	}
}

func TestDuplicateEmail(t *testing.T) {
	setupTest(t)

	// Create first user
	user := &domain.User{
		Email:    "test@example.com",
		Password: "Test@123",
		Name:     "Test User",
	}

	w := httptest.NewRecorder()
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/users", bytes.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	// Try to create user with same email
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/users", bytes.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "already registered")
}
