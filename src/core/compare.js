import fs from "node:fs";
import path from "node:path";

export async function compareTools(toolDirs, options = {}) {
  if (!Array.isArray(toolDirs) || toolDirs.length < 2) {
    const error = new Error("compare requires at least two tool directories");
    error.exitCode = 2;
    throw error;
  }
  const operation = options.operation;
  if (!operation) {
    const error = new Error("compare requires --op <operationId>");
    error.exitCode = 2;
    throw error;
  }

  const params = options.params || {};
  const results = await Promise.all(toolDirs.map((dir) => callOneTool(path.resolve(dir), operation, params, options)));

  return {
    schemaVersion: "gutenberg.compare.v1",
    operation,
    params,
    results
  };
}

async function callOneTool(toolDir, operationId, params, options) {
  const manifestPath = path.join(toolDir, "gutenberg.manifest.json");
  if (!fs.existsSync(manifestPath)) {
    return { tool: path.basename(toolDir), error: `no manifest in ${toolDir}` };
  }
  const manifest = JSON.parse(fs.readFileSync(manifestPath, "utf8"));
  const op = (manifest.operations || []).find((candidate) => candidate.id === operationId || candidate.id?.endsWith(operationId));
  if (!op) {
    return { tool: manifest.slug || path.basename(toolDir), error: `operation not found: ${operationId}` };
  }

  const baseUrl = options.baseUrl || (manifest.baseUrls && manifest.baseUrls[0]);
  if (!baseUrl) {
    return { tool: manifest.slug, error: "missing base URL" };
  }

  let apiPath = op.path;
  for (const param of op.parameters || []) {
    if (param.in === "path") {
      const value = params[param.name];
      if (value === undefined && param.required) {
        return { tool: manifest.slug, error: `missing path parameter: ${param.name}` };
      }
      apiPath = apiPath.replace(`{${param.name}}`, encodeURIComponent(String(value ?? "")));
    }
  }

  const url = new URL(baseUrl.replace(/\/$/, "") + apiPath);
  for (const param of op.parameters || []) {
    if (param.in === "query") {
      const value = params[param.name];
      if (value !== undefined && value !== null) url.searchParams.set(param.name, String(value));
    }
  }

  const headers = { Accept: "application/json" };
  if (options.token) headers.Authorization = `Bearer ${options.token}`;

  const startedAt = Date.now();
  try {
    const response = await fetch(url.toString(), { method: op.method, headers });
    const text = await response.text();
    let body;
    try { body = JSON.parse(text); } catch { body = text; }
    return {
      tool: manifest.slug,
      url: url.toString(),
      status: response.status,
      elapsedMs: Date.now() - startedAt,
      body: options.compact ? compactValue(body) : body
    };
  } catch (error) {
    return { tool: manifest.slug, url: url.toString(), error: error.message };
  }
}

function compactValue(value) {
  if (Array.isArray(value)) return value.slice(0, 3);
  if (value && typeof value === "object") {
    const keys = Object.keys(value).slice(0, 6);
    const out = {};
    for (const key of keys) out[key] = value[key];
    if (Object.keys(value).length > 6) out.__truncated = true;
    return out;
  }
  return value;
}
