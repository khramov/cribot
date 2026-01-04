package moex

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetStockPrice(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL format
		if r.URL.Path != "/engines/stock/markets/shares/boards/TQBR/securities/SBER.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.Error(w, "should not happen", http.StatusBadRequest)
			return
		}

		resp := `{
			"marketdata": {
				"columns": ["SECID", "LAST", "CHANGE"],
				"data": [
					["SBER", 255.50, 1.5]
				]
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	defer ts.Close()

	client := NewClient(time.Minute)
	client.SetBaseURL(ts.URL)

	data, err := client.GetStockPrice(context.Background(), "SBER")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Last != 255.50 {
		t.Errorf("expected price 255.50, got %f", data.Last)
	}
	if data.Change != 1.5 {
		t.Errorf("expected change 1.5, got %f", data.Change)
	}
}

func TestGetCurrencyRate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/engines/currency/markets/selt/boards/CETS/securities/USDRUB_TOM.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			return
		}

		resp := `{
			"marketdata": {
				"columns": ["SECID", "LAST"],
				"data": [
					["USDRUB_TOM", 92.45]
				]
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	defer ts.Close()

	client := NewClient(time.Minute)
	client.SetBaseURL(ts.URL)

	data, err := client.GetCurrencyRate(context.Background(), "USDRUB_TOM")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Last != 92.45 {
		t.Errorf("expected rate 92.45, got %f", data.Last)
	}
}

func TestCaching(t *testing.T) {
	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := `{
			"marketdata": {
				"columns": ["LAST"],
				"data": [[100.0]]
			}
		}`
		w.Write([]byte(resp))
	}))
	defer ts.Close()

	client := NewClient(time.Minute)
	client.SetBaseURL(ts.URL)

	// First call - should hit server
	_, err := client.GetStockPrice(context.Background(), "TEST")
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// Second call - should hit cache
	_, err = client.GetStockPrice(context.Background(), "TEST")
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (cache hit), got %d", callCount)
	}
}
