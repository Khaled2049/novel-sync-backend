// File: internal/auth/firebase.go
package auth

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

// FirebaseVerifier defines the interface for verifying Firebase ID tokens.
// This helps in testing by allowing mocking.
type FirebaseVerifier interface {
	VerifyFirebaseIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

// firebaseVerifier implements FirebaseVerifier using the Firebase Admin SDK.
type firebaseVerifier struct {
	authClient *auth.Client
}

// NewFirebaseVerifier creates a new verifier instance.
// Pass the initialized auth.Client from main.go here.
func NewFirebaseVerifier(client *auth.Client) FirebaseVerifier {
	if client == nil {
		// Or handle this more gracefully depending on your DI setup
		panic("authClient cannot be nil")
	}
	return &firebaseVerifier{authClient: client}
}

// VerifyFirebaseIDToken uses the Firebase Admin SDK to verify the given ID token.
func (v *firebaseVerifier) VerifyFirebaseIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if idToken == "" {
		return nil, fmt.Errorf("ID token cannot be empty")
	}

	token, err := v.authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		// Handle specific Firebase errors if needed, e.g., token expired
		return nil, fmt.Errorf("error verifying Firebase ID token: %w", err)
	}

	return token, nil
}