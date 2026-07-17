# CODING STANDARDS — LT2 Copilot

## General

1. **Language**: Code, comments, commit messages, ADRs — all in English. UI strings may be in Russian (target audience).
2. **No TODOs in committed code**: Use BACKLOG.md or issue tracker for todos. Committed TODOs are tech debt.
3. **Test coverage**: Every new feature must include tests. Minimum: unit tests for business logic, integration for boundaries.
4. **Error handling**: Never swallow errors. Use structured logging (`slog`) with context. Every error path must be logged.
5. **No panics**: Use `error` returns, never `panic`/`recover` except in `main()` for fatal errors.

## Code Style

6. **Go**: Follow `gofmt` and `golangci-lint` defaults. Use `slog` for logging, `std` library for testing.
7. **Python**: Follow PEP 8. Use `ruff` for linting, `black` for formatting.
8. **JS**: Standard JS (no framework for now). Clean functions, no jQuery.
9. **No comments in code**: Write clean, self-documenting code. Architecture decisions go in ADRs.
10. **Imports**: Group: std → third-party → internal. Separate groups with blank line.

## Process

11. **ADR first**: Any architectural decision requires an ADR before implementation.
12. **TDD cycle**: spec → test → code → refactor.
13. **Atomic commits**: ~100 lines per commit, single concern.
14. **Pre-merge**: run review checklist.
