package request

type FirebaseLoginRequest struct {
	IDToken string `json:"idToken" binding:"required"` // Use 'binding' tag for validation if using Gin
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"` // Use email validation
	Password string `json:"password" binding:"required"`
}