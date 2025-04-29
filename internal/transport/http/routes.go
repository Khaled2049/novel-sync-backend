// File: internal/transport/http/routes.go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaled2049/server/internal/transport/http/handlers"
)

// RegisterAllRoutes initializes all handlers and registers their routes.

func RegisterAllRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler, // Assuming authHandler is created in server.go or main.go
	helloHandler *handlers.HelloHandler,
) {
	// Initialize handlers
	helloHandler.RegisterRoutes(router) // Or pass 'api' group if using grouping
	authHandler.RegisterRoutes(router)  // Or pass 'api' group
	// novelHandler.RegisterRoutes(api)

	// Add health check endpoint (common practice)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
}

