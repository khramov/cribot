//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/antonkhramov/cribot/internal/config"
	"github.com/antonkhramov/cribot/internal/core"
	
	// Import plugins to register them for the test
	_ "github.com/antonkhramov/cribot/internal/plugins/fx"
	_ "github.com/antonkhramov/cribot/internal/plugins/price"
)

// MockNotifier for integration test
type MockNotifier struct{}

func (m *MockNotifier) Send(ctx context.Context, message string) error {
	return nil
}

func TestIntegration(t *testing.T) {
	// Load test configuration
	cfg, err := config.LoadFromFile("test_config.csv")
	if err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}

	// Create engine with mock notifier
	engine := core.New(cfg, &MockNotifier{}, nil)

	// Run check cycle with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := engine.Run(ctx)
	if err != nil {
		t.Fatalf("engine run failed: %v", err)
	}

	// Verification
	if result.Stats.TotalTickers != 2 {
		t.Errorf("expected 2 tickers, got %d", result.Stats.TotalTickers)
	}

	// Check individual results
	tickersFound := make(map[string]bool)
	for _, res := range result.Results {
		tickersFound[res.Ticker] = true
		
		if res.Error != nil {
			t.Errorf("ticker %s failed: %v", res.Ticker, res.Error)
			continue
		}

		if res.Value <= 0 {
			t.Errorf("ticker %s value should be > 0, got %f", res.Ticker, res.Value)
		}

		t.Logf("Verified %s: %.2f", res.Ticker, res.Value)
	}

	if !tickersFound["SBER"] {
		t.Error("SBER not found in results")
	}
	if !tickersFound["USDRUB_TOM"] {
		t.Error("USDRUB_TOM not found in results")
	}
}
