## MODIFIED Requirements

### Requirement: Price Plugin
The bot SHALL support checking asset prices from MOEX using the ISS API with caching and fallback.

#### Scenario: Fetch real price from MOEX
- **WHEN** configured to track a stock ticker (e.g., SBER)
- **THEN** the plugin fetches current price from MOEX ISS API
- **AND** returns the LAST price from marketdata

#### Scenario: Cache hit
- **WHEN** the same ticker is requested within cache TTL
- **THEN** the plugin returns cached price
- **AND** does not make an API request

#### Scenario: API error fallback
- **WHEN** MOEX API returns an error or times out
- **THEN** the plugin returns cached data if available
- **OR** returns an error if no cache exists
- **AND** logs the error with context

---

### Requirement: FX Plugin
The bot SHALL support fetching currency exchange rates from MOEX CETS market.

#### Scenario: Fetch real FX rate
- **WHEN** configured to track a currency pair (e.g., USDRUB)
- **THEN** the plugin maps it to MOEX format (USDRUB_TOM)
- **AND** fetches the rate from MOEX ISS API
- **AND** returns the LAST rate

#### Scenario: Unsupported pair
- **WHEN** configured with an unknown currency pair
- **THEN** the plugin returns an error
- **AND** logs the unsupported pair

---

## ADDED Requirements

### Requirement: MOEX Client
The system SHALL provide a reusable HTTP client for MOEX ISS API with caching.

#### Scenario: Stock price request
- **WHEN** `GetStockPrice("SBER")` is called
- **THEN** the client makes a request to MOEX ISS API
- **AND** parses the JSON response
- **AND** returns the current price

#### Scenario: Currency rate request
- **WHEN** `GetCurrencyRate("USDRUB_TOM")` is called
- **THEN** the client makes a request to MOEX CETS market
- **AND** returns the exchange rate

#### Scenario: Caching
- **WHEN** a request is made for a ticker
- **THEN** the response is cached with TTL
- **AND** subsequent requests within TTL return cached data
