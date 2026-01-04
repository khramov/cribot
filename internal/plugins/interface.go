// Package plugins defines the Source interface for all data source plugins.
package plugins

import (
	"context"

	"github.com/antonkhramov/cribot/internal/config"
)

// Result represents the outcome of a plugin check.
type Result struct {
	Triggered    bool    // Whether the condition was met
	Message      string  // Human-readable description
	CurrentValue float64 // The actual value that was checked
}

// Source is the interface that all data source plugins must implement.
type Source interface {
	// Name returns the plugin identifier used in CSV config.
	Name() string

	// Check evaluates the condition for the given ticker and config.
	// Returns a Result indicating whether the trigger fired.
	Check(ctx context.Context, ticker string, cfg config.TickerConfig) (*Result, error)
}
