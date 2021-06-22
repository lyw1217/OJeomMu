import requests
import os
import json
import configparser

BASE_DIR = os.getcwd()
config_file = os.path.join(BASE_DIR , "config/apiKeys.ini")
config 		= configparser.ConfigParser()
config.read(config_file ,encoding='UTF8')
kakao_api_key       = config['KAKAO']['API_KEY']
google_api_key      = config['GOOGLE']['KEY']
naver_api_id        = config['NAVER']['ID']
naver_api_secert    = config['NAVER']['SECRET']

# 구글 Geolocation API (https://developers.google.com/maps/documentation/geolocation/overview)
url = 'https://www.googleapis.com/geolocation/v1/geolocate?key={}'.format(google_api_key)
data = {
    'considerIp': True,
}

result = requests.post(url, data).json()

print(result)

#json_cord = json.loads(result.text)
glat = result['location']['lat'] # 위도
glng = result['location']['lng'] # 경도
print("구글 Geolocation")
print("google lat(위도, y) = " + str(glat))
print("google lng(경도, x) = " + str(glng))
print("------------------\n")


# 카카오 좌표로 주소 변환하기 (https://developers.kakao.com/docs/latest/ko/local/dev-guide#coord-to-address)
url = 'https://dapi.kakao.com/v2/local/geo/coord2address.json?x={}&y={}'.format(glng,glat)
headers = {
    "Authorization": "KakaoAK " + kakao_api_key
}
places = requests.get(url, headers = headers).json()
print("카카오 좌표로 주소 변환하기")
print(places)
print("------------------\n")


# 네이버 지역(local) 검색 (https://developers.naver.com/docs/search/local/)
searching = '판교 스타벅스'
url = 'https://openapi.naver.com/v1/search/local.json?query={}'.format(searching).encode('UTF-8')
headers = {
    "X-Naver-Client-Id": naver_api_id,
    "X-Naver-Client-Secret": naver_api_secert
}
places = requests.get(url, headers = headers).json()
print("네이버 지역(local) 검색")
print(places)
print("------------------\n")

# 카카오 키워드로 장소 검색 API (https://developers.kakao.com/docs/latest/ko/local/dev-guide#search-by-keyword)
searching = '유스페이스'
cat_grp_code = 'FD6'
url = 'https://dapi.kakao.com/v2/local/search/keyword.json?query={}&category_group_code={}'.format(searching, cat_grp_code)
headers = {
    "Authorization": "KakaoAK " + kakao_api_key
}
places = requests.get(url, headers = headers).json()
klat = places['documents'][0]['y']
klng = places['documents'][0]['x']
print("카카오 키워드로 장소 검색")
print(places)
print("kakao lat(위도, y) = " + str(klat))
print("kakao lng(경도, x) = " + str(klng))
print("------------------\n")

# 카카오 카테고리로 장소 검색 API (https://developers.kakao.com/docs/latest/ko/local/dev-guide#search-by-category)
'''
Name	Description
MT1	    대형마트
CS2	    편의점
PS3	    어린이집, 유치원
SC4	    학교
AC5	    학원
PK6	    주차장
OL7	    주유소, 충전소
SW8	    지하철역
BK9	    은행
CT1	    문화시설
AG2	    중개업소
PO3	    공공기관
AT4	    관광명소
AD5	    숙박
FD6	    음식점
CE7	    카페
HP8	    병원
PM9	    약국
'''
cat_grp_code = 'FD6'
radius = 500
sort = 'distance'
page = 1
url = 'https://dapi.kakao.com/v2/local/search/category.json?category_group_code={}&x={}&y={}&radius={}&sort={}&page={}'.format(cat_grp_code, klng, klat, radius, sort, str(page))
headers = {
    "Authorization": "KakaoAK " + kakao_api_key
}
places = requests.get(url, headers = headers).json()
print("카카오 카테고리로 장소 검색")
print(places)
print("------------------\n")
