import requests
import os
import json
import configparser
from flask import Flask, render_template
import config

# https://wikidocs.net/book/4542 참고

def page_not_found(e):
    return render_template('404.html'), 404

def create_app() :

    app = Flask(__name__)

    from .views import main_views
    app.register_blueprint(main_views.bp)

     # 오류페이지
    app.register_error_handler(404, page_not_found)

    return app
'''
# flask run 명령어 수행하면 create_app() 함수가 수행됨 (애플리케이션 팩토리)
if __name__ == '__main__' :
    app = create_app()
    app.run(debug=True)
'''