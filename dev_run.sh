#!/bin/bash
cd ~/app/OJeomMu
export APP_CONFIG_FILE=/home/ec2-user/app/OJeomMu/config/production.py
export FLASK_ENV=development
export FLASK_APP=app
source ~/app/OJeomMu/venv/bin/activate
echo "FLASK RUN [DEV]"
sleep 1
flask run
