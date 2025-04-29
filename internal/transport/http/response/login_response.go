// File: internal/transport/http/response/auth_response.go
package response

// LoginResponse defines the structure of the response after successful login.
type LoginResponse struct {
	Message       string `json:"message"`
	BackendToken  string `json:"backendToken,omitempty"` // Your backend's session token (e.g., JWT)
	UserID        string `json:"userId,omitempty"`       // Your internal user ID (optional)
	FirebaseUID   string `json:"firebaseUid,omitempty"`  // Firebase User ID (optional)
}

type ErrorResponse struct {
	Error string `json:"error"`
}