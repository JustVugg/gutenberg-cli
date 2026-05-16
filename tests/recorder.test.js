import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import path from "node:path";
import os from "node:os";
import { recordBrowser } from "../src/core/recorder.js";

import { loginInteractive } from "../src/core/recorder.js";

test("login requires --out", async () => {
  await assert.rejects(
    () => loginInteractive("https://example.com", {}),
    /requires --out/
  );
});

test("rejects missing storage-state file before launching", async () => {
  const missing = path.join(os.tmpdir(), `gutenberg-missing-${Date.now()}.json`);
  await assert.rejects(
    () => recordBrowser("https://example.com", { storageState: missing }),
    /storage-state file not found/
  );
});

test("accepts an existing storage-state file (smoke; opt-in via GUTENBERG_BROWSER_TESTS=1)", { skip: process.env.GUTENBERG_BROWSER_TESTS !== "1" }, async () => {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-record-"));
  const statePath = path.join(dir, "state.json");
  fs.writeFileSync(statePath, JSON.stringify({ cookies: [], origins: [] }));
  await recordBrowser("about:blank", {
    storageState: statePath,
    timeout: 5000,
    out: path.join(dir, "out.har.json")
  });
  assert.ok(fs.existsSync(path.join(dir, "out.har.json")));
});
