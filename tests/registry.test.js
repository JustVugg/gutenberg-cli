import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { buildRegistryFromCatalog } from "../src/core/catalog.js";
import { writeJson, writeText } from "../src/core/fs.js";

test("buildRegistryFromCatalog derives registry entries from manifests", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-registry-"));
  const toolDir = path.join(tmp, "library", "tools", "demo");
  writeJson(path.join(toolDir, "gutenberg.manifest.json"), {
    schemaVersion: "gutenberg.blueprint.v1",
    name: "Demo",
    slug: "demo",
    kind: "tool",
    description: "Demo tool",
    operations: [
      { id: "getDemo", method: "GET", path: "/demo", risk: "read", cacheable: true }
    ],
    heroes: [{ alias: "demo", operationId: "getDemo" }],
    tags: ["demo"],
    verification: { ok: true, proofFile: "proofs/verification.json" }
  });
  writeText(path.join(toolDir, "README.md"), "# Demo\n");

  const registry = buildRegistryFromCatalog(tmp);
  assert.equal(registry.schemaVersion, "gutenberg.registry.v1");
  assert.equal(registry.tools.length, 1);
  assert.equal(registry.tools[0].slug, "demo");
  assert.equal(registry.tools[0].package, "library/tools/demo");
  assert.equal(registry.tools[0].entrypoints.install, "gutenberg install library/tools/demo");
});
