import test from "node:test";
import assert from "node:assert/strict";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { buildBlueprint, loadOpenApi } from "../src/core/openapi.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sample = path.join(root, "samples", "petstore-openapi.json");
const sampleYaml = path.join(root, "samples", "petstore-openapi.yaml");

test("builds a blueprint from OpenAPI JSON", () => {
  const blueprint = buildBlueprint(loadOpenApi(sample), sample, "Petstore");
  assert.equal(blueprint.slug, "petstore");
  assert.equal(blueprint.operations.length, 6);
  assert.deepEqual(blueprint.tags, ["orders", "pets"]);
  assert.equal(blueprint.operations.find((operation) => operation.id === "deletePet").risk, "destructive");
  assert.ok(blueprint.insights.thesis.includes("local"));
});

test("loads OpenAPI YAML", () => {
  const blueprint = buildBlueprint(loadOpenApi(sampleYaml), sampleYaml, "YAML Petstore");
  assert.equal(blueprint.slug, "yaml-petstore");
  assert.equal(blueprint.operations.length, 2);
});
