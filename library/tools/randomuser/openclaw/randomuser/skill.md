# randomuser — OpenClaw skill

> randomuser should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `randomuser`
- **MCP:** `randomuser mcp`
- **Base URL:** https://randomuser.me
- **Operations:** 1 (1 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `randomuser api` — GET /api

## All operations (first 200)
- `randomuser call getRoot` (GET /api) — read-list — GET /api

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.