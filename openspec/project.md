# Project Context

## Purpose
CriBot — персональный бот для отслеживания финансовых идей и рыночных индикаторов. Автоматически проверяет условия по заданным тикерам и отправляет уведомления в Telegram при срабатывании триггеров. Построен на плагинной архитектуре для лёгкого добавления новых источников данных.

## Tech Stack
- **Language**: Go 1.21+
- **Deployment**: Yandex Cloud Functions (serverless)
- **Configuration**: CSV files
- **Notifications**: Telegram Bot API
- **Architecture**: Plugin-based, stateless functions

## Project Conventions

### Code Style
- Follow standard Go conventions (gofmt, golint)
- Use meaningful package names (e.g., `plugins`, `config`, `notify`)
- Error handling: wrap errors with context using `fmt.Errorf`

### Architecture Patterns
- **Plugin Interface**: All data sources implement a common `Source` interface
- **Stateless Functions**: Each invocation reads config, checks conditions, sends notifications
- **Configuration-Driven**: Behavior controlled via CSV, not hardcoded

### Testing Strategy
- Unit tests for core logic (config parsing, condition evaluation)
- Integration tests with mocked external APIs
- Manual testing with real Telegram bot in dev environment

### Git Workflow
- Main branch: `main`
- Feature branches: `feature/<name>`
- Commit messages: conventional commits (feat:, fix:, docs:)

## Domain Context
- **Тикер**: Идентификатор актива (SBER, VTBR, USDRUB и т.д.)
- **Плагин**: Модуль для получения данных (RSI, курсы валют, и т.д.)
- **Триггер**: Условие срабатывания (цена достигла X, RSI < 30)
- **Порог (threshold)**: Числовое значение для сравнения

## Important Constraints
- Serverless execution time limits (~10 sec for Yandex Cloud Functions)
- Single user system (no multi-tenancy)
- No persistent database in MVP (CSV-based config)
- OpenAI API access via European proxy

## External Dependencies
- Telegram Bot API (notifications)
- Market data APIs (MOEX, etc.) — specific providers TBD
- Yandex Cloud Functions runtime
