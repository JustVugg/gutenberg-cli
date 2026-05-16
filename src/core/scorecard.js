import fs from "node:fs";
import path from "node:path";
import { readJson } from "./fs.js";

const COMMON_REQUIRED_FILES = [
  "README.md",
  "docs/COOKBOOK.md",
  ".env.example"
];

const NODE_REQUIRED_FILES = [
  "package.json",
  "bin",
  "src/client.js",
  "src/cli.js",
  "src/mcp-server.js",
  "src/store.js",
  "tests/smoke.test.js"
];

const GO_REQUIRED_FILES = [
  "go.mod",
  "cmd",
  "internal/forge/manifest.go",
  "internal/forge/client.go",
  "internal/forge/store.go",
  "internal/forge/auth.go",
  "internal/forge/mcp.go",
  "internal/forge/forge_test.go"
];

const AGGREGATOR_REQUIRED_FILES = [
  "go.mod",
  "cmd",
  "internal/aggr/manifests.go",
  "internal/aggr/client.go",
  "internal/aggr/merge.go",
  "internal/aggr/rank.go",
  "internal/aggr/types.go",
  "internal/aggr/aggr_test.go"
];

export function scoreProject(projectDir) {
  const checks = [];
  const isGo = fs.existsSync(path.join(projectDir, "go.mod"));
  const manifestPathProbe = path.join(projectDir, "gutenberg.manifest.json");
  let isAggregator = false;
  if (fs.existsSync(manifestPathProbe)) {
    try {
      const probe = readJson(manifestPathProbe);
      isAggregator = probe.kind === "aggregator";
    } catch {}
  }
  const requiredFiles = isAggregator
    ? ["README.md", ".env.example", ...AGGREGATOR_REQUIRED_FILES]
    : [
        ...COMMON_REQUIRED_FILES,
        ...(isGo ? GO_REQUIRED_FILES : NODE_REQUIRED_FILES)
      ];

  for (const relative of requiredFiles) {
    const absolute = path.join(projectDir, relative);
    checks.push({
      id: `file:${relative}`,
      label: `Required artifact: ${relative}`,
      passed: fs.existsSync(absolute),
      points: 5
    });
  }

  const manifestPath = fs.existsSync(path.join(projectDir, "gutenberg.manifest.json"))
    ? path.join(projectDir, "gutenberg.manifest.json")
    : path.join(projectDir, "blackforge.manifest.json");
  let manifest = null;
  if (fs.existsSync(manifestPath)) {
    manifest = readJson(manifestPath);
    const operations = Array.isArray(manifest.operations) ? manifest.operations : [];
    const isAggregator = manifest.kind === "aggregator";
    checks.push({
      id: "manifest:file",
      label: "Manifest is present",
      passed: true,
      points: 5
    });
    checks.push({
      id: "manifest:operations",
      label: isAggregator ? "Aggregator manifest lists sources" : "Manifest contains operations",
      passed: isAggregator ? Array.isArray(manifest.sources) && manifest.sources.length > 0 : operations.length > 0,
      points: 10
    });
    checks.push({
      id: "manifest:insights",
      label: "Manifest contains non-obvious product insights",
      passed: Boolean(manifest.insights?.thesis) || isAggregator,
      points: 10
    });
    checks.push({
      id: "manifest:risk",
      label: "Mutating operations are risk classified",
      passed: isAggregator || (operations.length > 0 && operations.every((operation) => ["read", "write", "destructive"].includes(operation.risk))),
      points: 10
    });
    checks.push({
      id: "manifest:cache",
      label: "Cacheable read operations are identified",
      passed: isAggregator || operations.some((operation) => operation.cacheable),
      points: 10
    });
  } else {
    checks.push({
      id: "manifest:missing",
      label: "Manifest can be parsed",
      passed: false,
      points: 40
    });
  }

  if (isGo && !isAggregator) {
    const goModPath = path.join(projectDir, "go.mod");
    const goMod = fs.readFileSync(goModPath, "utf8");
    checks.push({
      id: "go:module",
      label: "Go module is present",
      passed: goMod.includes("module "),
      points: 10
    });
    checks.push({
      id: "go:mcp",
      label: "Go MCP runtime is present",
      passed: fs.existsSync(path.join(projectDir, "internal/forge/mcp.go")),
      points: 10
    });
    checks.push({
      id: "go:sqlite",
      label: "Go SQLite/FTS5 store is present",
      passed: fs.existsSync(path.join(projectDir, "internal/forge/store.go")) && fs.readFileSync(path.join(projectDir, "internal/forge/store.go"), "utf8").includes("fts5"),
      points: 10
    });
    checks.push({
      id: "go:oauth",
      label: "Go OAuth helper is present",
      passed: fs.existsSync(path.join(projectDir, "internal/forge/auth.go")),
      points: 10
    });
  } else if (isAggregator) {
    const goModPath = path.join(projectDir, "go.mod");
    const goMod = fs.existsSync(goModPath) ? fs.readFileSync(goModPath, "utf8") : "";
    checks.push({ id: "go:module", label: "Go module is present", passed: goMod.includes("module "), points: 10 });
    checks.push({ id: "go:aggregator", label: "Aggregator client present", passed: fs.existsSync(path.join(projectDir, "internal/aggr/client.go")), points: 10 });
    checks.push({ id: "go:merge", label: "Aggregator merge.go present", passed: fs.existsSync(path.join(projectDir, "internal/aggr/merge.go")), points: 10 });
    checks.push({ id: "go:rank", label: "Aggregator rank.go present", passed: fs.existsSync(path.join(projectDir, "internal/aggr/rank.go")), points: 10 });
  }

  const packagePath = path.join(projectDir, "package.json");
  if (fs.existsSync(packagePath)) {
    const pkg = readJson(packagePath);
    checks.push({
      id: "package:bin",
      label: "Package exposes a CLI binary",
      passed: Boolean(pkg.bin && Object.keys(pkg.bin).length > 0),
      points: 10
    });
    checks.push({
      id: "package:mcp",
      label: "Package exposes an MCP script",
      passed: Boolean(pkg.scripts?.mcp),
      points: 10
    });
  }

  const verification = readVerification(projectDir);
  checks.push({
    id: "verify:proof",
    label: "Verification proof is present",
    passed: Boolean(verification),
    points: 10
  });
  checks.push({
    id: "verify:build",
    label: "Build is verified",
    passed: verificationCheckPassed(verification, "go-build") || verificationCheckPassed(verification, "node-build"),
    points: 10
  });
  checks.push({
    id: "verify:cli",
    label: "CLI smoke test is verified",
    passed: verificationCheckPassed(verification, "cli-smoke"),
    points: 10
  });
  checks.push({
    id: "verify:mcp",
    label: "MCP handshake is verified",
    passed: verificationCheckPassed(verification, "mcp-handshake"),
    points: 10
  });
  checks.push({
    id: "verify:test",
    label: "Generated tests are verified",
    passed: verificationCheckPassed(verification, "go-test") || verificationCheckPassed(verification, "node-test"),
    points: 10
  });

  const maxScore = checks.reduce((sum, check) => sum + check.points, 0);
  const score = checks
    .filter((check) => check.passed)
    .reduce((sum, check) => sum + check.points, 0);

  const dimensions = computeDimensions(projectDir, manifest, isGo, checks);

  const ratio = maxScore > 0 ? score / maxScore : 0;
  return {
    projectDir,
    score,
    maxScore,
    ratio,
    grade: grade(ratio, verification),
    checks,
    dimensions,
    badges: scoreBadges(checks, verification, manifest),
    summary: manifest
      ? summarizeManifest(manifest)
      : "No manifest found."
  };
}

function computeDimensions(projectDir, manifest, isGo, checks) {
  const dims = {};
  dims.structure = scoreFromChecks(checks, (check) => check.id.startsWith("file:"));
  dims.manifest = scoreFromChecks(checks, (check) => check.id.startsWith("manifest:"));
  dims.runtime = scoreFromChecks(checks, (check) => check.id.startsWith("go:") || check.id.startsWith("package:"));
  dims.verification = scoreFromChecks(checks, (check) => check.id.startsWith("verify:"));

  if (manifest) {
    dims.coverage = scoreCoverage(manifest);
    dims.safety = scoreSafety(projectDir, manifest, isGo);
    dims.examples = scoreExamples(projectDir);
    dims.skill = scoreSkill(projectDir, manifest);
  }
  return dims;
}

function readVerification(projectDir) {
  const file = path.join(projectDir, "proofs", "verification.json");
  if (!fs.existsSync(file)) return null;
  try {
    return readJson(file);
  } catch {
    return null;
  }
}

function verificationCheckPassed(verification, id) {
  if (!verification || !Array.isArray(verification.checks)) return false;
  const check = verification.checks.find((item) => item.id === id);
  return check?.passed === true;
}

function scoreBadges(checks, verification, manifest) {
  const passed = (id) => checks.find((check) => check.id === id)?.passed === true;
  const badges = [];
  if (verification?.ok) badges.push("Gutenberg Verified");
  if (passed("verify:build")) badges.push("Build Verified");
  if (passed("verify:mcp")) badges.push("MCP Ready");
  const ops = manifest?.operations || [];
  if (ops.length > 0 && ops.every((op) => ["read", "write", "destructive"].includes(op.risk))) badges.push("Risk Gated");
  return badges;
}

function scoreFromChecks(checks, predicate) {
  const subset = checks.filter(predicate);
  if (subset.length === 0) return { score: 0, max: 0, ratio: 0 };
  const max = subset.reduce((sum, check) => sum + check.points, 0);
  const score = subset.filter((check) => check.passed).reduce((sum, check) => sum + check.points, 0);
  return { score, max, ratio: max > 0 ? score / max : 0, count: subset.length, passed: subset.filter((c) => c.passed).length };
}

function scoreCoverage(manifest) {
  if (manifest.kind === "aggregator") {
    const sources = manifest.sources || [];
    const max = 30;
    let score = 0;
    if (sources.length >= 1) score += 10;
    if (sources.length >= 2) score += 10;
    if (sources.every((source) => source.slug && source.operation)) score += 10;
    return { score, max, ratio: score / max, sources: sources.length };
  }
  const ops = manifest.operations || [];
  const tagged = ops.filter((op) => op.tags && op.tags.length > 0).length;
  const cacheable = ops.filter((op) => op.cacheable).length;
  const reads = ops.filter((op) => op.risk === "read").length;
  const max = 30;
  let score = 0;
  if (ops.length >= 1) score += 10;
  if (tagged / Math.max(ops.length, 1) >= 0.5) score += 5;
  if (cacheable >= 1) score += 5;
  if (reads >= 1) score += 5;
  if (ops.length >= 5) score += 5;
  return { score, max, ratio: score / max, operations: ops.length, cacheable, reads, tagged };
}

function scoreSafety(projectDir, manifest, isGo) {
  if (manifest.kind === "aggregator") {
    return { score: 20, max: 20, ratio: 1, note: "aggregator: fan-out only, no write ops" };
  }
  const ops = manifest.operations || [];
  const writes = ops.filter((op) => op.risk !== "read");
  const allClassified = ops.every((op) => ["read", "write", "destructive"].includes(op.risk));
  let guardSourcePresent = false;
  if (isGo) {
    const clientPath = path.join(projectDir, "internal", "forge", "client.go");
    if (fs.existsSync(clientPath)) {
      const source = fs.readFileSync(clientPath, "utf8");
      guardSourcePresent = /risk\s*!=\s*"read"|requires.*--yes|requireYes|writeGuard/i.test(source) || source.includes("yes");
    }
  } else {
    const clientPath = path.join(projectDir, "src", "client.js");
    if (fs.existsSync(clientPath)) {
      const source = fs.readFileSync(clientPath, "utf8");
      guardSourcePresent = source.includes('risk !== "read"') && source.includes("yes");
    }
  }
  const max = 20;
  let score = 0;
  if (allClassified) score += 10;
  if (writes.length === 0 || guardSourcePresent) score += 10;
  return { score, max, ratio: score / max, writes: writes.length, classified: allClassified, guardSourcePresent };
}

function scoreExamples(projectDir) {
  const cookbookPath = path.join(projectDir, "docs", "COOKBOOK.md");
  const max = 15;
  if (!fs.existsSync(cookbookPath)) return { score: 0, max, ratio: 0 };
  const content = fs.readFileSync(cookbookPath, "utf8");
  let score = 0;
  if (/```/.test(content)) score += 5;
  if (/sync|cache/i.test(content)) score += 5;
  if (/mcp/i.test(content)) score += 5;
  return { score, max, ratio: score / max };
}

function scoreSkill(projectDir, manifest) {
  const skillPath = path.join(projectDir, "skills", manifest.slug, "SKILL.md");
  const max = 10;
  if (!fs.existsSync(skillPath)) return { score: 0, max, ratio: 0, present: false };
  const content = fs.readFileSync(skillPath, "utf8");
  let score = 0;
  if (/^---\nname:/.test(content)) score += 4;
  if (/description:/.test(content)) score += 3;
  if (/Operations index/.test(content)) score += 3;
  return { score, max, ratio: score / max, present: true };
}

function summarizeManifest(manifest) {
  if (manifest.kind === "aggregator") {
    const sources = manifest.sources || [];
    return `${manifest.name}: aggregator with ${sources.length} source(s).`;
  }
  const ops = Array.isArray(manifest.operations) ? manifest.operations : [];
  const tags = Array.isArray(manifest.tags) ? manifest.tags : [];
  return `${manifest.name}: ${ops.length} operation(s), ${tags.length} domain group(s).`;
}

function grade(ratio, verification = null) {
  if (!verification?.ok && ratio >= 0.9) return "B";
  if (ratio >= 0.9) return "A";
  if (ratio >= 0.75) return "B";
  if (ratio >= 0.6) return "C";
  if (ratio >= 0.45) return "D";
  return "F";
}
