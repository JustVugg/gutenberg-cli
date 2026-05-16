export function slugify(value) {
  return String(value || "")
    .trim()
    .toLowerCase()
    .replace(/['"]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "") || "tool";
}

export function camelCase(value) {
  const parts = slugify(value).split("-");
  return parts
    .map((part, index) => index === 0 ? part : part.charAt(0).toUpperCase() + part.slice(1))
    .join("");
}

export function pascalCase(value) {
  return slugify(value)
    .split("-")
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join("");
}

export function envPrefix(value) {
  return slugify(value).replace(/-/g, "_").toUpperCase();
}

const NOISE_SEGMENTS = new Set([
  "api", "apis", "rest", "public", "json", "service", "services",
  "v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9",
  "site", "www", "graph", "graphql"
]);

function isAllNoise(text) {
  if (!text) return true;
  const parts = text.toLowerCase().split(/[_\-]+/).filter(Boolean);
  if (parts.length === 0) return true;
  return parts.every((part) => NOISE_SEGMENTS.has(part) || /^v\d+(\.\d+)?$/.test(part));
}

export function safeOperationId(method, path) {
  const segmentsWithMeta = path.split("/").filter(Boolean).map((raw) => ({
    raw,
    text: raw.replace(/[{}]/g, ""),
    isParam: /^\{[^}]+\}$/.test(raw)
  }));

  const meaningful = segmentsWithMeta.filter(({ text }) => {
    if (isAllNoise(text)) return false;
    return true;
  });

  // Prefer the trailing run of non-parameter segments. If the path ends with
  // parameters (e.g. /feed/featured/{year}/{month}/{day}), keep the literal
  // segments before them so we get "feedFeatured" instead of "yearMonthDay".
  const lastNonParam = (() => {
    for (let i = meaningful.length - 1; i >= 0; i--) {
      if (!meaningful[i].isParam) return i;
    }
    return -1;
  })();

  let tail;
  if (lastNonParam >= 0) {
    const start = Math.max(0, lastNonParam - 2);
    tail = meaningful.slice(start, lastNonParam + 1);
  } else {
    tail = meaningful.slice(-3);
  }

  const pathName = tail.map(({ text }) => text).join("_") || "root";
  return camelCase(`${method}-${pathName}`);
}
