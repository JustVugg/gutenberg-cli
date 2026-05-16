import { readJson } from "./fs.js";
import { safeOperationId } from "./sanitize.js";

const BODY_METHODS = new Set(["POST", "PUT", "PATCH"]);

export function loadHar(filePath) {
  const har = readJson(filePath);
  if (!har.log || !Array.isArray(har.log.entries)) {
    const error = new Error("input does not look like a HAR file");
    error.exitCode = 2;
    throw error;
  }
  return har;
}

export function harToOpenApi(har, options = {}) {
  const name = options.name || "Captured API";
  const entries = har.log.entries
    .map((entry) => normalizeEntry(entry))
    .filter(Boolean)
    .filter((entry) => !isAsset(entry));

  const origins = [...new Set(entries.map((entry) => entry.origin))];
  const primaryOrigin = options.origin || origins[0] || "https://api.example.com";
  const paths = {};
  const usedOperationIds = new Set();

  for (const entry of entries) {
    if (options.origin && entry.origin !== options.origin) {
      continue;
    }
    const openapiPath = normalizePath(entry.pathname);
    const method = entry.method.toLowerCase();
    if (!paths[openapiPath]) {
      paths[openapiPath] = {};
    }
    if (paths[openapiPath][method]) {
      mergeParameters(paths[openapiPath][method], entry);
      continue;
    }
    const responseSchema = entry.responseJSON ? inferSchema(entry.responseJSON) : null;
    const responseContent = entry.responseJSON
      ? {
          [entry.mimeType.split(";")[0] || "application/json"]: { schema: responseSchema }
        }
      : undefined;
    paths[openapiPath][method] = {
      operationId: uniqueOperationId(method, openapiPath, usedOperationIds),
      tags: [inferTag(openapiPath)],
      summary: `${entry.method} ${openapiPath}`,
      parameters: [
        ...pathParameters(openapiPath),
        ...queryParameters(entry)
      ],
      responses: {
        [String(entry.status || 200)]: responseContent
          ? { description: entry.statusText || "Captured response", content: responseContent }
          : { description: entry.statusText || "Captured response" }
      }
    };
    if (BODY_METHODS.has(entry.method) || entry.hasBody) {
      paths[openapiPath][method].requestBody = {
        required: Boolean(entry.hasBody),
        content: {
          "application/json": {
            schema: { type: "object" }
          }
        }
      };
    }
  }

  const authDetected = sniffAuthFromEntries(entries);
  const spec = {
    openapi: "3.0.3",
    info: {
      title: name,
      version: "0.1.0",
      description: "Generated from a browser HAR capture by Gutenberg."
    },
    servers: [{ url: primaryOrigin }],
    paths
  };
  if (authDetected.schemes && Object.keys(authDetected.schemes).length > 0) {
    spec.components = { securitySchemes: authDetected.schemes };
    spec.security = authDetected.security;
  }
  return spec;
}

function sniffAuthFromEntries(entries) {
  const schemes = {};
  let bearerSeen = false;
  let cookieSeen = false;
  const apiKeyHeaders = new Map();
  for (const entry of entries) {
    for (const header of entry.headers || []) {
      const name = String(header.name || "").toLowerCase();
      const value = String(header.value || "");
      if (name === "authorization") {
        if (/^bearer\s+/i.test(value)) bearerSeen = true;
        else if (/^basic\s+/i.test(value)) schemes.basicAuth = { type: "http", scheme: "basic" };
      } else if (name === "cookie" && value) {
        cookieSeen = true;
      } else if (/^(x-api[-_]?key|x-auth[-_]?token|x-access[-_]?token|api[-_]?key)$/i.test(name)) {
        apiKeyHeaders.set(header.name, (apiKeyHeaders.get(header.name) || 0) + 1);
      }
    }
  }
  if (bearerSeen) schemes.bearerAuth = { type: "http", scheme: "bearer" };
  if (cookieSeen) schemes.cookieAuth = { type: "apiKey", in: "cookie", name: "session" };
  for (const [headerName, count] of apiKeyHeaders) {
    if (count >= 1) {
      const id = headerName.toLowerCase().replace(/[^a-z0-9]/g, "") + "Auth";
      schemes[id] = { type: "apiKey", in: "header", name: headerName };
    }
  }
  const security = Object.keys(schemes).map((name) => ({ [name]: [] }));
  return { schemes, security };
}

function uniqueOperationId(method, openapiPath, used) {
  const base = safeOperationId(method, openapiPath);
  if (!used.has(base)) {
    used.add(base);
    return base;
  }
  // Collision: prepend a disambiguator from earlier path segments.
  const segments = openapiPath.replace(/[{}]/g, "").split("/").filter(Boolean);
  for (let prefix = 1; prefix <= segments.length; prefix++) {
    const candidate = safeOperationId(method, "/" + segments.slice(-3 - prefix).join("/"));
    if (!used.has(candidate)) {
      used.add(candidate);
      return candidate;
    }
  }
  let counter = 2;
  while (used.has(`${base}${counter}`)) counter += 1;
  const fallback = `${base}${counter}`;
  used.add(fallback);
  return fallback;
}

function normalizeEntry(entry) {
  const request = entry.request || {};
  if (!request.url || !request.method) {
    return null;
  }
  let parsed;
  try {
    parsed = new URL(request.url);
  } catch {
    return null;
  }
  const responseText = entry.response?.content?.text || "";
  let responseJSON = null;
  if (responseText && /json/.test(entry.response?.content?.mimeType || "")) {
    try {
      responseJSON = JSON.parse(responseText);
    } catch {}
  }
  return {
    method: request.method.toUpperCase(),
    origin: parsed.origin,
    pathname: parsed.pathname,
    queryString: Array.isArray(request.queryString) ? request.queryString : [],
    headers: Array.isArray(request.headers) ? request.headers : [],
    responseHeaders: Array.isArray(entry.response?.headers) ? entry.response.headers : [],
    mimeType: entry.response?.content?.mimeType || "",
    status: entry.response?.status || 200,
    statusText: entry.response?.statusText || "",
    hasBody: Boolean(request.postData?.text),
    responseJSON
  };
}

export function inferSchema(value, depth = 0) {
  if (depth > 6) return { type: "object" };
  if (value === null) return { type: "null" };
  if (typeof value === "boolean") return { type: "boolean" };
  if (typeof value === "number") return { type: Number.isInteger(value) ? "integer" : "number" };
  if (typeof value === "string") {
    if (/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/.test(value)) return { type: "string", format: "date-time" };
    if (/^\d{4}-\d{2}-\d{2}$/.test(value)) return { type: "string", format: "date" };
    if (/^https?:\/\//.test(value)) return { type: "string", format: "uri" };
    if (/^[\w.+-]+@[\w-]+\.[\w.-]+$/.test(value)) return { type: "string", format: "email" };
    return { type: "string" };
  }
  if (Array.isArray(value)) {
    if (value.length === 0) return { type: "array", items: { type: "object" } };
    const merged = value.slice(0, 5).reduce((acc, item) => mergeSchema(acc, inferSchema(item, depth + 1)), null);
    return { type: "array", items: merged || { type: "object" } };
  }
  if (typeof value === "object") {
    const properties = {};
    const required = [];
    for (const [key, item] of Object.entries(value)) {
      properties[key] = inferSchema(item, depth + 1);
      if (item !== null && item !== undefined) required.push(key);
    }
    const schema = { type: "object" };
    if (Object.keys(properties).length > 0) schema.properties = properties;
    if (required.length > 0 && required.length <= 8) schema.required = required;
    return schema;
  }
  return { type: "object" };
}

function mergeSchema(a, b) {
  if (!a) return b;
  if (!b) return a;
  if (a.type !== b.type) return { type: a.type === "null" ? b.type : a.type, nullable: true };
  if (a.type === "object" && b.type === "object") {
    const properties = { ...(a.properties || {}) };
    for (const [key, schema] of Object.entries(b.properties || {})) {
      properties[key] = properties[key] ? mergeSchema(properties[key], schema) : schema;
    }
    const required = [...new Set([...(a.required || []), ...(b.required || [])])];
    const merged = { type: "object" };
    if (Object.keys(properties).length) merged.properties = properties;
    if (required.length) merged.required = required;
    return merged;
  }
  if (a.type === "array" && b.type === "array") {
    return { type: "array", items: mergeSchema(a.items, b.items) };
  }
  return a;
}

function isAsset(entry) {
  const path = entry.pathname.toLowerCase();
  const mime = entry.mimeType.toLowerCase();
  if (path.match(/\.(png|jpg|jpeg|gif|webp|svg|ico|css|js|map|woff|woff2|ttf|mp4|webm|pdf)$/)) {
    return true;
  }
  if (mime.startsWith("image/") || mime.includes("font") || mime.includes("css") || mime.includes("javascript")) {
    return true;
  }
  return false;
}

function normalizePath(pathname) {
  const rawParts = pathname.split("/").filter(Boolean);
  const decoded = rawParts.map((part) => decodeURIComponent(part));
  const parts = decoded.map((segment, index) => {
    const dateRole = detectDateRole(decoded, index);
    if (dateRole) return `{${dateRole}}`;
    if (/^\d+$/.test(segment)) return "{id}";
    if (/^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(segment)) return "{id}";
    if (/^[0-9a-f]{24}$/i.test(segment)) return "{id}";
    return rawParts[index] || "segment";
  });
  return `/${parts.join("/")}`;
}

function detectDateRole(decoded, index) {
  const segment = decoded[index];
  if (!segment) return null;
  if (/^(19|20)\d{2}$/.test(segment) && looksLikeDateContext(decoded, index)) return "year";
  if (/^(0?[1-9]|1[0-2])$/.test(segment) && /^(19|20)\d{2}$/.test(decoded[index - 1] || "")) return "month";
  if (/^(0?[1-9]|[12]\d|3[01])$/.test(segment) && /^(0?[1-9]|1[0-2])$/.test(decoded[index - 1] || "") && /^(19|20)\d{2}$/.test(decoded[index - 2] || "")) return "day";
  return null;
}

function looksLikeDateContext(decoded, index) {
  if (index + 1 < decoded.length && /^(0?[1-9]|1[0-2])$/.test(decoded[index + 1])) return true;
  return false;
}

function inferTag(openapiPath) {
  const first = openapiPath.split("/").filter(Boolean)[0];
  return first || "default";
}

function pathParameters(openapiPath) {
  const matches = [...openapiPath.matchAll(/\{([^}]+)\}/g)];
  return matches.map((match) => ({
    name: match[1],
    in: "path",
    required: true,
    schema: { type: "string" }
  }));
}

function queryParameters(entry) {
  return entry.queryString
    .filter((item) => item.name)
    .map((item) => ({
      name: item.name,
      in: "query",
      required: false,
      schema: { type: inferPrimitive(item.value) }
    }));
}

function mergeParameters(operation, entry) {
  const existing = new Set((operation.parameters || []).map((parameter) => `${parameter.in}:${parameter.name}`));
  for (const parameter of queryParameters(entry)) {
    const key = `${parameter.in}:${parameter.name}`;
    if (!existing.has(key)) {
      operation.parameters.push(parameter);
      existing.add(key);
    }
  }
}

function inferPrimitive(value) {
  if (value === "true" || value === "false") return "boolean";
  if (/^-?\d+(\.\d+)?$/.test(String(value))) return "number";
  return "string";
}
