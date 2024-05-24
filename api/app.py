import json
import os
import threading

from flask import Flask, redirect, render_template, request, session, url_for

if not os.environ.get("DEBUG") or os.environ.get("DEBUG").lower() != "true":
    from gevent import monkey

    monkey.patch_all()
    import grpc.experimental.gevent

    grpc.experimental.gevent.init_gevent()

import logging
import sys
from logging.handlers import RotatingFileHandler

from config import Config

# from events import event_handlers
from extensions import (
    ext_celery,
    ext_compress,
    ext_database,
    ext_migrate,
    ext_redis,
    ext_storage,
)
from extensions.ext_database import db
from flask import Response

# from extensions.ext_login import login_manager


def initialize_extensions(app):
    # Since the application instance is now created, pass it to each Flask
    # extension instance to bind it to the Flask application instance (app)
    ext_compress.init_app(app)
    ext_database.init_app(app)
    ext_migrate.init(app, db)
    ext_redis.init_app(app)
    ext_storage.init_app(app)
    ext_celery.init_app(app)


class NewsInboxApp(Flask):
    pass


def create_app() -> Flask:
    app = NewsInboxApp(__name__)
    app.config.from_object(Config())

    app.secret_key = app.config["SECRET_KEY"]

    log_handlers = None
    log_file = app.config.get("LOG_FILE")
    if log_file:
        log_dir = os.path.dirname(log_file)
        os.makedirs(log_dir, exist_ok=True)
        log_handlers = [
            RotatingFileHandler(filename=log_file, maxBytes=1024 * 1024 * 1024, backupCount=5),
            logging.StreamHandler(sys.stdout),
        ]

    logging.basicConfig(
        level=app.config.get("LOG_LEVEL"),
        format=app.config.get("LOG_FORMAT"),
        datefmt=app.config.get("LOG_DATEFORMAT"),
        handlers=log_handlers,
    )

    initialize_extensions(app)

    return app


# create app
app = create_app()
celery = app.extensions["celery"]

if app.config["TESTING"]:
    print("App is running in TESTING mode")


@app.route("/", methods=["GET", "POST"])
def index():
    if request.method == "GET":
        return render_template("index.html")

    return redirect(url_for("index"))



@app.route("/health")
def health():
    return Response(
        json.dumps({"status": "ok", "version": app.config["CURRENT_VERSION"]}),
        status=200,
        content_type="application/json",
    )


@app.route("/threads")
def threads():
    num_threads = threading.active_count()
    threads = threading.enumerate()

    thread_list = []
    for thread in threads:
        thread_name = thread.name
        thread_id = thread.ident
        is_alive = thread.is_alive()

        thread_list.append({"name": thread_name, "id": thread_id, "is_alive": is_alive})

    return {"thread_num": num_threads, "threads": thread_list}


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=6001)
