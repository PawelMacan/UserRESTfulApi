package internal

import (
	"UserRESTfulApi/internal/handlers"
	"UserRESTfulApi/internal/repository/postgres"
	"UserRESTfulApi/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Router handles all routing for the application
type Router struct {
	engine *gin.Engine
}

// NewRouter creates a new router instance
func NewRouter(db *gorm.DB) *Router {
	engine := SetupRouter(db)
	return &Router{engine: engine}
}

// SetupRouter sets up the router with all routes
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Create dependencies
	userRepo := postgres.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// API routes
	api := router.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("", userHandler.ListUsers)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}

// Run starts the HTTP server
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
