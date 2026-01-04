package core

import (
	"context"
	"strings"
	"testing"

	"github.com/antonkhramov/cribot/internal/config"
	"github.com/antonkhramov/cribot/internal/plugins"
)

// mockNotifier collects sent messages for testing.
type mockNotifier struct {
	messages []string
}

func (m *mockNotifier) Send(ctx context.Context, message string) error {
	m.messages = append(m.messages, message)
	return nil
}

// mockPlugin is a test plugin that always triggers.
type mockPlugin struct {
	name    string
	trigger bool
}

func (m *mockPlugin) Name() string {
	return m.name
}

func (m *mockPlugin) Check(ctx context.Context, ticker string, cfg config.TickerConfig) (*plugins.Result, error) {
	return &plugins.Result{
		Triggered:    m.trigger,
		Message:      ticker + " triggered",
		CurrentValue: 100,
	}, nil
}

func TestEngine_Run(t *testing.T) {
	// Setup
	plugins.Clear()
	plugins.Register(&mockPlugin{name: "test-plugin", trigger: true})

	csvData := `ticker,plugin,enabled,threshold_type,threshold_value,target_value,notes
SBER,test-plugin,true,below,250,300,
VTBR,test-plugin,true,below,30,,`

	cfg, err := config.LoadFromReader(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	notifier := &mockNotifier{}
	engine := New(cfg, notifier, nil)

	// Run
	ctx := context.Background()
	result, err := engine.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify
	if result.Stats.TotalTickers != 2 {
		t.Errorf("expected 2 tickers, got %d", result.Stats.TotalTickers)
	}
	if result.Stats.TriggeredCount != 2 {
		t.Errorf("expected 2 triggers, got %d", result.Stats.TriggeredCount)
	}
	if result.Stats.NotificationsSent != 2 {
		t.Errorf("expected 2 notifications, got %d", result.Stats.NotificationsSent)
	}
	if len(notifier.messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(notifier.messages))
	}

	plugins.Clear()
}

func TestEngine_SkipsDisabled(t *testing.T) {
	plugins.Clear()
	plugins.Register(&mockPlugin{name: "test-plugin", trigger: true})

	csvData := `ticker,plugin,enabled,threshold_type,threshold_value,target_value,notes
SBER,test-plugin,true,below,250,300,
VTBR,test-plugin,false,below,30,,`

	cfg, err := config.LoadFromReader(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	notifier := &mockNotifier{}
	engine := New(cfg, notifier, nil)

	ctx := context.Background()
	result, err := engine.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only 1 ticker should be processed (VTBR is disabled)
	if result.Stats.TotalTickers != 1 {
		t.Errorf("expected 1 enabled ticker, got %d", result.Stats.TotalTickers)
	}

	plugins.Clear()
}

func TestEngine_MissingPlugin(t *testing.T) {
	plugins.Clear()
	// Don't register any plugins

	csvData := `ticker,plugin,enabled,threshold_type,threshold_value,target_value,notes
SBER,nonexistent,true,below,250,300,`

	cfg, err := config.LoadFromReader(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	notifier := &mockNotifier{}
	engine := New(cfg, notifier, nil)

	ctx := context.Background()
	result, err := engine.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should count as an error
	if result.Stats.ErrorCount != 1 {
		t.Errorf("expected 1 error, got %d", result.Stats.ErrorCount)
	}
	if result.Stats.NotificationsSent != 0 {
		t.Errorf("expected 0 notifications, got %d", result.Stats.NotificationsSent)
	}

	plugins.Clear()
}

func TestEngine_NoTrigger(t *testing.T) {
	plugins.Clear()
	plugins.Register(&mockPlugin{name: "test-plugin", trigger: false})

	csvData := `ticker,plugin,enabled,threshold_type,threshold_value,target_value,notes
SBER,test-plugin,true,below,250,300,`

	cfg, err := config.LoadFromReader(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	notifier := &mockNotifier{}
	engine := New(cfg, notifier, nil)

	ctx := context.Background()
	result, err := engine.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No triggers, no notifications
	if result.Stats.TriggeredCount != 0 {
		t.Errorf("expected 0 triggers, got %d", result.Stats.TriggeredCount)
	}
	if result.Stats.NotificationsSent != 0 {
		t.Errorf("expected 0 notifications, got %d", result.Stats.NotificationsSent)
	}

	plugins.Clear()
}
