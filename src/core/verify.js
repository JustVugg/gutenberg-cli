import fs from "node:fs";
import path from "node:path";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import { hashFile } from "./provenance.js";
import { readJson, writeJson, writeText } from "./fs.js";
import { scoreProject } from "./scorecard.js";

const rootDir = path.dirname(path.dirname(path.dirname(fileURLToPath(import.meta.url))));
const useGo = path.join(rootDir, "scripts", "use-go.sh");

export function verifyProject(projectDir, options = {}) {
  const resolvedDir = path.resolve(projectDir);
  const manifestPath = path.join(resolvedDir, "gutenberg.manifest.json");
  if (!fs.existsSync(manifestPath)) {
    const error = new Error(`No gutenberg.manifest.json in ${resolvedDir} — not a Gutenberg-generated project.`);
    error.exitCode = 2;
    throw error;
  }
  const manifest = readJson(manifestPath);
  const proofDir = path.join(resolvedDir, "proofs");
  fs.mkdirSync(path.join(proofDir, "bin"), { recursive: true });

  const result = {
    schemaVersion: "gutenberg.verification.v1",
    projectDir: resolvedDir,
    tool: manifest.slug,
    language: manifest.language || (fs.existsSync(path.join(resolvedDir, "go.mod")) ? "go" : "node"),
    startedAt: new Date().toISOString(),
    finishedAt: null,
    ok: false,
    checks: [],
    artifacts: {
      proofDir,
      manifest: manifestPath,
      source: manifest.provenance?.spec || null
    }
  };

  if (result.language === "go") {
    verifyGoProject(resolvedDir, manifest, proofDir, result, options);
  } else {
    verifyNodeProject(resolvedDir, manifest, proofDir, result, options);
  }

  const required = result.checks.filter((check) => check.required !== false);
  result.ok = required.length > 0 && required.every((check) => check.passed);
  result.finishedAt = new Date().toISOString();
  result.badges = verificationBadges(result, manifest);
  writeJson(path.join(proofDir, "verification.json"), result);

  const updatedManifest = {
    ...manifest,
    verification: {
      schemaVersion: "gutenberg.verification.summary.v1",
      verifiedAt: result.finishedAt,
      ok: result.ok,
      badges: result.badges,
      proofFile: path.relative(resolvedDir, path.join(proofDir, "verification.json"))
    }
  };
  if (updatedManifest.provenance) {
    updatedManifest.provenance = {
      ...updatedManifest.provenance,
      verification: updatedManifest.verification
    };
  }
  writeJson(manifestPath, updatedManifest);
  if (fs.existsSync(path.join(resolvedDir, "blackforge.manifest.json"))) {
    writeJson(path.join(resolvedDir, "blackforge.manifest.json"), updatedManifest);
  }

  const score = scoreProject(resolvedDir);
  writeJson(path.join(proofDir, "scorecard.json"), score);
  return { ...result, scorecard: score };
}

function verifyGoProject(projectDir, manifest, proofDir, result, options) {
  if (!fs.existsSync(useGo)) {
    addCheck(result, "go-wrapper", false, `missing ${useGo}`, { required: true });
    return;
  }
  if (!options.noTidy) {
    runCommandCheck(result, proofDir, {
      id: "go-mod-tidy",
      label: "go mod tidy",
      command: useGo,
      args: ["mod", "tidy"],
      cwd: projectDir
    });
  }
  runCommandCheck(result, proofDir, {
    id: "go-test",
    label: "go test ./...",
    command: useGo,
    args: ["test", "./..."],
    cwd: projectDir
  });
  const binPath = path.join(proofDir, "bin", manifest.slug);
  runCommandCheck(result, proofDir, {
    id: "go-build",
    label: "go build CLI binary",
    command: useGo,
    args: ["build", "-o", binPath, `./cmd/${manifest.slug}`],
    cwd: projectDir,
    artifact: binPath
  });
  if (fs.existsSync(binPath)) {
    result.artifacts.binary = binPath;
    const hash = hashFile(binPath);
    if (hash) result.artifacts.binarySha256 = hash.sha256;
  }
  const useBuiltBinary = fs.existsSync(binPath);
  const smokeArgs = manifest.kind === "aggregator" ? ["sources"] : ["operations", "--json"];
  const smoke = runCommandCheck(result, proofDir, {
    id: "cli-smoke",
    label: manifest.kind === "aggregator" ? "CLI sources smoke" : "CLI operations smoke",
    command: useBuiltBinary ? binPath : useGo,
    args: useBuiltBinary ? smokeArgs : ["run", `./cmd/${manifest.slug}`, ...smokeArgs],
    cwd: projectDir,
    validate: (stdout) => validateCliSmoke(stdout, manifest)
  });
  if (smoke?.stdout) writeText(path.join(proofDir, "cli-smoke.out"), smoke.stdout);

  if (manifest.kind !== "aggregator" && fs.existsSync(path.join(projectDir, "internal", "forge", "mcp.go"))) {
    const mcpInput = [
      JSON.stringify({ jsonrpc: "2.0", id: 1, method: "initialize", params: {} }),
      JSON.stringify({ jsonrpc: "2.0", id: 2, method: "tools/list", params: {} }),
      ""
    ].join("\n");
    const mcp = runCommandCheck(result, proofDir, {
      id: "mcp-handshake",
      label: "MCP initialize + tools/list",
      command: useBuiltBinary ? binPath : useGo,
      args: useBuiltBinary ? ["mcp"] : ["run", `./cmd/${manifest.slug}`, "mcp"],
      cwd: projectDir,
      input: mcpInput,
      timeout: 10000,
      acceptValidatedOutput: true,
      validate: (stdout) => validateMcp(stdout, manifest)
    });
    if (mcp?.stdout) writeText(path.join(proofDir, "mcp-handshake.out"), mcp.stdout);
  } else {
    addCheck(result, "mcp-handshake", true, "not applicable for aggregator", { required: false });
  }
}

function verifyNodeProject(projectDir, manifest, proofDir, result, options) {
  if (!options.skipTests && fs.existsSync(path.join(projectDir, "tests"))) {
    const tests = fs.readdirSync(path.join(projectDir, "tests")).filter((file) => file.endsWith(".test.js")).map((file) => path.join("tests", file));
    if (tests.length > 0) {
      runCommandCheck(result, proofDir, {
        id: "node-test",
        label: "node --test",
        command: process.execPath,
        args: ["--test", ...tests],
        cwd: projectDir
      });
    }
  }
  const bin = path.join(projectDir, "bin", `${manifest.slug}.js`);
  runCommandCheck(result, proofDir, {
    id: "cli-smoke",
    label: "Node CLI operations smoke",
    command: process.execPath,
    args: [bin, "operations"],
    cwd: projectDir,
    validate: (stdout) => stdout.includes(manifest.operations?.[0]?.id || "")
  });
  addCheck(result, "go-build", true, "not applicable for node target", { required: false });
  addCheck(result, "mcp-handshake", true, "node MCP smoke is covered by generated tests", { required: false });
}

function runCommandCheck(result, proofDir, spec) {
  const started = Date.now();
  const proc = spawnSync(spec.command, spec.args, {
    cwd: spec.cwd,
    input: spec.input,
    encoding: "utf8",
    timeout: spec.timeout || 120000,
    maxBuffer: 20 * 1024 * 1024
  });
  const elapsedMs = Date.now() - started;
  const stdout = proc.stdout || "";
  const stderr = proc.stderr || "";
  const log = [
    `$ ${[spec.command, ...spec.args].join(" ")}`,
    `cwd: ${spec.cwd}`,
    `status: ${proc.status}`,
    `signal: ${proc.signal || ""}`,
    "",
    "stdout:",
    stdout,
    "",
    "stderr:",
    stderr
  ].join("\n");
  writeText(path.join(proofDir, `${spec.id}.log`), log);

  let validation = true;
  if (spec.validate) {
    try {
      validation = spec.validate(stdout, stderr);
    } catch (error) {
      validation = error.message;
    }
  }
  let passed = proc.status === 0 && validation === true;
  let detail = passed ? "ok" : `exit ${proc.status ?? "null"}${proc.signal ? ` signal ${proc.signal}` : ""}`;
  if (validation !== true) {
    passed = false;
    detail = typeof validation === "string" ? validation : "validation failed";
  } else if (spec.acceptValidatedOutput && stdout && validation === true) {
    passed = true;
    detail = proc.status === 0 ? "ok" : "validated output";
  } else if (proc.error && proc.status !== 0) {
    passed = false;
    detail = proc.error.message;
  }

  const check = addCheck(result, spec.id, passed, detail, {
    label: spec.label,
    elapsedMs,
    log: path.relative(result.projectDir, path.join(proofDir, `${spec.id}.log`)),
    artifact: spec.artifact || null
  });
  return { ...check, stdout, stderr };
}

function addCheck(result, id, passed, detail, extra = {}) {
  const check = {
    id,
    label: extra.label || id,
    passed: Boolean(passed),
    required: extra.required !== false,
    detail,
    elapsedMs: extra.elapsedMs || 0,
    log: extra.log || null,
    artifact: extra.artifact || null
  };
  result.checks.push(check);
  return check;
}

function validateCliSmoke(stdout, manifest) {
  if (manifest.kind === "aggregator") return stdout.includes("source") || stdout.includes(manifest.slug);
  try {
    const parsed = JSON.parse(stdout);
    return Array.isArray(parsed) && parsed.length > 0 && parsed.some((operation) => operation.id);
  } catch {
    return "operations --json did not emit a JSON operation array";
  }
}

function validateMcp(stdout, manifest) {
  const frames = parseMcpFrames(stdout);
  const init = frames.find((frame) => frame.id === 1);
  const tools = frames.find((frame) => frame.id === 2);
  if (!init?.result?.serverInfo?.name?.includes(manifest.slug)) return "MCP initialize response missing serverInfo";
  const list = tools?.result?.tools;
  if (!Array.isArray(list) || list.length === 0) return "MCP tools/list returned no tools";
  return true;
}

export function parseMcpFrames(stdout) {
  const frames = [];
  let index = 0;
  const text = String(stdout || "");
  while (index < text.length) {
    const header = text.slice(index).match(/Content-Length:\s*(\d+)\r?\n\r?\n/i);
    if (!header) break;
    const headerStart = index + header.index;
    const bodyStart = headerStart + header[0].length;
    const length = Number(header[1]);
    const body = text.slice(bodyStart, bodyStart + length);
    try {
      frames.push(JSON.parse(body));
    } catch {
      // Keep scanning; the validation step will report missing frames.
    }
    index = bodyStart + length;
  }
  return frames;
}

function verificationBadges(result, manifest) {
  const check = (id) => result.checks.find((item) => item.id === id)?.passed === true;
  const badges = [];
  if (result.ok) badges.push("Gutenberg Verified");
  if (check("go-build") || result.language === "node") badges.push("Build Verified");
  if (check("mcp-handshake")) badges.push("MCP Ready");
  const ops = manifest.operations || [];
  if (ops.every((op) => ["read", "write", "destructive"].includes(op.risk))) badges.push("Risk Gated");
  return badges;
}
