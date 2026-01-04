// Package core provides the main orchestration engine for CriBot.
package core

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/antonkhramov/cribot/internal/config"
	"github.com/antonkhramov/cribot/internal/plugins"
)

// Notifier is the interface for sending notifications.
type Notifier interface {
	Send(ctx context.Context, message string) error
}

// CheckResult represents the result of checking a single ticker.
type CheckResult struct {
	Ticker    string
	Plugin    string
	Triggered bool
	Message   string
	Value     float64
	Error     error
}

// RunResult contains the full results of a check cycle.
type RunResult struct {
	Stats   RunStats
	Results []CheckResult
}

// RunStats contains execution statistics.
type RunStats struct {
	TotalTickers      int
	TriggeredCount    int
	NotificationsSent int
	ErrorCount        int
	Duration          time.Duration
}

// Engine orchestrates the checking of all configured tickers.
type Engine struct {
	config   *config.Config
	notifier Notifier
	logger   *slog.Logger
}

// New creates a new Engine instance.
func New(cfg *config.Config, notifier Notifier, logger *slog.Logger) *Engine {
	if logger == nil {
		logger = slog.Default()
	}
	return &Engine{
		config:   cfg,
		notifier: notifier,
		logger:   logger,
	}
}

// Run executes the full check cycle for all enabled tickers.
func (e *Engine) Run(ctx context.Context) (*RunResult, error) {
	start := time.Now()
	stats := RunStats{}

	enabledTickers := e.config.EnabledTickers()
	stats.TotalTickers = len(enabledTickers)

	e.logger.Info("starting check cycle",
		"total_tickers", stats.TotalTickers,
	)

	// Check all tickers in parallel
	results := e.checkAllParallel(ctx, enabledTickers)

	// Process results and send notifications
	for _, result := range results {
		if result.Error != nil {
			stats.ErrorCount++
			e.logger.Error("plugin check failed",
				"ticker", result.Ticker,
				"plugin", result.Plugin,
				"error", result.Error,
			)
			continue
		}

		if result.Triggered {
			stats.TriggeredCount++
			e.logger.Info("trigger fired",
				"ticker", result.Ticker,
				"plugin", result.Plugin,
				"value", result.Value,
			)

			if e.notifier != nil {
				if err := e.notifier.Send(ctx, result.Message); err != nil {
					e.logger.Error("failed to send notification",
						"ticker", result.Ticker,
						"error", err,
					)
				} else {
					stats.NotificationsSent++
				}
			}
		}
	}

	stats.Duration = time.Since(start)

	e.logger.Info("check cycle complete",
		"total", stats.TotalTickers,
		"triggered", stats.TriggeredCount,
		"notifications", stats.NotificationsSent,
		"errors", stats.ErrorCount,
		"duration_ms", stats.Duration.Milliseconds(),
	)

	return &RunResult{
		Stats:   stats,
		Results: results,
	}, nil
}

// checkAllParallel checks all tickers concurrently.
func (e *Engine) checkAllParallel(ctx context.Context, tickers []config.TickerConfig) []CheckResult {
	var wg sync.WaitGroup
	results := make([]CheckResult, len(tickers))

	for i, tc := range tickers {
		wg.Add(1)
		go func(idx int, tickerCfg config.TickerConfig) {
			defer wg.Done()
			results[idx] = e.checkOne(ctx, tickerCfg)
		}(i, tc)
	}

	wg.Wait()
	return results
}

// checkOne checks a single ticker against its configured plugin.
func (e *Engine) checkOne(ctx context.Context, tc config.TickerConfig) CheckResult {
	result := CheckResult{
		Ticker: tc.Ticker,
		Plugin: tc.Plugin,
	}

	// Get the plugin
	plugin := plugins.Get(tc.Plugin)
	if plugin == nil {
		e.logger.Error("plugin not found", "plugin", tc.Plugin)
		result.Error = fmt.Errorf("plugin '%s' not found", tc.Plugin)
		return result
	}

	e.logger.Info("checking ticker", "ticker", tc.Ticker, "plugin", tc.Plugin, "threshold", fmt.Sprintf("%s %f", tc.ThresholdType, tc.ThresholdValue))

	// Run the check
	pluginResult, err := plugin.Check(ctx, tc.Ticker, tc)
	if err != nil {
		e.logger.Error("check failed", "ticker", tc.Ticker, "error", err)
		result.Error = err
		return result
	}

	result.Triggered = pluginResult.Triggered
	result.Message = pluginResult.Message
	result.Value = pluginResult.CurrentValue

	e.logger.Info("check result", "ticker", tc.Ticker, "value", result.Value, "triggered", result.Triggered)

	return result
}
