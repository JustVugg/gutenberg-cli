import fs from "node:fs";
import path from "node:path";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import { buildBlueprint, loadOpenApi } from "./openapi.js";
import { generateProject } from "./render.js";
import { readLock } from "./lockfile.js";

const gutenbergRoot = path.dirname(path.dirname(path.dirname(fileURLToPath(import.meta.url))));
const useGo = path.join(gutenbergRoot, "scripts", "use-go.sh");

const KEEP_START = /\/\/\s*gutenberg:keep\s+([^\s]+)\s*$/m;
const ALL_KEEP_BLOCKS = /\/\/\s*gutenberg:keep\s+([^\s]+)[\s\S]*?\/\/\s*gutenberg:end-keep/g;

export function upgradeProject(projectDir, options = {}) {
  const resolvedDir = path.resolve(projectDir);
  const manifestPath = path.join(resolvedDir, "gutenberg.manifest.json");
  if (!fs.existsSync(manifestPath)) {
    const error = new Error(`No gutenberg.manifest.json in ${resolvedDir} — not a Gutenberg-generated project.`);
    error.exitCode = 2;
    throw error;
  }
  const manifest = JSON.parse(fs.readFileSync(manifestPath, "utf8"));
  const lock = readLock(resolvedDir);
  const specPath = options.spec
    ? path.resolve(options.spec)
    : manifest?.provenance?.spec?.path || lock?.spec?.path || manifest?.source;
  if (!specPath || !fs.existsSync(specPath)) {
    const error = new Error(`Cannot locate source spec (provenance.spec.path or lock.spec.path). Pass --spec <path>.`);
    error.exitCode = 2;
    throw error;
  }

  const preserved = collectKeepBlocks(resolvedDir);

  const blueprint = buildBlueprint(loadOpenApi(specPath), specPath, manifest.name || manifest.slug);
  const result = generateProject(blueprint, resolvedDir, {
    force: true,
    name: manifest.slug,
    displayName: manifest.name,
    targets: options.targets,
    specPath
  });

  const restoreReport = restoreKeepBlocks(resolvedDir, preserved);

  let tidyResult = null;
  if (!options.noTidy && fs.existsSync(path.join(resolvedDir, "go.mod")) && fs.existsSync(useGo)) {
    const tidy = spawnSync(useGo, ["mod", "tidy"], { cwd: resolvedDir });
    tidyResult = { status: tidy.status, stderr: tidy.stderr ? tidy.stderr.toString() : "" };
  }

  return { ...result, preservedBlocks: preserved.length, restoredBlocks: restoreReport.restored, orphanedBlocks: restoreReport.orphaned, tidy: tidyResult };
}

function collectKeepBlocks(projectDir) {
  const blocks = [];
  walk(projectDir, (file) => {
    if (!/\.(go|js|md)$/.test(file)) return;
    const content = fs.readFileSync(file, "utf8");
    const matches = content.matchAll(ALL_KEEP_BLOCKS);
    for (const match of matches) {
      const blockId = match[1];
      blocks.push({
        relative: path.relative(projectDir, file),
        blockId,
        content: match[0]
      });
    }
  });
  return blocks;
}

function restoreKeepBlocks(projectDir, blocks) {
  const restored = [];
  const orphaned = [];
  for (const block of blocks) {
    const file = path.join(projectDir, block.relative);
    if (!fs.existsSync(file)) {
      orphaned.push({ ...block, reason: "file-missing" });
      continue;
    }
    const content = fs.readFileSync(file, "utf8");
    const startPattern = new RegExp(`\\/\\/\\s*gutenberg:keep\\s+${escapeRegex(block.blockId)}[\\s\\S]*?\\/\\/\\s*gutenberg:end-keep`);
    if (!startPattern.test(content)) {
      orphaned.push({ ...block, reason: "marker-missing" });
      continue;
    }
    const updated = content.replace(startPattern, block.content);
    fs.writeFileSync(file, updated);
    restored.push(block);
  }
  return { restored, orphaned };
}

function walk(dir, callback) {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    if (entry.name === "node_modules" || entry.name === ".git") continue;
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      walk(full, callback);
    } else if (entry.isFile()) {
      callback(full);
    }
  }
}

function escapeRegex(value) {
  return String(value).replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}
