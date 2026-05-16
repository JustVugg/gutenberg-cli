import path from "node:path";
import { readJson, readText, writeJson } from "./fs.js";
import { safeOperationId, slugify } from "./sanitize.js";

export const INTROSPECTION_QUERY = `
query GutenbergIntrospection {
  __schema {
    queryType { name }
    mutationType { name }
    types {
      kind
      name
      fields {
        name
        description
        args {
          name
          description
          type { kind name ofType { kind name ofType { kind name } } }
        }
        type { kind name ofType { kind name ofType { kind name } } }
      }
    }
  }
}
`;

export async function loadGraphQLSource(source, options = {}) {
  if (/^https?:\/\//i.test(source)) {
    return introspectEndpoint(source, options);
  }
  const extension = path.extname(source).toLowerCase();
  if (extension === ".json") {
    return readJson(source);
  }
  return { sdl: readText(source) };
}

export async function introspectEndpoint(endpoint, options = {}) {
  const headers = {
    "Content-Type": "application/json",
    Accept: "application/json"
  };
  if (options.token) {
    headers.Authorization = `Bearer ${options.token}`;
  }
  const response = await fetch(endpoint, {
    method: "POST",
    headers,
    body: JSON.stringify({ query: INTROSPECTION_QUERY })
  });
  if (!response.ok) {
    throw new Error(`GraphQL introspection failed: ${response.status} ${response.statusText}`);
  }
  return response.json();
}

export function graphqlToOpenApi(input, options = {}) {
  if (input.sdl) {
    return sdlToOpenApi(input.sdl, options);
  }
  return introspectionToOpenApi(input, options);
}

export function writeGraphQLOpenApi(input, outPath, options = {}) {
  const spec = graphqlToOpenApi(input, options);
  writeJson(outPath, spec);
  return spec;
}

function introspectionToOpenApi(input, options) {
  const schema = input.data?.__schema || input.__schema;
  if (!schema) {
    throw new Error("GraphQL input does not contain __schema introspection data");
  }
  const typeMap = new Map((schema.types || []).map((type) => [type.name, type]));
  const queryType = typeMap.get(schema.queryType?.name);
  const mutationType = schema.mutationType?.name ? typeMap.get(schema.mutationType.name) : null;
  const paths = {};

  for (const field of queryType?.fields || []) {
    addGraphQLOperation(paths, "query", field);
  }
  for (const field of mutationType?.fields || []) {
    addGraphQLOperation(paths, "mutation", field);
  }

  return {
    openapi: "3.0.3",
    info: {
      title: options.name || "GraphQL API",
      version: "0.1.0",
      description: "Generated from GraphQL introspection by Gutenberg."
    },
    servers: [{ url: options.endpoint || "https://graphql.example.com/graphql" }],
    paths
  };
}

function sdlToOpenApi(sdl, options) {
  const paths = {};
  for (const [kind, body] of [
    ["query", extractTypeBody(sdl, "Query")],
    ["mutation", extractTypeBody(sdl, "Mutation")]
  ]) {
    if (!body) continue;
    for (const field of parseSDLFields(body)) {
      addGraphQLOperation(paths, kind, field);
    }
  }

  return {
    openapi: "3.0.3",
    info: {
      title: options.name || "GraphQL API",
      version: "0.1.0",
      description: "Generated from GraphQL SDL by Gutenberg."
    },
    servers: [{ url: options.endpoint || "https://graphql.example.com/graphql" }],
    paths
  };
}

function addGraphQLOperation(paths, kind, field) {
  const apiPath = `/graphql/${kind}/${slugify(field.name)}`;
  const method = "post";
  const operationId = safeOperationId(method, apiPath);
  paths[apiPath] = {
    [method]: {
      operationId,
      tags: [kind],
      summary: `${kind} ${field.name}`,
      description: field.description || "",
      requestBody: {
        required: true,
        content: {
          "application/json": {
            schema: {
              type: "object",
              properties: {
                query: { type: "string" },
                variables: {
                  type: "object",
                  additionalProperties: true,
                  properties: argsToProperties(field.args || [])
                }
              },
              required: ["query"]
            }
          }
        }
      },
      responses: {
        "200": { description: "GraphQL response" },
        default: { description: "GraphQL error" }
      },
      "x-graphql": {
        kind,
        field: field.name,
        args: (field.args || []).map((arg) => ({
          name: arg.name,
          type: typeName(arg.type || arg.typeName || "String")
        }))
      }
    }
  };
}

function argsToProperties(args) {
  const properties = {};
  for (const arg of args) {
    properties[arg.name] = {
      type: scalarToJsonType(typeName(arg.type || arg.typeName || "String")),
      description: arg.description || ""
    };
  }
  return properties;
}

function typeName(type) {
  if (typeof type === "string") return type.replace(/[!\[\]]/g, "");
  if (!type) return "String";
  if (type.name) return type.name;
  return typeName(type.ofType);
}

function scalarToJsonType(name) {
  if (["Int", "Float"].includes(name)) return "number";
  if (name === "Boolean") return "boolean";
  return "string";
}

function extractTypeBody(sdl, typeNameValue) {
  const match = sdl.match(new RegExp(`type\\s+${typeNameValue}\\s*\\{([\\s\\S]*?)\\}`, "m"));
  return match ? match[1] : "";
}

function parseSDLFields(body) {
  return body
    .split("\n")
    .map((line) => line.replace(/#.*/, "").trim())
    .filter(Boolean)
    .map((line) => {
      const match = line.match(/^([A-Za-z_][A-Za-z0-9_]*)(?:\(([^)]*)\))?\s*:/);
      if (!match) return null;
      return {
        name: match[1],
        args: parseSDLArgs(match[2] || "")
      };
    })
    .filter(Boolean);
}

function parseSDLArgs(argsText) {
  return argsText
    .split(",")
    .map((part) => part.trim())
    .filter(Boolean)
    .map((part) => {
      const match = part.match(/^([A-Za-z_][A-Za-z0-9_]*)\s*:\s*([^=]+)(?:=.*)?$/);
      return match ? { name: match[1], typeName: match[2].trim() } : null;
    })
    .filter(Boolean);
}
