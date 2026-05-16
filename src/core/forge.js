import path from "node:path";
import { fileURLToPath } from "node:url";
import { buildBlueprint } from "./openapi.js";
import { generateProject } from "./render.js";
import { installTool } from "./install.js";
import { materializeSource, planSource, resolvePolicy } from "./source.js";
import { slugify } from "./sanitize.js";
import { verifyProject } from "./verify.js";

const rootDir = path.dirname(path.dirname(path.dirname(fileURLToPath(import.meta.url))));

export async function runPlan(source, options = {}) {
  return planSource(source, options);
}

export async function forgeProject(source, options = {}) {
  const name = options.name || options.display || inferNameFromSource(source);
  const slug = slugify(name);
  const outDir = path.resolve(options.out || path.join("generated", `${slug}-go`));
  const resolved = await materializeSource(source, {
    ...options,
    name,
    outDir
  });
  const blueprint = buildBlueprint(resolved.spec, resolved.sourcePath || resolved.input, name);
  const policy = resolvePolicy(options);
  blueprint.policy = policy;
  const generated = generateProject(blueprint, outDir, {
    force: Boolean(options.force),
    lang: options.lang || "go",
    name,
    displayName: options.display || blueprint.name,
    targets: options.targets,
    defaultHeaders: options["default-header"],
    specPath: resolved.sourcePath || null,
    policy
  });
  const verification = verifyProject(generated.outDir, {
    noTidy: Boolean(options["no-tidy"]),
    skipTests: Boolean(options["skip-tests"])
  });

  let install = null;
  if (options.install) {
    install = installTool(generated.outDir, {
      prefix: options.prefix,
      skipTidy: true
    });
  }

  const command = install
    ? (install.onPath ? `${install.slug} operations` : `${install.binPath} operations`)
    : `cd ${generated.outDir} && ${path.join(rootDir, "scripts", "use-go.sh")} run ./cmd/${generated.manifest.slug} operations`;

  return {
    schemaVersion: "gutenberg.forge.result.v1",
    source: resolved.input,
    kind: resolved.kind,
    normalizedSpec: resolved.sourcePath || null,
    outDir: generated.outDir,
    entrypoint: generated.entrypoint,
    manifest: generated.manifest,
    verification,
    install,
    command
  };
}

function inferNameFromSource(source) {
  const text = String(source || "generated-api");
  if (/^curl\s/i.test(text)) {
    try {
      const url = text.match(/https?:\/\/[^\s'"]+/i)?.[0];
      if (url) return new URL(url).hostname.replace(/^www\./, "");
    } catch {}
    return "curl-api";
  }
  if (/^https?:\/\//i.test(text)) {
    try {
      return new URL(text).hostname.replace(/^www\./, "");
    } catch {
      return "site-api";
    }
  }
  return path.basename(text, path.extname(text)) || "generated-api";
}
