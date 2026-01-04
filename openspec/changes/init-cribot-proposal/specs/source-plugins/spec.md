## ADDED Requirements

### Requirement: Source Interface
The system SHALL define a common interface for all data source plugins to ensure uniform behavior.

#### Scenario: Interface contract
- **WHEN** a new plugin is implemented
- **THEN** it MUST implement the `Source` interface with methods:
  - `Name() string` — returns plugin identifier
  - `Check(ctx, ticker, config) (*Result, error)` — checks condition and returns result

#### Scenario: Result structure
- **WHEN** a plugin check completes
- **THEN** the `Result` struct contains:
  - `Triggered bool` — whether the condition was met
  - `Message string` — human-readable description
  - `CurrentValue float64` — the actual value checked

---

### Requirement: Plugin Registry
The system SHALL maintain a registry of available plugins that can be accessed by name.

#### Scenario: Register plugin
- **WHEN** a plugin package is imported
- **THEN** it registers itself with the global registry via `init()`
- **AND** becomes available by its name

#### Scenario: Get plugin by name
- **WHEN** the engine requests a plugin by name
- **THEN** the registry returns the plugin instance
- **OR** returns `nil` if the plugin is not registered

---

### Requirement: Price Plugin
The system SHALL include a plugin to check asset price against thresholds.

#### Scenario: Price below threshold
- **WHEN** configured with `threshold_type=below` and `threshold_value=250`
- **AND** the current price is 240
- **THEN** the plugin returns `Triggered=true`
- **AND** `Message` includes "цена 240 ниже порога 250"

#### Scenario: Price above threshold
- **WHEN** configured with `threshold_type=above` and `threshold_value=300`
- **AND** the current price is 310
- **THEN** the plugin returns `Triggered=true`
- **AND** `Message` includes "цена 310 выше порога 300"

#### Scenario: Price within threshold
- **WHEN** configured with `threshold_type=below` and `threshold_value=250`
- **AND** the current price is 260
- **THEN** the plugin returns `Triggered=false`

---

### Requirement: RSI Plugin
The system SHALL include a plugin to check RSI (Relative Strength Index) values.

#### Scenario: RSI oversold
- **WHEN** configured with `threshold_type=below` and `threshold_value=30`
- **AND** the current RSI is 25
- **THEN** the plugin returns `Triggered=true`
- **AND** `Message` includes "RSI 25 (перепроданность)"

#### Scenario: RSI overbought
- **WHEN** configured with `threshold_type=above` and `threshold_value=70`
- **AND** the current RSI is 78
- **THEN** the plugin returns `Triggered=true`
- **AND** `Message` includes "RSI 78 (перекупленность)"

---

### Requirement: FX Plugin
The system SHALL include a plugin to check currency exchange rates.

#### Scenario: Rate crosses threshold
- **WHEN** tracking USDRUB with `threshold_type=above` and `threshold_value=95`
- **AND** the current rate is 97.5
- **THEN** the plugin returns `Triggered=true`
- **AND** `Message` includes "курс 97.5 выше порога 95"

---

### Requirement: Graceful Plugin Errors
The system SHALL handle plugin errors gracefully without crashing.

#### Scenario: Plugin returns error
- **WHEN** a plugin fails to fetch data (network error, API limit, etc.)
- **THEN** the engine logs the error with context
- **AND** continues processing other tickers
- **AND** does not send a notification for this ticker
