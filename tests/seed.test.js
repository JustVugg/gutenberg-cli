import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import http from "node:http";
import os from "node:os";
import path from "node:path";
import { seedHarFromUrls } from "../src/core/seed.js";
import { skipIfNoLocalListen } from "./helpers/local-listen.js";

test("seed-har produces a HAR from plain HTTP fetches", async (t) => {
  if (await skipIfNoLocalListen(t)) return;
  const server = http.createServer((req, res) => {
    res.setHeader("Content-Type", "application/json");
    res.end(JSON.stringify({ path: req.url, ok: true }));
  });
  await new Promise((resolve) => server.listen(0, resolve));
  const port = server.address().port;
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-seed-"));
  const out = path.join(tmp, "capture.har.json");
  try {
    const har = await seedHarFromUrls([
      `http://127.0.0.1:${port}/v1/foo?a=1`,
      `http://127.0.0.1:${port}/v1/bar`
    ], { out });
    assert.equal(har.log.entries.length, 2);
    assert.equal(har.log.entries[0].request.method, "GET");
    assert.equal(har.log.entries[0].response.status, 200);
    assert.match(har.log.entries[0].response.content.text, /"ok":true/);
    assert.ok(fs.existsSync(out));
    const persisted = JSON.parse(fs.readFileSync(out, "utf8"));
    assert.equal(persisted.log.entries.length, 2);
  } finally {
    server.close();
  }
});

test("seed-har rejects empty url list", async () => {
  await assert.rejects(() => seedHarFromUrls([]), /at least one URL/);
});
