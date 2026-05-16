# tvmaze — OpenClaw skill

> tvmaze should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `tvmaze`
- **MCP:** `tvmaze mcp`
- **Base URL:** https://api.tvmaze.com
- **Operations:** 5 (5 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `tvmaze shows` — GET /search/shows
- `tvmaze schedule` — GET /schedule

## All operations (first 200)
- `tvmaze call getSearchShows` (GET /search/shows) — search — GET /search/shows
- `tvmaze call getShows` (GET /shows/{id}) — read-detail — GET /shows/{id}
- `tvmaze call getShowsIdEpisodes` (GET /shows/{id}/episodes) — read — GET /shows/{id}/episodes
- `tvmaze call getSchedule` (GET /schedule) — read-list — GET /schedule
- `tvmaze call getPeople` (GET /people/{id}) — read-detail — GET /people/{id}

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.