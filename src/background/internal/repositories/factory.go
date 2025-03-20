package repositories

import (
	"context"

	"github.com/cardonator/copilot-metrics-dashboard/internal/config"
	"go.uber.org/zap"
)

// CreateRepository creates a repository based on the provided configuration
func CreateRepository(cfg *config.Config, logger *zap.Logger) (Repository, error) {
	var repo Repository
	var err error

	switch cfg.StorageType {
	case config.StorageSQLite:
		logger.Info("Creating SQLite repository", zap.String("path", cfg.SQLitePath))
		repo, err = NewSQLiteRepository(cfg.SQLitePath, logger)
	case config.StorageCosmos:
		logger.Info("Creating Cosmos DB repository", zap.String("endpoint", cfg.CosmosDBEndpoint))
		repo, err = NewCosmosRepository(cfg.CosmosDBEndpoint, cfg.CosmosDBKey, logger)
	}

	if err != nil {
		return nil, err
	}

	if repo != nil {
		// Initialize repository
		ctx := context.Background()
		if err := repo.Initialize(ctx); err != nil {
			return nil, err
		}
	}

	return repo, nil
}
