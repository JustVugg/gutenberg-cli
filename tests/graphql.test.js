import test from "node:test";
import assert from "node:assert/strict";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { buildBlueprint } from "../src/core/openapi.js";
import { graphqlToOpenApi, loadGraphQLSource } from "../src/core/graphql.js";

const root = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const sample = path.join(root, "samples", "github-introspection.json");

test("converts GraphQL introspection into OpenAPI", async () => {
  const input = await loadGraphQLSource(sample);
  const spec = graphqlToOpenApi(input, { name: "GitHub GraphQL", endpoint: "https://api.github.com/graphql" });
  assert.equal(spec.openapi, "3.0.3");
  assert.ok(spec.paths["/graphql/query/repository"]);
  assert.ok(spec.paths["/graphql/mutation/addstar"]);

  const blueprint = buildBlueprint(spec, sample, "GitHub GraphQL");
  assert.equal(blueprint.operations.length, 3);
  assert.equal(blueprint.operations.find((operation) => operation.id === "postMutationAddstar").risk, "write");
});

test("converts simple GraphQL SDL into OpenAPI", () => {
  const spec = graphqlToOpenApi({
    sdl: `
      type Query {
        user(id: ID!): User
        search(q: String): [Result]
      }
      type Mutation {
        createUser(name: String!): User
      }
    `
  }, { name: "SDL API" });
  assert.ok(spec.paths["/graphql/query/user"]);
  assert.ok(spec.paths["/graphql/query/search"]);
  assert.ok(spec.paths["/graphql/mutation/createuser"]);
});
