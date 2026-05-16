import fs from "node:fs";
import path from "node:path";

export function searchCatalog(rootDir, intent, options = {}) {
  const tokens = normalizeTokens(intent);
  if (tokens.length === 0) {
    return { intent, results: [] };
  }
  const requiredTokens = tokens.filter((token) => !GENERIC_INTENT_TOKENS.has(token));

  const toolsDir = path.join(rootDir, "library", "tools");
  if (!fs.existsSync(toolsDir)) return { intent, results: [] };

  const results = [];
  for (const dirent of fs.readdirSync(toolsDir, { withFileTypes: true })) {
    if (!dirent.isDirectory()) continue;
    const manifestPath = path.join(toolsDir, dirent.name, "gutenberg.manifest.json");
    if (!fs.existsSync(manifestPath)) continue;
    let manifest;
    try {
      manifest = JSON.parse(fs.readFileSync(manifestPath, "utf8"));
    } catch {
      continue;
    }

    const toolMatch = scoreMatch(tokens, [manifest.name, manifest.slug, manifest.description || "", (manifest.tags || []).join(" ")]);
    const heroes = manifest.heroes || [];
    const operations = manifest.operations || [];

    for (const hero of heroes) {
      const heroScore = scoreMatch(tokens, [hero.alias, hero.summary, hero.path]);
      const match = combineMatches(toolMatch, heroScore);
      if (isRelevantMatch(match, requiredTokens)) {
        results.push({
          tool: manifest.slug,
          toolName: manifest.name,
          kind: "hero",
          alias: hero.alias,
          operationId: hero.operationId,
          summary: hero.summary,
          score: heroScore.score * 2 + toolMatch.score,
          command: `${manifest.slug} ${hero.alias}`
        });
      }
    }
    for (const op of operations.slice(0, 200)) {
      const opScore = scoreMatch(tokens, [op.id, op.summary, op.path, op.tag, op.kind || ""]);
      const match = combineMatches(toolMatch, opScore);
      if (isRelevantMatch(match, requiredTokens)) {
        results.push({
          tool: manifest.slug,
          toolName: manifest.name,
          kind: op.kind || "operation",
          operationId: op.id,
          summary: op.summary,
          method: op.method,
          path: op.path,
          score: opScore.score + toolMatch.score,
          command: `${manifest.slug} call ${op.id}`
        });
      }
    }
  }

  results.sort((a, b) => b.score - a.score);
  return {
    intent,
    results: results.slice(0, options.limit || 5)
  };
}

const GENERIC_INTENT_TOKENS = new Set([
  "api",
  "call",
  "data",
  "fetch",
  "find",
  "get",
  "list",
  "read",
  "result",
  "results",
  "score",
  "search",
  "show"
]);

function normalizeTokens(intent) {
  return String(intent || "")
    .toLowerCase()
    .split(/\s+/)
    .map((t) => t.trim().replace(/[^a-z0-9_-]/g, ""))
    .filter((t) => t.length > 1)
    .map((t) => t.endsWith("s") && t.length > 3 ? t.slice(0, -1) : t);
}

function scoreMatch(tokens, fields) {
  let score = 0;
  const matched = new Set();
  const haystack = fields.filter(Boolean).join(" ").toLowerCase();
  const words = haystack.split(/[^a-z0-9]+/).filter(Boolean);
  for (const token of tokens) {
    if (token.length <= 3) {
      if (words.includes(token)) {
        score += 3;
        matched.add(token);
      }
    } else if (haystack.split(/\s+/).some((word) => word.startsWith(token))) {
      score += 1;
      matched.add(token);
    } else if (haystack.includes(token)) {
      score += 2;
      matched.add(token);
    }
  }
  return { score, matched };
}

function combineMatches(...matches) {
  const matched = new Set();
  let score = 0;
  for (const match of matches) {
    score += match.score;
    for (const token of match.matched) matched.add(token);
  }
  return { score, matched };
}

function isRelevantMatch(match, requiredTokens) {
  if (match.score <= 0) return false;
  if (requiredTokens.length === 0) return true;
  return requiredTokens.some((token) => match.matched.has(token));
}
