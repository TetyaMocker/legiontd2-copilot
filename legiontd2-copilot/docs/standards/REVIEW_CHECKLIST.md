# REVIEW CHECKLIST

## Pre-Merge Review

### Architecture
- [ ] Does this change violate TZ_SEC1 (injection into game process)?
- [ ] Does this change touch the Go↔Python boundary? If yes, is gRPC used?
- [ ] Is there a corresponding ADR for any architectural decision?
- [ ] Is the change documented in PROJECT_STATE.md / CHANGELOG_DEV.md?

### Security
- [ ] Does this change expose user data? (screenshots, game state, API keys)
- [ ] Is input validated for all user/external inputs?
- [ ] Are API keys handled via environment variables, never hardcoded?
- [ ] Does this change affect `internal/http/`, `internal/ws/`, `internal/api/`? If yes → mandatory security review.

### Code Quality
- [ ] Follows Go/Python guidelines?
- [ ] No commented-out code
- [ ] No `TODO`/`FIXME`/`XXX` in committed code
- [ ] No unnecessary dependencies
- [ ] No dead code or unused variables

### Testing
- [ ] Unit tests added for new logic?
- [ ] Do existing tests pass?
- [ ] Manual test: what happens when game is closed / API down / patcher disconnected?

### Performance
- [ ] Not blocking the 2-second tick loop with heavy operations
- [ ] gRPC calls are async or in separate goroutine where appropriate
- [ ] No busy loops (use tickers, not `for {}`)

### CI
- [ ] `golangci-lint` passes
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
