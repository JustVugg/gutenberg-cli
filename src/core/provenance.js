import fs from "node:fs";
import crypto from "node:crypto";
import path from "node:path";
import { fileURLToPath } from "node:url";

const PROVENANCE_SCHEMA = "gutenberg.provenance.v1";

let cachedVersion = null;

export function getGutenbergVersion() {
  if (cachedVersion) return cachedVersion;
  try {
    const here = path.dirname(fileURLToPath(import.meta.url));
    const pkgPath = path.resolve(here, "..", "..", "package.json");
    const pkg = JSON.parse(fs.readFileSync(pkgPath, "utf8"));
    cachedVersion = pkg.version || "0.0.0";
  } catch {
    cachedVersion = "0.0.0";
  }
  return cachedVersion;
}

export function hashFile(filePath) {
  if (!filePath || !fs.existsSync(filePath)) return null;
  const content = fs.readFileSync(filePath);
  const sha = crypto.createHash("sha256").update(content).digest("hex");
  return { sha256: sha, size: content.length };
}

export function computeProvenance({ specPath, recipePath, targets, name }) {
  const spec = specPath ? hashFile(specPath) : null;
  const recipe = recipePath ? hashFile(recipePath) : null;
  return {
    schemaVersion: PROVENANCE_SCHEMA,
    generatedAt: new Date().toISOString(),
    gutenbergVersion: getGutenbergVersion(),
    name: name || null,
    spec: spec ? { path: specPath, ...spec } : null,
    recipe: recipe ? { path: recipePath, ...recipe } : null,
    targets: Array.isArray(targets) ? [...targets] : []
  };
}

export function attachScorecard(provenance, scorecard) {
  return {
    ...provenance,
    scorecard: scorecard
      ? {
          score: scorecard.score,
          maxScore: scorecard.maxScore,
          grade: scorecard.grade,
          dimensions: scorecard.dimensions || null
        }
      : null
  };
}
