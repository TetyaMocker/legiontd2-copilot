# SYSTEM ARCHITECTURE — LT2 Copilot

## High-Level Design

```
┌──────────────────────────────────────────────────────────────┐
│                   Windows Desktop Application                  │
│                   (Go + Wails Native UI)                       │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐  │
│  │  Wails UI    │  │    Advisor   │  │    SQLite Storage   │  │
│  │  (native     │◄─┤  (heuristic  │◄─┤  (matches, snaps,  │  │
│  │   window)    │  │   + ML)      │  │   recommendations)  │  │
│  └──────────────┘  └──────▲───────┘  └────────────────────┘  │
│                           │                                   │
│                     ┌─────┴──────┐                            │
│                     │ GameState  │                            │
│                     │   Hub      │                            │
│                     └─────▲──────┘                            │
└───────────────────────────┼──────────────────────────────────┘
                            │ WebSocket
             ┌──────────────┴──────────────┐
             │  Game Client                  │
             │  (Coherent GT HudApi)         │
             │  + copilot-patcher.js         │
             └──────────────────────────────┘

┌──────────────────────────────────────────────────────────────┐
│  Offline ML Pipeline (Python, separate from runtime)          │
│  - cmd/dataset: collect games from LT2 API v2                │
│  - Train model: placement recommender + wave-hold predictor  │
│  - Export model → load into Go Advisor                        │
└──────────────────────────────────────────────────────────────┘
```

## Architecture Principles

1. **Desktop-first**: Wails native window, no browser, no embedded web server.
2. **Single data source**: HudApi (Coherent GT) for live game state.
3. **SQLite** for runtime match logs and API data cache.
4. **ML offline**: Python training pipeline, model loaded into Go Advisor.
5. **Two-layer skill system**: generic engineering discipline (addyosmani) + project-specific (ltd2-*).

## Component Responsibilities

| Component | Responsibility |
|-----------|---------------|
| **Go Core** | Orchestrator, business logic (advisor), storage, UI bridge |
| **Wails UI** | Desktop window, renders game state and recommendations |
| **SQLite** | Match logs, wave snapshots, recommendations, cached API unit data |
| **Offline ML** | Python scripts for dataset collection and model training |
| **copilot-patcher.js** | Captures game events via HudApi, sends via WebSocket |

## Data Flow

1. Game (Coherent GT) → HudApi events → `copilot-patcher.js` → WebSocket → Hub
2. Hub → Advisor (heuristics + ML) → recommendations
3. Recommendations → Wails UI + SQLite
4. No data leaves the local machine
