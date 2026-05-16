import test from "node:test";
import assert from "node:assert/strict";
import { classifyOperation, detectPagination } from "../src/core/openapi.js";

test("classifyOperation: DELETE", () => {
  assert.equal(classifyOperation({ method: "DELETE", path: "/items/{id}", summary: "" }), "delete");
});

test("classifyOperation: POST create", () => {
  assert.equal(classifyOperation({ method: "POST", path: "/items", summary: "Create item" }), "create");
});

test("classifyOperation: POST search", () => {
  assert.equal(classifyOperation({ method: "POST", path: "/items/search", summary: "Search items" }), "search");
});

test("classifyOperation: PUT update with id", () => {
  assert.equal(classifyOperation({ method: "PUT", path: "/items/{id}", summary: "" }), "update");
});

test("classifyOperation: GET list vs detail", () => {
  assert.equal(classifyOperation({ method: "GET", path: "/items", summary: "List items" }), "read-list");
  assert.equal(classifyOperation({ method: "GET", path: "/items/{id}", summary: "Get item" }), "read-detail");
});

test("classifyOperation: GET search", () => {
  assert.equal(classifyOperation({ method: "GET", path: "/search", summary: "Search items" }), "search");
});

test("detectPagination: offset/limit", () => {
  const p = detectPagination({ parameters: [{ name: "offset", in: "query" }, { name: "limit", in: "query" }] });
  assert.equal(p.style, "offset-limit");
});

test("detectPagination: cursor", () => {
  const p = detectPagination({ parameters: [{ name: "cursor", in: "query" }] });
  assert.equal(p.style, "cursor");
});

test("detectPagination: page", () => {
  const p = detectPagination({ parameters: [{ name: "page", in: "query" }, { name: "per_page", in: "query" }] });
  assert.equal(p.style, "page");
  assert.equal(p.perPageParam, "per_page");
});

test("detectPagination: none", () => {
  assert.equal(detectPagination({ parameters: [{ name: "filter", in: "query" }] }), null);
});
