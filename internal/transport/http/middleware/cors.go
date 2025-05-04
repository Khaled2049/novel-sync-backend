// File: internal/transport/http/middleware/cors.go
package middleware

import (
	"time"

	"github.com/gin-contrib/cors" // Need to run: go get github.com/gin-contrib/cors
	"github.com/gin-gonic/gin"
)

// CORSMiddleware sets up Cross-Origin Resource Sharing policies.
func CORSMiddleware() gin.HandlerFunc {
	// Configure this carefully for production!
	// Use specific origins instead of AllowAllOrigins.
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowAllOrigins:  true, // Use this only for very open APIs or local testing
		MaxAge: 12 * time.Hour,
	})
}