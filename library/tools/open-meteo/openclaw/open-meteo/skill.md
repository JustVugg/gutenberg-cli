# open-meteo — OpenClaw skill

> open-meteo should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `open-meteo`
- **MCP:** `open-meteo mcp`
- **Base URL:** https://api.open-meteo.com
- **Operations:** 1 (1 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `open-meteo forecast` — GET /v1/forecast

## All operations (first 200)
- `open-meteo call getForecast` (GET /v1/forecast) — read-list — GET /v1/forecast

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.