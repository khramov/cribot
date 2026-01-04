## 1. MOEX Client Package
- [x] 1.1 Create `internal/moex/types.go` with response structs
- [x] 1.2 Create `internal/moex/client.go` with HTTP client and caching
- [x] 1.3 Implement `GetStockPrice(ticker string)` method
- [x] 1.4 Implement `GetCurrencyRate(pair string)` method
- [x] 1.5 Write unit tests for client (with mocked HTTP)

## 2. Update Price Plugin
- [x] 2.1 Inject MOEX client into price plugin
- [x] 2.2 Replace `mockPrice()` with real API call
- [x] 2.3 Add fallback to cached/mock data on error
- [x] 2.4 Update unit tests

## 3. Update FX Plugin
- [x] 3.1 Inject MOEX client into fx plugin
- [x] 3.2 Replace `mockFXRate()` with real API call
- [x] 3.3 Map ticker format (USDRUB → USDRUB_TOM)
- [x] 3.4 Add fallback to cached/mock data on error
- [x] 3.5 Update unit tests

## 4. Integration
- [x] 4.1 Update `cmd/function/main.go` to initialize MOEX client
- [x] 4.2 Add optional `MOEX_CACHE_TTL` environment variable
- [x] 4.3 Manual integration test with real API

---

**Dependencies:**
- Task 1.x must complete before 2.x and 3.x
- Tasks 2.x and 3.x can run in parallel
- Task 4.x requires 2.x and 3.x complete

**Verification:**
- Unit tests: `go test ./internal/moex/...`
- Integration: Run locally with real MOEX API and verify prices match moex.com
