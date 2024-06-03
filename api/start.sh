nohup celery -A app.celery worker --loglevel INFO -P gevent &
nohup celery -A app.celery beat --loglevel INFO &
nohup gunicorn  --bind 0.0.0.0:6001 --workers 2 --worker-class gevent  --timeout 200  --preload app:app &