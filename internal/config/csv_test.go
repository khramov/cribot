package config

import (
	"strings"
	"testing"
)

func TestLoadFromReader_ValidCSV(t *testing.T) {
	csvData := `ticker,plugin,enabled,threshold_type,threshold_value,target_value,notes
SBER,price,true,below,250,300,Брать на просадке
USDRUB,fx,true,above,95,,Алерт на ослабление
VTBR,rsi,false,below,30,,Перепроданность`

	cfg, err := LoadFromReader(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Tickers) != 3 {
		t.Errorf("expected 3 tickers, got %d", len(cfg.Tickers))
	}

	// Check first ticker
	tc := cfg.Tickers[0]
	if tc.Ticker != "SBER" {
		t.Errorf("expected SBER, got %s", tc.Ticker)
	}
	if tc.Plugin != "price" {
		t.Errorf("expected price plugin, got %s", tc.Plugin)
	}
	if !tc.Enabled {
		t.Error("expected SBER to be enabled")
	}
	if tc.ThresholdType != ThresholdBelow {
		t.Errorf("expected below, got %s", tc.ThresholdType)
	}
	if tc.ThresholdValue != 250 {
		t.Errorf("expected 250, got %f", tc.ThresholdValue)
	}
	if tc.TargetValue != 300 {
		t.Errorf("expected 300, got %f", tc.TargetValue)
	}
}

func TestLoadFromReader_EnabledTickers(t *testing.T) {
	csvData := `ticker,plugin,enabled,threshold_type,threshold_value,target_value,notes
SBER,price,true,below,250,300,
VTBR,rsi,false,below,30,,`

	cfg, err := LoadFromReader(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	enabled := cfg.EnabledTickers()
	if len(enabled) != 1 {
		t.Errorf("expected 1 enabled ticker, got %d", len(enabled))
	}
	if enabled[0].Ticker != "SBER" {
		t.Errorf("expected SBER, got %s", enabled[0].Ticker)
	}
}

func TestLoadFromReader_InvalidRow(t *testing.T) {
	csvData := `ticker,plugin,enabled,threshold_type,threshold_value,target_value,notes
SBER,price,true,below,250,300,
,price,true,below,250,,Missing ticker`

	cfg, err := LoadFromReader(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Invalid row should be skipped
	if len(cfg.Tickers) != 1 {
		t.Errorf("expected 1 valid ticker, got %d", len(cfg.Tickers))
	}
}

func TestLoadFromReader_MissingColumn(t *testing.T) {
	csvData := `ticker,plugin,enabled
SBER,price,true`

	_, err := LoadFromReader(strings.NewReader(csvData))
	if err == nil {
		t.Error("expected error for missing columns")
	}
}

func TestTickerConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  TickerConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: TickerConfig{
				Ticker:        "SBER",
				Plugin:        "price",
				ThresholdType: ThresholdBelow,
			},
			wantErr: false,
		},
		{
			name: "missing ticker",
			config: TickerConfig{
				Plugin:        "price",
				ThresholdType: ThresholdBelow,
			},
			wantErr: true,
		},
		{
			name: "missing plugin",
			config: TickerConfig{
				Ticker:        "SBER",
				ThresholdType: ThresholdBelow,
			},
			wantErr: true,
		},
		{
			name: "invalid threshold type",
			config: TickerConfig{
				Ticker:        "SBER",
				Plugin:        "price",
				ThresholdType: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
