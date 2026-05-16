# Gutenberg Product Spec

Gutenberg is an agent-native tool factory. It turns API descriptions, GraphQL schemas, and browser HAR captures into complete Go tool packages: CLI, MCP server, SQLite/FTS5 cache, OAuth helpers, tests, scorecard, registry metadata and web catalog entries.

## Better Than Printing Press

Gutenberg is designed around these improvements:

- A stable intermediate blueprint before code generation.
- Go-first generated tools, with Node used only for the factory itself.
- One shared runtime core for CLI and MCP instead of two separate implementations.
- Risk-aware execution: write and destructive operations dry-run unless explicitly confirmed.
- First-class quality gates: every generated tool ships with tests and can be scored.
- Product catalog from day one: generated tools are meant to be listed, compared and installed.
- Browser HAR capture and import.
- GraphQL import with query/mutation risk classification.
- OAuth token flows and local token storage.

## MVP Scope

- Input: OpenAPI JSON/YAML, browser HAR export, GraphQL SDL/introspection, or discovered website OpenAPI.
- Output: Go-based generated package with SQLite/FTS5 and OAuth helpers.
- Surfaces: CLI, MCP stdio server, local JSON cache, manifest, tests and docs.
- Catalog: registry JSON plus static web UI.

## Expansion Scope

- OAuth authorization-code PKCE browser flow.
- Hosted registry API.
- Hosted registry API.
- Installer package.
- Generated skill packs for Codex, Claude Code and other agent shells.
