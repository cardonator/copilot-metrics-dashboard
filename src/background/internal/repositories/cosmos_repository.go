package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/cardonator/copilot-metrics-dashboard/internal/models"
	"go.uber.org/zap"
)

// CosmosRepository implements Repository using Azure Cosmos DB
type CosmosRepository struct {
	client *azcosmos.Client
	logger *zap.Logger
}

// NewCosmosRepository creates a new Cosmos DB repository
func NewCosmosRepository(endpoint, key string, logger *zap.Logger) (*CosmosRepository, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("Cosmos DB endpoint is not specified")
	}

	if key == "" {
		return nil, fmt.Errorf("Cosmos DB key is not specified")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cosmos DB credentials: %w", err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cosmos DB client: %w", err)
	}

	return &CosmosRepository{
		client: client,
		logger: logger,
	}, nil
}

// Initialize sets up the necessary containers
func (r *CosmosRepository) Initialize(ctx context.Context) error {
	// Cosmos DB containers are created automatically in the handlers
	// via the CreateIfNotExists option
	return nil
}

// SaveMetrics stores metrics data in Cosmos DB
func (r *CosmosRepository) SaveMetrics(ctx context.Context, metrics []models.Metrics) error {
	container, err := r.client.NewContainer("platform-engineering", "metrics_history")
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		// Set ID if not already set
		if metric.ID == "" {
			metric.ID = metric.GetID()
		}

		data, err := json.Marshal(metric)
		if err != nil {
			r.logger.Warn("Failed to marshal metric", zap.Error(err))
			continue
		}

		_, err = container.UpsertItem(ctx, azcosmos.PartitionKey{}, data, nil)
		if err != nil {
			r.logger.Warn("Failed to upsert metric", zap.String("id", metric.ID), zap.Error(err))
			continue
		}

		r.logger.Info("Saved metric", zap.String("id", metric.ID))
	}

	return nil
}

// SaveSeats stores seats data in Cosmos DB
func (r *CosmosRepository) SaveSeats(ctx context.Context, seats *models.CopilotAssignedSeats) error {
	container, err := r.client.NewContainer("platform-engineering", "seats_history")
	if err != nil {
		return err
	}

	if seats.ID == "" {
		seats.ID = seats.GetID()
	}

	data, err := json.Marshal(seats)
	if err != nil {
		return err
	}

	_, err = container.UpsertItem(ctx, azcosmos.PartitionKey{}, data, nil)
	if err != nil {
		return err
	}

	r.logger.Info("Saved seats", zap.String("id", seats.ID), zap.Int("totalSeats", seats.TotalSeats))

	return nil
}

// SaveUsage stores usage data in Cosmos DB
func (r *CosmosRepository) SaveUsage(ctx context.Context, usageData []models.CopilotUsage) error {
	container, err := r.client.NewContainer("platform-engineering", "usage_history")
	if err != nil {
		return err
	}

	for _, usage := range usageData {
		if usage.ID == "" {
			usage.ID = usage.GetID()
		}

		data, err := json.Marshal(usage)
		if err != nil {
			r.logger.Warn("Failed to marshal usage data", zap.Error(err))
			continue
		}

		_, err = container.UpsertItem(ctx, azcosmos.PartitionKey{}, data, nil)
		if err != nil {
			r.logger.Warn("Failed to upsert usage data", zap.String("id", usage.ID), zap.Error(err))
			continue
		}

		r.logger.Info("Saved usage data", zap.String("id", usage.ID), zap.String("day", usage.Day))
	}

	return nil
}

// Close cleans up resources
func (r *CosmosRepository) Close() error {
	// Cosmos DB client doesn't require explicit closing
	return nil
}
