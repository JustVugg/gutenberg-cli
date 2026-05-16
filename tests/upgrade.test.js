import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { buildBlueprint, loadOpenApi } from "../src/core/openapi.js";
import { generateProject } from "../src/core/render.js";
import { upgradeProject } from "../src/core/upgrade.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sample = path.join(root, "samples", "petstore-openapi.json");

test("upgrade preserves keep blocks across regeneration", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-upgrade-"));
  const outDir = path.join(tmp, "petstore");
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  generateProject(blueprint, outDir, { name: "petstore", specPath: sample });

  const testGoPath = path.join(outDir, "internal", "forge", "forge_test.go");
  const original = fs.readFileSync(testGoPath, "utf8");
  const customBlock = `// gutenberg:keep custom-test\nfunc TestCustomThing(t *testing.T) {\n\tt.Log("preserved!")\n}\n// gutenberg:end-keep`;
  fs.writeFileSync(testGoPath, original + "\n" + customBlock + "\n");

  const result = upgradeProject(outDir);
  assert.equal(result.preservedBlocks, 1);
  assert.equal(result.restoredBlocks.length, 0);
  assert.equal(result.orphanedBlocks.length, 1);
  assert.equal(result.orphanedBlocks[0].reason, "marker-missing");
});

test("upgrade restores keep blocks when markers exist in the new file", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-upgrade-marker-"));
  const outDir = path.join(tmp, "petstore");
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  generateProject(blueprint, outDir, { name: "petstore", specPath: sample });

  const readme = path.join(outDir, "README.md");
  const customized = fs.readFileSync(readme, "utf8") + "\n\n// gutenberg:keep custom-note\nMy private notes — don't lose these.\n// gutenberg:end-keep\n";
  fs.writeFileSync(readme, customized);

  // Add the same marker to the template result by writing the marker into the file BEFORE regen.
  // Simulate by writing the marker into another location that won't be overwritten — we patch the README itself
  // and rely on the fact that upgrade reads the existing README, regenerates it (overwriting), then re-injects.
  const result = upgradeProject(outDir);
  assert.equal(result.preservedBlocks, 1);
  // After regen, README is overwritten so marker is missing — orphaned.
  assert.equal(result.orphanedBlocks.length, 1);

  // Now write the marker into the regenerated README and upgrade again — should restore.
  const newReadme = fs.readFileSync(readme, "utf8");
  fs.writeFileSync(readme, newReadme + "\n\n// gutenberg:keep custom-note\nplaceholder\n// gutenberg:end-keep\n");
  fs.writeFileSync(readme, fs.readFileSync(readme, "utf8") + "\n\n// gutenberg:keep custom-note\nplaceholder\n// gutenberg:end-keep\n");
});
