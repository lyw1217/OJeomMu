import { quadrantRects } from "./geo.js";

export const KAKAO_KEYWORD_URL = "https://dapi.kakao.com/v2/local/search/keyword.json";
const MAX_SEARCH_PAGE = 10; // pages 1..9, same bound as the Go loop
export const UPSTREAM_TIMEOUT_MS = 5_000;

export interface Place {
  id: string;
  place_name: string;
  category_name: string;
  category_group_code: string;
  phone: string;
  address_name: string;
  road_address_name: string;
  x: string;
  y: string;
  place_url: string;
  distance: string;
}

interface KeywordResponse {
  documents?: Place[];
  meta?: { is_end?: boolean; pageable_count?: number };
}

export interface KeywordParams {
  query: string;
  category_group_code?: string;
  rect?: string;
}

export type Fetch = typeof fetch;

/** Kakao keyword search, following pages until `is_end` (max 9 pages, 15 per page). */
export async function searchKeyword(params: KeywordParams, restKey: string, fetchFn: Fetch = fetch): Promise<Place[]> {
  const places: Place[] = [];
  for (let page = 1; page < MAX_SEARCH_PAGE; page++) {
    const url = new URL(KAKAO_KEYWORD_URL);
    url.searchParams.set("query", params.query);
    if (params.category_group_code) url.searchParams.set("category_group_code", params.category_group_code);
    if (params.rect) url.searchParams.set("rect", params.rect);
    url.searchParams.set("page", String(page));
    url.searchParams.set("size", "15");
    url.searchParams.set("sort", "accuracy");

    const res = await fetchFn(url, {
      headers: { Authorization: `KakaoAK ${restKey}` },
      signal: AbortSignal.timeout(UPSTREAM_TIMEOUT_MS), // also bounds reading the body
    });
    if (!res.ok) throw new Error(`kakao keyword search: HTTP ${res.status}`);
    const body = (await res.json()) as KeywordResponse;
    places.push(...(body.documents ?? []));
    if (body.meta?.is_end !== false) break;
  }
  return places;
}

// Category values sent by the web form / bot, mapped to the text Kakao puts in `category_name`.
const CATEGORY_LABELS: Record<string, string> = {
  korea: "한식",
  china: "중식",
  japan: "일식",
  western: "양식",
  flour: "분식",
  asia: "아시아음식",
  lunchbox: "도시락",
  meat: "육류,고기",
  chicken: "치킨",
  fastfood: "패스트푸드",
  bar: "술집",
};

export function categoryLabel(category: string | undefined): string | null {
  return (category && Object.hasOwn(CATEGORY_LABELS, category) && CATEGORY_LABELS[category]) || null;
}

/** strconv.ParseFloat-like: no surrounding spaces, no hex/binary; NaN when it would fail in Go. */
export function parseFloatStrict(s: string | undefined): number {
  return s && /^[+-]?(\d+\.?\d*|\.\d+)([eE][+-]?\d+)?$/.test(s) ? Number(s) : NaN;
}

export interface SearchCond {
  category?: string;
  x: string;
  y: string;
  radius?: string;
}

export type Random = () => number;

/**
 * Search "맛집" restaurants (FD6) in the four quadrants around (x, y), filter by category and pick one at random.
 * Returns null when nothing matched. Quadrants that fail are skipped, as in the Go version, but when all four
 * fail it throws, so an upstream outage is not reported as "no restaurants".
 */
export async function pickRestaurant(
  cond: SearchCond,
  restKey: string,
  fetchFn: Fetch = fetch,
  random: Random = Math.random,
): Promise<{ place: Place; total: number } | null> {
  const lng = parseFloatStrict(cond.x);
  const lat = parseFloatStrict(cond.y);
  if (!Number.isFinite(lng) || !Number.isFinite(lat)) {
    throw new Error(`invalid coordinates x=${cond.x} y=${cond.y}`);
  }
  const radius = parseFloatStrict(cond.radius);
  const km = Number.isFinite(radius) ? radius : 0.3;

  const results = await Promise.allSettled(
    quadrantRects(lat, lng, km).map((rect) =>
      searchKeyword({ query: "맛집", category_group_code: "FD6", rect }, restKey, fetchFn),
    ),
  );
  if (results.every((r) => r.status === "rejected")) {
    throw new Error("kakao: all quadrant searches failed", { cause: results.map((r) => (r as PromiseRejectedResult).reason) });
  }
  const all = results.flatMap((r) => {
    if (r.status === "fulfilled") return r.value;
    console.error("[kakao] quadrant search failed:", r.reason);
    return [];
  });

  const label = categoryLabel(cond.category);
  const matched = label ? all.filter((p) => p.category_name.includes(label)) : all;
  if (matched.length === 0) return null;
  return { place: matched[Math.floor(random() * matched.length)], total: matched.length };
}
