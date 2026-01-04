// Package rsi provides a plugin for checking RSI (Relative Strength Index) values.
package rsi

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/antonkhramov/cribot/internal/config"
	"github.com/antonkhramov/cribot/internal/plugins"
)

func init() {
	plugins.Register(&Plugin{})
}

// Plugin checks RSI values against thresholds.
type Plugin struct{}

// Name returns the plugin identifier.
func (p *Plugin) Name() string {
	return "rsi"
}

// Check evaluates the RSI condition for the given ticker.
func (p *Plugin) Check(ctx context.Context, ticker string, cfg config.TickerConfig) (*plugins.Result, error) {
	// TODO: Replace with real RSI calculation from market data
	if os.Getenv("CRIBOT_MOCK") != "true" {
		return nil, fmt.Errorf("RSI calculation not implemented (use CRIBOT_MOCK=true for dev data)")
	}

	currentRSI := mockRSI(ticker)

	var triggered bool
	var message string

	switch cfg.ThresholdType {
	case config.ThresholdBelow:
		triggered = currentRSI < cfg.ThresholdValue
		if triggered {
			message = fmt.Sprintf("🔻 %s: RSI %.1f (перепроданность, порог %.0f)", ticker, currentRSI, cfg.ThresholdValue)
		}
	case config.ThresholdAbove:
		triggered = currentRSI > cfg.ThresholdValue
		if triggered {
			message = fmt.Sprintf("🔺 %s: RSI %.1f (перекупленность, порог %.0f)", ticker, currentRSI, cfg.ThresholdValue)
		}
	}

	if !triggered {
		message = fmt.Sprintf("%s: RSI %.1f (порог %.0f %s)", ticker, currentRSI, cfg.ThresholdValue, cfg.ThresholdType)
	}

	return &plugins.Result{
		Triggered:    triggered,
		Message:      message,
		CurrentValue: currentRSI,
	}, nil
}

// mockRSI returns a random RSI value between 0 and 100.
func mockRSI(ticker string) float64 {
	// Simulate various RSI states for testing
	switch ticker {
	case "VTBR":
		// Simulate oversold
		return 20 + rand.Float64()*15
	case "SBER":
		// Simulate neutral
		return 40 + rand.Float64()*20
	default:
		// Random RSI
		return rand.Float64() * 100
	}
}
