import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import http from "node:http";
import os from "node:os";
import path from "node:path";
import { extractFromUrl, clearExtractCache } from "../src/core/extract.js";
import { skipIfNoLocalListen } from "./helpers/local-listen.js";

let llmCalls = 0;
let httpServer;
let llmServer;
const cacheDir = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-extract-cache-"));

test("extract uses cache on second call when TTL is set", async (t) => {
  if (await skipIfNoLocalListen(t)) return;
  process.env.GUTENBERG_EXTRACT_CACHE_DIR = cacheDir;
  process.env.GUTENBERG_LLM_PROVIDER = "ollama";

  httpServer = http.createServer((req, res) => {
    res.setHeader("Content-Type", "text/html");
    res.end("<html><body><h1>Cache test</h1></body></html>");
  });
  await new Promise((resolve) => httpServer.listen(0, resolve));
  const httpPort = httpServer.address().port;

  llmServer = http.createServer((req, res) => {
    let body = "";
    req.on("data", (chunk) => { body += chunk; });
    req.on("end", () => {
      llmCalls += 1;
      res.setHeader("Content-Type", "application/json");
      res.end(JSON.stringify({ message: { content: '{"hits": 1}' } }));
    });
  });
  await new Promise((resolve) => llmServer.listen(0, resolve));
  const llmPort = llmServer.address().port;
  process.env.OLLAMA_HOST = `http://127.0.0.1:${llmPort}`;

  t.after(() => {
    httpServer.close();
    llmServer.close();
    clearExtractCache();
    delete process.env.GUTENBERG_EXTRACT_CACHE_DIR;
    delete process.env.GUTENBERG_LLM_PROVIDER;
    delete process.env.OLLAMA_HOST;
  });

  const url = `http://127.0.0.1:${httpPort}/`;
  const first = await extractFromUrl(url, { prompt: "Extract", cache: "1h", verbose: false });
  assert.equal(first.cached, undefined);
  assert.equal(llmCalls, 1);

  const second = await extractFromUrl(url, { prompt: "Extract", cache: "1h", verbose: false });
  assert.equal(second.cached, true);
  assert.equal(llmCalls, 1, "second call must NOT hit the LLM");
});

test("clearExtractCache removes entries", async () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-extract-clear-"));
  process.env.GUTENBERG_EXTRACT_CACHE_DIR = tmp;
  fs.writeFileSync(path.join(tmp, "a.json"), "{}");
  fs.writeFileSync(path.join(tmp, "b.json"), "{}");
  const result = clearExtractCache();
  assert.equal(result.cleared, 2);
  delete process.env.GUTENBERG_EXTRACT_CACHE_DIR;
});
