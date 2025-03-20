package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// StorageType defines the type of storage to use
type StorageType string

const (
	StorageCosmos StorageType = "cosmos"
	StorageSQLite StorageType = "sqlite"
)

// Config holds the application configuration
type Config struct {
	GithubToken        string
	GithubApiBaseUrl   string
	GithubApiVersion   string
	GithubApiScope     string
	GithubEnterprise   string
	GithubOrganization string
	CosmosDBEndpoint   string
	CosmosDBKey        string
	Teams              []string
	UseTestData        bool
	StorageType        StorageType
	SQLitePath         string
}

// Load loads the configuration from environment variables
func Load(logger *zap.Logger) (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		GithubToken:        os.Getenv("GITHUB_TOKEN"),
		GithubApiBaseUrl:   os.Getenv("GITHUB_API_BASEURL"),
		GithubApiVersion:   os.Getenv("GITHUB_API_VERSION"),
		GithubApiScope:     os.Getenv("GITHUB_API_SCOPE"),
		GithubEnterprise:   os.Getenv("GITHUB_ENTERPRISE"),
		GithubOrganization: os.Getenv("GITHUB_ORGANIZATION"),
		CosmosDBEndpoint:   os.Getenv("AZURE_COSMOSDB_ENDPOINT"),
		CosmosDBKey:        os.Getenv("AZURE_COSMOSDB_KEY"),
	}

	// Set default values if not provided
	if config.GithubApiBaseUrl == "" {
		config.GithubApiBaseUrl = "https://api.github.com"
	}

	if config.GithubApiVersion == "" {
		config.GithubApiVersion = "2022-11-28"
	}

	// Parse teams from environment variable
	teamsStr := os.Getenv("GITHUB_METRICS_TEAMS")
	if teamsStr != "" {
		config.Teams = strings.Split(teamsStr, ",")
		for i := range config.Teams {
			config.Teams[i] = strings.TrimSpace(config.Teams[i])
		}
	}

	// Check if using test data
	useTestDataStr := os.Getenv("GITHUB_METRICS_USE_TESTDATA")
	config.UseTestData = strings.ToLower(useTestDataStr) == "true"

	// Configure storage type
	storageType := os.Getenv("STORAGE_TYPE")
	switch strings.ToLower(storageType) {
	case "sqlite":
		config.StorageType = StorageSQLite
		config.SQLitePath = os.Getenv("SQLITE_DB_PATH")
		if config.SQLitePath == "" {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				// Create default path in home directory
				dbDir := filepath.Join(homeDir, ".copilot-metrics")
				if err := os.MkdirAll(dbDir, 0755); err == nil {
					config.SQLitePath = filepath.Join(dbDir, "copilot-metrics.db")
					logger.Info("SQLite path not specified, using default", zap.String("path", config.SQLitePath))
				} else {
					logger.Warn("Failed to create default SQLite directory", zap.Error(err))
				}
			} else {
				logger.Warn("Failed to determine home directory for default SQLite path", zap.Error(err))
			}
		}
	default:
		config.StorageType = StorageCosmos
	}

	// Validate required configuration
	if config.GithubToken == "" {
		logger.Warn("GITHUB_TOKEN not set")
	}

	if config.GithubApiScope == "enterprise" && config.GithubEnterprise == "" {
		logger.Warn("GITHUB_ENTERPRISE not set but GITHUB_API_SCOPE is 'enterprise'")
	}

	if config.GithubApiScope != "enterprise" && config.GithubOrganization == "" {
		logger.Warn("GITHUB_ORGANIZATION not set and GITHUB_API_SCOPE is not 'enterprise'")
	}

	if config.StorageType == StorageCosmos {
		if config.CosmosDBEndpoint == "" {
			logger.Warn("AZURE_COSMOSDB_ENDPOINT not set")
		}

		if config.CosmosDBKey == "" {
			logger.Warn("AZURE_COSMOSDB_KEY not set")
		}
	}

	return config, nil
}
