import fs from "node:fs";
import path from "node:path";
import os from "node:os";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";

const gutenbergRoot = path.dirname(path.dirname(path.dirname(fileURLToPath(import.meta.url))));
const useGo = path.join(gutenbergRoot, "scripts", "use-go.sh");

export const STARTER_PACK = [
  "hacker-news",
  "espn",
  "wikipedia-it",
  "tvmaze",
  "open-meteo",
  "github"
];

export function installTool(toolDir, options = {}) {
  const resolvedDir = path.resolve(toolDir);
  const manifestPath = path.join(resolvedDir, "gutenberg.manifest.json");
  if (!fs.existsSync(manifestPath)) {
    const error = new Error(`No gutenberg.manifest.json in ${resolvedDir} — not a Gutenberg-generated project.`);
    error.exitCode = 2;
    throw error;
  }
  const manifest = JSON.parse(fs.readFileSync(manifestPath, "utf8"));
  const slug = manifest.slug;
  const binDir = path.resolve(options.prefix || path.join(os.homedir(), ".local", "bin"));
  fs.mkdirSync(binDir, { recursive: true });
  const binPath = path.join(binDir, slug);

  if (!options.skipTidy) {
    const tidy = spawnSync(useGo, ["mod", "tidy"], { cwd: resolvedDir, stdio: "inherit" });
    if (tidy.status !== 0) {
      const error = new Error(`go mod tidy failed in ${resolvedDir} (exit ${tidy.status})`);
      error.exitCode = tidy.status || 1;
      throw error;
    }
  }

  const cmdDir = path.join(resolvedDir, "cmd", slug);
  if (!fs.existsSync(cmdDir)) {
    const error = new Error(`Expected cmd/${slug} in ${resolvedDir}`);
    error.exitCode = 2;
    throw error;
  }
  const build = spawnSync(useGo, ["build", "-o", binPath, `./cmd/${slug}`], { cwd: resolvedDir, stdio: "inherit" });
  if (build.status !== 0) {
    const error = new Error(`go build failed (exit ${build.status})`);
    error.exitCode = build.status || 1;
    throw error;
  }

  const pathEntries = (process.env.PATH || "").split(path.delimiter);
  const onPath = pathEntries.includes(binDir);
  return { slug, binPath, binDir, onPath };
}

export function installStarterPack(rootDir, options = {}) {
  const tools = options.tools || STARTER_PACK;
  const results = [];
  for (const slug of tools) {
    const toolDir = path.join(rootDir, "library", "tools", slug);
    if (!fs.existsSync(path.join(toolDir, "gutenberg.manifest.json"))) {
      const error = new Error(`Starter-pack tool is missing: ${slug} (${toolDir})`);
      error.exitCode = 2;
      throw error;
    }
    if (options.dryRun) {
      results.push({ slug, toolDir, skipped: true });
      continue;
    }
    results.push(installTool(toolDir, {
      prefix: options.prefix,
      skipTidy: Boolean(options.skipTidy)
    }));
  }
  return {
    schemaVersion: "gutenberg.install.starter-pack.v1",
    tools,
    installed: results
  };
}
