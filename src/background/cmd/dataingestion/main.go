package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cardonator/copilot-metrics-dashboard/internal/config"
	"github.com/cardonator/copilot-metrics-dashboard/internal/handlers"
	"github.com/cardonator/copilot-metrics-dashboard/internal/repositories"
	"github.com/cardonator/copilot-metrics-dashboard/internal/services"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting GitHub Copilot Metrics Dashboard data ingestion")

	// Load configuration
	cfg, err := config.Load(logger)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Log if using test data
	if cfg.UseTestData {
		logger.Info("Running in test mode with test data")
	}

	// Set up GitHub client
	githubClient := services.NewGitHubClient(
		cfg.GithubApiBaseUrl,
		cfg.GithubToken,
		cfg.GithubApiVersion,
		logger,
	)

	// Set up service clients
	metricsClient := services.NewCopilotMetricsClient(githubClient, logger)
	seatsClient := services.NewCopilotSeatsClient(githubClient, logger)
	usageClient := services.NewCopilotUsageClient(githubClient, logger)

	// Set up repository based on configuration
	var repo repositories.Repository

	switch cfg.StorageType {
	case config.StorageSQLite:
		if cfg.SQLitePath == "" {
			logger.Warn("SQLite path is not set.")
		} else {
			logger.Info("Using SQLite repository", zap.String("path", cfg.SQLitePath))
			sqlite, err := repositories.NewSQLiteRepository(cfg.SQLitePath, logger)
			if err != nil {
				logger.Error("Failed to create SQLite repository", zap.Error(err))
			} else {
				repo = sqlite
			}
		}
	case config.StorageCosmos:
		if cfg.CosmosDBEndpoint == "" || cfg.CosmosDBKey == "" {
			logger.Warn("Cosmos DB endpoint or key is not set.")
		} else {
			logger.Info("Using Cosmos DB repository", zap.String("endpoint", cfg.CosmosDBEndpoint))
			cosmos, err := repositories.NewCosmosRepository(cfg.CosmosDBEndpoint, cfg.CosmosDBKey, logger)
			if err != nil {
				logger.Error("Failed to create Cosmos DB repository", zap.Error(err))
			} else {
				repo = cosmos
			}
		}
	default:
		logger.Warn("No storage type specified. Data will be collected but not persisted.")
	}

	// Initialize repository if it exists
	if repo != nil {
		ctx := context.Background()
		if err := repo.Initialize(ctx); err != nil {
			logger.Error("Failed to initialize repository, data will not be persisted", zap.Error(err))
			repo = nil
		} else {
			defer repo.Close()
		}
	}

	if repo == nil {
		logger.Warn("No valid repository configured. Data will be collected but not persisted.")
	}

	// Set up handlers
	metricsHandler := handlers.NewMetricsHandler(
		logger,
		metricsClient,
		repo,
		cfg.Teams,
		cfg.UseTestData,
	)

	seatsHandler := handlers.NewSeatsHandler(
		logger,
		seatsClient,
		repo,
		cfg.UseTestData,
	)

	usageHandler := handlers.NewUsageHandler(
		logger,
		usageClient,
		repo,
		cfg.UseTestData,
	)

	// Set up scheduler
	scheduler := gocron.NewScheduler(time.UTC)

	// Schedule metrics ingestion - every hour
	_, err = scheduler.Every(1).Hour().Do(func() {
		ctx := context.Background()
		if err := metricsHandler.Run(ctx); err != nil {
			logger.Error("Metrics ingestion failed", zap.Error(err))
		}
	})
	if err != nil {
		logger.Fatal("Failed to schedule metrics ingestion", zap.Error(err))
	}

	// Schedule seats ingestion - every hour
	_, err = scheduler.Every(1).Hour().Do(func() {
		ctx := context.Background()
		if err := seatsHandler.Run(ctx); err != nil {
			logger.Error("Seats ingestion failed", zap.Error(err))
		}
	})
	if err != nil {
		logger.Fatal("Failed to schedule seats ingestion", zap.Error(err))
	}

	// Schedule usage ingestion - every hour
	_, err = scheduler.Every(1).Hour().Do(func() {
		ctx := context.Background()
		if err := usageHandler.Run(ctx); err != nil {
			logger.Error("Usage ingestion failed", zap.Error(err))
		}
	})
	if err != nil {
		logger.Fatal("Failed to schedule usage ingestion", zap.Error(err))
	}

	// Start the scheduler in a non-blocking manner
	scheduler.StartAsync()

	// Run once immediately
	ctx := context.Background()
	logger.Info("Running initial data collection")

	if err := metricsHandler.Run(ctx); err != nil {
		logger.Error("Initial metrics ingestion failed", zap.Error(err))
	} else {
		logger.Info("Initial metrics ingestion completed successfully")
	}

	if err := seatsHandler.Run(ctx); err != nil {
		logger.Error("Initial seats ingestion failed", zap.Error(err))
	} else {
		logger.Info("Initial seats ingestion completed successfully")
	}

	if err := usageHandler.Run(ctx); err != nil {
		logger.Error("Initial usage ingestion failed", zap.Error(err))
	} else {
		logger.Info("Initial usage ingestion completed successfully")
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigChan
	logger.Info("Received signal, shutting down", zap.String("signal", sig.String()))

	// Stop the scheduler
	scheduler.Stop()

	logger.Info("Shutdown complete")
}
