## Context
CriBot — персональный serverless-бот для мониторинга финансовых идей. Пользователь ведёт CSV-таблицу с тикерами и условиями, бот периодически проверяет условия через плагины и шлёт уведомления в Telegram.

**Constraints:**
- Yandex Cloud Functions (10 sec timeout)
- Single user, no auth needed
- Минимальная инфраструктура (no DB in MVP)
- Go для хорошей поддержки serverless и concurrency

## Goals / Non-Goals

### Goals
- ✅ Расширяемая плагинная система для источников данных
- ✅ CSV-конфигурация для управления тикерами и порогами
- ✅ Serverless деплой с минимальным maintenance
- ✅ Telegram-уведомления при срабатывании условий
- ✅ Возможность включать/выключать плагины per-ticker

### Non-Goals
- ❌ Веб-интерфейс (Phase 2)
- ❌ AI-интеграция для парсинга идей (Phase 2)
- ❌ Telegram-бот для ввода идей (Phase 2)
- ❌ Автоматическое исполнение сделок
- ❌ Multi-user / аутентификация

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                  Yandex Cloud Function                  │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │   Config    │───▶│    Core     │───▶│  Telegram   │ │
│  │  (CSV file) │    │   Engine    │    │   Notifier  │ │
│  └─────────────┘    └──────┬──────┘    └─────────────┘ │
│                            │                            │
│         ┌──────────────────┼──────────────────┐        │
│         ▼                  ▼                  ▼        │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │ RSI Plugin  │    │ Price Plugin│    │ FX Plugin   │ │
│  └─────────────┘    └─────────────┘    └─────────────┘ │
└─────────────────────────────────────────────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │  External APIs  │
                    │  (MOEX, etc.)   │
                    └─────────────────┘
```

## Decisions

### Decision 1: Go для реализации
**Rationale:** 
- Отличная поддержка serverless (быстрый cold start)
- Встроенная concurrency для параллельных запросов к API
- Статическая типизация для надёжности
- Компиляция в единый бинарник

### Decision 2: Плагины как Go interfaces (compile-time)
**Rationale:**
- Serverless не подходит для dynamic plugin loading (.so files)
- Все плагины компилируются в один бинарник
- Простота деплоя — один артефакт
- Добавление плагина = добавление пакета + регистрация

**Interface:**
```go
type Source interface {
    Name() string
    Check(ctx context.Context, ticker string, config TickerConfig) (*Result, error)
}

type Result struct {
    Triggered   bool
    Message     string
    CurrentValue float64
}
```

### Decision 3: CSV для конфигурации
**Rationale:**
- Легко редактировать в Excel/Google Sheets
- Можно хранить в Object Storage или прямо в репозитории
- Версионируется через Git
- Достаточно для single-user MVP

**CSV Structure:**
```csv
ticker,plugin,enabled,threshold_type,threshold_value,target_value,notes
SBER,price,true,below,250,300,Брать на просадке
USDRUB,fx,true,above,95,,Алерт на ослабление
VTBR,rsi,true,below,30,,Перепроданность
```

### Decision 4: Cron-triggered execution
**Rationale:**
- Большинство данных (котировки, RSI) не push-based
- Yandex Cloud Functions поддерживает Timer triggers
- Configurable frequency (каждые 5 мин, каждый час)
- Stateless — каждый запуск независим

## Risks / Trade-offs

### Risk 1: Rate limits внешних API
**Mitigation:** 
- Кеширование в Object Storage между вызовами
- Батчинг запросов где возможно
- Graceful degradation при ошибках

### Risk 2: Cold start latency
**Mitigation:**
- Go имеет быстрый cold start (~100ms)
- Минимизация зависимостей
- Provisioned concurrency если критично (платно)

### Risk 3: CSV corruption / sync issues
**Mitigation:**
- Валидация при загрузке
- Git для версионирования
- Backup перед изменениями

## File Structure

```
cribot/
├── cmd/
│   └── function/
│       └── main.go          # Yandex Cloud Function entry point
├── internal/
│   ├── config/
│   │   └── csv.go           # CSV parsing
│   ├── core/
│   │   └── engine.go        # Main orchestration logic
│   ├── notify/
│   │   └── telegram.go      # Telegram notification
│   └── plugins/
│       ├── interface.go     # Source interface
│       ├── registry.go      # Plugin registration
│       ├── price/
│       │   └── price.go     # Price check plugin
│       ├── rsi/
│       │   └── rsi.go       # RSI indicator plugin
│       └── fx/
│           └── fx.go        # Currency rates plugin
├── config/
│   └── tickers.csv          # Sample configuration
├── go.mod
├── go.sum
└── deploy/
    └── yc-function.yaml     # Yandex Cloud deployment config
```

## Open Questions
- [ ] Какой API использовать для данных MOEX? (tinkoff-invest-api, moex-iss, etc.)
- [ ] Где хранить CSV в проде? (Object Storage vs Git repo)
- [ ] Нужен ли state между вызовами? (last notified time, cooldowns)
