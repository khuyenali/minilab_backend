#!/bin/sh
set -e # Exit immediately if a command exits with a non-zero status

# Construct Database URL for golang-migrate
# Format: postgresql://user:password@host:port/dbname?sslmode=sslmode
DATABASE_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

MIGRATIONS_PATH="/app/migrations"

echo "Waiting for database to be ready..."
# Simple loop to wait for DB to be available. 
# For production, a more robust solution like wait-for-it.sh or dockerize -wait is recommended.
c=0
while ! pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -q; do
  c=$((c+1))
  if [ "$c" -ge 30 ]; then # timeout after 30 seconds (30 * 1s sleep)
    echo "Error: Database not ready after 30 seconds."
    exit 1
  fi
  sleep 1
done
echo "Database is ready."

echo "Running database migrations from $MIGRATIONS_PATH..."

# Run migrations
# The migrate binary will be in /app/migrate in the container
/app/migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" up

echo "Database migrations complete."

# Now execute the main command (passed as CMD in Dockerfile or arguments to this script)
echo "Starting application..."
exec "$@" 