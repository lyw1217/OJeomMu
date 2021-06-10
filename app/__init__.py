import requests
import os
import json
import configparser
from flask import Flask

# https://wikidocs.net/book/4542 참고

BASE_DIR = os.getcwd()
config_file = BASE_DIR + "/apiKeys.ini"
config 		= configparser.ConfigParser()
config.read(config_file , encoding='UTF8')
kakao_api_key       = config['KAKAO']['API_KEY']
google_api_key      = config['GOOGLE']['KEY']
naver_api_id        = config['NAVER']['ID']
naver_api_secert    = config['NAVER']['SECRET']

def create_app() :

    app = Flask(__name__)

    @app.route('/')
    @app.route('/home')
    def home():
        return '''
        <h1>이건 h1 제목</h1>
        <p>이건 p 본문 </p>
        <a href="https://flask.palletsprojects.com">Flask 홈페이지 바로가기</a>
        '''

    @app.route('/user/<user_name>/<int:user_id>')
    def user(user_name, user_id):
        return f'Hello, {user_name}({user_id})!'

    return app
'''
# flask run 명령어 수행하면 create_app() 함수가 수행됨 (애플리케이션 팩토리)
if __name__ == '__main__' :
    app = create_app()
    app.run(debug=True)
'''