# hacker-news — OpenClaw skill

> hacker-news should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `hacker-news`
- **MCP:** `hacker-news mcp`
- **Base URL:** https://hacker-news.firebaseio.com
- **Operations:** 5 (5 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `hacker-news topstories-json` — GET /v0/topstories.json
- `hacker-news newstories-json` — GET /v0/newstories.json
- `hacker-news beststories-json` — GET /v0/beststories.json
- `hacker-news 8863-json` — GET /v0/item/8863.json
- `hacker-news maxitem-json` — GET /v0/maxitem.json

## All operations (first 200)
- `hacker-news call getTopstoriesJson` (GET /v0/topstories.json) — read-list — GET /v0/topstories.json
- `hacker-news call getNewstoriesJson` (GET /v0/newstories.json) — read-list — GET /v0/newstories.json
- `hacker-news call getBeststoriesJson` (GET /v0/beststories.json) — read-list — GET /v0/beststories.json
- `hacker-news call getItem8863Json` (GET /v0/item/8863.json) — read-list — GET /v0/item/8863.json
- `hacker-news call getMaxitemJson` (GET /v0/maxitem.json) — read-list — GET /v0/maxitem.json

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.