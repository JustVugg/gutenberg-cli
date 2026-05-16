import path from "node:path";
import YAML from "yaml";
import { readJson, readText } from "./fs.js";
import { envPrefix, safeOperationId, slugify } from "./sanitize.js";

const HTTP_METHODS = new Set(["get", "put", "post", "delete", "patch", "head", "options", "trace"]);

export function loadOpenApi(specPath) {
  const extension = path.extname(specPath).toLowerCase();
  if (extension === ".yaml" || extension === ".yml") {
    const spec = YAML.parse(readText(specPath));
    assertOpenApi(spec);
    return spec;
  }

  const spec = readJson(specPath);
  assertOpenApi(spec);
  return spec;
}

function assertOpenApi(spec) {
  if (!spec.openapi && !spec.swagger) {
    const error = new Error("input does not look like an OpenAPI/Swagger document");
    error.exitCode = 2;
    throw error;
  }
}

export function buildBlueprint(spec, sourcePath, requestedName) {
  const title = requestedName || spec.info?.title || "Generated API";
  const slug = slugify(title);
  const operations = [];
  const securitySchemes = spec.components?.securitySchemes || spec.securityDefinitions || {};
  const globalParameters = [];

  for (const [apiPath, pathItem] of Object.entries(spec.paths || {})) {
    const inheritedParameters = Array.isArray(pathItem.parameters) ? pathItem.parameters : globalParameters;
    for (const [method, operation] of Object.entries(pathItem)) {
      if (!HTTP_METHODS.has(method.toLowerCase())) {
        continue;
      }
      const parameters = [
        ...inheritedParameters,
        ...(Array.isArray(operation.parameters) ? operation.parameters : [])
      ].map((parameter) => normalizeParameter(resolveParameter(spec, parameter)));
      const id = operation.operationId || safeOperationId(method, apiPath);
      const op = {
        id,
        method: method.toUpperCase(),
        path: apiPath,
        tag: (operation.tags && operation.tags[0]) || "default",
        summary: operation.summary || operation.description || id,
        description: operation.description || "",
        parameters,
        hasRequestBody: Boolean(operation.requestBody),
        risk: classifyRisk(method, operation),
        cacheable: method.toLowerCase() === "get" || operation["x-graphql"]?.kind === "query",
        responseCodes: Object.keys(operation.responses || {}),
        inputHints: inferInputHints(parameters, operation),
        graphql: operation["x-graphql"] || null
      };
      op.kind = classifyOperation(op);
      op.pagination = detectPagination(op);
      if (operation["x-gutenberg-hero"]) {
        op.heroOverride = operation["x-gutenberg-hero"];
      }
      operations.push(op);
    }
  }

  const servers = Array.isArray(spec.servers) && spec.servers.length > 0
    ? spec.servers.map((server) => server.url).filter(Boolean)
    : spec.host
      ? [`${spec.schemes?.[0] || "https"}://${spec.host}${spec.basePath || ""}`]
      : [];

  linkRelatedOperations(operations);
  const tags = [...new Set(operations.map((operation) => operation.tag))].sort();
  const mutating = operations.filter((operation) => operation.risk !== "read").length;
  const cacheable = operations.filter((operation) => operation.cacheable).length;
  const heroes = pickHeroes(operations);

  return {
    schemaVersion: "gutenberg.blueprint.v1",
    source: sourcePath,
    name: title,
    slug,
    envPrefix: envPrefix(slug),
    description: spec.info?.description || "",
    version: spec.info?.version || "0.0.0",
    baseUrls: servers,
    auth: summarizeAuth(securitySchemes, spec.security),
    tags,
    operations,
    heroes,
    insights: buildInsights({ title, slug, operations, tags, cacheable, mutating }),
    generatedAt: new Date().toISOString()
  };
}

function pickHeroes(operations) {
  // Explicit heroes via x-gutenberg-hero spec extension take priority.
  const explicit = [];
  const usedAliases = new Set();
  for (const op of operations) {
    if (!op.heroOverride) continue;
    const alias = slugify(op.heroOverride.alias || op.id);
    if (!alias || usedAliases.has(alias)) continue;
    usedAliases.add(alias);
    explicit.push({
      alias,
      operationId: op.id,
      summary: op.heroOverride.summary || op.summary,
      method: op.method,
      path: op.path,
      defaultParams: op.heroOverride["default-params"] || op.heroOverride.defaultParams || {},
      explicit: true
    });
  }

  const autoCandidates = [];
  for (const op of operations) {
    if (op.heroOverride) continue;
    if (op.method !== "GET") continue;
    const hasRequiredPath = (op.parameters || []).some((p) => p.in === "path" && p.required);
    if (hasRequiredPath) continue;
    const requiredQueries = (op.parameters || []).filter((p) => p.in === "query" && p.required).length;
    if (requiredQueries > 1) continue;
    autoCandidates.push(op);
  }

  const limited = autoCandidates.slice(0, 20);
  const auto = assignHeroAliases(limited).map((entry) => ({
    alias: entry.alias,
    operationId: entry.op.id,
    summary: entry.op.summary,
    method: entry.op.method,
    path: entry.op.path,
    defaultParams: {},
    explicit: false
  })).filter((hero) => !usedAliases.has(hero.alias));

  return [...explicit, ...auto];
}

function meaningfulSegments(p) {
  const NOISE = new Set(["api", "apis", "rest", "public", "json", "service", "services", "site", "www", "graph", "graphql"]);
  return p
    .replace(/[{}]/g, "")
    .split("/")
    .filter(Boolean)
    .filter((segment) => !NOISE.has(segment.toLowerCase()) && !/^v\d+(\.\d+)?$/i.test(segment))
    .map((segment) => slugify(segment));
}

function assignHeroAliases(candidates) {
  const entries = candidates.map((op) => ({ op, segments: meaningfulSegments(op.path || op.id) }));

  // Strip a common trailing segment if shared by all (e.g. "scoreboard" across 4 endpoints).
  // This pulls the discriminator forward so aliases become 'nba', 'nfl', 'mlb' instead of all
  // collapsing into 'scoreboard'.
  if (entries.length > 1) {
    while (true) {
      const last = entries[0].segments[entries[0].segments.length - 1];
      if (!last) break;
      const sharedSuffix = entries.every((entry) => entry.segments[entry.segments.length - 1] === last);
      if (!sharedSuffix) break;
      // Don't strip if doing so would leave any entry empty.
      if (entries.some((entry) => entry.segments.length <= 1)) break;
      for (const entry of entries) entry.segments = entry.segments.slice(0, -1);
    }
  }

  const used = new Set();
  for (const entry of entries) {
    const segments = entry.segments.length > 0 ? entry.segments : [slugify(entry.op.tag || "op")];
    let alias = "";
    for (let take = 1; take <= segments.length; take++) {
      const candidate = segments.slice(-take).join("-");
      if (!used.has(candidate) && candidate) {
        alias = candidate;
        break;
      }
    }
    if (!alias) {
      let counter = 2;
      const base = segments.join("-") || "op";
      alias = base;
      while (used.has(alias)) {
        alias = `${base}-${counter++}`;
      }
    }
    used.add(alias);
    entry.alias = alias;
  }
  return entries;
}

function resolveParameter(spec, parameter) {
  if (!parameter || !parameter.$ref) {
    return parameter || {};
  }
  const resolved = resolveRef(spec, parameter.$ref);
  return resolved || parameter;
}

function resolveRef(spec, ref) {
  if (!ref.startsWith("#/")) {
    return null;
  }
  const parts = ref.slice(2).split("/").map((part) => part.replace(/~1/g, "/").replace(/~0/g, "~"));
  let current = spec;
  for (const part of parts) {
    current = current?.[part];
    if (current === undefined) {
      return null;
    }
  }
  return current;
}

function normalizeParameter(parameter) {
  return {
    name: parameter.name || "",
    in: parameter.in || "query",
    required: Boolean(parameter.required),
    description: parameter.description || "",
    type: parameter.schema?.type || parameter.type || "string"
  };
}

function inferInputHints(parameters, operation) {
  const path = parameters.filter((parameter) => parameter.in === "path").map((parameter) => parameter.name);
  const query = parameters.filter((parameter) => parameter.in === "query").map((parameter) => parameter.name);
  const headers = parameters.filter((parameter) => parameter.in === "header").map((parameter) => parameter.name);
  return {
    path,
    query,
    headers,
    body: Boolean(operation.requestBody)
  };
}

function classifyRisk(method, operation) {
  if (operation["x-graphql"]?.kind === "query") {
    return "read";
  }
  if (operation["x-graphql"]?.kind === "mutation") {
    return "write";
  }
  const normalized = method.toLowerCase();
  if (normalized === "get" || normalized === "head" || normalized === "options") {
    return "read";
  }
  if (normalized === "delete") {
    return "destructive";
  }
  if (operation.summary && /delete|remove|cancel|void|refund/i.test(operation.summary)) {
    return "destructive";
  }
  return "write";
}

function summarizeAuth(schemes, globalSecurity) {
  const names = Object.keys(schemes);
  if (names.length === 0 && (!globalSecurity || globalSecurity.length === 0)) {
    return { mode: "none", schemes: [], env: null };
  }

  const schemeSummaries = names.map((name) => {
    const scheme = schemes[name] || {};
    return {
      name,
      type: scheme.type || "unknown",
      in: scheme.in || null,
      header: scheme.name || null,
      scheme: scheme.scheme || null,
      flows: normalizeOAuthFlows(scheme.flows)
    };
  });

  return {
    mode: "detected",
    schemes: schemeSummaries,
    env: "API_KEY",
    oauth: schemeSummaries.some((scheme) => scheme.type === "oauth2")
  };
}

function normalizeOAuthFlows(flows) {
  if (!flows || typeof flows !== "object") {
    return {};
  }
  const output = {};
  for (const [name, flow] of Object.entries(flows)) {
    output[name] = {
      authorizationUrl: flow.authorizationUrl || "",
      tokenUrl: flow.tokenUrl || "",
      refreshUrl: flow.refreshUrl || "",
      scopes: flow.scopes || {}
    };
  }
  return output;
}

function buildInsights({ title, slug, operations, tags, cacheable, mutating }) {
  const collectionReads = operations.filter((operation) => {
    return operation.cacheable && !operation.path.match(/\{[^}]+\}$/);
  });
  const entityReads = operations.filter((operation) => {
    return operation.cacheable && operation.path.match(/\{[^}]+\}/);
  });
  const destructive = operations.filter((operation) => operation.risk === "destructive");

  return {
    thesis: `${title} should become a local, searchable operational workspace, not only an endpoint wrapper.`,
    generatedAdvantages: [
      "One command surface for humans and agents",
      "Shared core between CLI and MCP server",
      "Local cache with search and replay metadata",
      "Risk-aware dry-run defaults for write/destructive calls",
      "Machine-readable manifest for catalogs and scorecards"
    ],
    recommendedCommands: [
      collectionReads.length > 0 ? "sync: cache list endpoints locally" : null,
      entityReads.length > 0 ? "inspect: fetch entity details by id" : null,
      cacheable > 0 ? "search: query cached API responses offline" : null,
      mutating > 0 ? "plan: preview write operations before execution" : null
    ].filter(Boolean),
    domainMap: tags.map((tag) => ({
      tag,
      operations: operations.filter((operation) => operation.tag === tag).length
    })),
    riskNotes: destructive.length > 0
      ? [`${destructive.length} destructive operation(s) require explicit --yes in generated clients.`]
      : []
  };
}

export function classifyOperation(op) {
  const method = op.method.toUpperCase();
  const path = op.path;
  const lastSegment = path.split("/").filter(Boolean).pop() || "";
  const hasPathParam = /\{[^}]+\}/.test(path);
  const endsWithParam = /\}$/.test(path);
  const lowerSummary = (op.summary || "").toLowerCase();

  if (method === "DELETE") return "delete";
  if (method === "POST") {
    if (/search|query|find|lookup|filter/.test(lastSegment) || /search|query/.test(lowerSummary)) return "search";
    if (/login|auth|token|signup|register/.test(lastSegment)) return "action";
    return "create";
  }
  if (method === "PUT" || method === "PATCH") {
    return endsWithParam || hasPathParam ? "update" : "upsert";
  }
  if (method === "GET") {
    if (/search|query|find|lookup/.test(lastSegment) || /search|query|find/.test(lowerSummary)) return "search";
    if (endsWithParam) return "read-detail";
    if (!hasPathParam) return "read-list";
    return "read";
  }
  return "action";
}

export function detectPagination(op) {
  const queryNames = new Set((op.parameters || []).filter((p) => p.in === "query").map((p) => (p.name || "").toLowerCase()));
  const hasOffsetLimit = (queryNames.has("offset") || queryNames.has("skip")) && (queryNames.has("limit") || queryNames.has("count") || queryNames.has("size"));
  if (hasOffsetLimit) return { style: "offset-limit", offsetParam: queryNames.has("offset") ? "offset" : "skip", limitParam: queryNames.has("limit") ? "limit" : queryNames.has("count") ? "count" : "size" };
  if (queryNames.has("cursor") || queryNames.has("after") || queryNames.has("next") || queryNames.has("page_token")) return { style: "cursor", cursorParam: queryNames.has("cursor") ? "cursor" : queryNames.has("after") ? "after" : queryNames.has("page_token") ? "page_token" : "next" };
  if (queryNames.has("page") || queryNames.has("page_number")) return { style: "page", pageParam: queryNames.has("page") ? "page" : "page_number", perPageParam: queryNames.has("per_page") ? "per_page" : queryNames.has("page_size") ? "page_size" : queryNames.has("limit") ? "limit" : null };
  return null;
}

export function linkRelatedOperations(operations) {
  const byPath = new Map();
  for (const op of operations) {
    const key = op.path.replace(/\{[^}]+\}/g, "{*}");
    if (!byPath.has(key)) byPath.set(key, []);
    byPath.get(key).push(op);
  }

  for (const op of operations) {
    const segments = op.path.split("/").filter(Boolean);
    const lastIsParam = /^\{[^}]+\}$/.test(segments[segments.length - 1] || "");
    const parentPath = lastIsParam ? "/" + segments.slice(0, -1).join("/") : null;
    const childPath = lastIsParam ? null : op.path + "/{id}";

    const related = [];
    for (const other of operations) {
      if (other === op) continue;
      if (parentPath && other.path === parentPath) related.push({ id: other.id, role: "list-parent" });
      if (childPath && other.path.startsWith(op.path + "/{")) related.push({ id: other.id, role: "detail-child" });
      if (other.path === op.path && other.method !== op.method) related.push({ id: other.id, role: "same-resource" });
    }
    op.related = related;
  }
  return operations;
}
