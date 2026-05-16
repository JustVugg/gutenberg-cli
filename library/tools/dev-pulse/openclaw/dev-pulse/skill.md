# Dev Pulse aggregator — OpenClaw skill

> Second real proof of the aggregator pattern: fan-out over GitHub REST and ESPN scoreboards. Demonstrates that the same generator handles wildly different source shapes.

- **CLI:** `dev-pulse`
- **Base URL:** (none)
- **Operations:** 2 (0 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `dev-pulse fan-out` — github
- `dev-pulse fan-out` — espn

## All operations (first 200)
- `dev-pulse fan-out` (GET meta/root) — source — github
- `dev-pulse fan-out` (GET getApisSiteV2SportsBasketballNbaScoreboard) — source — espn

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.