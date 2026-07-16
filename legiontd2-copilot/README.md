# LT2 Copilot — интеллектуальный помощник для Legion TD 2

Desktop-приложение под Windows для рекомендаций во время матча Legion TD 2: расстановка юнитов, покупки/продажи, найм Mercenaries, управление Mythium.

**Статус:** Phase 1.3 — OCR Mythium/HP/таймер через EasyOCR, веб-интерфейс.

Полное ТЗ: [TZ_LegionTD2_Assistant.md](TZ_LegionTD2_Assistant.md)

---

## Структура проекта

```ascii
cmd/orchestrator/             Точка входа Go-оркестратора
internal/
  api/                        HTTP-клиент Legion TD 2 API v2
  advisor/                    Эвристический Advisor (правила трат/накопления)
  perceptionclient/           gRPC-клиент к Python-сервису распознавания
  perceptionclient/pb/        Go-стабы из proto
  storage/                    SQLite-хранилище
  webserver/                  HTTP-сервер + встроенный веб-интерфейс
  webserver/static/           HTML/CSS/JS фронтенд
migrations/                   SQL-схема (справочно)
proto/perception.proto        gRPC-контракт Go <-> Python
services/perception/          Python: захват экрана, EasyOCR, gRPC-сервер
```

## Зависимости

- **Go 1.22+** — `google.golang.org/grpc`, `modernc.org/sqlite`, `google.golang.org/protobuf`
- **Python 3.12+** — `easyocr`, `opencv-python-headless`, `mss`, `grpcio`, `numpy`
- **Protoc** — для перегенерации gRPC-стабов при изменении proto (не обязательно для запуска)

---

## Быстрый старт

### 1. Установить Python-зависимости

```bash
cd services/perception
pip install -r requirements.txt
```

### 2. Запустить Perception Service (OCR)

```bash
python server.py
# слушает localhost:50051
```

> **Важно:** EasyOCR при первом запуске скачает модель (~100 МБ).  
> Регионы захвата настроены на 1920×1080 — для другого разрешения отредактируй `self.regions` в `server.py`.

### 3. Запустить оркестратор с веб-интерфейсом

```bash
# Терминал 2
go run ./cmd/orchestrator
# веб-интерфейс на http://localhost:8080
```

Переменные окружения:
| Переменная | По умолчанию | Описание |
|-----------|-------------|---------|
| `LT2_DB_PATH` | `lt2_copilot.db` | Путь к SQLite |
| `LT2_PERCEPTION_ADDR` | `localhost:50051` | Адрес Perception Service |
| `LT2_WEB_ADDR` | `:8080` | Порт веб-интерфейса |
| `LT2_API_KEY` | — | API-ключ Legion TD 2 для `/units` и `/games` |

### 4. Открыть веб-интерфейс

Перейди в браузере на **http://localhost:8080**  
(рядом с запущенной игрой на втором мониторе или в оконном режиме)

---

## Что работает сейчас

| Компонент | Статус |
|-----------|--------|
| Go-оркестратор | ✅ Запускается, логирует, вызывает gRPC, веб-сервер |
| SQLite (internal/storage) | ✅ Миграции, таблицы готовы |
| gRPC-клиент (internal/perceptionclient) | ✅ ReadEconomy/HealthCheck |
| Эвристики (internal/advisor) | ✅ Правила трат/накопления Mythium |
| API-клиент (internal/api) | ✅ HTTP-вызовы /units/byVersion, /players/matchHistory |
| Perception Service (Python) | ✅ EasyOCR, захват mss, gRPC-сервер |
| Веб-интерфейс (internal/webserver) | ✅ HTML/CSS/JS, автообновление, /api/state |
| Wails UI | ⏳ Факультативно (замена веб-интерфейса) |

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

- **Phase 2** — сбор датасета через API, ML-модель рекомендаций расстановки
- **Phase 3** — прогноз удержания волны, пост-матчевый отчёт
