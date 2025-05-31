# Swagger Documentation Guide

This guide explains how to add Swagger/OpenAPI documentation to your Mini Lab API endpoints.

## Overview

The Mini Lab API uses `swaggo/swag` to automatically generate OpenAPI/Swagger documentation from Go comments in your code.

## Adding Documentation to New Endpoints

### Step 1: Add Swagger Comments

Use special comments above your handler functions to document the API:

```go
// @Summary Short description of the endpoint
// @Description Detailed description of what the endpoint does
// @Tags tag-name
// @Accept json (for endpoints that accept JSON)
// @Produce json (for endpoints that return JSON)
// @Param name path/query/body type required "description" format(optional)
// @Success 200 {object} ResponseType "Success description"
// @Failure 400 {object} ErrorType "Error description"
// @Router /path/to/endpoint [method]
func YourHandler(c *gin.Context) {
    // Handler implementation
}
```

### Step 2: Document Request/Response Models

Add example tags to your structs:

```go
type User struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"John Doe"`
    Email string `json:"email" example:"john@example.com"`
}

type CreateUserRequest struct {
    Name  string `json:"name" binding:"required" example:"John Doe"`
    Email string `json:"email" binding:"required,email" example:"john@example.com"`
}
```

### Step 3: Regenerate Documentation

After adding new endpoints or modifying existing ones:

```bash
make swagger-generate
```

## Common Swagger Annotations

### Basic Endpoint Documentation

```go
// @Summary Get user profile
// @Description Retrieve user profile information by user ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User "User profile"
// @Failure 404 {object} ApiResponse "User not found"
// @Router /api/v1/users/{id} [get]
```

### POST Endpoint with Body

```go
// @Summary Create new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} ApiResponse{data=User} "User created"
// @Failure 400 {object} ApiResponse "Invalid input"
// @Router /api/v1/users [post]
```

### Query Parameters

```go
// @Summary List users
// @Description Get list of users with optional filtering
// @Tags users
// @Produce json
// @Param limit query int false "Limit results" minimum(1) maximum(100)
// @Param offset query int false "Offset for pagination" minimum(0)
// @Param search query string false "Search term"
// @Success 200 {object} ApiResponse{data=[]User} "List of users"
// @Router /api/v1/users [get]
```

### Headers

```go
// @Summary Protected endpoint
// @Description Access protected resource
// @Tags auth
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} ApiResponse "Success"
// @Failure 401 {object} ApiResponse "Unauthorized"
// @Router /api/v1/protected [get]
```

## Main Application Documentation

At the top of your main.go file, add general API information:

```go
// @title Your API Title
// @version 1.0
// @description Your API description
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url http://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
    // Your main function
}
```

## Response Types

### Standard API Response

```go
type ApiResponse struct {
    Data    interface{} `json:"data,omitempty"`
    Status  string      `json:"status" example:"success"`
    Message string      `json:"message" example:"Operation completed successfully"`
    Error   string      `json:"error,omitempty" example:"Error message"`
}
```

### Using Generic Response Types

For endpoints returning different data types:

```go
// Success with User data
// @Success 200 {object} ApiResponse{data=User}

// Success with array of Users
// @Success 200 {object} ApiResponse{data=[]User}

// Success with custom data
// @Success 200 {object} ApiResponse{data=map[string]interface{}}
```

## Testing Your Documentation

1. **Start the server**: `make run`
2. **Visit Swagger UI**: `http://localhost:8080/swagger/index.html`
3. **Test endpoints**: Use the interactive interface to test your API

## Best Practices

1. **Be Descriptive**: Use clear, descriptive summaries and descriptions
2. **Include Examples**: Add realistic examples to your models
3. **Document All Parameters**: Include all required and optional parameters
4. **Use Proper HTTP Status Codes**: Document all possible response codes
5. **Group Related Endpoints**: Use consistent tags to group related endpoints
6. **Keep Documentation Updated**: Regenerate docs whenever you change the API

## Troubleshooting

### Common Issues

1. **Documentation not updating**: Run `make swagger-generate` after changes
2. **Missing endpoints**: Ensure your handler functions have proper comments
3. **Import errors**: Make sure to import the generated docs package in main.go

### Validation

Check that your generated documentation is valid:

```bash
# View the generated JSON
cat docs/swagger.json | jq .

# Validate the OpenAPI spec online
# Copy contents of docs/swagger.json to https://editor.swagger.io/
```

## Example Complete Handler

```go
// @Summary Create user account
// @Description Create a new user account with email verification
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User registration data"
// @Success 201 {object} ApiResponse{data=User} "User created successfully"
// @Failure 400 {object} ApiResponse "Invalid input data"
// @Failure 409 {object} ApiResponse "Email already exists"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid request body",
            "status":  "error",
            "message": err.Error(),
        })
        return
    }
    
    // Implementation...
}
``` 