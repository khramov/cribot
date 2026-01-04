package moex

// ISSResponse represents the generic JSON response structure from MOEX ISS.
// It uses a generic map because column names are provided in the "columns" field,
// and data in the "data" field (array of arrays).
type ISSResponse map[string]TableData

// TableData holds the columns and data for a specific section (e.g., "marketdata").
type TableData struct {
	Columns []string        `json:"columns"`
	Data    [][]interface{} `json:"data"`
}

// MarketData represents parsed market data.
type MarketData struct {
	Last       float64
	Change     float64
	ChangePerc float64
	UpdateTime string
}
