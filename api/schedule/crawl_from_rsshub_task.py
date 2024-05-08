
import app
from .crawl_rsshub import run

@app.celery.task
def crawl_from_rsshub_task():
    run()
