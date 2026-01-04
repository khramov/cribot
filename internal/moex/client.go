package moex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	defaultBaseURL = "https://iss.moex.com/iss"
	defaultTTL     = 5 * time.Minute
)

// Client interface for mocking
type Client interface {
	GetStockPrice(ctx context.Context, ticker string) (*MarketData, error)
	GetCurrencyRate(ctx context.Context, pair string) (*MarketData, error)
	SetBaseURL(url string)
}

// ISSClient handles requests to MOEX ISS API.
type ISSClient struct {
	httpClient *http.Client
	cache      sync.Map
	ttl        time.Duration
	baseURL    string
}

type cachedItem struct {
	data      *MarketData
	timestamp time.Time
}

// SetBaseURL allows overriding the base API URL (useful for testing).
func (c *ISSClient) SetBaseURL(url string) {
	c.baseURL = url
}

// NewClient creates a new MOEX client.
func NewClient(ttl time.Duration) *ISSClient {
	if ttl <= 0 {
		ttl = defaultTTL
	}
	return &ISSClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		ttl:        ttl,
		baseURL:    defaultBaseURL,
	}
}

// GetStockPrice fetches current price for a stock on TQBR board.
func (c *ISSClient) GetStockPrice(ctx context.Context, ticker string) (*MarketData, error) {
	cacheKey := "stock:" + ticker
	if item, ok := c.getFromCache(cacheKey); ok {
		return item, nil
	}

	url := fmt.Sprintf("%s/engines/stock/markets/shares/boards/TQBR/securities/%s.json?iss.meta=off", c.baseURL, ticker)
	data, err := c.fetchMarketData(ctx, url)
	if err != nil {
		return nil, err
	}

	c.saveToCache(cacheKey, data)
	return data, nil
}

// GetCurrencyRate fetches exchange rate for a pair on CETS market (e.g., USDRUB_TOM).
func (c *ISSClient) GetCurrencyRate(ctx context.Context, pair string) (*MarketData, error) {
	cacheKey := "fx:" + pair
	if item, ok := c.getFromCache(cacheKey); ok {
		return item, nil
	}

	url := fmt.Sprintf("%s/engines/currency/markets/selt/boards/CETS/securities/%s.json?iss.meta=off", c.baseURL, pair)
	data, err := c.fetchMarketData(ctx, url)
	if err != nil {
		return nil, err
	}

	c.saveToCache(cacheKey, data)
	return data, nil
}

func (c *ISSClient) getFromCache(key string) (*MarketData, bool) {
	val, ok := c.cache.Load(key)
	if !ok {
		return nil, false
	}
	item := val.(cachedItem)
	if time.Since(item.timestamp) > c.ttl {
		c.cache.Delete(key)
		return nil, false
	}
	return item.data, true
}

func (c *ISSClient) saveToCache(key string, data *MarketData) {
	c.cache.Store(key, cachedItem{
		data:      data,
		timestamp: time.Now(),
	})
}

func (c *ISSClient) fetchMarketData(ctx context.Context, url string) (*MarketData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	var issResp ISSResponse
	if err := json.Unmarshal(body, &issResp); err != nil {
		return nil, fmt.Errorf("parse json failed: %w", err)
	}

	return parseFirstRow(issResp, "marketdata")
}

func parseFirstRow(resp ISSResponse, tableName string) (*MarketData, error) {
	table, ok := resp[tableName]
	if !ok {
		return nil, fmt.Errorf("table '%s' not found in response", tableName)
	}

	if len(table.Data) == 0 {
		return nil, fmt.Errorf("no data found for '%s'", tableName)
	}

	// Find column indices
	colMap := make(map[string]int)
	for i, name := range table.Columns {
		colMap[name] = i
	}

	row := table.Data[0]
	
	val, err := getFloat(row, colMap, "LAST")
	if err != nil {
		// Try LASTVALUE if LAST is missing (sometimes happens)
		val, err = getFloat(row, colMap, "LASTVALUE")
		if err != nil {
			return nil, fmt.Errorf("LAST price not found")
		}
	}

	change, _ := getFloat(row, colMap, "CHANGE") // optional
	
	// Try multiple fields for change percentage
	changePerc, err := getFloat(row, colMap, "LASTTOPREVPRICE") 
	if err != nil {
		changePerc, _ = getFloat(row, colMap, "CHANGEPCT")
	}

	return &MarketData{
		Last:       val,
		Change:     change,
		ChangePerc: changePerc,
		UpdateTime: time.Now().Format(time.RFC3339),
	}, nil
}

func getFloat(row []interface{}, colMap map[string]int, colName string) (float64, error) {
	idx, ok := colMap[colName]
	if !ok || idx >= len(row) {
		return 0, fmt.Errorf("col %s not found", colName)
	}

	val := row[idx]
	if val == nil {
		return 0, fmt.Errorf("col %s is nil", colName)
	}

	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("col %s is not a number: %T", colName, v)
	}
}
