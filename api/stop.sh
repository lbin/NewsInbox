ps -ef | grep celery | grep -v grep | awk '{print $2}' | xargs kill -9
ps -ef | grep gunicorn | grep -v grep | awk '{print $2}' | xargs kill -9
rm celery.log summary.log run.log celerybeat-schedule