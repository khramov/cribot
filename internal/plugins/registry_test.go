package plugins

import (
	"context"
	"testing"

	"github.com/antonkhramov/cribot/internal/config"
)

// mockSource is a test implementation of Source.
type mockSource struct {
	name string
}

func (m *mockSource) Name() string {
	return m.name
}

func (m *mockSource) Check(ctx context.Context, ticker string, cfg config.TickerConfig) (*Result, error) {
	return &Result{Triggered: true, Message: "test", CurrentValue: 100}, nil
}

func TestRegistry(t *testing.T) {
	// Clear registry before test
	Clear()

	// Test empty registry
	if len(List()) != 0 {
		t.Error("expected empty registry")
	}

	// Test Get on empty registry
	if Get("test") != nil {
		t.Error("expected nil for non-existent plugin")
	}

	// Register a plugin
	mock := &mockSource{name: "test"}
	Register(mock)

	// Test List
	plugins := List()
	if len(plugins) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(plugins))
	}

	// Test Get
	retrieved := Get("test")
	if retrieved == nil {
		t.Error("expected to retrieve plugin")
	}
	if retrieved.Name() != "test" {
		t.Errorf("expected 'test', got '%s'", retrieved.Name())
	}

	// Test multiple registrations
	Register(&mockSource{name: "another"})
	if len(List()) != 2 {
		t.Errorf("expected 2 plugins, got %d", len(List()))
	}

	// Clean up
	Clear()
}

func TestSourceInterface(t *testing.T) {
	Clear()

	mock := &mockSource{name: "interface-test"}
	Register(mock)

	s := Get("interface-test")
	if s == nil {
		t.Fatal("expected plugin to be registered")
	}

	// Test Check method
	ctx := context.Background()
	cfg := config.TickerConfig{
		Ticker:         "TEST",
		Plugin:         "interface-test",
		Enabled:        true,
		ThresholdType:  config.ThresholdBelow,
		ThresholdValue: 50,
	}

	result, err := s.Check(ctx, "TEST", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Triggered {
		t.Error("expected triggered to be true")
	}
	if result.CurrentValue != 100 {
		t.Errorf("expected 100, got %f", result.CurrentValue)
	}

	Clear()
}
