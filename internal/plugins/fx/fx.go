// Package fx provides a plugin for checking currency exchange rates.
package fx

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
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

	plugins.Register(&Plugin{
		client: moex.NewClient(ttl),
	})
}

// Plugin checks currency exchange rates against thresholds.
type Plugin struct {
	client moex.Client
}

// Name returns the plugin identifier.
func (p *Plugin) Name() string {
	return "fx"
}

// Check evaluates the FX rate condition for the given pair.
// Ticker format: USDRUB, EURRUB, etc.
func (p *Plugin) Check(ctx context.Context, ticker string, cfg config.TickerConfig) (*plugins.Result, error) {
	currentRate, err := p.getRate(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate for %s: %w", ticker, err)
	}

	var triggered bool
	var message string

	switch cfg.ThresholdType {
	case config.ThresholdBelow:
		triggered = currentRate < cfg.ThresholdValue
		if triggered {
			message = fmt.Sprintf("💱 %s: курс %.2f ниже порога %.2f", ticker, currentRate, cfg.ThresholdValue)
		}
	case config.ThresholdAbove:
		triggered = currentRate > cfg.ThresholdValue
		if triggered {
			message = fmt.Sprintf("💱 %s: курс %.2f выше порога %.2f", ticker, currentRate, cfg.ThresholdValue)
		}
	}

	if !triggered {
		message = fmt.Sprintf("%s: курс %.2f (порог %.2f %s)", ticker, currentRate, cfg.ThresholdValue, cfg.ThresholdType)
	}

	return &plugins.Result{
		Triggered:    triggered,
		Message:      message,
		CurrentValue: currentRate,
	}, nil
}

// getRate tries to fetch from MOEX, falling back to mock on error.
func (p *Plugin) getRate(ctx context.Context, ticker string) (float64, error) {
	// Map common pairs to MOEX format (e.g. USDRUB -> USDRUB_TOM)
	moexTicker := ticker
	if len(ticker) == 6 && !strings.Contains(ticker, "_") {
		moexTicker = ticker + "_TOM"
	}

	// Try MOEX API
	data, err := p.client.GetCurrencyRate(ctx, moexTicker)
	if err == nil {
		return data.Last, nil
	}

	// Log error but continue with fallback
	slog.Warn("moex api failed, using mock data", "ticker", ticker, "moex_ticker", moexTicker, "error", err)

	// Fallback to mock data
	return mockFXRate(ticker), nil
}

// mockFXRate returns a mock exchange rate.
func mockFXRate(pair string) float64 {
	pair = strings.ToUpper(pair)

	// Base rates
	baseRates := map[string]float64{
		"USDRUB": 92.5,
		"EURRUB": 100.3,
		"CNYRUB": 12.8,
		"GBPRUB": 117.5,
	}

	base, ok := baseRates[pair]
	if !ok {
		base = 1.0
	}

	// Add some random variation (-2% to +2%)
	variation := (rand.Float64() - 0.5) * 0.04 * base
	return base + variation
}
