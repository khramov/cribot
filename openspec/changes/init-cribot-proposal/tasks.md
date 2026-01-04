## 1. Project Setup
- [x] 1.1 Initialize Go module (`go mod init github.com/antonkhramov/cribot`)
- [x] 1.2 Create directory structure (`cmd/`, `internal/`, `config/`, `deploy/`)
- [x] 1.3 Add `.gitignore` for Go projects
- [x] 1.4 Create sample `config/tickers.csv` with test data

## 2. Configuration System
- [x] 2.1 Define `TickerConfig` struct with all CSV fields
- [x] 2.2 Implement CSV parser in `internal/config/csv.go`
- [x] 2.3 Add validation for required fields and value ranges
- [x] 2.4 Write unit tests for CSV parsing

## 3. Plugin System
- [x] 3.1 Define `Source` interface in `internal/plugins/interface.go`
- [x] 3.2 Implement plugin registry with `Register()` and `Get()` functions
- [x] 3.3 Create stub `price` plugin (returns mock data)
- [x] 3.4 Create stub `rsi` plugin (returns mock data)
- [x] 3.5 Create stub `fx` plugin (returns mock data)
- [x] 3.6 Write unit tests for registry and interface compliance

## 4. Core Engine
- [x] 4.1 Implement `Engine` struct in `internal/core/engine.go`
- [x] 4.2 Add `Run()` method: load config → iterate tickers → call plugins → collect results
- [x] 4.3 Implement parallel plugin execution with `errgroup`
- [x] 4.4 Add filtering: skip disabled tickers, handle missing plugins gracefully
- [x] 4.5 Write unit tests for engine logic

## 5. Telegram Notifications
- [x] 5.1 Implement `Notifier` interface in `internal/notify/`
- [x] 5.2 Implement Telegram notifier using Bot API
- [x] 5.3 Format messages with ticker, plugin, triggered value
- [x] 5.4 Add error handling and retries
- [ ] 5.5 Write integration test with Telegram (manual verification)

## 6. Function Entry Point
- [x] 6.1 Create `cmd/function/main.go` with Yandex Cloud Function handler
- [x] 6.2 Wire together: config loading → engine → notifier
- [x] 6.3 Add environment variable support for secrets (TELEGRAM_TOKEN, CHAT_ID)
- [x] 6.4 Add structured logging

## 7. Deployment
- [x] 7.1 Create `deploy/yc-function.yaml` with function configuration
- [x] 7.2 Write deployment script using `yc` CLI
- [ ] 7.3 Set up Timer trigger (cron) in Yandex Cloud
- [x] 7.4 Document deployment steps in README

## 8. Real Data Integration (post-MVP validation)
- [ ] 8.1 Connect `price` plugin to real MOEX API
- [ ] 8.2 Connect `fx` plugin to currency rate API (CBR or similar)
- [ ] 8.3 Connect `rsi` plugin to market data source

---

**Dependencies:**
- Tasks 2.x must complete before 4.x (engine needs config)
- Tasks 3.x must complete before 4.x (engine needs plugins)
- Tasks 4.x and 5.x must complete before 6.x (entry point needs engine and notifier)
- Task 6.x must complete before 7.x (deployment needs working function)

**Parallelizable:**
- Tasks 2.x, 3.x, 5.x can be developed in parallel
- Task 8.x can happen after MVP validation
