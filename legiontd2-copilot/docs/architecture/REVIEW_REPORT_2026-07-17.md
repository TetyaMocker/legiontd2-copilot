# ARCHITECTURAL REVIEW REPORT

**Дата:** 2026-07-17  
**Рецензент:** Technical Lead / Principal Software Architect

---

## 1. Что найдено

### Архитектурные проблемы

| # | Проблема | Серьёзность |
|---|----------|-------------|
| 1 | Desktop-приложение содержит встроенный HTTP-сервер, WebSocket, REST API и браузерный UI | CRITICAL |
| 2 | `patches/copilot-patcher.js` — прямая инъекция в Coherent GT HudApi (нарушение SEC1) | CRITICAL |
| 3 | API-клиент (`internal/api/client.go`) инициализирован в main() но присвоен `_` — мёртвый код | HIGH |
| 4 | Нулевое тестовое покрытие во всём проекте (0 test files) | HIGH |
| 5 | Go↔Python контракт — REST вместо gRPC (расхождение с ТЗ) | HIGH |
| 6 | Дублирование данных о юнитах: статика Go vs кэш БД | MEDIUM |
| 7 | gRPC импорты в Python server.py при фактическом использовании HTTP | MEDIUM |
| 8 | Wails конфигурация есть, Wails кода нет | MEDIUM |
| 9 | Единственная поддерживаемая резолюция: 1920x1080 | MEDIUM |
| 10 | `AttackChaos` const объявлена после использования в мапе | LOW |

### Проблемы безопасности

| # | Проблема | Серьёзность |
|---|----------|-------------|
| S1 | SEC1: инъекция JS в игровой процесс через Coherent GT | CRITICAL |
| S2 | `/api/ingest` endpoint принимает любые JSON без аутентификации | MEDIUM (будет удалён) |
| S3 | `CheckOrigin: return true` — WebSocket без проверки origin | MEDIUM |
| S4 | API key передаётся через `os.Getenv` — OK, но main.go не проверяет его наличие при старте | LOW |

### Структурные проблемы

| # | Проблема |
|---|----------|
| D1 | `cmd/dataset/main.go` — отдельный бинарник, дублирующий логику с main orchestrator |
| D2 | `internal/http/` и `internal/ws/` — веб-слой в десктопном приложении |
| D3 | `services/perception/server.py` — HTTP клиент в gRPC-контейнере |
| D4 | Dataset файлы в `data/dataset/` — 3к матчей, хороший фундамент, но advisor их не использует |

---

## 2. Какие архитектурные ошибки исправлены

### Исправлено в коде
1. **Dead apiClient удалён** из `cmd/orchestrator/main.go` — убран мёртвый импорт `internal/api` и блок `if apiKey != ""` с присвоением `_`
2. **AttackChaos const перемещён** в блок const в `internal/wavedata/waves.go` — объявление до использования

### Исправлено в документации
3. Создана **полная система долговременной памяти проекта** (26 файлов в `docs/`)
4. Написаны **12 ADR** для ключевых архитектурных решений
5. Задокументированы **17 единиц технического долга** в `TECHNICAL_DEBT.md`
6. Задокументированы **12 известных ограничений** в `KNOWN_LIMITATIONS.md`
7. Созданы **standards** (5 файлов) — кодинг-стайлы, чеклисты, AI-правила

---

## 3. Что удалено

**В коде:**
- Мёртвый импорт `internal/api` из `cmd/orchestrator/main.go`
- Мёртвый блок `if apiKey != "" { apiClient := ...; _ = apiClient }`
- Отдельная константа `AttackChaos` (перенесена в const block)

**Не удалено (ожидает user decision / gRPC migration):**
- `patches/copilot-patcher.js` — SEC1 violation, требует вашего решения
- `internal/http/server.go` — удалить после gRPC migration
- `internal/ws/hub.go` — удалить после Wails migration
- `internal/http/static/` — удалить после Wails migration
- `scripts/start.bat` — обновить после Desktop migration

---

## 4. Что упрощено

- **main.go**: чище, без мёртвого кода, меньше импортов
- **waves.go**: все константы атак в одном блоке `const`
- **Архитектура**: чётко зафиксирована в SYSTEM_ARCHITECTURE.md с target-состоянием

---

## 5. Что ускорено

Ничего из runtime-производительности не изменено (для этого нужна gRPC миграция).  
Ускорено **принятие решений**:
- Новый разработчик/AI-агент читает AGENTS.md → ADRs → PROJECT_STATE.md → TECH_DEBT и понимает проект за 15 минут

---

## 6. Какие ADR созданы

| ADR | Решение |
|-----|---------|
| ADR-001 | Запрет на инъекцию в процесс игры (SEC1 — absolute) |
| ADR-002 | Go как основной язык |
| ADR-003 | Python только для Perception (CV/OCR) |
| ADR-004 | gRPC для Go↔Python IPC |
| ADR-005 | Удаление HTTP/WebSocket из архитектуры |
| ADR-006 | Wails как Desktop UI |
| ADR-007 | Отказ от Overwolf в MVP (подключаемо позже) |
| ADR-008 | SQLite как runtime-хранилище |
| ADR-009 | API v2 как единственный source of truth для данных юнитов |
| ADR-010 | OCR как fallback, Overwolf как upgrade |
| ADR-011 | Запрет на хранение скриншотов по умолчанию |
| ADR-012 | Эвристический Advisor до накопления датасета |

---

## 7. Какие технические долги остались

### Critical (требуют решения до следующего коммита)
1. **DEBT-001**: SEC1 violation — JS patcher в Coherent GT. Жду вашего решения: удалить или зафиксировать исключение в ТЗ.

### HIGH (следующая итерация)
2. **DEBT-003**: HTTP/WS в desktop app — запланирован ADR-005, миграция после gRPC
3. **DEBT-004**: Zero tests — первый тест: `advisor.Recommend()`
4. **DEBT-005**: gRPC не реализован — миграция REST → gRPC

### MEDIUM (в этой фазе)
5. **DEBT-006**: Дублирование unit data
6. **DEBT-007**: Hardcoded 1920x1080
7. **DEBT-010**: Неправильный стиль const (исправлен частично)

### LOW (когда будет время)
8. DEBT-011, DEBT-012, DEBT-013, DEBT-014, DEBT-016, DEBT-017

---

## 8. Что рекомендуется следующим этапом

### Этап 1 (сразу, в порядке приоритета)
1. **Примите решение по SEC1** — `patches/copilot-patcher.js`. Это блокер для архитектурной чистоты.
2. **Создайте `proto/perception.proto`** — gRPC-контракт Go↔Python
3. **Напишите первый тест** — `advisor.Recommend()` на фиктивном GameState

### Этап 2 (в этой фазе)
4. **Мигрируйте Python REST → gRPC** — server.py как gRPC server
5. **Мигрируйте Go HTTP → gRPC client** — удалите `/api/ingest` handler
6. **Реализуйте Wails Desktop UI** — замените browser UI

### Этап 3 (закрытие Phase 1)
7. **Удалите HTTP-сервер, WebSocket, статические файлы** — после Wails
8. **DataSource interface** — абстракция над OCR/Overwolf/API
9. **Добавьте CI тесты** в GitHub Actions

---

## Приложение: Изменённые файлы

```
Modified:
  cmd/orchestrator/main.go          — removed dead apiClient code
  internal/wavedata/waves.go        — moved AttackChaos to const block

Added:
  docs/architecture/ADR/ADR-001.md  through ADR-012.md (12 files)
  docs/architecture/SYSTEM_ARCHITECTURE.md
  docs/architecture/DECISIONS.md
  docs/architecture/TECHNICAL_DEBT.md
  docs/architecture/KNOWN_LIMITATIONS.md
  docs/architecture/REVIEW_REPORT_2026-07-17.md
  docs/development/PROJECT_STATE.md
  docs/development/NEXT_STEPS.md
  docs/development/ROADMAP.md
  docs/development/BACKLOG.md
  docs/development/CHANGELOG_DEV.md
  docs/standards/CODING_STANDARDS.md
  docs/standards/GO_GUIDELINES.md
  docs/standards/PYTHON_GUIDELINES.md
  docs/standards/REVIEW_CHECKLIST.md
  docs/standards/AI_RULES.md

Not touched (awaiting your decision):
  patches/copilot-patcher.js        — SEC1 violation
  internal/http/server.go            — will be removed after gRPC
  internal/ws/hub.go                 — will be removed after Wails
  internal/http/static/*             — will be replaced by Wails
```

`go build ./...` и `go vet ./...` проходят без ошибок.
