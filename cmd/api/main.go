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

	firebase "firebase.google.com/go/v4"
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

	"google.golang.org/api/option"
)

var firebaseAuthClient *auth.Client
var dbPool *pgxpool.Pool // Store the pool globally or pass via DI

func initializeFirebase(cfg *config.FirebaseConfig) error {
	
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
	
	err := godotenv.Load() // Loads .env from current directory or parent dirs
	if err != nil {
		
		log.Printf("Warning: could not load .env file: %v", err)
	}

	log.Println("Starting application...")
	// Create context for initialization steps with potential timeouts
	initCtx, initCancel := context.WithTimeout(context.Background(), 15*time.Second) // 15 sec timeout for init
	defer initCancel()

	
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	jwtGenerator, err := jwt.NewGenerator(&cfg.JWT)
    if err != nil {
        log.Fatalf("Failed to initialize JWT Generator: %v", err)
    }

	
	if err := initializeFirebase(&cfg.Firebase); err != nil {
		log.Printf("Firebase initialization failed: %v. Continuing...", err)
	}

	
	dbPool, err = postgres.NewConnectionPool(&cfg.Database, initCtx) // Use initCtx
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	defer func() {
		log.Println("Closing database connection pool...")
		dbPool.Close()
	}()

	
	userRepo := postgres.NewUserRepository(dbPool)
	novelRepo := postgres.NewNovelRepository(dbPool)
	chapterRepo := postgres.NewChapterRepository(dbPool)

	
	var firebaseVerifier fbAuth.FirebaseVerifier
	if firebaseAuthClient != nil {
		firebaseVerifier = fbAuth.NewFirebaseVerifier(firebaseAuthClient)
	} else {
		log.Println("Firebase Auth Client is nil, Firebase verification will not work.")
		// firebaseVerifier = fbAuth.NewNoopVerifier() // Placeholder if needed
	}

	
	
	authService := service.NewAuthService(firebaseVerifier, userRepo, jwtGenerator)
	novelService := service.NewNovelService(novelRepo, chapterRepo)

	
	authHandler := handlers.NewAuthHandler(authService)
	helloHandler := handlers.NewHelloHandler()
	novelHandler := handlers.NewNovelHandler(novelService) 

	
	srv := http.NewServer(cfg, authHandler, helloHandler, novelHandler)

	
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

	log.Println("Application shut down gracefully.")
}