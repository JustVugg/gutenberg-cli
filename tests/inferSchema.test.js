import test from "node:test";
import assert from "node:assert/strict";
import { inferSchema } from "../src/core/har.js";

test("inferSchema: primitives", () => {
  assert.deepEqual(inferSchema(null), { type: "null" });
  assert.deepEqual(inferSchema(true), { type: "boolean" });
  assert.deepEqual(inferSchema(42), { type: "integer" });
  assert.deepEqual(inferSchema(3.14), { type: "number" });
  assert.deepEqual(inferSchema("foo"), { type: "string" });
});

test("inferSchema: string formats", () => {
  assert.deepEqual(inferSchema("2026-05-15T10:00:00Z"), { type: "string", format: "date-time" });
  assert.deepEqual(inferSchema("2026-05-15"), { type: "string", format: "date" });
  assert.deepEqual(inferSchema("https://example.com"), { type: "string", format: "uri" });
  assert.deepEqual(inferSchema("a@b.com"), { type: "string", format: "email" });
});

test("inferSchema: arrays", () => {
  const schema = inferSchema([1, 2, 3]);
  assert.equal(schema.type, "array");
  assert.equal(schema.items.type, "integer");
});

test("inferSchema: object with required keys", () => {
  const schema = inferSchema({ id: 1, name: "x", optional: null });
  assert.equal(schema.type, "object");
  assert.ok(schema.required.includes("id"));
  assert.ok(schema.required.includes("name"));
  assert.ok(!schema.required.includes("optional"));
});

test("inferSchema: merges array element schemas", () => {
  const schema = inferSchema([{ id: 1, name: "a" }, { id: 2, label: "b" }]);
  assert.equal(schema.items.type, "object");
  assert.ok("id" in schema.items.properties);
  assert.ok("name" in schema.items.properties);
  assert.ok("label" in schema.items.properties);
});
