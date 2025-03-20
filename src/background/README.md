# Copilot Metrics Dashboard - Data Ingestion Service

This Go application collects GitHub Copilot metrics, seats, and usage data and stores it in either Azure Cosmos DB or SQLite.

## Prerequisites

- Go 1.22 or later
- GitHub API token with appropriate permissions
- Azure Cosmos DB instance (optional, if using Cosmos DB storage)

## Configuration

The application is configured through environment variables:

- `GITHUB_TOKEN` - GitHub API token
- `GITHUB_API_BASEURL` - GitHub API base URL (default: https://api.github.com)
- `GITHUB_API_VERSION` - GitHub API version (default: 2022-11-28)
- `GITHUB_API_SCOPE` - Scope of data collection (enterprise or organization)
- `GITHUB_ENTERPRISE` - Enterprise name (when scope is enterprise)
- `GITHUB_ORGANIZATION` - Organization name (when scope is not enterprise)
- `STORAGE_TYPE` - Storage type to use: "cosmos" or "sqlite" (default: "cosmos")
- `AZURE_COSMOSDB_ENDPOINT` - Azure Cosmos DB endpoint (required if storage type is cosmos)
- `AZURE_COSMOSDB_KEY` - Azure Cosmos DB key (required if storage type is cosmos)
- `SQLITE_DB_PATH` - Path to SQLite database file (optional, default: ~/.copilot-metrics/copilot-metrics.db)
- `GITHUB_METRICS_TEAMS` - Comma-separated list of teams to collect metrics for
- `GITHUB_METRICS_USE_TESTDATA` - Set to "true" to use test data instead of calling the GitHub API
- `ENABLE_SEATS_INGESTION` - Set to "false" to disable seats ingestion

You can set these variables in a `.env` file in the project root.

## Building

```bash
# Build with CGO disabled (pure Go, no C dependencies)
CGO_ENABLED=0 go build -o dataingestion ./cmd/dataingestion
```

## Running

```bash
./dataingestion
```

The application will:
1. Collect metrics, seats, and usage data immediately upon startup
2. Schedule hourly collection of this data
3. Store the data in the configured storage (Azure Cosmos DB or SQLite)

## Development

To run with test data:

```bash
GITHUB_METRICS_USE_TESTDATA=true ./dataingestion
```

To use SQLite instead of Cosmos DB:

```bash
STORAGE_TYPE=sqlite ./dataingestion
```

The test data is located in the `testdata/` directory.
