# catfact — OpenClaw skill

> catfact should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `catfact`
- **MCP:** `catfact mcp`
- **Base URL:** https://catfact.ninja
- **Operations:** 3 (3 read, 0 write, 0 destructive)

## Auth setup
```bash
export CATFACT_API_KEY=<your-key>
```

## Actions (zero-friction)
- `catfact breeds` — Get a list of breeds
- `catfact fact` — Get Random Fact
- `catfact facts` — Get a list of facts

## All operations (first 200)
- `catfact call getBreeds` (GET /breeds) — read-list — Get a list of breeds
- `catfact call getRandomFact` (GET /fact) — read-list — Get Random Fact
- `catfact call getFacts` (GET /facts) — read-list — Get a list of facts

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.