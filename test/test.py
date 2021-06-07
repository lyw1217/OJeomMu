import requests
import os
import json
import configparser

BASE_DIR = os.getcwd()

config_file = BASE_DIR + "/keys/apiKeys.ini"
config 		= configparser.ConfigParser()
config.read(config_file ,encoding='UTF8')
kakao_api_key       = config['KAKAO']['KEY']
google_api_key      = config['GOOGLE']['KEY']
naver_api_id        = config['NAVER']['ID']
naver_api_secert    = config['NAVER']['SECRET']

# 카카오 키워드로 장소 검색 API
searching = '합정 스타벅스'
url = 'https://dapi.kakao.com/v2/local/search/keyword.json?query={}'.format(searching)
headers = {
    "Authorization": "KakaoAK " + kakao_api_key
}
places = requests.get(url, headers = headers).json()
print(places)
print("------------------\n")

# 구글 Geolocation API
url = f'https://www.googleapis.com/geolocation/v1/geolocate?key={google_api_key}'
data = {
    'considerIp': True,
}

result = requests.post(url, data)

print(result.text)
print("------------------\n")

json_cord = json.loads(result.text)
lat = json_cord['location']['lat'] # 위도
lng = json_cord['location']['lng'] # 경도
print("lat = " + str(lat))
print("lng = " + str(lng))
print("------------------\n")

# 카카오 좌표로 주소 변환하기
searching = '합정 스타벅스'
url = 'https://dapi.kakao.com/v2/local/geo/coord2address.json?x={}&y={}'.format(lng,lat)
headers = {
    "Authorization": "KakaoAK " + kakao_api_key
}
places = requests.get(url, headers = headers).json()
print(places)
print("------------------\n")

# 네이버 지역(local) 검색
searching = '합정 스타벅스'
url = 'https://openapi.naver.com/v1/search/local.json?query={}'.format(searching).encode('UTF-8')
headers = {
    "X-Naver-Client-Id": naver_api_id,
    "X-Naver-Client-Secret": naver_api_secert
}
places = requests.get(url, headers = headers).json()
print(places)
print("------------------\n")