# LT2 Copilot — интеллектуальный помощник для Legion TD 2

Desktop-приложение под Windows для рекомендаций во время матча Legion TD 2: расстановка юнитов, покупки/продажи, найм Mercenaries, управление Mythium.

**Статус:** Phase 1.1 — Feature Matrix + расширенный патчер (scoreboard, team gold, king upgrades, moneylender). Следующий этап: ML-матрица и разметка датасетов.

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
patches/
  copilot-patcher.js         JS-патч для вставки в UI игры (manual-install копия)
internal/
  api/                       HTTP-клиент Legion TD 2 API v2
  advisor/                   Эвристический Advisor
  deploy/                    Пакет автоустановки патча в игру
  deploy/patcher/            Канонический copilot-patcher.js (встраивается бинарно)
  http/                      HTTP-сервер + встроенный веб-интерфейс
  http/static/               HTML/CSS/JS фронтенд
  matrix/                    Feature Matrix (AI-вход: 7 секций, board/coverage/opponent)
  storage/                   SQLite-хранилище
  unitdata/                  Статические данные юнитов (тип атаки, стоимость)
  ws/                        WebSocket Hub (приём GameState, рассылка реком.)
docs/
  ml-training/               Документация ML-матрицы и данных
  architecture/              ADR, ограничения, архитектура
  development/               Чейнджлоги, роадмап, бэклог
frontend/                    Wails UI (будущий in-app интерфейс)
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
| Feature Matrix (internal/matrix) | ✅ 7 секций: economy/board/available/waves/coverage/context/opponent |
| Парсинг стоимости (subheader) | ✅ Извлечение gold/mythium/supply из HTML-подписи действия |
| Имя юнита из иконки (CamelCase) | ✅ Fallback при пустом имени |
| scoreboardInfo (все игроки) | ✅ Gold/mythium/income/workers/supply/value/grid |
| teamGold/KingUpgrades/Moneylender | ✅ События подключены и разбираются |
| Эвристики (internal/advisor) | ✅ Правила трат/накопления Mythium |
| API-клиент (internal/api) | ✅ HTTP-вызовы /units/byVersion, /players/matchHistory |
| SQLite (internal/storage) | ✅ Миграции, таблицы готовы |
| Веб-интерфейс (internal/http/static) | ✅ HTML/CSS/JS, WebSocket + polling, debug-панель с raw dump |
| Copilot Patcher (JS) | ✅ Каноническая версия в internal/deploy/patcher, копия в patches/ |
| Автоустановка патча (deploy) | ✅ Go-пакет для копирования .js в директорию игры |
| Wails UI | ⏳ Факультативно (нужен gcc+node) |

---

## План дальнейшей разработки

- **Phase 1.1** ✅ — Feature Matrix + расширенный патчер + debug UI
- **Phase 1.2** — ML-матрица: разбор правил, датасеты, разметка, компоновка в единое приложение
- **Phase 1.3** — Wails Desktop UI (замена браузера)
- **Phase 2** — ML-модель для предсказаний состава на волну
- **Phase 3** — Refinement: объяснения, пост-матчевые отчёты, мерк-рекомендации
