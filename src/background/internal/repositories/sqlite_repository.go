package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cardonator/copilot-metrics-dashboard/internal/models"
	"go.uber.org/zap"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

const schema = `
CREATE TABLE IF NOT EXISTS metrics_history (
    id TEXT PRIMARY KEY,
    date TEXT NOT NULL,
    data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS seats_history (
    id TEXT PRIMARY KEY,
    date TEXT NOT NULL,
    data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS usage_history (
    id TEXT PRIMARY KEY,
    day TEXT NOT NULL,
    data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

// SQLiteRepository implements Repository using SQLite
type SQLiteRepository struct {
	db     *sql.DB
	logger *zap.Logger
	path   string
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string, logger *zap.Logger) (*SQLiteRepository, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("SQLite database path is not specified")
	}

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Use a connection string with the necessary pragmas for modernc.org/sqlite
	connString := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)", dbPath)
	db, err := sql.Open("sqlite", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable foreign keys and set pragmas that might improve performance
	if _, err := db.Exec("PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL; PRAGMA synchronous = NORMAL"); err != nil {
		logger.Warn("Failed to set SQLite pragmas", zap.Error(err))
	}

	return &SQLiteRepository{
		db:     db,
		logger: logger,
		path:   dbPath,
	}, nil
}

// Initialize sets up the database schema
func (r *SQLiteRepository) Initialize(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}
	return nil
}

// SaveMetrics stores metrics data in SQLite
func (r *SQLiteRepository) SaveMetrics(ctx context.Context, metrics []models.Metrics) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR REPLACE INTO metrics_history (id, date, data)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

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

		_, err = stmt.ExecContext(ctx, metric.ID, metric.Date, string(data))
		if err != nil {
			r.logger.Warn("Failed to insert metric", zap.String("id", metric.ID), zap.Error(err))
			continue
		}

		r.logger.Info("Saved metric", zap.String("id", metric.ID))
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// SaveSeats stores seats data in SQLite
func (r *SQLiteRepository) SaveSeats(ctx context.Context, seats *models.CopilotAssignedSeats) error {
	if seats.ID == "" {
		seats.ID = seats.GetID()
	}

	data, err := json.Marshal(seats)
	if err != nil {
		return fmt.Errorf("failed to marshal seats: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT OR REPLACE INTO seats_history (id, date, data)
		VALUES (?, ?, ?)
	`, seats.ID, seats.Date, string(data))

	if err != nil {
		return fmt.Errorf("failed to insert seats: %w", err)
	}

	r.logger.Info("Saved seats", zap.String("id", seats.ID), zap.Int("totalSeats", seats.TotalSeats))
	return nil
}

// SaveUsage stores usage data in SQLite
func (r *SQLiteRepository) SaveUsage(ctx context.Context, usageData []models.CopilotUsage) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR REPLACE INTO usage_history (id, day, data)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, usage := range usageData {
		if usage.ID == "" {
			usage.ID = usage.GetID()
		}

		data, err := json.Marshal(usage)
		if err != nil {
			r.logger.Warn("Failed to marshal usage", zap.Error(err))
			continue
		}

		_, err = stmt.ExecContext(ctx, usage.ID, usage.Day, string(data))
		if err != nil {
			r.logger.Warn("Failed to insert usage", zap.String("id", usage.ID), zap.Error(err))
			continue
		}

		r.logger.Info("Saved usage data", zap.String("id", usage.ID), zap.String("day", usage.Day))
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Close cleans up resources
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
