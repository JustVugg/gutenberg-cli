# Multi-source demo aggregator — OpenClaw skill

> Real proof of the aggregator pattern: fan-out over ESPN scoreboards and Open-Meteo weather. Both sources are public, no auth. Generated to validate that gutenberg can fan-out N OpenAPI sources, merge results, and ship a compiling Go CLI.

- **CLI:** `multi-demo`
- **Base URL:** (none)
- **Operations:** 2 (0 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `multi-demo fan-out` — espn
- `multi-demo fan-out` — open-meteo

## All operations (first 200)
- `multi-demo fan-out` (GET getApisSiteV2SportsBasketballNbaScoreboard) — source — espn
- `multi-demo fan-out` (GET getV1Forecast) — source — open-meteo

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.