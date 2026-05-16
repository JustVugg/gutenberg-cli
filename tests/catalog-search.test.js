import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { searchCatalog } from "../src/core/catalog-search.js";

function setupFakeCatalog() {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "gutenberg-cat-"));
  const toolsDir = path.join(tmp, "library", "tools");
  fs.mkdirSync(path.join(toolsDir, "fake"), { recursive: true });
  fs.writeFileSync(path.join(toolsDir, "fake", "gutenberg.manifest.json"), JSON.stringify({
    slug: "fake",
    name: "Fake API",
    description: "Test fake catalog tool with sports operations",
    tags: ["sports", "fixtures"],
    heroes: [{ alias: "list-fixtures", operationId: "fixturesList", summary: "list sport fixtures", method: "GET", path: "/fixtures" }],
    operations: [
      { id: "fixturesList", method: "GET", path: "/fixtures", summary: "list sport fixtures", tag: "sports", kind: "read-list" },
      { id: "playerById", method: "GET", path: "/players/{id}", summary: "fetch a player", tag: "players", kind: "read-detail" }
    ]
  }));
  return tmp;
}

test("searchCatalog ranks heroes higher than ops", () => {
  const root = setupFakeCatalog();
  const result = searchCatalog(root, "sport fixtures");
  assert.ok(result.results.length > 0);
  assert.equal(result.results[0].kind, "hero");
  assert.equal(result.results[0].alias, "list-fixtures");
});

test("searchCatalog returns empty for irrelevant intent", () => {
  const root = setupFakeCatalog();
  const result = searchCatalog(root, "quantum spaghetti monsters");
  assert.equal(result.results.length, 0);
});
