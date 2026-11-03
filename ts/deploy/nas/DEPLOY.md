# NAS 배포 절차 (TypeScript ojeommu)

Go 버전(`lyw1217-ojeommu-1`, 2026-11-02 22:44부터 정지)을 TS 버전으로 바꾸는 절차입니다.
다른 NAS 서비스와 같은 방식으로 DSM 리버스 프록시 뒤에 둡니다.

```
https://ojeommu.lyw1217.synology.me ──(DSM 리버스 프록시, TLS)──▶ http://localhost:3023 ──▶ 컨테이너 :8080
```

- 컨테이너는 `127.0.0.1:3023`에만 열립니다. 외부에서는 리버스 프록시로만 접근합니다.
- HTTPS는 DSM이 처리하므로 컨테이너에 인증서를 넣지 않습니다. 기존 Go 인증서는 2024-08-20에 만료됐습니다.
- 키는 기존 `/volume1/docker/ojeommu/config/keys.json`을 읽기 전용으로 마운트합니다. 파일은 옮기지도 바꾸지도 않습니다.
- 이 레포에는 NAS webhook이 없습니다. main에 머지해도 자동 배포되지 않으며, 아래 절차를 직접 실행해야 합니다.

## 0. 사전 준비

1. PR을 **squash merge**합니다. 첫 커밋의 작성자 이메일이 기록에 남지 않게 하려는 것입니다.
2. [카카오 개발자 콘솔](https://developers.kakao.com) → 내 애플리케이션 → 플랫폼 → Web 사이트 도메인에 `https://ojeommu.lyw1217.synology.me`를 추가합니다.
   - 이 도메인을 추가하지 않으면 페이지는 열려도 지도가 뜨지 않습니다.
   - REST API(검색)는 도메인 제한과 무관합니다.

## 1. 레포 받기 (NAS SSH, 처음 한 번)

NAS 셸에는 git이 없어서 git 컨테이너로 받습니다. 컨테이너는 Synology ACL 때문에 `/volume1/docker`에 폴더를 만들 수 없으므로, 폴더는 셸에서 먼저 만듭니다.

```bash
mkdir /volume1/docker/ojeommu-ts && cd /volume1/docker/ojeommu-ts
docker run --rm -u "$(id -u):$(id -g)" -e HOME=/tmp -v "$PWD":/w -w /w alpine/git:v2.52.0 \
  clone https://github.com/lyw1217/OJeomMu.git .
test -f /volume1/docker/ojeommu/config/keys.json && echo keys OK
```

`keys OK`가 나오지 않으면 진행하지 마세요. 파일이 없으면 Docker가 그 경로에 빈 디렉터리를 만들어 마운트하고, 앱은 시작하자마자 종료됩니다.

## 2. 빌드·실행

```bash
cd /volume1/docker/ojeommu-ts
docker compose -f ts/deploy/nas/docker-compose.yml up -d --build
docker ps --filter name=ojeommu --format '{{.Names}} {{.Status}}'   # 1분 뒤 (healthy)
curl -s -o /dev/null -w '%{http_code}\n' http://127.0.0.1:3023/     # 200
```

- 포트를 바꾸려면 `OJEOMMU_PORT`를, 키 파일 경로를 바꾸려면 `OJEOMMU_KEYS_FILE`을 환경변수로 줍니다.
- 기존 Go 컨테이너는 지우지 말고 정지 상태로 둡니다. 되돌릴 때 사용합니다.

## 3. DSM 리버스 프록시 규칙 추가 (DSM 웹 화면)

제어판 → 로그인 포털 → 고급 → 역방향 프록시 → 생성:

| 항목 | 값 |
|---|---|
| 설명 | `ojeommu` |
| 소스 | HTTPS · `ojeommu.lyw1217.synology.me` · 443 |
| 대상 | HTTP · `localhost` · 3023 |

인증서는 다른 `*.lyw1217.synology.me` 서비스에 쓰는 것과 같은 인증서가 적용되는지 확인하세요.

## 4. 확인

- 브라우저에서 `https://ojeommu.lyw1217.synology.me`를 열고, 지도가 뜨는지 봅니다. 이어서 거리와 카테고리를 고르고 "오늘 뭐 먹을까?"를 눌러 결과 창이 나오는지 확인합니다.
- 봇 API:
  ```bash
  curl "https://ojeommu.lyw1217.synology.me/ojeommu?auth=<base64(apikey)>&query=강남역"
  ```
  실제 키 값은 셸 기록에 남지 않게 주의하세요.
- 로그: `docker logs -f ojeommu`. 인증값은 로그에 `auth=***`로 가려져 나옵니다.

## 갱신

main에 새 커밋이 생기면 다음을 실행합니다.

```bash
docker run --rm -u "$(id -u):$(id -g)" -e HOME=/tmp -v /volume1/docker/ojeommu-ts:/w -w /w alpine/git:v2.52.0 pull --ff-only origin main
cd /volume1/docker/ojeommu-ts && docker compose -f ts/deploy/nas/docker-compose.yml up -d --build
```

나중에 webhook으로 자동화하려면 `deploy/nas/update.sh`를 쓰면 됩니다. 이 스크립트는 git이 있는 webhook 컨테이너 안에서 실행하는 다른 서비스의 스크립트와 같은 구조입니다. 다만 자동화하려면 세 곳을 바꿔야 합니다.

- `hooks.json`에 이 레포용 hook 추가
- webhook 컨테이너에 `/volume1/docker/ojeommu-ts` 마운트 추가
- GitHub 레포 설정에 webhook 등록

## 되돌리기

```bash
cd /volume1/docker/ojeommu-ts && docker compose -f ts/deploy/nas/docker-compose.yml down
docker start lyw1217-ojeommu-1   # Go 버전 (호스트 8080/8443/9090)
```

DSM 리버스 프록시 규칙은 지워도 되고 그대로 둬도 됩니다. Go 버전은 포트 3023을 쓰지 않으므로, 규칙을 남겨 두면 그 주소는 502를 반환합니다.

## 카카오톡 봇 (messengerbot)

휴대폰의 봇 스크립트는 이미 사라진 `lunchtoday.site`를 호출하고 있습니다. 봇을 다시 쓰려면 스크립트의 주소를 `https://ojeommu.lyw1217.synology.me`로 바꿔야 합니다.
