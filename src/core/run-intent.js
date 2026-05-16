import fs from "node:fs";
import path from "node:path";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import { searchCatalog } from "./catalog-search.js";
import { extractFromText } from "./extract.js";

const gutenbergRoot = path.dirname(path.dirname(path.dirname(fileURLToPath(import.meta.url))));
const useGo = path.join(gutenbergRoot, "scripts", "use-go.sh");

export async function runIntent(rootDir, intent, options = {}) {
  if (!intent) {
    const error = new Error("Missing intent.");
    error.exitCode = 2;
    throw error;
  }
  const candidates = searchCatalog(rootDir, intent, { limit: 8 }).results;
  if (candidates.length === 0) {
    return { intent, plan: null, message: "No candidate operation in catalog." };
  }

  const plan = await pickPlan(intent, candidates, options);
  if (options.dryRun || !plan?.command) {
    return { intent, candidates, plan };
  }

  const result = executePlan(rootDir, plan);
  return { intent, plan, result };
}

async function pickPlan(intent, candidates, options) {
  const hasKey = process.env.ANTHROPIC_API_KEY || process.env.OPENAI_API_KEY;
  if (!hasKey || options.noLlm) {
    const best = candidates[0];
    return planFromCandidate(best);
  }
  const summary = candidates.map((hit, i) => `${i + 1}. ${hit.command} — ${hit.summary || ""}`).join("\n");
  const prompt = `User intent: "${intent}"\n\nAvailable commands in the local catalog:\n${summary}\n\nChoose ONE command and any required parameters as JSON.`;
  const schema = {
    type: "object",
    required: ["command"],
    properties: {
      command: { type: "string", description: "Full command to run, e.g. 'espn nba' or 'github call meta/root'" },
      params: { type: "object" },
      headers: { type: "object" },
      reason: { type: "string" }
    }
  };
  try {
    const { parsed } = await extractFromText(prompt, { prompt: "Pick the best command for the intent.", schema, provider: options.provider, model: options.model });
    if (parsed?.command) return parsed;
  } catch (error) {
    return { ...planFromCandidate(candidates[0]), reason: `LLM unavailable (${error.message}); falling back to top match` };
  }
  return planFromCandidate(candidates[0]);
}

function planFromCandidate(candidate) {
  return { command: candidate.command, params: {}, reason: `top catalog match (score=${candidate.score})` };
}

function executePlan(rootDir, plan) {
  const parts = String(plan.command).split(/\s+/).filter(Boolean);
  const toolSlug = parts[0];
  const args = parts.slice(1);
  for (const [key, value] of Object.entries(plan.params || {})) {
    args.push("--param", `${key}=${value}`);
  }
  for (const [key, value] of Object.entries(plan.headers || {})) {
    args.push("--header", `${key}: ${value}`);
  }

  const toolDir = path.join(rootDir, "library", "tools", toolSlug);
  if (!fs.existsSync(toolDir)) {
    return { error: `tool not found: ${toolSlug}` };
  }
  const child = spawnSync(useGo, ["run", `./cmd/${toolSlug}`, ...args], {
    cwd: toolDir,
    encoding: "utf8"
  });
  return { stdout: child.stdout, stderr: child.stderr, status: child.status };
}
