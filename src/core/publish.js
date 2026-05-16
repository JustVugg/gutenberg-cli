import fs from "node:fs";
import path from "node:path";
import { readJson } from "./fs.js";
import { scoreProject } from "./scorecard.js";
import { generateClaudeSkill } from "./render-skill.js";
import { generateOpenClawSkill } from "./render-openclaw.js";
import { buildRegistryFromCatalog, writeRegistryFile } from "./catalog.js";

export function ensureAgentAssets(toolDir, options = {}) {
  const resolvedDir = path.resolve(toolDir);
  const manifestPath = path.join(resolvedDir, "gutenberg.manifest.json");
  if (!fs.existsSync(manifestPath)) {
    const error = new Error(`No gutenberg.manifest.json in ${resolvedDir}`);
    error.exitCode = 2;
    throw error;
  }
  const manifest = readJson(manifestPath);
  const changed = [];
  const skillFile = path.join(resolvedDir, "skills", manifest.slug, "SKILL.md");
  const openclawJson = path.join(resolvedDir, "openclaw", manifest.slug, "skill.json");
  const openclawMd = path.join(resolvedDir, "openclaw", manifest.slug, "skill.md");

  if (options.force || !fs.existsSync(skillFile)) {
    const skill = generateClaudeSkill(manifest, resolvedDir);
    changed.push(skill.file);
  }
  if (options.force || !fs.existsSync(openclawJson) || !fs.existsSync(openclawMd)) {
    const openclaw = generateOpenClawSkill(manifest, resolvedDir);
    changed.push(openclaw.jsonFile, openclaw.mdFile);
  }
  return { toolDir: resolvedDir, slug: manifest.slug, changed };
}

export function publishTool(toolDir, options = {}) {
  const resolvedDir = path.resolve(toolDir);
  const manifestPath = path.join(resolvedDir, "gutenberg.manifest.json");
  if (!fs.existsSync(manifestPath)) {
    const error = new Error(`No gutenberg.manifest.json in ${resolvedDir}`);
    error.exitCode = 2;
    throw error;
  }
  const manifest = readJson(manifestPath);
  const fixed = options.fixAssets ? ensureAgentAssets(resolvedDir, { force: Boolean(options.forceAssets) }) : { changed: [] };
  const proofPath = path.join(resolvedDir, "proofs", "verification.json");
  const skillPath = path.join(resolvedDir, "skills", manifest.slug, "SKILL.md");
  const openclawPath = path.join(resolvedDir, "openclaw", manifest.slug, "skill.json");
  const errors = [];
  let scorecard = null;
  let proof = null;

  try {
    scorecard = scoreProject(resolvedDir);
  } catch (error) {
    errors.push(`scorecard failed: ${error.message}`);
  }
  if (fs.existsSync(proofPath)) {
    try {
      proof = readJson(proofPath);
      if (!proof.ok) errors.push("proofs/verification.json exists but ok is false");
    } catch (error) {
      errors.push(`invalid proofs/verification.json: ${error.message}`);
    }
  } else {
    errors.push("missing proofs/verification.json");
  }
  if (!fs.existsSync(skillPath)) errors.push(`missing skills/${manifest.slug}/SKILL.md`);
  if (!fs.existsSync(openclawPath)) errors.push(`missing openclaw/${manifest.slug}/skill.json`);

  return {
    schemaVersion: "gutenberg.publish.result.v1",
    slug: manifest.slug,
    toolDir: resolvedDir,
    ready: errors.length === 0,
    errors,
    fixedAssets: fixed.changed,
    score: scorecard ? {
      score: scorecard.score,
      maxScore: scorecard.maxScore,
      grade: scorecard.grade,
      badges: scorecard.badges || []
    } : null,
    proof: proof ? {
      ok: Boolean(proof.ok),
      verifiedAt: proof.verifiedAt || null,
      checks: Array.isArray(proof.checks) ? proof.checks.length : 0
    } : null
  };
}

export function publishAll(rootDir, options = {}) {
  const toolsDir = path.join(rootDir, "library", "tools");
  const results = [];
  if (!fs.existsSync(toolsDir)) {
    return { schemaVersion: "gutenberg.publish.all.v1", results, registry: null };
  }
  for (const dirent of fs.readdirSync(toolsDir, { withFileTypes: true })) {
    if (!dirent.isDirectory()) continue;
    const toolDir = path.join(toolsDir, dirent.name);
    if (!fs.existsSync(path.join(toolDir, "gutenberg.manifest.json"))) continue;
    results.push(publishTool(toolDir, options));
  }
  let registry = null;
  if (options.syncRegistry) {
    registry = buildRegistryFromCatalog(rootDir);
    writeRegistryFile(registry, options.registryPath || path.join(rootDir, "library", "registry.json"));
  }
  return {
    schemaVersion: "gutenberg.publish.all.v1",
    ready: results.every((result) => result.ready),
    results,
    registry: registry ? {
      tools: registry.tools.length,
      path: options.registryPath || path.join(rootDir, "library", "registry.json")
    } : null
  };
}
