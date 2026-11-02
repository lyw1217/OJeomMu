# 오늘 무먹? (OJeomMu) — TypeScript + Hono

`go/`의 Gin 서비스를 Node 22 + Hono로 옮긴 버전입니다. 템플릿(`templates/`)과 정적 파일(`public/`)은 Go 버전 것을 그대로 복사했습니다.

```bash
npm ci
npm test          # node:test (tsx), 업스트림은 fake fetch
npm run typecheck
npm run dev       # http://localhost:8080
docker buildx build --platform linux/amd64 -t lyw1217/ojeommu:ts .
```

## 설정

| 이름 | 기본값 | 설명 |
|---|---|---|
| `HTTP_PORT` | `8080` | HTTP 포트 (컨테이너는 non-root라 80 대신 8080) |
| `HTTPS_PORT` | `8443` | `TLS_CERT_FILE`·`TLS_KEY_FILE`이 모두 있을 때만 엽니다 |
| `TLS_CERT_FILE`, `TLS_KEY_FILE` | (없음) | PEM 경로. 비우면 HTTP만 → TLS는 DSM 리버스 프록시에서 처리 권장 |
| `KEYS_FILE` | `./config/keys.json` (이미지: `/app/config/keys.json`) | Go 버전과 같은 형식의 keys.json. 기존 파일을 읽기 전용으로 마운트하면 됩니다 |
| `KAKAO_REST_API_KEY`, `KAKAO_JAVASCRIPT_KEY`, `NAVER_CLIENT_ID`, `NAVER_CLIENT_SECRET`, `NEWYO_API_KEY` | keys.json 값 | 있으면 keys.json보다 우선 |
| `REDIRECT_ALLOWED_HOSTS` | 뉴스·지도 호스트 목록 (`src/config.ts`) | `/?redirect=`가 허용하는 호스트 (쉼표 구분) |
| `TZ` | `Asia/Seoul` | |

로그는 stdout으로만 냅니다(`docker logs`). Go 버전의 `log/*.log` 파일 로테이션은 없앴습니다.

기능 대조표는 [PARITY.md](PARITY.md)에 있습니다.
