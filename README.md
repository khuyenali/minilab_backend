# Mini Lab API

A simple RESTful API server built with Go, Gin, sqlc, and golang-migrate using clean architecture principles.

## Features

- **Gin** - Fast HTTP web framework
- **sqlc** - Generate type-safe Go code from SQL
- **golang-migrate** - Database migrations
- **PostgreSQL** - Database
- **Adminer** - Web-based database management
- **Swagger** - API documentation
- **Clean Architecture** - Layered design with proper separation of concerns
- RESTful API design
- Environment-based configuration
- CORS support
- JSON logging

## Project Structure

```
mini_lab_backend/
├── cmd/api/                 # Application entry point
│   ├── main.go             # HTTP server setup and dependency injection
│   └── main_test.go        # HTTP layer tests
├── internal/                # Private application code
│   ├── config/             # Configuration management
│   │   ├── config.go       # Environment configuration
│   │   └── database.go     # Database connection setup
│   ├── db/                 # Generated database code (sqlc)
│   │   ├── db.go           # Database interface
│   │   ├── models.go       # Database models
│   │   ├── querier.go      # Query interface
│   │   └── users.sql.go    # Generated user queries
│   ├── handlers/           # HTTP request handlers
│   │   ├── handler.go      # Handler struct and shared types
│   │   ├── ping.go         # Health check handlers
│   │   └── users.go        # User-specific handlers
│   ├── models/             # Domain models and DTOs
│   │   └── user.go         # User domain model
│   ├── repository/         # Data access layer
│   │   ├── errors.go       # Repository errors
│   │   └── user_repository.go # User data access
│   └── service/            # Business logic layer
│       ├── errors.go       # Service errors
│       └── user_service.go # User business logic
├── migrations/             # Database migrations
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── sql/                    # SQL files for sqlc
│   ├── queries/           # SQL queries
│   │   └── users.sql      # User queries
│   └── schema/            # Database schema
│       └── 001_users.sql  # User table schema
├── docs/                   # Generated Swagger documentation
├── sqlc.yaml              # sqlc configuration
├── docker-compose.yml     # PostgreSQL setup
├── Makefile              # Build automation
├── ARCHITECTURE.md       # Architecture documentation
├── SWAGGER_GUIDE.md      # Swagger documentation guide
└── README.md             # This file
```

## Architecture

This API follows a clean, layered architecture:

- **HTTP Layer** (`cmd/api/`): Server setup and routing
- **Handler Layer** (`internal/handlers/`): HTTP request/response handling
- **Service Layer** (`internal/service/`): Business logic and validation
- **Repository Layer** (`internal/repository/`): Data access abstraction
- **Database Layer** (`internal/db/`): Generated database code

For detailed architecture information, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- sqlc (for code generation)
- golang-migrate (for migrations)
- swag (for Swagger documentation)

### Quick Setup

1. **Clone and install dependencies:**
```bash
git clone <repository-url>
cd mini_lab_backend
go mod download
```

2. **Install development tools:**
```bash
make install-tools
```

3. **Setup environment:**
```bash
cp env.example .env
# Edit .env with your database configuration
```

4. **Start PostgreSQL:**
```bash
# Using Docker (recommended)
docker compose up -d

# Or use your own PostgreSQL instance
```

5. **Run migrations:**
```bash
make migrate-up
```

6. **Generate code and documentation:**
```bash
make sqlc-generate
make swagger-generate
```

7. **Start the server:**
```bash
make run
```

The API will be available at:
- **API**: `http://localhost:8080`
- **Swagger Documentation**: `http://localhost:8080/swagger/index.html`
- **Database Management (Adminer)**: `http://localhost:8081`

### Complete Setup (One Command)

```bash
make setup
```

This will install all tools, generate code, and prepare the project for development.

## API Endpoints

### Health & Documentation
- `GET /health` - Health check endpoint
- `GET /swagger/index.html` - Interactive Swagger UI
- `GET /swagger/doc.json` - OpenAPI JSON specification

### Ping
- `GET /api/v1/ping` - Simple connectivity test

### Users (Database-backed)
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Example Requests

Create a user:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

Get all users:
```bash
curl http://localhost:8080/api/v1/users
```

Update a user:
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john.updated@example.com"}'
```

## Development

### Quick Development Setup

For a completely automated development setup:

```bash
make dev-setup
```

This command will:
- Install all development tools (including Air)
- Start PostgreSQL with Docker
- Run database migrations
- Prepare the project for development

### Hot Reload Development

For the best development experience, use **Air** for hot reloading:

```bash
make dev
```

**Air** will:
- Automatically rebuild and restart your server when Go files change
- Watch for changes in all `.go` files
- Exclude test files and temporary directories
- Show colored output for better debugging
- Log build errors to `build-errors.log`

**Air Configuration** (`.air.toml`):
- **Build command**: `go build -o ./tmp/main cmd/api/main.go`
- **Watched extensions**: `.go`, `.html`, `.tpl`, `.tmpl`
- **Excluded directories**: `tmp/`, `vendor/`, `migrations/`, `docs/`, `bin/`
- **Build delay**: 1000ms (to avoid excessive rebuilds)

### Common Commands

```bash
# Development
make run                    # Run the application (standard mode)
make dev                    # Run with hot reload using Air
make dev-setup              # Complete development environment setup
make build                  # Build the application

# Database & Services
make docker-up              # Start PostgreSQL and Adminer with Docker Compose
make docker-down            # Stop Docker Compose services
make docker-logs            # Show Docker Compose logs
make migrate-up             # Run database migrations
make migrate-down           # Rollback migrations
make migrate-create NAME=migration_name # Create new migration

# Code Generation
make sqlc-generate          # Generate Go code from SQL
make swagger-generate       # Generate/update API documentation

# Testing & Quality
make test                   # Run tests
make test-cover            # Run tests with coverage
make fmt                   # Format code
make lint                  # Lint code

# Utilities
make clean                 # Clean build artifacts
make help                  # Show all available commands
```

### Development Workflow

1. **Start development environment:**
```bash
make dev-setup  # One-time setup
make dev        # Start hot-reload server
```

2. **Make changes** to any `.go` file - Air will automatically:
   - Detect the change
   - Rebuild the application
   - Restart the server
   - Show you any build errors

3. **Test your changes** at `http://localhost:8080`

4. **View logs** in the terminal or check `build-errors.log` for build issues

### Development Tips

- **Air rebuilds are fast** - typically under 1-2 seconds
- **Database changes** require manual migration runs (`make migrate-up`)
- **SQL changes** require code regeneration (`make sqlc-generate`)
- **API doc changes** require doc regeneration (`make swagger-generate`)
- **Environment changes** require server restart (Ctrl+C and `make dev`)

### Database Management with Adminer

**Adminer** provides a web-based interface to manage your PostgreSQL database:

1. **Start services:**
```bash
make docker-up
```

2. **Access Adminer:** Visit `http://localhost:8081`

3. **Login credentials:**
   - **System**: `PostgreSQL`
   - **Server**: `postgres` (container name)
   - **Username**: `postgres`
   - **Password**: `password`
   - **Database**: `minilab`

**Adminer Features:**
- Browse tables and data
- Execute SQL queries
- View database schema
- Import/export data
- Monitor database performance
- Create and modify tables

**Quick Actions:**
```bash
make docker-logs    # View database and Adminer logs
make docker-down    # Stop all services
make docker-up      # Restart all services
```

### Troubleshooting Air

If Air isn't working as expected:

1. **Check Air installation:**
```bash
air -v
```

2. **Clean temporary files:**
```bash
make clean
```

3. **Restart development server:**
```bash
# Stop with Ctrl+C, then restart
make dev
```

4. **Check Air configuration:**
```bash
cat .air.toml
```

### Database Operations

**Create a new migration:**
```bash
make migrate-create NAME=add_products_table
```

**Add SQL queries** in `sql/queries/` and regenerate code:
```bash
make sqlc-generate
```

**Update API documentation:**
```bash
make swagger-generate
```

### Adding New Features

When adding a new entity (e.g., "Products"):

1. **Create SQL schema** in `sql/schema/002_products.sql`
2. **Create migration** with `make migrate-create NAME=create_products_table`
3. **Add SQL queries** in `sql/queries/products.sql`
4. **Generate code** with `make sqlc-generate`
5. **Create domain model** in `internal/models/product.go`
6. **Implement repository** in `internal/repository/product_repository.go`
7. **Implement service** in `internal/service/product_service.go`
8. **Create handlers** in `internal/handlers/products.go`
9. **Add routes** in `cmd/api/main.go`
10. **Update documentation** with `make swagger-generate`

## API Documentation

The API is fully documented using OpenAPI/Swagger:

- **Interactive UI**: Visit `http://localhost:8080/swagger/index.html`
- **JSON Spec**: Available at `http://localhost:8080/swagger/doc.json`
- **Local Files**: Generated in the `docs/` directory

For adding documentation to new endpoints, see [SWAGGER_GUIDE.md](SWAGGER_GUIDE.md).

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Application environment | `development` |
| `PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `minilab` |
| `DB_SSLMODE` | Database SSL mode | `disable` |
| `DATABASE_URL` | Full database URL | - |

## Error Handling

The API uses consistent error responses:

```json
{
    "error": "Error type",
    "status": "error",
    "message": "Detailed error message"
}
```

**HTTP Status Codes:**
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `409` - Conflict (duplicate email, etc.)
- `500` - Internal Server Error

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Run specific package tests
go test ./internal/service/...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the architecture guidelines
4. Add tests for new functionality
5. Update documentation
6. Submit a pull request

## License

This project is licensed under the MIT License. 