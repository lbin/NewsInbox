import time

import click
from crawl_rsshub import run

import app


@app.celery.task()
def crawl_from_rsshub_task():
    click.echo(click.style("Start Crawl From RSSHub.", fg="green"))
    start_at = time.perf_counter()
    
    run()
    
    end_at = time.perf_counter()
    click.echo(
        click.style(
            "Crawled From RSSHub success latency: {}".format(end_at - start_at),
            fg="green",
        )
    )
