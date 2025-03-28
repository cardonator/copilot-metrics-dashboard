import {
  formatResponseError,
  unknownResponseError,
} from "@/features/common/response-error";
import { ServerActionResponse, ServerActionSuccess } from "@/features/common/server-action-response";
import {
  ensureGitHubEnvConfig,
  storageEnvConfig,
  stringIsNullOrEmpty,
  StorageConfig,
  GitHubConfig,
} from "./env-service";
import {
  CopilotSeatsData,
  CopilotSeatManagementData,
  CopilotSeatsWithRaw,
} from "@/features/common/models";
import { cosmosClient, cosmosConfiguration } from "./cosmos-db-service";
import { format } from "date-fns";
import { SqlQuerySpec } from "@azure/cosmos";
import { queryDb } from "./sqlite-db-service";

export interface IFilter {
  date?: Date;
  enterprise: string;
  organization: string;
  team: string;
}

// Update to use the interface from models.tsx
export type CopilotSeatsWithRawData = CopilotSeatsWithRaw;

export const getCopilotSeats = async (
  filter: IFilter
): Promise<ServerActionResponse<CopilotSeatsWithRawData>> => {
  const storageConfig = storageEnvConfig();
  if (storageConfig.status === "ERROR") {
    return storageConfig;
  }

  // Type assertion to ServerActionSuccess<StorageConfig>
  const storageSuccessConfig = storageConfig as ServerActionSuccess<StorageConfig>;
  const storageType = storageSuccessConfig.response.type;
  
  const env = ensureGitHubEnvConfig();
  if (env.status === "ERROR") {
    return env;
  }
  
  // Type assertion to ServerActionSuccess<GitHubConfig>
  const envSuccessConfig = env as ServerActionSuccess<GitHubConfig>;
  const { enterprise, organization } = envSuccessConfig.response;

  try {
    // Use the storage type from configuration
    if (storageType === 'cosmosdb' || storageType === 'sqlite') {
      if (storageType === 'cosmosdb') {
        const isCosmosConfig = cosmosConfiguration();
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
        if (isCosmosConfig) {
          const result = await getCopilotSeatsFromDatabase(filter);
          if (result.status === "OK" && result.response) {
            return result;
          }
        }
      } else {
        // SQLite implementation
        let query = "SELECT data FROM seats_history";
        const whereConditions = [];
        const params = [];
        
        // Add conditions to whereConditions array
        if (filter.date) {
          whereConditions.push("date = ?");
          params.push(format(filter.date, "yyyy-MM-dd"));
        }
        
        // Add enterprise filter if provided
        if (!stringIsNullOrEmpty(filter.enterprise)) {
          whereConditions.push("json_extract(data, '$.enterprise') = ?");
          params.push(filter.enterprise);
        }
        
        // Add organization filter if provided
        if (!stringIsNullOrEmpty(filter.organization)) {
          whereConditions.push("json_extract(data, '$.organization') = ?");
          params.push(filter.organization);
        }
        
        // Add WHERE clause if there are any conditions
        if (whereConditions.length > 0) {
          query += " WHERE " + whereConditions.join(" AND ");
        }
        
        // Add ORDER BY clause at the end
        if (!filter.date) {
          // Get the most recent data if no date specified
          query += " ORDER BY date DESC LIMIT 1";
        }
        
        const result = await queryDb<{ data: string }>(query, params);
        
        if (result.status === "OK" && result.response.length > 0) {
          return {
            status: "OK",
            response: JSON.parse(result.response[0].data) as CopilotSeatsData,
          };
        }
      }
    }
    
    // If no database or no data in database, fall back to GitHub API
    if (enterprise && !stringIsNullOrEmpty(filter.enterprise)) {
      return getCopilotSeatsFromApi(filter);
    }
    else if (organization && !stringIsNullOrEmpty(filter.organization)) {
      return getCopilotSeatsFromApi(filter);
    } else {
      // If no enterprise or organization is specified, return an error
      return {
        status: "ERROR",
        errors: [{
          message: "No enterprise or organization specified"
        }]
      };
    }
  } catch (e) {
    return unknownResponseError(e);
  }
};

const getCopilotSeatsFromDatabase = async (
  filter: IFilter
): Promise<ServerActionResponse<CopilotSeatsData>> => {
  const client = cosmosClient();
  const database = client.database("platform-engineering");
  const container = database.container("seats_history");

  let date = "";
  const maxDays = 365 * 2; // maximum 2 years of data

  if (filter.date) {
    date = format(filter.date, "yyyy-MM-dd");
  } else {
    const today = new Date();
    date = format(today, "yyyy-MM-dd");
  }

  let querySpec: SqlQuerySpec = {
    query: `SELECT * FROM c WHERE c.date = @date`,
    parameters: [{ name: "@date", value: date }],
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
    .query<CopilotSeatsData>(querySpec, {
      maxItemCount: maxDays,
    })
    .fetchAll();

  if (resources.length === 0) {
    return {
      status: "NOT_FOUND",
      errors: [{ message: "No data found for the specified filter" }],
    };
  }

  return {
    status: "OK",
    response: resources[0],
  };
};

// Export this function so it can be used by raw-data-service
export const getCopilotSeatsFromApi = async (
  filter: IFilter
): Promise<ServerActionResponse<CopilotSeatsWithRawData>> => {
  const env = ensureGitHubEnvConfig();

  if (env.status !== "OK") {
    return env;
  }

  const { enterprise, organization, token, version } = env.response;

  if (enterprise && !stringIsNullOrEmpty(filter.enterprise)) {
    // Define with proper type that includes rawApiResponse
    let enterpriseSeats: CopilotSeatsWithRawData = {
      seats: [] as any[],
      total_seats: 0,
      last_update: format(new Date(), "yyyy-MM-ddTHH:mm:ss"),
      date: format(new Date(), "yyyy-MM-dd"),
      id: `${new Date()}-ENT-${filter.enterprise}`,
      enterprise: filter.enterprise,
      organization: null,
      rawApiResponse: '',  // Initialize the property
    };
    
    // Raw response will be stored here
    const rawResponses: string[] = [];
    
    try {
      let url = `https://api.github.com/enterprises/${filter.enterprise}/copilot/billing/seats`;
      do {
        const enterpriseResponse = await fetch(url, {
          cache: "no-store",
          headers: {
            Accept: `application/vnd.github+json`,
            Authorization: `Bearer ${token}`,
            "X-GitHub-Api-Version": version || "2022-11-28",
          },
        });

        if (!enterpriseResponse.ok) {
          return formatResponseError(filter.enterprise, enterpriseResponse);
        }

        const enterpriseText = await enterpriseResponse.text();
        rawResponses.push(enterpriseText);
        
        const enterpriseData = JSON.parse(enterpriseText);
        enterpriseSeats.seats.push(...enterpriseData.seats);
        enterpriseSeats.total_seats += enterpriseData.total_seats;

        // Check if there's a 'next' link in the Link header
        const linkHeader = enterpriseResponse.headers.get('link');
        url = getNextUrlFromLinkHeader(linkHeader);
      } while (url);

      // Store the raw API response
      enterpriseSeats.rawApiResponse = rawResponses.join('\n\n--- NEXT PAGE ---\n\n');
      
      return {
        status: "OK",
        response: enterpriseSeats
      };
    } catch (e) {
      return unknownResponseError(e);
    }
  } else if (!stringIsNullOrEmpty(filter.organization)) {
    // Define with proper type that includes rawApiResponse
    let organizationSeats: CopilotSeatsWithRawData = {
      seats: [] as any[],
      total_seats: 0,
      last_update: format(new Date(), "yyyy-MM-ddTHH:mm:ss"),
      date: format(new Date(), "yyyy-MM-dd"),
      id: `${new Date()}-ORG-${filter.organization}`,
      enterprise: null,
      organization: filter.organization,
      rawApiResponse: '',  // Initialize the property
    };
    
    const rawResponses: string[] = [];
    
    try {
      let url = `https://api.github.com/orgs/${filter.organization}/copilot/billing/seats`;
      do {
        const organizationResponse = await fetch(url, {
          cache: "no-store",
          headers: {
            Accept: `application/vnd.github+json`,
            Authorization: `Bearer ${token}`,
            "X-GitHub-Api-Version": version || "2022-11-28",
          },
        });

        if (!organizationResponse.ok) {
          return formatResponseError(filter.organization, organizationResponse);
        }

        const organizationText = await organizationResponse.text();
        rawResponses.push(organizationText);
        
        const organizationData = JSON.parse(organizationText);
        organizationSeats.seats.push(...organizationData.seats);
        organizationSeats.total_seats += organizationData.total_seats;

        // Check if there's a 'next' link in the Link header
        const linkHeader = organizationResponse.headers.get('link');
        url = getNextUrlFromLinkHeader(linkHeader);
      } while (url);

      // Store the raw API response
      organizationSeats.rawApiResponse = rawResponses.join('\n\n--- NEXT PAGE ---\n\n');
      
      return {
        status: "OK",
        response: organizationSeats
      };
    } catch (e) {
      return unknownResponseError(e);
    }
  }

  return {
    status: "ERROR",
    errors: [{ message: "No enterprise or organization specified" }],
  };
};

export const getCopilotSeatsManagement = async (
  filter: IFilter
): Promise<ServerActionResponse<CopilotSeatManagementData>> => {
  const env = ensureGitHubEnvConfig();

  if (env.status !== "OK") {
    return env;
  }

  const { enterprise, organization } = env.response;

  try {
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

    const data = await getCopilotSeats(filter);
    if (data.status !== "OK" || !data.response) {
      return unknownResponseError(filter.enterprise);
    }
    const seatsData = data.response;

    // Copilot seats are considered active if they have been active in the last 30 days
    const activeSeats = seatsData.seats.filter((seat) => {
      const lastActivityDate = new Date(seat.last_activity_at);
      const thirtyDaysAgo = new Date();
      thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);
      return lastActivityDate >= thirtyDaysAgo;
    });
    const seatManagementData: CopilotSeatManagementData = {
      enterprise: seatsData.enterprise,
      organization: seatsData.organization,
      date: seatsData.date,
      id: seatsData.id,
      last_update: seatsData.last_update,
      total_seats: seatsData.total_seats,
      seats: {
        seat_breakdown: {
          total: seatsData.seats.length,
          active_this_cycle: activeSeats.length,
          inactive_this_cycle: seatsData.seats.length - activeSeats.length,
          added_this_cycle: 0,
          pending_invitation: 0,
          pending_cancellation: 0,
        },
        seat_management_setting: "",
        public_code_suggestions: "",
        ide_chat: "",
        platform_chat: "",
        cli: "",
        plan_type: "",
      },
    };

    return {
      status: "OK",
      response: seatManagementData as CopilotSeatManagementData,
    };
  } catch (e) {
    return unknownResponseError(e);
  }
};

const getNextUrlFromLinkHeader = (linkHeader: string | null): string => {
  if (!linkHeader) return "";

  const links = linkHeader.split(',');
  for (const link of links) {
    const match = link.match(/<([^>]+)>;\s*rel="([^"]+)"/);
    if (match && match[2] === 'next') {
      return match[1];
    }
  }
  return "";
}
