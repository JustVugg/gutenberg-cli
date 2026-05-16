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

test("generates a Go CLI/MCP package by default", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "black-forge-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore" });

  assert.ok(fs.existsSync(path.join(result.outDir, "gutenberg.manifest.json")));
  assert.ok(fs.existsSync(path.join(result.outDir, "go.mod")));
  assert.ok(fs.existsSync(path.join(result.outDir, "cmd", "petstore", "main.go")));
  assert.ok(fs.existsSync(path.join(result.outDir, "internal", "forge", "mcp.go")));
  assert.ok(fs.existsSync(path.join(result.outDir, "internal", "forge", "auth.go")));
  assert.ok(fs.existsSync(path.join(result.outDir, "internal", "forge", "forge_test.go")));
  assert.equal(result.manifest.language, "go");
});

test("can still generate the Node compatibility target", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "black-forge-node-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore", lang: "node" });

  assert.ok(fs.existsSync(path.join(result.outDir, "src", "cli.js")));
  assert.ok(fs.existsSync(path.join(result.outDir, "src", "mcp-server.js")));
  assert.ok(fs.existsSync(path.join(result.outDir, "tests", "smoke.test.js")));
});

test("emits a Claude Skill by default", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-skill-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore" });

  const skillPath = path.join(result.outDir, "skills", "petstore", "SKILL.md");
  assert.ok(fs.existsSync(skillPath), "SKILL.md should exist");
  const body = fs.readFileSync(skillPath, "utf8");
  assert.match(body, /^---\nname: petstore\n/);
  assert.match(body, /description: /);
  assert.match(body, /## Operations index/);
  assert.match(body, /petstore call /);
  assert.deepEqual(result.targets, ["go", "mcp", "skill"]);
  assert.equal(result.skill.file, skillPath);
});

test("writes extended provenance to the manifest", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-prov-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore" });

  const manifest = JSON.parse(fs.readFileSync(path.join(result.outDir, "gutenberg.manifest.json"), "utf8"));
  assert.ok(manifest.provenance, "manifest.provenance must exist");
  assert.equal(manifest.provenance.schemaVersion, "gutenberg.provenance.v1");
  assert.match(manifest.provenance.generatedAt, /^\d{4}-\d{2}-\d{2}T/);
  assert.equal(typeof manifest.provenance.gutenbergVersion, "string");
  assert.ok(manifest.provenance.spec);
  assert.match(manifest.provenance.spec.sha256, /^[0-9a-f]{64}$/);
  assert.deepEqual(manifest.provenance.targets, ["go", "mcp", "skill"]);
  assert.ok(manifest.provenance.scorecard);
  assert.equal(typeof manifest.provenance.scorecard.score, "number");
});

test("respects --targets to skip the skill", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-noskill-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore", targets: "go,mcp" });

  assert.ok(!fs.existsSync(path.join(result.outDir, "skills", "petstore", "SKILL.md")));
  assert.equal(result.skill, undefined);
  assert.deepEqual(result.targets, ["go", "mcp"]);
});
