nohup celery -A app.celery worker --loglevel INFO -P gevent &
nohup celery -A app.celery beat --loglevel INFO &