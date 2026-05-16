import fs from "node:fs";
import path from "node:path";
import { writeJson } from "./fs.js";

const browserbaseApi = () => process.env.BROWSERBASE_API_URL || "https://api.browserbase.com";

export function browserbaseStatus() {
  return {
    apiKey: Boolean(process.env.BROWSERBASE_API_KEY),
    projectId: Boolean(process.env.BROWSERBASE_PROJECT_ID),
    apiUrl: browserbaseApi(),
    ready: Boolean(process.env.BROWSERBASE_API_KEY) && Boolean(process.env.BROWSERBASE_PROJECT_ID)
  };
}

export async function recordViaBrowserbase(url, options = {}) {
  const apiKey = options.apiKey || process.env.BROWSERBASE_API_KEY;
  const projectId = options.projectId || process.env.BROWSERBASE_PROJECT_ID;
  if (!apiKey || !projectId) {
    const error = new Error(
      "Browserbase requires BROWSERBASE_API_KEY and BROWSERBASE_PROJECT_ID env vars.\n" +
      "Get them at https://www.browserbase.com/ → Dashboard → Settings."
    );
    error.exitCode = 2;
    throw error;
  }
  const BB_API = options.apiUrl || browserbaseApi();

  let playwright;
  try {
    playwright = await import("playwright");
  } catch {
    const error = new Error("Playwright not installed. Run `npm install` in the Gutenberg project.");
    error.exitCode = 2;
    throw error;
  }

  const sessionResponse = await fetch(`${BB_API.replace(/\/$/, "")}/v1/sessions`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-BB-API-Key": apiKey
    },
    body: JSON.stringify({ projectId, browserSettings: { viewport: { width: 1440, height: 1000 } } })
  });
  if (!sessionResponse.ok) {
    const text = await sessionResponse.text();
    const error = new Error(`Browserbase session create failed: ${sessionResponse.status} ${text}`);
    error.exitCode = 1;
    throw error;
  }
  const session = await sessionResponse.json();
  const wsEndpoint = session.connectUrl || session.wsEndpoint;
  if (!wsEndpoint) {
    const error = new Error(`Browserbase response missing connect URL: ${JSON.stringify(session)}`);
    error.exitCode = 1;
    throw error;
  }

  const browser = await playwright.chromium.connectOverCDP(wsEndpoint);
  const context = browser.contexts()[0] || await browser.newContext();
  const page = await context.newPage();
  const entries = [];
  page.on("requestfinished", async (request) => {
    const response = await request.response().catch(() => null);
    if (!response) return;
    entries.push(await toHarEntry(request, response));
  });

  try {
    await page.goto(url, { waitUntil: "domcontentloaded", timeout: Number(options.timeout || 60000) });
    if (options.wait) await page.waitForTimeout(Number(options.wait));
  } finally {
    await context.close().catch(() => {});
    await browser.close().catch(() => {});
    await fetch(`${BB_API.replace(/\/$/, "")}/v1/sessions/${session.id}`, {
      method: "DELETE",
      headers: { "X-BB-API-Key": apiKey }
    }).catch(() => {});
  }

  const har = {
    log: {
      version: "1.2",
      creator: { name: "Gutenberg Browserbase recorder", version: "0.1.0" },
      pages: [{ startedDateTime: new Date().toISOString(), id: "page_1", title: url, pageTimings: {} }],
      entries
    }
  };
  if (options.out) {
    const target = path.resolve(options.out);
    fs.mkdirSync(path.dirname(target), { recursive: true });
    writeJson(target, har);
  }
  return har;
}

async function toHarEntry(request, response) {
  const requestUrl = new URL(request.url());
  const requestHeaders = await request.allHeaders().catch(() => ({}));
  const responseHeaders = await response.allHeaders().catch(() => ({}));
  return {
    pageref: "page_1",
    startedDateTime: new Date().toISOString(),
    time: 0,
    request: {
      method: request.method(),
      url: request.url(),
      httpVersion: "HTTP/2",
      headers: Object.entries(requestHeaders).map(([name, value]) => ({ name, value })),
      queryString: [...requestUrl.searchParams.entries()].map(([name, value]) => ({ name, value })),
      cookies: [],
      headersSize: -1,
      bodySize: request.postDataBuffer()?.length || 0
    },
    response: {
      status: response.status(),
      statusText: response.statusText(),
      httpVersion: "HTTP/2",
      headers: Object.entries(responseHeaders).map(([name, value]) => ({ name, value })),
      cookies: [],
      content: {
        size: Number(responseHeaders["content-length"] || 0),
        mimeType: responseHeaders["content-type"] || ""
      },
      redirectURL: responseHeaders.location || "",
      headersSize: -1,
      bodySize: Number(responseHeaders["content-length"] || 0)
    },
    cache: {},
    timings: { send: 0, wait: 0, receive: 0 }
  };
}
