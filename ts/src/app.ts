import { createHash, timingSafeEqual } from "node:crypto";
import { existsSync, readFileSync } from "node:fs";
import { join } from "node:path";
import { Hono, type Context } from "hono";
import { logger } from "hono/logger";
import { secureHeaders } from "hono/secure-headers";
import { serveStatic } from "@hono/node-server/serve-static";
import type { Config } from "./config.js";
import { distanceMeters } from "./geo.js";
import { pickRestaurant, searchKeyword, type Fetch, type Random, type SearchCond } from "./kakao.js";
import { pickNaverPlace } from "./naver.js";
import { loadTemplates, render } from "./templates.js";

export const TITLE = "오늘 무먹?";

export interface AppDeps {
  config: Config;
  rootDir: string; // folder holding public/ and templates/
  fetch?: Fetch;
  random?: Random;
  log?: boolean;
}

const jsonError = (c: Context, status: 400 | 401 | 404 | 500, reason: string) =>
  c.json({ status, reason }, status);

export function checkAuth(auth: string | undefined, apiKey: string): boolean {
  if (!auth || !apiKey) return false;
  // Clients send base64(apikey); the Go service trimmed surrounding spaces after decoding.
  const decoded = Buffer.from(auth, "base64").toString("utf8").replace(/^ +| +$/g, "");
  const a = createHash("sha256").update(decoded).digest();
  const b = createHash("sha256").update(apiKey).digest();
  return timingSafeEqual(a, b);
}

export function isAllowedRedirect(target: string, allowedHosts: string[]): boolean {
  try {
    const u = new URL(target);
    return (u.protocol === "http:" || u.protocol === "https:") && allowedHosts.includes(u.hostname.toLowerCase());
  } catch {
    return false;
  }
}

export function createApp(deps: AppDeps): Hono {
  const { config } = deps;
  const fetchFn = deps.fetch ?? fetch;
  const random = deps.random ?? Math.random;
  const publicDir = join(deps.rootDir, "public");
  const templates = loadTemplates(join(deps.rootDir, "templates"));
  const page = (name: string, key?: string) => render(templates, name, { title: TITLE, key });

  const app = new Hono();
  if (deps.log !== false) app.use(logger());
  // The pages load Bootstrap/jQuery/Kakao SDK from CDNs and use inline scripts, so no CSP here.
  app.use(secureHeaders({ contentSecurityPolicy: undefined, crossOriginEmbedderPolicy: false }));

  // Static files (same URL layout as the Gin service)
  app.use("/css/*", serveStatic({ root: publicDir }));
  app.use("/js/*", serveStatic({ root: publicDir }));
  app.use("/assets/*", serveStatic({ root: publicDir }));
  app.get("/favicon.ico", serveStatic({ path: join(publicDir, "assets/favicon.ico") }));
  app.get("/robots.txt", serveStatic({ path: join(publicDir, "robots.txt") }));

  // Pages
  const home = (c: Context) => {
    const target = c.req.query("redirect");
    if (target && isAllowedRedirect(target, config.redirectAllowedHosts)) return c.redirect(target, 301);
    return c.html(page("views/index.html", config.keys.kakaoJs));
  };
  app.get("/", home);
  app.get("/index.html", home);
  app.get("/info.html", (c) => c.html(page("views/info.html")));
  app.get("/test.html", (c) => c.html(page("views/test.html", config.keys.kakaoJs)));

  // Web form: random restaurant around the map marker (or around `?query=` place)
  app.post("/sendToGo", async (c) => {
    const body = await c.req.json().catch(() => null);
    if (!body || typeof body !== "object") return jsonError(c, 400, "Bad Request");
    const str = (v: unknown) => (v === undefined || v === null ? "" : String(v));
    const cond: SearchCond = { category: str(body.category), x: str(body.x), y: str(body.y), radius: str(body.radius) };
    let qryPlaceName = "";

    const qry = c.req.query("query");
    try {
      if (qry) {
        const found = await searchKeyword({ query: qry }, config.keys.kakaoRest, fetchFn);
        if (found.length === 0) return jsonError(c, 404, "Not Found");
        Object.assign(cond, { category: "FD6", radius: "0.5", x: found[0].x, y: found[0].y });
        qryPlaceName = found[0].place_name;
      }
      const picked = await pickRestaurant(cond, config.keys.kakaoRest, fetchFn, random);
      if (!picked) return jsonError(c, 404, "Not Found");
      const p = picked.place;
      return c.json({
        ID: p.id,
        NAME: p.place_name,
        CAT_NAME: p.category_name,
        PHONE: p.phone,
        ADDRESS: p.address_name,
        ROAD_ADDRESS: p.road_address_name,
        X: p.x,
        Y: p.y,
        URL: p.place_url,
        DISTANCE: distanceMeters(cond.x, cond.y, p.x, p.y),
        TOTAL_NUMS: picked.total,
        QRY_PLACE_NAME: qryPlaceName,
      });
    } catch (err) {
      console.error("[sendToGo]", err);
      return jsonError(c, 500, "Internal Server Error");
    }
  });

  // Chat-bot API
  app.get("/ojeommu", async (c) => {
    if (!checkAuth(c.req.query("auth"), config.keys.newyoApiKey)) return jsonError(c, 401, "Unauthorized API Key");
    const qry = c.req.query("query");
    try {
      if (c.req.query("target") === "naver") {
        if (!qry) return jsonError(c, 404, "Not Found");
        const item = await pickNaverPlace(qry, { id: config.keys.naverId, secret: config.keys.naverSecret }, fetchFn, random);
        if (!item) return jsonError(c, 404, "Not Found");
        return c.json({ hdr: item.title, place: qry, lnk: item.link, cat: item.category });
      }

      if (!qry) return jsonError(c, 400, "Bad Request");
      const found = await searchKeyword({ query: qry }, config.keys.kakaoRest, fetchFn);
      if (found.length === 0) return jsonError(c, 404, `Query(${qry}) Place Not Found`);
      const origin = found[0];
      const cond: SearchCond = { category: c.req.query("cat") || "anything", radius: "0.5", x: origin.x, y: origin.y };
      const picked = await pickRestaurant(cond, config.keys.kakaoRest, fetchFn, random);
      if (!picked) {
        return jsonError(c, 404, `Not Found Restaurant 500m around (query(${qry}), searchd_place(${origin.place_name}))`);
      }
      return c.json({
        hdr: picked.place.place_name,
        place: origin.place_name,
        d: distanceMeters(cond.x, cond.y, picked.place.x, picked.place.y),
        lnk: picked.place.place_url,
        cat: picked.place.category_name,
      });
    } catch (err) {
      console.error("[ojeommu]", err);
      return jsonError(c, 500, "Internal Server Error");
    }
  });

  // Chat-bot weather icon: /weather?auth=..&query=<n> -> assets/img/wt/wt<n>.png
  app.get("/weather", (c) => {
    if (!checkAuth(c.req.query("auth"), config.keys.newyoApiKey)) return jsonError(c, 401, "Unauthorized API Key");
    const qry = c.req.query("query") ?? "";
    if (!qry) return jsonError(c, 400, "Bad Request");
    const file = join(publicDir, "assets/img/wt", `wt${qry}.png`);
    if (!/^\d{1,3}$/.test(qry) || !existsSync(file)) return jsonError(c, 404, "Not Found");
    return c.body(readFileSync(file), 200, { "Content-Type": "image/png" });
  });

  app.notFound((c) => c.html(page("views/404.html"), 404));
  app.onError((err, c) => {
    console.error("[unhandled]", err);
    return c.html(page("views/500.html"), 500);
  });
  return app;
}
