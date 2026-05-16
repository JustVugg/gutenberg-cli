# Hero Commands

Gutenberg has two layers:

- Factory commands create Go CLI/MCP packages from OpenAPI, GraphQL, HAR, or websites.
- Hero commands are curated shortcuts with compact output for agents.

## Sports

```bash
gutenberg sports nba today
gutenberg sports serie-a today
gutenberg sports serie-a today --team juventus
```

The sports adapter uses ESPN public scoreboard JSON and returns compact event rows: time, teams, status, score when present, broadcast, tickets, and venue.

## Travel

```bash
gutenberg travel rom-par june --adults 1 --currency EUR
gutenberg travel rom-par june --web
gutenberg travel rom-par --date 2026-06-15 --return 2026-06-20
```

The travel adapter plans an Amadeus Flight Offers request. If `AMADEUS_ACCESS_TOKEN` is set, it calls Amadeus and summarizes offers. If no token is set, it prints the exact request, browser search URLs, and a `record` command for HAR capture.

## Recipes

```bash
gutenberg recipes list
gutenberg recipes show travel-amadeus-flight-offers
gutenberg recipes scaffold my-tool --kind openapi --source openapi.json
```

Recipes are intentionally plain JSON. A generated recipe records source, kind, intended output, and the next commands needed to turn that source into a Go CLI/MCP package.

## Create

```bash
gutenberg create my-site --from https://example.com --kind har
gutenberg create my-api --from openapi.json --kind openapi
gutenberg create my-graphql --from schema.graphql --kind graphql
```

`create` scaffolds a project brief and recipe. It does not pretend every website is automatically stable. If a site exposes a real API or a visible network API, the HAR path can generate a tool; if it requires credentials or blocks automation, the resulting recipe documents that boundary.
