import test from "node:test";
import assert from "node:assert/strict";
import http from "node:http";
import { extractFromText } from "../src/core/extract.js";
import { skipIfNoLocalListen } from "./helpers/local-listen.js";

test("ollama provider hits /api/chat and returns parsed JSON", async (t) => {
  if (await skipIfNoLocalListen(t)) return;
  const server = http.createServer((req, res) => {
    let body = "";
    req.on("data", (chunk) => { body += chunk; });
    req.on("end", () => {
      const payload = JSON.parse(body);
      assert.equal(payload.format, "json");
      assert.equal(payload.stream, false);
      assert.equal(typeof payload.messages[0].content, "string");
      res.setHeader("Content-Type", "application/json");
      res.end(JSON.stringify({ message: { content: '{"title":"Demo","items":[1,2]}' } }));
    });
  });
  await new Promise((resolve) => server.listen(0, resolve));
  const port = server.address().port;
  try {
    const result = await extractFromText("Some HTML content here", {
      provider: "ollama",
      host: `http://127.0.0.1:${port}`,
      model: "test-model",
      prompt: "Extract"
    });
    assert.deepEqual(result.parsed, { title: "Demo", items: [1, 2] });
    assert.equal(result.errors.length, 0);
  } finally {
    server.close();
  }
});

test("ollama provider reports clear error when unreachable", async () => {
  await assert.rejects(
    () => extractFromText("text", { provider: "ollama", host: "http://127.0.0.1:1" }),
    /Ollama unreachable/
  );
});
