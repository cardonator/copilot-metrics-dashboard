#!/bin/bash
set -e

echo "Starting GitHub Copilot Metrics Dashboard..."

# Ensure the SQLite database directory exists
mkdir -p "$(dirname "$SQLITE_DB_PATH")"

# Initialize the SQLite database if it doesn't exist
if [ "$ENABLE_SQLITE" = "true" ] && [ ! -f "$SQLITE_DB_PATH" ]; then
  echo "Initializing SQLite database at $SQLITE_DB_PATH"
  touch "$SQLITE_DB_PATH"
fi

# Start the background data ingestion job in the background
echo "Starting background data ingestion job..."
cd /app/background
./dataingestion &
BACKGROUND_PID=$!

# Wait a moment to ensure the background job has started
sleep 4

# Start the dashboard application
echo "Starting dashboard application..."
cd /app/dashboard

# Handle termination signals to gracefully shut down both processes
trap 'kill $BACKGROUND_PID; exit' SIGINT SIGTERM

exec node server.js
