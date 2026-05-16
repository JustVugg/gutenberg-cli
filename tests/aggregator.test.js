import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { generateAggregatorProject, loadAggregatorRecipe } from "../src/core/render-aggregator.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sample = path.join(root, "samples", "petstore-openapi.json");

test("generates a Go aggregator over multiple OpenAPI sources", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-aggr-"));
  const recipePath = path.join(tmp, "flights.recipe.json");
  fs.writeFileSync(recipePath, JSON.stringify({
    schemaVersion: "gutenberg.aggregator.v1",
    name: "flights",
    displayName: "Flights Aggregator",
    description: "Demo aggregator over two petstore copies",
    sources: [
      { name: "alpha", spec: sample, operation: "listPets" },
      { name: "beta", spec: sample, operation: "listPets" }
    ],
    merge: "concatenate",
    rank: "by-source-order"
  }));

  const recipe = loadAggregatorRecipe(recipePath, fs);
  const outDir = path.join(tmp, "out");
  const result = generateAggregatorProject(recipe, outDir, { name: "flights", baseDir: tmp });

  assert.equal(result.manifest.kind, "aggregator");
  assert.equal(result.manifest.sources.length, 2);
  assert.ok(fs.existsSync(path.join(outDir, "go.mod")));
  assert.ok(fs.existsSync(path.join(outDir, "cmd", "flights", "main.go")));
  for (const file of ["manifests.go", "client.go", "merge.go", "rank.go", "types.go", "aggr_test.go"]) {
    assert.ok(fs.existsSync(path.join(outDir, "internal", "aggr", file)), `missing ${file}`);
  }

  const manifestsGo = fs.readFileSync(path.join(outDir, "internal", "aggr", "manifests.go"), "utf8");
  assert.match(manifestsGo, /SourceDescriptor\{/);
  assert.match(manifestsGo, /Slug:\s*"alpha"/);
  assert.match(manifestsGo, /Slug:\s*"beta"/);

  const mainGo = fs.readFileSync(path.join(outDir, "cmd", "flights", "main.go"), "utf8");
  assert.match(mainGo, /aggr\.FanOut/);
  assert.match(mainGo, /aggr\.Sources/);
});

test("rejects an empty sources array", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-aggr-bad-"));
  const recipePath = path.join(tmp, "bad.recipe.json");
  fs.writeFileSync(recipePath, JSON.stringify({
    schemaVersion: "gutenberg.aggregator.v1",
    name: "bad",
    sources: []
  }));
  assert.throws(() => loadAggregatorRecipe(recipePath, fs), /sources/);
});
