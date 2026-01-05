package models

import (
	"time"
)

// FeatureFlag represents a feature flag in the system
type FeatureFlag struct {
	Key       string    `json:"key"`
	Enabled   bool      `json:"enabled"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Known feature flags
const (
	// FlagRestrictionsEnabled controls whether all restrictions are enforced
	// When false: no email verification, no identity verification, no gender restrictions, no wealth status requirements
	// When true: all restrictions are enforced
	FlagRestrictionsEnabled = "restrictions_enabled"
)

// DefaultFeatureFlags returns the default values for all feature flags
func DefaultFeatureFlags() map[string]bool {
	return map[string]bool{
		FlagRestrictionsEnabled: false, // Start with no restrictions
	}
}

