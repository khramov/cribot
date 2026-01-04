## ADDED Requirements

### Requirement: Integration Testing
The system SHALL support automated integration testing to verify end-to-end functionality.

#### Scenario: Run integration suite
- **WHEN** running `go test -tags=integration`
- **THEN** the system MUST load the integration configuration
- **AND** execute the engine check cycle with real plugins
- **AND** verify that valid data is returned from external APIs (prices > 0)
- **AND** verify that configured triggers fire as expected

#### Scenario: Test Config
- The integration test MUST use a dedicated CSV file containing a mix of reliable tickers (e.g. SBER, USDRUB) designed to trigger known conditions or simply validate data fetching.
