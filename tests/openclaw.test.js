import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { buildBlueprint, loadOpenApi } from "../src/core/openapi.js";
import { generateProject } from "../src/core/render.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sample = path.join(root, "samples", "petstore-openapi.json");

test("emits an OpenClaw skill when targets include openclaw", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-openclaw-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), {
    name: "petstore",
    targets: "go,mcp,skill,openclaw"
  });

  const jsonFile = path.join(result.outDir, "openclaw", "petstore", "skill.json");
  const mdFile = path.join(result.outDir, "openclaw", "petstore", "skill.md");
  assert.ok(fs.existsSync(jsonFile), "openclaw skill.json must exist");
  assert.ok(fs.existsSync(mdFile), "openclaw skill.md must exist");

  const skill = JSON.parse(fs.readFileSync(jsonFile, "utf8"));
  assert.equal(skill.schemaVersion, "openclaw.skill.v1");
  assert.equal(skill.name, "petstore");
  assert.ok(skill.runtime.cli);
  assert.ok(typeof skill.capabilities.total === "number");
  assert.ok(Array.isArray(skill.actions));
  assert.equal(skill.safety.dryRunByDefault, true);

  const md = fs.readFileSync(mdFile, "utf8");
  assert.match(md, /OpenClaw skill/);
  assert.match(md, /Operations:/);
});

test("openclaw target is opt-in (skipped by default)", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-no-openclaw-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore" });
  assert.equal(result.openclaw, undefined);
  assert.equal(fs.existsSync(path.join(result.outDir, "openclaw")), false);
});
