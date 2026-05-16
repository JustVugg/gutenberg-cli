import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { buildBlueprint, loadOpenApi } from "../src/core/openapi.js";
import { generateProject } from "../src/core/render.js";
import { readLock } from "../src/core/lockfile.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sample = path.join(root, "samples", "petstore-openapi.json");

test("writes a lockfile with spec hash and version", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-lock-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore" });

  const lock = readLock(result.outDir);
  assert.ok(lock, "lock file should exist");
  assert.equal(lock.schemaVersion, "gutenberg.lock.v1");
  assert.match(lock.spec.sha256, /^[0-9a-f]{64}$/);
  assert.deepEqual(lock.targets, ["go", "mcp", "skill"]);
  assert.equal(typeof lock.gutenbergVersion, "string");
});

test("detects drift when the spec content changes", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-lock-drift-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const first = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore" });
  assert.equal(first.drift.drifted, false);

  const tweaked = path.join(tmp, "tweaked.json");
  const original = JSON.parse(fs.readFileSync(sample, "utf8"));
  original.info.title = `${original.info.title} (modified)`;
  fs.writeFileSync(tweaked, JSON.stringify(original, null, 2));
  const blueprint2 = buildBlueprint(loadOpenApi(tweaked), tweaked, "Petstore");
  const second = generateProject(blueprint2, path.join(tmp, "petstore"), { name: "petstore", force: true });

  assert.equal(second.drift.drifted, true);
  const specDrift = second.drift.drifts.find((d) => d.kind === "spec");
  assert.ok(specDrift);
  assert.notEqual(specDrift.locked, specDrift.current);
});
