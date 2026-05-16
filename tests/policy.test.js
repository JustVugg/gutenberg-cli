import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { resolvePolicy } from "../src/core/source.js";

test("resolvePolicy loads and normalizes a policy file", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-policy-"));
  const file = path.join(tmp, "policy.json");
  fs.writeFileSync(file, JSON.stringify({
    rules: [
      { risk: "read", action: "allow" },
      { risk: "write", action: "deny" }
    ],
    redaction: ["Authorization"]
  }));

  const policy = resolvePolicy({ policy: file });
  assert.equal(policy.schemaVersion, "gutenberg.policy.v1");
  assert.equal(policy.rules.find((rule) => rule.risk === "write").action, "deny");
  assert.equal(policy.rules.find((rule) => rule.risk === "destructive").action, "confirm");
  assert.deepEqual(policy.redaction, ["authorization"]);
});

test("resolvePolicy rejects duplicate risk rules", () => {
  assert.throws(
    () => resolvePolicy({ policy: { rules: [
      { risk: "read", action: "allow" },
      { risk: "read", action: "deny" }
    ] } }),
    /duplicate rule/
  );
});
