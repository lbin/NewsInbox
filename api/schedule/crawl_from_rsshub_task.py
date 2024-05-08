import time

import click
from .crawl_rsshub import run

import app

from .logger import logger

@app.celery.task
def crawl_from_rsshub_task():
    run()
