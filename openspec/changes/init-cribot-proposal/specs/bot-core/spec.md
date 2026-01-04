## ADDED Requirements

### Requirement: CSV Configuration
The system SHALL load ticker configurations from a CSV file with the following columns: ticker, plugin, enabled, threshold_type, threshold_value, target_value, notes.

#### Scenario: Load valid configuration
- **WHEN** the function starts
- **THEN** it reads the CSV file from the configured path
- **AND** parses each row into a `TickerConfig` struct
- **AND** validates required fields (ticker, plugin, threshold_type)

#### Scenario: Handle invalid CSV row
- **WHEN** a CSV row has missing required fields
- **THEN** the system logs a warning
- **AND** skips the invalid row
- **AND** continues processing remaining rows

#### Scenario: Handle missing CSV file
- **WHEN** the CSV file does not exist
- **THEN** the system returns an error
- **AND** does not proceed with execution

---

### Requirement: Ticker Filtering
The system SHALL only process tickers where `enabled` is `true`.

#### Scenario: Skip disabled ticker
- **WHEN** a ticker has `enabled=false`
- **THEN** no plugin check is performed for that ticker
- **AND** no notification is sent

---

### Requirement: Telegram Notifications
The system SHALL send notifications to a configured Telegram chat when plugin conditions are triggered.

#### Scenario: Send triggered notification
- **WHEN** a plugin returns `Triggered=true`
- **THEN** the system sends a Telegram message containing:
  - Ticker symbol
  - Plugin name
  - Current value
  - Threshold that was crossed
- **AND** the message is formatted for readability

#### Scenario: Handle Telegram API error
- **WHEN** the Telegram API returns an error
- **THEN** the system logs the error
- **AND** continues processing remaining tickers
- **AND** does not crash the function

---

### Requirement: Environment Configuration
The system SHALL read sensitive configuration from environment variables.

#### Scenario: Load Telegram credentials
- **WHEN** the function starts
- **THEN** it reads `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID` from environment
- **AND** fails fast if either is missing

#### Scenario: Load config path
- **WHEN** the function starts
- **THEN** it reads `CONFIG_PATH` from environment
- **AND** defaults to `./config/tickers.csv` if not set
