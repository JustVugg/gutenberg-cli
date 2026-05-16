import fs from "node:fs";
import path from "node:path";
import { writeJson } from "./fs.js";

export async function recordBrowser(url, options = {}) {
  const contextOptions = {
    viewport: { width: 1440, height: 1000 },
    userAgent: options.userAgent
  };
  if (options.storageState) {
    const resolved = path.resolve(options.storageState);
    if (!fs.existsSync(resolved)) {
      const error = new Error(`storage-state file not found: ${resolved}`);
      error.exitCode = 2;
      throw error;
    }
    contextOptions.storageState = resolved;
  }

  let playwright;
  try {
    playwright = await import("playwright");
  } catch {
    const error = new Error("Playwright is not installed. Run `npm install` in the Gutenberg project.");
    error.exitCode = 2;
    throw error;
  }

  const browser = await playwright.chromium.launch({
    headless: options.headless !== false
  });
  const context = await browser.newContext(contextOptions);
  const page = await context.newPage();
  const entries = [];
  let navigationWarning = null;

  page.on("requestfinished", async (request) => {
    const response = await request.response().catch(() => null);
    if (!response) return;
    entries.push(await requestToHarEntry(request, response));
  });

  try {
    await navigate(page, url, options);
    if (options.wait) {
      await page.waitForTimeout(Number(options.wait));
    }
  } catch (error) {
    navigationWarning = error.message || String(error);
    if (options.wait) {
      await page.waitForTimeout(Number(options.wait)).catch(() => {});
    }
  } finally {
    if (options.saveStorageState) {
      const target = path.resolve(options.saveStorageState);
      fs.mkdirSync(path.dirname(target), { recursive: true });
      await context.storageState({ path: target }).catch(() => {});
    }
    await context.close();
    await browser.close();
  }

  const har = {
    log: {
      version: "1.2",
      creator: {
        name: "Gutenberg Playwright recorder",
        version: "0.1.0"
      },
      pages: [
        {
          startedDateTime: new Date().toISOString(),
          id: "page_1",
          title: url,
          pageTimings: {}
        }
      ],
      entries,
      comment: navigationWarning ? `Navigation warning: ${navigationWarning}` : undefined
    }
  };

  if (options.out) {
    writeJson(options.out, har);
  }
  return har;
}

async function navigate(page, url, options) {
  const timeout = Number(options.timeout || 60000);
  const preferred = options.waitUntil || "domcontentloaded";
  const waitUntilOptions = unique([preferred, "domcontentloaded", "load"]);
  let lastError = null;

  for (const waitUntil of waitUntilOptions) {
    try {
      await page.goto(url, { waitUntil, timeout });
      return;
    } catch (error) {
      lastError = error;
      if (page.isClosed()) {
        throw error;
      }
    }
  }
  throw lastError;
}

function unique(values) {
  return values.filter((value, index) => value && values.indexOf(value) === index);
}

async function requestToHarEntry(request, response) {
  const requestUrl = new URL(request.url());
  const responseHeaders = await response.allHeaders().catch(() => ({}));
  const requestHeaders = await request.allHeaders().catch(() => ({}));
  return {
    pageref: "page_1",
    startedDateTime: new Date().toISOString(),
    time: response.timing?.()?.responseEnd || 0,
    request: {
      method: request.method(),
      url: request.url(),
      httpVersion: "HTTP/2",
      headers: objectHeaders(requestHeaders),
      queryString: [...requestUrl.searchParams.entries()].map(([name, value]) => ({ name, value })),
      cookies: [],
      headersSize: -1,
      bodySize: request.postDataBuffer()?.length || 0,
      postData: request.postData()
        ? {
            mimeType: requestHeaders["content-type"] || "application/octet-stream",
            text: request.postData()
          }
        : undefined
    },
    response: {
      status: response.status(),
      statusText: response.statusText(),
      httpVersion: "HTTP/2",
      headers: objectHeaders(responseHeaders),
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
    timings: {
      send: 0,
      wait: 0,
      receive: 0
    }
  };
}

function objectHeaders(headers) {
  return Object.entries(headers).map(([name, value]) => ({ name, value }));
}

export async function loginInteractive(url, options = {}) {
  const out = options.out;
  if (!out) {
    const error = new Error("login requires --out <state.json> to save the storage state");
    error.exitCode = 2;
    throw error;
  }
  const target = path.resolve(out);
  fs.mkdirSync(path.dirname(target), { recursive: true });

  let playwright;
  try {
    playwright = await import("playwright");
  } catch {
    const error = new Error("Playwright is not installed. Run `npm install` in the Gutenberg project.");
    error.exitCode = 2;
    throw error;
  }

  const browser = await playwright.chromium.launch({ headless: false });
  const contextOptions = {
    viewport: { width: 1440, height: 1000 },
    userAgent: options.userAgent
  };
  if (options.storageState) {
    const resolved = path.resolve(options.storageState);
    if (fs.existsSync(resolved)) contextOptions.storageState = resolved;
  }
  const context = await browser.newContext(contextOptions);
  const page = await context.newPage();
  await page.goto(url, { waitUntil: "domcontentloaded", timeout: Number(options.timeout || 60000) }).catch(() => {});

  const closed = new Promise((resolve) => {
    let done = false;
    const finish = () => {
      if (done) return;
      done = true;
      resolve();
    };
    page.on("close", finish);
    context.on("close", finish);
    if (options.waitForUrl) {
      page.waitForURL(options.waitForUrl, { timeout: 0 }).then(finish).catch(() => {});
    }
    if (typeof options.onReady === "function") {
      options.onReady({ done: finish });
    } else if (process.stdin.isTTY) {
      process.stdout.write(`\nLog in inside the browser. Press Enter here when done.\n`);
      process.stdin.once("data", finish);
      process.stdin.resume();
    }
  });

  await closed;
  await context.storageState({ path: target }).catch((error) => {
    console.warn(`Could not save storage-state: ${error.message}`);
  });
  if (process.stdin.isTTY) process.stdin.pause();
  await context.close().catch(() => {});
  await browser.close().catch(() => {});
  return { out: target };
}
