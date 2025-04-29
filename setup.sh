#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Define Project Name (Optional - assumes execution in the desired root folder)
# Read project name if needed, e.g.: read -p "Enter project name (module path like github.com/youruser/yourproject): " project_name
# Or get from argument: project_name=$1

echo "Setting up Go project structure..."

# --- Root Files ---
echo "Creating root files..."
touch README.md
touch .gitignore
touch .env.example
touch docker-compose.yml
# Initialize Go module (Replace 'your_module_path' with your actual Go module path)
# Example: github.com/yourusername/novel-platform-backend
go mod init your_module_path_placeholder
echo "Placeholder: Remember to replace 'your_module_path_placeholder' in go.mod"

# --- Top Level Directories ---
echo "Creating top-level directories..."
mkdir -p api       # For OpenAPI/Swagger specs, protobuf definitions, etc.
mkdir -p cmd       # Application entry points (main functions)
mkdir -p configs   # Configuration file templates or defaults
mkdir -p deployments # Deployment configs (Kubernetes manifests, Terraform, etc.)
mkdir -p docs      # Project documentation
mkdir -p internal  # Private application code (core logic)
mkdir -p migrations # Database migration files
mkdir -p scripts   # Helper scripts (like this one, build scripts, etc.)
mkdir -p test      # Additional tests (E2E, integration beyond _test.go)
# mkdir -p web       # If serving web assets directly (JS, CSS, templates)
# mkdir -p pkg       # Only if you have library code intended for external use

# --- API Directory ---
echo "Setting up api/ directory..."
touch api/.gitkeep # Or add your spec files here later (e.g., api/openapi.yaml)

# --- CMD Directory (Application Binaries) ---
echo "Setting up cmd/ directory..."
# API Server
mkdir -p cmd/api
touch cmd/api/main.go
echo "// File: cmd/api/main.go" > cmd/api/main.go
# AI Worker
mkdir -p cmd/aiworker
touch cmd/aiworker/main.go
echo "// File: cmd/aiworker/main.go" > cmd/aiworker/main.go
# Notification Worker (Optional)
mkdir -p cmd/notifier
touch cmd/notifier/main.go
echo "// File: cmd/notifier/main.go" > cmd/notifier/main.go

# --- Configs Directory ---
echo "Setting up configs/ directory..."
touch configs/config.example.yaml # Example config structure

# --- Deployments Directory ---
echo "Setting up deployments/ directory..."
touch deployments/README.md # Explain deployment strategies

# --- Docs Directory ---
echo "Setting up docs/ directory..."
touch docs/architecture.md # Place for architecture diagrams/notes
touch docs/api_endpoints.md # Place for detailed API docs

# --- Internal Directory (Core Logic) ---
echo "Setting up internal/ directory..."
mkdir -p internal/auth              # Authentication logic (Firebase integration)
mkdir -p internal/config            # Config loading/management (e.g., Viper)
mkdir -p internal/domain            # Core domain models/structs (Novel, User, Chapter, etc.)
mkdir -p internal/llm               # Client for interacting with external LLM APIs
mkdir -p internal/mq                # Message Queue interfaces and clients (e.g., RabbitMQ)
mkdir -p internal/repository        # Data Repositories (interfaces define contracts)
mkdir -p internal/service           # Business logic layer (orchestrates repositories, MQ, etc.)
mkdir -p internal/transport         # Network transport handlers (HTTP, gRPC, etc.)
mkdir -p internal/util              # Shared utility functions

touch internal/auth/firebase.go
touch internal/config/config.go
touch internal/domain/novel.go # Add other domain files as needed
touch internal/domain/user.go
touch internal/llm/client.go
touch internal/mq/producer.go
touch internal/mq/consumer.go
touch internal/repository/novel_repo.go # Repository Interface definition
touch internal/repository/user_repo.go  # Add other interfaces as needed
touch internal/service/novel_service.go # Business logic for Novels
touch internal/service/ai_service.go    # Business logic for AI features
touch internal/transport/README.md      # Explain transport layers used
touch internal/util/README.md

# --- Internal Sub-directories ---
echo "Setting up internal/ sub-directories..."
# Message Queue Implementation (RabbitMQ)
mkdir -p internal/mq/rabbitmq
touch internal/mq/rabbitmq/connection.go
touch internal/mq/rabbitmq/publisher.go
touch internal/mq/rabbitmq/consumer.go

# Repository Implementation (Postgres)
mkdir -p internal/repository/postgres
touch internal/repository/postgres/db.go           # DB connection setup
touch internal/repository/postgres/novel_repo.go   # Postgres implementation of NovelRepository
touch internal/repository/postgres/user_repo.go    # Postgres implementation of UserRepository

# Transport Implementation (HTTP)
mkdir -p internal/transport/http
mkdir -p internal/transport/http/handlers   # API request handlers (Gin/Echo/etc.)
mkdir -p internal/transport/http/middleware # HTTP middleware (Auth, Logging, CORS, etc.)
mkdir -p internal/transport/http/request    # Request struct definitions/validation
mkdir -p internal/transport/http/response   # Response struct definitions/formatting
touch internal/transport/http/server.go       # HTTP server setup
touch internal/transport/http/routes.go       # Route definitions
touch internal/transport/http/handlers/novel_handler.go
touch internal/transport/http/handlers/auth_handler.go
touch internal/transport/http/middleware/auth_middleware.go
touch internal/transport/http/request/novel_request.go
touch internal/transport/http/response/novel_response.go

# --- Migrations Directory ---
echo "Setting up migrations/ directory..."
touch migrations/README.md # Explain how to run migrations
# Migration files (e.g., 000001_create_users_table.up.sql) will go here

# --- Scripts Directory ---
echo "Setting up scripts/ directory..."
# Copy this script into the scripts directory if desired
# cp setup_structure.sh scripts/
touch scripts/README.md # Document available scripts
touch scripts/build.sh  # Example build script placeholder
touch scripts/run.sh    # Example run script placeholder

# --- Test Directory ---
echo "Setting up test/ directory..."
touch test/README.md # Explain testing strategy (E2E, integration)

# --- .gitignore Content ---
echo "Creating .gitignore..."
cat << EOF > .gitignore
# Go build outputs
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test

# Output folder (if compiling locally)
/bin/
/dist/

# Go workspace history
.history/

# Env files
.env
*.env

# IDE and OS files
.idea/
.vscode/
*.iml
*.DS_Store
Thumbs.db

# Dependency directories (though 'go mod vendor' is less common now)
/vendor/

# Log files
*.log

# Temporary files
*.tmp
EOF

echo "Project structure setup complete!"
echo "Next steps:"
echo "1. Replace 'your_module_path_placeholder' in go.mod with your actual module path."
echo "2. Review the structure and adjust if needed."
echo "3. Populate the placeholder files with actual code."
echo "4. Configure your .gitignore further if necessary."
echo "5. Initialize Git if you haven't already ('git init', 'git add .', 'git commit -m \"Initial project structure\"')."

exit 0