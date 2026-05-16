import test from "node:test";
import assert from "node:assert/strict";
import { extractMainContent, htmlToMarkdown } from "../src/core/scrape.js";

test("extractMainContent prefers the longest <main>/<article>", () => {
  const html = `<html><body><nav>menu</nav><main><h1>Title</h1><p>Body content here that is fairly long to win over short articles. Body content here that is fairly long to win over short articles.</p></main><footer>bye</footer></body></html>`;
  const inner = extractMainContent(html);
  assert.match(inner, /<h1>Title<\/h1>/);
  assert.doesNotMatch(inner, /menu|footer/);
});

test("htmlToMarkdown handles headings, lists, links, code", () => {
  const html = `<h1>Title</h1><p>Hello <a href="/x">there</a></p><ul><li>One</li><li>Two</li></ul><pre><code>x = 1</code></pre>`;
  const md = htmlToMarkdown(html, { baseUrl: "https://example.com" });
  assert.match(md, /^# Title/m);
  assert.match(md, /\[there\]\(https:\/\/example\.com\/x\)/);
  assert.match(md, /- One/);
  assert.match(md, /- Two/);
  assert.match(md, /```\nx = 1\n```/);
});
