import assert from "node:assert/strict";
import test from "node:test";
import { planTravelSearch, RECIPES, SPORTS_LEAGUES, summarizeAmadeusOffers, summarizeScoreboard } from "../src/core/hero.js";

test("plans compact travel searches from route and month", () => {
  const plan = planTravelSearch(["rom-par", "june"], { adults: 2, currency: "EUR", max: 3 }, new Date("2026-05-09T10:00:00Z"));
  assert.equal(plan.origin, "ROM");
  assert.equal(plan.destination, "PAR");
  assert.equal(plan.departureDate, "2026-06-15");
  assert.equal(plan.adults, 2);
  assert.match(plan.amadeusUrl, /originLocationCode=ROM/);
  assert.match(plan.webUrls.kayak, /ROM-PAR/);
});

test("summarizes ESPN scoreboard payloads into agent rows", () => {
  const summary = summarizeScoreboard({
    day: { date: "2026-05-09" },
    events: [{
      id: "1",
      name: "Oklahoma City Thunder at Los Angeles Lakers",
      shortName: "OKC @ LAL",
      date: "2026-05-10T00:30Z",
      status: { type: { shortDetail: "8:30 PM EDT", state: "pre" } },
      competitions: [{
        venue: { fullName: "crypto.com Arena" },
        competitors: [
          { homeAway: "away", score: "0", team: { abbreviation: "OKC", shortDisplayName: "Thunder" } },
          { homeAway: "home", score: "0", team: { abbreviation: "LAL", shortDisplayName: "Lakers" } }
        ],
        tickets: [{ summary: "Tickets as low as $25" }]
      }]
    }]
  }, SPORTS_LEAGUES.nba);

  assert.equal(summary.count, 1);
  assert.equal(summary.events[0].away.name, "Thunder");
  assert.equal(summary.events[0].home.name, "Lakers");
  assert.equal(summary.events[0].tickets, "Tickets as low as $25");
});

test("summarizes Amadeus offers without dumping raw JSON", () => {
  const plan = planTravelSearch(["rom-par"], { date: "2026-06-15", max: 1 }, new Date("2026-05-09T10:00:00Z"));
  const summary = summarizeAmadeusOffers({
    data: [{
      price: { grandTotal: "123.45", currency: "EUR" },
      itineraries: [{
        duration: "PT2H10M",
        segments: [
          { carrierCode: "AZ", departure: { at: "2026-06-15T08:00:00" }, arrival: { at: "2026-06-15T10:10:00" } }
        ]
      }]
    }]
  }, plan);

  assert.equal(summary.count, 1);
  assert.equal(summary.offers[0].total, "123.45");
  assert.equal(summary.offers[0].carrier, "AZ");
});

test("ships curated recipes for sports, travel, and generic creation", () => {
  const ids = RECIPES.map((recipe) => recipe.id);
  assert.ok(ids.includes("sports-nba-espn"));
  assert.ok(ids.includes("travel-amadeus-flight-offers"));
  assert.ok(ids.includes("generic-browser-har"));
});
