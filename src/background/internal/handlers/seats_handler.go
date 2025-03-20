package handlers

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/cardonator/copilot-metrics-dashboard/internal/models"
	"github.com/cardonator/copilot-metrics-dashboard/internal/repositories"
	"github.com/cardonator/copilot-metrics-dashboard/internal/services"
	"go.uber.org/zap"
)

// SeatsHandler handles the processing of Copilot seats
type SeatsHandler struct {
	logger      *zap.Logger
	seatsClient *services.CopilotSeatsClient
	repository  repositories.Repository
	useTestData bool
}

// NewSeatsHandler creates a new seats handler
func NewSeatsHandler(
	logger *zap.Logger,
	seatsClient *services.CopilotSeatsClient,
	repository repositories.Repository,
	useTestData bool,
) *SeatsHandler {
	return &SeatsHandler{
		logger:      logger,
		seatsClient: seatsClient,
		repository:  repository,
		useTestData: useTestData,
	}
}

// Run runs the seats ingestion process
func (h *SeatsHandler) Run(ctx context.Context) error {
	h.logger.Info("Running GitHub Copilot seats ingestion")

	// Check if seats ingestion is enabled
	enableSeatsIngestionStr := os.Getenv("ENABLE_SEATS_INGESTION")
	if enableSeatsIngestionStr == "" {
		enableSeatsIngestionStr = "true"
	}

	enableSeatsIngestion, err := strconv.ParseBool(enableSeatsIngestionStr)
	if err != nil {
		h.logger.Warn("Failed to parse ENABLE_SEATS_INGESTION, defaulting to true", zap.Error(err))
		enableSeatsIngestion = true
	}

	if !enableSeatsIngestion {
		h.logger.Info("Seats ingestion is disabled")
		return nil
	}

	// Get seats data
	var seats *models.CopilotAssignedSeats

	if h.useTestData {
		h.logger.Info("Using test data for seats ingestion")
		isEnterprise := strings.ToLower(os.Getenv("GITHUB_API_SCOPE")) == "enterprise"
		seats, err = h.seatsClient.LoadTestSeatsData(isEnterprise)
	} else {
		scope := os.Getenv("GITHUB_API_SCOPE")
		if strings.ToLower(scope) == "enterprise" {
			enterprise := os.Getenv("GITHUB_ENTERPRISE")
			h.logger.Info("Fetching GitHub Copilot seats for enterprise", zap.String("enterprise", enterprise))
			seats, err = h.seatsClient.GetEnterpriseAssignedSeats(enterprise)
		} else {
			organization := os.Getenv("GITHUB_ORGANIZATION")
			h.logger.Info("Fetching GitHub Copilot seats for organization", zap.String("organization", organization))
			seats, err = h.seatsClient.GetOrganizationAssignedSeats(organization)
		}
	}

	if err != nil {
		h.logger.Error("Failed to get seats", zap.Error(err))
		return err
	}

	// Set ID if not already set
	if seats.ID == "" {
		seats.ID = seats.GetID()
	}

	// Save to repository if available
	if h.repository != nil {
		if err := h.repository.SaveSeats(ctx, seats); err != nil {
			h.logger.Error("Failed to save seats", zap.Error(err))
			return err
		}
	} else {
		h.logger.Info("Repository not available, skipping save operation")
	}

	return nil
}
