# PROJECT STATE

**Date:** 2026-07-17  
**Phase:** Phase 1 (Desktop Application) — в процессе

## What Works

- Go orchestrator with 2-second tick loop
- WebSocket Hub receives game state from HudApi patcher
- Heuristic advisor: counter-wave detection, build/send/upgrade recommendations
- SQLite storage with migrations (matches, wave_snapshots, recommendations, units_reference)
- HTTP server with REST endpoints and browser-based UI (временное решение)
- API v2 client: players, matches, units, waves, legions
- Dataset collected: 3000 games, 11,899 player entries, 50.5% win rate
- Static wave data (21 waves) and unit data (fighter/merc attack types)
- Web UI (temporary, browser-based)

## What's Removed (this session)

- **OCR/CV полностью удалено** — Python Perception Service, EasyOCR, OpenCV, MSS
- **Docker композ** удалён (не нужен без Python-сервиса)
- **JSON-конфиг регионов** (`data/config/regions.json`) удалён
- **internal/config/** удалён (был нужен только для OCR)
- `services/perception/` — удалено целиком

## Current Architecture

```
Game (Coherent GT) → HudApi → copilot-patcher.js → WS → Go Hub → HTTP REST → Browser UI
                                                          ↓
                                                   Heuristic Advisor
                                                          ↓
                                                   SQLite Storage
```

## What's Next

1. Заменить браузерный UI на Wails Desktop
2. Удалить HTTP-сервер и WebSocket (Wails предоставит IPC)
3. Начать ML-тренировку на датасете 3000 игр
4. Написать unit-тесты
