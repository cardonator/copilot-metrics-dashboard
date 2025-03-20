# Multi-stage build for GitHub Copilot Metrics Dashboard

# Stage 1: Build the React dashboard
FROM node:20-alpine AS dashboard-builder
WORKDIR /app

# Copy package files and install dependencies
COPY src/dashboard/package*.json ./
RUN npm ci

# Copy dashboard source code
COPY src/dashboard/ ./

# Build the dashboard app
RUN npm run build

# Stage 2: Build the Go background job
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app

# Copy Go source code
COPY src/background/ ./

# Build the Go application
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o dataingestion ./cmd/dataingestion

# Stage 3: Final runtime image
FROM alpine:3.19

# Install Node.js and necessary packages
RUN apk add --no-cache nodejs npm bash sqlite

# Set working directory
WORKDIR /app

# Copy built assets from previous stages
COPY --from=dashboard-builder /app/.next/standalone ./dashboard
COPY --from=dashboard-builder /app/.next/static ./dashboard/.next/static
COPY --from=dashboard-builder /app/public ./dashboard/public
COPY --from=backend-builder /app/dataingestion ./background/dataingestion

# Create a directory for the SQLite database if using SQLite
RUN mkdir -p /app/data

# Set environment variables (these can be overridden at runtime)
ENV PORT=3000
ENV NODE_ENV=production
ENV STORAGE_TYPE=sqlite
ENV SQLITE_DB_PATH=/app/data/copilot-metrics.db

# Copy startup script
COPY docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh

# Expose port for the dashboard
EXPOSE 3000

# Set entrypoint to the startup script
ENTRYPOINT ["./docker-entrypoint.sh"]
