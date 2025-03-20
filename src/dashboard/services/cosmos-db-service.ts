import { ServerActionResponse } from "@/features/common/server-action-response";
import { CosmosClient } from "@azure/cosmos";
import { DefaultAzureCredential } from "@azure/identity";
import { stringIsNullOrEmpty } from "../utils/helpers";

export const cosmosClient = () => {
  const endpoint = process.env.AZURE_COSMOSDB_ENDPOINT;

  if(stringIsNullOrEmpty(endpoint)) {
    throw new Error("Missing required environment variable for CosmosDB endpoint");
  }

  const credential = new DefaultAzureCredential();
  return new CosmosClient({ endpoint, aadCredentials: credential });
};

export const cosmosConfiguration = (): boolean => {
  const endpoint = process.env.AZURE_COSMOSDB_ENDPOINT;

  return (
    endpoint !== undefined &&
    endpoint.trim() !== ""
  );
};

// New function to determine which database to use
export const getDatabaseType = (): 'cosmos' | 'sqlite' | null => {
  if (cosmosConfiguration()) {
    return 'cosmos';
  }
  
  // Import dynamically to avoid circular dependencies
  const { sqliteConfiguration } = require('./sqlite-db-service');
  if (sqliteConfiguration()) {
    return 'sqlite';
  }
  
  return null;
};