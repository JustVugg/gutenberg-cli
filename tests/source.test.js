import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { curlToRequest, detectSourceKind, planSource, postmanToOpenApi } from "../src/core/source.js";

test("detectSourceKind recognizes OpenAPI, HAR, curl, Postman and GraphQL", () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-source-"));
  const openapi = path.join(tmp, "openapi.json");
  const har = path.join(tmp, "capture.json");
  const postman = path.join(tmp, "postman.json");
  const gql = path.join(tmp, "schema.graphql");
  fs.writeFileSync(openapi, JSON.stringify({ openapi: "3.0.3", info: { title: "A", version: "1" }, paths: {} }));
  fs.writeFileSync(har, JSON.stringify({ log: { entries: [] } }));
  fs.writeFileSync(postman, JSON.stringify({ info: { name: "P", schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json" }, item: [] }));
  fs.writeFileSync(gql, "type Query { ping: String }");

  assert.equal(detectSourceKind(openapi), "openapi");
  assert.equal(detectSourceKind(har), "har");
  assert.equal(detectSourceKind(postman), "postman");
  assert.equal(detectSourceKind(gql), "graphql");
  assert.equal(detectSourceKind("curl https://api.example.com/v1/pets"), "curl");
});

test("curlToRequest parses method, URL, headers and body without leaking header values into operation metadata", () => {
  const request = curlToRequest("curl -X POST https://api.example.com/v1/pets -H 'Authorization: Bearer secret' -H 'X-Api-Key: abc' --data '{\"name\":\"Milo\"}'");
  assert.equal(request.method, "POST");
  assert.equal(request.url, "https://api.example.com/v1/pets");
  assert.equal(request.headers.length, 2);
  assert.equal(request.body, "{\"name\":\"Milo\"}");
});

test("Postman collection converts to OpenAPI with auth schemes and request bodies", () => {
  const spec = postmanToOpenApi({
    info: { name: "Pets", schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json" },
    item: [{
      name: "Create pet",
      request: {
        method: "POST",
        url: "https://api.example.com/v1/pets",
        header: [{ key: "Authorization", value: "Bearer secret" }],
        body: { mode: "raw", raw: "{\"name\":\"Milo\"}" }
      }
    }]
  });
  assert.equal(spec.openapi, "3.0.3");
  assert.equal(spec.servers[0].url, "https://api.example.com");
  assert.ok(spec.paths["/v1/pets"].post.requestBody);
  assert.ok(spec.components.securitySchemes.bearerAuth);
});

test("planSource returns a reviewable blueprint summary", async () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-plan-"));
  const specPath = path.join(tmp, "openapi.json");
  fs.writeFileSync(specPath, JSON.stringify({
    openapi: "3.0.3",
    info: { title: "Plan API", version: "1" },
    servers: [{ url: "https://api.example.com" }],
    paths: { "/items": { get: { operationId: "listItems", responses: { 200: { description: "ok" } } } } }
  }));
  const plan = await planSource(specPath, { name: "plan-api" });
  assert.equal(plan.schemaVersion, "gutenberg.plan.v1");
  assert.equal(plan.operations.length, 1);
  assert.equal(plan.policy.rules.some((rule) => rule.risk === "write" && rule.requiresYes), true);
});
