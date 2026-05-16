import { writeJson } from "./fs.js";

const COMMON_SPEC_PATHS = [
  "/openapi.json",
  "/swagger.json",
  "/api-docs",
  "/v3/api-docs",
  "/docs/openapi.json",
  "/.well-known/openapi.json"
];

export async function discoverOpenApi(url, options = {}) {
  const input = new URL(url).toString();
  const root = normalizeRoot(url);
  const attempts = [input];
  for (const candidate of COMMON_SPEC_PATHS) {
    attempts.push(new URL(candidate, root).toString());
  }

  const html = await fetchText(root).catch(() => null);
  if (html) {
    for (const candidate of extractSpecLinks(html, root)) {
      if (!attempts.includes(candidate)) {
        attempts.push(candidate);
      }
    }
  }

  for (const candidate of attempts) {
    const result = await fetchJson(candidate).catch((error) => ({ error: error.message }));
    if (result && result.openapi || result && result.swagger) {
      normalizeServers(result, candidate);
      if (options.out) {
        writeJson(options.out, result);
      }
      return {
        found: true,
        url: candidate,
        attempts,
        spec: result
      };
    }
  }

  return {
    found: false,
    url: null,
    attempts,
    spec: null
  };
}

function normalizeServers(spec, specURL) {
  if (!Array.isArray(spec.servers)) {
    return;
  }
  spec.servers = spec.servers.map((server) => {
    if (!server || !server.url) {
      return server;
    }
    try {
      return {
        ...server,
        url: new URL(server.url, specURL).toString().replace(/\/$/, "")
      };
    } catch {
      return server;
    }
  });
}


function normalizeRoot(value) {
  const parsed = new URL(value);
  return `${parsed.protocol}//${parsed.host}/`;
}

async function fetchText(url) {
  const response = await fetch(url, { headers: { Accept: "text/html,application/json" } });
  if (!response.ok) {
    throw new Error(`${response.status} ${response.statusText}`);
  }
  return response.text();
}

async function fetchJson(url) {
  const response = await fetch(url, { headers: { Accept: "application/json" } });
  if (!response.ok) {
    throw new Error(`${response.status} ${response.statusText}`);
  }
  const text = await response.text();
  return JSON.parse(text);
}

function extractSpecLinks(html, root) {
  const links = [];
  const patterns = [
    /href=["']([^"']*(?:openapi|swagger)[^"']*\.json[^"']*)["']/gi,
    /url:\s*["']([^"']*(?:openapi|swagger|api-docs)[^"']*)["']/gi,
    /["']([^"']*(?:openapi|swagger|api-docs)[^"']*\.json)["']/gi
  ];
  for (const pattern of patterns) {
    for (const match of html.matchAll(pattern)) {
      try {
        links.push(new URL(match[1], root).toString());
      } catch {
        // Ignore malformed links.
      }
    }
  }
  return [...new Set(links)];
}
