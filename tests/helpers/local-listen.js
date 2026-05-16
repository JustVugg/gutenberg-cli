import net from "node:net";

let cached = null;

export async function localListenAvailable() {
  if (cached !== null) return cached;
  cached = await new Promise((resolve) => {
    const server = net.createServer();
    server.once("error", () => resolve(false));
    server.listen(0, "127.0.0.1", () => {
      server.close(() => resolve(true));
    });
  });
  return cached;
}

export async function skipIfNoLocalListen(t) {
  if (await localListenAvailable()) return false;
  t.skip("local TCP listen is unavailable in this sandbox");
  return true;
}
