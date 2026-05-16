# espn — OpenClaw skill

> espn should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `espn`
- **MCP:** `espn mcp`
- **Base URL:** https://site.api.espn.com
- **Operations:** 4 (4 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `espn nba` — GET /apis/site/v2/sports/basketball/nba/scoreboard
- `espn ita-1` — GET /apis/site/v2/sports/soccer/ita.1/scoreboard
- `espn nfl` — GET /apis/site/v2/sports/football/nfl/scoreboard
- `espn mlb` — GET /apis/site/v2/sports/baseball/mlb/scoreboard

## All operations (first 200)
- `espn call getBasketballNbaScoreboard` (GET /apis/site/v2/sports/basketball/nba/scoreboard) — read-list — GET /apis/site/v2/sports/basketball/nba/scoreboard
- `espn call getSoccerIta1Scoreboard` (GET /apis/site/v2/sports/soccer/ita.1/scoreboard) — read-list — GET /apis/site/v2/sports/soccer/ita.1/scoreboard
- `espn call getFootballNflScoreboard` (GET /apis/site/v2/sports/football/nfl/scoreboard) — read-list — GET /apis/site/v2/sports/football/nfl/scoreboard
- `espn call getBaseballMlbScoreboard` (GET /apis/site/v2/sports/baseball/mlb/scoreboard) — read-list — GET /apis/site/v2/sports/baseball/mlb/scoreboard

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.