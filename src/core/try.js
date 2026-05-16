import { discoverOpenApi } from "./discover.js";
import { seedHarFromUrls, detectChallenge } from "./seed.js";

export async function tryUrl(url, options = {}) {
  const report = {
    url,
    timestamp: new Date().toISOString(),
    steps: [],
    verdict: null,
    confidence: "low",
    nextSteps: []
  };

  // Step 1: OpenAPI discovery.
  try {
    const discovery = await discoverOpenApi(url, { out: null });
    report.steps.push({ name: "discover", found: Boolean(discovery.found), url: discovery.url || null });
    if (discovery.found) {
      report.verdict = "openapi-published";
      report.confidence = "high";
      report.nextSteps = [
        `gutenberg generate ${discovery.url} --out library/tools/<name> --name <name> --force --strict`,
        `gutenberg install library/tools/<name>`
      ];
      return report;
    }
  } catch (error) {
    report.steps.push({ name: "discover", found: false, error: error.message });
  }

  // Step 2: plain seed-har.
  let seedHar;
  try {
    seedHar = await seedHarFromUrls([url], { out: null });
  } catch (error) {
    report.steps.push({ name: "seed-har", ok: false, error: error.message });
    report.verdict = "unreachable";
    report.nextSteps = [`Check the URL is correct and accessible: ${url}`];
    return report;
  }

  const entry = seedHar.log.entries[0];
  const status = entry.response.status;
  const text = entry.response.content.text || "";
  const mime = entry.response.content.mimeType || "";
  const challenge = detectChallenge(text);
  const isJson = /json/i.test(mime);
  const isHtml = /html/i.test(mime);

  report.steps.push({
    name: "seed-har",
    status,
    contentType: mime,
    bytes: text.length,
    challenge: challenge || null
  });

  // Step 3: classify.
  if (challenge) {
    report.verdict = "anti-bot-challenge";
    report.confidence = "high";
    report.nextSteps = [
      `Site is behind ${challenge.includes("Cloudflare") ? "Cloudflare" : "an anti-bot"} challenge. Headless tools won't get clean content.`,
      `Option A: gutenberg login ${url} --out state.json   (manual solve, headed browser)`,
      `Option B: gutenberg record ${url} --backend browserbase --key $BROWSERBASE_API_KEY`,
      `Option C: skip this site — use an official API/affiliate program instead.`
    ];
    return report;
  }
  if (status >= 400) {
    report.verdict = `http-${status}`;
    report.confidence = "high";
    report.nextSteps = [
      `Got HTTP ${status} ${entry.response.statusText}. The default browser-like UA was rejected.`,
      `Try: gutenberg record ${url} --out capture.har.json  (uses real Playwright)`
    ];
    return report;
  }
  if (isJson) {
    report.verdict = "json-endpoint";
    report.confidence = "high";
    report.nextSteps = [
      `Single JSON endpoint detected. Seed a few related endpoints, then generate:`,
      `  gutenberg seed-har ${url} <related-url-1> <related-url-2> --out /tmp/spec.har.json`,
      `  gutenberg import-har /tmp/spec.har.json --out /tmp/spec.openapi.json --name <name>`,
      `  gutenberg generate /tmp/spec.openapi.json --out library/tools/<name> --name <name> --force --strict`
    ];
    return report;
  }
  if (isHtml) {
    const spa = looksLikeSpa(text);
    if (spa) {
      report.verdict = "spa";
      report.confidence = "medium";
      report.steps.push({ name: "spa-detect", markers: spa });
      report.nextSteps = [
        `SPA detected (${spa.join(", ")}). The homepage HTML is mostly bootstrap; real data arrives via XHR.`,
        `Use a real browser to capture the XHR traffic:`,
        `  gutenberg record ${url} --out capture.har.json --wait 5000`,
        `Then filter the HAR for JSON entries and run import-har on a deep-link route URL.`
      ];
      return report;
    }
    // Low-content-density check: when the page is large HTML but visible text is tiny,
    // the page is rendered lazy/JS-side even without classic SPA markers.
    const density = lowDensityRatio(text);
    if (density && density.ratio < 0.05 && text.length > 50000) {
      report.verdict = "lazy-rendered";
      report.confidence = "medium";
      report.steps.push({ name: "density-check", ratio: density.ratio.toFixed(3), htmlBytes: text.length, textBytes: density.textBytes });
      report.nextSteps = [
        `Page is ${text.length} bytes of HTML but only ${density.textBytes} bytes are visible text (${(density.ratio * 100).toFixed(1)}% density).`,
        `Content is likely loaded lazily via JS. Try a headless browser capture:`,
        `  gutenberg record ${url} --out capture.har.json --wait 5000`,
        `If the site exposes a search/category URL with results inline, try that URL too.`
      ];
      return report;
    }
    report.verdict = "html-content";
    report.confidence = "high";
    report.steps.push({ name: "density-check", ratio: density ? density.ratio.toFixed(3) : null });
    report.nextSteps = [
      `Plain HTML content. Best path is scrape/extract:`,
      `  gutenberg scrape ${url}                              # narrative markdown`,
      `  gutenberg scrape ${url} --structured                 # tabular if patterns repeat`,
      `  gutenberg extract ${url} -p "<prompt>" -s schema.json --cache 1h`
    ];
    return report;
  }
  report.verdict = "unknown";
  report.nextSteps = [`Content type ${mime} not classified. Inspect the HAR manually.`];
  return report;
}

function lowDensityRatio(html) {
  const bodyMatch = html.match(/<body[^>]*>([\s\S]*?)<\/body>/i);
  const body = bodyMatch ? bodyMatch[1] : html;
  const text = body
    .replace(/<script[\s\S]*?<\/script>/gi, "")
    .replace(/<style[\s\S]*?<\/style>/gi, "")
    .replace(/<noscript[\s\S]*?<\/noscript>/gi, "")
    .replace(/<[^>]+>/g, " ")
    .replace(/\s+/g, " ")
    .trim();
  return { textBytes: text.length, ratio: text.length / Math.max(html.length, 1) };
}

function looksLikeSpa(html) {
  const markers = [];
  if (/<div id=["'](root|app|__next|__nuxt|svelte)/i.test(html)) markers.push("root mount element");
  if (/window\.__INITIAL_STATE__|window\.__NUXT__|window\.__NEXT_DATA__/i.test(html)) markers.push("hydration payload");
  if (/<script[^>]+src=[^>]*runtime[^>]*\.m?js/i.test(html)) markers.push("runtime bundle");
  if (/data-reactroot|data-reactid/i.test(html)) markers.push("React markers");
  if (/ng-version=|ng-app=/i.test(html)) markers.push("Angular markers");
  // SPA shells usually have <body> with very little text; the bulk is JS.
  const bodyMatch = html.match(/<body[^>]*>([\s\S]*?)<\/body>/i);
  if (bodyMatch) {
    const bodyText = bodyMatch[1].replace(/<script[\s\S]*?<\/script>/gi, "").replace(/<[^>]+>/g, " ").replace(/\s+/g, " ").trim();
    if (bodyText.length < 500) markers.push("body has <500 chars of text");
  }
  return markers.length >= 2 ? markers : null;
}

export function formatTryReport(report) {
  const lines = [];
  lines.push(`Gutenberg try report for ${report.url}`);
  lines.push(`Verdict: ${report.verdict} (${report.confidence} confidence)`);
  lines.push("");
  lines.push("Steps:");
  for (const step of report.steps) {
    const detail = Object.entries(step).filter(([key]) => key !== "name").map(([k, v]) => `${k}=${typeof v === "object" ? JSON.stringify(v) : v}`).join(" ");
    lines.push(`  - ${step.name}: ${detail}`);
  }
  lines.push("");
  lines.push("Next steps:");
  for (const next of report.nextSteps) lines.push(`  ${next}`);
  return lines.join("\n");
}
