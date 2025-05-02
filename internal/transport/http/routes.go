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
	authHandler *handlers.AuthHandler, 
	helloHandler *handlers.HelloHandler,
	novelHandler *handlers.NovelHandler, 
) {
	// Initialize handlers
	helloHandler.RegisterRoutes(router) 
	authHandler.RegisterRoutes(router)  
	novelHandler.RegisterRoutes(router)


	// Add health check endpoint (common practice)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
}

