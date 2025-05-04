// File: internal/transport/http/server.go
package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khaled2049/server/internal/config"
	"github.com/khaled2049/server/internal/transport/http/handlers"
	"github.com/khaled2049/server/internal/transport/http/middleware"
	// Import services if needed to create handlers here
	// "your_module_path_placeholder/internal/service"
)

// Server wraps the Gin engine and server configuration.
type Server struct {
	config  *config.ServerConfig
	engine  *gin.Engine
	httpSrv *http.Server
	// Add dependencies needed to create handlers (services, etc.)
	authHandler  *handlers.AuthHandler
	helloHandler *handlers.HelloHandler
	novelHandler *handlers.NovelHandler
}

// NewServer creates and configures a new HTTP server instance.
func NewServer(
	cfg *config.Config,
	// Pass dependencies needed by handlers
	authHandler *handlers.AuthHandler,
	helloHandler *handlers.HelloHandler,
	novelHandler *handlers.NovelHandler,

) *Server {
	// Set Gin mode (e.g., debug, release, test)
	// gin.SetMode(gin.ReleaseMode) // Uncomment for production

	engine := gin.Default() // Includes Logger and Recovery middleware
	engine.Use(middleware.CORSMiddleware())

	// Create server instance
	server := &Server{
		config:  &cfg.Server, // Store only server config
		engine:  engine,
		// Store handlers
		authHandler:  authHandler,
		helloHandler: helloHandler,
		novelHandler: novelHandler,
	}

	// --- Register Routes ---
	// Pass the engine and handlers to the central registration function
	RegisterAllRoutes(engine, authHandler, helloHandler, novelHandler)

	return server
}

// Run starts the HTTP server and listens for connections.
// It also handles graceful shutdown.
func (s *Server) Run() error {
	// Configure the underlying http.Server
	s.httpSrv = &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      s.engine,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	// --- Start server in a goroutine ---
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("HTTP server listening on port %s", s.config.Port)
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("failed to start server: %w", err)
		}
	}()

	// --- Wait for interrupt signal for graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received or the server fails to start
	select {
	case err := <-serverErr:
		return err // Return startup error immediately
	case sig := <-quit:
		log.Printf("Received signal %s. Shutting down server gracefully...", sig)
	}

	// --- Attempt graceful shutdown ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Shutdown timeout
	defer cancel()

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server exiting gracefully.")
	return nil
}
