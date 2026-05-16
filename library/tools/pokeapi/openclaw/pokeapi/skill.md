# pokeapi — OpenClaw skill

> pokeapi should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `pokeapi`
- **MCP:** `pokeapi mcp`
- **Base URL:** https://pokeapi.co
- **Operations:** 4 (4 read, 0 write, 0 destructive)

## Actions (zero-friction)
- `pokeapi pokemon` — GET /api/v2/pokemon
- `pokeapi pikachu` — GET /api/v2/pokemon/pikachu
- `pokeapi type` — GET /api/v2/type
- `pokeapi ability` — GET /api/v2/ability

## All operations (first 200)
- `pokeapi call getPokemon` (GET /api/v2/pokemon) — read-list — GET /api/v2/pokemon
- `pokeapi call getPokemonPikachu` (GET /api/v2/pokemon/pikachu) — read-list — GET /api/v2/pokemon/pikachu
- `pokeapi call getType` (GET /api/v2/type) — read-list — GET /api/v2/type
- `pokeapi call getAbility` (GET /api/v2/ability) — read-list — GET /api/v2/ability

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.