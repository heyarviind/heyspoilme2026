package repository

import (
	"database/sql"
	"time"

	"heyspoilme/internal/models"
)

type FeatureFlagRepository struct {
	db *sql.DB
}

func NewFeatureFlagRepository(db *sql.DB) *FeatureFlagRepository {
	return &FeatureFlagRepository{db: db}
}

// EnsureTable creates the feature_flags table if it doesn't exist
func (r *FeatureFlagRepository) EnsureTable() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS feature_flags (
			key VARCHAR(100) PRIMARY KEY,
			enabled BOOLEAN NOT NULL DEFAULT false,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

// Get retrieves a feature flag by key
func (r *FeatureFlagRepository) Get(key string) (*models.FeatureFlag, error) {
	var flag models.FeatureFlag
	err := r.db.QueryRow(`
		SELECT key, enabled, updated_at FROM feature_flags WHERE key = $1
	`, key).Scan(&flag.Key, &flag.Enabled, &flag.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &flag, nil
}

// GetAll retrieves all feature flags
func (r *FeatureFlagRepository) GetAll() ([]models.FeatureFlag, error) {
	rows, err := r.db.Query(`
		SELECT key, enabled, updated_at FROM feature_flags ORDER BY key
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flags []models.FeatureFlag
	for rows.Next() {
		var flag models.FeatureFlag
		if err := rows.Scan(&flag.Key, &flag.Enabled, &flag.UpdatedAt); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}
	return flags, nil
}

// Set creates or updates a feature flag
func (r *FeatureFlagRepository) Set(key string, enabled bool) error {
	_, err := r.db.Exec(`
		INSERT INTO feature_flags (key, enabled, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET enabled = $2, updated_at = $3
	`, key, enabled, time.Now().UTC())
	return err
}

// InitializeDefaults ensures all default feature flags exist
func (r *FeatureFlagRepository) InitializeDefaults() error {
	defaults := models.DefaultFeatureFlags()
	for key, defaultEnabled := range defaults {
		// Only insert if not exists (don't overwrite existing values)
		_, err := r.db.Exec(`
			INSERT INTO feature_flags (key, enabled, updated_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (key) DO NOTHING
		`, key, defaultEnabled, time.Now().UTC())
		if err != nil {
			return err
		}
	}
	return nil
}

