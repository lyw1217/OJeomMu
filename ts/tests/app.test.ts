import { test } from "node:test";
import assert from "node:assert/strict";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { createApp, checkAuth, isAllowedRedirect } from "../src/app.js";
import { DEFAULT_REDIRECT_HOSTS, type Config } from "../src/config.js";
import type { Place } from "../src/kakao.js";

const rootDir = join(dirname(fileURLToPath(import.meta.url)), "..");
const API_KEY = "test-api-key";
const AUTH = Buffer.from(API_KEY).toString("base64");

const config: Config = {
  keys: { kakaoRest: "rest", kakaoJs: "js-key", naverId: "nid", naverSecret: "nsecret", newyoApiKey: API_KEY },
  httpPort: 0,
  httpsPort: 0,
  tlsCertFile: "",
  tlsKeyFile: "",
  redirectAllowedHosts: DEFAULT_REDIRECT_HOSTS,
};

const place = (id: string, category: string, x = "127.0", y = "37.5"): Place => ({
  id,
  place_name: `place-${id}`,
  category_name: `음식점 > ${category}`,
  category_group_code: "FD6",
  phone: "02-000-0000",
  address_name: "addr",
  road_address_name: "road",
  x,
  y,
  place_url: `http://place.map.kakao.com/${id}`,
  distance: "",
});

interface Call {
  url: URL;
  headers: Record<string, string>;
}

/** Fake upstreams: Kakao keyword search with/without `rect`, Naver local search. */
function fakeFetch(opts: { origin?: Place[]; around?: Place[]; naver?: unknown; kakaoStatus?: number } = {}) {
  const calls: Call[] = [];
  const fn = (async (input: URL | string, init?: RequestInit) => {
    const url = new URL(String(input));
    calls.push({ url, headers: (init?.headers ?? {}) as Record<string, string> });
    if (url.hostname === "dapi.kakao.com") {
      if (opts.kakaoStatus) return new Response("{}", { status: opts.kakaoStatus });
      const docs = url.searchParams.has("rect") ? (opts.around ?? []) : (opts.origin ?? []);
      return Response.json({ documents: docs, meta: { is_end: true, pageable_count: docs.length } });
    }
    if (url.hostname === "openapi.naver.com") return Response.json(opts.naver ?? { items: [] });
    throw new Error(`unexpected fetch ${url}`);
  }) as typeof fetch;
  return { fn, calls };
}

const make = (f = fakeFetch()) => createApp({ config, rootDir, fetch: f.fn, random: () => 0, log: false });

test("home and index render the Gin templates with the Kakao JS key", async () => {
  const app = make();
  for (const path of ["/", "/index.html"]) {
    const res = await app.request(path);
    assert.equal(res.status, 200);
    const html = await res.text();
    assert.match(html, /<title>오늘 무먹\?<\/title>/);
    assert.match(html, /sdk\.js\?appkey=js-key&libraries=services/);
    assert.doesNotMatch(html, /\{\{/);
  }
});

test("info and test pages render", async () => {
  const app = make();
  for (const path of ["/info.html", "/test.html"]) {
    const res = await app.request(path);
    assert.equal(res.status, 200);
    assert.doesNotMatch(await res.text(), /\{\{/);
  }
});

test("unknown path renders the 404 page", async () => {
  const res = await make().request("/nope");
  assert.equal(res.status, 404);
  assert.match(await res.text(), /404/);
});

test("static files are served at the Gin paths", async () => {
  const app = make();
  for (const [path, type] of [
    ["/css/styles.css", /css/],
    ["/js/scripts.js", /javascript/],
    ["/js/map.js", /javascript/],
    ["/assets/img/gps.png", /png/],
    ["/favicon.ico", /icon/],
    ["/robots.txt", /text/],
  ] as const) {
    const res = await app.request(path);
    assert.equal(res.status, 200, path);
    assert.match(res.headers.get("content-type") ?? "", type, path);
  }
  assert.equal((await app.request("/assets/../src/app.ts")).status, 404);
});

test("?redirect= only follows allow-listed hosts", async () => {
  const app = make();
  const ok = await app.request("/?redirect=" + encodeURIComponent("https://www.mk.co.kr/premium/series/20007/"));
  assert.equal(ok.status, 301);
  assert.equal(ok.headers.get("location"), "https://www.mk.co.kr/premium/series/20007/");
  const evil = await app.request("/?redirect=" + encodeURIComponent("https://evil.example/"));
  assert.equal(evil.status, 200);
  assert.equal(isAllowedRedirect("javascript:alert(1)", DEFAULT_REDIRECT_HOSTS), false);
  assert.equal(isAllowedRedirect("//www.mk.co.kr/", DEFAULT_REDIRECT_HOSTS), false);
});

test("POST /sendToGo picks a restaurant matching the category", async () => {
  const f = fakeFetch({ around: [place("1", "중식 > 짜장면"), place("2", "한식 > 국밥", "127.001", "37.5")] });
  const res = await make(f).request("/sendToGo", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ category: "korea", radius: "0.5", x: "127.0", y: "37.5" }),
  });
  assert.equal(res.status, 200);
  const body = await res.json();
  assert.equal(body.ID, "2");
  assert.equal(body.NAME, "place-2");
  assert.equal(body.URL, "http://place.map.kakao.com/2");
  assert.equal(body.QRY_PLACE_NAME, "");
  // 4 quadrants x identical fake results -> each place appears 4 times; only the 한식 one matches.
  assert.equal(body.TOTAL_NUMS, 4);
  assert.ok(body.DISTANCE > 80 && body.DISTANCE < 90, `distance ${body.DISTANCE}`);
  const rects = f.calls.filter((c) => c.url.searchParams.has("rect"));
  assert.equal(rects.length, 4);
  assert.equal(rects[0].headers.Authorization, "KakaoAK rest");
  assert.equal(rects[0].url.searchParams.get("query"), "맛집");
  assert.equal(rects[0].url.searchParams.get("category_group_code"), "FD6");
});

test("POST /sendToGo?query= searches around the named place", async () => {
  const f = fakeFetch({ origin: [place("o", "역", "126.9", "37.4")], around: [place("1", "일식")] });
  const res = await make(f).request("/sendToGo?query=" + encodeURIComponent("강남역"), {
    method: "POST",
    body: JSON.stringify({ category: "none", radius: "1", x: "0", y: "0" }),
  });
  assert.equal(res.status, 200);
  assert.equal((await res.json()).QRY_PLACE_NAME, "place-o");
  const rect = f.calls.find((c) => c.url.searchParams.has("rect"))!.url.searchParams.get("rect")!;
  const [x1] = rect.split(",").map(Number);
  assert.ok(Math.abs(x1 - 126.9) < 0.01, rect); // around the origin, radius forced to 0.5 km
});

test("POST /sendToGo error cases", async () => {
  const post = (app: ReturnType<typeof make>, body: string, q = "") =>
    app.request("/sendToGo" + q, { method: "POST", body });
  assert.equal((await post(make(), "not json")).status, 400);
  assert.equal((await post(make(), JSON.stringify({ category: "korea", x: "a", y: "b" }))).status, 500);
  assert.equal((await post(make(), JSON.stringify({ category: "korea", x: "127", y: "37" }))).status, 404);
  assert.equal((await post(make(), JSON.stringify({}), "?query=x")).status, 404);
  assert.equal((await post(make(fakeFetch({ kakaoStatus: 401 })), JSON.stringify({}), "?query=x")).status, 500);
});

test("GET /ojeommu requires base64(apikey)", async () => {
  const app = make();
  const res = await app.request("/ojeommu?query=x");
  assert.equal(res.status, 401);
  assert.deepEqual(await res.json(), { status: 401, reason: "Unauthorized API Key" });
  assert.equal((await app.request("/ojeommu?auth=" + Buffer.from("wrong").toString("base64"))).status, 401);
  assert.equal(checkAuth(Buffer.from(` ${API_KEY} `).toString("base64"), API_KEY), true);
  assert.equal(checkAuth("", ""), false);
});

test("GET /ojeommu (kakao) returns the bot payload", async () => {
  const f = fakeFetch({ origin: [place("o", "역")], around: [place("1", "치킨")] });
  const res = await make(f).request(`/ojeommu?auth=${AUTH}&query=x&cat=chicken`);
  assert.equal(res.status, 200);
  assert.deepEqual(await res.json(), {
    hdr: "place-1",
    place: "place-o",
    d: 0,
    lnk: "http://place.map.kakao.com/1",
    cat: "음식점 > 치킨",
  });
  assert.equal((await make().request(`/ojeommu?auth=${AUTH}`)).status, 400);
  const notFound = await make().request(`/ojeommu?auth=${AUTH}&query=zz`);
  assert.equal(notFound.status, 404);
  assert.equal((await notFound.json()).reason, "Query(zz) Place Not Found");
  const noFood = await make(fakeFetch({ origin: [place("o", "역")] })).request(`/ojeommu?auth=${AUTH}&query=zz`);
  assert.match((await noFood.json()).reason, /Not Found Restaurant 500m around/);
});

test("GET /ojeommu?target=naver picks a Naver local result", async () => {
  const item = { title: "<b>국밥</b>집", link: "https://n.example", category: "한식", description: "", address: "", roadAddress: "", mapx: "", mapy: "" };
  const f = fakeFetch({ naver: { items: [item] } });
  const res = await make(f).request(`/ojeommu?auth=${AUTH}&target=naver&query=${encodeURIComponent("강남")}`);
  assert.equal(res.status, 200);
  assert.deepEqual(await res.json(), { hdr: "<b>국밥</b>집", place: "강남", lnk: "https://n.example", cat: "한식" });
  assert.equal(f.calls.length, 2);
  assert.equal(f.calls[0].headers["X-Naver-Client-Id"], "nid");
  assert.equal(f.calls[0].url.searchParams.get("query"), "강남 맛집");
  assert.equal(f.calls[1].url.searchParams.get("query"), "강남 맛집");
  assert.deepEqual(f.calls.map((c) => c.url.searchParams.get("sort")), ["random", "comment"]);
  assert.equal((await make().request(`/ojeommu?auth=${AUTH}&target=naver&query=x`)).status, 404);
  assert.equal((await make().request(`/ojeommu?auth=${AUTH}&target=naver`)).status, 404);
});

test("GET /weather serves wt<n>.png behind auth", async () => {
  const app = make();
  assert.equal((await app.request("/weather?query=1")).status, 401);
  const res = await app.request(`/weather?auth=${AUTH}&query=1`);
  assert.equal(res.status, 200);
  assert.equal(res.headers.get("content-type"), "image/png");
  assert.equal((await app.request(`/weather?auth=${AUTH}&query=999`)).status, 404);
  assert.equal((await app.request(`/weather?auth=${AUTH}&query=${encodeURIComponent("../../x")}`)).status, 404);
  assert.equal((await app.request(`/weather?auth=${AUTH}`)).status, 400);
});
