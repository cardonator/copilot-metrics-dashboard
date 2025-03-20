package repositories

import (
	"context"
	"encoding/json"

	"github.com/cardonator/copilot-metrics-dashboard/internal/models"
)

// Repository defines the interface for data storage implementations
type Repository interface {
	// Initialize prepares the repository for use
	Initialize(ctx context.Context) error

	// SaveMetrics stores metrics data
	SaveMetrics(ctx context.Context, metrics []models.Metrics) error

	// SaveSeats stores seats data
	SaveSeats(ctx context.Context, seats *models.CopilotAssignedSeats) error

	// SaveUsage stores usage data
	SaveUsage(ctx context.Context, usage []models.CopilotUsage) error

	// Close closes the repository
	Close() error
}

// DataMarshaler provides a common way to marshal data for storage
func DataMarshaler(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
