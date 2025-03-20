import { getCopilotMetrics, getCopilotMetricsFromApi } from "./copilot-metrics-service";
import { getCopilotSeats, getCopilotSeatsFromApi } from "./copilot-seat-service";
import { isDatabaseConfigured } from "./env-service";
import { ensureGitHubEnvConfig } from "./env-service";

export const getRawCopilotMetrics = async () => {
  const env = ensureGitHubEnvConfig();
  
  // Create a default filter based on env config
  const defaultFilter = {
    enterprise: env.status === "OK" ? env.response.enterprise || "" : "",
    organization: env.status === "OK" ? env.response.organization || "" : "",
    team: ""
  };
  
  // If database is configured, we need to fetch directly from API to get raw data
  if (isDatabaseConfigured()) {
    return await getCopilotMetricsFromApi(defaultFilter);
  } else {
    // Otherwise, we use the regular service which will call the API
    return await getCopilotMetrics(defaultFilter);
  }
};

export const getRawCopilotSeats = async () => {
  const env = ensureGitHubEnvConfig();
  
  // Create a default filter based on env config
  const defaultFilter = {
    enterprise: env.status === "OK" ? env.response.enterprise || "" : "",
    organization: env.status === "OK" ? env.response.organization || "" : "",
    team: ""
  };
  
  // If database is configured, we need to fetch directly from API to get raw data
  if (isDatabaseConfigured()) {
    return await getCopilotSeatsFromApi(defaultFilter);
  } else {
    // Otherwise, we use the regular service which will call the API
    return await getCopilotSeats(defaultFilter);
  }
};
