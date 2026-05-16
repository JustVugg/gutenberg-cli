# public-holidays — OpenClaw skill

> public-holidays should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `public-holidays`
- **MCP:** `public-holidays mcp`
- **Base URL:** https://date.nager.at
- **Operations:** 4 (4 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `public-holidays availablecountries` — GET /api/v3/AvailableCountries
- `public-holidays it` — GET /api/v3/NextPublicHolidays/IT

## All operations (first 200)
- `public-holidays call getPublicholidaysIdIt` (GET /api/v3/PublicHolidays/{id}/IT) — read — GET /api/v3/PublicHolidays/{id}/IT
- `public-holidays call getPublicholidaysIdUs` (GET /api/v3/PublicHolidays/{id}/US) — read — GET /api/v3/PublicHolidays/{id}/US
- `public-holidays call getAvailablecountries` (GET /api/v3/AvailableCountries) — read-list — GET /api/v3/AvailableCountries
- `public-holidays call getNextpublicholidaysIt` (GET /api/v3/NextPublicHolidays/IT) — read-list — GET /api/v3/NextPublicHolidays/IT

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.