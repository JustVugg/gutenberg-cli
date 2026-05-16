import fs from "node:fs";
import path from "node:path";
import os from "node:os";
import YAML from "yaml";
import { discoverOpenApi } from "./discover.js";
import { graphqlToOpenApi, loadGraphQLSource } from "./graphql.js";
import { harToOpenApi, loadHar } from "./har.js";
import { buildBlueprint, loadOpenApi } from "./openapi.js";
import { readJson, readText, writeJson } from "./fs.js";
import { safeOperationId, slugify } from "./sanitize.js";

const BODY_METHODS = new Set(["POST", "PUT", "PATCH"]);

export async function planSource(source, options = {}) {
  const resolved = await resolveSource(source, options);
  const blueprint = buildBlueprint(resolved.spec, resolved.sourcePath || resolved.input, options.name || resolved.name);
  const policy = resolvePolicy(options);
  return {
    schemaVersion: "gutenberg.plan.v1",
    source: resolved.input,
    kind: resolved.kind,
    normalizedSpec: resolved.sourcePath || null,
    name: blueprint.name,
    slug: blueprint.slug,
    baseUrls: blueprint.baseUrls,
    auth: blueprint.auth,
    operations: blueprint.operations.map((operation) => ({
      id: operation.id,
      method: operation.method,
      path: operation.path,
      risk: operation.risk,
      cacheable: operation.cacheable,
      tag: operation.tag,
      summary: operation.summary
    })),
    heroes: blueprint.heroes,
    policy,
    insights: blueprint.insights,
    nextCommands: nextCommandsForPlan(source, blueprint.slug)
  };
}

export async function materializeSource(source, options = {}) {
  const outDir = options.outDir ? path.resolve(options.outDir) : null;
  const resolved = await resolveSource(source, options);
  if (resolved.kind === "openapi" && resolved.sourcePath) {
    return resolved;
  }

  const slug = slugify(options.name || resolved.name || "generated-api");
  const specPath = path.resolve(
    options.specOut ||
      (outDir
        ? path.join(path.dirname(outDir), `${slug}.openapi.json`)
        : path.join(process.cwd(), "generated", `${slug}.openapi.json`))
  );
  writeJson(specPath, resolved.spec);
  return { ...resolved, sourcePath: specPath, normalized: true };
}

export async function resolveSource(source, options = {}) {
  if (!source) {
    const error = new Error("missing source. Example: gutenberg forge samples/petstore-openapi.json --name petstore");
    error.exitCode = 2;
    throw error;
  }
  const input = String(source);
  const kind = options.kind || detectSourceKind(input);

  if (kind === "openapi") {
    const sourcePath = path.resolve(input);
    return { input, kind, sourcePath, spec: loadOpenApi(sourcePath), name: options.name };
  }
  if (kind === "har") {
    const sourcePath = path.resolve(input);
    return {
      input,
      kind,
      sourcePath,
      spec: harToOpenApi(loadHar(sourcePath), { name: options.name, origin: options.origin }),
      name: options.name || path.basename(input, path.extname(input))
    };
  }
  if (kind === "graphql") {
    const gql = await loadGraphQLSource(input, {
      token: options.token || process.env.GUTENBERG_GRAPHQL_TOKEN
    });
    return {
      input,
      kind,
      sourcePath: fs.existsSync(input) ? path.resolve(input) : null,
      spec: graphqlToOpenApi(gql, {
        name: options.name,
        endpoint: /^https?:\/\//i.test(input) ? input : options.endpoint
      }),
      name: options.name || path.basename(input, path.extname(input))
    };
  }
  if (kind === "postman") {
    const collection = readJson(path.resolve(input));
    return {
      input,
      kind,
      sourcePath: path.resolve(input),
      spec: postmanToOpenApi(collection, { name: options.name }),
      name: options.name || collection.info?.name || path.basename(input, path.extname(input))
    };
  }
  if (kind === "insomnia") {
    const exportData = loadJsonOrYaml(path.resolve(input));
    return {
      input,
      kind,
      sourcePath: path.resolve(input),
      spec: insomniaToOpenApi(exportData, { name: options.name }),
      name: options.name || exportData.__export_source || path.basename(input, path.extname(input))
    };
  }
  if (kind === "curl") {
    const request = curlToRequest(input);
    return {
      input,
      kind,
      sourcePath: null,
      spec: requestsToOpenApi([request], { name: options.name || request.summary || "curl capture" }),
      name: options.name || request.summary || "curl capture"
    };
  }
  if (kind === "url") {
    const discovered = await discoverOpenApi(input, { out: null });
    if (!discovered.found || !discovered.spec) {
      const error = new Error(
        `No OpenAPI/Swagger document was discovered at ${input}. Use 'gutenberg record ${input} --out capture.har.json' or 'gutenberg forge <capture.har.json> --kind har'.`
      );
      error.exitCode = 2;
      error.attempts = discovered.attempts || [];
      throw error;
    }
    return {
      input,
      kind: "openapi",
      sourcePath: discovered.file || null,
      spec: discovered.spec,
      name: options.name || discovered.spec.info?.title
    };
  }

  const error = new Error(`unsupported source kind: ${kind}`);
  error.exitCode = 2;
  throw error;
}

export function detectSourceKind(input) {
  const value = String(input).trim();
  if (/^curl(\s|$)/i.test(value)) return "curl";
  if (/^https?:\/\//i.test(value)) return "url";

  if (!fs.existsSync(value)) {
    const error = new Error(`source does not exist and is not a URL/curl command: ${value}`);
    error.exitCode = 2;
    throw error;
  }

  const ext = path.extname(value).toLowerCase();
  if (ext === ".har") return "har";
  if (ext === ".graphql" || ext === ".gql") return "graphql";
  if (ext === ".yaml" || ext === ".yml") {
    const parsed = YAML.parse(readText(value));
    if (parsed?.openapi || parsed?.swagger) return "openapi";
    if (Array.isArray(parsed?.resources) || parsed?.__export_format) return "insomnia";
    return "openapi";
  }
  if (ext === ".json") {
    const parsed = readJson(value);
    if (parsed?.openapi || parsed?.swagger) return "openapi";
    if (parsed?.log && Array.isArray(parsed.log.entries)) return "har";
    if (isPostmanCollection(parsed)) return "postman";
    if (isInsomniaExport(parsed)) return "insomnia";
    if (parsed?.__schema || parsed?.data?.__schema) return "graphql";
  }
  return "openapi";
}

export function defaultPolicy() {
  return {
    schemaVersion: "gutenberg.policy.v1",
    rules: [
      { risk: "read", action: "allow", requiresYes: false },
      { risk: "write", action: "confirm", requiresYes: true },
      { risk: "destructive", action: "confirm", requiresYes: true }
    ],
    redaction: [
      "authorization",
      "cookie",
      "token",
      "secret",
      "api-key",
      "apikey",
      "client-secret",
      "session"
    ]
  };
}

export function resolvePolicy(options = {}) {
  if (options.policyPath) return loadPolicyFile(options.policyPath);
  if (typeof options.policy === "string") return loadPolicyFile(options.policy);
  if (options.policy && typeof options.policy === "object") return normalizePolicy(options.policy);
  return defaultPolicy();
}

export function loadPolicyFile(filePath) {
  const resolved = path.resolve(filePath);
  if (!fs.existsSync(resolved)) {
    const error = new Error(`Policy file not found: ${resolved}`);
    error.exitCode = 2;
    throw error;
  }
  return normalizePolicy(loadJsonOrYaml(resolved));
}

export function normalizePolicy(policy) {
  if (!policy || typeof policy !== "object") {
    const error = new Error("policy must be a JSON/YAML object");
    error.exitCode = 2;
    throw error;
  }
  const rules = Array.isArray(policy.rules) ? policy.rules : [];
  if (rules.length === 0) {
    const error = new Error("policy.rules must include at least one rule");
    error.exitCode = 2;
    throw error;
  }
  const allowedRisks = new Set(["read", "write", "destructive"]);
  const allowedActions = new Set(["allow", "confirm", "deny"]);
  const normalized = [];
  const seen = new Set();
  for (const rule of rules) {
    const risk = String(rule.risk || "").toLowerCase();
    const action = String(rule.action || "").toLowerCase();
    if (!allowedRisks.has(risk)) {
      const error = new Error(`policy rule has invalid risk: ${rule.risk}`);
      error.exitCode = 2;
      throw error;
    }
    if (!allowedActions.has(action)) {
      const error = new Error(`policy rule for ${risk} has invalid action: ${rule.action}`);
      error.exitCode = 2;
      throw error;
    }
    if (seen.has(risk)) {
      const error = new Error(`policy has duplicate rule for risk: ${risk}`);
      error.exitCode = 2;
      throw error;
    }
    seen.add(risk);
    normalized.push({
      risk,
      action,
      requiresYes: Boolean(rule.requiresYes || action === "confirm")
    });
  }
  for (const fallback of defaultPolicy().rules) {
    if (!seen.has(fallback.risk)) normalized.push(fallback);
  }
  return {
    schemaVersion: policy.schemaVersion || "gutenberg.policy.v1",
    rules: normalized,
    redaction: Array.isArray(policy.redaction) && policy.redaction.length > 0
      ? policy.redaction.map((item) => String(item).toLowerCase())
      : defaultPolicy().redaction,
    allowDomains: Array.isArray(policy.allowDomains) ? policy.allowDomains.map(String) : undefined,
    denyDomains: Array.isArray(policy.denyDomains) ? policy.denyDomains.map(String) : undefined
  };
}

export function postmanToOpenApi(collection, options = {}) {
  if (!isPostmanCollection(collection)) {
    const error = new Error("input does not look like a Postman collection");
    error.exitCode = 2;
    throw error;
  }
  const requests = [];
  collectPostmanItems(collection.item || [], requests);
  return requestsToOpenApi(requests, {
    name: options.name || collection.info?.name || "Postman Collection",
    description: "Generated from a Postman collection by Gutenberg."
  });
}

export function insomniaToOpenApi(exportData, options = {}) {
  if (!isInsomniaExport(exportData)) {
    const error = new Error("input does not look like an Insomnia export");
    error.exitCode = 2;
    throw error;
  }
  const requests = (exportData.resources || [])
    .filter((resource) => resource._type === "request" && resource.url)
    .map((resource) => ({
      method: String(resource.method || "GET").toUpperCase(),
      url: resource.url,
      headers: normalizeHeaderList(resource.headers),
      body: resource.body?.text || resource.body?.params || null,
      summary: resource.name || resource.url
    }));
  return requestsToOpenApi(requests, {
    name: options.name || "Insomnia Export",
    description: "Generated from an Insomnia export by Gutenberg."
  });
}

export function curlToRequest(command) {
  const tokens = splitCommandLine(command);
  if (tokens[0]?.toLowerCase() !== "curl") {
    const error = new Error("curl source must start with 'curl'");
    error.exitCode = 2;
    throw error;
  }
  let method = "";
  let url = "";
  let body = null;
  const headers = [];
  for (let i = 1; i < tokens.length; i += 1) {
    const token = tokens[i];
    const next = tokens[i + 1];
    if (token === "-X" || token === "--request") {
      method = String(next || "").toUpperCase();
      i += 1;
    } else if (token.startsWith("-X") && token.length > 2) {
      method = token.slice(2).toUpperCase();
    } else if (token === "-H" || token === "--header") {
      if (next) headers.push(parseHeader(next));
      i += 1;
    } else if (token === "-d" || token === "--data" || token === "--data-raw" || token === "--data-binary" || token === "--json") {
      body = next || "";
      if (!method) method = "POST";
      if (token === "--json") headers.push({ name: "Content-Type", value: "application/json" });
      i += 1;
    } else if (token === "--url") {
      url = next || "";
      i += 1;
    } else if (!token.startsWith("-") && /^https?:\/\//i.test(token)) {
      url = token;
    }
  }
  if (!url) {
    const error = new Error("curl command does not include an http(s) URL");
    error.exitCode = 2;
    throw error;
  }
  return { method: method || "GET", url, headers, body, summary: `${method || "GET"} ${new URL(url).pathname}` };
}

export function requestsToOpenApi(requests, options = {}) {
  const normalized = requests.map(normalizeRequest).filter(Boolean);
  if (normalized.length === 0) {
    const error = new Error("no HTTP requests found in source");
    error.exitCode = 2;
    throw error;
  }
  const origins = [...new Set(normalized.map((request) => request.origin))];
  const primaryOrigin = origins[0] || "https://api.example.com";
  const paths = {};
  const used = new Set();

  for (const request of normalized) {
    if (!paths[request.path]) paths[request.path] = {};
    const methodKey = request.method.toLowerCase();
    if (paths[request.path][methodKey]) continue;
    const operation = {
      operationId: uniqueOperationId(request.method, request.path, used),
      tags: [inferTag(request.path)],
      summary: request.summary || `${request.method} ${request.path}`,
      parameters: [
        ...pathParameters(request.path),
        ...request.query.map((item) => ({
          name: item.name,
          in: "query",
          required: false,
          schema: { type: "string" }
        }))
      ],
      responses: {
        200: { description: "Captured or imported response" }
      }
    };
    if (BODY_METHODS.has(request.method) || request.body) {
      operation.requestBody = {
        required: Boolean(request.body),
        content: {
          "application/json": { schema: { type: "object" } }
        }
      };
    }
    paths[request.path][methodKey] = operation;
  }

  const auth = sniffAuth(normalized);
  const spec = {
    openapi: "3.0.3",
    info: {
      title: options.name || "Imported API",
      version: "0.1.0",
      description: options.description || "Generated from imported HTTP requests by Gutenberg."
    },
    servers: [{ url: primaryOrigin }],
    paths
  };
  if (Object.keys(auth.schemes).length > 0) {
    spec.components = { securitySchemes: auth.schemes };
    spec.security = auth.security;
  }
  return spec;
}

function nextCommandsForPlan(source, slug) {
  const quoted = quoteForShell(source);
  return [
    `gutenberg forge ${quoted} --name ${slug} --install`,
    `gutenberg verify generated/${slug}-go`,
    `gutenberg scorecard generated/${slug}-go`
  ];
}

function isPostmanCollection(value) {
  return Boolean(value?.info && Array.isArray(value?.item) && /postman/i.test(String(value.info.schema || "")));
}

function isInsomniaExport(value) {
  return Boolean(Array.isArray(value?.resources) && (value.__export_format || value.__export_source || value.resources.some((r) => r._type === "request")));
}

function loadJsonOrYaml(filePath) {
  const ext = path.extname(filePath).toLowerCase();
  if (ext === ".yaml" || ext === ".yml") return YAML.parse(readText(filePath));
  return readJson(filePath);
}

function collectPostmanItems(items, out) {
  for (const item of items) {
    if (Array.isArray(item.item)) {
      collectPostmanItems(item.item, out);
      continue;
    }
    if (!item.request) continue;
    const request = item.request;
    const url = typeof request.url === "string" ? request.url : request.url?.raw;
    if (!url) continue;
    out.push({
      method: String(request.method || "GET").toUpperCase(),
      url,
      headers: normalizeHeaderList(request.header),
      body: request.body?.raw || null,
      summary: item.name || `${request.method || "GET"} ${url}`
    });
  }
}

function normalizeRequest(request) {
  let parsed;
  try {
    parsed = new URL(request.url);
  } catch {
    return null;
  }
  return {
    method: String(request.method || "GET").toUpperCase(),
    origin: parsed.origin,
    path: normalizeImportedPath(parsed.pathname),
    query: [...parsed.searchParams.entries()].map(([name, value]) => ({ name, value })),
    headers: normalizeHeaderList(request.headers),
    body: request.body || null,
    summary: request.summary
  };
}

function normalizeImportedPath(pathname) {
  const parts = pathname.split("/").filter(Boolean).map((part) => {
    const decoded = decodeURIComponent(part);
    if (decoded.startsWith(":") && decoded.length > 1) return `{${slugify(decoded.slice(1)).replaceAll("-", "_")}}`;
    if (/^\d+$/.test(decoded)) return "{id}";
    if (/^[0-9a-f]{24}$/i.test(decoded)) return "{id}";
    return part;
  });
  return `/${parts.join("/")}`;
}

function pathParameters(openapiPath) {
  const names = [...openapiPath.matchAll(/\{([^}]+)\}/g)].map((match) => match[1]);
  return names.map((name) => ({
    name,
    in: "path",
    required: true,
    schema: { type: "string" }
  }));
}

function normalizeHeaderList(headers) {
  if (!headers) return [];
  if (Array.isArray(headers)) {
    return headers
      .map((header) => {
        if (typeof header === "string") return parseHeader(header);
        return { name: header.name || header.key, value: header.value };
      })
      .filter((header) => header.name);
  }
  if (typeof headers === "object") {
    return Object.entries(headers).map(([name, value]) => ({ name, value: String(value) }));
  }
  return [];
}

function parseHeader(text) {
  const colon = String(text).indexOf(":");
  if (colon === -1) return { name: String(text).trim(), value: "" };
  return {
    name: String(text).slice(0, colon).trim(),
    value: String(text).slice(colon + 1).trim()
  };
}

function sniffAuth(requests) {
  const schemes = {};
  const apiKeyHeaders = new Set();
  let bearer = false;
  let basic = false;
  let cookie = false;
  for (const request of requests) {
    for (const header of request.headers || []) {
      const name = String(header.name || "").toLowerCase();
      const value = String(header.value || "");
      if (name === "authorization" && /^bearer\s+/i.test(value)) bearer = true;
      if (name === "authorization" && /^basic\s+/i.test(value)) basic = true;
      if (name === "cookie" && value) cookie = true;
      if (/^(x-api[-_]?key|x-auth[-_]?token|x-access[-_]?token|api[-_]?key)$/i.test(name)) apiKeyHeaders.add(header.name);
    }
  }
  if (bearer) schemes.bearerAuth = { type: "http", scheme: "bearer" };
  if (basic) schemes.basicAuth = { type: "http", scheme: "basic" };
  if (cookie) schemes.cookieAuth = { type: "apiKey", in: "cookie", name: "session" };
  for (const header of apiKeyHeaders) {
    const id = `${String(header).toLowerCase().replace(/[^a-z0-9]/g, "")}Auth`;
    schemes[id] = { type: "apiKey", in: "header", name: header };
  }
  return { schemes, security: Object.keys(schemes).map((name) => ({ [name]: [] })) };
}

function uniqueOperationId(method, openapiPath, used) {
  const base = safeOperationId(method, openapiPath);
  if (!used.has(base)) {
    used.add(base);
    return base;
  }
  let counter = 2;
  while (used.has(`${base}${counter}`)) counter += 1;
  const id = `${base}${counter}`;
  used.add(id);
  return id;
}

function inferTag(openapiPath) {
  return openapiPath.split("/").filter(Boolean)[0] || "default";
}

function splitCommandLine(input) {
  const tokens = [];
  let current = "";
  let quote = "";
  let escaped = false;
  for (const char of String(input)) {
    if (escaped) {
      current += char;
      escaped = false;
      continue;
    }
    if (char === "\\") {
      escaped = true;
      continue;
    }
    if (quote) {
      if (char === quote) quote = "";
      else current += char;
      continue;
    }
    if (char === "'" || char === "\"") {
      quote = char;
      continue;
    }
    if (/\s/.test(char)) {
      if (current) {
        tokens.push(current);
        current = "";
      }
      continue;
    }
    current += char;
  }
  if (current) tokens.push(current);
  return tokens;
}

function quoteForShell(value) {
  const text = String(value);
  if (/^curl\s/i.test(text)) return JSON.stringify(text);
  if (/^[A-Za-z0-9_./:@-]+$/.test(text)) return text;
  return `'${text.replaceAll("'", "'\\''")}'`;
}

export function tempSpecPath(name = "source") {
  return path.join(os.tmpdir(), `gutenberg-${slugify(name)}-${Date.now()}.openapi.json`);
}
