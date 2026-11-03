import { existsSync, readFileSync } from "node:fs";

// Same shape as the Go service's config/keys.json, so the existing file can be mounted as-is.
export interface Keys {
  kakaoRest: string;
  kakaoJs: string;
  naverId: string;
  naverSecret: string;
  newyoApiKey: string;
}

export interface Config {
  keys: Keys;
  httpPort: number;
  httpsPort: number;
  tlsCertFile: string;
  tlsKeyFile: string;
  redirectAllowedHosts: string[];
}

// Hosts the old scraper sent through `/?redirect=` (news links) plus the place pages ojeommu links to.
export const DEFAULT_REDIRECT_HOSTS = [
  "mobile.hankyung.com",
  "www.hankyung.com",
  "www.mk.co.kr",
  "quicknews.co.kr",
  "place.map.kakao.com",
  "map.kakao.com",
  "map.naver.com",
  "m.place.naver.com",
];

function readKeysFile(path: string): Record<string, Record<string, string>> {
  if (!existsSync(path)) return {};
  return JSON.parse(readFileSync(path, "utf8"));
}

export function loadConfig(env: NodeJS.ProcessEnv = process.env): Config {
  const file = readKeysFile(env.KEYS_FILE ?? "./config/keys.json");
  const keys: Keys = {
    kakaoRest: env.KAKAO_REST_API_KEY ?? file.kakao?.rest_api ?? "",
    kakaoJs: env.KAKAO_JAVASCRIPT_KEY ?? file.kakao?.javascript ?? "",
    naverId: env.NAVER_CLIENT_ID ?? file.naver?.id ?? "",
    naverSecret: env.NAVER_CLIENT_SECRET ?? file.naver?.secret ?? "",
    newyoApiKey: env.NEWYO_API_KEY ?? file.newyo?.apikey ?? "",
  };
  const missing = Object.entries(keys).filter(([, v]) => !v).map(([k]) => k);
  if (missing.length) console.warn(`[config] missing keys: ${missing.join(", ")}`);

  return {
    keys,
    httpPort: Number(env.HTTP_PORT ?? 8080),
    httpsPort: Number(env.HTTPS_PORT ?? 8443),
    tlsCertFile: env.TLS_CERT_FILE ?? "",
    tlsKeyFile: env.TLS_KEY_FILE ?? "",
    redirectAllowedHosts: env.REDIRECT_ALLOWED_HOSTS
      ? env.REDIRECT_ALLOWED_HOSTS.split(",").map((h) => h.trim().toLowerCase()).filter(Boolean)
      : DEFAULT_REDIRECT_HOSTS,
  };
}
