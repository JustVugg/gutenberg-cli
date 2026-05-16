# npm-registry — OpenClaw skill

> npm-registry should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `npm-registry`
- **MCP:** `npm-registry mcp`
- **Base URL:** https://registry.npmjs.org
- **Operations:** 2 (2 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `npm-registry express` — GET /express
- `npm-registry search` — GET /-/v1/search

## All operations (first 200)
- `npm-registry call getExpress` (GET /express) — read-list — GET /express
- `npm-registry call getSearch` (GET /-/v1/search) — search — GET /-/v1/search

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.