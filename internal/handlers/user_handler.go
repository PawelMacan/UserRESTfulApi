package handlers

import (
	"UserRESTfulApi/internal/domain"
	"UserRESTfulApi/internal/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service domain.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser handles user creation
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Create(&user)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		switch appErr.Type {
		case errors.InvalidEmail, errors.InvalidPassword, errors.InvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		case errors.DuplicateEmail:
			c.JSON(http.StatusConflict, gin.H{"error": appErr.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles user retrieval
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.Get(uint(id))
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		switch appErr.Type {
		case errors.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles user updates
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = uint(id)
	err = h.service.Update(&user)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		switch appErr.Type {
		case errors.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Error()})
		case errors.InvalidEmail, errors.InvalidPassword, errors.InvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		case errors.DuplicateEmail:
			c.JSON(http.StatusConflict, gin.H{"error": appErr.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles user deletion
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		switch appErr.Type {
		case errors.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": appErr.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ListUsers handles user listing with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, err := h.service.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, users)
}
