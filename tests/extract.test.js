import test from "node:test";
import assert from "node:assert/strict";
import { stripHtml, validateAgainstSchema } from "../src/core/extract.js";

test("stripHtml removes scripts, styles, and tags", () => {
  const html = `<html><head><style>.x{}</style><script>alert(1)</script></head><body><h1>Hi</h1><p>World &amp; you</p></body></html>`;
  const cleaned = stripHtml(html);
  assert.equal(cleaned, "Hi World & you");
});

test("validateAgainstSchema flags missing required + type errors", () => {
  const schema = { type: "object", required: ["title"], properties: { title: { type: "string" }, price: { type: "number" } } };
  const errors = validateAgainstSchema({ price: "not a number" }, schema);
  assert.ok(errors.some((e) => e.message.includes("missing required property: title")));
  assert.ok(errors.some((e) => e.path === "$.price" && e.message.includes("expected number")));
});

test("validateAgainstSchema accepts well-formed payloads", () => {
  const schema = { type: "object", required: ["items"], properties: { items: { type: "array", items: { type: "object", required: ["name"], properties: { name: { type: "string" }, qty: { type: "integer" } } } } } };
  const errors = validateAgainstSchema({ items: [{ name: "a", qty: 1 }, { name: "b", qty: 2 }] }, schema);
  assert.equal(errors.length, 0);
});
