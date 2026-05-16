import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { buildBlueprint, loadOpenApi } from "../src/core/openapi.js";
import { generateProject } from "../src/core/render.js";
import { ensureAgentAssets, publishTool } from "../src/core/publish.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sample = path.join(root, "samples", "petstore-openapi.json");

test("ensureAgentAssets adds skill and OpenClaw assets", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-publish-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), {
    name: "petstore",
    targets: "go"
  });

  const fixed = ensureAgentAssets(result.outDir);
  assert.equal(fixed.slug, "petstore");
  assert.ok(fixed.changed.some((file) => file.endsWith("SKILL.md")));
  assert.ok(fs.existsSync(path.join(result.outDir, "skills", "petstore", "SKILL.md")));
  assert.ok(fs.existsSync(path.join(result.outDir, "openclaw", "petstore", "skill.json")));
});

test("publishTool blocks when verification proof is missing", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-publish-blocked-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), {
    name: "petstore",
    targets: "go,mcp,skill,openclaw"
  });

  const published = publishTool(result.outDir);
  assert.equal(published.ready, false);
  assert.ok(published.errors.some((error) => /verification\.json/.test(error)));
});
