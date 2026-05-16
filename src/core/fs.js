import fs from "node:fs";
import path from "node:path";

export function readText(filePath) {
  return fs.readFileSync(filePath, "utf8");
}

export function readJson(filePath) {
  return JSON.parse(readText(filePath));
}

export function writeText(filePath, content) {
  fs.mkdirSync(path.dirname(filePath), { recursive: true });
  fs.writeFileSync(filePath, content, "utf8");
}

export function writeJson(filePath, value) {
  writeText(filePath, `${JSON.stringify(value, null, 2)}\n`);
}

export function pathExists(filePath) {
  return fs.existsSync(filePath);
}

export function assertCanWriteDirectory(outDir, force = false) {
  if (!fs.existsSync(outDir)) {
    return;
  }
  const entries = fs.readdirSync(outDir);
  if (entries.length > 0 && !force) {
    const error = new Error(`output directory already exists and is not empty: ${outDir}. Use --force to overwrite generated Gutenberg files.`);
    error.exitCode = 2;
    throw error;
  }
}
