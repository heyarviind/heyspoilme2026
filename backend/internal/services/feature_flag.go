package services

import (
	"log"
	"sync"
	"time"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
)

type FeatureFlagService struct {
	repo  *repository.FeatureFlagRepository
	cache map[string]bool
	mu    sync.RWMutex
}

func NewFeatureFlagService(repo *repository.FeatureFlagRepository) *FeatureFlagService {
	s := &FeatureFlagService{
		repo:  repo,
		cache: make(map[string]bool),
	}

	// Ensure table exists and initialize defaults
	if err := repo.EnsureTable(); err != nil {
		log.Printf("Warning: Failed to ensure feature_flags table: %v", err)
	}
	if err := repo.InitializeDefaults(); err != nil {
		log.Printf("Warning: Failed to initialize default feature flags: %v", err)
	}

	// Load initial cache
	s.refreshCache()

	// Start background refresh
	go s.backgroundRefresh()

	return s
}

// refreshCache loads all flags from database into cache
func (s *FeatureFlagService) refreshCache() {
	flags, err := s.repo.GetAll()
	if err != nil {
		log.Printf("Error refreshing feature flag cache: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Start with defaults
	s.cache = models.DefaultFeatureFlags()

	// Override with database values
	for _, flag := range flags {
		s.cache[flag.Key] = flag.Enabled
	}
}

// backgroundRefresh periodically refreshes the cache
func (s *FeatureFlagService) backgroundRefresh() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.refreshCache()
	}
}

// IsEnabled checks if a feature flag is enabled
func (s *FeatureFlagService) IsEnabled(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if enabled, ok := s.cache[key]; ok {
		return enabled
	}

	// Return default if not in cache
	defaults := models.DefaultFeatureFlags()
	if defaultVal, ok := defaults[key]; ok {
		return defaultVal
	}

	return false
}

// RestrictionsEnabled is a convenience method to check the main restrictions flag
func (s *FeatureFlagService) RestrictionsEnabled() bool {
	return s.IsEnabled(models.FlagRestrictionsEnabled)
}

// SetFlag updates a feature flag
func (s *FeatureFlagService) SetFlag(key string, enabled bool) error {
	if err := s.repo.Set(key, enabled); err != nil {
		return err
	}

	// Update cache immediately
	s.mu.Lock()
	s.cache[key] = enabled
	s.mu.Unlock()

	log.Printf("[FeatureFlag] %s set to %v", key, enabled)
	return nil
}

// GetAll returns all feature flags
func (s *FeatureFlagService) GetAll() ([]models.FeatureFlag, error) {
	return s.repo.GetAll()
}

// GetAllWithDefaults returns all flags including defaults for any missing
func (s *FeatureFlagService) GetAllWithDefaults() []models.FeatureFlag {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var flags []models.FeatureFlag
	for key, enabled := range s.cache {
		flags = append(flags, models.FeatureFlag{
			Key:       key,
			Enabled:   enabled,
			UpdatedAt: time.Now(), // Not accurate but good enough for display
		})
	}
	return flags
}

