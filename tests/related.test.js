import test from "node:test";
import assert from "node:assert/strict";
import { linkRelatedOperations } from "../src/core/openapi.js";

test("linkRelatedOperations: detail/list pairing", () => {
  const ops = [
    { id: "listItems", method: "GET", path: "/items" },
    { id: "getItem", method: "GET", path: "/items/{id}" },
    { id: "createItem", method: "POST", path: "/items" }
  ];
  linkRelatedOperations(ops);
  const getItem = ops.find((op) => op.id === "getItem");
  assert.ok(getItem.related.some((r) => r.id === "listItems" && r.role === "list-parent"));
  const listItems = ops.find((op) => op.id === "listItems");
  assert.ok(listItems.related.some((r) => r.id === "getItem" && r.role === "detail-child"));
  assert.ok(listItems.related.some((r) => r.id === "createItem" && r.role === "same-resource"));
});
