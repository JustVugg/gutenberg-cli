# pypi — OpenClaw skill

> pypi should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `pypi`
- **MCP:** `pypi mcp`
- **Base URL:** https://pypi.org
- **Operations:** 2 (2 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `pypi requests` — GET /pypi/requests/json
- `pypi numpy` — GET /pypi/numpy/json

## All operations (first 200)
- `pypi call getPypiRequests` (GET /pypi/requests/json) — read-list — GET /pypi/requests/json
- `pypi call getPypiNumpy` (GET /pypi/numpy/json) — read-list — GET /pypi/numpy/json

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.