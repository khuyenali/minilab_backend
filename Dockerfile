# Stage 1: Build the Go application and download migrate CLI
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install curl for downloading, git is not strictly needed for this download method
RUN apk update && apk add --no-cache curl tar

# Download and install golang-migrate CLI
ENV MIGRATE_VERSION v4.17.1
RUN curl -fsSL "https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz" -o migrate.tar.gz && \
    echo "--- migrate.tar.gz details:" && \
    ls -lh migrate.tar.gz && \
    echo "--- migrate.tar.gz contents:" && \
    tar -tvzf migrate.tar.gz && \
    echo "--- Extracting migrate.tar.gz:" && \
    tar -xzf migrate.tar.gz && \
    echo "--- Contents of /app after extraction:" && \
    ls -la && \
    echo "--- Moving migrate binary:" && \
    mv migrate /usr/local/bin/migrate && \
    echo "--- migrate binary moved successfully." && \
    rm migrate.tar.gz

# Copy go.mod and go.sum to leverage Docker cache for dependency downloads
COPY go.mod go.sum ./

# Download dependencies based on go.mod and go.sum. This is primarily for caching.
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Now that all source code is present, tidy and vendor for highest accuracy
RUN go mod tidy
RUN go mod vendor

# Build the Go app using the vendored modules
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o /app/api ./cmd/api

# Stage 2: Create the minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install postgresql-client for pg_isready, and ca-certificates for https downloads if migrate CLI needs it (should be static)
RUN apk add --no-cache postgresql-client ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/api /app/api

# Copy the migrate CLI from the builder stage
COPY --from=builder /usr/local/bin/migrate /app/migrate

# Copy migration files
# Assuming your migrations are in a directory named 'migrations' at the root of minilab_backend
COPY migrations ./migrations

# Copy the entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Set Gin to release mode (good practice for production/staging)
ENV GIN_MODE=release

# Expose port 8080 to the outside world
EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]

# Command to run the executable (will be passed as arguments to entrypoint.sh)
CMD ["/app/api"] 