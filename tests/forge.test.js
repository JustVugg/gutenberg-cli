import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { forgeProject, runPlan } from "../src/core/forge.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sample = path.join(root, "samples", "petstore-openapi.json");

test("runPlan exposes a safe generation plan", async () => {
  const plan = await runPlan(sample, { name: "petstore-plan" });
  assert.equal(plan.kind, "openapi");
  assert.equal(plan.auth.mode, "detected");
  assert.ok(plan.operations.some((op) => op.risk === "destructive"));
  assert.ok(plan.nextCommands[0].includes("gutenberg forge"));
});

test("forgeProject generates, verifies and returns an executable command", { timeout: 180000 }, async () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-forge-"));
  const out = path.join(tmp, "petstore-go");
  const result = await forgeProject(sample, {
    name: "petstore-forge",
    out,
    force: true
  });
  assert.equal(result.verification.ok, true);
  assert.ok(fs.existsSync(path.join(out, "proofs", "verification.json")));
  assert.ok(result.command.includes("petstore-forge"));
  assert.equal(result.verification.scorecard.grade, "A");
});
