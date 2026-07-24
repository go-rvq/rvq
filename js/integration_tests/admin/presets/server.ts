// Boots the Go preset test app (in-memory SQLite, seeded) as a real HTTP server
// for the bun integration tests. bun-only (Bun.spawn); the server binary is built
// once and reused.

// server.ts lives at rvq/js/integration_tests/admin/presets/ — four directories
// down from the repo root (presets → admin → integration_tests → js → rvq).
const REPO_ROOT = new URL("../../../../", import.meta.url).pathname;
const SERVER_PKG = "./js/integration_tests/admin/presets/server";
const SERVER_BIN = REPO_ROOT + ".tmp/presets_testserver";

export interface TestServer {
  url: string;
  stop(): void;
}

// build the server binary once per process, so each startServer() spawns the
// prebuilt binary instead of paying `go run`'s recompile every time.
let builtOnce: Promise<void> | null = null;
function ensureBuilt(): Promise<void> {
  if (!builtOnce) {
    builtOnce = (async () => {
      const p = Bun.spawn(["go", "build", "-o", SERVER_BIN, SERVER_PKG], {
        cwd: REPO_ROOT,
        stdout: "inherit",
        stderr: "inherit",
      });
      if ((await p.exited) !== 0) throw new Error("go build server failed");
    })();
  }
  return builtOnce;
}

// startServer boots the Go test server on an ephemeral port. `app` selects the
// fixture: "listeditor" (default) or "nested" (the four-level nested list editor).
export async function startServer(
  app: "listeditor" | "nested" = "listeditor",
): Promise<TestServer> {
  await ensureBuilt();
  const proc = Bun.spawn([SERVER_BIN], {
    cwd: REPO_ROOT,
    env: { ...process.env, PORT: "0", APP: app },
    stdout: "pipe",
    stderr: "inherit",
  });

  const reader = proc.stdout.getReader();
  const decoder = new TextDecoder();
  let buffer = "";
  const deadline = Date.now() + 60_000;
  while (Date.now() < deadline) {
    const { value, done } = await reader.read();
    if (done) break;
    buffer += decoder.decode(value, { stream: true });
    const m = buffer.match(/LISTENING (\S+)/);
    if (m) {
      reader.releaseLock();
      return { url: m[1], stop: () => proc.kill() };
    }
  }
  proc.kill();
  throw new Error("server did not become ready:\n" + buffer);
}
