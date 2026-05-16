import test from "node:test";
import assert from "node:assert/strict";
import { harToOpenApi } from "../src/core/har.js";

function makeEntry(url, headers) {
  return {
    request: { method: "GET", url, headers, queryString: [] },
    response: { status: 200, statusText: "OK", content: { mimeType: "application/json", text: "{}" }, headers: [] }
  };
}

test("HAR sniffer: detects Bearer token", () => {
  const har = { log: { entries: [makeEntry("https://api.example.com/x", [{ name: "Authorization", value: "Bearer abc123" }]) ] } };
  const spec = harToOpenApi(har, { name: "test" });
  assert.ok(spec.components.securitySchemes.bearerAuth);
  assert.equal(spec.components.securitySchemes.bearerAuth.scheme, "bearer");
});

test("HAR sniffer: detects X-API-Key header", () => {
  const har = { log: { entries: [makeEntry("https://api.example.com/x", [{ name: "X-API-Key", value: "abc123" }]) ] } };
  const spec = harToOpenApi(har, { name: "test" });
  const schemes = spec.components.securitySchemes;
  const apiKeyScheme = Object.values(schemes).find((s) => s.type === "apiKey" && s.in === "header");
  assert.ok(apiKeyScheme);
  assert.equal(apiKeyScheme.name, "X-API-Key");
});

test("HAR sniffer: no auth on bare requests", () => {
  const har = { log: { entries: [makeEntry("https://api.example.com/x", []) ] } };
  const spec = harToOpenApi(har, { name: "test" });
  assert.equal(spec.components, undefined);
});
