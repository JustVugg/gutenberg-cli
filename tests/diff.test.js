import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { diffSpecs } from "../src/core/diff.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sample = path.join(root, "samples", "petstore-openapi.json");

test("diff reports identical specs as empty", () => {
  const result = diffSpecs(sample, sample);
  assert.equal(result.counts.added, 0);
  assert.equal(result.counts.removed, 0);
  assert.equal(result.counts.changed, 0);
});

test("diff detects added/removed/changed operations", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-diff-"));
  const original = JSON.parse(fs.readFileSync(sample, "utf8"));
  const tweaked = JSON.parse(JSON.stringify(original));

  const firstPath = Object.keys(tweaked.paths)[0];
  delete tweaked.paths[firstPath];
  tweaked.paths["/newpath"] = {
    get: {
      operationId: "newOperation",
      summary: "A new operation",
      responses: { "200": { description: "ok" } }
    }
  };

  const tweakedPath = path.join(tmp, "tweaked.json");
  fs.writeFileSync(tweakedPath, JSON.stringify(tweaked));

  const result = diffSpecs(sample, tweakedPath);
  assert.ok(result.counts.added >= 1);
  assert.ok(result.counts.removed >= 1);
  assert.ok(result.added.some((op) => op.id === "newOperation"));
});
