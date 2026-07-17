# DEVELOPMENT CHANGELOG

## 2026-07-17 — Session 2: OCR removal + architecture restructure

### Added
- TZ v2.0 — новая структура: 3 фазы (Desktop → ML → Refinement)
- Phase 1 checklist (1.1–1.6)
- Phase 2 checklist (2.1–2.6)
- Phase 3 checklist (3.1–3.5)

### Removed
- `services/perception/` — весь Python OCR/CV сервис (EasyOCR, OpenCV, MSS, Docker)
- `docker-compose.yml` — больше не нужен
- `data/config/regions.json` — OCR регионы
- `internal/config/config.go` — загрузчик конфигов (был только для OCR)
- OCR-импорты и мёртвый код

### Changed
- `TZ_LegionTD2_Assistant.md` — полностью переписан (v2.0)
- `AGENTS.md` — убраны ссылки на OCR, python-perception агента
- `internal/storage/migrations.go` — source default: 'ocr' → 'hudapi'
- `docs/architecture/ADR/ADR-001.md` — revised (HudApi легитимен)
- `docs/architecture/ADR/ADR-002.md` — updated
- `docs/architecture/ADR/ADR-003.md` — SUPERSEDED
- `docs/architecture/ADR/ADR-004.md` — rewritten (gRPC cancelled)
- `docs/architecture/ADR/ADR-007.md` — updated
- `docs/architecture/ADR/ADR-010.md` — SUPERSEDED
- `docs/architecture/ADR/ADR-011.md` — SUPERSEDED
- All docs: removed OCR/perception references

### Known Issues
- HTTP/WS + browser UI всё ещё существуют (временное решение до Wails)
- ML-модель не выбрана (Phase 2)
- Нулевое тестовое покрытие
