import type { Fetch, Random } from "./kakao.js";

export const NAVER_LOCAL_URL = "https://openapi.naver.com/v1/search/local.json";

const QUERY_SUFFIXES = ["맛집", "한식", "일식", "중식", "분식", "아시아음식", "도시락", "육류", "치킨", "패스트푸드", "술집"];
const SORTS = ["random", "comment"];

export interface NaverItem {
  title: string;
  link: string;
  category: string;
  description: string;
  address: string;
  roadAddress: string;
  mapx: string;
  mapy: string;
}

interface LocalResponse {
  items?: NaverItem[];
  errorCode?: string;
  errorMessage?: string;
}

/**
 * Naver local search for "<query> <random food word>", fetched with both sort orders; returns one random item
 * from one of the successful result sets, or null when nothing came back.
 */
export async function pickNaverPlace(
  query: string,
  creds: { id: string; secret: string },
  fetchFn: Fetch = fetch,
  random: Random = Math.random,
): Promise<NaverItem | null> {
  const fullQuery = `${query} ${QUERY_SUFFIXES[Math.floor(random() * QUERY_SUFFIXES.length)]}`;
  const sets: NaverItem[][] = [];
  for (const sort of SORTS) {
    const url = new URL(NAVER_LOCAL_URL);
    url.searchParams.set("query", fullQuery);
    url.searchParams.set("display", "5");
    url.searchParams.set("start", "1");
    url.searchParams.set("sort", sort);
    const res = await fetchFn(url, {
      headers: { "X-Naver-Client-Id": creds.id, "X-Naver-Client-Secret": creds.secret },
    });
    const body = (await res.json().catch(() => ({}))) as LocalResponse;
    if (!res.ok || body.errorCode) {
      console.error(`[naver] local search failed: HTTP ${res.status} ${body.errorCode ?? ""}`);
      continue;
    }
    if (body.items?.length) sets.push(body.items);
  }
  if (sets.length === 0) return null;
  const items = sets[Math.floor(random() * sets.length)];
  return items[Math.floor(random() * items.length)];
}
