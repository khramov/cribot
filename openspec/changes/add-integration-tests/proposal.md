# Change: Add Integration Tests Infrastructure

## Why
Пользователю требуется способ верифицировать работу бота на реальных данных и сценариях. Юнит-тесты проверяют компоненты изолированно, но необходим сквозной тест, который:
1. Загрузит тестовую конфигурацию (CSV).
2. Выполнит реальный цикл проверки (включая запросы к API).
3. Проверит, что результаты (триггеры, цены) корректны.

## What Changes
- Создать директорию `tests/integration/` для тестов.
- Создать тестовый конфиг `tests/integration/test_config.csv`.
- Реализовать `TestIntegration` в `tests/integration/main_test.go`, который запускает `Engine` с реальными плагинами.
- Добавить флаг `-integration` для запуска этих тестов (чтобы они не бежали в unit-сборке).

## Impact
- **Affected specs**: `testing` (ADDED capability).
- **Affected code**: Новая директория `tests/`.
