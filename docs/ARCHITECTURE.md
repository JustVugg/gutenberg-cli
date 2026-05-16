# Architecture

Gutenberg uses a pipeline with explicit artifacts.

1. Ingest an API source: OpenAPI JSON/YAML, GraphQL, discovered website OpenAPI, or browser HAR capture.
2. Normalize it into a Gutenberg blueprint.
3. Generate a package from the blueprint.
4. Score the generated package.
5. Publish metadata into the registry.
6. Present the registry on the website.

## Core Concepts

### Blueprint

The blueprint is the stable internal representation. It contains operations, parameters, auth hints, risk classification, cacheability and product insights.

### Generated Package

Each generated package contains:

- `gutenberg.manifest.json`
- `go.mod`
- `cmd/<tool>/main.go`
- `internal/forge/client.go`
- `internal/forge/mcp.go`
- `internal/forge/store.go` with SQLite/FTS5
- `internal/forge/auth.go` with OAuth helpers
- `internal/forge/manifest.go`
- `internal/forge/forge_test.go`
- `docs/COOKBOOK.md`

### Risk Model

Read operations run normally. Write and destructive operations return a dry-run plan until `--yes` is passed.

### Store Model

The dependency-free MVP uses a local JSON cache. The intended production adapter is SQLite with FTS5 and optional vector search.

## Command Flow

```text
OpenAPI JSON/YAML / GraphQL / HAR / website discovery
  -> buildBlueprint()
  -> generateProject()
  -> generated Go CLI/MCP/cache package
  -> scoreProject()
  -> registry.json
  -> web catalog
```
