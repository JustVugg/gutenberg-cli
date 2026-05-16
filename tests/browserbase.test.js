import test from "node:test";
import assert from "node:assert/strict";
import http from "node:http";
import { recordViaBrowserbase, browserbaseStatus } from "../src/core/recorder-browserbase.js";
import { skipIfNoLocalListen } from "./helpers/local-listen.js";

test("browserbaseStatus reports missing credentials", () => {
  const previousKey = process.env.BROWSERBASE_API_KEY;
  const previousProject = process.env.BROWSERBASE_PROJECT_ID;
  delete process.env.BROWSERBASE_API_KEY;
  delete process.env.BROWSERBASE_PROJECT_ID;
  try {
    const status = browserbaseStatus();
    assert.equal(status.apiKey, false);
    assert.equal(status.projectId, false);
    assert.equal(status.ready, false);
  } finally {
    if (previousKey) process.env.BROWSERBASE_API_KEY = previousKey;
    if (previousProject) process.env.BROWSERBASE_PROJECT_ID = previousProject;
  }
});

test("recordViaBrowserbase requires API key and project id", async () => {
  const previousKey = process.env.BROWSERBASE_API_KEY;
  const previousProject = process.env.BROWSERBASE_PROJECT_ID;
  delete process.env.BROWSERBASE_API_KEY;
  delete process.env.BROWSERBASE_PROJECT_ID;
  try {
    await assert.rejects(
      () => recordViaBrowserbase("https://example.com"),
      /BROWSERBASE_API_KEY and BROWSERBASE_PROJECT_ID/
    );
  } finally {
    if (previousKey) process.env.BROWSERBASE_API_KEY = previousKey;
    if (previousProject) process.env.BROWSERBASE_PROJECT_ID = previousProject;
  }
});

test("recordViaBrowserbase POSTs to /v1/sessions and surfaces failure", async (t) => {
  if (await skipIfNoLocalListen(t)) return;
  let captured = null;
  const server = http.createServer((req, res) => {
    let body = "";
    req.on("data", (chunk) => { body += chunk; });
    req.on("end", () => {
      captured = { method: req.method, url: req.url, headers: req.headers, body };
      res.statusCode = 400;
      res.setHeader("Content-Type", "application/json");
      res.end(JSON.stringify({ error: "invalid project" }));
    });
  });
  await new Promise((resolve) => server.listen(0, resolve));
  const port = server.address().port;
  t.after(() => server.close());

  await assert.rejects(
    () => recordViaBrowserbase("https://example.com", {
      apiKey: "test-key",
      projectId: "proj_xxx",
      apiUrl: `http://127.0.0.1:${port}`
    }),
    /Browserbase session create failed: 400/
  );

  assert.ok(captured);
  assert.equal(captured.method, "POST");
  assert.equal(captured.url, "/v1/sessions");
  assert.equal(captured.headers["x-bb-api-key"], "test-key");
  const payload = JSON.parse(captured.body);
  assert.equal(payload.projectId, "proj_xxx");
});
