// File: cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/jackc/pgx/v5/pgxpool" // Import pgxpool
	"github.com/joho/godotenv"

	// Import necessary packages from your project
	fbAuth "github.com/khaled2049/server/internal/auth"
	"github.com/khaled2049/server/internal/config"
	"github.com/khaled2049/server/internal/repository/postgres"
	"github.com/khaled2049/server/internal/service"
	"github.com/khaled2049/server/internal/transport/http"
	"github.com/khaled2049/server/internal/transport/http/handlers"
	"github.com/khaled2049/server/internal/util/jwt"

	// --- Add firebase imports ---
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var firebaseAuthClient *auth.Client
var dbPool *pgxpool.Pool // Store the pool globally or pass via DI

func initializeFirebase(cfg *config.FirebaseConfig) error {
	// ... (Firebase init logic as before) ...
    serviceAccountKeyPath := cfg.ServiceAccountKeyPath
	if serviceAccountKeyPath == "" {
		serviceAccountKeyPath = os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY_PATH")
		if serviceAccountKeyPath == "" {
			log.Println("Warning: FIREBASE_SERVICE_ACCOUNT_KEY_PATH not set. Firebase Admin SDK will not be initialized.")
			return nil
		}
	}
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing Firebase app: %w", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		return fmt.Errorf("error getting Firebase Auth client: %w", err)
	}
	firebaseAuthClient = client
	log.Println("Firebase Admin SDK initialized successfully.")
	return nil
}

func main() {
	// Load .env file variables into the environment for this process
	err := godotenv.Load() // Loads .env from current directory or parent dirs
	if err != nil {
		// It's often okay to just log a warning if .env is optional (e.g., relying on system env vars)
		log.Printf("Warning: could not load .env file: %v", err)
	}

	log.Println("Starting application...")
	// Create context for initialization steps with potential timeouts
	initCtx, initCancel := context.WithTimeout(context.Background(), 15*time.Second) // 15 sec timeout for init
	defer initCancel()

	// --- 1. Load Configuration ---
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	jwtGenerator, err := jwt.NewGenerator(&cfg.JWT)
    if err != nil {
        log.Fatalf("Failed to initialize JWT Generator: %v", err)
    }

	// --- 2. Initialize Firebase ---
	if err := initializeFirebase(&cfg.Firebase); err != nil {
		log.Printf("Firebase initialization failed: %v. Continuing...", err)
	}

	// --- 3. Initialize Database Connection Pool ---
	dbPool, err = postgres.NewConnectionPool(&cfg.Database, initCtx) // Use initCtx
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Defer closing the pool until the main function exits
	defer func() {
		log.Println("Closing database connection pool...")
		dbPool.Close()
	}()

	// --- 4. Setup Repositories ---
	userRepo := postgres.NewUserRepository(dbPool)
	// novelRepo := postgres.NewNovelRepository(dbPool) // Create this later

	// --- 5. Create Firebase Verifier ---
	var firebaseVerifier fbAuth.FirebaseVerifier
	if firebaseAuthClient != nil {
		firebaseVerifier = fbAuth.NewFirebaseVerifier(firebaseAuthClient)
	} else {
		log.Println("Firebase Auth Client is nil, Firebase verification will not work.")
		// firebaseVerifier = fbAuth.NewNoopVerifier() // Placeholder if needed
	}

	// --- 6. Create Services ---
	// Inject the concrete userRepo implementation
	authService := service.NewAuthService(firebaseVerifier, userRepo, jwtGenerator)
	// novelService := service.NewNovelService(novelRepo)

	// --- 7. Create HTTP Handlers ---
	authHandler := handlers.NewAuthHandler(authService)
	helloHandler := handlers.NewHelloHandler()
	// novelHandler := handlers.NewNovelHandler(novelService)

	// --- 8. Create and Prepare HTTP Server ---
	srv := http.NewServer(cfg, authHandler, helloHandler /*, novelHandler */)

	// --- 9. Start Server and Handle Graceful Shutdown ---
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("Starting HTTP server...")
		serverErrors <- srv.Run() // srv.Run now handles its own shutdown logic internally
	}()

	// Wait for shutdown signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server runtime error: %v", err) // If server fails to start or runtime error occurs
	case sig := <-quit:
		log.Printf("Received signal %s. Application shutting down...", sig)
	}

	// Server shutdown is handled within srv.Run(),
	// Database pool closure is handled by the defer statement.

	log.Println("Application shut down gracefully.")
}