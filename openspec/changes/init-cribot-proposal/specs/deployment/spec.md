## ADDED Requirements

### Requirement: Yandex Cloud Function Compatibility
The system SHALL be deployable as a Yandex Cloud Function with HTTP or Timer trigger.

#### Scenario: Timer trigger execution
- **WHEN** the Timer trigger fires (e.g., every 5 minutes)
- **THEN** the function executes the full check cycle
- **AND** completes within the 10-second timeout
- **AND** returns HTTP 200 on success

#### Scenario: HTTP trigger execution
- **WHEN** an HTTP request is received
- **THEN** the function executes the full check cycle
- **AND** returns HTTP 200 with execution summary on success
- **AND** returns HTTP 500 with error details on failure

---

### Requirement: Environment Variables
The system SHALL use environment variables for all deployment-specific configuration.

#### Scenario: Required variables
- **WHEN** the function starts
- **THEN** it reads the following environment variables:
  - `TELEGRAM_BOT_TOKEN` — Telegram Bot API token (required)
  - `TELEGRAM_CHAT_ID` — Target chat ID for notifications (required)
  - `CONFIG_PATH` — Path to CSV config (optional, defaults to `./config/tickers.csv`)
  - `LOG_LEVEL` — Logging verbosity (optional, defaults to `info`)

---

### Requirement: Structured Logging
The system SHALL use structured logging compatible with Yandex Cloud logging.

#### Scenario: Log execution summary
- **WHEN** the function completes
- **THEN** it logs:
  - Total tickers processed
  - Number of triggered conditions
  - Number of notifications sent
  - Execution duration in milliseconds

#### Scenario: Log errors with context
- **WHEN** an error occurs
- **THEN** the log entry includes:
  - Error message
  - Affected ticker (if applicable)
  - Plugin name (if applicable)
  - Stack trace (for critical errors)

---

### Requirement: Single Binary Deployment
The system SHALL compile into a single binary for simplified serverless deployment.

#### Scenario: Build for deployment
- **WHEN** running `go build -o function ./cmd/function`
- **THEN** a single executable is produced
- **AND** it contains all plugins and dependencies
- **AND** it is suitable for upload to Yandex Cloud Functions
