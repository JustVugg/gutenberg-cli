import test from "node:test";
import assert from "node:assert/strict";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { harToOpenApi, loadHar } from "../src/core/har.js";
import { buildBlueprint } from "../src/core/openapi.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sampleHar = path.join(root, "samples", "simple.har.json");

test("converts HAR captures into OpenAPI", () => {
  const spec = harToOpenApi(loadHar(sampleHar), { name: "Captured CRM" });
  assert.equal(spec.openapi, "3.0.3");
  assert.equal(spec.servers[0].url, "https://api.example.com");
  assert.ok(spec.paths["/v1/customers"]);
  assert.ok(spec.paths["/v1/customers/{id}"]);
  assert.equal(spec.paths["/assets/app_js"], undefined);

  const blueprint = buildBlueprint(spec, sampleHar, "Captured CRM");
  assert.equal(blueprint.operations.length, 3);
  assert.equal(blueprint.operations.find((operation) => operation.id === "postCustomers").risk, "write");
});
