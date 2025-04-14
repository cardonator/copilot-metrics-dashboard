# GitHub Copilot Metrics Dashboard

A dashboard for tracking GitHub Copilot usage metrics across your organization.

## Prerequisites

- Node.js (v18 or higher)
- npm or yarn
- A GitHub personal access token (PAT) with appropriate permissions
- Access to a GitHub organization with Copilot subscriptions

## Setup

1. Clone the repository:
```
git clone https://github.com/microsoft/copilot-metrics-dashboard.git
```

2. Navigate to the project directory:
```
cd copilot-metrics-dashboard
```

3. Install dependencies for both the background service and dashboard:
```
npm install
```

4. Configure environment variables:
- Copy the example environment files in both directories:
  ```
  cp src/background/.env.example src/background/.env.local
  cp src/dashboard/.env.example src/dashboard/.env.local
  ```
- Update the `.env.local` files with your specific settings

5. **Important: GitHub Authentication**
- You must uncomment and set a valid GitHub token in `src/dashboard/.env.local`:
  ```
  GITHUB_ORGANIZATION=your-organization
  GITHUB_TOKEN=your-personal-access-token
  ```
- Ensure your token has `read:org` permissions at minimum
- If you encounter authorization errors, check that your token is valid and has the correct permissions

6. Database Configuration:
- The system uses SQLite by default (`ENABLE_SQLITE=true`)
- For SQLite, set the correct path to your database:
  ```
  SQLITE_DB_PATH=../background/copilot-metrics.db
  ```
- Relative paths are based on the execution directory, absolute paths are recommended for reliability

## Running the Application

### Start the Background Service

The background service collects metrics from GitHub and stores them in the database.

This service should be run periodically (daily or weekly) to refresh metrics data.

### Start the Dashboard
```
npm start
```

The dashboard will be available at http://localhost:3000 by default.

## Troubleshooting

### Authentication Errors

If you see "Authorization failed" errors, check:
- Your GitHub token is uncommented in the `.env.local` file
- The token has not expired
- The token has proper permissions for your organization
- The organization name is correctly specified

### Database Issues

- Ensure the path to the SQLite database is correct 
- Check that the background service has successfully created the database
- For first-time setup, run the background service before starting the dashboard

## License

[License information]

# GitHub Copilot Metrics - Dashboard

1. [Introduction](#introduction)
2. [Deploy to Azure](#deploy-to-azure)

# Introduction

GitHub Copilot Metrics Dashboard is a reimplementation and fork of the official microsoft/copilot-metrics-dashboard repository. The background job has been ported to golang. The frontend has had numerous changes and enhancements. Much of this README remains from that project for now, however more detailed README files specific to this fork can be found in the src/background and src/dashboard folders.

**NOTE**: This has not been tested deploying to Azure using the accelerator or az library discussed below.

---

The GitHub Copilot Metrics Dashboard is a solution accelerator designed to visualize metrics from GitHub Copilot using the [GitHub Copilot Metrics API](https://docs.github.com/en/enterprise-cloud@latest/rest/copilot/copilot-metrics?apiVersion=2022-11-28) and [GitHub Copilot User Management API](https://docs.github.com/en/enterprise-cloud@latest/rest/copilot/copilot-user-management?apiVersion=2022-11-28).

## Dashboard

![GitHub Copilot Metrics - Dashboard](/docs/dashboard.jpeg "GitHub Copilot Metrics - Dashboard")

The dashboard showcases a range of features:

**Filters:**
Ability to filter metrics by date range, languages, code editors and visualise data by time frame (daily, weekly, monthly).

**Acceptance Average:** Percentage of suggestions accepted by users for given date range and group by time range (daily, weekly, monthly).

**Active Users:** Number of active users for the last cycle.

**Adoption Rate:** Number of active users who are using GitHub Copilot in relation to the total number of licensed users.

**Seat Information:** Number of active, inactive, and total users.

**Language:** Breakdown of languages which can be used to filter the data.

**Code Editors:** Breakdown of code editors which can be used to filter the data.

## Seats

Seats feature shows the list of user having a Copilot licence assigned.
This feature is can be enabled or disabled by setting the `ENABLE_SEATS_FEATURE` environment variable to `true` or `false` respectively (default value is `true`).

> Assigned seats ingestion is enabled by default, is possbile to disable by setting the `ENABLE_SEATS_INGESTION` environment variable to `false`

# Deploy to Azure

The solution accelerator is a web application that uses Azure App Service, Azure Functions, Azure Cosmos DB, Azure Storage and Azure Key Vault. The deployment template will automatically populate the required environment variables in Azure Key Vault and configure the application settings in Azure App Service and Azure Functions.
![GitHub Copilot Metrics - Architecture ](/docs/CopilotDashboard.png "GitHub Copilot Metrics - Architecture")

The following steps will automatically provision Azure resources and deploy the solution accelerator to Azure App Service and Azure Functions using the Azure Developer CLI.

> [!IMPORTANT]
> 🚨🚨🚨 You must setup [authentication](https://learn.microsoft.com/en-us/azure/app-service/overview-authentication-authorization) using the built-in authentication and authorization capabilities of Azure App Service.

#### Prerequisites

You will be prompted to provide the following information:

```
- GitHub Enterprise name
- GitHub Organization name
- GitHub Token
- GitHub API Scope
- Team Names (if you choose to use the new metrics API)
```

> More details here for the [GA Metrics API](https://github.blog/changelog/2024-10-30-github-copilot-metrics-api-ga-release-now-available/)

> Team Names must be a valid JSON array, e.g. ``["team-1", "team-2]``

GitHub API Scope define the GITHUB_API_SCOPE environment variable that can be "enterprise" or "organization". It is used to define at which level the GitHub APIs will gather data. If not specified, the default value is "organization".

1. Download the [Azure Developer CLI](https://learn.microsoft.com/en-us/azure/developer/azure-developer-cli/overview)
2. If you have not cloned this repo, run `azd init -t microsoft/copilot-metrics-dashboard`. If you have cloned this repo, just run 'azd init' from the repo root directory.
3. Run `azd up` to provision and deploy the application

```pwsh
azd init -t microsoft/copilot-metrics-dashboard
azd up

# if you are wanting to see logs run with debug flag
azd up --debug
```

# Contributing

This project welcomes contributions and suggestions. 

# Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft
trademarks or logos is subject to and must follow
[Microsoft&#39;s Trademark &amp; Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.
