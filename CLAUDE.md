# ATHYLPS Backend API

## Project Overview
ATHYLPS is a Go-based backend API for a mobile application. This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) and uses modern Go practices and tooling.

## Tech Stack
- **Language**: Go 1.25.0
- **Router**: [Chi](https://github.com/go-chi/chi) - Lightweight, idiomatic HTTP router
- **Database**: PostgreSQL with [pgxpool](https://github.com/jackc/pgx) - High-performance PostgreSQL driver
- **Migrations**: [Goose](https://github.com/pressly/goose) - Database migration tool
- **Formatter**: [gofumpt](https://github.com/mvdan/gofumpt) - Stricter gofmt

## Project Structure

The project follows the Standard Go Project Layout:

```
athylps-backend-go/
├── cmd/
│   └── athylps/           # Main application entry point
│       └── main.go
├── internal/              # Private application code
│   └── .gitkeep
├── pkg/                   # Public library code
│   └── .gitkeep
├── .github/               # GitHub workflows and configurations
├── go.mod                 # Go module definition
├── Makefile              # Build and development tasks
└── .gitignore
```

### Directory Descriptions

- **`cmd/athylps/`**: Main application entry point. Contains `main.go` which initializes and runs the API server.
- **`internal/`**: Private application code that cannot be imported by other projects. This is where most of the business logic, handlers, services, and repositories should live.
- **`pkg/`**: Public library code that can be imported by external applications. Use sparingly and only for truly reusable components.
- **`.github/`**: CI/CD workflows and GitHub-specific configurations.

## Development

### Prerequisites
- Go 1.25.0 or later
- PostgreSQL database
- Make (optional, but recommended)

### Available Make Commands

```bash
make run    # Run the application
make test   # Run all tests
make vet    # Run go vet for static analysis
```

### Running the Application

```bash
# Using Make
make run

# Direct Go command
go run cmd/athylps/main.go
```

### Testing

```bash
# Using Make
make test

# Direct Go command
go test ./...
```

## Code Style and Formatting

This project uses **gofumpt** for code formatting, which is a stricter version of `gofmt`.

```bash
# Install gofumpt
go install mvdan.cc/gofumpt@latest

# Format all code
gofumpt -l -w .
```

## Database

### PostgreSQL with pgxpool
The project uses `pgxpool` for database connection pooling, which provides:
- High performance
- Built-in connection pooling
- Better PostgreSQL feature support than database/sql

### Migrations with Goose
Database migrations are managed using Goose. Migration files should be placed in a migrations directory (typically `internal/db/migrations/` or similar).

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Example migration commands
goose -dir ./migrations postgres "connection-string" up
goose -dir ./migrations postgres "connection-string" down
```

## Recommended Project Structure (As You Build)

```
athylps-backend-go/
├── cmd/
│   └── athylps/
│       └── main.go
├── internal/
│   ├── api/              # HTTP handlers and routing
│   │   ├── handlers/     # HTTP handlers
│   │   ├── middleware/   # Custom middleware
│   │   └── router/       # Chi router setup
│   ├── config/           # Configuration management
│   ├── db/
│   │   ├── migrations/   # Goose migration files
│   │   └── queries/      # SQL queries or repository interfaces
│   ├── models/           # Domain models
│   ├── repository/       # Database access layer
│   └── service/          # Business logic layer
├── pkg/
│   └── utils/            # Shared utilities (if needed)
├── .env.example          # Example environment variables
├── go.mod
├── go.sum
└── Makefile
```

## Environment Variables

Create a `.env` file in the project root (already in `.gitignore`). Common variables:

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=athylps
DB_SSLMODE=disable

# Optional: Connection pool settings
DB_MAX_CONNS=25
DB_MIN_CONNS=5
DB_MAX_CONN_LIFETIME=1h
DB_MAX_CONN_IDLE_TIME=30m
```

## API Design Principles

When building the API:
1. Use RESTful conventions where appropriate
2. Structure routes logically with Chi router groups
3. Apply middleware consistently (logging, CORS, auth, etc.)
4. Return proper HTTP status codes
5. Use JSON for request/response bodies
6. Implement proper error handling and validation
7. Consider versioning API endpoints (e.g., `/api/v1/...`)

## Testing Strategy

- **Unit Tests**: Test individual functions and methods in isolation
- **Integration Tests**: Test database operations and handler logic
- **Table-Driven Tests**: Use Go's table-driven test pattern
- **Test Coverage**: Aim for meaningful coverage, not just high percentages

## CI/CD

GitHub Actions workflows are configured in `.github/`. These should include:
- Running tests
- Code linting and formatting checks
- Building the application
- Deployment steps (if applicable)

## Security Considerations

- Never commit `.env` files or secrets
- Use environment variables for sensitive configuration
- Implement proper authentication and authorization
- Validate and sanitize all user inputs
- Use prepared statements/parameterized queries (pgx handles this)
- Keep dependencies up to date
- Use HTTPS in production

## Deployment

(Add specific deployment instructions as your infrastructure is set up)

## Contributing

1. Follow the Standard Go Project Layout
2. Use `gofumpt` for formatting
3. Write tests for new features
4. Run `make vet` and `make test` before committing
5. Keep commits atomic and write clear commit messages

## Resources

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Chi Router Documentation](https://go-chi.io)
- [pgx Documentation](https://github.com/jackc/pgx)
- [Goose Migrations](https://github.com/pressly/goose)
- [Effective Go](https://golang.org/doc/effective_go)
