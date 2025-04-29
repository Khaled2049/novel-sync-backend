// File: internal/repository/postgres/db.go
package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool" // Using pgx v5 pool

	"github.com/khaled2049/server/internal/config"
)

// NewConnectionPool creates a new PostgreSQL connection pool.
func NewConnectionPool(cfg *config.DatabaseConfig, ctx context.Context) (*pgxpool.Pool, error) {
	// Example DSN: postgres://user:password@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Configure the pool settings
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config from DSN: %w", err)
	}

	// Example: Set max connections, idle time, etc.
	poolConfig.MaxConns = 10                  // Example: Limit max connections
	poolConfig.MaxConnIdleTime = 5 * time.Minute // Example: Close idle connections

	log.Println("Connecting to database...")

	// Connect using the configured pool settings
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig) // Use background context for pool creation
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ping the database to verify connection (use the provided context with timeout)
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // Use context passed to main/init
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close() // Close the pool if ping fails
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection successful.")
	return pool, nil
}