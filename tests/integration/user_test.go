package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"UserRESTfulApi/internal/domain"
)

func TestUserAPI(t *testing.T) {
	setupTest(t)

	t.Run("Create User Flow", func(t *testing.T) {
		testUser := domain.User{
			Email:    "create@example.com",
			Password: "Test123!@#",
			Name:     "Test User",
		}

		// Create user
		rr := makeRequest(t, http.MethodPost, "/api/users", testUser)
		
		if rr.Code != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
		}

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Verify user was created
		userID := uint(response["id"].(float64))
		if userID == 0 {
			t.Error("Expected user ID to be returned")
		}

		// Get created user
		rr = makeRequest(t, http.MethodGet, fmt.Sprintf("/api/users/%d", userID), nil)
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
		}

		var user domain.User
		err = json.Unmarshal(rr.Body.Bytes(), &user)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if user.Email != testUser.Email {
			t.Errorf("handler returned wrong email: got %v want %v", user.Email, testUser.Email)
		}
	})

	t.Run("Update User Flow", func(t *testing.T) {
		// Create user first
		testUser := domain.User{
			Email:    "update@example.com",
			Password: "Test123!@#",
			Name:     "Test User",
		}

		rr := makeRequest(t, http.MethodPost, "/api/users", testUser)
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		userID := uint(response["id"].(float64))

		// Update user
		updatedUser := domain.User{
			ID:       userID,
			Email:    "updated@example.com",
			Password: "UpdatedTest123!@#",
			Name:     "Updated User",
		}

		rr = makeRequest(t, http.MethodPut, fmt.Sprintf("/api/users/%d", userID), updatedUser)
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
		}

		// Verify update
		rr = makeRequest(t, http.MethodGet, fmt.Sprintf("/api/users/%d", userID), nil)
		var user domain.User
		err = json.Unmarshal(rr.Body.Bytes(), &user)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if user.Email != updatedUser.Email {
			t.Errorf("handler returned wrong email: got %v want %v", user.Email, updatedUser.Email)
		}
	})

	t.Run("Delete User Flow", func(t *testing.T) {
		// Create user first
		testUser := domain.User{
			Email:    "delete@example.com",
			Password: "Test123!@#",
			Name:     "Test User",
		}

		rr := makeRequest(t, http.MethodPost, "/api/users", testUser)
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		userID := uint(response["id"].(float64))

		// Delete user
		rr = makeRequest(t, http.MethodDelete, fmt.Sprintf("/api/users/%d", userID), nil)
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
		}

		// Verify deletion
		rr = makeRequest(t, http.MethodGet, fmt.Sprintf("/api/users/%d", userID), nil)
		if rr.Code != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("List Users Flow", func(t *testing.T) {
		// Create multiple users
		for i := 0; i < 3; i++ {
			user := domain.User{
				Email:    fmt.Sprintf("list%d@example.com", i),
				Password: "Test123!@#",
				Name:     fmt.Sprintf("Test User %d", i),
			}
			makeRequest(t, http.MethodPost, "/api/users", user)
		}

		// Get user list
		rr := makeRequest(t, http.MethodGet, "/api/users?page=1&limit=10", nil)
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
		}

		var users []*domain.User
		err := json.Unmarshal(rr.Body.Bytes(), &users)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(users) < 3 {
			t.Errorf("handler returned wrong number of users: got %v want at least 3", len(users))
		}
	})

	t.Run("Invalid Input Tests", func(t *testing.T) {
		// Test invalid email
		invalidUser := domain.User{
			Email:    "invalid-email",
			Password: "Test123!@#",
			Name:     "Test User",
		}
		rr := makeRequest(t, http.MethodPost, "/api/users", invalidUser)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("handler should return bad request for invalid email: got %v want %v", rr.Code, http.StatusBadRequest)
		}

		// Test invalid password
		invalidUser = domain.User{
			Email:    "test@example.com",
			Password: "weak",
			Name:     "Test User",
		}
		rr = makeRequest(t, http.MethodPost, "/api/users", invalidUser)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("handler should return bad request for invalid password: got %v want %v", rr.Code, http.StatusBadRequest)
		}

		// Test empty name
		invalidUser = domain.User{
			Email:    "test@example.com",
			Password: "Test123!@#",
			Name:     "",
		}
		rr = makeRequest(t, http.MethodPost, "/api/users", invalidUser)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("handler should return bad request for empty name: got %v want %v", rr.Code, http.StatusBadRequest)
		}

		// Test duplicate email
		makeRequest(t, http.MethodPost, "/api/users", domain.User{
			Email:    "duplicate@example.com",
			Password: "Test123!@#",
			Name:     "Test User",
		})
		rr = makeRequest(t, http.MethodPost, "/api/users", domain.User{
			Email:    "duplicate@example.com",
			Password: "Test123!@#",
			Name:     "Test User",
		})
		if rr.Code != http.StatusConflict {
			t.Errorf("handler should return conflict for duplicate email: got %v want %v", rr.Code, http.StatusConflict)
		}
	})
}
