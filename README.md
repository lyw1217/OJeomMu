# 오늘 점심 뭐 먹지? 오점무? (OJM)

오늘 점심 뭐 먹을지 골라볼까

Domain : https://mumeog.site/

## 1. 개발 환경

- CentOS Linux release 7.9.2009 (Core)
- go version go1.17.3 linux/amd64
  - github.com/gin-gonic/gin v1.7.7
  - github.com/google/go-querystring v1.1.0

## 2. 기능 설명

오늘 점심 뭐 먹을지 골라주는 웹페이지

현재 위치(지도 상 마커) 기준 동서남북으로 거리(m)만큼을 한 변으로 하는 4개의 정사각형 구역에서 음식점을 검색합니다. (사분면, 현재 위치가 (0,0))

한 구역당 최대 45개, 전체 구역에서 최대 180개의 음식점을 검색하고 그 중 카테고리와 일치하는 음식점을 랜덤하게 하나 골라줍니다.

지도의 GPS 버튼을 누르면 현재 위치를 다시 가져올 수 있습니다.

결과에 나오는 거리는 현재 위치와 음식점 사이의 대략적인 직선 거리를 의미합니다.

배달 가능 여부, 영업 시간 확인 기능 추가 예정

### 2.1. 메인 페이지

![main pages](go/assets/img/main_page.jpg)

### 2.2. 설명 페이지

![info pages](go/assets/img/info_page.jpg)

### 2.3. 검색 결과 화면

![result alert](go/assets/img/result_alert.jpg)


## 3. Let's Encrypt(Certbot) 설치 및 인증서 발급

### 3.1. [Certbot 설치](https://certbot.eff.org/instructions)
### 3.2. SSL 인증서 발급 및 설정
    
        $ sudo certbot certonly --standalone -d {domain}

### 3.3. 인증서 갱신 테스트

        $ sudo certbot renew --dry-run

## 4. 스크립트

### 4.1. `monitor.sh`

서버 기동 확인 및 기동

- 현재 구동 중인 애플리케이션 확인
  - 정상 실행 중이라면 스크립트 종료
  - 구동 중인 애플리케이션이 없으면 애플리케이션 재구동

### 4.2. `start.sh`

서버 시작

- 순서
    1. git pull
    2. go 패키지 생성
    3. 현재 구동 중인 애플리케이션 확인
        - 실행 중인 애플리케이션이 있다면 종료
    4. 생성한 패키지 실행

### 4.3. `stop.sh`

서버 중단

- 순서
  1.  현재 구동 중인 애플리케이션 확인
      - 실행 중인 애플리케이션이 있다면 종료

### - `renew.sh`

인증서 갱신

- 순서
    1. 현재 애플리케이션이 구동 중이라면 종료
    2. SSL 인증서 갱신
    3. 애플리케이션 재실행