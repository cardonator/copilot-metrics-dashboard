# Copilot Metrics Dashboard

A visualization dashboard for GitHub Copilot metrics built with Next.js.

## Features

- Interactive data visualizations
- Filter metrics by date range
- Export metrics as CSV
- Team and individual performance insights
- Integration with GitHub Copilot API

## Prerequisites

- Node.js 18+ and npm/yarn
- GitHub Copilot license and access
- Authorization credentials for the Copilot API

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/copilot-metrics-dashboard.git
cd copilot-metrics-dashboard/src/dashboard
```

2. Install dependencies:
```bash
npm install
# or
yarn install
```

## Environment Setup

1. Create a `.env.local` file in the project root:
```bash
# GitHub API settings (required only if not using database options below)
GITHUB_ORGANIZATION=your_github_organization
GITHUB_ENTERPRISE=your_github_enterprise_name  # can be omitted if not using GitHub Enterprise
GITHUB_TOKEN=your_github_token
GITHUB_API_VERSION=your_github_api_version  # e.g., "2022-11-28"
GITHUB_API_SCOPE=organization  # or "enterprise" if using GitHub Enterprise

# Database settings - if either of these options is configured, 
# the GitHub API settings above are not required:

# Option 1: Azure CosmosDB
AZURE_COSMOSDB_ENDPOINT=your_cosmosdb_endpoint
AZURE_COSMOSDB_KEY=your_cosmosdb_key

# Option 2: SQLite (local storage)
ENABLE_SQLITE=true
# Optional: Specify a custom path for the SQLite database
# SQLITE_DB_PATH=/custom/path/to/copilot-metrics.db

# Feature flags (optional)
ENABLE_DASHBOARD_FEATURE=true
ENABLE_SEATS_FEATURE=true
```

2. Replace the placeholder values with your actual credentials.

## Running the Dashboard

Start the development server:

```bash
npm run dev
# or
yarn dev
```

The dashboard will be available at [http://localhost:3000](http://localhost:3000).

## Usage Guide

### Viewing Metrics

1. Navigate to the main dashboard page
2. Select the desired date range using the date picker
3. Browse through different metric categories:
   - Acceptance Rate
   - Time Saved
   - Suggestions Count
   - Team Performance

### Filtering Data

Use the filters in the sidebar to narrow down metrics by:
- Team members
- Projects
- Programming languages
- Time period

### Exporting Data

1. Apply desired filters
2. Click the "Export" button in the top-right corner
3. Choose your preferred format (CSV, JSON)

## Configuration Options

The dashboard offers several configuration options:

- **Environment Variables**: Adjust settings in your `.env.local` file as described in the Environment Setup section
- **Feature Flags**: Enable/disable features using the environment variables like `ENABLE_DASHBOARD_FEATURE` and `ENABLE_SEATS_FEATURE`
- **Database Selection**: Choose between CosmosDB and SQLite by setting the appropriate environment variables

## Troubleshooting

- **API Connection Issues**: Verify your GitHub token has the appropriate scopes
- **No Data Appearing**: Check date range filters and API permissions
- **Slow Performance**: Consider reducing the date range or applying more filters

## Contributing

Contributions are welcome! Please see our [contributing guide](github.com/cardonator/copilot-metrics-dashboard/CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
