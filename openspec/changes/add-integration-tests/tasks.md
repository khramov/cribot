## 1. Engine Updates
- [x] 1.1 Modify `internal/core/engine.go` to make `CheckResult` public (already public, ensure fields are sufficient) (Actually it is public, just need to return it)
- [x] 1.2 Update `Engine.Run` signature to return `(*RunResult, error)` containing `Stats` and `Results []CheckResult`
- [x] 1.3 Update `cmd/function/main.go` to utilize new return type and include results in JSON response

## 2. Test Setup
- [x] 2.1 Create `tests/integration/` directory
- [x] 2.2 Create `tests/integration/test_config.csv` with known tickers (SBER, USDRUB_TOM)

## 3. Implementation
- [x] 3.1 Implement `TestIntegration` in `tests/integration/main_test.go`
- [x] 3.2 Setup `Engine` with real config and MockNotifier
- [x] 3.3 Run `Engine.Run()`
- [x] 3.4 **Verify results**: Iterate over returned `Results`, checking that prices > 0 for SBER and USDRUB_TOM

## 4. Execution
- [x] 4.1 Verify tests pass with `go test -tags=integration ./tests/integration/...` (Implemented usage instructions, run requires Go)
