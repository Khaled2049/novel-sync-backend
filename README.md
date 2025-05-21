# NovelCraft: System Architecture & API Documentation

## Table of Contents
1. [System Overview](#system-overview)
2. [Architecture](#architecture)
    - [High-Level Architecture](#high-level-architecture)
    - [Component Interaction](#component-interaction)
3. [Core Components](#core-components)
    - [API Server](#api-server)
    - [AI Worker](#ai-worker)
    - [Notification Service](#notification-service)
    - [Database Design](#database-design)
4. [API Reference](#api-reference)
    - [User Management](#user-management)
    - [Novel Management](#novel-management)
    - [Chapters & Content](#chapters--content)
    - [Worldbuilding](#worldbuilding)
    - [AI Integration](#ai-integration)
    - [Collaboration](#collaboration)
5. [Directory Structure](#directory-structure)
    - [Component Placement](#component-placement)
    - [File Responsibilities](#file-responsibilities)
6. [Development Guidelines](#development-guidelines)
    - [Best Practices](#best-practices)
    - [Extension Points](#extension-points)

## System Overview

NovelCraft is a comprehensive novel writing application designed to assist authors in creating and managing their literary works. The system combines traditional writing tools with AI assistance capabilities through large language models (LLMs).

**Key Features:**
- Novel management with chapters, revisions, and content tracking
- Character, place, and worldbuilding management
- Collaborative writing with role-based permissions
- AI-powered suggestions for continuations, dialogue, descriptions, and more
- Timeline management for story events
- Comments and annotations for collaborative feedback
- User private notes for personal insights

## Architecture

### High-Level Architecture

NovelCraft follows a modern, scalable architecture pattern with clear separation of concerns:

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│                 │     │                 │     │                 │
│  Web Frontend   │◄────┤    API Server   │◄────┤    Database     │
│                 │     │                 │     │   (PostgreSQL)  │
└────────┬────────┘     └────────┬────────┘     └─────────────────┘
         │                       │                       ▲
         │                       │                       │
         │                       ▼                       │
         │              ┌─────────────────┐              │
         │              │                 │              │
         └─────────────►│  Message Queue  │◄─────────────┘
                        │   (RabbitMQ)    │
                        └────────┬────────┘
                                 │
                                 ▼
                      ┌──────────────────────┐
                      │                      │
                      │     AI Worker(s)     │
                      │                      │
                      └──────────┬───────────┘
                                 │
                                 ▼
                      ┌──────────────────────┐
                      │                      │
                      │     LLM Service      │
                      │                      │
                      └──────────────────────┘
```

### Component Interaction

The system operates through the following flow:

1. **User interaction**: Frontend interfaces with API Server through RESTful HTTP endpoints
2. **Data persistence**: API Server reads/writes to PostgreSQL database
3. **Asynchronous processing**: Long-running tasks like AI generation are handled through a message queue
4. **AI processing**: Dedicated AI workers consume messages, process them through LLM services, and store results
5. **Notifications**: Updates and real-time events are managed through the notification system

## Core Components

### API Server

The primary interface for client applications, providing RESTful endpoints for all system functions:

- **Authentication**: Integrated with Firebase Authentication
- **Resource management**: CRUD operations for novels, characters, places, etc.
- **Collaboration**: User permissions, comments, and sharing
- **AI request submission**: Endpoints to request AI assistance (processed asynchronously)

### AI Worker

Dedicated service for processing AI generation requests:

- **Context assembly**: Gathers relevant novel content, characters, settings, etc.
- **Prompt engineering**: Crafts effective prompts based on user requests
- **LLM interaction**: Communicates with external LLM APIs
- **Response processing**: Validates, filters, and stores generated content

### Notification Service

Manages real-time updates and notifications:

- **Collaboration alerts**: Notifies users of comments, edits, etc.
- **AI completion notifications**: Alerts when AI suggestions are ready
- **System messages**: Provides status updates on operations

### Database Design

The PostgreSQL schema is designed to support all aspects of novel writing:

- **User management**: Authentication and profile data
- **Content management**: Novels, chapters, revisions with version history
- **Worldbuilding**: Characters, places, notes, relationships
- **Collaboration**: Comments, permissions, shared resources
- **AI integration**: Suggestion tracking, context storage, and feedback

## API Reference

### User Management

#### `POST /auth/register`
- **Purpose**: Register a new user with Firebase auth integration
- **File**: `internal/transport/http/handlers/auth_handler.go`
- **Implementation**: Integrates with Firebase to register and create user record

#### `POST /auth/login`
- **Purpose**: Authenticate user and return JWT token
- **File**: `internal/transport/http/handlers/auth_handler.go`
- **Implementation**: Verifies Firebase credentials and issues JWT

#### `GET /users/me`
- **Purpose**: Retrieve current user profile
- **File**: `internal/transport/http/handlers/user_handler.go`
- **Implementation**: Returns user data from database based on JWT claims

#### `PUT /users/me`
- **Purpose**: Update user profile information
- **File**: `internal/transport/http/handlers/user_handler.go`
- **Implementation**: Updates user record in database

### Novel Management

#### `POST /novels`
- **Purpose**: Create a new novel
- **File**: `internal/transport/http/handlers/novel_handler.go`
- **Implementation**: Creates a novel record and establishes ownership

#### `GET /novels`
- **Purpose**: List all novels accessible to the user
- **File**: `internal/transport/http/handlers/novel_handler.go`
- **Implementation**: Queries novels user owns or collaborates on

#### `GET /novels/:id`
- **Purpose**: Get detailed information about a specific novel
- **File**: `internal/transport/http/handlers/novel_handler.go`
- **Implementation**: Returns novel data with permission verification

#### `PUT /novels/:id`
- **Purpose**: Update novel metadata
- **File**: `internal/transport/http/handlers/novel_handler.go`
- **Implementation**: Updates novel title, description, etc.

#### `DELETE /novels/:id`
- **Purpose**: Delete a novel
- **File**: `internal/transport/http/handlers/novel_handler.go`
- **Implementation**: Soft/hard delete with permission checks

#### `POST /novels/:id/collaborators`
- **Purpose**: Add a collaborator to a novel
- **File**: `internal/transport/http/handlers/collaborator_handler.go`
- **Implementation**: Creates collaboration record with role

### Chapters & Content

#### `POST /novels/:id/chapters`
- **Purpose**: Create a new chapter
- **File**: `internal/transport/http/handlers/chapter_handler.go`
- **Implementation**: Creates chapter record with order management

#### `GET /novels/:id/chapters`
- **Purpose**: List all chapters in a novel
- **File**: `internal/transport/http/handlers/chapter_handler.go`
- **Implementation**: Returns ordered chapters with metadata

#### `GET /chapters/:id`
- **Purpose**: Get chapter content and metadata
- **File**: `internal/transport/http/handlers/chapter_handler.go`
- **Implementation**: Returns chapter with content and linked entities

#### `PUT /chapters/:id`
- **Purpose**: Update chapter content
- **File**: `internal/transport/http/handlers/chapter_handler.go`
- **Implementation**: Updates content and creates revision record

#### `PUT /chapters/:id/order`
- **Purpose**: Reorder chapter position
- **File**: `internal/transport/http/handlers/chapter_handler.go`
- **Implementation**: Updates order indices with transaction

#### `GET /chapters/:id/revisions`
- **Purpose**: List chapter revision history
- **File**: `internal/transport/http/handlers/chapter_handler.go`
- **Implementation**: Returns revision history with timestamps

#### `GET /chapters/:id/revisions/:revision_id`
- **Purpose**: Get specific revision content
- **File**: `internal/transport/http/handlers/chapter_handler.go`
- **Implementation**: Returns historical content

### Worldbuilding

#### `POST /novels/:id/characters`
- **Purpose**: Create a character
- **File**: `internal/transport/http/handlers/character_handler.go`
- **Implementation**: Creates character record

#### `GET /novels/:id/characters`
- **Purpose**: List all characters in a novel
- **File**: `internal/transport/http/handlers/character_handler.go`
- **Implementation**: Returns all characters with basic info

#### `GET /characters/:id`
- **Purpose**: Get detailed character information
- **File**: `internal/transport/http/handlers/character_handler.go`
- **Implementation**: Returns full character data with relationships

#### `PUT /characters/:id`
- **Purpose**: Update character information
- **File**: `internal/transport/http/handlers/character_handler.go`
- **Implementation**: Updates character record

#### `POST /novels/:id/places`
- **Purpose**: Create a place
- **File**: `internal/transport/http/handlers/place_handler.go`
- **Implementation**: Creates place record

#### `GET /novels/:id/places`
- **Purpose**: List all places in a novel
- **File**: `internal/transport/http/handlers/place_handler.go`
- **Implementation**: Returns places with basic info

#### `GET /places/:id`
- **Purpose**: Get detailed place information
- **File**: `internal/transport/http/handlers/place_handler.go`
- **Implementation**: Returns place with full details

#### `POST /novels/:id/notes`
- **Purpose**: Create worldbuilding or research note
- **File**: `internal/transport/http/handlers/note_handler.go`
- **Implementation**: Creates note with proper type and links

#### `GET /novels/:id/notes`
- **Purpose**: List notes by type
- **File**: `internal/transport/http/handlers/note_handler.go`
- **Implementation**: Returns filtered notes by type

#### `POST /novels/:id/timeline`
- **Purpose**: Create timeline event
- **File**: `internal/transport/http/handlers/timeline_handler.go`
- **Implementation**: Creates event with proper ordering

#### `GET /novels/:id/timeline`
- **Purpose**: Get timeline events
- **File**: `internal/transport/http/handlers/timeline_handler.go`
- **Implementation**: Returns ordered events with links

### AI Integration

#### `POST /ai/suggestions`
- **Purpose**: Request AI assistance 
- **File**: `internal/transport/http/handlers/ai_handler.go`
- **Implementation**: Creates suggestion record and publishes message to queue

#### `GET /ai/suggestions/:id`
- **Purpose**: Get AI suggestion status and content
- **File**: `internal/transport/http/handlers/ai_handler.go`
- **Implementation**: Returns suggestion with status and generated content

#### `PUT /ai/suggestions/:id/accept`
- **Purpose**: Accept and apply AI suggestion
- **File**: `internal/transport/http/handlers/ai_handler.go`
- **Implementation**: Updates suggestion status and applies content (if to chapter)

#### `PUT /ai/suggestions/:id/reject`
- **Purpose**: Reject AI suggestion
- **File**: `internal/transport/http/handlers/ai_handler.go`
- **Implementation**: Updates status with feedback

#### `POST /ai/suggestions/:id/edit`
- **Purpose**: Edit AI suggestion before accepting
- **File**: `internal/transport/http/handlers/ai_handler.go`
- **Implementation**: Updates suggestion content before applying

### Collaboration

#### `POST /chapters/:id/comments`
- **Purpose**: Add comment to a chapter
- **File**: `internal/transport/http/handlers/comment_handler.go`
- **Implementation**: Creates comment with proper links

#### `GET /chapters/:id/comments`
- **Purpose**: Get comments for a chapter
- **File**: `internal/transport/http/handlers/comment_handler.go`
- **Implementation**: Returns comments with threading

#### `POST /comments/:id/replies`
- **Purpose**: Reply to a comment
- **File**: `internal/transport/http/handlers/comment_handler.go`
- **Implementation**: Creates threaded reply

#### `PUT /comments/:id`
- **Purpose**: Edit a comment
- **File**: `internal/transport/http/handlers/comment_handler.go`
- **Implementation**: Updates comment content

#### `POST /novels/:id/private-notes`
- **Purpose**: Create private note visible only to user
- **File**: `internal/transport/http/handlers/private_note_handler.go`
- **Implementation**: Creates note with user-only visibility

## Directory Structure

### Component Placement

```
novel-writing-platform/
├── api/                 # API specifications (OpenAPI/Swagger)
├── cmd/                 # Application entry points
│   ├── api/             # API server main
│   ├── aiworker/        # AI worker service main
│   └── notifier/        # Notification service main
├── configs/             # Configuration templates
├── deployments/         # Deployment configurations
├── docs/                # Documentation
├── internal/            # Core application code
│   ├── auth/            # Authentication logic
│   ├── config/          # Configuration loading
│   ├── domain/          # Core business models
│   ├── llm/             # LLM integration client
│   ├── mq/              # Message queue abstractions
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic layer
│   ├── transport/       # Communication layer (HTTP)
│   └── util/            # Shared utilities
├── migrations/          # Database migrations
├── scripts/             # Helper scripts
└── test/                # Integration/E2E tests
```

### File Responsibilities

#### `/cmd` Directory - Entry Points

- `cmd/api/main.go`
  - **Role**: Main entry point for API server
  - **Responsibility**: Initializes and runs HTTP server, connects dependencies

- `cmd/aiworker/main.go`
  - **Role**: Main entry point for AI generation worker
  - **Responsibility**: Sets up message consumers, processes AI requests

- `cmd/notifier/main.go`
  - **Role**: Main entry point for notification service
  - **Responsibility**: Processes events, sends notifications

#### `/internal/domain` - Core Models

- `internal/domain/novel.go`
  - **Role**: Core Novel domain model
  - **Responsibility**: Defines Novel entity and related structures

- `internal/domain/user.go`
  - **Role**: User domain model
  - **Responsibility**: Defines User entity and authentication structures

- `internal/domain/chapter.go`
  - **Role**: Chapter model and content structures
  - **Responsibility**: Defines Chapter entity with revision tracking

- `internal/domain/character.go`, `internal/domain/place.go`, etc.
  - **Role**: Worldbuilding entity models
  - **Responsibility**: Define entity structures for worldbuilding

#### `/internal/repository` - Data Access

- `internal/repository/novel_repo.go`
  - **Role**: Novel repository interface
  - **Responsibility**: Defines data access contract for novels

- `internal/repository/postgres/novel_repo.go`
  - **Role**: PostgreSQL implementation of novel repository
  - **Responsibility**: Implements data access for novels using SQL

- (Similar pattern for user, chapter, character, etc.)

#### `/internal/service` - Business Logic

- `internal/service/novel_service.go`
  - **Role**: Novel business logic
  - **Responsibility**: Orchestrates repositories and implements domain logic

- `internal/service/ai_service.go`
  - **Role**: AI integration service
  - **Responsibility**: Handles prompt generation and LLM responses

#### `/internal/transport` - API Layer

- `internal/transport/http/server.go`
  - **Role**: HTTP server configuration
  - **Responsibility**: Sets up routes, middleware, server options

- `internal/transport/http/handlers/novel_handler.go`
  - **Role**: Novel API endpoints
  - **Responsibility**: Handles HTTP requests for novel resources

- `internal/transport/http/middleware/auth_middleware.go`
  - **Role**: Authentication middleware
  - **Responsibility**: Validates JWT tokens, enforces permissions

#### `/internal/llm` - LLM Integration

- `internal/llm/client.go`
  - **Role**: LLM client interface
  - **Responsibility**: Defines contract for LLM interactions

- `internal/llm/openai.go`
  - **Role**: OpenAI implementation
  - **Responsibility**: Implements LLM client for OpenAI

#### `/internal/mq` - Message Queue

- `internal/mq/producer.go`
  - **Role**: Message producer interface
  - **Responsibility**: Defines contract for sending messages

- `internal/mq/rabbitmq/publisher.go`
  - **Role**: RabbitMQ implementation of producer
  - **Responsibility**: Implements message publishing

## Development Guidelines

### Best Practices

1. **Dependency Injection**
   - Use constructor injection for dependencies
   - Create interfaces for testability and flexibility

2. **Error Handling**
   - Return domain-specific errors
   - Log errors with context
   - Map errors to appropriate HTTP status codes

3. **Context Propagation**
   - Pass context.Context through all layers
   - Use for cancellation, timeouts, and tracing

4. **Testing**
   - Unit test all repository implementations
   - Mock external dependencies
   - Use integration tests for data flow verification

### Extension Points

1. **LLM Providers**
   - Extend `internal/llm` with additional provider implementations
   - Implement the LLM client interface for the new provider

2. **AI Capabilities**
   - Add new suggestion types in the `ai_suggestion_type` enum
   - Implement corresponding context gathering and prompt generation

3. **Export Formats**
   - Create exporters in a new `internal/export` package
   - Add endpoints for different export formats (PDF, EPUB, etc.)

4. **Analytics**
   - Add analytics collectors in a new `internal/analytics` package
   - Implement event tracking for user actions and system metrics

5. **Media Support**
   - Extend schema for image/media storage
   - Add media handling services for upload/processing
