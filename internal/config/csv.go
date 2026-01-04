// Package config handles loading and parsing of ticker configurations.
package config

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// ThresholdType defines the comparison type for triggers.
type ThresholdType string

const (
	ThresholdBelow ThresholdType = "below"
	ThresholdAbove ThresholdType = "above"
)

// TickerConfig represents a single ticker configuration from CSV.
type TickerConfig struct {
	Ticker         string        // Ticker symbol (e.g., SBER, USDRUB)
	Plugin         string        // Plugin name to use (e.g., price, rsi, fx)
	Enabled        bool          // Whether this ticker is active
	ThresholdType  ThresholdType // Comparison type: above or below
	ThresholdValue float64       // Threshold value to compare against
	TargetValue    float64       // Target price (optional, for notes)
	Notes          string        // User notes
}

// Config holds the entire configuration.
type Config struct {
	Tickers []TickerConfig
}

// Validate checks that the TickerConfig has all required fields.
func (tc *TickerConfig) Validate() error {
	if tc.Ticker == "" {
		return errors.New("ticker is required")
	}
	if tc.Plugin == "" {
		return errors.New("plugin is required")
	}
	if tc.ThresholdType != ThresholdBelow && tc.ThresholdType != ThresholdAbove {
		return fmt.Errorf("threshold_type must be 'above' or 'below', got '%s'", tc.ThresholdType)
	}
	return nil
}

// LoadFromFile reads a CSV file and parses it into a Config.
func LoadFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	return LoadFromReader(file)
}

// LoadFromReader parses CSV from an io.Reader.
func LoadFromReader(r io.Reader) (*Config, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Map column names to indices
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[strings.ToLower(strings.TrimSpace(col))] = i
	}

	// Verify required columns
	requiredCols := []string{"ticker", "plugin", "enabled", "threshold_type", "threshold_value"}
	for _, col := range requiredCols {
		if _, ok := colIndex[col]; !ok {
			return nil, fmt.Errorf("missing required column: %s", col)
		}
	}

	var tickers []TickerConfig
	lineNum := 1 // header is line 1

	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV line %d: %w", lineNum, err)
		}

		tc, err := parseRecord(record, colIndex, lineNum)
		if err != nil {
			// Log warning and skip invalid row
			fmt.Printf("warning: skipping line %d: %v\n", lineNum, err)
			continue
		}

		tickers = append(tickers, tc)
	}

	return &Config{Tickers: tickers}, nil
}

// parseRecord converts a CSV record to TickerConfig.
func parseRecord(record []string, colIndex map[string]int, lineNum int) (TickerConfig, error) {
	getCol := func(name string) string {
		idx, ok := colIndex[name]
		if !ok || idx >= len(record) {
			return ""
		}
		return strings.TrimSpace(record[idx])
	}

	tc := TickerConfig{
		Ticker:        getCol("ticker"),
		Plugin:        strings.ToLower(getCol("plugin")),
		ThresholdType: ThresholdType(strings.ToLower(getCol("threshold_type"))),
		Notes:         getCol("notes"),
	}

	// Parse enabled
	enabledStr := strings.ToLower(getCol("enabled"))
	tc.Enabled = enabledStr == "true" || enabledStr == "1" || enabledStr == "yes"

	// Parse threshold_value
	thresholdStr := getCol("threshold_value")
	if thresholdStr != "" {
		val, err := strconv.ParseFloat(thresholdStr, 64)
		if err != nil {
			return tc, fmt.Errorf("invalid threshold_value '%s': %w", thresholdStr, err)
		}
		tc.ThresholdValue = val
	}

	// Parse target_value (optional)
	targetStr := getCol("target_value")
	if targetStr != "" {
		val, err := strconv.ParseFloat(targetStr, 64)
		if err != nil {
			return tc, fmt.Errorf("invalid target_value '%s': %w", targetStr, err)
		}
		tc.TargetValue = val
	}

	// Validate
	if err := tc.Validate(); err != nil {
		return tc, err
	}

	return tc, nil
}

// EnabledTickers returns only tickers that are enabled.
func (c *Config) EnabledTickers() []TickerConfig {
	var enabled []TickerConfig
	for _, tc := range c.Tickers {
		if tc.Enabled {
			enabled = append(enabled, tc)
		}
	}
	return enabled
}
