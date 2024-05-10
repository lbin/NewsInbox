ps -ef | grep celery | grep -v grep | awk '{print $2}' | xargs kill -9
rm nohup.out celerybeat-schedule