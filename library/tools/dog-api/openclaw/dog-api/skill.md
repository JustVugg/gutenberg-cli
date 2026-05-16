# dog-api — OpenClaw skill

> dog-api should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `dog-api`
- **MCP:** `dog-api mcp`
- **Base URL:** https://dog.ceo
- **Operations:** 3 (3 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `dog-api all` — GET /api/breeds/list/all
- `dog-api random` — GET /api/breeds/image/random
- `dog-api list` — GET /api/breed/hound/list

## All operations (first 200)
- `dog-api call getBreedsListAll` (GET /api/breeds/list/all) — read-list — GET /api/breeds/list/all
- `dog-api call getBreedsImageRandom` (GET /api/breeds/image/random) — read-list — GET /api/breeds/image/random
- `dog-api call getBreedHoundList` (GET /api/breed/hound/list) — read-list — GET /api/breed/hound/list

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.