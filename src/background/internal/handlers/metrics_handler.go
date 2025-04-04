package handlers

import (
	"context"
	"os"
	"strings"

	"github.com/cardonator/copilot-metrics-dashboard/internal/models"
	"github.com/cardonator/copilot-metrics-dashboard/internal/repositories"
	"github.com/cardonator/copilot-metrics-dashboard/internal/services"
	"go.uber.org/zap"
)

// MetricsHandler handles the processing of Copilot metrics
type MetricsHandler struct {
	logger        *zap.Logger
	metricsClient *services.CopilotMetricsClient
	repository    repositories.Repository
	teams         []string
	useTestData   bool
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(
	logger *zap.Logger,
	metricsClient *services.CopilotMetricsClient,
	repository repositories.Repository,
	teams []string,
	useTestData bool,
) *MetricsHandler {
	return &MetricsHandler{
		logger:        logger,
		metricsClient: metricsClient,
		repository:    repository,
		teams:         teams,
		useTestData:   useTestData,
	}
}

// ProcessUsage processes and stores usage data derived from metrics
func (h *MetricsHandler) ProcessUsage(ctx context.Context) error {
	h.logger.Info("Processing GitHub Copilot usage from metrics")

	var usageData []models.CopilotUsage
	var err error

	if h.useTestData {
		h.logger.Info("Using test data for usage processing")
		usageData, err = h.metricsClient.LoadTestUsageData()
	} else {
		scope := os.Getenv("GITHUB_API_SCOPE")
		if strings.ToLower(scope) == "enterprise" {
			enterprise := os.Getenv("GITHUB_ENTERPRISE")
			h.logger.Info("Processing GitHub Copilot usage for enterprise", zap.String("enterprise", enterprise))
			usageData, err = h.metricsClient.GetCopilotUsageForEnterprise(enterprise)
		} else {
			organization := os.Getenv("GITHUB_ORGANIZATION")
			h.logger.Info("Processing GitHub Copilot usage for organization", zap.String("organization", organization))
			usageData, err = h.metricsClient.GetCopilotUsageForOrganization(organization)
		}
	}

	if err != nil {
		h.logger.Error("Failed to process usage data", zap.Error(err))
		return err
	}

	h.logger.Info("Saving usage data", zap.Int("count", len(usageData)))
	if err := h.repository.SaveUsage(ctx, usageData); err != nil {
		h.logger.Error("Failed to save usage data", zap.Error(err))
		return err
	}

	h.logger.Info("Successfully processed and saved usage data")
	return nil
}

// Run runs the metrics ingestion process and also processes usage data
func (h *MetricsHandler) Run(ctx context.Context) error {
	// Run the original metrics ingestion
	if err := h.runMetricsIngestion(ctx); err != nil {
		return err
	}

	// Process usage data from the metrics
	return h.ProcessUsage(ctx)
}

// runMetricsIngestion handles the original metrics ingestion process
func (h *MetricsHandler) runMetricsIngestion(ctx context.Context) error {
	h.logger.Info("Running GitHub Copilot metrics ingestion")

	metrics := []models.Metrics{}

	// Process organization/enterprise level metrics
	orgMetrics, err := h.extractMetrics("")
	if err != nil {
		h.logger.Error("Failed to extract metrics", zap.Error(err))
		return err
	}
	metrics = append(metrics, orgMetrics...)

	// Process team metrics
	if len(h.teams) > 0 {
		for _, team := range h.teams {
			teamMetrics, err := h.extractMetrics(team)
			if err != nil {
				h.logger.Warn("Failed to extract metrics for team", zap.String("team", team), zap.Error(err))
				continue
			}
			metrics = append(metrics, teamMetrics...)
		}
	}

	h.logger.Info("Metrics extracted", zap.Int("count", len(metrics)))

	// Save metrics to repository if available
	if h.repository != nil {
		if err := h.repository.SaveMetrics(ctx, metrics); err != nil {
			h.logger.Error("Failed to save metrics", zap.Error(err))
			return err
		}
	} else {
		h.logger.Info("Repository not available, skipping save operation")
	}

	return nil
}

// extractMetrics extracts Copilot metrics for the given team or organization/enterprise
func (h *MetricsHandler) extractMetrics(team string) ([]models.Metrics, error) {
	if h.useTestData {
		return h.metricsClient.LoadTestMetrics(team)
	}

	scope := os.Getenv("GITHUB_API_SCOPE")
	if strings.ToLower(scope) == "enterprise" {
		enterprise := os.Getenv("GITHUB_ENTERPRISE")
		h.logger.Info("Fetching GitHub Copilot metrics for enterprise", zap.String("enterprise", enterprise), zap.String("team", team))
		return h.metricsClient.GetCopilotMetricsForEnterprise(enterprise, team)
	}

	organization := os.Getenv("GITHUB_ORGANIZATION")
	h.logger.Info("Fetching GitHub Copilot metrics for organization", zap.String("organization", organization), zap.String("team", team))
	return h.metricsClient.GetCopilotMetricsForOrganization(organization, team)
}
