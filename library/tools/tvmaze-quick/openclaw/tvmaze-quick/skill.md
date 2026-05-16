# tvmaze-quick — OpenClaw skill

> tvmaze-quick should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `tvmaze-quick`
- **MCP:** `tvmaze-quick mcp`
- **Base URL:** https://api.tvmaze.com
- **Operations:** 1 (1 read, 0 write, 0 destructive)

## All operations (first 200)
- `tvmaze-quick call getShows` (GET /shows/{id}) — read-detail — GET /shows/{id}

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.