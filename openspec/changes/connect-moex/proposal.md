# Change: Connect MOEX ISS API for Real Market Data

## Why
Текущие плагины (price, fx) возвращают mock-данные. Для реальной работы бота необходимо подключиться к MOEX ISS API для получения актуальных котировок и курсов валют.

**MOEX ISS API:**
- Бесплатный доступ к delayed data (15 мин задержка)
- JSON формат ответов
- Не требует авторизации для базовых запросов

## What Changes

- Создать пакет `internal/moex` — HTTP-клиент для MOEX ISS API
- Модифицировать `price` плагин — использовать реальные данные с MOEX
- Модифицировать `fx` плагин — получать курсы валют через MOEX
- Добавить кеширование для снижения нагрузки на API
- Добавить graceful fallback на mock-данные при ошибках

## Impact

- **Affected specs**: `source-plugins` (MODIFIED: Price Plugin, FX Plugin)
- **Affected code**:
  - `internal/moex/` — новый пакет
  - `internal/plugins/price/price.go` — модификация
  - `internal/plugins/fx/fx.go` — модификация
