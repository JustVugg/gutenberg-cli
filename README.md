<p align="center">
  <img src="gutenberg-cli.jpg" alt="gutenberg-cli_logo" width="500"/>
</p>

**A verified tool factory for AI agents.** From any API surface to safe, verified tools any agent runtime can use — Go CLI + MCP server + agent skills + SQLite/FTS5 cache + dry-run-by-default policy + hash-verifiable proof artifacts.

```bash
$ gutenberg quick https://catfact.ninja/fact
gutenberg quick: probing https://catfact.ninja/fact
gutenberg quick: verdict = openapi-published (high)
gutenberg quick: downloading spec from https://catfact.ninja/docs?api-docs.json
gutenberg quick: generating tool at library/tools/catfact
gutenberg quick: verifying catfact   ← go build · cli-smoke · mcp-handshake · go test
gutenberg quick: verify ok
gutenberg quick: installing catfact  → ~/.local/bin/catfact

$ catfact fact                       ← Grade A · 100/100 · ~30 s end-to-end
GET https://catfact.ninja/fact
status: 200 OK
{
  "fact": "In just 7 years, one un-spayed female cat and one un-neutered male cat...",
  "length": 121
}
```

**One command. Any public API. A working tool surface for AI agents: CLI, MCP, skills, cache, policy and proofs. Verified by build, cli-smoke, MCP handshake, and Go tests.** Nothing else does this.

## What "quick" does, in 6 numbered steps

1. **probe** the URL (`try`) — classifies it as `openapi-published / json-endpoint / html-content / spa / lazy-rendered / anti-bot-challenge / http-error`
2. **fetch the spec** (or seed a HAR if no spec is exposed; auto-infers baseUrl/schemas from real response bodies)
3. **generate** Go CLI + MCP server + Claude skill + OpenClaw skill
4. **verify** with build + cli smoke + MCP handshake + Go tests → `proofs/verification.json`
5. **install** into `~/.local/bin/<slug>`
6. **expose heroes** — zero-friction aliases auto-detected from the spec (e.g. `catfact fact`, `github meta`, `tvmaze shows`)

## Or pass a natural language intent and route to the catalog

```bash
$ gutenberg quick "top hacker news stories"
gutenberg quick: searching catalog for "top hacker news stories"
gutenberg quick: top match: hacker-news topstories-json (score=5)
[...executes against the cataloged tool...]
```

---

## Why Gutenberg

Most AI-agent infrastructure today gives you a model and tells you "bring your own tools". Gutenberg makes the tools — and makes sure they actually work before any agent calls them.

| Primitive | What it does |
|---|---|
| **`forge` / `generate`** | OpenAPI / GraphQL / HAR / OpenAPI-from-Postman/Insomnia/curl → Go CLI + MCP server + Claude/OpenClaw skill |
| **`verify`** | Runs `go build`, CLI smoke test, MCP handshake, `go test` — writes `proofs/verification.json`. **No proof = no Grade A.** |
| **`try <url>`** | Autopilot: classifies the site (`openapi-published / json-endpoint / html-content / spa / lazy-rendered / anti-bot-challenge / http-error`) and tells you the exact next command |
| **`extract` / `scrape`** | LLM-powered structured extraction (JSON-Schema validated, cache TTL, Anthropic/OpenAI/Ollama provider) + Firecrawl-style markdown scraping |
| **`snapshot/replay`** | Generated tools support deterministic record + offline replay for tests and CI |
| **`vault`** | Shared, AES-GCM-encrypted token vault across all generated tools |
| **`lockfile` + `diff` + `upgrade` + `watch`** | Pin spec hashes, diff semantically, regenerate while preserving `// gutenberg:keep` blocks, poll upstream for drift |
| **`compare`** | Call N tools on the same intent, see results side by side |
| **`aggregator` kind** | One recipe → a fan-out Go CLI over N JSON sources with `merge.go` + `rank.go` |
| **`heroes`** | Zero-friction aliases auto-detected (or curated via `x-gutenberg-hero`) — e.g. `espn nba`, `github meta`, `tvmaze shows` |

## Quick Start

```bash
npm install                                            # install deps (playwright, yaml)
node bin/gutenberg.js doctor                           # check runtime + LLM providers + Browserbase

# Generate from a public OpenAPI
curl -fsSL -o /tmp/github.json \
  https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.json
node bin/gutenberg.js generate /tmp/github.json --out library/tools/github --name github --force --strict
node bin/gutenberg.js verify library/tools/github
node bin/gutenberg.js install library/tools/github     # builds + installs to ~/.local/bin/

github meta                                            # zero-friction hero, returns digest of GitHub API root
github call search/repos --param q=claude-code --select '[*].full_name'
```

Probe a site you've never seen:

```bash
gutenberg try https://www.thetrainline.com/it          # → "spa, use record"
gutenberg try https://it.wikipedia.org/wiki/Bologna    # → "html-content, use scrape/extract"
gutenberg try https://api.tvmaze.com/shows/1           # → "json-endpoint, seed-har + import + generate"
```

## What Gutenberg works on (honest table)

| Site shape | Recommended path | Success rate | Examples we've shipped |
|---|---|---|---|
| **API with public OpenAPI** | `generate` directly | ✅ ~95% | GitHub (1183 ops), Stripe (587), Sentry (209) |
| **API without spec but JSON endpoints** | `seed-har` → `import-har` → `generate` | ✅ ~90% | ESPN, TVMaze, Wikipedia REST, public JSON APIs |
| **Plain HTML content site** | `scrape` (+ `--structured`) or `extract` with schema | ✅ ~80% | Wikipedia and content/documentation pages |
| **SPA with reachable XHR** | `record` (Playwright) → `import-har` | ⚠️ ~40% | needs interactive triggering |
| **SPA + Cloudflare Turnstile / strong anti-bot** | `record --backend browserbase` (paid) — or use an official partner API | ❌ ~5% with free tooling | thetrainline, lastminute, OpenTable, Subito.it |

## Non-goals (what Gutenberg deliberately does NOT do)

- **Bypass anti-bot protections.** Cloudflare Turnstile, JS challenges, fingerprint-based bot detection: Gutenberg does not try to defeat them. When a site invests in stopping bots, that's a product decision we respect. For legitimate browser-driven scraping use the official Browserbase backend (`record --backend browserbase`).
- **Solve CAPTCHAs.** No 2Captcha-style integrations.
- **Scrape against ToS.** The README of every generated tool tells users to respect rate limits, robots.txt, and terms of service.
- **Be a model.** Gutenberg routes LLM calls (Anthropic / OpenAI / Ollama) but does not train or fine-tune.

## Commands

```bash
gutenberg doctor                            Probe runtime, LLM providers, Browserbase
gutenberg try <url>                         Autopilot — classify a site and pick the next command
gutenberg discover <url>                    Look for a published OpenAPI/Swagger
gutenberg seed-har <url> [<url>...]         Fetch URL(s) via plain HTTP, emit HAR (no browser)
gutenberg record <url>                      Playwright HAR recording  [--storage-state state.json]
                                                                       [--backend browserbase --key $KEY --project-id $PROJ]
gutenberg login <url> --out state.json      Headed login flow, save storage-state for later record/extract
gutenberg import-har <har>                  HAR → OpenAPI (with schema inference + auth sniffer)
gutenberg import-graphql <source>           GraphQL SDL / introspection → OpenAPI
gutenberg generate <spec> --out <dir>       OpenAPI → Go CLI + MCP + Claude skill (+ openclaw, optional)
                                            [--strict] [--min-score 0.7] [--default-header 'k: v']
                                            [--policy policy.json]
gutenberg generate <recipe> --kind aggregator --out <dir>
gutenberg verify <project-dir>              Runs build + cli-smoke + mcp-handshake + tests; writes proofs/
gutenberg install <project-dir>             go build + symlink into ~/.local/bin
gutenberg install starter-pack              Install Hacker News, ESPN, Wikipedia, TVMaze, Open-Meteo, GitHub
gutenberg scorecard <project-dir>           Multi-dim scorecard with verification badges
gutenberg publish <tool-dir>|--all          Validate proofs + skills + OpenClaw; sync registry
gutenberg registry sync                     Rebuild library/registry.json from real manifests
gutenberg upgrade <project-dir> [--spec ...] Regenerate, preserve `// gutenberg:keep` blocks, auto go mod tidy
gutenberg watch <project-dir> <spec-url>    Poll spec URL, regenerate on drift
gutenberg diff <old-spec> <new-spec>        Semantic spec diff
gutenberg extract <url> -p PROMPT -s SCHEMA LLM extraction with JSON-Schema validation + --cache 1h
gutenberg scrape <url> [--structured]       Main-content markdown (table detector for repeated patterns)
gutenberg search <intent>                   Catalog discovery: top matching operations
gutenberg run <intent> [--no-llm]           LLM-routed (or top-match) execute against catalog
gutenberg compare <tool-a> <tool-b> --op X  Side-by-side call across tools
gutenberg site build                        Rebuild web/data.js with scorecard + heroes
```

Generated tools share a uniform sub-command surface:

```bash
<tool> operations [--json]                  List operations
<tool> heroes [--json]                      List auto-detected hero aliases
<tool> <hero-alias>                         Shortcut to call --digest
<tool> call <op-id> [options]               Call with full control
  --param k=v --path k=v --header 'k: v' --data '{...}' --yes
  --json | --compact | --digest | --select 'data.events[*].name'
  --stream                                  For SSE/NDJSON endpoints
<tool> walk <op-id> --max N                 Auto-pagination (offset/limit, cursor, page)
<tool> sync                                 Cache read operations locally (SQLite/FTS5)
<tool> search <query>                       Full-text over cached responses
<tool> resources [query]                    Search projected resource rows in SQLite
<tool> auth <status|client-credentials|device|pkce-start|pkce-finish|refresh>
<tool> mcp                                  Run as MCP stdio server
```

## Catalog included

21 verified catalog entries with `proofs/verification.json`:

| Tool | Ops | Source | Verified checks |
|---|---:|---|---|
| github | 1183 | GitHub OpenAPI ufficiale | build • cli • mcp • tests |
| stripe | 587 | Stripe OpenAPI ufficiale | build • cli • mcp • tests |
| sentry | 209 | Sentry OpenAPI ufficiale | build • cli • mcp • tests |
| espn | 4 | HAR-seeded | build • cli • mcp • tests |
| wikipedia-it | 5 | HAR-seeded (User-Agent default) | build • cli • mcp • tests |
| tvmaze | 5 | HAR-seeded | build • cli • mcp • tests |
| multi-demo | 2 sources | aggregator (ESPN + Open-Meteo) | build • cli • tests |
| dev-pulse | 2 sources | aggregator (GitHub + ESPN) | build • cli • tests |

Plus CatFact, CoinGecko, Dog API, Exchange Rates, Hacker News, npm registry, Open-Meteo, PokeAPI, public holidays, PyPI, RandomUser, REST Countries, and TVMaze quick. The complete source of truth is `library/registry.json`, generated by `gutenberg registry sync`.

Plus 16+ recipe templates in `library/recipes/` for common sources (Hacker News, CoinGecko, PokeAPI, Open-Meteo, REST Countries, Nominatim, npm registry, PyPI, exchange rates, public holidays, Reddit, Spotify, Slack, …).

## Output targets

`--targets go,mcp,skill,openclaw` controls what gets emitted next to the Go module:

- **`go`** + **`mcp`** — Go binary + stdio MCP server (default)
- **`skill`** — `skills/<slug>/SKILL.md` (Claude Code skill format)
- **`openclaw`** — `openclaw/<slug>/skill.{json,md}` (agent-framework-neutral schema)

## Architecture diff vs other open-source factories

Compared to **site2cli** (PyPI) and **printingpress.dev** (~82 CLIs, npm), Gutenberg uniquely ships:

- Verification proofs gate (no Grade A without `proofs/verification.json`)
- `gutenberg try <url>` autopilot
- Aggregator-as-recipe (declarative)
- Snapshot/replay deterministic mode (built into every generated tool)
- Shared AES-GCM token vault
- `lockfile` + `diff` + `upgrade` + `watch`
- `cross-tool compare`
- LLM-provider-agnostic extract (Anthropic / OpenAI / Ollama) + cache TTL
- HAR → OpenAPI with **response schema inference** and **auth scheme sniffer**
- Built-in JSONPath `--select` and `--stream` for SSE/NDJSON

What we don't have yet: their catalog scale (82 vs 21) and packaged npm distribution. The catalog grows with recipes and `gutenberg publish --all --fix-assets`.

## Run tests

```bash
node --test tests/*.test.js
```

Local-network-dependent tests skip cleanly in sandboxes that deny `127.0.0.1` listeners. Browser tests remain opt-in (`GUTENBERG_BROWSER_TESTS=1` to exercise live Playwright).

## Environment variables

| Var | What |
|---|---|
| `ANTHROPIC_API_KEY` | Use Anthropic for `extract` / `run` |
| `OPENAI_API_KEY` | Use OpenAI for `extract` / `run` |
| `OLLAMA_HOST` (default `http://localhost:11434`) + `GUTENBERG_LLM_PROVIDER=ollama` | Use local or remote Ollama for `extract` |
| `BROWSERBASE_API_KEY` + `BROWSERBASE_PROJECT_ID` | Required for `record --backend browserbase` |
| `GUTENBERG_VAULT_KEY` (32-byte hex) | Encrypt the shared OAuth vault with AES-GCM |
| `GUTENBERG_TELEMETRY=1` + `GUTENBERG_TELEMETRY_FILE` | Opt-in local JSONL log of operation calls (never sent anywhere) |
| `<TOOL_ENV>_SNAPSHOT_MODE=record\|replay` + `<TOOL_ENV>_SNAPSHOT_DIR` | Deterministic HTTP record/replay per tool |
| `<TOOL_ENV>_AUTO_REFRESH=0` | Disable automatic OAuth refresh in a generated tool |

## License

MIT.

## Status

`v0.1.0` — public. Honest about its boundaries:

- ✅ Verified end-to-end on GitHub, Stripe, Sentry, ESPN, Wikipedia IT, TVMaze and the public-API starter catalog
- ✅ 21 catalog entries shipped, 16+ recipe templates
- ❌ Cannot bypass Cloudflare Turnstile / strong anti-bot. Use `record --backend browserbase` or an official partner API for that class of site.
