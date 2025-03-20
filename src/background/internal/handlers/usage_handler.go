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

// UsageHandler handles the processing of Copilot usage data
type UsageHandler struct {
	logger      *zap.Logger
	usageClient *services.CopilotUsageClient
	repository  repositories.Repository
	useTestData bool
}

// NewUsageHandler creates a new usage handler
func NewUsageHandler(
	logger *zap.Logger,
	usageClient *services.CopilotUsageClient,
	repository repositories.Repository,
	useTestData bool,
) *UsageHandler {
	return &UsageHandler{
		logger:      logger,
		usageClient: usageClient,
		repository:  repository,
		useTestData: useTestData,
	}
}

// Run runs the usage data ingestion process
func (h *UsageHandler) Run(ctx context.Context) error {
	h.logger.Info("Running GitHub Copilot usage ingestion")

	var usageData []models.CopilotUsage
	var err error

	if h.useTestData {
		h.logger.Info("Using test data for usage ingestion")
		usageData, err = h.usageClient.LoadTestUsageData()
	} else {
		scope := os.Getenv("GITHUB_API_SCOPE")
		if strings.ToLower(scope) == "enterprise" {
			enterprise := os.Getenv("GITHUB_ENTERPRISE")
			h.logger.Info("Fetching GitHub Copilot usage for enterprise", zap.String("enterprise", enterprise))
			usageData, err = h.usageClient.GetCopilotUsageForEnterprise(enterprise)
		} else {
			organization := os.Getenv("GITHUB_ORGANIZATION")
			h.logger.Info("Fetching GitHub Copilot usage for organization", zap.String("organization", organization))
			usageData, err = h.usageClient.GetCopilotUsageForOrganization(organization)
		}
	}

	if err != nil {
		h.logger.Error("Failed to get usage data", zap.Error(err))
		return err
	}

	h.logger.Info("Usage data extracted", zap.Int("count", len(usageData)))

	// Save usage data to repository if available
	if h.repository != nil {
		if err := h.repository.SaveUsage(ctx, usageData); err != nil {
			h.logger.Error("Failed to save usage data", zap.Error(err))
			return err
		}
	} else {
		h.logger.Info("Repository not available, skipping save operation")
	}

	return nil
}
