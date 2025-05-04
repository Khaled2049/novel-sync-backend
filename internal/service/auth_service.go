// File: internal/service/auth_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"log" // Use a proper logger in production

	"github.com/khaled2049/server/internal/auth"
	"github.com/khaled2049/server/internal/domain"
	"github.com/khaled2049/server/internal/repository"
	"github.com/khaled2049/server/internal/util/jwt"
	"github.com/khaled2049/server/internal/util/password"
)

// AuthService handles authentication logic.
type AuthService struct {
	firebaseVerifier auth.FirebaseVerifier
	userRepo         repository.UserRepository
	jwtGenerator     *jwt.Generator
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	verifier auth.FirebaseVerifier,
	userRepo repository.UserRepository,
	jwtGen *jwt.Generator, // Inject JWT Generator
) *AuthService {
	return &AuthService{
		firebaseVerifier: verifier,
		userRepo:         userRepo,
		jwtGenerator:     jwtGen, 
	}
}


// LoginWithFirebaseToken verifies the Firebase token, finds or creates a user,
// and generates a backend session token.
func (s *AuthService) LoginWithFirebaseToken(ctx context.Context, idToken string) ( map[string]interface{}, error) {

	firebaseToken, err := s.firebaseVerifier.VerifyFirebaseIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid Firebase token here: %w", err)
	}

	// 2. Find or Create User in your Database
	// Use Firebase UID as the unique identifier linking to your internal user
	firebaseUID := firebaseToken.UID
	user, err := s.userRepo.FindByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		// Example: Handle "not found" by creating the user
		if err == repository.ErrUserNotFound { // Assume your repo defines/returns specific errors
			log.Printf("User with Firebase UID %s not found, creating new user...", firebaseUID)
			newUser := &domain.User{
				FirebaseUID: firebaseUID,
				Email:       firebaseToken.Claims["email"].(string), // Extract relevant info
				// Name:     firebaseToken.Claims["name"].(string), // Be careful with type assertions
				// Populate other fields as needed
			}

			createdUser, createErr := s.userRepo.Create(ctx, newUser)
			if createErr != nil {
				return nil, fmt.Errorf("failed to create user: %w", createErr)
			}
			user = createdUser // Use the newly created user
		} else {
			// Handle other potential database errors
			return nil, fmt.Errorf("error finding user by Firebase UID: %w", err)
		}
	} else {
		// Optional: Update user details if they've changed in Firebase
		// log.Printf("User found: %s", user.ID)
		// s.userRepo.Update(ctx, user) // Example update logic
	}

	// 3. Generate Backend Session Token (e.g., JWT)
	// This token will be used for authenticating subsequent requests to *your* API.
	// You would typically implement JWT generation/signing logic separately.
	// backendToken, err := s.jwtGenerator.GenerateToken(user.ID, user.Roles) // Example
	// if err != nil {
	// 	 return nil, fmt.Errorf("failed to generate backend token: %w", err)
	// }

	// For now, let's just return success data without a backend token
	log.Printf("User %s (Firebase UID: %s) logged in successfully.", user.ID, user.FirebaseUID)
	responseData := map[string]interface{}{
		"message":     "Login successful",
		"userId":      user.ID, // Your internal DB User ID
		"firebaseUid": user.FirebaseUID,
		// "backendToken": backendToken, // Include the generated token
	}

	return responseData, nil
}

func (s *AuthService) Login(ctx context.Context, email, plainPassword string) (map[string]interface{}, error) {
	// Basic Input Validation (optional here if done via binding)
	if email == "" || plainPassword == "" {
		return nil, fmt.Errorf("email and password are required") // Or a more generic error
	}

	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// IMPORTANT: Don't reveal if the user doesn't exist vs. other errors
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Printf("Login attempt failed: user not found for email %s", email)
			return nil, fmt.Errorf("invalid email or password")
		}
		// Log the real error but return a generic one
		log.Printf("Error fetching user by email %s during login: %v", email, err)
		return nil, fmt.Errorf("invalid email or password") // Or "internal server error" ?
	}

	// 2. Verify password
	// Ensure user object and password hash are not nil/empty before checking
	if user == nil || user.PasswordHash == "" {
        log.Printf("Error during login for email %s: user data or password hash missing", email)
        return nil, fmt.Errorf("authentication failed") // More generic internal-like error
    }
	// Verify the password hash using the password utility
	match := password.CheckPasswordHash(plainPassword, user.PasswordHash)

	if !match {
		log.Printf("Login attempt failed: incorrect password for email %s", email)
		return nil, fmt.Errorf("invalid email or password")
	}

	// 3. Authentication successful - Generate JWT token
	tokenString, err := s.jwtGenerator.GenerateToken(user.ID)
	if err != nil {
		// Log the real error
		log.Printf("Error generating JWT token for user %s after login: %v", user.ID, err)
		return nil, fmt.Errorf("login failed: could not generate session") // Internal error
	}

	// 4. Prepare response
	// (Optional: update last login time)
	log.Printf("User %s logged in successfully via standard login.", user.ID)
	responseData := map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString, // Send the JWT token back to the client
		"userId":  user.ID,     // Optionally send user ID
		// Add other non-sensitive user info if needed
	}

	return responseData, nil
}
