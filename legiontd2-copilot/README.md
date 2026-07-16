# LT2 Copilot — интеллектуальный помощник для Legion TD 2

Полное ТЗ: см. `TZ_LegionTD2_Assistant.md` в корне проекта (выдан отдельно от этого репозитория).

## Статус

Это скелет проекта после Phase 0 (валидация источников данных). Реализации ещё нет —
только структура, контракты и TODO с привязкой к фазам ТЗ.

**Закрыто в Phase 0:**
- официальный API v2 подтверждён, схема разобрана (`internal/api/client.go` содержит найденные детали);
- клиентские логи игры — отсутствуют, исключены из архитектуры;
- EULA/ToS — прочитан, блокеров для личного использования не найдено;
- reply-файлов у игры нет.

**Открыто:** проверка Overwolf GEP (не блокирует разработку, см. ТЗ раздел 20).

## Структура

```
cmd/orchestrator/       — точка входа Go-приложения (ядро: бизнес-логика, хранение, UI)
internal/api/            — клиент официального Legion TD 2 API v2 (офлайн-контур)
internal/advisor/        — рекомендации (Phase 1: эвристики, Phase 3: + ML)
internal/perceptionclient/ — сгенерированный gRPC-клиент к Perception Service (появится после protoc)
internal/storage/        — SQLite (миграции в /migrations)
proto/perception.proto   — контракт Go <-> Python
services/perception/     — Python: захват экрана, OCR, CV, ONNX (единственное место с изображениями)
migrations/               — SQL-схема runtime-хранилища и кэша справочников API
```

## Как запустить (после реализации Phase 1)

```bash
# 1. Сгенерировать gRPC-код из proto (Go)
protoc --go_out=. --go-grpc_out=. proto/perception.proto

# 2. Сгенерировать gRPC-код из proto (Python)
cd services/perception
python -m grpc_tools.protoc -I../../proto --python_out=. --grpc_python_out=. ../../proto/perception.proto

# 3. Поднять Perception Service
docker compose up perception

# 4. Собрать и запустить оркестратор
go run ./cmd/orchestrator
```

## Следующий шаг

Phase 1.1–1.2 (см. ТЗ, раздел 20): реализовать SQLite-хранилище, HTTP-вызовы в
`internal/api/client.go` (сейчас `panic("not implemented")`), и минимальный
`ReadEconomy` в Perception Service поверх любого готового OCR-движка.
