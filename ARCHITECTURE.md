# Mini Lab API Architecture

This document explains the layered architecture of the Mini Lab API and what should go in each layer.

## Overview

The Mini Lab API follows a clean, layered architecture pattern that separates concerns and makes the codebase maintainable, testable, and scalable. The architecture consists of the following layers:

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer                             │
│                   (cmd/api/main.go)                        │
│                    Gin Router, Middleware                   │
└─────────────────────────────┬───────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                    Handler Layer                            │
│                (internal/handlers/)                         │
│           HTTP Request/Response, Validation                 │
└─────────────────────────────┬───────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                   Service Layer                             │
│                 (internal/service/)                         │
│              Business Logic, Validation                     │
└─────────────────────────────┬───────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                 Repository Layer                            │
│               (internal/repository/)                        │
│               Data Access, Database Operations              │
└─────────────────────────────┬───────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                   Database Layer                            │
│            sqlc Generated Code (internal/db/)               │
│               SQL Queries, Database Models                  │
└─────────────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### 1. HTTP Layer (`cmd/api/main.go`)

**Purpose**: Application entry point and HTTP server setup

**What it contains:**
- Application initialization
- Database connection setup
- Dependency injection (wiring layers together)
- Gin router configuration
- Middleware setup (CORS, logging, recovery)
- Route definitions
- Server startup

**Example:**
```go
func main() {
    // Initialize config
    cfg := config.New()
    
    // Setup database
    db, err := config.NewDatabase(cfg)
    
    // Wire dependencies
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    handlers := handlers.New(userService)
    
    // Setup routes
    r := gin.Default()
    api := r.Group("/api/v1")
    api.GET("/users", handlers.GetUsers)
}
```

### 2. Handler Layer (`internal/handlers/`)

**Purpose**: HTTP request/response handling and basic validation

**What it contains:**
- HTTP request parsing
- Input validation (JSON binding)
- Calling service layer methods
- HTTP response formatting
- Error handling and status codes
- Swagger documentation annotations

**File Organization:**
- `handler.go` - Main Handler struct, constructor, and shared types
- `ping.go` - Health check and general endpoints
- `users.go` - User-specific HTTP handlers
- `products.go` - Product-specific HTTP handlers (when added)
- `orders.go` - Order-specific HTTP handlers (when added)

**Naming Convention:**
- Use entity names in plural form (e.g., `users.go`, `products.go`)
- Group related endpoints by domain/entity
- Keep shared functionality in `handler.go`

**What it should NOT contain:**
- Business logic
- Database operations
- Complex validation rules

**Example:**
```go
func (h *Handler) CreateUser(c *gin.Context) {
    var req models.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.userService.CreateUser(c.Request.Context(), req)
    if err != nil {
        // Handle different error types
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, gin.H{"data": user})
}
```

### 3. Service Layer (`internal/service/`)

**Purpose**: Business logic implementation and orchestration

**What it contains:**
- Business rules and validation
- Data transformation and normalization
- Orchestrating multiple repository calls
- Transaction management
- Business-specific error handling
- Cross-cutting concerns (logging, metrics)

**What it should NOT contain:**
- HTTP-specific code
- Database-specific code
- UI concerns

**Example:**
```go
func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
    // Business validation
    if err := s.validateCreateUserRequest(req); err != nil {
        return nil, err
    }
    
    // Data normalization
    req.Email = strings.ToLower(strings.TrimSpace(req.Email))
    
    // Business logic: check if email already exists
    _, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err == nil {
        return nil, ErrUserEmailExists
    }
    
    // Create user
    return s.userRepo.Create(ctx, req)
}
```

### 4. Repository Layer (`internal/repository/`)

**Purpose**: Data access abstraction and database operations

**What it contains:**
- Database query execution
- Data mapping between database and domain models
- Database-specific error handling
- Transaction management
- CRUD operations

**What it should NOT contain:**
- Business logic
- HTTP-specific code
- Complex business validation

**Example:**
```go
func (r *userRepository) Create(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
    dbUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
        Name:  req.Name,
        Email: req.Email,
    })
    if err != nil {
        return nil, err
    }
    
    return models.FromDBUser(dbUser), nil
}
```

### 5. Database Layer (`internal/db/`)

**Purpose**: Generated database code and raw SQL operations

**What it contains:**
- sqlc generated Go code
- Database models (structs)
- SQL query functions
- Database connection management

**Note**: This layer is auto-generated by sqlc and should not be modified manually.

### 6. Models Layer (`internal/models/`)

**Purpose**: Domain models and data transfer objects

**What it contains:**
- Domain entities (business objects)
- Request/Response DTOs
- Data validation tags
- Conversion functions between layers
- Business-related structs

**Example:**
```go
type User struct {
    ID        int32     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}
```

## Data Flow

### Request Flow (Top to Bottom)
1. **HTTP Request** → Handler receives HTTP request
2. **Handler** → Validates input, calls Service
3. **Service** → Applies business logic, calls Repository
4. **Repository** → Executes database operations via sqlc
5. **Database** → Returns raw data

### Response Flow (Bottom to Top)
1. **Database** → Returns raw data to Repository
2. **Repository** → Converts to domain models, returns to Service
3. **Service** → Applies business transformations, returns to Handler
4. **Handler** → Formats HTTP response, sends to client

## Error Handling Strategy

Each layer has its own error types and responsibilities:

### Repository Errors
```go
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrUserEmailExists = errors.New("user with this email already exists")
)
```

### Service Errors
```go
var (
    ErrInvalidUserID     = errors.New("invalid user ID")
    ErrInvalidUserName   = errors.New("invalid user name")
    ErrUserNameTooLong   = errors.New("user name is too long")
)
```

### Handler Error Responses
- Repository errors → HTTP status codes
- Service errors → HTTP status codes with appropriate messages
- Validation errors → 400 Bad Request
- Not found errors → 404 Not Found
- Business logic errors → 409 Conflict or 422 Unprocessable Entity

## Benefits of This Architecture

### 1. Separation of Concerns
- Each layer has a single responsibility
- Changes in one layer don't affect others
- Easy to understand and maintain

### 2. Testability
- Each layer can be tested independently
- Easy to mock dependencies
- Unit tests are focused and fast

### 3. Flexibility
- Can swap implementations (e.g., different databases)
- Can add new features without affecting existing code
- Can change API format without changing business logic

### 4. Scalability
- Clear boundaries for team development
- Easy to add new endpoints
- Business logic is reusable across different interfaces

## Best Practices

### 1. Dependency Direction
- Dependencies should point inward (toward business logic)
- Outer layers depend on inner layers, not vice versa
- Use interfaces to decouple layers

### 2. Error Handling
- Each layer should handle errors appropriate to its level
- Convert errors between layers (repository → service → handler)
- Provide meaningful error messages to users

### 3. Context Usage
- Pass context through all layers for cancellation and timeouts
- Use context for request-scoped data
- Don't store business data in context

### 4. Testing Strategy
- Unit tests for each layer
- Mock dependencies for isolated testing
- Integration tests for database operations
- E2E tests for complete request flow

## Example Implementation Checklist

When adding a new feature (e.g., "Products"):

1. **SQL Layer**: Create SQL schema and queries
2. **Generate Code**: Run `sqlc generate`
3. **Models**: Create domain models and DTOs
4. **Repository**: Implement data access interface
5. **Service**: Implement business logic
6. **Handlers**: Create HTTP endpoints
7. **Routes**: Wire up in main.go
8. **Tests**: Add tests for each layer
9. **Documentation**: Update Swagger annotations

This architecture ensures that your API remains maintainable and scalable as it grows in complexity.