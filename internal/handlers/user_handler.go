package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/user-api/internal/domain"
	"github.com/user-api/internal/errors"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// handleError converts application errors to appropriate HTTP responses
func handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch appErr.Type {
	case errors.NotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": appErr.Message})
	case errors.InvalidInput, errors.InvalidEmail, errors.InvalidPassword:
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Message})
	case errors.DuplicateEmail:
		c.JSON(http.StatusConflict, gin.H{"error": appErr.Message})
	case errors.Unauthorized:
		c.JSON(http.StatusUnauthorized, gin.H{"error": appErr.Message})
	case errors.DatabaseOperation:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": appErr.Message})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Message})
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handleError(c, errors.InvalidInputError("request body", err.Error()))
		return
	}

	if err := h.service.CreateUser(&user); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errors.InvalidInputError("id", "must be a positive integer"))
		return
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errors.InvalidInputError("id", "must be a positive integer"))
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handleError(c, errors.InvalidInputError("request body", err.Error()))
		return
	}

	user.ID = uint(id)
	if err := h.service.UpdateUser(&user); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errors.InvalidInputError("id", "must be a positive integer"))
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		handleError(c, errors.InvalidInputError("page", "must be a positive integer"))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		handleError(c, errors.InvalidInputError("limit", "must be a positive integer"))
		return
	}

	users, err := h.service.ListUsers(page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}
