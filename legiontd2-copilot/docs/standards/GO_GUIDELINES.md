# GO GUIDELINES

## Project Conventions

- Module path: `github.com/yourname/legiontd2-copilot`
- Go version: 1.22
- Logging: `log/slog` with JSON handler
- Testing: standard `testing` package

## Package Structure

```
internal/
  advisor/      — business logic, heuristic rules
  api/          — LT2 API v2 HTTP client
  dataset/      — dataset collection (offline)
  storage/      — SQLite wrapper
  unitdata/     — static unit meta-data (attack types)
  wavedata/     — static wave data
  http/         — временно, будет заменён на Wails
  ws/           — временно, будет заменён на Wails IPC
```
