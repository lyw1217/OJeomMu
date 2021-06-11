from flask import Blueprint, render_template
from config.default import *

bp = Blueprint('main', __name__, url_prefix='/')

@bp.route('/')
@bp.route('/main')
def home():
    key = KAKAO_JS_KEY
    return render_template('main.html', key=key)

@bp.route('/user/<user_name>/<int:user_id>')
def user(user_name, user_id):
    return f'Hello, {user_name}({user_id})!'
