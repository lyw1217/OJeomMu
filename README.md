# 오늘 점심 뭐먹지? 오점무? (OJM)

오늘 점심 뭐먹을지 골라볼까

## 1. 개발 환경

- CentOS Linux release 7.9.2009 (Core)
- go version go1.17.3 linux/amd64
  - github.com/gin-gonic/gin v1.7.7
  - github.com/google/go-querystring v1.1.0
    
## 2. Let's Encrypt(Certbot) 설치 및 인증서 발급

### 1. [Certbot 설치](https://certbot.eff.org/instructions)
### 2. SSL 인증서 발급 및 설정
    
        $ sudo certbot certonly --standalone -d {domain}

### 3. 인증서 갱신 테스트

        $ sudo certbot renew --dry-run

## 3. 스크립트

### - `monitor.sh`

서버 기동 확인 및 기동

- 현재 구동중인 애플리케이션 확인
  - 정상 실행 중이라면 스크립트 종료
  - 구동중인 애플리케이션이 없으면 애플리케이션 재구동

### - `start.sh`

서버 시작

- git pull
- go 패키지 생성
- 현재 구동중인 애플리케이션 확인
  - 실행 중인 애플리케이션이 있다면 종료
- 생성한 패키지 실행

### - `stop.sh`

서버 중단

- 현재 구동중인 애플리케이션 확인
  - 실행 중인 애플리케이션이 있다면 종료

### - `renew.sh`

인증서 갱신

- SSL 인증서 갱신

## 4. 기능 설명

오늘 점심 뭐먹을지 골라주는 웹페이지

현재 위치(지도 상 마커)를 기준으로 상하좌우 거리(m) 만큼의 정사각형 모양 4개 구역에서 음식점을 검색함 (사분면) 

한 구역당 최대 45개, 전체 구역에서 최대 180개의 음식점을 검색하고 그 중 카테고리와 일치하는 음식점을 랜덤하게 하나 골라줌

지도의 GPS 버튼을 누르면 현재 위치를 다시 가져올 수 있음

결과에 나오는 거리는 현재 위치와 음식점 사이의 대략적인 직선 거리

배달 가능 여부, 영업 시간 확인 기능 추가 예정