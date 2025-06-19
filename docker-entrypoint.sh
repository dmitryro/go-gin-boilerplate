#!/bin/sh
# docker-entrypoint.sh

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h postgres -U postgres -d platform_db; do
  echo "PostgreSQL is not ready yet, waiting..."
  sleep 5  # Increase the wait time to 10 seconds
done

echo "PostgreSQL is up, starting the Go application..."

exec /app/main
