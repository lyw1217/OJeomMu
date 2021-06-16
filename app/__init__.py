# https://wikidocs.net/book/4542 참고

from flask import Flask, render_template
from flask_migrate import Migrate
from flask_sqlalchemy import SQLAlchemy
import config

# SQLAlchemy 적용하기
db = SQLAlchemy()
migrate = Migrate()

def page_not_found(e):
    return render_template('404.html'), 404

def server_error(e):
    return render_template('500.html'), 500

def create_app() :

    app = Flask(__name__)

    app.config.from_object(config)

    # ORM
    db.init_app(app)
    migrate.init_app(app, db)
    
    from . import models

    app.config.from_envvar('APP_CONFIG_FILE')
    
    # 블루프린트 적용
    from .views import main_views, question_views, answer_views, auth_views
    app.register_blueprint(main_views.bp)
    app.register_blueprint(question_views.bp)
    app.register_blueprint(answer_views.bp)
    app.register_blueprint(auth_views.bp)

    # 필터
    from .filter import format_datetime
    app.jinja_env.filters['datetime'] = format_datetime

    # 오류페이지
    app.register_error_handler(404, page_not_found)
    app.register_error_handler(500, server_error)

    return app

'''
# flask run 명령어 수행하면 create_app() 함수가 수행됨 (애플리케이션 팩토리)
if __name__ == '__main__' :
    app = create_app()
    app.run(debug=True)
'''