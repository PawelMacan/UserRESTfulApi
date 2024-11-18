package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Metrics middleware adds request timing and error tracking
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log timing and status
		duration := time.Since(start)
		status := c.Writer.Status()
		path := c.Request.URL.Path
		method := c.Request.Method

		if status >= 400 {
			log.Printf("ERROR [%s] %s %d %v", method, path, status, duration)
		} else if duration > time.Millisecond*500 {
			log.Printf("SLOW [%s] %s %d %v", method, path, status, duration)
		} else {
			log.Printf("INFO [%s] %s %d %v", method, path, status, duration)
		}
	}
}
