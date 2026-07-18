# PROJECT STATE

**Date:** 2026-07-18  
**Phase:** Phase 1.1 (Feature Matrix) — завершён

## What Works

- Go orchestrator with 2-second tick loop
- WebSocket Hub receives game state from HudApi patcher
- Heuristic advisor: counter-wave detection, build/send/upgrade recommendations
- SQLite storage with migrations (matches, wave_snapshots, recommendations, units_reference)
- HTTP server with REST endpoints and browser-based UI (временное решение)
- API v2 client: players, matches, units, waves, legions
- Dataset collected: 3000 games, 11,899 player entries, 50.5% win rate
- Static wave data (21 waves) and unit data (fighter/merc attack types)
- Web UI (temporary, browser-based) with debug panel

## New in Phase 1.1

- **Feature Matrix** (`internal/matrix/matrix.go`): 7 sections — economy, board, available,
  waves, coverage, context, opponent. Built from GameState on each /api/state request.
- **Patcher v2** (`internal/deploy/patcher/copilot-patcher.js`): cost extraction from
  subheader HTML (Gold.png/Mythium.png/Supply.png), icon-based CamelCase name fallback,
  raw action dump (first hand unit + first merc). New listeners:
  - `refreshScoreboardInfo` — all players economy + grid
  - `refreshTeamGold` — left/right team gold
  - `refreshKingUpgrades` — attack/regen/spell per side
  - `refreshMoneylender` — gold/cost/enabled
- **FighterAttack map**: 60+ fighters, all major mercs.
- **Board analysis**: only field units (not hand), total/max HP.
- **Coverage**: only affordable units counted by damage type.
- **Deploy package**: Go-пакет для бинарной вставки патчера.
- **Debug panel**: raw JSON data, copy button in Web UI.
- **start.bat** fix: correct working directory.

## What's Removed (previous session)

- OCR/CV полностью удалено — Python Perception Service, EasyOCR, OpenCV, MSS
- Docker композ удалён
- JSON-конфиг регионов удалён
- internal/config/ удалён
- services/perception/ — удалено целиком

## Current Architecture

```
Game (Coherent GT) → HudApi → copilot-patcher.js → WS → Go Hub → HTTP REST → Browser UI
                                                          ↓
                                              Feature Matrix (7 sec.)
                                                          ↓
                                                  Heuristic Advisor
                                                          ↓
                                                  SQLite Storage
```

## What's Next

1. ML-матрица: разбор правил, датасеты, разметка (Phase 1.2)
2. Заменить браузерный UI на Wails Desktop
3. Начать ML-тренировку на датасете 3000 игр
4. Написать unit-тесты
