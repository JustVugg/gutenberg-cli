import test from "node:test";
import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const cli = path.join(root, "bin", "gutenberg.js");
const sample = path.join(root, "samples", "petstore-openapi.json");

test("--strict passes when scorecard ratio >= min-score", (t) => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-strict-"));
  const out = path.join(tmp, "petstore-strict");
  const result = spawnSync(process.execPath, [
    cli, "generate", sample,
    "--out", out, "--name", "petstore", "--force",
    "--strict", "--min-score", "0.7"
  ], { encoding: "utf8" });
  if (result.error?.code === "EPERM") {
    t.skip("child process spawn is unavailable in this sandbox");
    return;
  }
  assert.equal(result.status, 0, `stderr: ${result.stderr}\nstdout: ${result.stdout}`);
});

test("--strict fails with exit 3 when min-score is unreachable", (t) => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-strict-fail-"));
  const out = path.join(tmp, "petstore-strict-fail");
  const result = spawnSync(process.execPath, [
    cli, "generate", sample,
    "--out", out, "--name", "petstore", "--force",
    "--strict", "--min-score", "1.5"
  ], { encoding: "utf8" });
  if (result.error?.code === "EPERM") {
    t.skip("child process spawn is unavailable in this sandbox");
    return;
  }
  assert.equal(result.status, 3, `expected exit 3, got ${result.status}\nstderr: ${result.stderr}`);
  assert.match(result.stderr, /Strict mode failed/);
});
