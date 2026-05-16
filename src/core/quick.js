import fs from "node:fs";
import path from "node:path";
import os from "node:os";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import { tryUrl } from "./try.js";
import { seedHarFromUrls } from "./seed.js";
import { harToOpenApi, loadHar } from "./har.js";
import { buildBlueprint, loadOpenApi } from "./openapi.js";
import { generateProject } from "./render.js";
import { installTool } from "./install.js";
import { searchCatalog } from "./catalog-search.js";
import { runIntent } from "./run-intent.js";
import { writeJson } from "./fs.js";

const gutenbergRoot = path.dirname(path.dirname(path.dirname(fileURLToPath(import.meta.url))));

// `quick` is the magic front-door. Pass either a URL (it builds you a tool from
// scratch — try → seed → import → generate → verify → install) or a natural
// language intent (it routes to an existing tool via the catalog).
export async function quick(input, options = {}) {
  if (!input) {
    const error = new Error("Missing input. Pass a URL or a natural language intent.");
    error.exitCode = 2;
    throw error;
  }
  const isUrl = /^https?:\/\//i.test(input);
  if (isUrl) {
    return await quickFromUrl(input, options);
  }
  return await quickFromIntent(input, options);
}

async function quickFromUrl(url, options) {
  const log = (msg) => process.stderr.write(`gutenberg quick: ${msg}\n`);
  log(`probing ${url}`);
  const probe = await tryUrl(url);
  log(`verdict = ${probe.verdict} (${probe.confidence})`);

  if (probe.verdict !== "json-endpoint" && probe.verdict !== "openapi-published") {
    return {
      mode: "url",
      url,
      probe,
      message: `Cannot auto-build for verdict "${probe.verdict}". See report.nextSteps for the next manual command.`,
      autoBuilt: false
    };
  }

  const slug = options.name || guessSlugFromUrl(url);
  const projectDir = path.resolve(options.outDir || path.join(gutenbergRoot, "library", "tools", slug));

  let specPath;
  if (probe.verdict === "openapi-published") {
    const discoveryUrl = probe.steps.find((s) => s.name === "discover")?.url || url;
    specPath = path.join(os.tmpdir(), `${slug}-${Date.now()}.openapi.json`);
    log(`downloading spec from ${discoveryUrl}`);
    const response = await fetch(discoveryUrl);
    if (!response.ok) {
      const error = new Error(`failed to download spec: HTTP ${response.status}`);
      error.exitCode = 1;
      throw error;
    }
    const body = await response.text();
    fs.writeFileSync(specPath, body, "utf8");

    // Patch in baseUrl from the original URL if the spec didn't declare any servers.
    try {
      const spec = JSON.parse(body);
      const inferredOrigin = new URL(url).origin;
      const noServers = !Array.isArray(spec.servers) || spec.servers.length === 0;
      const noHost = !spec.host;
      if (noServers && noHost) {
        spec.servers = [{ url: inferredOrigin }];
        fs.writeFileSync(specPath, JSON.stringify(spec, null, 2), "utf8");
        log(`spec had no servers; inferred ${inferredOrigin}`);
      }
    } catch {}
  } else {
    const harPath = path.join(os.tmpdir(), `${slug}-${Date.now()}.har.json`);
    const extras = options.relatedUrls || [];
    log(`seeding HAR (${1 + extras.length} url${extras.length === 0 ? "" : "s"})`);
    await seedHarFromUrls([url, ...extras], { out: harPath });
    const har = loadHar(harPath);
    const spec = harToOpenApi(har, { name: slug });
    specPath = path.join(os.tmpdir(), `${slug}-${Date.now()}.openapi.json`);
    writeJson(specPath, spec);
  }

  log(`generating tool at ${projectDir}`);
  const blueprint = buildBlueprint(loadOpenApi(specPath), specPath, slug);
  const result = generateProject(blueprint, projectDir, {
    force: true,
    name: slug,
    targets: options.targets,
    specPath,
    defaultHeaders: options.defaultHeaders
  });

  let verified = null;
  if (!options.noVerify) {
    log(`verifying ${slug}`);
    const useGo = path.join(gutenbergRoot, "scripts", "use-go.sh");
    const verifyCmd = spawnSync(process.execPath, [
      path.join(gutenbergRoot, "bin", "gutenberg.js"),
      "verify",
      projectDir
    ], { encoding: "utf8", env: { ...process.env } });
    verified = verifyCmd.status === 0;
    log(`verify ${verified ? "ok" : "failed (exit " + verifyCmd.status + ")"}`);
  }

  let installed = null;
  if (verified && !options.noInstall) {
    try {
      log(`installing ${slug}`);
      installed = installTool(projectDir, {});
    } catch (error) {
      log(`install skipped: ${error.message}`);
    }
  }

  return {
    mode: "url",
    url,
    slug,
    projectDir,
    scorecard: result.scorecard ? { score: result.scorecard.score, grade: result.scorecard.grade } : null,
    verified,
    installed: installed ? { binPath: installed.binPath, onPath: installed.onPath } : null,
    heroes: (result.manifest?.heroes || []).slice(0, 5).map((h) => h.alias),
    autoBuilt: true
  };
}

async function quickFromIntent(intent, options) {
  const log = (msg) => process.stderr.write(`gutenberg quick: ${msg}\n`);
  log(`searching catalog for "${intent}"`);
  const catalog = searchCatalog(gutenbergRoot, intent, { limit: 5 });
  if (catalog.results.length === 0) {
    return {
      mode: "intent",
      intent,
      message: "No matches in catalog. Try `gutenberg <provider URL>` to build a new tool, or `gutenberg search <intent>`.",
      candidates: []
    };
  }
  log(`top match: ${catalog.results[0].command}`);
  if (options.dryRun) {
    return { mode: "intent", intent, dryRun: true, candidates: catalog.results };
  }
  const run = await runIntent(gutenbergRoot, intent, {
    dryRun: false,
    noLlm: options.noLlm,
    provider: options.provider,
    model: options.model
  });
  return { mode: "intent", intent, plan: run.plan, candidates: catalog.results, result: run.result };
}

function guessSlugFromUrl(url) {
  try {
    const parsed = new URL(url);
    const host = parsed.hostname.replace(/^www\./, "").replace(/\.(com|org|net|io|dev|app|it|ai|co|us|uk|de|fr|es)$/i, "").replace(/[^a-z0-9]+/gi, "-").toLowerCase();
    return host || "tool";
  } catch {
    return "tool";
  }
}
