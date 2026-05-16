import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";
import { buildBlueprint, loadOpenApi } from "./core/openapi.js";
import { generateProject } from "./core/render.js";
import { readJson, writeJson, writeText } from "./core/fs.js";
import { scoreProject } from "./core/scorecard.js";
import { slugify } from "./core/sanitize.js";
import { harToOpenApi, loadHar } from "./core/har.js";
import { discoverOpenApi } from "./core/discover.js";
import { graphqlToOpenApi, loadGraphQLSource } from "./core/graphql.js";
import { recordBrowser, loginInteractive } from "./core/recorder.js";
import { generateAggregatorProject, loadAggregatorRecipe } from "./core/render-aggregator.js";
import { seedHarFromUrls } from "./core/seed.js";
import { diffSpecs } from "./core/diff.js";
import { upgradeProject } from "./core/upgrade.js";
import { buildCatalog, buildRegistryFromCatalog, writeCatalogDataFile, writeRegistryFile } from "./core/catalog.js";
import { compareTools } from "./core/compare.js";
import { installStarterPack, installTool, STARTER_PACK } from "./core/install.js";
import { searchCatalog } from "./core/catalog-search.js";
import { extractFromUrl, loadSchemaFromFlag, clearExtractCache } from "./core/extract.js";
import { tryUrl, formatTryReport } from "./core/try.js";
import { quick } from "./core/quick.js";
import { scrapeMarkdown } from "./core/scrape.js";
import { runIntent } from "./core/run-intent.js";
import { watchSpec } from "./core/watch.js";
import { runCreate, runRecipes, runSports, runTravel } from "./core/hero.js";
import { forgeProject, runPlan } from "./core/forge.js";
import { verifyProject } from "./core/verify.js";
import { resolvePolicy } from "./core/source.js";
import { publishAll, publishTool } from "./core/publish.js";

const rootDir = path.dirname(path.dirname(fileURLToPath(import.meta.url)));

export async function main(argv) {
  const command = argv[0] || "help";
  const parsed = parseArgs(argv.slice(1));

  if (["help", "--help", "-h"].includes(command)) {
    printHelp();
    return;
  }
  if (command === "doctor") {
    await doctor();
    return;
  }
  if (command === "sports") {
    await runSports(parsed.positionals, parsed.options);
    return;
  }
  if (command === "travel") {
    await runTravel(parsed.positionals, parsed.options);
    return;
  }
  if (command === "recipes") {
    await runRecipes(parsed.positionals, parsed.options, rootDir);
    return;
  }
  if (command === "create") {
    runCreate(parsed.positionals, parsed.options, rootDir);
    return;
  }
  if (command === "plan") {
    const source = parsed.positionals.join(" ");
    if (!source) throw usage("Missing source. Example: gutenberg plan samples/petstore-openapi.json --name petstore");
    const plan = await runPlan(source, parsed.options);
    if (parsed.options.json) {
      console.log(JSON.stringify(plan, null, 2));
      return;
    }
    printPlan(plan);
    return;
  }
  if (command === "forge") {
    const source = parsed.positionals.join(" ");
    if (!source) throw usage("Missing source. Example: gutenberg forge samples/petstore-openapi.json --name petstore --install");
    const result = await forgeProject(source, parsed.options);
    if (parsed.options.json) {
      console.log(JSON.stringify(result, null, 2));
      return;
    }
    printForgeResult(result);
    return;
  }
  if (command === "analyze") {
    const specPath = parsed.positionals[0];
    if (!specPath) throw usage("Missing spec path. Example: gutenberg analyze openapi.json");
    const blueprint = buildBlueprint(loadOpenApi(specPath), specPath, parsed.options.name);
    if (parsed.options.out) writeJson(path.resolve(parsed.options.out), blueprint);
    printAnalysis(blueprint, Boolean(parsed.options.json));
    return;
  }
  if (command === "import-har") {
    const harPath = parsed.positionals[0];
    if (!harPath) throw usage("Missing HAR path. Example: gutenberg import-har capture.har --out openapi.json --name app");
    const spec = harToOpenApi(loadHar(harPath), { name: parsed.options.name, origin: parsed.options.origin });
    if (parsed.options.out) {
      writeJson(path.resolve(parsed.options.out), spec);
      console.log(`Wrote OpenAPI spec: ${path.resolve(parsed.options.out)}`);
      return;
    }
    console.log(JSON.stringify(spec, null, 2));
    return;
  }
  if (command === "import-graphql") {
    const source = parsed.positionals[0];
    if (!source) throw usage("Missing GraphQL source. Example: gutenberg import-graphql schema.json --out openapi.json --name github");
    const input = await loadGraphQLSource(source, {
      token: parsed.options.token || process.env.GUTENBERG_GRAPHQL_TOKEN
    });
    const spec = graphqlToOpenApi(input, {
      name: parsed.options.name,
      endpoint: /^https?:\/\//i.test(source) ? source : parsed.options.endpoint
    });
    if (parsed.options.out) {
      writeJson(path.resolve(parsed.options.out), spec);
      console.log(`Wrote OpenAPI spec: ${path.resolve(parsed.options.out)}`);
      return;
    }
    console.log(JSON.stringify(spec, null, 2));
    return;
  }
  if (command === "seed-har") {
    const urls = parsed.positionals;
    if (urls.length === 0) throw usage("Missing URL(s). Example: gutenberg seed-har https://api.example.com/v1/foo --out capture.har.json");
    const out = parsed.options.out ? path.resolve(parsed.options.out) : path.resolve("capture.har.json");
    const headers = headerOptions(parsed.options.header);
    const har = await seedHarFromUrls(urls, {
      out,
      method: parsed.options.method,
      headers,
      body: parsed.options.data,
      title: parsed.options.title,
      userAgent: parsed.options["user-agent"],
      goldenDir: parsed.options["save-goldens"]
    });
    console.log(`Wrote HAR with ${har.log.entries.length} entry/entries: ${out}`);
    const showWarnings = !parsed.options["allow-error"];
    if (showWarnings && Array.isArray(har.warnings) && har.warnings.length > 0) {
      const httpErrors = har.warnings.filter((w) => w.kind === "http-error");
      const challenges = har.warnings.filter((w) => w.kind === "anti-bot");
      if (httpErrors.length > 0) {
        console.warn(`WARN: ${httpErrors.length} response(s) returned HTTP errors:`);
        for (const w of httpErrors) console.warn(`  ${w.status} ${w.statusText} — ${w.url}`);
      }
      if (challenges.length > 0) {
        console.warn(`WARN: anti-bot challenge detected on ${challenges.length} URL(s) (Cloudflare/JS verification).`);
        console.warn(`      The captured HAR likely contains the challenge page, not real content.`);
        console.warn(`      Try: gutenberg record <url> --out capture.har.json  (uses Playwright)`);
        console.warn(`      Or:  gutenberg record <url> --backend browserbase --key $BB_KEY`);
      }
      if (httpErrors.length > 0 && !challenges.length) {
        console.warn(`      (suppress with --allow-error)`);
      }
    }
    return;
  }
  if (command === "record" && parsed.options.backend === "browserbase") {
    const { recordViaBrowserbase } = await import("./core/recorder-browserbase.js");
    const targetUrl = parsed.positionals[0];
    if (!targetUrl) throw usage("Missing URL.");
    const out = parsed.options.out ? path.resolve(parsed.options.out) : path.resolve("capture.har.json");
    const har = await recordViaBrowserbase(targetUrl, {
      out,
      apiKey: parsed.options.key,
      projectId: parsed.options["project-id"],
      wait: parsed.options.wait,
      timeout: parsed.options.timeout
    });
    console.log(`Recorded ${har.log.entries.length} request(s) via Browserbase: ${out}`);
    return;
  }
  if (command === "record") {
    const targetUrl = parsed.positionals[0];
    if (!targetUrl) throw usage("Missing URL. Example: gutenberg record https://example.com --out capture.har.json");
    const out = parsed.options.out ? path.resolve(parsed.options.out) : path.resolve("capture.har.json");
    const har = await recordBrowser(targetUrl, {
      out,
      headless: !parsed.options.headed,
      wait: parsed.options.wait,
      timeout: parsed.options.timeout,
      storageState: parsed.options["storage-state"],
      saveStorageState: parsed.options["save-storage-state"]
    });
    console.log(`Recorded ${har.log.entries.length} request(s): ${out}`);
    if (har.log.comment) {
      console.warn(har.log.comment);
    }
    return;
  }
  if (command === "login") {
    const targetUrl = parsed.positionals[0];
    if (!targetUrl) throw usage("Missing URL. Example: gutenberg login https://example.com --out state.json");
    if (!parsed.options.out) throw usage("Missing --out. Example: gutenberg login https://example.com --out state.json");
    const result = await loginInteractive(targetUrl, {
      out: parsed.options.out,
      storageState: parsed.options["storage-state"],
      timeout: parsed.options.timeout,
      waitForUrl: parsed.options["wait-for-url"]
    });
    console.log(`Saved storage-state: ${result.out}`);
    return;
  }
  if (command === "discover") {
    const targetUrl = parsed.positionals[0];
    if (!targetUrl) throw usage("Missing URL. Example: gutenberg discover https://example.com --out openapi.json");
    const result = await discoverOpenApi(targetUrl, { out: parsed.options.out ? path.resolve(parsed.options.out) : null });
    if (parsed.options.json) {
      console.log(JSON.stringify(result, null, 2));
      return;
    }
    if (result.found) {
      console.log(`Found OpenAPI spec: ${result.url}`);
      if (parsed.options.out) console.log(`Wrote: ${path.resolve(parsed.options.out)}`);
    } else {
      console.log("No OpenAPI spec found.");
      console.log("Attempts:");
      for (const attempt of result.attempts) console.log(`  - ${attempt}`);
    }
    return;
  }
  if (command === "generate" && parsed.options.kind === "aggregator") {
    const recipePath = parsed.positionals[0];
    if (!recipePath) throw usage("Missing recipe path. Example: gutenberg generate flights.recipe.json --kind aggregator --out generated/flights");
    const recipe = loadAggregatorRecipe(path.resolve(recipePath), fs);
    const outDir = path.resolve(parsed.options.out || path.join("generated", slugify(recipe.name || "aggregator")));
    const result = generateAggregatorProject(recipe, outDir, {
      force: Boolean(parsed.options.force),
      name: parsed.options.name,
      displayName: parsed.options.display,
      baseDir: path.dirname(path.resolve(recipePath))
    });
    console.log(`Generated aggregator: ${result.manifest.name}`);
    console.log(`Directory: ${result.outDir}`);
    console.log(`Sources: ${result.manifest.sources.map((source) => source.slug).join(", ")}`);
    console.log(`Next: cd ${result.outDir} && ${path.join(rootDir, "scripts", "use-go.sh")} run ./cmd/${result.manifest.slug} sources`);
    return;
  }
  if (command === "generate") {
    const specPath = parsed.positionals[0];
    if (!specPath) throw usage("Missing spec path. Example: gutenberg generate openapi.json --out tools/acme");
    const spec = loadOpenApi(specPath);
    const blueprint = buildBlueprint(spec, specPath, parsed.options.name);
    const policy = resolvePolicy(parsed.options);
    blueprint.policy = policy;
    const outDir = path.resolve(parsed.options.out || path.join("generated", blueprint.slug));
    const result = generateProject(blueprint, outDir, {
      force: Boolean(parsed.options.force),
      lang: parsed.options.lang || "go",
      name: parsed.options.name,
      displayName: parsed.options.display || blueprint.name,
      targets: parsed.options.targets,
      defaultHeaders: parsed.options["default-header"],
      specPath: path.resolve(specPath),
      policy
    });
    if (parsed.options.strict) {
      const minScore = parsed.options["min-score"] !== undefined ? Number(parsed.options["min-score"]) : 0.7;
      const ratio = result.scorecard?.ratio ?? 0;
      if (!Number.isFinite(minScore) || minScore < 0) {
        const error = new Error(`--min-score must be a non-negative number (got ${parsed.options["min-score"]})`);
        error.exitCode = 2;
        throw error;
      }
      if (ratio < minScore) {
        const error = new Error(
          `Strict mode failed: scorecard ratio ${ratio.toFixed(2)} < min ${minScore.toFixed(2)} (score ${result.scorecard?.score}/${result.scorecard?.maxScore}, grade ${result.scorecard?.grade}).`
        );
        error.exitCode = 3;
        throw error;
      }
    }
    console.log(`Generated ${result.manifest.name}`);
    console.log(`Directory: ${result.outDir}`);
    console.log(`Entrypoint: ${result.entrypoint}`);
    if ((parsed.options.lang || "go") === "node") {
      console.log(`Next: node ${result.entrypoint} operations`);
    } else {
      console.log(`Next: cd ${result.outDir} && ${path.join(rootDir, "scripts", "use-go.sh")} run ./cmd/${result.manifest.slug} operations`);
    }
    if (result.skill) {
      console.log(`Skill:    ${result.skill.file}`);
    }
    if (result.drift?.drifted) {
      console.warn(`\nWARN: gutenberg.lock drift detected (${result.drift.drifts.length} change(s)):`);
      for (const drift of result.drift.drifts) {
        console.warn(`  - ${drift.kind}: locked=${drift.locked} current=${drift.current}${drift.path ? ` (${drift.path})` : ""}`);
      }
      console.warn("Review changes with `gutenberg diff` (when available) before shipping.");
    }
    return;
  }
  if (command === "init") {
    const name = parsed.positionals[0];
    if (!name) throw usage("Missing project name. Example: gutenberg init acme-crm --spec openapi.json");
    const outDir = path.resolve(parsed.options.out || slugify(name));
    fs.mkdirSync(outDir, { recursive: true });
    const config = {
      schemaVersion: "gutenberg.project.v1",
      name,
      slug: slugify(name),
      spec: parsed.options.spec || null,
      goals: [
        "Generate one shared API core",
        "Expose CLI and MCP from the same core",
        "Classify risky operations",
        "Cache read data locally",
        "Publish to library registry and web catalog"
      ],
      createdAt: new Date().toISOString()
    };
    writeJson(path.join(outDir, "gutenberg.project.json"), config);
    writeText(path.join(outDir, "RESEARCH.md"), researchTemplate(name));
    if (parsed.options.spec) {
      const spec = loadOpenApi(parsed.options.spec);
      const blueprint = buildBlueprint(spec, parsed.options.spec, name);
      generateProject(blueprint, path.join(outDir, "tool"), { force: Boolean(parsed.options.force), name, lang: parsed.options.lang || "go" });
    }
    console.log(`Initialized ${name} at ${outDir}`);
    return;
  }
  if (command === "upgrade") {
    const dir = parsed.positionals[0];
    if (!dir) throw usage("Missing project dir. Example: gutenberg upgrade library/tools/github");
    const result = upgradeProject(dir, { spec: parsed.options.spec, targets: parsed.options.targets, noTidy: Boolean(parsed.options["no-tidy"]) });
    console.log(`Upgraded ${result.manifest.name}`);
    console.log(`Directory: ${result.outDir}`);
    console.log(`Preserved blocks: ${result.preservedBlocks}, restored: ${result.restoredBlocks.length}, orphaned: ${result.orphanedBlocks.length}`);
    if (result.tidy) {
      console.log(`go mod tidy: ${result.tidy.status === 0 ? "ok" : `FAILED exit=${result.tidy.status}`}`);
      if (result.tidy.status !== 0 && result.tidy.stderr) console.warn(result.tidy.stderr);
    }
    for (const block of result.orphanedBlocks) {
      console.warn(`  WARN orphaned keep block: ${block.relative}:${block.blockId} (${block.reason})`);
    }
    return;
  }
  if (command === "install") {
    const dir = parsed.positionals[0];
    if (!dir) throw usage("Missing project dir. Example: gutenberg install library/tools/github");
    if (dir === "starter-pack") {
      const result = installStarterPack(rootDir, {
        prefix: parsed.options.prefix,
        skipTidy: Boolean(parsed.options["no-tidy"]),
        dryRun: Boolean(parsed.options["dry-run"])
      });
      if (parsed.options.json) {
        if (!result.ready) process.exitCode = 1;
        console.log(JSON.stringify(result, null, 2));
        return;
      }
      console.log(`Starter pack: ${STARTER_PACK.join(", ")}`);
      for (const item of result.installed) {
        if (item.skipped) console.log(`DRY ${item.slug} -> ${item.toolDir}`);
        else console.log(`Installed ${item.slug} -> ${item.binPath}`);
      }
      return;
    }
    const result = installTool(dir, { prefix: parsed.options.prefix, skipTidy: Boolean(parsed.options["no-tidy"]) });
    console.log(`Installed ${result.slug} -> ${result.binPath}`);
    if (!result.onPath) {
      console.warn(`WARN: ${result.binDir} is not on PATH. Add it with:`);
      console.warn(`  export PATH="${result.binDir}:$PATH"`);
    } else {
      console.log(`Try: ${result.slug} operations`);
    }
    return;
  }
  if (command === "scrape") {
    const url = parsed.positionals[0];
    if (!url) throw usage("Missing URL. Example: gutenberg scrape https://example.com");
    const markdown = await scrapeMarkdown(url, { structured: Boolean(parsed.options.structured) });
    if (parsed.options.out) {
      fs.writeFileSync(path.resolve(parsed.options.out), markdown, "utf8");
      console.log(`Wrote markdown: ${path.resolve(parsed.options.out)} (${markdown.length} chars)`);
    } else {
      console.log(markdown);
    }
    return;
  }
  if (command === "extract") {
    if (parsed.options["clear-cache"]) {
      const stats = clearExtractCache();
      console.log(`Cleared ${stats.cleared} cached entries${stats.dir ? ` from ${stats.dir}` : ""}.`);
      return;
    }
    const url = parsed.positionals[0];
    if (!url) throw usage("Missing URL. Example: gutenberg extract https://example.com -p 'extract title and price' -s schema.json");
    const schema = loadSchemaFromFlag(parsed.options.schema);
    const result = await extractFromUrl(url, {
      prompt: parsed.options.prompt,
      schema,
      provider: parsed.options.provider,
      model: parsed.options.model,
      cache: parsed.options.cache
    });
    if (result.errors.length > 0) {
      console.warn(`Schema validation warnings (${result.errors.length}):`);
      for (const err of result.errors) console.warn(`  ${err.path}: ${err.message}`);
    }
    console.log(JSON.stringify(result.parsed, null, 2));
    return;
  }
  if (command === "watch") {
    const [toolDir, specUrl] = parsed.positionals;
    if (!toolDir || !specUrl) throw usage("Usage: gutenberg watch <tool-dir> <spec-url> [--regenerate]");
    const result = await watchSpec(toolDir, specUrl, { regenerate: Boolean(parsed.options.regenerate), noTidy: Boolean(parsed.options["no-tidy"]) });
    console.log(JSON.stringify(result, null, 2));
    return;
  }
  if (command === "run") {
    const intent = parsed.positionals.join(" ");
    const result = await runIntent(rootDir, intent, {
      dryRun: Boolean(parsed.options["dry-run"]),
      noLlm: Boolean(parsed.options["no-llm"]),
      provider: parsed.options.provider,
      model: parsed.options.model
    });
    if (parsed.options.json || result.result === undefined) {
      console.log(JSON.stringify(result, null, 2));
      return;
    }
    if (result.plan) console.log(`> ${result.plan.command}`);
    if (result.result.stdout) process.stdout.write(result.result.stdout);
    if (result.result.stderr) process.stderr.write(result.result.stderr);
    if (result.result.status) process.exitCode = result.result.status;
    return;
  }
  if (command === "quick" || command === "new") {
    const input = parsed.positionals.join(" ").trim();
    if (!input) throw usage("Missing input. Example: gutenberg quick https://api.example.com  or  gutenberg quick 'top hacker news stories'");
    const result = await quick(input, {
      noLlm: Boolean(parsed.options["no-llm"]),
      noInstall: Boolean(parsed.options["no-install"]),
      noVerify: Boolean(parsed.options["no-verify"]),
      dryRun: Boolean(parsed.options["dry-run"]),
      outDir: parsed.options.out,
      name: parsed.options.name,
      relatedUrls: Array.isArray(parsed.options.also) ? parsed.options.also : (parsed.options.also ? [parsed.options.also] : []),
      defaultHeaders: parsed.options["default-header"],
      provider: parsed.options.provider,
      model: parsed.options.model
    });
    console.log(JSON.stringify(result, null, 2));
    return;
  }
  if (command === "try") {
    const url = parsed.positionals[0];
    if (!url) throw usage("Missing URL. Example: gutenberg try https://example.com");
    const report = await tryUrl(url);
    if (parsed.options.json) {
      console.log(JSON.stringify(report, null, 2));
      return;
    }
    console.log(formatTryReport(report));
    return;
  }
  if (command === "search") {
    const intent = parsed.positionals.join(" ");
    if (!intent) throw usage("Missing intent. Example: gutenberg search 'serie a fixtures'");
    const result = searchCatalog(rootDir, intent, { limit: Number(parsed.options.limit) || 5 });
    if (parsed.options.json) {
      console.log(JSON.stringify(result, null, 2));
      return;
    }
    if (result.results.length === 0) {
      console.log("No matches in the catalog.");
      return;
    }
    console.log(`Top matches for: "${intent}"`);
    for (const hit of result.results) {
      console.log(`  ${hit.command}`);
      console.log(`    score=${hit.score} kind=${hit.kind} ${hit.summary || ""}`);
    }
    return;
  }
  if (command === "compare") {
    const tools = parsed.positionals;
    if (tools.length < 2) throw usage("Need 2+ tool dirs. Example: gutenberg compare library/tools/github library/tools/sentry --op meta/root");
    const params = {};
    const paramOpts = parsed.options.param ? (Array.isArray(parsed.options.param) ? parsed.options.param : [parsed.options.param]) : [];
    for (const item of paramOpts) {
      const eq = String(item).indexOf("=");
      if (eq !== -1) params[item.slice(0, eq)] = item.slice(eq + 1);
    }
    const result = await compareTools(tools, {
      operation: parsed.options.op || parsed.options.operation,
      params,
      compact: Boolean(parsed.options.compact),
      token: parsed.options.token
    });
    console.log(JSON.stringify(result, null, 2));
    return;
  }
  if (command === "diff") {
    const [oldSpec, newSpec] = parsed.positionals;
    if (!oldSpec || !newSpec) throw usage("Missing spec paths. Example: gutenberg diff old.openapi.json new.openapi.json");
    const result = diffSpecs(path.resolve(oldSpec), path.resolve(newSpec));
    if (parsed.options.json) {
      console.log(JSON.stringify(result, null, 2));
      return;
    }
    console.log(`Diff: ${oldSpec} -> ${newSpec}`);
    console.log(`Operations: ${result.old.operations} -> ${result.new.operations}`);
    console.log(`Added: ${result.counts.added}, Removed: ${result.counts.removed}, Changed: ${result.counts.changed}`);
    for (const op of result.added) console.log(`  + ${op.method} ${op.path} (${op.id})`);
    for (const op of result.removed) console.log(`  - ${op.method} ${op.path} (${op.id})`);
    for (const op of result.changed) {
      console.log(`  ~ ${op.id}`);
      for (const change of op.changes) console.log(`      ${change.kind}: ${change.from ?? change.name} -> ${change.to ?? ""}`);
    }
    return;
  }
  if (command === "scorecard") {
    const projectDir = path.resolve(parsed.positionals[0] || ".");
    if (parsed.options.verify) {
      verifyProject(projectDir, { noTidy: Boolean(parsed.options["no-tidy"]) });
    }
    const score = scoreProject(projectDir);
    if (parsed.options.json) {
      console.log(JSON.stringify(score, null, 2));
      return;
    }
    printScore(score);
    return;
  }
  if (command === "verify") {
    const projectDir = path.resolve(parsed.positionals[0] || ".");
    const result = verifyProject(projectDir, { noTidy: Boolean(parsed.options["no-tidy"]), skipTests: Boolean(parsed.options["skip-tests"]) });
    if (parsed.options.json) {
      console.log(JSON.stringify(result, null, 2));
      return;
    }
    printVerification(result);
    return;
  }
  if (command === "publish") {
    const target = parsed.positionals[0];
    if (parsed.options.all || target === "--all") {
      const result = publishAll(rootDir, {
        fixAssets: Boolean(parsed.options["fix-assets"]),
        forceAssets: Boolean(parsed.options.force),
        syncRegistry: !parsed.options["no-registry"],
        registryPath: path.resolve(parsed.options.registry || path.join(rootDir, "library", "registry.json"))
      });
      if (parsed.options.json) {
        console.log(JSON.stringify(result, null, 2));
        return;
      }
      printPublishAll(result);
      return;
    }
    if (!target) throw usage("Missing tool dir. Example: gutenberg publish library/tools/github --fix-assets");
    const result = publishTool(target, {
      fixAssets: Boolean(parsed.options["fix-assets"]),
      forceAssets: Boolean(parsed.options.force)
    });
    if (parsed.options.json) {
      if (!result.ready) process.exitCode = 1;
      console.log(JSON.stringify(result, null, 2));
      return;
    }
    printPublishResult(result);
    return;
  }
  if (command === "registry") {
    registry(parsed.positionals, parsed.options);
    return;
  }
  if (command === "site") {
    const action = parsed.positionals[0] || "open";
    if (action === "build") {
      const catalog = buildCatalog(rootDir);
      const outFile = path.join(rootDir, "web", "data.js");
      writeCatalogDataFile(catalog, outFile);
      console.log(`Built catalog: ${catalog.tools.length} tool(s) -> ${outFile}`);
      for (const tool of catalog.tools) {
        console.log(`  ${tool.slug.padEnd(20)} ${tool.kind.padEnd(12)} score=${tool.score ?? "n/a"} ops=${tool.operations}`);
      }
      return;
    }
    console.log(`Open this file in your browser: ${path.join(rootDir, "web", "index.html")}`);
    return;
  }

  throw usage(`Unknown command: ${command}`);
}

const SHORT_FLAGS = { p: "prompt", s: "schema", H: "header", o: "out" };

function parseArgs(args) {
  const options = {};
  const positionals = [];
  for (let index = 0; index < args.length; index += 1) {
    const item = args[index];
    const isLong = item.startsWith("--");
    const isShort = !isLong && item.length === 2 && item.startsWith("-") && /^-[a-zA-Z]$/.test(item);
    if (!isLong && !isShort) {
      positionals.push(item);
      continue;
    }
    let key;
    let value;
    if (isShort) {
      key = SHORT_FLAGS[item[1]] || item.slice(1);
      value = args[index + 1];
      if (value === undefined || value.startsWith("--") || /^-[a-zA-Z]$/.test(value)) {
        options[key] = true;
        continue;
      }
      index += 1;
    } else {
      const equals = item.indexOf("=");
      key = item.slice(2, equals === -1 ? undefined : equals);
      value = equals === -1 ? args[index + 1] : item.slice(equals + 1);
      if (equals === -1 && (value === undefined || value.startsWith("--"))) {
        options[key] = true;
        continue;
      }
      if (equals === -1) index += 1;
    }
    if (options[key] === undefined) {
      options[key] = value;
    } else if (Array.isArray(options[key])) {
      options[key].push(value);
    } else {
      options[key] = [options[key], value];
    }
  }
  return { options, positionals };
}

function printHelp() {
  console.log(`Gutenberg

Build verified tool surfaces for AI agents: CLI, MCP, skills, stores, scorecards, and a web catalog from APIs.

Commands:
  doctor                                  Check local runtime
  sports nba|serie-a today [--json]       Agent-compact live sports schedule
  travel rom-par june [--web|--json]      Agent-compact travel search plan/results
  recipes list|show|run|scaffold          Curated adapters and creation recipes
  create <name> --from <source>           Scaffold a new AI-agent-ready tool
  plan <source> [--json]                  Review detected source, auth, risk, heroes, and policy
  forge <source> --name NAME [--install]  Detect/import/generate/verify/install an AI-agent-ready tool
  analyze <openapi.json> [--json]         Analyze an OpenAPI spec and print a blueprint
  import-har <capture.har> --out <json>   Convert a browser HAR capture to OpenAPI
  import-graphql <source> --out <json>    Convert GraphQL SDL/introspection/endpoint to OpenAPI
  discover <url> [--out openapi.json]     Discover OpenAPI/Swagger from a website
  seed-har <url> [<url> ...] --out <har>  Fetch URL(s) via plain HTTP and emit a HAR (no browser)
                                          [--method GET] [--header 'k: v'] [--data '{...}'] [--title T]
  record <url> --out <capture.har.json>   Record browser network traffic with Playwright
                                          [--storage-state state.json] [--save-storage-state state.json]
  login <url> --out <state.json>          Open a headed browser, let user log in, save storage-state
                                          [--storage-state state.json] [--wait-for-url pattern]
  generate <openapi.json> --out <dir>     Generate a Go CLI/MCP package
                                          [--targets go,mcp,skill,openclaw] [--policy policy.json]
                                          [--strict] [--min-score 0.7]
  init <name> [--spec openapi.json]       Create a product workspace
  verify <project-dir> [--json]           Build, test, smoke, MCP-handshake, and write proofs
  scorecard <project-dir> [--verify]      Score a generated package
  extract <url> -p PROMPT [-s SCHEMA]     LLM-powered extraction (Anthropic/OpenAI/Ollama)
                                          [--provider anthropic|openai|ollama] [--model NAME]
  scrape <url> [--out file.md]            Main-content → markdown (Firecrawl alternative)
  try <url> [--json]                      Classify a URL and print the exact next command
  quick <url|intent>                      URL-to-tool or intent-to-catalog fast path
  search <intent>                         Catalog discovery: top matching ops
  run <intent> [--dry-run] [--no-llm]     LLM-routed (or top-match) execute against catalog
  compare <tool-a> <tool-b> --op OP       Side-by-side call across tools
  diff <old-spec> <new-spec>              Semantic OpenAPI diff
  install starter-pack                    Build and install the curated starter tools
  publish <tool-dir>|--all [--fix-assets] Validate catalog-ready proofs/skills/OpenClaw
  registry list|validate|sync [file]      Work with the library registry
  site [build]                            Print path or rebuild web/data.js

Examples:
  node bin/gutenberg.js sports nba today
  node bin/gutenberg.js sports serie-a today --team juventus
  node bin/gutenberg.js travel rom-par june --adults 1 --currency EUR
  node bin/gutenberg.js recipes list
  node bin/gutenberg.js plan samples/petstore-openapi.json --name petstore
  node bin/gutenberg.js forge samples/petstore-openapi.json --out generated/petstore-go --name petstore --install --force
  node bin/gutenberg.js quick "top hacker news stories" --no-llm
  node bin/gutenberg.js create my-site --from https://example.com --kind har
  node bin/gutenberg.js analyze samples/petstore-openapi.yaml
  node bin/gutenberg.js import-har samples/simple.har.json --out generated/capture.openapi.json --name captured-api
  node bin/gutenberg.js import-graphql samples/github-introspection.json --out generated/github.openapi.json --name github
  node bin/gutenberg.js generate samples/petstore-openapi.json --out generated/petstore-go --name petstore --force
  node bin/gutenberg.js scorecard generated/petstore-go
`);
}

function printAnalysis(blueprint, json) {
  if (json) {
    console.log(JSON.stringify(blueprint, null, 2));
    return;
  }
  console.log(`${blueprint.name} (${blueprint.slug})`);
  console.log(`Operations: ${blueprint.operations.length}`);
  console.log(`Groups: ${blueprint.tags.join(", ") || "none"}`);
  console.log(`Base URLs: ${blueprint.baseUrls.join(", ") || "none detected"}`);
  console.log(`Auth: ${blueprint.auth.mode}`);
  console.log("");
  console.log("Thesis:");
  console.log(`  ${blueprint.insights.thesis}`);
  console.log("");
  console.log("Recommended generated commands:");
  for (const command of blueprint.insights.recommendedCommands) {
    console.log(`  - ${command}`);
  }
}

function printPlan(plan) {
  console.log(`${plan.name} (${plan.slug})`);
  console.log(`Source: ${plan.kind} — ${plan.source}`);
  console.log(`Operations: ${plan.operations.length}`);
  console.log(`Base URLs: ${plan.baseUrls.join(", ") || "none detected"}`);
  console.log(`Auth: ${plan.auth.mode}`);
  console.log("");
  console.log("Risk:");
  const counts = plan.operations.reduce((acc, operation) => {
    acc[operation.risk] = (acc[operation.risk] || 0) + 1;
    return acc;
  }, {});
  console.log(`  read=${counts.read || 0} write=${counts.write || 0} destructive=${counts.destructive || 0}`);
  console.log("");
  console.log("Policy:");
  for (const rule of plan.policy.rules) {
    console.log(`  ${rule.risk}: ${rule.action}${rule.requiresYes ? " (--yes required)" : ""}`);
  }
  console.log("");
  console.log("Hero candidates:");
  for (const hero of plan.heroes.slice(0, 8)) {
    console.log(`  ${hero.alias} -> ${hero.operationId}`);
  }
  console.log("");
  console.log("Next:");
  for (const command of plan.nextCommands) console.log(`  ${command}`);
}

function printForgeResult(result) {
  console.log(`Forged ${result.manifest.name}`);
  console.log(`Directory: ${result.outDir}`);
  if (result.normalizedSpec) console.log(`Spec:      ${result.normalizedSpec}`);
  console.log(`Proofs:    ${result.verification.artifacts.proofDir}`);
  console.log(`Verified:  ${result.verification.ok ? "yes" : "no"}`);
  if (result.verification.scorecard) {
    console.log(`Score:     ${result.verification.scorecard.score}/${result.verification.scorecard.maxScore} (${result.verification.scorecard.grade})`);
    if (result.verification.scorecard.badges?.length) {
      console.log(`Badges:    ${result.verification.scorecard.badges.join(", ")}`);
    }
  }
  if (result.install) {
    console.log(`Installed: ${result.install.binPath}`);
  }
  console.log(`Run:       ${result.command}`);
}

function printVerification(result) {
  console.log(`${result.tool}: ${result.ok ? "verified" : "FAILED"}`);
  console.log(`Proofs: ${result.artifacts.proofDir}`);
  for (const check of result.checks) {
    console.log(`${check.passed ? "PASS" : "FAIL"}  ${check.id}: ${check.detail}`);
  }
  if (result.scorecard) {
    console.log(`Score: ${result.scorecard.score}/${result.scorecard.maxScore} (${result.scorecard.grade}, ratio=${(result.scorecard.ratio || 0).toFixed(2)})`);
  }
}

function printScore(score) {
  console.log(`${score.summary}`);
  console.log(`Score: ${score.score}/${score.maxScore} (${score.grade}, ratio=${(score.ratio || 0).toFixed(2)})`);
  if (score.badges?.length) {
    console.log(`Badges: ${score.badges.join(", ")}`);
  }
  if (score.dimensions) {
    console.log("");
    console.log("Dimensions:");
    for (const [name, dim] of Object.entries(score.dimensions)) {
      const ratio = dim.max > 0 ? `${dim.score}/${dim.max}` : "n/a";
      console.log(`  ${name.padEnd(12)} ${ratio.padEnd(8)} (${(dim.ratio || 0).toFixed(2)})`);
    }
    console.log("");
  }
  for (const check of score.checks) {
    console.log(`${check.passed ? "PASS" : "FAIL"}  ${check.label}`);
  }
}

function printPublishResult(result) {
  console.log(`${result.slug}: ${result.ready ? "catalog-ready" : "blocked"}`);
  if (result.score) {
    console.log(`Score: ${result.score.score}/${result.score.maxScore} (${result.score.grade})`);
  }
  if (result.fixedAssets?.length) {
    console.log(`Updated assets: ${result.fixedAssets.length}`);
    for (const file of result.fixedAssets) console.log(`  ${file}`);
  }
  for (const error of result.errors) console.log(`FAIL ${error}`);
  if (!result.ready) process.exitCode = 1;
}

function printPublishAll(result) {
  console.log(`Publish check: ${result.ready ? "ready" : "blocked"} (${result.results.length} tool(s))`);
  if (result.registry) console.log(`Registry synced: ${result.registry.path} (${result.registry.tools} tool(s))`);
  for (const item of result.results) {
    const suffix = item.fixedAssets?.length ? ` fixed=${item.fixedAssets.length}` : "";
    console.log(`${item.ready ? "PASS" : "FAIL"} ${item.slug}${suffix}`);
    for (const error of item.errors) console.log(`  - ${error}`);
  }
  if (!result.ready) process.exitCode = 1;
}

function registry(positionals, options) {
  const action = positionals[0] || "list";
  const registryPath = path.resolve(positionals[1] || path.join(rootDir, "library", "registry.json"));
  if (action === "sync") {
    const data = buildRegistryFromCatalog(rootDir);
    writeRegistryFile(data, registryPath);
    if (options.json) {
      console.log(JSON.stringify(data, null, 2));
      return;
    }
    console.log(`Registry synced: ${registryPath}`);
    console.log(`Tools: ${data.tools.length}`);
    return;
  }
  const data = readJson(registryPath);
  if (action === "list") {
    if (options.json) {
      console.log(JSON.stringify(data, null, 2));
      return;
    }
    for (const tool of data.tools || []) {
      console.log(`${tool.slug.padEnd(16)} ${tool.status.padEnd(10)} ${tool.description}`);
    }
    return;
  }
  if (action === "validate") {
    const errors = validateRegistry(data);
    if (options.json) {
      if (errors.length > 0) process.exitCode = 1;
      console.log(JSON.stringify({ valid: errors.length === 0, errors }, null, 2));
      return;
    }
    if (errors.length > 0) {
      for (const error of errors) console.log(`FAIL ${error}`);
      process.exitCode = 1;
    } else {
      console.log(`Registry valid: ${registryPath}`);
    }
    return;
  }
  throw usage(`Unknown registry action: ${action}`);
}

function validateRegistry(data) {
  const errors = [];
  if (!["gutenberg.registry.v1", "black-forge.registry.v1"].includes(data.schemaVersion)) errors.push("schemaVersion must be gutenberg.registry.v1");
  if (!Array.isArray(data.tools)) errors.push("tools must be an array");
  const seen = new Set();
  for (const tool of data.tools || []) {
    if (!tool.slug) errors.push("tool missing slug");
    if (seen.has(tool.slug)) errors.push(`duplicate tool slug: ${tool.slug}`);
    seen.add(tool.slug);
    if (!tool.package) errors.push(`tool ${tool.slug} missing package`);
    if (!["prototype", "verified", "deprecated"].includes(tool.status)) errors.push(`tool ${tool.slug} has invalid status`);
  }
  return errors;
}

async function doctor() {
  const goPath = path.join(rootDir, "scripts", "use-go.sh");
  const generatedDir = path.join(rootDir, "generated");
  const libraryDir = path.join(rootDir, "library", "tools");
  const playwrightDir = path.join(rootDir, "node_modules", "playwright");
  const checks = [
    ["node", process.version, Number(process.versions.node.split(".")[0]) >= 20],
    ["fetch", typeof fetch, typeof fetch === "function"],
    ["go-wrapper", goPath, fs.existsSync(goPath)],
    ["cwd", process.cwd(), true],
    ["root", rootDir, fs.existsSync(rootDir)],
    ["generated:writable", generatedDir, canWriteDirectory(generatedDir)],
    ["library:tools", libraryDir, fs.existsSync(libraryDir)],
    ["playwright", fs.existsSync(playwrightDir) ? "installed" : "not installed", fs.existsSync(playwrightDir)]
  ];
  if (fs.existsSync(goPath)) {
    const go = spawnSync(goPath, ["version"], { encoding: "utf8", timeout: 10000 });
    checks.push(["go", go.status === 0 ? go.stdout.trim() : (go.stderr || "failed").trim(), go.status === 0]);
  }
  const anthropic = Boolean(process.env.ANTHROPIC_API_KEY);
  const openai = Boolean(process.env.OPENAI_API_KEY);
  checks.push(["llm:anthropic", anthropic ? "ANTHROPIC_API_KEY set" : "not set", anthropic]);
  checks.push(["llm:openai", openai ? "OPENAI_API_KEY set" : "not set", openai]);

  const ollamaHost = process.env.OLLAMA_HOST || "http://localhost:11434";
  let ollamaInfo = "unreachable";
  let ollamaOk = false;
  try {
    const response = await fetch(`${ollamaHost.replace(/\/$/, "")}/api/tags`, { signal: AbortSignal.timeout(2000) });
    if (response.ok) {
      const data = await response.json();
      const count = Array.isArray(data.models) ? data.models.length : 0;
      ollamaInfo = `${ollamaHost} (${count} model${count === 1 ? "" : "s"})`;
      ollamaOk = true;
    } else {
      ollamaInfo = `${ollamaHost} HTTP ${response.status}`;
    }
  } catch (error) {
    ollamaInfo = `${ollamaHost} ${error.name === "TimeoutError" ? "timeout" : "down"}`;
  }
  checks.push(["llm:ollama", ollamaInfo, ollamaOk]);

  const llmReady = anthropic || openai || ollamaOk;
  checks.push(["llm:any-provider", llmReady ? "yes" : "no — extract/run require a provider", llmReady]);

  const { browserbaseStatus } = await import("./core/recorder-browserbase.js");
  const bb = browserbaseStatus();
  checks.push(["browserbase", bb.ready ? `ready (${bb.apiUrl})` : "credentials missing — record --backend browserbase needs BROWSERBASE_API_KEY + BROWSERBASE_PROJECT_ID", bb.ready]);
  checks.push(["forge:llm-independent", "plan/forge/generate do not require an LLM provider", true]);

  for (const [name, value, passed] of checks) {
    console.log(`${passed ? "PASS" : "FAIL"}  ${name}: ${value}`);
  }
}

function canWriteDirectory(dir) {
  try {
    fs.mkdirSync(dir, { recursive: true });
    const probe = path.join(dir, `.doctor-${process.pid}`);
    fs.writeFileSync(probe, "ok");
    fs.unlinkSync(probe);
    return true;
  } catch {
    return false;
  }
}

function researchTemplate(name) {
  return `# ${name} Research

## Target Users

- Operators who need compact, scriptable workflows.
- Agents that need reliable tools with predictable JSON output.

## Competitor Sweep

- Official CLI:
- Existing MCP server:
- SDK:
- No-code connector:

## Non-Obvious Insight

What can this API become once its data is local, searchable, and composable?

## Risk Register

- Auth model:
- Rate limits:
- Destructive operations:
- Terms of service:
- Sensitive data:
`;
}

function usage(message) {
  const error = new Error(message);
  error.exitCode = 2;
  return error;
}

function headerOptions(value) {
  const list = Array.isArray(value) ? value : value ? [value] : [];
  const headers = {};
  for (const item of list) {
    const colon = String(item).indexOf(":");
    if (colon === -1) continue;
    headers[item.slice(0, colon).trim()] = item.slice(colon + 1).trim();
  }
  return headers;
}
