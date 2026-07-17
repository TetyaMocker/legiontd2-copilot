# LT2 Copilot — интеллектуальный помощник для Legion TD 2

Desktop-приложение под Windows для рекомендаций во время матча Legion TD 2: расстановка юнитов, покупки/продажи, найм Mercenaries, управление Mythium.

**Статус:** Phase 0 — WebSocket-приём событий из игры через Coherent GT HudApi.

Полное ТЗ: [TZ_LegionTD2_Assistant.md](TZ_LegionTD2_Assistant.md)

---

## Архитектура

```
┌─────────────────────────────────────────┐
│           Игра (Legion TD 2)            │
│  ┌─────────────────────────────────┐   │
│  │  Copilot Patcher (JS)           │   │
│  │  - engine.on() хуки             │   │
│  │  - WebSocket клиент             │   │
│  └──────────┬──────────────────────┘   │
└─────────────┼──────────────────────────┘
              │ WebSocket (ws://localhost:8080/ws)
              ▼
┌─────────────────────────────────────────┐
│      Go Orchestrator (ядро)             │
│  ┌──────────┐ ┌──────────┐ ┌────────┐ │
│  │ WS Hub   │ │ Heuristic│ │ HTTP   │ │
│  │ (приём   │ │ Advisor  │ │ Server │ │
│  │ событий) │ │          │ │ + stat │ │
│  └──────────┘ └──────────┘ └────────┘ │
│  ┌──────────┐ ┌──────────┐            │
│  │ SQLite   │ │ API      │            │
│  │ Storage  │ │ Client   │            │
│  └──────────┘ └──────────┘            │
└─────────────────────────────────────────┘
              │ WebSocket / HTTP
              ▼
        ┌──────────┐
        │ Browser  │
        │ / Web UI │
        └──────────┘
```

## Структура проекта

```
cmd/orchestrator/            Точка входа Go-оркестратора
config/
  regions.json               JSON-конфиг регионов (OCR fallback)
patches/
  copilot-patcher.js         JS-патч для вставки в UI игры
internal/
  config/                    Загрузчик JSON-конфигов
  api/                       HTTP-клиент Legion TD 2 API v2
  advisor/                   Эвристический Advisor
  ws/                        WebSocket Hub (приём GameState, рассылка реком.)
  http/                      HTTP-сервер + встроенный веб-интерфейс
  http/static/               HTML/CSS/JS фронтенд
  storage/                   SQLite-хранилище
frontend/                    Wails UI (будущий in-app интерфейс)
services/perception/         Python: захват экрана, EasyOCR, HTTP POST
```

## Зависимости

- **Go 1.22+** — `modernc.org/sqlite`, `gorilla/websocket`
- **Python 3.12+** — `easyocr`, `opencv-python-headless`, `mss`, `numpy` (OCR fallback)

---

## Быстрый старт

### 1. Запустить оркестратор

```bash
go run ./cmd/orchestrator
# сервер на http://localhost:8080
```

Переменные окружения:
| Переменная | По умолчанию | Описание |
|-----------|-------------|---------|
| `LT2_DB_PATH` | `lt2_copilot.db` | Путь к SQLite |
| `LT2_WEB_ADDR` | `:8080` | Порт (HTTP + WebSocket) |
| `LT2_API_KEY` | — | API-ключ Legion TD 2 для `/units` и `/games` |

### 2. Установить патч в игру

Скопируй `patches/copilot-patcher.js` в папку:
```
Legion TD 2_Data/uiresources/AeonGT/hud/js/
```

И добавь в `uiresources/AeonGT/gateway.html` перед закрывающим `</body>`:
```html
<script src="hud/js/copilot-patcher.js"></script>
```

> **Внимание:** При обновлении игры файлы сбрасываются — потребуется повторная установка.

### 3. Открыть веб-интерфейс

Перейди в браузере на **http://localhost:8080**  
(рядом с запущенной игрой на втором мониторе или в оконном режиме)

---

## Что работает сейчас

| Компонент | Статус |
|-----------|--------|
| Go-оркестратор | ✅ WebSocket Hub + HTTP + эвристики |
| WebSocket-приём GameState | ✅ Приём данных из игры через Coherent GT HudApi |
| Эвристики (internal/advisor) | ✅ Правила трат/накопления Mythium |
| API-клиент (internal/api) | ✅ HTTP-вызовы /units/byVersion, /players/matchHistory |
| SQLite (internal/storage) | ✅ Миграции, таблицы готовы |
| Веб-интерфейс (internal/http/static) | ✅ HTML/CSS/JS, WebSocket + polling |
| Json-конфиги (config/regions.json) | ✅ Регионы OCR вынесены в конфиг |
| Copilot Patcher (patches/copilot-patcher.js) | ✅ JS-патч для хуков engine.on() |
| OCR Fallback (Python) | ✅ EasyOCR, POST /api/ingest (резервный канал) |
| Wails UI | ⏳ Факультативно (нужен gcc+node) |
| Overwolf overlay | 🔍 Исследуется |

---

## План дальнейшей разработки

- **Phase 1** — Полноценный патч с WebSocket + эвристики + автоустановка
- **Phase 2** — Интеграция API статики юнитов, улучшение эвристик
- **Phase 3** — ML-модель для предсказаний состава на волну
- **Overwolf overlay** — вывод рекомендаций поверх игры (если подтвердится Game ID)
- **Wails in-app UI** — замена браузера на нативное окно
