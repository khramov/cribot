## MODIFIED Requirements

### Requirement: Execution Results
The Engine SHALL return detailed check results for each ticker, including value and error status, to support observability and testing.

#### Scenario: Return check details
- **WHEN** the `Engine.Run` method completes
- **THEN** it returns a `RunResult` containing aggregated statistics
- **AND** a list of `CheckResult` objects for all processed tickers
- **AND** each `CheckResult` contains the Ticker, Plugin used, Triggered status, current Value, and any Error
