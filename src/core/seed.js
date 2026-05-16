import fs from "node:fs";
import path from "node:path";
import { writeJson } from "./fs.js";

const DEFAULT_METHOD = "GET";
const DEFAULT_USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36";

const CHALLENGE_MARKERS = [
  /Just a moment\.\.\./i,
  /Verifying you are human/i,
  /Checking your browser before/i,
  /cf-browser-verification/i,
  /<title>Attention Required! \| Cloudflare<\/title>/i,
  /Please enable JavaScript to view/i
];

export function detectChallenge(text) {
  if (!text) return null;
  for (const marker of CHALLENGE_MARKERS) {
    if (marker.test(text)) {
      return marker.source;
    }
  }
  return null;
}

export async function seedHarFromUrls(urls, options = {}) {
  if (!Array.isArray(urls) || urls.length === 0) {
    const error = new Error("seed-har requires at least one URL");
    error.exitCode = 2;
    throw error;
  }
  const userAgent = options.userAgent || DEFAULT_USER_AGENT;
  const headers = options.headers || {};
  const method = (options.method || DEFAULT_METHOD).toUpperCase();
  const entries = [];
  const warnings = [];

  for (const target of urls) {
    const startedAt = new Date();
    const response = await fetch(target, {
      method,
      headers: {
        Accept: "application/json, text/html;q=0.9, */*;q=0.8",
        "User-Agent": userAgent,
        ...headers
      },
      body: options.body && method !== "GET" ? options.body : undefined
    });
    const text = await response.text();
    const elapsed = Date.now() - startedAt.getTime();
    if (response.status >= 400) {
      warnings.push({ url: target, status: response.status, statusText: response.statusText, kind: "http-error" });
    }
    const challenge = detectChallenge(text);
    if (challenge) {
      warnings.push({ url: target, kind: "anti-bot", marker: challenge });
    }
    entries.push({
      pageref: "page_1",
      startedDateTime: startedAt.toISOString(),
      time: elapsed,
      request: {
        method,
        url: target,
        httpVersion: "HTTP/1.1",
        headers: headersToArray({ "User-Agent": userAgent, Accept: "application/json, text/html;q=0.9, */*;q=0.8", ...headers }),
        queryString: queryStringEntries(target),
        cookies: [],
        headersSize: -1,
        bodySize: options.body ? options.body.length : 0,
        postData: options.body && method !== "GET"
          ? {
              mimeType: headers["Content-Type"] || "application/json",
              text: options.body
            }
          : undefined
      },
      response: {
        status: response.status,
        statusText: response.statusText,
        httpVersion: "HTTP/1.1",
        headers: headersToArray(Object.fromEntries(response.headers.entries())),
        cookies: [],
        content: {
          size: text.length,
          mimeType: response.headers.get("content-type") || "application/json",
          text
        },
        redirectURL: response.headers.get("location") || "",
        headersSize: -1,
        bodySize: text.length
      },
      cache: {},
      timings: { send: 0, wait: elapsed, receive: 0 }
    });
  }

  const har = {
    log: {
      version: "1.2",
      creator: { name: "Gutenberg seed-har", version: "0.1.0" },
      pages: [
        {
          startedDateTime: new Date().toISOString(),
          id: "page_1",
          title: options.title || `Gutenberg seed (${urls.length} url${urls.length === 1 ? "" : "s"})`,
          pageTimings: {}
        }
      ],
      entries
    }
  };

  if (options.out) {
    const target = path.resolve(options.out);
    fs.mkdirSync(path.dirname(target), { recursive: true });
    writeJson(target, har);
  }
  if (options.goldenDir && entries.length > 0) {
    const dir = path.resolve(options.goldenDir);
    fs.mkdirSync(dir, { recursive: true });
    for (const entry of entries) {
      const slug = entry.request.url.replace(/^https?:\/\//, "").replace(/[^a-zA-Z0-9]+/g, "-").slice(0, 80);
      const file = path.join(dir, `${slug}.json`);
      fs.writeFileSync(file, entry.response.content.text || "", "utf8");
    }
  }
  har.warnings = warnings;
  return har;
}

function headersToArray(headers) {
  return Object.entries(headers).map(([name, value]) => ({ name, value: String(value) }));
}

function queryStringEntries(target) {
  try {
    const url = new URL(target);
    return [...url.searchParams.entries()].map(([name, value]) => ({ name, value }));
  } catch {
    return [];
  }
}
