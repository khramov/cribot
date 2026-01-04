## Context
MOEX ISS (Informational & Statistical Server) предоставляет бесплатный REST API для получения рыночных данных. Delayed data (15 мин) доступна без авторизации.

## Goals / Non-Goals

### Goals
- ✅ Получение реальных цен акций с MOEX
- ✅ Получение курсов валют (USD/RUB, EUR/RUB)
- ✅ Кеширование для минимизации запросов
- ✅ Graceful degradation при ошибках API

### Non-Goals
- ❌ Real-time данные (требуют подписки)
- ❌ Исторические данные / свечи
- ❌ Подключение к другим биржам

## API Endpoints

### Акции (TQBR board)
```
GET https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities/{ticker}.json?iss.meta=off
```

Response содержит:
- `marketdata` — текущие данные (LAST, BID, OFFER, CHANGE)
- `securities` — справочная информация

### Валюта (CETS market)
```
GET https://iss.moex.com/iss/engines/currency/markets/selt/boards/CETS/securities/{pair}.json?iss.meta=off
```

Пары: USDRUB_TOM, EURRUB_TOM, CNYRUB_TOM

## Decisions

### Decision 1: HTTP Client в отдельном пакете
**Rationale:** Изоляция API-логики, легче тестировать и мокать.

### Decision 2: In-memory кеш с TTL
**Rationale:**
- Serverless function может вызываться часто
- Данные с задержкой 15 мин — можно кешировать на 1-5 мин
- Простой sync.Map с timestamp

### Decision 3: Fallback на mock при ошибках
**Rationale:**
- Не ломать весь цикл при временных проблемах API
- Логировать ошибку, но продолжать работу
- В Telegram сообщать что данные устаревшие

## File Structure

```
internal/
└── moex/
    ├── client.go      # HTTP client, caching
    ├── client_test.go # Unit tests
    └── types.go       # Response types
```
