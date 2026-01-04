// Package price provides a plugin for checking asset prices.
package price

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/antonkhramov/cribot/internal/config"
	"github.com/antonkhramov/cribot/internal/moex"
	"github.com/antonkhramov/cribot/internal/plugins"
)

func init() {
	ttl := 5 * time.Minute
	if val := os.Getenv("MOEX_CACHE_TTL"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			ttl = d
		}
	}

	// Register plugin with default MOEX client
	plugins.Register(&Plugin{
		client: moex.NewClient(ttl),
	})
}

// Plugin checks asset prices against thresholds.
type Plugin struct {
	client moex.Client
}

// Name returns the plugin identifier.
func (p *Plugin) Name() string {
	return "price"
}

// Check evaluates the price condition for the given ticker.
func (p *Plugin) Check(ctx context.Context, ticker string, cfg config.TickerConfig) (*plugins.Result, error) {
	currentPrice, err := p.getPrice(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to get price for %s: %w", ticker, err)
	}

	var triggered bool
	var message string

	switch cfg.ThresholdType {
	case config.ThresholdBelow:
		triggered = currentPrice < cfg.ThresholdValue
		if triggered {
			message = fmt.Sprintf("📉 %s: цена %.2f ниже порога %.2f", ticker, currentPrice, cfg.ThresholdValue)
		}
	case config.ThresholdAbove:
		triggered = currentPrice > cfg.ThresholdValue
		if triggered {
			message = fmt.Sprintf("📈 %s: цена %.2f выше порога %.2f", ticker, currentPrice, cfg.ThresholdValue)
		}
	}

	if !triggered {
		message = fmt.Sprintf("%s: цена %.2f (порог %.2f %s)", ticker, currentPrice, cfg.ThresholdValue, cfg.ThresholdType)
	}

	return &plugins.Result{
		Triggered:    triggered,
		Message:      message,
		CurrentValue: currentPrice,
	}, nil
}

// getPrice tries to fetch from MOEX. Fallback to mock only if CRIBOT_MOCK is true.
func (p *Plugin) getPrice(ctx context.Context, ticker string) (float64, error) {
	// Try MOEX API
	data, err := p.client.GetStockPrice(ctx, ticker)
	if err == nil {
		return data.Last, nil
	}

	// Falls back to mock only if explicitly enabled
	if os.Getenv("CRIBOT_MOCK") == "true" {
		slog.Warn("moex api failed, using mock data (CRIBOT_MOCK=true)", "ticker", ticker, "error", err)
		return mockPrice(ticker), nil
	}

	return 0, fmt.Errorf("moex api error: %w", err)
}

// mockPrice returns a random price for testing.
// This simulates API responses during development.
func mockPrice(ticker string) float64 {
	// Base prices for known tickers
	basePrices := map[string]float64{
		"SBER": 260,
		"GAZP": 155,
		"LKOH": 7200,
		"VTBR": 0.025,
		"YNDX": 3500,
	}

	base, ok := basePrices[ticker]
	if !ok {
		base = 100
	}

	// Add some random variation (-10% to +10%)
	variation := (rand.Float64() - 0.5) * 0.2 * base
	return base + variation
}
