import { ErrorPage } from "@/features/common/error-page";
import { RawDataPage } from "@/features/raw/raw-page";
import { getFeatures } from "@/utils/helpers";
import { cosmosConfiguration } from "@/services/cosmos-db-service";
import {
  getRawCopilotMetrics,
  getRawCopilotSeats,
} from "@/services/raw-data-service";
import { redirect } from "next/navigation";
import { isDatabaseConfigured } from "@/services/env-service";

export default async function RawData() {
  const features = getFeatures();
  const isCosmosDb = cosmosConfiguration();
  const useDb = isDatabaseConfigured();

  // If Cosmos DB or other database is being used, redirect to home
  // since we won't have raw API responses
  if (useDb) {
    redirect("/");
  }

  // Fetch both types of raw data
  const metricsDataPromise = getRawCopilotMetrics();
  const seatsDataPromise = features.seats ? getRawCopilotSeats() : Promise.resolve(null);
  
  const [metricsData, seatsData] = await Promise.all([metricsDataPromise, seatsDataPromise]);
  
  if (metricsData.status !== "OK") {
    return <ErrorPage error={metricsData.errors[0].message} />;
  }

  return <RawDataPage 
    metricsData={metricsData.response} 
    seatsData={seatsData && 'response' in seatsData ? seatsData.response : null} 
  />;
}
