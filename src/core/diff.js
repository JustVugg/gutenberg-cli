import { buildBlueprint, loadOpenApi } from "./openapi.js";

export function diffSpecs(oldPath, newPath) {
  const oldBlueprint = buildBlueprint(loadOpenApi(oldPath), oldPath, "old");
  const newBlueprint = buildBlueprint(loadOpenApi(newPath), newPath, "new");

  const oldOps = new Map(oldBlueprint.operations.map((op) => [op.id, op]));
  const newOps = new Map(newBlueprint.operations.map((op) => [op.id, op]));

  const added = [];
  const removed = [];
  const changed = [];

  for (const [id, op] of newOps) {
    if (!oldOps.has(id)) {
      added.push(operationSummary(op));
    }
  }
  for (const [id, op] of oldOps) {
    if (!newOps.has(id)) {
      removed.push(operationSummary(op));
    }
  }
  for (const [id, newOp] of newOps) {
    const oldOp = oldOps.get(id);
    if (!oldOp) continue;
    const ops = compareOperation(oldOp, newOp);
    if (ops.length > 0) {
      changed.push({ id, summary: newOp.summary, changes: ops });
    }
  }

  return {
    schemaVersion: "gutenberg.diff.v1",
    old: { path: oldPath, operations: oldBlueprint.operations.length },
    new: { path: newPath, operations: newBlueprint.operations.length },
    counts: { added: added.length, removed: removed.length, changed: changed.length },
    added,
    removed,
    changed
  };
}

function operationSummary(op) {
  return {
    id: op.id,
    method: op.method,
    path: op.path,
    risk: op.risk,
    summary: op.summary || ""
  };
}

function compareOperation(oldOp, newOp) {
  const changes = [];
  if (oldOp.method !== newOp.method) changes.push({ kind: "method", from: oldOp.method, to: newOp.method });
  if (oldOp.path !== newOp.path) changes.push({ kind: "path", from: oldOp.path, to: newOp.path });
  if (oldOp.risk !== newOp.risk) changes.push({ kind: "risk", from: oldOp.risk, to: newOp.risk });

  const oldParams = parameterMap(oldOp.parameters);
  const newParams = parameterMap(newOp.parameters);
  for (const [key, param] of newParams) {
    if (!oldParams.has(key)) changes.push({ kind: "parameter-added", name: param.name, in: param.in, required: param.required });
  }
  for (const [key, param] of oldParams) {
    if (!newParams.has(key)) changes.push({ kind: "parameter-removed", name: param.name, in: param.in });
  }
  for (const [key, newParam] of newParams) {
    const oldParam = oldParams.get(key);
    if (!oldParam) continue;
    if (oldParam.required !== newParam.required) {
      changes.push({ kind: "parameter-required-changed", name: newParam.name, in: newParam.in, from: oldParam.required, to: newParam.required });
    }
  }
  return changes;
}

function parameterMap(parameters) {
  const map = new Map();
  for (const param of parameters || []) {
    map.set(`${param.in}:${param.name}`, param);
  }
  return map;
}
