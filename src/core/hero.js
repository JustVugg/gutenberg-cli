import fs from "node:fs";
import path from "node:path";
import { writeJson, writeText } from "./fs.js";
import { slugify } from "./sanitize.js";

export const SPORTS_LEAGUES = {
  nba: {
    id: "nba",
    label: "NBA",
    endpoint: "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard"
  },
  "serie-a": {
    id: "serie-a",
    label: "Serie A",
    endpoint: "https://site.api.espn.com/apis/site/v2/sports/soccer/ita.1/scoreboard"
  }
};

const LEAGUE_ALIASES = {
  basketball: "nba",
  "italian-serie-a": "serie-a",
  "ita.1": "serie-a",
  seriea: "serie-a",
  "serie_a": "serie-a",
  soccer: "serie-a",
  calcio: "serie-a"
};

const MONTHS = {
  jan: 0,
  january: 0,
  gennaio: 0,
  feb: 1,
  february: 1,
  febbraio: 1,
  mar: 2,
  march: 2,
  marzo: 2,
  apr: 3,
  april: 3,
  aprile: 3,
  may: 4,
  maggio: 4,
  jun: 5,
  june: 5,
  giugno: 5,
  jul: 6,
  july: 6,
  luglio: 6,
  aug: 7,
  august: 7,
  agosto: 7,
  sep: 8,
  sept: 8,
  september: 8,
  settembre: 8,
  oct: 9,
  october: 9,
  ottobre: 9,
  nov: 10,
  november: 10,
  novembre: 10,
  dec: 11,
  december: 11,
  dicembre: 11
};

export const RECIPES = [
  {
    id: "sports-nba-espn",
    kind: "sports",
    status: "ready",
    source: "ESPN public scoreboard JSON",
    command: "gutenberg sports nba today",
    creates: ["compact live schedule", "SQLite-ready generated HAR tool", "agent-readable event rows"],
    notes: "No API key required. Best for NBA schedule, scores, venues, broadcasts, tickets, and links."
  },
  {
    id: "sports-serie-a-espn",
    kind: "sports",
    status: "ready",
    source: "ESPN public Serie A scoreboard JSON",
    command: "gutenberg sports serie-a today",
    creates: ["compact live fixtures", "SQLite-ready generated HAR tool", "agent-readable match rows"],
    notes: "No API key required. Uses league slug ita.1 and preserves dots in imported HAR paths."
  },
  {
    id: "travel-amadeus-flight-offers",
    kind: "travel",
    status: "auth-required",
    source: "Amadeus Flight Offers Search API",
    command: "gutenberg travel rom-par june --adults 1 --currency EUR",
    creates: ["flight search plan", "live price results when AMADEUS_ACCESS_TOKEN is set", "agent summary rows"],
    notes: "Use official Amadeus OAuth credentials for real prices. Without a token it prints the exact request and setup."
  },
  {
    id: "travel-web-flight-search",
    kind: "travel",
    status: "browser-assisted",
    source: "Kayak/Google Flights-like browser search URLs",
    command: "gutenberg travel rom-par june --web",
    creates: ["Kayak-style URL", "Google Flights query URL", "record command for HAR capture"],
    notes: "Does not bypass anti-bot, login, paywalls, or terms. It gives reproducible browser-recording entry points."
  },
  {
    id: "generic-openapi",
    kind: "factory",
    status: "ready",
    source: "OpenAPI JSON/YAML",
    command: "gutenberg generate openapi.json --out generated/my-tool --name my-tool --force",
    creates: ["Go CLI", "MCP server", "SQLite/FTS5 cache", "OAuth helper", "scorecard"],
    notes: "Best path when the service publishes OpenAPI or Swagger."
  },
  {
    id: "generic-graphql",
    kind: "factory",
    status: "ready",
    source: "GraphQL SDL, introspection JSON, or endpoint",
    command: "gutenberg import-graphql schema.json --out generated/api.openapi.json --name api",
    creates: ["OpenAPI bridge", "Go CLI", "MCP server"],
    notes: "Best path for GraphQL APIs and agent-friendly query wrappers."
  },
  {
    id: "generic-browser-har",
    kind: "factory",
    status: "ready",
    source: "Website or web app with no public spec",
    command: "gutenberg record https://example.com --out capture.har.json",
    creates: ["HAR capture", "OpenAPI bridge", "Go CLI/MCP package"],
    notes: "Best path for sites with visible network APIs. It records; it does not defeat access controls."
  }
];

export async function runSports(positionals, options = {}) {
  if ((positionals[0] || "").toLowerCase() === "list") {
    printSportsList();
    return;
  }
  const league = resolveLeague(positionals[0] || options.league || "nba");
  const date = resolveDateToken(options.date || positionals[1] || "today");
  const url = new URL(league.endpoint);
  if (date) url.searchParams.set("dates", date.replaceAll("-", ""));
  const payload = await fetchJson(url);
  const summary = summarizeScoreboard(payload, league, {
    date,
    team: options.team,
    limit: Number(options.limit || 12)
  });
  if (options.json) {
    console.log(JSON.stringify(summary, null, 2));
    return;
  }
  printSportsSummary(summary);
}

export async function runTravel(positionals, options = {}) {
  const plan = planTravelSearch(positionals, options);
  const token = options.token || process.env.AMADEUS_ACCESS_TOKEN || process.env.AMADEUS_FLIGHT_OFFERS_ACCESS_TOKEN;
  if (options.json && (options.plan || options.web || !token)) {
    console.log(JSON.stringify({ status: token ? "ready" : "needs_auth", plan }, null, 2));
    return;
  }
  if (options.web || !token) {
    printTravelPlan(plan, token ? "web" : "needs_auth");
    return;
  }
  const payload = await fetchJson(plan.amadeusUrl, {
    headers: {
      Accept: "application/json",
      Authorization: `Bearer ${token}`
    }
  });
  const summary = summarizeAmadeusOffers(payload, plan);
  if (options.json) {
    console.log(JSON.stringify(summary, null, 2));
    return;
  }
  printTravelSummary(summary);
}

export async function runRecipes(positionals, options = {}, rootDir = process.cwd()) {
  const action = positionals[0] || "list";
  if (action === "list") {
    if (options.json) {
      console.log(JSON.stringify({ recipes: RECIPES }, null, 2));
      return;
    }
    for (const recipe of RECIPES) {
      console.log(`${recipe.id.padEnd(30)} ${recipe.status.padEnd(16)} ${recipe.command}`);
    }
    return;
  }
  if (action === "show") {
    const recipe = findRecipe(positionals[1]);
    console.log(JSON.stringify(recipe, null, 2));
    return;
  }
  if (action === "run") {
    const recipe = findRecipe(positionals[1]);
    const rest = positionals.slice(2);
    if (recipe.id === "sports-nba-espn") return runSports(["nba", ...rest], options);
    if (recipe.id === "sports-serie-a-espn") return runSports(["serie-a", ...rest], options);
    if (recipe.id === "travel-amadeus-flight-offers" || recipe.id === "travel-web-flight-search") return runTravel(rest, options);
    throw new Error(`recipe ${recipe.id} is a factory recipe; use its command: ${recipe.command}`);
  }
  if (action === "scaffold") {
    const name = positionals[1];
    if (!name) throw new Error("missing recipe name. Example: gutenberg recipes scaffold my-site --kind har");
    const recipePath = scaffoldRecipe(rootDir, name, options);
    console.log(`Wrote recipe: ${recipePath}`);
    return;
  }
  throw new Error(`unknown recipes action: ${action}`);
}

export function runCreate(positionals, options = {}, rootDir = process.cwd()) {
  const name = positionals[0];
  if (!name) throw new Error("missing project name. Example: gutenberg create acme --from https://api.example.com/openapi.json");
  const slug = slugify(name);
  const outDir = path.resolve(options.out || path.join(rootDir, "created", slug));
  fs.mkdirSync(outDir, { recursive: true });
  const source = options.from || options.spec || options.url || null;
  const kind = options.kind || inferCreateKind(source);
  const project = {
    schemaVersion: "gutenberg.create.v1",
    name,
    slug,
    kind,
    source,
    goal: options.goal || "Create an agent-native CLI/MCP tool with compact hero commands.",
    nextCommands: nextCommandsForCreate(slug, source, kind),
    createdAt: new Date().toISOString()
  };
  writeJson(path.join(outDir, "gutenberg.create.json"), project);
  writeText(path.join(outDir, "README.md"), createReadme(project));
  const recipePath = scaffoldRecipe(rootDir, slug, { kind, source, out: path.join(outDir, `${slug}.recipe.json`) });
  console.log(`Created ${name} at ${outDir}`);
  console.log(`Recipe: ${recipePath}`);
  for (const command of project.nextCommands) console.log(`Next: ${command}`);
}

export function planTravelSearch(positionals = [], options = {}, now = new Date()) {
  const routeToken = options.route || positionals.find((item) => /^[a-z]{3}\s*[-:>]\s*[a-z]{3}$/i.test(item));
  const origin = (options.origin || parseRoute(routeToken)?.origin || "").toUpperCase();
  const destination = (options.destination || options.dest || parseRoute(routeToken)?.destination || "").toUpperCase();
  if (!origin || !destination) {
    throw new Error("missing route. Example: gutenberg travel rom-par june --date 2026-06-15");
  }
  const monthToken = positionals.find((item) => MONTHS[item.toLowerCase()] !== undefined);
  const dateToken = positionals.find((item) => /^\d{4}-\d{2}-\d{2}$/.test(item));
  const departureDate = options.date || dateToken || inferDepartureDate(monthToken, now);
  const returnDate = options.return || options["return-date"] || null;
  const query = new URLSearchParams({
    originLocationCode: origin,
    destinationLocationCode: destination,
    departureDate,
    adults: String(options.adults || 1),
    currencyCode: options.currency || "EUR",
    max: String(options.max || 5)
  });
  if (returnDate) query.set("returnDate", returnDate);
  if (options.nonstop || options["non-stop"]) query.set("nonStop", "true");
  if (options["max-price"]) query.set("maxPrice", String(options["max-price"]));
  const baseUrl = options.prod ? "https://api.amadeus.com/v2" : "https://test.api.amadeus.com/v2";
  const kayakPath = `${origin}-${destination}/${departureDate}${returnDate ? `/${returnDate}` : ""}`;
  return {
    origin,
    destination,
    departureDate,
    returnDate,
    adults: Number(options.adults || 1),
    currency: options.currency || "EUR",
    max: Number(options.max || 5),
    nonStop: Boolean(options.nonstop || options["non-stop"]),
    amadeusUrl: `${baseUrl}/shopping/flight-offers?${query.toString()}`,
    webUrls: {
      kayak: `https://www.kayak.com/flights/${kayakPath}?sort=bestflight_a`,
      googleFlights: `https://www.google.com/travel/flights?q=${encodeURIComponent(`Flights from ${origin} to ${destination} on ${departureDate}${returnDate ? ` returning ${returnDate}` : ""}`)}`
    },
    recordCommand: `node bin/gutenberg.js record "https://www.kayak.com/flights/${kayakPath}?sort=bestflight_a" --out generated/${origin.toLowerCase()}-${destination.toLowerCase()}-${departureDate}.har.json`
  };
}

export function summarizeScoreboard(payload, league, options = {}) {
  const teamNeedle = options.team ? String(options.team).toLowerCase() : "";
  const events = (payload.events || [])
    .map((event) => summarizeEvent(event))
    .filter((event) => !teamNeedle || JSON.stringify(event).toLowerCase().includes(teamNeedle))
    .slice(0, options.limit || 12);
  return {
    league: league.label,
    date: payload.day?.date || options.date || null,
    count: events.length,
    events
  };
}

export function summarizeAmadeusOffers(payload, plan) {
  const offers = (payload.data || []).slice(0, plan.max || 5).map((offer, index) => {
    const firstItinerary = offer.itineraries?.[0];
    const firstSegment = firstItinerary?.segments?.[0];
    const lastSegment = firstItinerary?.segments?.[firstItinerary.segments.length - 1];
    return {
      rank: index + 1,
      total: offer.price?.grandTotal || offer.price?.total || null,
      currency: offer.price?.currency || plan.currency,
      carrier: firstSegment?.carrierCode || null,
      departure: firstSegment?.departure?.at || null,
      arrival: lastSegment?.arrival?.at || null,
      duration: firstItinerary?.duration || null,
      stops: Math.max(0, (firstItinerary?.segments?.length || 1) - 1)
    };
  });
  return {
    route: `${plan.origin}-${plan.destination}`,
    departureDate: plan.departureDate,
    returnDate: plan.returnDate,
    count: offers.length,
    offers,
    warnings: payload.errors || []
  };
}

function resolveLeague(value) {
  const key = String(value || "nba").toLowerCase();
  const normalized = LEAGUE_ALIASES[key] || key;
  const league = SPORTS_LEAGUES[normalized];
  if (!league) throw new Error(`unknown sports league: ${value}. Try nba or serie-a.`);
  return league;
}

function resolveDateToken(token) {
  if (!token || token === "today") return ymd(new Date());
  if (token === "tomorrow") {
    const date = new Date();
    date.setDate(date.getDate() + 1);
    return ymd(date);
  }
  if (/^\d{4}-\d{2}-\d{2}$/.test(token)) return token;
  return ymd(new Date());
}

async function fetchJson(url, init = {}) {
  const response = await fetch(url, init);
  const text = await response.text();
  let data;
  try {
    data = text ? JSON.parse(text) : null;
  } catch {
    data = text;
  }
  if (!response.ok && !data?.errors) {
    return { errors: [{ status: response.status, title: response.statusText, detail: data }] };
  }
  return data;
}

function summarizeEvent(event) {
  const competition = event.competitions?.[0] || {};
  const competitors = competition.competitors || [];
  const home = competitors.find((item) => item.homeAway === "home") || competitors[0] || {};
  const away = competitors.find((item) => item.homeAway === "away") || competitors[1] || {};
  const status = event.status?.type || competition.status?.type || {};
  const tickets = competition.tickets?.[0]?.summary || event.tickets?.[0]?.summary || null;
  return {
    id: event.id,
    name: event.name,
    shortName: event.shortName,
    date: event.date || competition.date,
    status: status.shortDetail || status.detail || status.description || status.name || "unknown",
    state: status.state || null,
    completed: Boolean(status.completed),
    home: teamSummary(home),
    away: teamSummary(away),
    venue: competition.venue?.fullName || event.venue?.displayName || null,
    broadcast: competition.broadcast || competition.broadcasts?.[0]?.names?.join(", ") || null,
    tickets,
    link: event.links?.[0]?.href || null
  };
}

function teamSummary(competitor) {
  return {
    id: competitor.id || null,
    name: competitor.team?.shortDisplayName || competitor.team?.displayName || competitor.team?.name || null,
    abbreviation: competitor.team?.abbreviation || null,
    score: competitor.score || null,
    record: competitor.record || competitor.records?.[0]?.summary || null,
    winner: Boolean(competitor.winner)
  };
}

function printSportsSummary(summary) {
  console.log(`${summary.league} - ${summary.date || "live"}`);
  if (summary.count === 0) {
    console.log("No events found.");
    return;
  }
  for (const event of summary.events) {
    const teams = event.away?.name && event.home?.name ? `${event.away.name} @ ${event.home.name}` : event.shortName || event.name;
    const score = event.away?.score || event.home?.score ? ` | ${event.away.abbreviation || "AWAY"} ${event.away.score || 0}-${event.home.score || 0} ${event.home.abbreviation || "HOME"}` : "";
    const extras = [event.broadcast, event.tickets, event.venue].filter(Boolean).join(" | ");
    console.log(`- ${formatDateTime(event.date)} | ${teams} | ${event.status}${score}${extras ? ` | ${extras}` : ""}`);
  }
}

function printTravelPlan(plan, status) {
  console.log(`Travel ${plan.origin}-${plan.destination} ${plan.departureDate}${plan.returnDate ? ` to ${plan.returnDate}` : ""}`);
  console.log(status === "needs_auth" ? "Status: needs AMADEUS_ACCESS_TOKEN for live prices." : "Status: browser-assisted search plan.");
  console.log(`Amadeus: ${plan.amadeusUrl}`);
  console.log(`Kayak-like: ${plan.webUrls.kayak}`);
  console.log(`Google Flights-like: ${plan.webUrls.googleFlights}`);
  console.log(`Record: ${plan.recordCommand}`);
}

function printTravelSummary(summary) {
  console.log(`Flights ${summary.route} ${summary.departureDate}${summary.returnDate ? ` to ${summary.returnDate}` : ""}`);
  if (summary.warnings.length > 0) {
    for (const warning of summary.warnings) console.log(`- ${warning.status || "WARN"} ${warning.title || warning.message || warning.detail}`);
    return;
  }
  if (summary.count === 0) {
    console.log("No offers returned.");
    return;
  }
  for (const offer of summary.offers) {
    console.log(`- #${offer.rank} ${offer.total || "n/a"} ${offer.currency || ""} | ${offer.carrier || "carrier n/a"} | ${offer.departure || "departure n/a"} -> ${offer.arrival || "arrival n/a"} | ${offer.duration || "duration n/a"} | stops ${offer.stops}`);
  }
}

function printSportsList() {
  for (const league of Object.values(SPORTS_LEAGUES)) {
    console.log(`${league.id.padEnd(10)} ${league.label}`);
  }
}

function parseRoute(value) {
  if (!value) return null;
  const match = String(value).match(/^([a-z]{3})\s*[-:>]\s*([a-z]{3})$/i);
  if (!match) return null;
  return { origin: match[1], destination: match[2] };
}

function inferDepartureDate(monthToken, now) {
  if (!monthToken) {
    const date = new Date(now);
    date.setDate(date.getDate() + 30);
    return ymd(date);
  }
  const month = MONTHS[monthToken.toLowerCase()];
  const year = month < now.getMonth() ? now.getFullYear() + 1 : now.getFullYear();
  return `${year}-${pad2(month + 1)}-15`;
}

function ymd(date) {
  return `${date.getFullYear()}-${pad2(date.getMonth() + 1)}-${pad2(date.getDate())}`;
}

function pad2(value) {
  return String(value).padStart(2, "0");
}

function formatDateTime(value) {
  if (!value) return "time n/a";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("it-IT", {
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    timeZoneName: "short"
  }).format(date);
}

function findRecipe(id) {
  const recipe = RECIPES.find((item) => item.id === id);
  if (!recipe) throw new Error(`unknown recipe: ${id}`);
  return recipe;
}

function scaffoldRecipe(rootDir, name, options = {}) {
  const slug = slugify(name);
  const recipePath = path.resolve(options.out || path.join(rootDir, "recipes", `${slug}.recipe.json`));
  const recipe = {
    schemaVersion: "gutenberg.recipe.v1",
    id: slug,
    kind: options.kind || "har",
    source: options.source || null,
    goal: options.goal || "Turn this source into an agent-native CLI/MCP tool.",
    inputs: {
      openapi: options.kind === "openapi" ? options.source || "openapi.json" : null,
      graphql: options.kind === "graphql" ? options.source || "schema.graphql" : null,
      url: options.kind === "har" || options.kind === "browser" ? options.source || "https://example.com" : null
    },
    output: {
      name: slug,
      directory: `generated/${slug}-go`,
      compactCommands: true,
      sqliteFts5: true,
      mcp: true
    },
    nextCommands: nextCommandsForCreate(slug, options.source, options.kind || "har")
  };
  writeJson(recipePath, recipe);
  return recipePath;
}

function inferCreateKind(source) {
  if (!source) return "har";
  if (/graphql/i.test(source)) return "graphql";
  if (/openapi|swagger|\.ya?ml$|\.json$/i.test(source)) return "openapi";
  if (/^https?:\/\//i.test(source)) return "har";
  return "openapi";
}

function nextCommandsForCreate(slug, source, kind) {
  if (kind === "graphql") {
    return [
      `node bin/gutenberg.js import-graphql ${source || "schema.graphql"} --out generated/${slug}.openapi.json --name ${slug}`,
      `node bin/gutenberg.js generate generated/${slug}.openapi.json --out generated/${slug}-go --name ${slug} --force`
    ];
  }
  if (kind === "openapi") {
    return [
      source && /^https?:\/\//i.test(source)
        ? `node bin/gutenberg.js discover ${source} --out generated/${slug}.openapi.json`
        : `node bin/gutenberg.js analyze ${source || "openapi.json"}`,
      `node bin/gutenberg.js generate ${source && !/^https?:\/\//i.test(source) ? source : `generated/${slug}.openapi.json`} --out generated/${slug}-go --name ${slug} --force`
    ];
  }
  return [
    `node bin/gutenberg.js record ${source || "https://example.com"} --out generated/${slug}.har.json`,
    `node bin/gutenberg.js import-har generated/${slug}.har.json --out generated/${slug}.openapi.json --name ${slug}`,
    `node bin/gutenberg.js generate generated/${slug}.openapi.json --out generated/${slug}-go --name ${slug} --force`
  ];
}

function createReadme(project) {
  return `# ${project.name}

Generated by Gutenberg create.

Kind: ${project.kind}
Source: ${project.source || "not set yet"}

## Next Commands

${project.nextCommands.map((command) => `\`\`\`bash\n${command}\n\`\`\``).join("\n\n")}
`;
}
