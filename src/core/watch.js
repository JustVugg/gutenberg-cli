import fs from "node:fs";
import crypto from "node:crypto";
import path from "node:path";
import { diffSpecs } from "./diff.js";
import { upgradeProject } from "./upgrade.js";

export async function watchSpec(toolDir, specUrl, options = {}) {
  const resolvedDir = path.resolve(toolDir);
  const manifestPath = path.join(resolvedDir, "gutenberg.manifest.json");
  if (!fs.existsSync(manifestPath)) {
    const error = new Error(`No manifest in ${resolvedDir}`);
    error.exitCode = 2;
    throw error;
  }

  const response = await fetch(specUrl);
  if (!response.ok) {
    const error = new Error(`HTTP ${response.status} fetching ${specUrl}`);
    error.exitCode = 1;
    throw error;
  }
  const body = await response.text();
  const sha = crypto.createHash("sha256").update(body).digest("hex");

  const lockPath = path.join(resolvedDir, "gutenberg.lock.json");
  let lockedSha = null;
  if (fs.existsSync(lockPath)) {
    try {
      const lock = JSON.parse(fs.readFileSync(lockPath, "utf8"));
      lockedSha = lock.spec?.sha256 || null;
    } catch {}
  }

  const changed = !lockedSha || lockedSha !== sha;
  const stagingPath = options.stagingPath || path.join(resolvedDir, ".gutenberg-watch-spec.json");
  fs.writeFileSync(stagingPath, body, "utf8");

  if (!changed) {
    return { changed: false, sha, lockedSha, message: "Spec unchanged." };
  }

  let diffSummary = null;
  if (lockedSha) {
    try {
      const previousPath = path.join(resolvedDir, ".gutenberg-watch-spec.previous.json");
      if (fs.existsSync(previousPath)) {
        diffSummary = diffSpecs(previousPath, stagingPath);
      }
    } catch (error) {
      diffSummary = { error: error.message };
    }
  }

  if (!options.regenerate) {
    return { changed: true, sha, lockedSha, diff: diffSummary, message: "Spec changed. Re-run with --regenerate to upgrade." };
  }

  const result = upgradeProject(resolvedDir, { spec: stagingPath, noTidy: Boolean(options.noTidy) });
  fs.copyFileSync(stagingPath, path.join(resolvedDir, ".gutenberg-watch-spec.previous.json"));
  return { changed: true, sha, lockedSha, diff: diffSummary, upgrade: { manifestName: result.manifest.name, preservedBlocks: result.preservedBlocks } };
}
