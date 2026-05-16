# Comparison

Printing Press proved that API tools should be generated for agents, not only for humans. Gutenberg follows that direction and makes the product more explicit.

## Gutenberg Differences

- The blueprint is a first-class artifact.
- Generated tools are Go-first.
- Scorecard is part of the normal workflow.
- Generated write operations are guarded by default.
- The web catalog is included in the core repo.
- The generated MCP server and CLI share the same client module.
- Browser HAR import is included.
- Browser recording is included through Playwright.
- GraphQL import is included.
- OAuth helper commands are generated.
- SQLite/FTS5 search is generated.
- The generated Go package has no external runtime dependencies in the MVP.

## Non-Goals

- It is not a scraper-first product.
- It does not bypass API permissions.
- It does not pretend generated tools are production-ready without review.
