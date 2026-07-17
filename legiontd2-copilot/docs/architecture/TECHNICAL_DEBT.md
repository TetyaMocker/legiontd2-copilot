# TECHNICAL DEBT

## Critical

| ID | Issue | Location | Fix Priority |
|----|-------|----------|-------------|
| DEBT-001 | **Zero test coverage**: no test files anywhere | Entire project | HIGH |
| DEBT-002 | **HTTP/WS in desktop app**: browser UI вместо Wails | `internal/http/`, `internal/ws/` | HIGH |
| DEBT-003 | **copilot-patcher.js**: может сломаться при патче игры | `patches/copilot-patcher.js` | MEDIUM |

## Medium

| ID | Issue | Location | Fix Priority |
|----|-------|----------|-------------|
| DEBT-004 | **Duplicate unit data**: static Go map vs DB cache | `internal/unitdata/units.go` + DB | MEDIUM |
| DEBT-005 | **Race condition**: Hub.broadcast может блокировать tick | `internal/ws/hub.go` | MEDIUM |
| DEBT-006 | **No graceful shutdown**: HTTP server started but never shut down | `internal/http/server.go` | MEDIUM |
| DEBT-007 | **Wails config incomplete**: config exists, no implementation | `wails.json`, `frontend/` | MEDIUM |

## Low

| ID | Issue | Location | Fix Priority |
|----|-------|----------|-------------|
| DEBT-008 | **Dataset binary path**: default outDir is "dataset" not "data/dataset" | `cmd/dataset/main.go` | LOW |
| DEBT-009 | **start.bat outdated**: refers to web URL | `scripts/start.bat` | LOW |
| DEBT-010 | **frontend/ package.json**: only vite, no actual SPA | `frontend/package.json` | LOW |

## Already Fixed This Session

- OCR/CV полностью удалён из проекта
- Python Perception Service удалён
- Docker удалён
- JSON-конфиг регионов удалён
- internal/config/ удалён
- ТЗ обновлено до версии 2.0
- ADR-003, ADR-010, ADR-011 помечены как Superseded
- ADR-004 переписан (gRPC отменён)
- Все ADR и docs обновлены под новую архитектуру
