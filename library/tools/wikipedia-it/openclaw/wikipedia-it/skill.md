# wikipedia-it — OpenClaw skill

> wikipedia-it should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `wikipedia-it`
- **MCP:** `wikipedia-it mcp`
- **Base URL:** https://it.wikipedia.org
- **Operations:** 5 (5 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `wikipedia-it bologna` — GET /api/rest_v1/page/summary/Bologna
- `wikipedia-it roma` — GET /api/rest_v1/page/summary/Roma
- `wikipedia-it italia` — GET /api/rest_v1/page/related/Italia
- `wikipedia-it summary` — GET /api/rest_v1/page/random/summary

## All operations (first 200)
- `wikipedia-it call getPageSummaryBologna` (GET /api/rest_v1/page/summary/Bologna) — read-list — GET /api/rest_v1/page/summary/Bologna
- `wikipedia-it call getPageSummaryRoma` (GET /api/rest_v1/page/summary/Roma) — read-list — GET /api/rest_v1/page/summary/Roma
- `wikipedia-it call getPageRelatedItalia` (GET /api/rest_v1/page/related/Italia) — read-list — GET /api/rest_v1/page/related/Italia
- `wikipedia-it call getPageRandomSummary` (GET /api/rest_v1/page/random/summary) — read-list — GET /api/rest_v1/page/random/summary
- `wikipedia-it call getFeedFeatured` (GET /api/rest_v1/feed/featured/{year}/{month}/{day}) — read-detail — GET /api/rest_v1/feed/featured/{year}/{month}/{day}

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.