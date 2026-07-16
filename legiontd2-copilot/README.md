# LT2 Copilot — интеллектуальный помощник для Legion TD 2

Desktop-приложение под Windows для рекомендаций во время матча Legion TD 2: расстановка юнитов, покупки/продажи, найм Mercenaries, управление Mythium.

**Статус:** Phase 1.1 — каркас проекта. Код компилируется, структура готова, реализации OCR/CV и UI — заглушки.

Полное ТЗ: [TZ_LegionTD2_Assistant.md](TZ_LegionTD2_Assistant.md)

---

## Структура проекта

```ascii
cmd/orchestrator/           Точка входа Go-оркестратора
internal/
  api/                      HTTP-клиент официального Legion TD 2 API v2
  advisor/                  Эвристический Advisor (правила, не ML)
  perceptionclient/         gRPC-клиент к Python-сервису распознавания
  perceptionclient/pb/      Go-стабы из proto (сгенерированы protoc)
  storage/                  SQLite-хранилище (миграции встроены)
migrations/                 SQL-схема (для справки)
proto/perception.proto      gRPC-контракт Go <-> Python
services/perception/        Python: захват экрана, OCR, CV, ONNX
```

## Зависимости

- **Go 1.22+** (`go`, установлен)
- **Python 3.12+** (`python`, установлен)
- **Protoc** (`protoc`) — для перегенерации gRPC-кода при изменении proto

Go-зависимости (в `go.mod`): `google.golang.org/grpc`, `modernc.org/sqlite`, `google.golang.org/protobuf`  
Python-зависимости (в `requirements.txt`): `grpcio`, `opencv-python-headless`, `mss`, `numpy`

---

## Как запустить

### 1. Perception Service (Python) — в Docker или напрямую

Через Docker:
```bash
docker compose up perception
```
Или напрямую (быстрее для разработки):
```bash
cd services/perception
pip install -r requirements.txt
python server.py
```
Сервис слушает на `localhost:50051`.

### 2. Оркестратор (Go)

```bash
go run ./cmd/orchestrator
```

Оркестратор:
- Подключается к SQLite (`lt2_copilot.db` в текущей директории)
- Пытается соединиться с Perception Service на `localhost:50051`
- Каждые 2 секунды опрашивает `ReadEconomy` и прогоняет через эвристики Advisor'а
- Логирует результат в JSON (stdout)

Переменные окружения:
- `LT2_DB_PATH` — путь к SQLite (по умолч. `lt2_copilot.db`)
- `LT2_PERCEPTION_ADDR` — адрес Perception Service (по умолч. `localhost:50051`)

### 3. Или всё сразу (Go без Python)

Оркестратор штатно работает при недоступном Perception Service — пишет WARN и ждёт. Можно запустить только Go-часть для проверки сборки:
```bash
go build ./...
```

---

## Что работает сейчас

| Компонент | Статус |
|-----------|--------|
| Go-оркестратор (cmd/orchestrator) | ✅ Запускается, логирует, вызывает gRPC |
| SQLite (internal/storage) | ✅ Миграции, таблицы matches/wave_snapshots/recommendations/units_reference |
| gRPC-клиент (internal/perceptionclient) | ✅ Соединение, ReadEconomy/HealthCheck |
| Эвристики (internal/advisor) | ✅ Правила трат/накопления Mythium |
| API-клиент (internal/api) | ✅ HTTP-вызовы /units/byVersion, /players/matchHistory |
| Proto-контракт | ✅ Сгенерирован для Go и Python |
| Perception Service (Python) | ✅ gRPC-сервер, заглушки ReadEconomy/HealthCheck |
| OCR/CV | ❌ Чистый экран, возвращает 0 (будет в Phase 1.3) |
| Wails UI | ❌ (будет в Phase 1.5) |

---

## Как перегенерировать gRPC-стабы

Если меняется `proto/perception.proto`:

**Go:**
```bash
protoc --go_out=internal/perceptionclient/pb --go_opt=paths=source_relative \
  --go-grpc_out=internal/perceptionclient/pb --go-grpc_opt=paths=source_relative \
  --proto_path=proto proto/perception.proto
```

**Python:**
```bash
cd services/perception
python -m grpc_tools.protoc -I../../proto --python_out=. --grpc_python_out=. ../../proto/perception.proto
```

---

## План дальнейшей разработки (по ТЗ)

- **Phase 1.3** — OCR: Mythium/HP/таймер/EasyOCR или Tesseract
- **Phase 1.5** — Wails-окно с минимальным UI
- **Phase 2** — сбор датасета через API, ML-модель рекомендаций
- **Phase 3** — прогноз удержания волны, пост-матчевый отчёт
