# CriBot

Персональный serverless-бот для мониторинга финансовых идей и рыночных индикаторов.

## Возможности

- 📊 **Плагинная архитектура** — легко добавлять новые источники данных
- 📝 **CSV-конфигурация** — управление тикерами и порогами через Excel/Google Sheets
- 📱 **Telegram уведомления** — получай алерты только когда нужно действовать
- ☁️ **Serverless** — работает на Yandex Cloud Functions

## Быстрый старт

### 1. Настройка конфигурации

Отредактируй `config/tickers.csv`:

```csv
ticker,plugin,enabled,threshold_type,threshold_value,target_value,notes
SBER,price,true,below,250,300,Брать на просадке
USDRUB,fx,true,above,95,,Алерт на ослабление
VTBR,rsi,true,below,30,,Перепроданность
```

### 2. Создание Telegram бота

1. Напиши [@BotFather](https://t.me/BotFather) команду `/newbot`
2. Получи токен бота
3. Напиши своему боту, затем получи chat_id через `https://api.telegram.org/bot<TOKEN>/getUpdates`

### 3. Локальный запуск
 
**С уведомлениями в Telegram:**
```bash
export TELEGRAM_BOT_TOKEN="your-token"
export TELEGRAM_CHAT_ID="your-chat-id"
CGO_ENABLED=0 go run ./cmd/function
```

**Режим консоли (без Telegram):**
Если не указывать переменные окружения для Telegram, бот запустится в режиме консоли: он выведет все результаты проверок и алерты прямо в терминал.
```bash
CGO_ENABLED=0 go run ./cmd/function
```
*Примечание: `CGO_ENABLED=0` рекомендуется для стабильной работы на macOS.*


### 4. Деплой на Yandex Cloud

```bash
# Установи yc CLI: https://cloud.yandex.ru/docs/cli/quickstart
yc init

# Деплой
./deploy/deploy.sh

# Настрой Timer trigger (каждые 5 минут)
yc serverless trigger create timer \
    --name=cribot-timer \
    --cron-expression='0/5 * * * ? *' \
    --invoke-function-name=cribot
```

## Доступные плагины

| Plugin | Описание | Пример threshold |
|--------|----------|------------------|
| `price` | Цена актива | `below,250` — алерт когда цена < 250 |
| `rsi` | RSI индикатор | `below,30` — перепроданность |
| `fx` | Курсы валют | `above,95` — USDRUB > 95 |

## Структура CSV

| Колонка | Обязательно | Описание |
|---------|-------------|----------|
| `ticker` | ✅ | Тикер (SBER, USDRUB) |
| `plugin` | ✅ | Имя плагина |
| `enabled` | ✅ | true/false |
| `threshold_type` | ✅ | above/below |
| `threshold_value` | ✅ | Числовой порог |
| `target_value` | ❌ | Целевая цена (для заметок) |
| `notes` | ❌ | Комментарий |

## Переменные окружения

| Переменная | Обязательно | Описание |
|------------|-------------|----------|
| `TELEGRAM_BOT_TOKEN` | ✅ | Токен Telegram бота |
| `TELEGRAM_CHAT_ID` | ✅ | ID чата для уведомлений |
| `CONFIG_PATH` | ❌ | Путь к CSV (default: `./config/tickers.csv`) |
| `LOG_LEVEL` | ❌ | debug/info/warn/error (default: info) |

## Разработка

```bash
# Запуск unit-тестов
go test ./...

# Запуск интеграционных тестов (требует интернет)
go test -tags=integration -v ./tests/integration/...

# Сборка
go build -o cribot ./cmd/function
```

## Лицензия

MIT
