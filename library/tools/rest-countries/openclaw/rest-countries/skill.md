# rest-countries — OpenClaw skill

> rest-countries should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `rest-countries`
- **MCP:** `rest-countries mcp`
- **Base URL:** https://restcountries.com
- **Operations:** 3 (3 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `rest-countries italy` — GET /v3.1/name/italy
- `rest-countries europe` — GET /v3.1/region/europe
- `rest-countries independent` — GET /v3.1/independent

## All operations (first 200)
- `rest-countries call getNameItaly` (GET /v3.1/name/italy) — read-list — GET /v3.1/name/italy
- `rest-countries call getRegionEurope` (GET /v3.1/region/europe) — read-list — GET /v3.1/region/europe
- `rest-countries call getIndependent` (GET /v3.1/independent) — read-list — GET /v3.1/independent

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.