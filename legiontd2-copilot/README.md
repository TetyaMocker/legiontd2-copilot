# LT2 Copilot — интеллектуальный помощник для Legion TD 2

Desktop-приложение под Windows для рекомендаций во время матча Legion TD 2: расстановка юнитов, покупки/продажи, найм Mercenaries, управление Mythium.

**Статус:** Phase 1.4 — Переход на REST+JSON-конфиги, подготовка Wails UI.

Полное ТЗ: [TZ_LegionTD2_Assistant.md](TZ_LegionTD2_Assistant.md)

---

## Архитектура (новая, актуальная)

```
┌─────────────────┐     POST /api/ingest     ┌──────────────────┐     fetch /api/state     ┌─────────┐
│  Perception      │ ──────────────────────▶  │  Go Orchestrator │ ◀──────────────────────  │ Browser │
│  Service         │   {"mythium":120,...}    │  (internal/web-  │      JSON + Recs         │  UI     │
│  (Python/OCR)    │                          │   server + ad-   │                          │         │
└─────────────────┘                          │   visor + conf)  │                          └─────────┘
                                             └──────────────────┘
                                                       │
                                                       ▼
                                             ┌──────────────────┐      ┌──────────────────┐
                                             │  SQLite Storage  │      │  LT2 API Client   │
                                             │  (матчи/логи)    │      │  (офиц. API v2)  │
                                             └──────────────────┘      └──────────────────┘
```

## Структура проекта

```
cmd/orchestrator/            Точка входа Go-оркестратора
config/
  regions.json               JSON-конфиг регионов OCR (паттерн)
internal/
  config/                    Загрузчик JSON-конфигов
  api/                       HTTP-клиент Legion TD 2 API v2
  advisor/                   Эвристический Advisor (правила трат/накопления)
  storage/                   SQLite-хранилище
  webserver/                 HTTP-сервер + встроенный веб-интерфейс
  webserver/static/          HTML/CSS/JS фронтенд (текущий)
frontend/                    Wails UI (будущий in-app интерфейс)
services/perception/         Python: захват экрана, EasyOCR, HTTP POST
```

## Зависимости

- **Go 1.22+** — `modernc.org/sqlite` (pure-Go, без CGO)
- **Python 3.12+** — `easyocr`, `opencv-python-headless`, `mss`, `numpy`

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
# отправляет POST на http://localhost:8080/api/ingest
```

> **Важно:** EasyOCR при первом запуске скачает модель (~100 МБ).  
> Регионы захвата настраиваются в `config/regions.json` (1920×1080).

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
| `LT2_WEB_ADDR` | `:8080` | Порт (API + веб-интерфейс) |
| `LT2_API_KEY` | — | API-ключ Legion TD 2 для `/units` и `/games` |

### 4. Открыть веб-интерфейс

Перейди в браузере на **http://localhost:8080**  
(рядом с запущенной игрой на втором мониторе или в оконном режиме)

---

## Что работает сейчас

| Компонент | Статус |
|-----------|--------|
| Go-оркестратор | ✅ Запускается, принимает POST-данные, веб-сервер |
| SQLite (internal/storage) | ✅ Миграции, таблицы готовы |
| REST-приём OCR (POST /api/ingest) | ✅ Замена gRPC |
| Эвристики (internal/advisor) | ✅ Правила трат/накопления Mythium |
| API-клиент (internal/api) | ✅ HTTP-вызовы /units/byVersion, /players/matchHistory |
| Perception Service (Python) | ✅ EasyOCR, захват mss, HTTP POST |
| Веб-интерфейс (internal/webserver) | ✅ HTML/CSS/JS, автообновление, /api/state |
| JSON-конфиги (config/regions.json) | ✅ Регионы OCR вынесены в конфиг |
| Трекинг OCR (Tracker в Python) | ✅ Кеширование, контроль частоты сканов |
| Wails UI | ⏳ Факультативно (замена веб-интерфейса, нужен gcc+node) |

---

## План дальнейшей разработки

- **Phase 2** — сбор датасета через API, ML-модель рекомендаций расстановки
- **Phase 3** — прогноз удержания волны, пост-матчевый отчёт
- **Overwolf overlay** — вывод рекомендаций прямо поверх игры (после подтверждения Game ID в Overwolf GEP)
- **Wails in-app UI** — замена браузера на нативное окно (требует MinGW-w64 + Node.js)
