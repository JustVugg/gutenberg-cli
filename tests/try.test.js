import test from "node:test";
import assert from "node:assert/strict";
import http from "node:http";
import { tryUrl } from "../src/core/try.js";
import { skipIfNoLocalListen } from "./helpers/local-listen.js";

async function startServer(handler) {
  const server = http.createServer(handler);
  await new Promise((resolve) => server.listen(0, resolve));
  return { server, port: server.address().port };
}

test("try classifies a JSON endpoint", async (t) => {
  if (await skipIfNoLocalListen(t)) return;
  const { server, port } = await startServer((req, res) => {
    res.setHeader("Content-Type", "application/json");
    res.end(JSON.stringify({ ok: true }));
  });
  t.after(() => server.close());
  const report = await tryUrl(`http://127.0.0.1:${port}/`);
  assert.equal(report.verdict, "json-endpoint");
  assert.equal(report.confidence, "high");
});

test("try classifies plain HTML content", async (t) => {
  if (await skipIfNoLocalListen(t)) return;
  const { server, port } = await startServer((req, res) => {
    res.setHeader("Content-Type", "text/html; charset=UTF-8");
    res.end("<html><body><h1>Article Title</h1><p>" + "Plenty of content. ".repeat(200) + "</p></body></html>");
  });
  t.after(() => server.close());
  const report = await tryUrl(`http://127.0.0.1:${port}/`);
  assert.equal(report.verdict, "html-content");
});

test("try detects anti-bot challenge", async (t) => {
  if (await skipIfNoLocalListen(t)) return;
  const { server, port } = await startServer((req, res) => {
    res.statusCode = 403;
    res.setHeader("Content-Type", "text/html");
    res.end("<html><head><title>Just a moment...</title></head><body>cf-browser-verification</body></html>");
  });
  t.after(() => server.close());
  const report = await tryUrl(`http://127.0.0.1:${port}/`);
  assert.equal(report.verdict, "anti-bot-challenge");
});

test("try detects SPA shell", async (t) => {
  if (await skipIfNoLocalListen(t)) return;
  const { server, port } = await startServer((req, res) => {
    res.setHeader("Content-Type", "text/html");
    res.end(`<html><body><div id="root"></div><script src="/runtime.app.mjs"></script></body></html>`);
  });
  t.after(() => server.close());
  const report = await tryUrl(`http://127.0.0.1:${port}/`);
  assert.equal(report.verdict, "spa");
});
