import test from "node:test";
import assert from "node:assert/strict";
import { buildBlueprint } from "../src/core/openapi.js";

test("x-gutenberg-hero overrides auto-detected hero alias", () => {
  const spec = {
    openapi: "3.0.3",
    info: { title: "Demo" },
    paths: {
      "/items": {
        get: {
          operationId: "listItems",
          summary: "list items",
          "x-gutenberg-hero": { alias: "list", "default-params": { limit: "5" } },
          responses: { "200": { description: "ok" } }
        }
      },
      "/users": {
        get: { operationId: "listUsers", summary: "list users", responses: { "200": { description: "ok" } } }
      }
    }
  };
  const blueprint = buildBlueprint(spec, "test.json", "Demo");
  const heroList = blueprint.heroes.find((h) => h.alias === "list");
  assert.ok(heroList, "explicit hero alias 'list' should be present");
  assert.equal(heroList.operationId, "listItems");
  assert.equal(heroList.explicit, true);
  assert.deepEqual(heroList.defaultParams, { limit: "5" });
});
