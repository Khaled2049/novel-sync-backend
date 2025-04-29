// File: internal/transport/http/handlers/auth_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaled2049/server/internal/service"
	"github.com/khaled2049/server/internal/transport/http/request"
	"github.com/khaled2049/server/internal/transport/http/response"
)

// AuthHandler handles authentication related HTTP requests.
type AuthHandler struct {
	authService *service.AuthService // Inject your AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRoutes registers authentication routes with the Gin engine.
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login/firebase", h.FirebaseLogin)
		authGroup.POST("/login", h.StandardLogin)
	}
}

func (h *AuthHandler) StandardLogin(c *gin.Context) {
    var req request.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request body: " + err.Error()})
        return
    }
    
    responseData, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, responseData)
}

// FirebaseLogin handles the POST /auth/login/firebase request.
func (h *AuthHandler) FirebaseLogin(c *gin.Context) {
	var req request.FirebaseLoginRequest

	// Bind the incoming JSON payload to the request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	// Call the authentication service
	responseData, err := h.authService.LoginWithFirebaseToken(c.Request.Context(), req.IDToken)
	if err != nil {
		// Handle different error types potentially (e.g., invalid token vs. internal error)
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
		return
	}

	// On success, return the response data (which might include your backend token)
	// You might set the backend token in an HttpOnly secure cookie here instead of the body
	c.JSON(http.StatusOK, responseData) // Use response.LoginResponse if you created one
}