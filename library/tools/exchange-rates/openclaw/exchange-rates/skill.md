# exchange-rates — OpenClaw skill

> exchange-rates should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `exchange-rates`
- **MCP:** `exchange-rates mcp`
- **Base URL:** https://open.er-api.com
- **Operations:** 2 (2 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `exchange-rates usd` — GET /v6/latest/USD
- `exchange-rates eur` — GET /v6/latest/EUR

## All operations (first 200)
- `exchange-rates call getLatestUsd` (GET /v6/latest/USD) — read-list — GET /v6/latest/USD
- `exchange-rates call getLatestEur` (GET /v6/latest/EUR) — read-list — GET /v6/latest/EUR

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.