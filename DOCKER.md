# Docker Setup for Copilot Metrics Dashboard

This repository includes Docker support for running both the dashboard UI and the background data ingestion job in a single container.

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository:**
   ```sh
   git clone https://github.com/your-repo/copilot-metrics-dashboard.git
   cd copilot-metrics-dashboard
   ```

2. **Create a `.env` file:**
   ```sh
   cp .env.example .env
   # Update the .env file with your configuration
   ```

3. **Build and start the services:**
   ```sh
   docker-compose up --build
   ```

4. **Access the dashboard:**
   Open your browser and go to `http://localhost:3000`.

### Using Docker CLI

1. **Clone the repository:**
   ```sh
   git clone https://github.com/your-repo/copilot-metrics-dashboard.git
   cd copilot-metrics-dashboard
   ```

2. **Create a `.env` file:**
   ```sh
   cp .env.example .env
   # Update the .env file with your configuration
   ```

3. **Build the Docker image:**
   ```sh
   docker build -t copilot-metrics-dashboard .
   ```

4. **Run the Docker container:**
   ```sh
   docker run -p 3000:3000 --env-file .env copilot-metrics-dashboard
   ```

5. **Access the dashboard:**
   Open your browser and go to `http://localhost:3000`.

### Environment Variables

Ensure the following environment variables are set in your `.env` file:

- `GITHUB_ENTERPRISE`
- `GITHUB_ORGANIZATION`
- `GITHUB_TOKEN`
- `GITHUB_API_VERSION` (default: `2022-11-28`)
- `GITHUB_API_SCOPE` (default: `organization`)
- `STORAGE_TYPE` (default: `sqlite`)
- `SQLITE_DB_PATH` (default: `/app/data/copilot-metrics.db`)
- `ENABLE_DASHBOARD_FEATURE` (default: `true`)
- `ENABLE_SEATS_FEATURE` (default: `true`)

### Volumes

The Docker Compose setup uses a named volume `copilot_data` to persist the SQLite database:

```yml
volumes:
  copilot_data:
```
