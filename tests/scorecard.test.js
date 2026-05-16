import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { buildBlueprint, loadOpenApi } from "../src/core/openapi.js";
import { generateProject } from "../src/core/render.js";
import { scoreProject } from "../src/core/scorecard.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sample = path.join(root, "samples", "petstore-openapi.json");

test("scores generated packages", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "black-forge-score-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore" });
  const score = scoreProject(result.outDir);
  assert.notEqual(score.grade, "A", "unverified generated packages must not receive grade A");
  assert.ok(score.score > 0);
  assert.equal(score.checks.some((check) => check.id === "manifest:insights" && check.passed), true);
  assert.equal(score.checks.some((check) => check.id === "verify:proof" && !check.passed), true);
});

test("scorecard exposes dimensions", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-dim-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore" });
  const score = scoreProject(result.outDir);

  assert.ok(score.dimensions);
  const expected = ["structure", "manifest", "runtime", "verification", "coverage", "safety", "examples", "skill"];
  for (const dim of expected) {
    assert.ok(score.dimensions[dim], `missing dimension: ${dim}`);
    assert.equal(typeof score.dimensions[dim].score, "number");
    assert.equal(typeof score.dimensions[dim].max, "number");
  }
  assert.equal(score.dimensions.skill.present, true);
  assert.ok(score.ratio >= 0.7, `expected ratio >= 0.7, got ${score.ratio}`);
});

test("verified packages receive verification badges and grade A", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-verify-score-"));
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  const result = generateProject(blueprint, path.join(tmp, "petstore"), { name: "petstore" });
  fs.mkdirSync(path.join(result.outDir, "proofs"), { recursive: true });
  fs.writeFileSync(path.join(result.outDir, "proofs", "verification.json"), JSON.stringify({
    schemaVersion: "gutenberg.verification.v1",
    ok: true,
    checks: [
      { id: "go-build", passed: true },
      { id: "cli-smoke", passed: true },
      { id: "mcp-handshake", passed: true },
      { id: "go-test", passed: true }
    ]
  }, null, 2));

  const score = scoreProject(result.outDir);
  assert.equal(score.grade, "A");
  assert.ok(score.badges.includes("Gutenberg Verified"));
  assert.ok(score.badges.includes("Build Verified"));
  assert.ok(score.badges.includes("MCP Ready"));
});
