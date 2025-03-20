import {
  formatResponseError,
  unknownResponseError,
} from "@/features/common/response-error";
import {
  CopilotMetrics,
  CopilotUsageOutput,
} from "@/features/common/models";
import { ServerActionResponse } from "@/features/common/server-action-response";
import { SqlQuerySpec } from "@azure/cosmos";
import { format } from "date-fns";
import { cosmosClient, cosmosConfiguration, getDatabaseType } from "./cosmos-db-service";
import { ensureGitHubEnvConfig } from "./env-service";
import { stringIsNullOrEmpty, applyTimeFrameLabel } from "../utils/helpers";
import { sampleData } from "./sample-data";
import { queryDb } from "./sqlite-db-service";

export interface IFilter {
  startDate?: Date;
  endDate?: Date;
  enterprise: string;
  organization: string;
  team: string;
}

export const getCopilotMetrics = async (
  filter: IFilter
): Promise<ServerActionResponse<CopilotUsageOutput[]>> => {
  
  const databaseType = getDatabaseType();
  

  try {
    // If we have a database configured, try to use it first
    if (databaseType === 'cosmos' || databaseType === 'sqlite') {
      if (databaseType === 'cosmos') {
        return getCopilotMetricsFromDatabase(filter);
      } else {
        // New SQLite implementation
        let query = "SELECT data FROM metrics_history";
        const params = [];
        
        // Add date filter if provided
        if (filter.startDate && filter.endDate) {
          query += " WHERE date BETWEEN ? AND ?";
          params.push(
            format(filter.startDate, "yyyy-MM-dd"), 
            format(filter.endDate, "yyyy-MM-dd")
          );
        }
        
        // Add enterprise filter if provided
        if (!stringIsNullOrEmpty(filter.enterprise)) {
          const jsonCondition = `json_extract(data, '$.enterprise') = ?`;
          query += params.length ? ` AND ${jsonCondition}` : ` WHERE ${jsonCondition}`;
          params.push(filter.enterprise);
        }
        
        // Add organization filter if provided
        if (!stringIsNullOrEmpty(filter.organization)) {
          const jsonCondition = `json_extract(data, '$.organization') = ?`;
          query += params.length ? ` AND ${jsonCondition}` : ` WHERE ${jsonCondition}`;
          params.push(filter.organization);
        }
        
        // Order by date
        query += " ORDER BY date DESC";
        
        const result = await queryDb<{ data: string }>(query, params);
        
        if (result.status === "OK" && result.response.length > 0) {
          const metricsData = result.response.map(item => JSON.parse(item.data) as CopilotMetrics);
          return {
            status: "OK",
            response: applyTimeFrameLabel(metricsData),
          };
        }
      }
    }

    // Fall back to GitHub API if no database or database query returned no results
    const env = ensureGitHubEnvConfig();
    
    if (env.status !== "OK") {
      return env;
    }

    const { enterprise, organization } = env.response;
    switch (process.env.GITHUB_API_SCOPE) {
      case "enterprise":
        if (stringIsNullOrEmpty(filter.enterprise)) {
          filter.enterprise = enterprise || '';
        }
        break;
      default:
        if (stringIsNullOrEmpty(filter.organization)) {
          filter.organization = organization || '';
        }
        break;
    }
    return getCopilotMetricsFromApi(filter);
  } catch (e) {
    return unknownResponseError(e);
  }
};

const fetchCopilotMetrics = async (
  url: string,
  token: string,
  version: string,
  entityName: string
): Promise<ServerActionResponse<CopilotUsageOutput[]>> => {
  const response = await fetch(url, {
    cache: "no-store",
    headers: {
      Accept: `application/vnd.github+json`,
      Authorization: `Bearer ${token}`,
      "X-GitHub-Api-Version": version,
    },
  });

  if (!response.ok) {
    return formatResponseError(entityName, response);
  }

  // Get the raw text to preserve it
  const rawResponse = await response.text();
  
  // Parse the response
  const data = JSON.parse(rawResponse);
  const dataWithTimeFrame = applyTimeFrameLabel(data);
  
  // Add the raw response to the result
  Object.defineProperty(dataWithTimeFrame, 'rawApiResponse', {
    enumerable: true,
    value: rawResponse
  });
  
  return {
    status: "OK",
    response: dataWithTimeFrame,
  };
};

export const getCopilotMetricsFromApi = async (
  filter: IFilter
): Promise<ServerActionResponse<CopilotUsageOutput[]>> => {
  const env = ensureGitHubEnvConfig();

  if (env.status !== "OK") {
    return env;
  }

  const { token, version } = env.response;

  try {
    const queryParams = new URLSearchParams();
    
    if (filter.startDate) {
      queryParams.append('since', format(filter.startDate, "yyyy-MM-dd"));
    }
    if (filter.endDate) {
      queryParams.append('until', format(filter.endDate, "yyyy-MM-dd"));
    }
    
    const queryString = queryParams.toString() ? `?${queryParams.toString()}` : '';

    if (filter.enterprise) {
      const url = `https://api.github.com/enterprises/${filter.enterprise}/copilot/metrics${queryString}`;
      return fetchCopilotMetrics(url, token, version, filter.enterprise || '');
    } else {
      const url = `https://api.github.com/orgs/${filter.organization}/copilot/metrics${queryString}`;
      return fetchCopilotMetrics(url, token, version, filter.organization || '');
    }
  } catch (e) {
    return unknownResponseError(e);
  }
};

export const getCopilotMetricsFromDatabase = async (
  filter: IFilter
): Promise<ServerActionResponse<CopilotUsageOutput[]>> => {
  const client = cosmosClient();
  const database = client.database("platform-engineering");
  const container = database.container("metrics_history");

  let start = "";
  let end = "";
  const maxDays = 365 * 2; // maximum 2 years of data
  const maximumDays = 31;

  if (filter.startDate && filter.endDate) {
    start = format(filter.startDate, "yyyy-MM-dd");
    end = format(filter.endDate, "yyyy-MM-dd");
  } else {
    // set the start date to today and the end date to 31 days ago
    const todayDate = new Date();
    const startDate = new Date(todayDate);
    startDate.setDate(todayDate.getDate() - maximumDays);

    start = format(startDate, "yyyy-MM-dd");
    end = format(todayDate, "yyyy-MM-dd");
  }

  let querySpec: SqlQuerySpec = {
    query: `SELECT * FROM c WHERE c.date >= @start AND c.date <= @end`,
    parameters: [
      { name: "@start", value: start },
      { name: "@end", value: end },
    ],
  };

  if (filter.enterprise) {
    querySpec.query += ` AND c.enterprise = @enterprise`;
    querySpec.parameters?.push({
      name: "@enterprise",
      value: filter.enterprise,
    });
  }

  if (filter.organization) {
    querySpec.query += ` AND c.organization = @organization`;
    querySpec.parameters?.push({
      name: "@organization",
      value: filter.organization,
    });
  }
  
  if (filter.team) {
    querySpec.query += ` AND c.team = @team`;
    querySpec.parameters?.push({ name: "@team", value: filter.team });
  }

  const { resources } = await container.items
    .query<CopilotMetrics>(querySpec, {
      maxItemCount: maxDays,
    })
    .fetchAll();

  const dataWithTimeFrame = applyTimeFrameLabel(resources);
  return {
    status: "OK",
    response: dataWithTimeFrame,
  };
};

export const _getCopilotMetrics = (): Promise<CopilotUsageOutput[]> => {
  const promise = new Promise<CopilotUsageOutput[]>((resolve) => {
    setTimeout(() => {
      const weekly = applyTimeFrameLabel(sampleData);
      resolve(weekly);
    }, 1000);
  });

  return promise;
};
