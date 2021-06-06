import requests
import os
import json

BASE_DIR = os.getcwd()

with open( BASE_DIR + '/keys/kakaoApi.txt', 'r') as f :
    kakao_api_key = f.read()
with open( BASE_DIR + '/keys/googleApi.txt', 'r') as f :
    google_api_key = f.read()

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

# 카카오 좌표로 주소 변환하기
searching = '합정 스타벅스'
url = 'https://dapi.kakao.com/v2/local/geo/coord2address.json?x=%s&y=%s'%(lng,lat)
headers = {
    "Authorization": "KakaoAK " + kakao_api_key
}
places = requests.get(url, headers = headers).json()
print(places)
print("------------------")