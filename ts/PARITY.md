# Go → TypeScript 기능 대조표 (ojeommu)

기준: `go/` (HEAD 88c1a50). NAS `/volume1/docker/ojeommu` 배포본과 동일합니다(`main.go`의 미사용 import 한 줄만 다름).
✅ 동일 · 🔧 동작은 같고 버그/보안만 수정 · ➖ 의도적으로 뺌

## 라우트

| Go (Gin) | TS (Hono) | 상태 | 테스트 |
|---|---|---|---|
| `GET /`, `/index.html` → index 템플릿, 카카오 JS 키 주입 | 같음 | ✅ | `home and index render…` |
| `GET /?redirect=<url>` → 아무 URL로 301 | 허용 호스트(http/https)만 301, 나머지는 홈 화면 | 🔧 오픈 리다이렉트 차단 | `?redirect= only follows…` |
| `GET /info.html`, `/test.html` | 같음 (`test.html`에도 JS 키를 주입. Go는 빈 키라 지도가 깨졌음) | ✅/🔧 | `info and test pages render` |
| NoRoute → `404.html` | `notFound` → `404.html` | ✅ | `unknown path…` |
| 처리 안 된 panic → Gin Recovery 500(빈 본문) | `500.html` 템플릿 | 🔧 | — |
| `/css/*`, `/js/*`, `/assets/*`, `/favicon.ico`, `/robots.txt` | 같은 경로, `public/`에서 제공 | ✅ | `static files…` |
| `POST /sendToGo` body `{category,x,y,radius}` (모두 문자열) | 같은 요청·응답 키 (`ID,NAME,CAT_NAME,PHONE,ADDRESS,ROAD_ADDRESS,X,Y,URL,DISTANCE,TOTAL_NUMS,QRY_PLACE_NAME`) | ✅ | `POST /sendToGo picks…` |
| body가 배열이거나 필드가 문자열이 아님 (예: `"x":127`) → Gin `BindJSON` 400 | 400 (`null`/누락은 Go처럼 빈 문자열) | ✅ | `…rejects non-object bodies…` |
| `POST /sendToGo?query=` → 장소 검색 후 그 좌표, 반경 0.5km | 같음 | ✅ | `…?query= searches around…` |
| `/sendToGo` JSON 파싱·바인딩 실패 → 실제로는 400 + `{"status":500}` 본문 | 400 `{"status":400,"reason":"Bad Request"}` | 🔧 | `error cases` |
| `GET /ojeommu?auth=&target=kakao&query=&cat=` → `{hdr,place,d,lnk,cat}` | 같음, 오류 문구도 같음 | ✅ | `GET /ojeommu (kakao)…` |
| `GET /ojeommu?target=naver&query=` → `{hdr,place,lnk,cat}` | 같음 | ✅/🔧 | `…target=naver…` |
| `auth` = base64(apikey), 앞뒤 공백 trim | 같음 + 상수 시간 비교, apikey가 비어 있으면 거부 | 🔧 (Go는 키와 auth가 둘 다 비면 통과) | `requires base64(apikey)` |
| `GET /weather?auth=&query=N` → `assets/img/wt/wtN.png` | 같음, `N`은 숫자만 | 🔧 경로 입력 검증 | `GET /weather serves…` |
| `/weather`에 query 없음 → 200 빈 본문 | 400 | 🔧 | 같음 |

## 검색 로직

| 항목 | 상태 | 비고 |
|---|---|---|
| 카카오 키워드 검색 (`KakaoAK` 헤더, size 15, accuracy, 1~9페이지 — Go 루프와 같은 상한) | 🔧 | `is_end`에서 멈춥니다. Go는 `pageable_count>2`를 페이지 수처럼 써서 마지막 페이지를 중복 조회할 수 있었음. 공식 노출 한도(45건)상 보통 3페이지 이내 |
| 4개 사분면 rect × "맛집" FD6 검색, 실패한 사분면은 건너뜀 | ✅/🔧 | 사분면 4개를 병렬로 호출 (Go는 순차 + 10ms sleep). **네 개가 모두 실패하면 500** (Go는 결과 없음 404) |
| 외부 API 제한 시간 없음 | 🔧 | 카카오·네이버 요청마다 5초 (본문 읽기 포함) |
| 반경 파싱 실패 시 0.3km | ✅ | Go `ParseFloat`처럼 엄격하게 파싱 (`" "`, `0x10` 등은 실패 → 0.3km). 좌표도 같은 규칙 |
| 카테고리 표 (korea→한식 … bar→술집, 그 외 anything) | ✅ | `constructor`·`__proto__` 같은 상속 키도 anything. `category labels match…`, `inherited object keys…` |
| 무작위 선택 | 🔧 | 균등 분포 (Go는 `rand(180) % n`이라 약간 치우침) |
| 거리 계산 (haversine, R=6378.1km, 미터 정수) | 🔧 | Go는 두 위도 모두 `lat2`를 써서 남북 거리가 0이 됐음 |
| 사분면 좌표 계산 | 🔧 | Go 공식의 `sin(lat)·sin(lat)` 오타(`sin(lat2)`여야 함) 수정. 결과 차이는 미터 단위 미만 |
| 네이버 지역 검색 (`"<query> <무작위 음식 단어>"`, random+comment 두 번, display 5) | 🔧 | Go는 두 번째 요청에 단어를 한 번 더 붙였고, 결과 0건이거나 한쪽이 실패하면 panic이 났음. 이제 한쪽이 실패(HTTP·네트워크·타임아웃)해도 다른 쪽 결과를 쓰고, 둘 다 없으면 404 |
| 카카오 업스트림이 HTTP 오류를 주면 | 🔧 | Go는 "결과 없음(404)"으로 취급했음. 이제 장소 검색(`?query=`) 실패 또는 사분면 4개 모두 실패 시 500. 일부 사분면만 실패하면 건너뜀 |

## 실행 환경

| Go | TS | 상태 |
|---|---|---|
| `:80` HTTP + `:443` HTTPS (포트 하드코딩, `HTTP_PORT` 등 환경변수는 무시됨) | `HTTP_PORT`(8080) + 선택적 `HTTPS_PORT`(8443) | 🔧 non-root라 1024 이상 포트. compose에서 `8080:8080`, `8443:8443`으로 매핑 |
| TLS: `./cert/archive/{cert,privkey}.pem` 고정 (CN=lyw1217.synology.me, **2024-08-20 만료**) | `TLS_CERT_FILE`/`TLS_KEY_FILE`을 줄 때만 HTTPS | 🔧 DSM 리버스 프록시 권장 |
| `REDIRECT_PORT=9090` 노출 | 없음 | ➖ Go 코드에 9090 리스너가 없었음(노출만 됨) |
| `config/keys.json` | 같은 파일·형식 + 환경변수 우선 | ✅ |
| `kakao.admin`, `naver.ncp_id/ncp_secret` | 읽지 않음 | ➖ Go에서도 미사용 |
| logrus + lumberjack 파일 로그 (`log/<host>_ojeommu.log`), `auth` 쿼리값이 그대로 남음 | stdout (Hono logger), `auth=***`로 가림 | ➖/🔧 |
| 종료 처리 없음 | SIGTERM 시 새 연결을 받지 않고 진행 중인 요청을 마친 뒤 종료 (최대 10초) | 🔧 |
| `GIN_MODE=debug`, root 실행, Debian 10 buster | `NODE_ENV=production`, `USER node`, `node:22.23.4-bookworm-slim`, npm 제거, HEALTHCHECK | 🔧 |
| 보안 헤더 없음 | `hono/secure-headers` (CSP 제외: CDN과 인라인 스크립트 사용) | 🔧 |

Trivy 0.69.1 (HIGH/CRITICAL): 수정 가능한 취약점 0건. Debian 12.15 기본 패키지에 수정판이 아직 없는 33건이 남아 있습니다.
