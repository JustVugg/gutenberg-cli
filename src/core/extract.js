import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import crypto from "node:crypto";

const anthropicModel = () => process.env.GUTENBERG_LLM_MODEL || "claude-sonnet-4-5";
const openaiModel = () => process.env.GUTENBERG_LLM_MODEL_OPENAI || "gpt-4o-mini";
const ollamaModel = () => process.env.GUTENBERG_LLM_MODEL_OLLAMA || "llama3.1";
const ollamaHost = () => process.env.OLLAMA_HOST || "http://localhost:11434";

export async function extractFromUrl(url, options = {}) {
  const cacheTtlMs = parseTtl(options.cache);
  if (cacheTtlMs > 0) {
    const cached = readExtractCache(url, options, cacheTtlMs);
    if (cached) {
      if (options.verbose !== false) {
        process.stderr.write(`extract: cache hit (age ${Math.round(cached.ageSeconds)}s)\n`);
      }
      return { parsed: cached.parsed, errors: cached.errors, raw: cached.raw, cached: true };
    }
  }
  const html = await fetchHtml(url);
  let preparedText;
  try {
    const { extractMainContent, htmlToMarkdown } = await import("./scrape.js");
    preparedText = htmlToMarkdown(extractMainContent(html), { baseUrl: url });
  } catch {
    preparedText = stripHtml(html);
  }
  const result = await extractFromText(preparedText, options);
  if (cacheTtlMs > 0) {
    writeExtractCache(url, options, result);
  }
  return result;
}

function extractCacheDir() {
  if (process.env.GUTENBERG_EXTRACT_CACHE_DIR) return process.env.GUTENBERG_EXTRACT_CACHE_DIR;
  return path.join(os.homedir(), ".gutenberg", "extract-cache");
}

function extractCacheKey(url, options) {
  const provider = options.provider || (process.env.ANTHROPIC_API_KEY ? "anthropic" : process.env.OPENAI_API_KEY ? "openai" : "ollama");
  const model = options.model || "";
  const promptText = options.prompt || "";
  const schemaText = options.schema ? JSON.stringify(options.schema) : "";
  return crypto.createHash("sha256").update(`${provider}|${model}|${url}|${promptText}|${schemaText}`).digest("hex");
}

function readExtractCache(url, options, ttlMs) {
  const file = path.join(extractCacheDir(), `${extractCacheKey(url, options)}.json`);
  if (!fs.existsSync(file)) return null;
  try {
    const stat = fs.statSync(file);
    const age = Date.now() - stat.mtimeMs;
    if (age > ttlMs) return null;
    const data = JSON.parse(fs.readFileSync(file, "utf8"));
    return { ...data, ageSeconds: age / 1000 };
  } catch {
    return null;
  }
}

function writeExtractCache(url, options, result) {
  const dir = extractCacheDir();
  fs.mkdirSync(dir, { recursive: true });
  const file = path.join(dir, `${extractCacheKey(url, options)}.json`);
  fs.writeFileSync(file, JSON.stringify({ url, parsed: result.parsed, errors: result.errors, raw: result.raw }, null, 2), "utf8");
}

export function clearExtractCache() {
  const dir = extractCacheDir();
  if (!fs.existsSync(dir)) return { cleared: 0 };
  let cleared = 0;
  for (const file of fs.readdirSync(dir)) {
    if (file.endsWith(".json")) {
      fs.unlinkSync(path.join(dir, file));
      cleared++;
    }
  }
  return { cleared, dir };
}

function parseTtl(value) {
  if (!value || value === "0") return 0;
  const match = String(value).match(/^(\d+)\s*([smhdw]?)$/i);
  if (!match) return 0;
  const num = Number(match[1]);
  const unit = (match[2] || "s").toLowerCase();
  const multiplier = { s: 1000, m: 60 * 1000, h: 60 * 60 * 1000, d: 24 * 60 * 60 * 1000, w: 7 * 24 * 60 * 60 * 1000 }[unit];
  return num * multiplier;
}

export async function extractFromText(text, options = {}) {
  const prompt = options.prompt || "Extract structured data from this content.";
  const schema = options.schema || null;
  const truncated = text.length > 60000 ? text.slice(0, 60000) + "\n…[truncated]" : text;

  const userMessage = schema
    ? `${prompt}\n\nReturn ONLY a JSON object that matches this schema:\n${JSON.stringify(schema, null, 2)}\n\nContent:\n${truncated}`
    : `${prompt}\n\nReturn ONLY valid JSON. Content:\n${truncated}`;

  const provider = pickProvider(options);
  if (options.verbose !== false) {
    process.stderr.write(`extract: using ${provider.name}${provider.model ? `@${provider.model}` : ""}\n`);
  }
  const raw = await provider.call(userMessage);
  const parsed = parseJSON(raw);
  if (schema) {
    const errors = validateAgainstSchema(parsed, schema);
    return { parsed, errors, raw };
  }
  return { parsed, errors: [], raw };
}

const BROWSER_UA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36";

export async function fetchHtml(url, options = {}) {
  const response = await fetch(url, {
    headers: {
      "User-Agent": options.userAgent || BROWSER_UA,
      "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
      "Accept-Language": options.acceptLanguage || "it-IT,it;q=0.9,en;q=0.8"
    }
  });
  if (!response.ok) {
    const error = new Error(`HTTP ${response.status} ${response.statusText} fetching ${url}. If anti-bot blocks you, try 'gutenberg record <url>' which uses a real browser.`);
    error.exitCode = 1;
    error.status = response.status;
    throw error;
  }
  const contentType = response.headers.get("content-type") || "";
  const charsetMatch = contentType.match(/charset=([^;]+)/i);
  const buffer = await response.arrayBuffer();
  let charset = charsetMatch ? charsetMatch[1].trim().toLowerCase() : null;
  if (!charset) {
    const sniff = new TextDecoder("ascii").decode(buffer.slice(0, 2048));
    const metaMatch = sniff.match(/<meta[^>]*charset=["']?([^"'>\s/]+)/i);
    if (metaMatch) charset = metaMatch[1].toLowerCase();
  }
  let text;
  try {
    text = new TextDecoder(charset || "utf-8", { fatal: false }).decode(buffer);
  } catch {
    text = new TextDecoder("utf-8", { fatal: false }).decode(buffer);
  }
  // Late-detect challenge pages (200 OK with anti-bot challenge body).
  const { detectChallenge } = await import("./seed.js");
  const marker = detectChallenge(text);
  if (marker) {
    const error = new Error(`Anti-bot challenge detected at ${url} (marker: ${marker}). Use 'gutenberg record <url>' (Playwright) or '--backend browserbase' to bypass.`);
    error.exitCode = 1;
    error.kind = "anti-bot";
    throw error;
  }
  return text;
}

export function stripHtml(html) {
  return html
    .replace(/<script[\s\S]*?<\/script>/gi, " ")
    .replace(/<style[\s\S]*?<\/style>/gi, " ")
    .replace(/<noscript[\s\S]*?<\/noscript>/gi, " ")
    .replace(/<!--[\s\S]*?-->/g, " ")
    .replace(/<[^>]+>/g, " ")
    .replace(/&nbsp;/g, " ")
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/\s+/g, " ")
    .trim();
}

function pickProvider(options) {
  if (options.provider === "ollama" || process.env.GUTENBERG_LLM_PROVIDER === "ollama") {
    return ollamaProvider(options.host || ollamaHost(), options.model || ollamaModel());
  }
  if (options.provider === "anthropic" || process.env.ANTHROPIC_API_KEY) {
    return anthropicProvider(process.env.ANTHROPIC_API_KEY, options.model || anthropicModel());
  }
  if (options.provider === "openai" || process.env.OPENAI_API_KEY) {
    return openaiProvider(process.env.OPENAI_API_KEY, options.model || openaiModel());
  }
  return ollamaProvider(ollamaHost(), ollamaModel());
}

function ollamaProvider(host, model) {
  return {
    name: "ollama",
    model,
    async call(message) {
      const response = await fetch(`${host.replace(/\/$/, "")}/api/chat`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          model,
          stream: false,
          format: "json",
          messages: [{ role: "user", content: message }]
        })
      }).catch((error) => {
        const wrapped = new Error(`Ollama unreachable at ${host}: ${error.message}. Run 'ollama serve' or set ANTHROPIC_API_KEY/OPENAI_API_KEY.`);
        wrapped.exitCode = 2;
        throw wrapped;
      });
      const data = await response.json();
      if (!response.ok) {
        const error = new Error(`Ollama error ${response.status}: ${JSON.stringify(data)}`);
        error.exitCode = 1;
        throw error;
      }
      return data.message?.content || data.response || "";
    }
  };
}

function anthropicProvider(apiKey, model) {
  if (!apiKey) {
    const error = new Error("ANTHROPIC_API_KEY is required");
    error.exitCode = 2;
    throw error;
  }
  return {
    name: "anthropic",
    model,
    async call(message) {
      const response = await fetch("https://api.anthropic.com/v1/messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "x-api-key": apiKey,
          "anthropic-version": "2023-06-01"
        },
        body: JSON.stringify({
          model,
          max_tokens: 4096,
          messages: [{ role: "user", content: message }]
        })
      });
      const data = await response.json();
      if (!response.ok) {
        const error = new Error(`Anthropic API error ${response.status}: ${JSON.stringify(data)}`);
        error.exitCode = 1;
        throw error;
      }
      const block = (data.content || []).find((item) => item.type === "text");
      return block?.text || "";
    }
  };
}

function openaiProvider(apiKey, model) {
  if (!apiKey) {
    const error = new Error("OPENAI_API_KEY is required");
    error.exitCode = 2;
    throw error;
  }
  return {
    name: "openai",
    model,
    async call(message) {
      const response = await fetch("https://api.openai.com/v1/chat/completions", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${apiKey}`
        },
        body: JSON.stringify({
          model,
          messages: [{ role: "user", content: message }],
          response_format: { type: "json_object" }
        })
      });
      const data = await response.json();
      if (!response.ok) {
        const error = new Error(`OpenAI API error ${response.status}: ${JSON.stringify(data)}`);
        error.exitCode = 1;
        throw error;
      }
      return data.choices?.[0]?.message?.content || "";
    }
  };
}

function parseJSON(text) {
  const trimmed = text.trim().replace(/^```(?:json)?\s*|\s*```$/g, "");
  try {
    return JSON.parse(trimmed);
  } catch {
    const match = trimmed.match(/\{[\s\S]*\}/);
    if (match) {
      try {
        return JSON.parse(match[0]);
      } catch {}
    }
    return { _raw: text };
  }
}

export function loadSchemaFromFlag(value) {
  if (!value) return null;
  if (value.startsWith("{") || value.startsWith("[")) return JSON.parse(value);
  return JSON.parse(fs.readFileSync(path.resolve(value), "utf8"));
}

export function validateAgainstSchema(value, schema) {
  const errors = [];
  validateNode(value, schema, "$", errors);
  return errors;
}

function validateNode(value, schema, ptr, errors) {
  if (!schema || typeof schema !== "object") return;
  if (Array.isArray(schema.required) && value && typeof value === "object" && !Array.isArray(value)) {
    for (const key of schema.required) {
      if (!(key in value)) errors.push({ path: ptr, message: `missing required property: ${key}` });
    }
  }
  if (schema.type) {
    const actual = typeOf(value);
    const expected = Array.isArray(schema.type) ? schema.type : [schema.type];
    if (!expected.includes(actual) && !(value === null && expected.includes("null"))) {
      errors.push({ path: ptr, message: `expected ${expected.join("|")}, got ${actual}` });
    }
  }
  if (schema.type === "object" && schema.properties && value && typeof value === "object") {
    for (const [key, sub] of Object.entries(schema.properties)) {
      if (key in value) validateNode(value[key], sub, `${ptr}.${key}`, errors);
    }
  }
  if (schema.type === "array" && schema.items && Array.isArray(value)) {
    for (let i = 0; i < value.length; i++) {
      validateNode(value[i], schema.items, `${ptr}[${i}]`, errors);
    }
  }
}

function typeOf(value) {
  if (value === null) return "null";
  if (Array.isArray(value)) return "array";
  if (typeof value === "number") return Number.isInteger(value) ? "integer" : "number";
  return typeof value;
}
