import os
import configparser

BASE_DIR = os.getcwd()
#BASE_DIR = os.path.dirname(os.path.dirname(__file__))
CONFIG_FILE = os.path.join(BASE_DIR , "apiKeys.ini")

config 		= configparser.ConfigParser()
config.read(CONFIG_FILE ,encoding='UTF8')

KAKAO_API_KEY       = config['KAKAO']['API_KEY']
KAKAO_JS_KEY        = config['KAKAO']['JS_KEY']
GOOGLE_API_KEY      = config['GOOGLE']['KEY']
NAVER_API_ID        = config['NAVER']['ID']
NAVER_API_SECERT    = config['NAVER']['SECRET']