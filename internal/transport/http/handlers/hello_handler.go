// File: internal/transport/http/handlers/hello_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloHandler handles simple, miscellaneous routes like /hello.
type HelloHandler struct {
	// No dependencies needed for this simple handler
}

// NewHelloHandler creates a new HelloHandler.
func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// RegisterRoutes registers the hello route.
func (h *HelloHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/hello", h.helloWorld)
}

// helloWorld is the handler function for GET /hello.
func (h *HelloHandler) helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
		"status":  "OK",
	})
}