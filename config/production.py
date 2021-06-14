from config.default import *

SQLALCHEMY_DATABASE_URI = 'sqlite:///{}'.format(os.path.join(BASE_DIR, 'ojeommu.db'))
SQLALCHEMY_TRACK_MODIFICATIONS = False
SECRET_KEY = b'b\xdf\n\xbe\x14X\xaal\x96\x1dbd\x92\xe8\x96\x89'