import { existsSync, readFileSync } from "node:fs";
import { createServer as createHttpsServer } from "node:https";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { serve } from "@hono/node-server";
import { createApp } from "./app.js";
import { loadConfig } from "./config.js";

const config = loadConfig();
const rootDir = join(dirname(fileURLToPath(import.meta.url)), "..");
const app = createApp({ config, rootDir });

const servers = [serve({ fetch: app.fetch, port: config.httpPort }, (i) => console.log(`[ojeommu] http on :${i.port}`))];

// HTTPS only when a cert pair is mounted; otherwise terminate TLS at the DSM reverse proxy.
if (config.tlsCertFile && config.tlsKeyFile) {
  if (!existsSync(config.tlsCertFile) || !existsSync(config.tlsKeyFile)) {
    console.error("[ojeommu] TLS_CERT_FILE / TLS_KEY_FILE not found");
    process.exit(1);
  }
  servers.push(
    serve(
      {
        fetch: app.fetch,
        port: config.httpsPort,
        createServer: createHttpsServer,
        serverOptions: { cert: readFileSync(config.tlsCertFile), key: readFileSync(config.tlsKeyFile) },
      },
      (i) => console.log(`[ojeommu] https on :${i.port}`),
    ),
  );
}

// Stop accepting connections, let in-flight requests finish, then exit (forced after SHUTDOWN_TIMEOUT_MS).
const SHUTDOWN_TIMEOUT_MS = 10_000;
for (const sig of ["SIGINT", "SIGTERM"] as const) {
  process.once(sig, () => {
    console.log(`[ojeommu] ${sig}, shutting down`);
    setTimeout(() => {
      console.error("[ojeommu] in-flight requests did not finish, forcing exit");
      process.exit(1);
    }, SHUTDOWN_TIMEOUT_MS).unref();
    let open = servers.length;
    for (const s of servers) s.close(() => --open === 0 && process.exit(0));
  });
}
