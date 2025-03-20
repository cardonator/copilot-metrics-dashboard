import { ServerActionResponse } from "@/features/common/server-action-response";

interface GitHubConfig {
  organization: string;
  enterprise: string;
  token: string;
  version: string;
  scope: string;
}

interface FeaturesConfig {
  dashboard: boolean;
  seats: boolean;
}

// Helper function to check if database is configured
export const isDatabaseConfigured = (): boolean => {
  // Check for CosmosDB configuration
  const cosmosDbConfigured = !stringIsNullOrEmpty(process.env.AZURE_COSMOSDB_ENDPOINT) && 
                            !stringIsNullOrEmpty(process.env.AZURE_COSMOSDB_KEY);
  
  // Check for SQLite configuration
  const sqliteEnabled = process.env.ENABLE_SQLITE === 'true';
  
  return cosmosDbConfigured || sqliteEnabled;
}

export const ensureGitHubEnvConfig = (): ServerActionResponse<GitHubConfig> => {
  const organization = process.env.GITHUB_ORGANIZATION || "";
  const enterprise = process.env.GITHUB_ENTERPRISE || "";
  const token = process.env.GITHUB_TOKEN || "";
  const version = process.env.GITHUB_API_VERSION || "2022-11-28";
  let scope = process.env.GITHUB_API_SCOPE;

  // If database is configured, all GitHub settings are optional
  const dbConfigured = isDatabaseConfigured();
  
  if (!dbConfigured) {
    // Only validate if database is not configured
    if (stringIsNullOrEmpty(organization)) {
      console.log("Missing required environment variable for organization");
      return {
        status: "ERROR",
        errors: [
          {
            message: "Missing required environment variable for organization",
          },
        ],
      };
    }

    if (stringIsNullOrEmpty(token)) {
      return {
        status: "ERROR",
        errors: [
          {
            message: "Missing required environment variable for GitHub token",
          },
        ],
      };
    }

    if (stringIsNullOrEmpty(version)) {
      return {
        status: "ERROR",
        errors: [
          {
            message:
              "Missing required environment variable for GitHub API version",
          },
        ],
      };
    }
  }

  // Only validate scope if it's provided
  if (!stringIsNullOrEmpty(scope) && validateScope(scope)) {
    return {
      status: "ERROR",
      errors: [
        {
          message:
            "Invalid GitHub API scope: " +
            scope +
            ". Value must be 'enterprise' or 'organization'",
        },
      ],
    };
  }

  if (stringIsNullOrEmpty(scope)) {
    scope = "organization";
  }

  return {
    status: "OK",
    response: {
      organization,
      enterprise,
      token,
      version,
      scope,
    },
  };
};

export const featuresEnvConfig = (): ServerActionResponse<FeaturesConfig> => {
  const enableDashboardFeature = process.env.ENABLE_DASHBOARD_FEATURE !== "false" ? true : false;
  const enableSeatsFeature = process.env.ENABLE_SEATS_FEATURE !== "false" ? true : false;
  return {
    status: "OK",
    response: {
      dashboard: enableDashboardFeature,
      seats: enableSeatsFeature,
    },
  };
};

export const stringIsNullOrEmpty = (str: string | null | undefined) => {
  return str === null || str === undefined || str === "";
};

export const validateScope = (str: string | null | undefined) => {
  return str !== "enterprise" && str !== "organization";
};
