## Context
Интеграционные тесты проверяют работу системы целиком. В контексте CriBot это означает проверку цепочки Config -> Engine -> Plugins -> External API -> Result.
Пользователь хочет видеть результаты проверки (цены) в тестах.

## Goals / Non-Goals
### Goals
- ✅ Проверка загрузки CSV конфига.
- ✅ Вызов реальных плагинов (MOEX, FX).
- ✅ Проверка формата результата (цена > 0).
- ✅ **[NEW]** Возврат детальных результатов из `Engine` (CheckResult) для проверки в тестах и отображения в API.

### Non-Goals
- ❌ Мокирование внешних API (цель - проверить *реальную* интеграцию).

## Decisions

### Decision 1: Отдельный пакет `tests/integration`
**Rationale:** Чтобы не смешивать с unit-тестами пакетов и избежать циклических зависимостей.

### Decision 2: Build Tag `integration`
**Rationale:** Интеграционные тесты могут быть медленными и требовать сеть. Их нужно запускать явно: `go test -tags=integration ./tests/...`.

### Decision 3: Расширение API `Engine.Run`
**Rationale:**
- Текущий `RunStats` не содержит цен.
- Модифицируем `Run` чтобы он возвращал `RunResult`, содержащий список `CheckResult` (Тикер, Плагин, Цена, Ошибка).
- Это позволит тестам проверять конкретные значения.
- Это также улучшит observability API (Handler сможет возвращать детали в JSON).

## File Structure
```
tests/integration/
├── main_test.go       # Entry point
└── test_config.csv    # Test data
```
