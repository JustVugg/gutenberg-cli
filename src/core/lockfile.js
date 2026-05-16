import fs from "node:fs";
import path from "node:path";
import { writeJson } from "./fs.js";

const LOCK_FILE = "gutenberg.lock.json";
const LOCK_SCHEMA = "gutenberg.lock.v1";

export function lockPath(projectDir) {
  return path.join(projectDir, LOCK_FILE);
}

export function readLock(projectDir) {
  const file = lockPath(projectDir);
  if (!fs.existsSync(file)) return null;
  try {
    return JSON.parse(fs.readFileSync(file, "utf8"));
  } catch {
    return null;
  }
}

export function writeLock(projectDir, provenance) {
  const lock = {
    schemaVersion: LOCK_SCHEMA,
    generatedAt: provenance.generatedAt,
    gutenbergVersion: provenance.gutenbergVersion,
    name: provenance.name || null,
    spec: provenance.spec || null,
    recipe: provenance.recipe || null,
    targets: provenance.targets || []
  };
  writeJson(lockPath(projectDir), lock);
  return lock;
}

export function compareLock(existingLock, currentProvenance) {
  if (!existingLock) return { hasLock: false, drifted: false, drifts: [] };
  const drifts = [];
  const lockedSpecHash = existingLock.spec?.sha256;
  const currentSpecHash = currentProvenance.spec?.sha256;
  if (lockedSpecHash && currentSpecHash && lockedSpecHash !== currentSpecHash) {
    drifts.push({
      kind: "spec",
      locked: lockedSpecHash.slice(0, 12),
      current: currentSpecHash.slice(0, 12),
      path: currentProvenance.spec?.path || existingLock.spec?.path
    });
  }
  if (existingLock.gutenbergVersion && currentProvenance.gutenbergVersion && existingLock.gutenbergVersion !== currentProvenance.gutenbergVersion) {
    drifts.push({
      kind: "gutenberg-version",
      locked: existingLock.gutenbergVersion,
      current: currentProvenance.gutenbergVersion
    });
  }
  return { hasLock: true, drifted: drifts.length > 0, drifts };
}
