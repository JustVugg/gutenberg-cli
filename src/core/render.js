import path from "node:path";
import { assertCanWriteDirectory, writeJson, writeText } from "./fs.js";
import { camelCase, pascalCase, slugify } from "./sanitize.js";
import { generateGoProject } from "./render-go.js";
import { generateClaudeSkill } from "./render-skill.js";
import { generateOpenClawSkill } from "./render-openclaw.js";
import { attachScorecard, computeProvenance } from "./provenance.js";
import { scoreProject } from "./scorecard.js";
import { compareLock, readLock, writeLock } from "./lockfile.js";
import { defaultPolicy } from "./source.js";

const DEFAULT_TARGETS = ["go", "mcp", "skill"];

function parseHeaderList(value) {
  if (!value) return {};
  const list = Array.isArray(value) ? value : [value];
  const out = {};
  for (const item of list) {
    const text = String(item);
    const colon = text.indexOf(":");
    if (colon === -1) continue;
    const name = text.slice(0, colon).trim();
    const val = text.slice(colon + 1).trim();
    if (name) out[name] = val;
  }
  return out;
}

export function parseTargets(value) {
  if (!value) return DEFAULT_TARGETS;
  const list = Array.isArray(value) ? value : String(value).split(",");
  const normalized = list.map((item) => String(item).trim().toLowerCase()).filter(Boolean);
  return normalized.length > 0 ? normalized : DEFAULT_TARGETS;
}

export function generateProject(blueprint, outDir, options = {}) {
  const targets = parseTargets(options.targets);
  const lang = options.lang || (targets.includes("go") ? "go" : "node");
  const provenance = computeProvenance({
    specPath: options.specPath || blueprint?.source || null,
    recipePath: options.recipePath || null,
    targets,
    name: options.name || blueprint?.slug
  });
  const defaultHeaders = parseHeaderList(options.defaultHeaders);
  const enrichedBlueprint = { ...blueprint, provenance, defaultHeaders, policy: options.policy || blueprint.policy || defaultPolicy() };

  const previousLock = readLock(outDir);
  const drift = compareLock(previousLock, provenance);

  const result = lang === "go"
    ? generateGoProject(enrichedBlueprint, outDir, options)
    : generateNodeProject(enrichedBlueprint, outDir, options);
  result.targets = [...targets];
  if (targets.includes("skill")) {
    const skill = generateClaudeSkill(result.manifest, outDir);
    result.skill = skill;
  }
  if (targets.includes("openclaw")) {
    const openclaw = generateOpenClawSkill(result.manifest, outDir);
    result.openclaw = openclaw;
  }

  let scorecard = null;
  try {
    scorecard = scoreProject(outDir);
  } catch {
    scorecard = null;
  }
  const finalProvenance = attachScorecard(provenance, scorecard);
  result.manifest.provenance = finalProvenance;
  writeJson(path.join(outDir, "gutenberg.manifest.json"), result.manifest);
  result.provenance = finalProvenance;
  result.scorecard = scorecard;
  result.lock = writeLock(outDir, finalProvenance);
  result.drift = drift;
  return result;
}

export function generateNodeProject(blueprint, outDir, options = {}) {
  const slug = slugify(options.name || blueprint.slug);
  const displayName = options.displayName || blueprint.name;
  const env = blueprint.envPrefix;
  const binName = slug;
  const manifest = {
    ...blueprint,
    name: displayName,
    slug,
    packageName: `${slug}-gutenberg`,
    generatedBy: "gutenberg",
    generatedByVersion: "0.1.0"
  };

  assertCanWriteDirectory(outDir, Boolean(options.force));

  writeJson(path.join(outDir, "gutenberg.manifest.json"), manifest);
  writeJson(path.join(outDir, "blackforge.manifest.json"), manifest);
  writeJson(path.join(outDir, "package.json"), generatedPackageJson(manifest));
  writeText(path.join(outDir, ".env.example"), `${env}_API_KEY=\n${env}_BASE_URL=${manifest.baseUrls[0] || "https://api.example.com"}\n`);
  writeText(path.join(outDir, "README.md"), generatedReadme(manifest));
  writeText(path.join(outDir, "docs", "COOKBOOK.md"), generatedCookbook(manifest));
  writeText(path.join(outDir, "bin", `${binName}.js`), generatedBin(manifest));
  writeText(path.join(outDir, "src", "client.js"), generatedClient(manifest));
  writeText(path.join(outDir, "src", "cli.js"), generatedCli(manifest));
  writeText(path.join(outDir, "src", "store.js"), generatedStore(manifest));
  writeText(path.join(outDir, "src", "mcp-server.js"), generatedMcpServer(manifest));
  writeText(path.join(outDir, "tests", "smoke.test.js"), generatedSmokeTest(manifest));

  return {
    outDir,
    manifest,
    entrypoint: path.join(outDir, "bin", `${binName}.js`)
  };
}

function generatedPackageJson(manifest) {
  return {
    name: manifest.packageName,
    version: "0.1.0",
    type: "module",
    private: true,
    description: `Gutenberg generated agent-native CLI and MCP server for ${manifest.name}.`,
    bin: {
      [manifest.slug]: `./bin/${manifest.slug}.js`
    },
    scripts: {
      test: "node --test tests/*.test.js",
      mcp: "node src/mcp-server.js",
      operations: `node bin/${manifest.slug}.js operations`
    },
    engines: {
      node: ">=20"
    }
  };
}

function generatedBin(manifest) {
  return `#!/usr/bin/env node
import { main } from "../src/cli.js";

main(process.argv.slice(2)).catch((error) => {
  console.error("${manifest.slug}:", error && error.message ? error.message : String(error));
  process.exitCode = error && Number.isInteger(error.exitCode) ? error.exitCode : 1;
});
`;
}

function generatedClient(manifest) {
  const manifestLiteral = JSON.stringify(manifest, null, 2);
  return `const manifest = ${manifestLiteral};

export { manifest };

export function listOperations() {
  return manifest.operations;
}

export function getOperation(operationId) {
  const operation = manifest.operations.find((item) => item.id === operationId);
  if (!operation) {
    const error = new Error("Unknown operation: " + operationId);
    error.exitCode = 2;
    throw error;
  }
  return operation;
}

export function buildUrl(operation, options = {}) {
  const baseUrl = String(options.baseUrl || process.env["${manifest.envPrefix}_BASE_URL"] || manifest.baseUrls[0] || "").replace(/\\/$/, "");
  if (!baseUrl) {
    const error = new Error("Missing base URL. Pass --base-url or set ${manifest.envPrefix}_BASE_URL.");
    error.exitCode = 2;
    throw error;
  }

  const pathParams = options.pathParams || {};
  const queryParams = options.queryParams || {};
  let apiPath = operation.path;
  for (const parameter of operation.parameters.filter((item) => item.in === "path")) {
    const value = pathParams[parameter.name] ?? queryParams[parameter.name];
    if (value === undefined || value === null || value === "") {
      const error = new Error("Missing path parameter: " + parameter.name);
      error.exitCode = 2;
      throw error;
    }
    apiPath = apiPath.replace("{" + parameter.name + "}", encodeURIComponent(String(value)));
  }

  const url = new URL(baseUrl + apiPath);
  for (const parameter of operation.parameters.filter((item) => item.in === "query")) {
    const value = queryParams[parameter.name];
    if (value !== undefined && value !== null && value !== "") {
      url.searchParams.set(parameter.name, String(value));
    } else if (parameter.required) {
      const error = new Error("Missing query parameter: " + parameter.name);
      error.exitCode = 2;
      throw error;
    }
  }
  return url.toString();
}

export function authHeaders(options = {}) {
  const headers = {};
  const apiKey = options.apiKey || process.env["${manifest.envPrefix}_API_KEY"];
  if (!apiKey || manifest.auth.mode === "none") {
    return headers;
  }

  const scheme = manifest.auth.schemes[0] || {};
  if (scheme.type === "http" && scheme.scheme === "bearer") {
    headers.Authorization = "Bearer " + apiKey;
  } else if (scheme.type === "apiKey" && scheme.in === "header" && scheme.header) {
    headers[scheme.header] = apiKey;
  } else {
    headers.Authorization = "Bearer " + apiKey;
  }
  return headers;
}

export async function callOperation(operationId, options = {}) {
  const operation = getOperation(operationId);
  const url = buildUrl(operation, options);
  const headers = {
    Accept: "application/json",
    ...authHeaders(options),
    ...(options.headers || {})
  };

  let body;
  if (options.body !== undefined && options.body !== null) {
    headers["Content-Type"] = "application/json";
    body = typeof options.body === "string" ? options.body : JSON.stringify(options.body);
  }

  if (operation.risk !== "read" && !options.yes) {
    return {
      dryRun: true,
      operation,
      request: {
        method: operation.method,
        url,
        headers: redactHeaders(headers),
        body: body ? JSON.parse(body) : undefined
      },
      note: "Write and destructive operations require yes=true or --yes."
    };
  }

  const controller = new AbortController();
  const timeout = Number(options.timeoutMs || 30000);
  const timer = setTimeout(() => controller.abort(), timeout);
  try {
    const response = await fetch(url, {
      method: operation.method,
      headers,
      body,
      signal: controller.signal
    });
    const text = await response.text();
    const data = parseResponse(text, response.headers.get("content-type"));
    return {
      dryRun: false,
      operation,
      request: {
        method: operation.method,
        url,
        headers: redactHeaders(headers)
      },
      response: {
        ok: response.ok,
        status: response.status,
        statusText: response.statusText,
        data
      }
    };
  } finally {
    clearTimeout(timer);
  }
}

function parseResponse(text, contentType) {
  if (!text) {
    return null;
  }
  if (contentType && contentType.includes("application/json")) {
    try {
      return JSON.parse(text);
    } catch {
      return text;
    }
  }
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

function redactHeaders(headers) {
  const redacted = {};
  for (const [key, value] of Object.entries(headers)) {
    redacted[key] = /authorization|token|secret|api[-_]?key|subscription[-_]?key|ocp-apim-subscription-key|^key$/i.test(key) ? "[redacted]" : value;
  }
  return redacted;
}
`;
}

function generatedStore(manifest) {
  return `import fs from "node:fs";
import path from "node:path";

const cacheFile = process.env["${manifest.envPrefix}_CACHE_FILE"] || path.join(process.cwd(), ".gutenberg", "${manifest.slug}-cache.json");

export function readCache() {
  if (!fs.existsSync(cacheFile)) {
    return { version: 1, records: [] };
  }
  return JSON.parse(fs.readFileSync(cacheFile, "utf8"));
}

export function writeCache(cache) {
  fs.mkdirSync(path.dirname(cacheFile), { recursive: true });
  fs.writeFileSync(cacheFile, JSON.stringify(cache, null, 2) + "\\n", "utf8");
}

export function saveRecord(record) {
  const cache = readCache();
  cache.records.push({
    ...record,
    cachedAt: new Date().toISOString()
  });
  writeCache(cache);
  return cache.records[cache.records.length - 1];
}

export function searchCache(query) {
  const needle = String(query || "").toLowerCase();
  return readCache().records.filter((record) => JSON.stringify(record).toLowerCase().includes(needle));
}

export function cacheStats() {
  const cache = readCache();
  const byOperation = {};
  for (const record of cache.records) {
    byOperation[record.operationId] = (byOperation[record.operationId] || 0) + 1;
  }
  return {
    file: cacheFile,
    records: cache.records.length,
    byOperation
  };
}
`;
}

function generatedCli(manifest) {
  return `import { spawn } from "node:child_process";
import { fileURLToPath } from "node:url";
import path from "node:path";
import { callOperation, listOperations, manifest } from "./client.js";
import { cacheStats, saveRecord, searchCache } from "./store.js";

export async function main(argv) {
  const command = argv[0] || "help";
  const parsed = parseArgs(argv.slice(1));

  if (command === "help" || command === "--help" || command === "-h") {
    printHelp();
    return;
  }
  if (command === "info") {
    printJsonOrText(manifest, parsed.options.json);
    return;
  }
  if (command === "operations") {
    const operations = listOperations();
    if (parsed.options.json) {
      console.log(JSON.stringify(operations, null, 2));
      return;
    }
    printOperations(operations);
    return;
  }
  if (command === "call") {
    const operationId = parsed.positionals[0] || parsed.options.operation;
    if (!operationId) {
      throw usage("Missing operation id. Example: ${manifest.slug} call listPets");
    }
    const result = await callOperation(operationId, {
      baseUrl: parsed.options["base-url"],
      apiKey: parsed.options["api-key"],
      yes: Boolean(parsed.options.yes),
      pathParams: objectOptions(parsed.options.path || []),
      queryParams: objectOptions(parsed.options.param || []),
      body: parsed.options.data ? JSON.parse(parsed.options.data) : undefined
    });
    if (parsed.options.cache && result.response) {
      saveRecord({ operationId, request: result.request, response: result.response });
    }
    printJsonOrText(result, parsed.options.json || result.dryRun);
    return;
  }
  if (command === "sync") {
    const selected = parsed.positionals[0];
    const operations = listOperations().filter((operation) => operation.cacheable && (!selected || operation.id === selected));
    const records = [];
    for (const operation of operations) {
      const result = await callOperation(operation.id, {
        baseUrl: parsed.options["base-url"],
        apiKey: parsed.options["api-key"],
        queryParams: objectOptions(parsed.options.param || [])
      });
      if (result.response) {
        records.push(saveRecord({ operationId: operation.id, request: result.request, response: result.response }));
      }
    }
    printJsonOrText({ synced: records.length, records }, parsed.options.json);
    return;
  }
  if (command === "search") {
    const query = parsed.positionals.join(" ");
    const results = searchCache(query);
    printJsonOrText({ query, results }, parsed.options.json);
    return;
  }
  if (command === "cache") {
    printJsonOrText(cacheStats(), parsed.options.json);
    return;
  }
  if (command === "mcp") {
    const dirname = path.dirname(fileURLToPath(import.meta.url));
    const child = spawn(process.execPath, [path.join(dirname, "mcp-server.js")], { stdio: "inherit" });
    child.on("exit", (code) => { process.exitCode = code || 0; });
    return;
  }

  throw usage("Unknown command: " + command);
}

function printHelp() {
  console.log(\`${manifest.name} (${manifest.slug})

Agent-native CLI generated by Gutenberg.

Commands:
  help                         Show this help
  info [--json]                Show manifest and product thesis
  operations [--json]          List API operations
  call <operation> [options]   Call an API operation
  sync [operation] [options]   Cache read operations locally
  search <query> [--json]      Search cached API responses
  cache [--json]               Show cache stats
  mcp                          Start the MCP stdio server

Options for call/sync:
  --base-url <url>             Override base URL
  --api-key <key>              Override ${manifest.envPrefix}_API_KEY
  --param name=value           Query parameter, repeatable
  --path name=value            Path parameter, repeatable
  --data '{...}'               JSON body
  --cache                      Store call result in local cache
  --yes                        Execute write/destructive operations
  --json                       Machine-readable output
\`);
}

function parseArgs(args) {
  const options = {};
  const positionals = [];
  for (let index = 0; index < args.length; index += 1) {
    const item = args[index];
    if (!item.startsWith("--")) {
      positionals.push(item);
      continue;
    }
    const equals = item.indexOf("=");
    const key = item.slice(2, equals === -1 ? undefined : equals);
    const value = equals === -1 ? args[index + 1] : item.slice(equals + 1);
    if (equals === -1 && (value === undefined || value.startsWith("--"))) {
      options[key] = true;
      continue;
    }
    if (equals === -1) {
      index += 1;
    }
    if (options[key] === undefined) {
      options[key] = value;
    } else if (Array.isArray(options[key])) {
      options[key].push(value);
    } else {
      options[key] = [options[key], value];
    }
  }
  return { options, positionals };
}

function objectOptions(values) {
  const list = Array.isArray(values) ? values : values ? [values] : [];
  const output = {};
  for (const item of list) {
    const equals = String(item).indexOf("=");
    if (equals === -1) {
      output[item] = "true";
    } else {
      output[item.slice(0, equals)] = item.slice(equals + 1);
    }
  }
  return output;
}

function printOperations(operations) {
  const rows = operations.map((operation) => [
    operation.id,
    operation.method,
    operation.path,
    operation.risk,
    operation.summary
  ]);
  printTable(["id", "method", "path", "risk", "summary"], rows);
}

function printTable(headers, rows) {
  const widths = headers.map((header, index) => Math.max(header.length, ...rows.map((row) => String(row[index] || "").length)));
  console.log(headers.map((header, index) => header.padEnd(widths[index])).join("  "));
  console.log(widths.map((width) => "-".repeat(width)).join("  "));
  for (const row of rows) {
    console.log(row.map((cell, index) => String(cell || "").padEnd(widths[index])).join("  "));
  }
}

function printJsonOrText(value, json) {
  if (json) {
    console.log(JSON.stringify(value, null, 2));
    return;
  }
  if (value.response) {
    console.log(JSON.stringify(value.response.data, null, 2));
    return;
  }
  console.log(JSON.stringify(value, null, 2));
}

function usage(message) {
  const error = new Error(message);
  error.exitCode = 2;
  return error;
}
`;
}

function generatedMcpServer(manifest) {
  return `import { callOperation, listOperations, manifest } from "./client.js";
import { searchCache } from "./store.js";

const tools = [
  {
    name: "${manifest.slug}_operations",
    description: "List available ${manifest.name} API operations.",
    inputSchema: { type: "object", properties: {}, additionalProperties: false }
  },
  {
    name: "${manifest.slug}_call",
    description: "Call a ${manifest.name} API operation. Write/destructive operations dry-run unless yes is true.",
    inputSchema: {
      type: "object",
      required: ["operationId"],
      properties: {
        operationId: { type: "string" },
        baseUrl: { type: "string" },
        params: { type: "object" },
        path: { type: "object" },
        body: {},
        yes: { type: "boolean" }
      }
    }
  },
  {
    name: "${manifest.slug}_search_cache",
    description: "Search cached ${manifest.name} API responses.",
    inputSchema: {
      type: "object",
      required: ["query"],
      properties: { query: { type: "string" } }
    }
  }
];

export function startServer() {
  let buffer = Buffer.alloc(0);
  process.stdin.on("data", async (chunk) => {
    buffer = Buffer.concat([buffer, chunk]);
    while (true) {
      const parsed = readFrame(buffer);
      if (!parsed) {
        break;
      }
      buffer = parsed.rest;
      try {
        const response = await handleMessage(parsed.message);
        if (response) {
          writeFrame(response);
        }
      } catch (error) {
        writeFrame({
          jsonrpc: "2.0",
          id: parsed.message.id || null,
          error: { code: -32000, message: error.message || String(error) }
        });
      }
    }
  });
}

async function handleMessage(message) {
  if (message.method === "initialize") {
    return {
      jsonrpc: "2.0",
      id: message.id,
      result: {
        protocolVersion: "2024-11-05",
        capabilities: { tools: {} },
        serverInfo: { name: "${manifest.slug}-mcp", version: "0.1.0" }
      }
    };
  }
  if (message.method === "tools/list") {
    return { jsonrpc: "2.0", id: message.id, result: { tools } };
  }
  if (message.method === "tools/call") {
    const name = message.params?.name;
    const args = message.params?.arguments || {};
    if (name === "${manifest.slug}_operations") {
      return toolResult(message.id, listOperations());
    }
    if (name === "${manifest.slug}_search_cache") {
      return toolResult(message.id, searchCache(args.query));
    }
    if (name === "${manifest.slug}_call") {
      const result = await callOperation(args.operationId, {
        baseUrl: args.baseUrl,
        queryParams: args.params || {},
        pathParams: args.path || {},
        body: args.body,
        yes: Boolean(args.yes)
      });
      return toolResult(message.id, result);
    }
    return {
      jsonrpc: "2.0",
      id: message.id,
      error: { code: -32602, message: "Unknown tool: " + name }
    };
  }
  if (message.method === "notifications/initialized") {
    return null;
  }
  return { jsonrpc: "2.0", id: message.id, error: { code: -32601, message: "Method not found" } };
}

function toolResult(id, value) {
  return {
    jsonrpc: "2.0",
    id,
    result: {
      content: [{ type: "text", text: JSON.stringify(value, null, 2) }]
    }
  };
}

function readFrame(buffer) {
  const separator = buffer.indexOf("\\r\\n\\r\\n");
  if (separator === -1) {
    const newline = buffer.indexOf("\\n");
    if (newline === -1) {
      return null;
    }
    const line = buffer.subarray(0, newline).toString("utf8").trim();
    if (!line.startsWith("{")) {
      return null;
    }
    return {
      message: JSON.parse(line),
      rest: buffer.subarray(newline + 1)
    };
  }
  const header = buffer.subarray(0, separator).toString("utf8");
  const match = header.match(/Content-Length: (\\d+)/i);
  if (!match) {
    throw new Error("Missing Content-Length header");
  }
  const length = Number(match[1]);
  const start = separator + 4;
  const end = start + length;
  if (buffer.length < end) {
    return null;
  }
  return {
    message: JSON.parse(buffer.subarray(start, end).toString("utf8")),
    rest: buffer.subarray(end)
  };
}

function writeFrame(message) {
  const body = JSON.stringify(message);
  process.stdout.write("Content-Length: " + Buffer.byteLength(body, "utf8") + "\\r\\n\\r\\n" + body);
}

if (import.meta.url === "file://" + process.argv[1]) {
  startServer();
}
`;
}

function generatedReadme(manifest) {
  return `# ${manifest.name}

Generated by Gutenberg.

This package contains:

- an agent-native CLI: \`${manifest.slug}\`
- a stdio MCP server: \`npm run mcp\`
- a local response cache with offline search
- a machine-readable manifest: \`gutenberg.manifest.json\`
- smoke tests and cookbook documentation

## Product Thesis

${manifest.insights.thesis}

## Quick Start

\`\`\`bash
cp .env.example .env
node bin/${manifest.slug}.js operations
node bin/${manifest.slug}.js call ${manifest.operations[0]?.id || "operationId"} --json
node bin/${manifest.slug}.js sync --json
node bin/${manifest.slug}.js search keyword --json
npm run mcp
\`\`\`

Write and destructive operations are dry-run by default. Add \`--yes\` only when you want to execute them.

## Generated Advantages

${manifest.insights.generatedAdvantages.map((item) => `- ${item}`).join("\n")}
`;
}

function generatedCookbook(manifest) {
  const firstRead = manifest.operations.find((operation) => operation.cacheable);
  const firstWrite = manifest.operations.find((operation) => operation.risk !== "read");
  return `# ${manifest.name} Cookbook

## List Operations

\`\`\`bash
node bin/${manifest.slug}.js operations
\`\`\`

## Call an Operation

\`\`\`bash
node bin/${manifest.slug}.js call ${manifest.operations[0]?.id || "operationId"} --json
\`\`\`

## Cache Read Data

\`\`\`bash
node bin/${manifest.slug}.js sync ${firstRead?.id || ""} --json
node bin/${manifest.slug}.js search "status" --json
\`\`\`

## Preview a Write

\`\`\`bash
node bin/${manifest.slug}.js call ${firstWrite?.id || manifest.operations[0]?.id || "operationId"} --data '{"example":true}' --json
\`\`\`

Add \`--yes\` to execute a write/destructive operation.

## MCP

\`\`\`bash
npm run mcp
\`\`\`

The MCP server exposes:

- \`${manifest.slug}_operations\`
- \`${manifest.slug}_call\`
- \`${manifest.slug}_search_cache\`
`;
}

function generatedSmokeTest(manifest) {
  return `import test from "node:test";
import assert from "node:assert/strict";
import { listOperations, manifest } from "../src/client.js";
import { cacheStats } from "../src/store.js";

test("manifest exposes operations", () => {
  assert.equal(manifest.slug, "${manifest.slug}");
  assert.ok(listOperations().length > 0);
});

test("cache stats are available", () => {
  const stats = cacheStats();
  assert.equal(typeof stats.file, "string");
  assert.equal(typeof stats.records, "number");
});
`;
}
