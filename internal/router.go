package internal

import (
	"github.com/gin-gonic/gin"
)

// Router handles all routing for the application
type Router struct {
	engine *gin.Engine
}

// NewRouter creates a new router instance
func NewRouter() *Router {
	engine := gin.Default()
	return &Router{engine: engine}
}

// SetupRoutes configures all the routes for the application
func (r *Router) SetupRoutes() {
	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/", nil)       // Create user
			users.GET("/:id", nil)     // Get user by ID
			users.PUT("/:id", nil)     // Update user
			users.DELETE("/:id", nil)  // Delete user
			users.GET("/", nil)        // List users
		}
	}
}

// Run starts the HTTP server
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
