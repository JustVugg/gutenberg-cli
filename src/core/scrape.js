import { fetchHtml } from "./extract.js";

export async function scrapeMarkdown(url, options = {}) {
  const html = await fetchHtml(url);
  const main = extractMainContent(html);
  if (options.structured) {
    const tables = detectStructuredLists(main);
    if (tables.length > 0) {
      const md = htmlToMarkdown(main, { baseUrl: url, ...options });
      return tables.map((t) => renderTable(t)).join("\n\n") + "\n\n---\n\n" + md;
    }
  }
  return htmlToMarkdown(main, { baseUrl: url, ...options });
}

// detectStructuredLists scans the body for sibling clusters that share
// the same inline shape (e.g. <li>HH:MM - title</li> patterns) and emits
// them as table rows with inferred columns.
export function detectStructuredLists(html) {
  const tables = [];

  // Pattern 1: <ul>/<ol> with repeated <li> entries having a time prefix.
  const lists = [...html.matchAll(/<(ul|ol)[^>]*>([\s\S]*?)<\/\1>/gi)];
  for (const [, , inner] of lists) {
    const items = [...inner.matchAll(/<li[^>]*>([\s\S]*?)<\/li>/gi)].map((m) => stripTags(m[1]));
    if (items.length < 3) continue;
    const rows = items.map(parseTimeTitleLine).filter(Boolean);
    if (rows.length >= 3 && rows.length / items.length >= 0.6) {
      tables.push({ columns: ["time", "title"], rows });
    }
  }

  // Pattern 2: bare lines like "HH:MM - Title" inside text blocks.
  const stripped = stripTags(html);
  const bareRows = stripped
    .split(/\n+/)
    .map((line) => line.trim())
    .map(parseTimeTitleLine)
    .filter(Boolean);
  if (bareRows.length >= 3) {
    tables.push({ columns: ["time", "title"], rows: bareRows });
  }
  return dedupeTables(tables);
}

function parseTimeTitleLine(line) {
  if (!line) return null;
  const match = line.match(/^(\d{1,2}[:.]\d{2})\s*[-–—|]?\s*(.+)$/);
  if (!match) return null;
  const time = match[1].replace(".", ":");
  const title = match[2].trim();
  if (!title || title.length < 2) return null;
  return { time, title };
}

function dedupeTables(tables) {
  const seen = new Set();
  const out = [];
  for (const table of tables) {
    const key = table.rows.map((r) => `${r.time}|${r.title}`).join("\n");
    if (seen.has(key)) continue;
    seen.add(key);
    out.push(table);
  }
  return out;
}

function renderTable(table) {
  const headers = table.columns;
  const lines = [];
  lines.push(`| ${headers.join(" | ")} |`);
  lines.push(`| ${headers.map(() => "---").join(" | ")} |`);
  for (const row of table.rows.slice(0, 50)) {
    lines.push(`| ${headers.map((col) => String(row[col] || "").replace(/\|/g, "\\|")).join(" | ")} |`);
  }
  if (table.rows.length > 50) lines.push(`| … +${table.rows.length - 50} more rows … | |`);
  return lines.join("\n");
}

function stripTags(value) {
  return String(value)
    .replace(/<script[\s\S]*?<\/script>/gi, " ")
    .replace(/<style[\s\S]*?<\/style>/gi, " ")
    .replace(/<[^>]+>/g, " ")
    .replace(/&nbsp;/g, " ")
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&#39;/g, "'");
}

export function extractMainContent(html) {
  let body = html;
  body = body.replace(/<head[\s\S]*?<\/head>/gi, "");
  body = body.replace(/<script[\s\S]*?<\/script>/gi, "");
  body = body.replace(/<style[\s\S]*?<\/style>/gi, "");
  body = body.replace(/<noscript[\s\S]*?<\/noscript>/gi, "");
  body = body.replace(/<!--[\s\S]*?-->/g, "");

  const tagsToStrip = ["nav", "footer", "header", "aside", "form", "iframe", "svg", "video", "audio", "canvas"];
  for (const tag of tagsToStrip) {
    const open = new RegExp(`<${tag}[^>]*>[\\s\\S]*?<\\/${tag}>`, "gi");
    body = body.replace(open, "");
  }

  const candidates = [];
  for (const selector of ["main", "article", "section"]) {
    const re = new RegExp(`<${selector}[^>]*>([\\s\\S]*?)<\\/${selector}>`, "gi");
    let m;
    while ((m = re.exec(body)) !== null) {
      candidates.push(m[1]);
    }
  }
  candidates.sort((a, b) => b.length - a.length);
  if (candidates.length > 0 && candidates[0].length > 200) return candidates[0];

  const bodyMatch = body.match(/<body[^>]*>([\s\S]*?)<\/body>/i);
  return bodyMatch ? bodyMatch[1] : body;
}

export function htmlToMarkdown(html, options = {}) {
  let out = html;

  out = out.replace(/<h1[^>]*>([\s\S]*?)<\/h1>/gi, (_, inner) => `\n\n# ${cleanInline(inner)}\n\n`);
  out = out.replace(/<h2[^>]*>([\s\S]*?)<\/h2>/gi, (_, inner) => `\n\n## ${cleanInline(inner)}\n\n`);
  out = out.replace(/<h3[^>]*>([\s\S]*?)<\/h3>/gi, (_, inner) => `\n\n### ${cleanInline(inner)}\n\n`);
  out = out.replace(/<h4[^>]*>([\s\S]*?)<\/h4>/gi, (_, inner) => `\n\n#### ${cleanInline(inner)}\n\n`);
  out = out.replace(/<h5[^>]*>([\s\S]*?)<\/h5>/gi, (_, inner) => `\n\n##### ${cleanInline(inner)}\n\n`);
  out = out.replace(/<h6[^>]*>([\s\S]*?)<\/h6>/gi, (_, inner) => `\n\n###### ${cleanInline(inner)}\n\n`);

  out = out.replace(/<blockquote[^>]*>([\s\S]*?)<\/blockquote>/gi, (_, inner) => `\n\n> ${cleanInline(inner)}\n\n`);

  out = out.replace(/<pre[^>]*><code[^>]*>([\s\S]*?)<\/code><\/pre>/gi, (_, code) => `\n\n\`\`\`\n${decodeEntities(code)}\n\`\`\`\n\n`);
  out = out.replace(/<code[^>]*>([\s\S]*?)<\/code>/gi, (_, code) => `\`${cleanInline(code)}\``);

  out = out.replace(/<ul[^>]*>([\s\S]*?)<\/ul>/gi, (_, inner) => `\n\n${renderList(inner, "-")}\n\n`);
  out = out.replace(/<ol[^>]*>([\s\S]*?)<\/ol>/gi, (_, inner) => `\n\n${renderList(inner, "1.")}\n\n`);

  out = out.replace(/<a[^>]*href=["']([^"']+)["'][^>]*>([\s\S]*?)<\/a>/gi, (_, href, text) => {
    const cleanedText = cleanInline(text);
    const absoluteHref = absolutize(href, options.baseUrl);
    return `[${cleanedText}](${absoluteHref})`;
  });

  out = out.replace(/<strong[^>]*>([\s\S]*?)<\/strong>/gi, (_, inner) => `**${cleanInline(inner)}**`);
  out = out.replace(/<b[^>]*>([\s\S]*?)<\/b>/gi, (_, inner) => `**${cleanInline(inner)}**`);
  out = out.replace(/<em[^>]*>([\s\S]*?)<\/em>/gi, (_, inner) => `*${cleanInline(inner)}*`);
  out = out.replace(/<i[^>]*>([\s\S]*?)<\/i>/gi, (_, inner) => `*${cleanInline(inner)}*`);

  out = out.replace(/<br\s*\/?\s*>/gi, "\n");
  out = out.replace(/<p[^>]*>/gi, "\n\n");
  out = out.replace(/<\/p>/gi, "\n\n");

  out = out.replace(/<[^>]+>/g, "");
  out = decodeEntities(out);
  out = out.replace(/[ \t]+\n/g, "\n").replace(/\n{3,}/g, "\n\n").trim();
  return out;
}

function renderList(inner, marker) {
  const items = [...inner.matchAll(/<li[^>]*>([\s\S]*?)<\/li>/gi)].map((m) => m[1]);
  return items.map((item) => `${marker} ${cleanInline(item)}`).join("\n");
}

function cleanInline(value) {
  return decodeEntities(value.replace(/<[^>]+>/g, "")).replace(/\s+/g, " ").trim();
}

function absolutize(href, baseUrl) {
  if (!baseUrl) return href;
  try {
    return new URL(href, baseUrl).toString();
  } catch {
    return href;
  }
}

function decodeEntities(value) {
  return String(value)
    .replace(/&nbsp;/g, " ")
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/&#x27;/g, "'")
    .replace(/&#x2F;/g, "/");
}
