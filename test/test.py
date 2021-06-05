import requests
import os

BASE_DIR = os.getcwd()

with open( BASE_DIR + '/keys/kakaoApi.txt', 'r') as f :
    api_key = f.read()

searching = '합정 스타벅스'
url = 'https://dapi.kakao.com/v2/local/search/keyword.json?query={}'.format(searching)
headers = {
    "Authorization": "KakaoAK " + api_key
}
places = requests.get(url, headers = headers).json()
print(places)
