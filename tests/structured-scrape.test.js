import test from "node:test";
import assert from "node:assert/strict";
import { detectStructuredLists } from "../src/core/scrape.js";

test("detectStructuredLists finds time-title patterns in <ul>", () => {
  const html = `<ul>
    <li>18:00 - TG1</li>
    <li>18:30 - Soliti Ignoti</li>
    <li>20:30 - Affari Tuoi</li>
    <li>21:30 - Sanremo Show</li>
  </ul>`;
  const tables = detectStructuredLists(html);
  assert.ok(tables.length > 0);
  assert.ok(tables[0].rows.length >= 4);
  assert.equal(tables[0].rows[0].time, "18:00");
  assert.equal(tables[0].rows[0].title, "TG1");
});

test("detectStructuredLists ignores short lists", () => {
  const html = `<ul><li>18:00 - One</li></ul>`;
  const tables = detectStructuredLists(html);
  assert.equal(tables.length, 0);
});

test("detectStructuredLists finds bare lines pattern", () => {
  const html = `<div>
    21:30 - Show A
    22:00 - Show B
    23:00 - News
  </div>`;
  const tables = detectStructuredLists(html);
  assert.ok(tables.length > 0);
  assert.ok(tables[0].rows.length >= 3);
});
