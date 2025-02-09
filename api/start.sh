cd ..
python3 setup.py  install --user
cd api
nohup celery -A app.celery worker --loglevel INFO -P gevent >> celery.log 2>&1 &
nohup celery -A app.celery beat --loglevel INFO >> celery.log 2>&1 &
nohup gunicorn  --bind 0.0.0.0:6001 --workers 2 --worker-class gevent  --timeout 200  --preload app:app >> summary.log 2>&1 &