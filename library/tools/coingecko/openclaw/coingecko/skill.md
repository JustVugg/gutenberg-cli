# coingecko — OpenClaw skill

> coingecko should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `coingecko`
- **MCP:** `coingecko mcp`
- **Base URL:** https://api.coingecko.com
- **Operations:** 4 (4 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `coingecko markets` — GET /api/v3/coins/markets
- `coingecko bitcoin` — GET /api/v3/coins/bitcoin
- `coingecko price` — GET /api/v3/simple/price
- `coingecko global` — GET /api/v3/global

## All operations (first 200)
- `coingecko call getCoinsMarkets` (GET /api/v3/coins/markets) — read-list — GET /api/v3/coins/markets
- `coingecko call getCoinsBitcoin` (GET /api/v3/coins/bitcoin) — read-list — GET /api/v3/coins/bitcoin
- `coingecko call getSimplePrice` (GET /api/v3/simple/price) — read-list — GET /api/v3/simple/price
- `coingecko call getGlobal` (GET /api/v3/global) — read-list — GET /api/v3/global

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.