import test from "node:test";
import assert from "node:assert/strict";
import { installTool } from "../src/core/install.js";

test("installTool rejects missing manifest", () => {
  assert.throws(() => installTool("/tmp/does-not-exist-gutenberg"), /No gutenberg.manifest.json/);
});
