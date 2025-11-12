FROM golang:1.25-alpine AS build-stage

WORKDIR /build

# Copy go mod files from parent directory
COPY ../go.mod ../go.sum ./
RUN go mod download

# Copy the entire project from parent directory
COPY ../ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/athylps ./cmd/athylps/main.go
# Build the migrate application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/migrate ./cmd/migrate/main.go

# Deploy stage
FROM alpine:latest AS production-stage

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from build stage
COPY --from=build-stage /app/athylps .

# Copy the migrate binary from build state
COPY --from=build-stage /app/migrate .

# Copy migrations if they exist (optional - remove this line if migrations don't exist yet)
COPY --from=build-stage /build/migrations/ ./migrations/

# Expose port (adjust if needed)
EXPOSE 8080

# Run as non-root user
RUN adduser -D -u 1000 appuser && chown -R appuser:appuser /app
USER appuser

# Run migrations
RUN "./migrate" "up"

CMD ["./athylps"]