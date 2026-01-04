// Package main is the entry point for Yandex Cloud Function.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/antonkhramov/cribot/internal/config"
	"github.com/antonkhramov/cribot/internal/core"
	"github.com/antonkhramov/cribot/internal/notify"

	// Import plugins to register them
	_ "github.com/antonkhramov/cribot/internal/plugins/fx"
	_ "github.com/antonkhramov/cribot/internal/plugins/price"
	_ "github.com/antonkhramov/cribot/internal/plugins/rsi"
)

// Response is the HTTP response structure.
type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body"`
}

// HandlerResult is returned in the response body.
type HandlerResult struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message"`
	Total      int                `json:"total_tickers"`
	Triggered  int                `json:"triggered"`
	Sent       int                `json:"notifications_sent"`
	Errors     int                `json:"errors"`
	DurationMs int64              `json:"duration_ms"`
	Results    []core.CheckResult `json:"results,omitempty"`
}

func main() {
	// For local testing, run the handler directly
	// In Yandex Cloud, the function runtime calls Handler
	if os.Getenv("YCF_FUNCTION_ID") == "" {
		// Local execution
		result := Handler(context.Background(), nil)
		fmt.Println(result.Body)
	}
}

// Handler is the Yandex Cloud Function entry point.
// It can be triggered by Timer or HTTP request.
func Handler(ctx context.Context, request interface{}) Response {
	start := time.Now()

	// Setup structured logging
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	logger.Info("function started", "request_type", fmt.Sprintf("%T", request))

	// Load configuration
	configPath := getEnvOrDefault("CONFIG_PATH", "./config/tickers.csv")
	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, fmt.Sprintf("failed to load config: %v", err))
	}

	// Setup Notifier
	tgConfig := notify.TelegramConfig{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	}

	var notifier notify.Notifier
	if tgConfig.BotToken != "" && tgConfig.ChatID != "" {
		if err := tgConfig.Validate(); err != nil {
			return errorResponse(http.StatusInternalServerError, err.Error())
		}
		notifier = notify.NewTelegram(tgConfig)
	} else {
		logger.Info("telegram credentials missing, using console notifier")
		notifier = notify.NewConsole()
	}

	// Create and run the engine
	engine := core.New(cfg, notifier, logger)
	runResult, err := engine.Run(ctx)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, fmt.Sprintf("engine error: %v", err))
	}

	// Build response
	result := HandlerResult{
		Success:    true,
		Message:    "Check cycle complete",
		Total:      runResult.Stats.TotalTickers,
		Triggered:  runResult.Stats.TriggeredCount,
		Sent:       runResult.Stats.NotificationsSent,
		Errors:     runResult.Stats.ErrorCount,
		DurationMs: time.Since(start).Milliseconds(),
		Results:    runResult.Results,
	}

	body, _ := json.Marshal(result)

	logger.Info("function completed",
		"duration_ms", result.DurationMs,
		"triggered", result.Triggered,
	)

	return Response{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}
}

func errorResponse(code int, message string) Response {
	result := HandlerResult{
		Success: false,
		Message: message,
	}
	body, _ := json.Marshal(result)

	return Response{
		StatusCode: code,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
