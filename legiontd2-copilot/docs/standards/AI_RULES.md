# AI RULES — For AI-assisted development

## Entry Protocol

1. Read `AGENTS.md`
2. Read all ADRs in `docs/architecture/ADR/`
3. Read `docs/development/PROJECT_STATE.md`
4. Read `docs/architecture/TECHNICAL_DEBT.md`
5. Read `docs/architecture/DECISIONS.md`
6. Only then start work

## Update Protocol

After any significant task:
1. Update `PROJECT_STATE.md`
2. Update `NEXT_STEPS.md`
3. Update `TECHNICAL_DEBT.md`
4. Update `CHANGELOG_DEV.md`
5. New architectural decision → new ADR

## Priority Rule

Project-specific skill (`ltd2-*`) wins over general engineering practice.

## Known Sensitive Points

- `internal/http/` + `internal/ws/` — временные, будут удалены после Wails
- `internal/unitdata/units.go` vs DB `units_reference` — не допускать расхождения
- `patches/copilot-patcher.js` — единственный источник live-данных, критичен
