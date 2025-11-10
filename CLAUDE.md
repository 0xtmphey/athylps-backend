# ATHYLPS Backend API

## Project Overview
ATHYLPS is a Go-based backend API for a mobile application. This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) and uses modern Go practices and tooling.

## Tech Stack
- **Language**: Go 1.25.0
- **Router**: [Chi](https://github.com/go-chi/chi) - Lightweight, idiomatic HTTP router
- **Database**: PostgreSQL with [pgxpool](https://github.com/jackc/pgx) - High-performance PostgreSQL driver
- **Migrations**: [Goose](https://github.com/pressly/goose) - Database migration tool
- **Formatter**: [gofumpt](https://github.com/mvdan/gofumpt) - Stricter gofmt
- **API Spec**: OpenAPI 3.0 with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) for code generation
- **Container**: Docker & Docker Compose for local development

## Project Structure

The project follows the Standard Go Project Layout:

```
athylps-backend-go/
├── api/                   # OpenAPI specification and codegen config
│   ├── openapi.yaml      # API specification
│   └── oapi-codegen.yaml # Code generation configuration
├── build/                 # Build files and scripts (Dockerfile)
├── cmd/
│   └── athylps/          # Main application entry point
│       └── main.go
├── internal/             # Private application code
│   ├── api/             # Generated API models from OpenAPI spec
│   ├── app/             # Application initialization (dependencies, HTTP server)
│   ├── config/          # Configuration management (env vars)
│   ├── handlers/        # HTTP request handlers
│   ├── services/        # Reusable services (external APIs, shared logic)
│   ├── usecases/        # Business logic use cases
│   └── repositories/    # Database access layer
├── migrations/           # Database migrations (Goose)
├── pkg/                  # Public library code
├── .github/              # GitHub workflows and configurations
├── docker-compose.yml    # Local database setup
├── go.mod               # Go module definition
├── Makefile             # Build and development tasks
└── .env.example         # Example environment variables
```

### Directory Descriptions

- **`api/`**: Contains OpenAPI 3.0 specification for API documentation and code generation. Used to auto-generate request/response models.
- **`build/`**: Build-related files including Dockerfile for containerization.
- **`cmd/athylps/`**: Main application entry point. Contains `main.go` which initializes and runs the API server.
- **`internal/api/`**: Auto-generated API models from OpenAPI specification.
- **`internal/app/`**: Application bootstrapping - dependency injection, HTTP server setup, router configuration.
- **`internal/config/`**: Configuration management, reading environment variables.
- **`internal/handlers/`**: HTTP request handlers - parse requests, validate input, send responses.
- **`internal/services/`**: Reusable services not tied to specific use cases (e.g., EmailService, TelegramApiService).
- **`internal/usecases/`**: Business logic use cases - encapsulate specific user scenarios.
- **`internal/repositories/`**: Database access layer - functions for working with database tables.
- **`migrations/`**: Database migration files managed by Goose.
- **`pkg/`**: Public library code that can be imported by external applications. Use sparingly and only for truly reusable components.
- **`.github/`**: CI/CD workflows and GitHub-specific configurations.

## Development

### Prerequisites
- Go 1.25.0 or later
- Docker & Docker Compose
- Make

### Initial Setup

1. **Install dependencies and tools:**
```bash
make setup
```

This installs:
- Go module dependencies
- goose (migration tool)
- gofumpt (formatter)
- oapi-codegen (OpenAPI code generator)

2. **Create environment file:**
```bash
cp .env.example .env
```

3. **Start local database:**
```bash
docker compose up -d
```

This starts a PostgreSQL container. You can monitor it in Docker Desktop.

4. **Generate API models from OpenAPI spec:**
```bash
make generate-api-models
```

### Available Make Commands

```bash
make help                   # Display all available commands
make setup                  # Install all required tools
make run                    # Run the application
make build                  # Build the application binary
make test                   # Run all tests
make test-cover            # Run tests with coverage report
make lint                  # Run all linting checks (format + vet)
make fmt                   # Format code with gofumpt
make vet                   # Run go vet for static analysis
make generate-api-models   # Generate API models from OpenAPI spec
make migrations-up         # Run database migrations
make migrations-down       # Rollback last migration
make migrations-status     # Check migration status
make migrations-create     # Create new migration (use NAME=migration_name)
make ci                    # Run all CI checks
make clean                 # Clean build artifacts
```

### Running the Application

```bash
make run
```

The application runs on port `8080` by default: `http://localhost:8080`

**Note:** Hot-reload is not currently supported. Restart the application after making changes.

### Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-cover
```

## Code Style and Formatting

This project uses **gofumpt** for code formatting, which is a stricter version of `gofmt`.

```bash
# Format all code
make fmt

# Check if code is properly formatted
make fmt-check
```

Run linting before committing:
```bash
make lint
```

## Database

### PostgreSQL with pgxpool
The project uses `pgxpool` for database connection pooling, which provides:
- High performance
- Built-in connection pooling
- Better PostgreSQL feature support than database/sql

### Migrations with Goose
Database migrations are managed using Goose. Migration files are in the `migrations/` directory.

```bash
# Create a new migration
make migrations-create NAME=create_users_table

# Run migrations
make migrations-up

# Rollback last migration
make migrations-down

# Check migration status
make migrations-status
```

**Note:** Migration commands require `DATABASE_URL` environment variable to be set.

## OpenAPI Code Generation

The project uses OpenAPI 3.0 specification for API documentation and automatic code generation.

1. **Define your API** in `api/openapi.yaml`
2. **Configure code generation** in `api/oapi-codegen.yaml`
3. **Generate models:**
```bash
make generate-api-models
```

Generated models will be placed in `internal/api/` and can be used in handlers for request/response types.

## Environment Variables

Create a `.env` file in the project root (already in `.gitignore`). See `.env.example` for reference:

```env
# Server
PORT=8080
ENV=development

# Database
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=athylps
DB_HOST=localhost
DB_PORT=5432

# RevenueCat
RC_BEARER=

# Telegram Bot
BOT_TOKEN=
NOTIFY_CHAT_ID=
```

## Architecture & Request Lifecycle

The project follows a **layered architecture** with clear separation of concerns:

### Request Flow
```
User → Handler → Usecase → Repository/Service → Database/External API
                     ↓
User ← Handler ← Usecase
```

1. **User** sends HTTP request
2. **Handler** parses request, validates input
3. **Handler** calls **Usecase** to execute business logic
4. **Usecase** uses **Repositories** (database) and **Services** (external APIs)
5. **Usecase** returns result to **Handler**
6. **Handler** formats and sends HTTP response

### Handlers
- Located in `internal/handlers/`
- Parse incoming requests
- Perform basic validation
- Send responses to clients
- One handler per endpoint
- Accept dependencies via closure

Example:
```go
func CreateUserHandler(logger *zap.Logger, usecase userUsecase) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse, validate, call usecase, respond
    }
}
```

### Usecases
- Located in `internal/usecases/`
- Encapsulate business logic scenarios
- Reusable across different contexts (HTTP, cron jobs, etc.)
- Implemented as structs with dependencies
- One public method to perform the use case

Example:
```go
type CreateUserUsecase struct {
    logger *zap.Logger
    repo   userRepository
}

func NewCreateUserUsecase(logger *zap.Logger, repo userRepository) *CreateUserUsecase {
    return &CreateUserUsecase{logger: logger, repo: repo}
}

func (u *CreateUserUsecase) Perform(userName string) (*User, error) {
    // Business logic here
}
```

### Repositories
- Located in `internal/repositories/`
- Encapsulate database access
- Usually one repository per entity/table
- Implemented as structs with database connection

Example:
```go
type UserRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(id string) (*User, error) {
    // Database query here
}
```

### Services
- Located in `internal/services/`
- Reusable functionality not tied to specific use cases
- External API integrations
- Shared business logic

Example:
```go
type TelegramService struct {
    logger *zap.Logger
    client *tg.Api
}

func NewTelegramService(logger *zap.Logger, client *tg.Api) *TelegramService {
    return &TelegramService{logger: logger, client: client}
}

func (s *TelegramService) SendMessage(msg string) error {
    // Send message via Telegram API
}
```

## Interfaces in Go

Go uses **implicit interface implementation**. Types automatically satisfy an interface if they implement all its methods.

**Key principle:** Define interfaces on the consumer side, not the provider side.

Example:
```go
// In handler file - define what you need
type userUsecase interface {
    CreateUser(name string) (*User, error)
}

func CreateUserHandler(u userUsecase) http.HandlerFunc {
    // Handler code
}

// In usecase file - just implement the method
type CreateUserUsecase struct {}

func (u *CreateUserUsecase) CreateUser(name string) (*User, error) {
    // Implementation
}
// Automatically satisfies userUsecase interface!
```

This approach:
- Keeps interfaces small and focused
- Reduces coupling
- Makes testing easier (easy to mock dependencies)

## Adding a New Feature

Example: Add a `GET /random` endpoint that returns a random number.

1. **Define the endpoint in OpenAPI spec** (`api/openapi.yaml`)

2. **Generate models:**
```bash
make generate-api-models
```

3. **Create handler** (`internal/handlers/get_random.go`):
```go
type randomUsecase interface {
    GetRandomInt() int
}

func GetRandomHandler(u randomUsecase) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        random := u.GetRandomInt()
        resp := api.GetRandomResponse{Random: random}
        json.NewEncoder(w).Encode(resp)
    }
}
```

4. **Create usecase** (`internal/usecases/get_random_usecase.go`):
```go
type GetRandomUsecase struct {}

func NewGetRandomUsecase() *GetRandomUsecase {
    return &GetRandomUsecase{}
}

func (u *GetRandomUsecase) GetRandomInt() int {
    return 42 // or use math/rand
}
```

5. **Wire it up in app** (`internal/app/app.go`):
```go
randomUsecase := usecases.NewGetRandomUsecase()
r.Get("/random", handlers.GetRandomHandler(randomUsecase))
```

6. **Test:**
```bash
curl http://localhost:8080/random
```

## API Design Principles

When building the API:
1. Define endpoints in OpenAPI spec first
2. Generate models with `make generate-api-models`
3. Use RESTful conventions where appropriate
4. Structure routes logically with Chi router groups
5. Apply middleware consistently (logging, CORS, auth, etc.)
6. Return proper HTTP status codes
7. Use JSON for request/response bodies
8. Implement proper error handling and validation
9. Keep handlers thin - business logic goes in usecases
10. Define interfaces on the consumer side (in handlers/usecases)

## Testing Strategy

- **Unit Tests**: Test individual functions and methods in isolation
- **Integration Tests**: Test database operations and handler logic
- **Table-Driven Tests**: Use Go's table-driven test pattern
- **Test Coverage**: Aim for meaningful coverage, not just high percentages

## Deployment

### Production URL
```
https://athylps-api.tmphey.dev
```

### Deploying to Production

The project uses GitHub Actions for CI/CD. To deploy to production:

1. **Create a new tag:**
```bash
git tag -a v1.0.0 -m "Release notes here"
```

2. **Push to main with tags:**
```bash
git push origin main --tags
```

This triggers the GitHub Actions pipeline which:
- Runs tests and linting
- Builds the Docker image
- Deploys to production

### CI/CD Pipeline

GitHub Actions workflows are configured in `.github/`. The pipeline includes:
- Running tests
- Code linting and formatting checks
- Building the Docker image
- Deploying to production on tag push

## Security Considerations

- Never commit `.env` files or secrets
- Use environment variables for sensitive configuration
- Implement proper authentication and authorization
- Validate and sanitize all user inputs
- Use prepared statements/parameterized queries (pgx handles this automatically)
- Keep dependencies up to date (`go mod tidy`)
- Use HTTPS in production

## Development Workflow

1. **Start a new feature:**
   - Pull latest changes from main
   - Create a feature branch (optional)
   - Define API endpoint in `api/openapi.yaml`
   - Run `make generate-api-models`

2. **Implement the feature:**
   - Create handler in `internal/handlers/`
   - Create usecase in `internal/usecases/`
   - Create repository/service if needed
   - Wire dependencies in `internal/app/app.go`
   - Format code: `make fmt`

3. **Test your changes:**
   - Write tests
   - Run `make test`
   - Run `make lint`
   - Test manually with the running application

4. **Before committing:**
   - Run `make ci` to run all checks
   - Ensure all tests pass
   - Ensure code is properly formatted

5. **Deploy to production:**
   - Merge to main
   - Create and push a tag: `git tag -a v1.0.0 -m "Release"` && `git push origin main --tags`

## Common Patterns

### File Organization
Files are organized by **responsibility zone** rather than by feature:
- `internal/handlers/` - all handlers
- `internal/usecases/` - all usecases
- `internal/repositories/` - all repositories

This keeps the structure simple and leverages Go's package naming conventions.

### Dependency Injection
Dependencies are manually wired in `internal/app/app.go`:
```go
// Create dependencies
db := connectDB()
logger := setupLogger()

// Create repositories
userRepo := repositories.NewUserRepository(db)

// Create services
tgService := services.NewTelegramService(logger)

// Create usecases
createUserUsecase := usecases.NewCreateUserUsecase(logger, userRepo)

// Create handlers and register routes
r.Post("/users", handlers.CreateUserHandler(logger, createUserUsecase))
```

### Error Handling
- Return errors from functions, don't panic
- Wrap errors with context: `fmt.Errorf("failed to create user: %w", err)`
- Log errors at the appropriate level
- Return appropriate HTTP status codes in handlers

### Logging
The project uses structured logging (likely `zap` or similar). Always include context:
```go
logger.Info("Creating user", zap.String("username", userName))
logger.Error("Failed to create user", zap.Error(err))
```

## Important Notes

1. **No hot-reload**: Restart the application after making changes
2. **Organization by responsibility**: Files organized by type (handlers, usecases, etc.), not by feature
3. **Interfaces on consumer side**: Define interfaces where you use them, not where you implement them
4. **OpenAPI-first**: Always define API endpoints in the spec before implementing
5. **Manual dependency injection**: All wiring happens in `internal/app/app.go`

## Troubleshooting

### Database connection issues
- Ensure Docker container is running: `docker ps`
- Check environment variables in `.env`
- Verify database credentials match Docker Compose settings

### Migration issues
- Ensure `DATABASE_URL` environment variable is set
- Check migration status: `make migrations-status`
- Review migration files in `migrations/` directory

### Code generation issues
- Ensure `oapi-codegen` is installed: `make setup`
- Check `api/openapi.yaml` for syntax errors
- Review `api/oapi-codegen.yaml` configuration

### Build/Run issues
- Clean build artifacts: `make clean`
- Update dependencies: `go mod tidy`
- Verify all tools installed: `make tools`

## Resources

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Chi Router Documentation](https://go-chi.io)
- [pgx Documentation](https://github.com/jackc/pgx)
- [Goose Migrations](https://github.com/pressly/goose)
- [OpenAPI Specification](https://swagger.io/specification/)
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
- [Effective Go](https://golang.org/doc/effective_go)
- [Using interfaces in Go the right way](https://medium.com/@mbinjamil/using-interfaces-in-go-the-right-way-99384bc69d39)
