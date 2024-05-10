from celery import Celery, Task
from celery.schedules import crontab
from flask import Flask


def init_app(app: Flask) -> Celery:
    class FlaskTask(Task):
        def __call__(self, *args: object, **kwargs: object) -> object:
            with app.app_context():
                return self.run(*args, **kwargs)

    celery_app = Celery(
        app.name,
        task_cls=FlaskTask,
        broker=app.config["CELERY_BROKER_URL"],
        backend=app.config["CELERY_BACKEND"],
        task_ignore_result=True,
    )

    # Add SSL options to the Celery configuration
    ssl_options = {
        "ssl_cert_reqs": None,
        "ssl_ca_certs": None,
        "ssl_certfile": None,
        "ssl_keyfile": None,
    }

    celery_app.conf.update(
        result_backend=app.config["CELERY_RESULT_BACKEND"],
        broker_connection_retry_on_startup=True,
    )

    if app.config["BROKER_USE_SSL"]:
        celery_app.conf.update(
            broker_use_ssl=ssl_options,  # Add the SSL options to the broker configuration
        )

    celery_app.set_default()
    app.extensions["celery"] = celery_app

    imports = [
        "schedule.crawl_from_rsshub_task",
    ]

    beat_schedule = {
        "crawl_from_rsshub_task": {
            "task": "schedule.crawl_from_rsshub_task.crawl_from_rsshub_task",
            "schedule": crontab(minute=0, hour="1,3,5,7,9,11,13,15,17,19,21,23"),
        }
    }
    celery_app.conf.update(beat_schedule=beat_schedule, imports=imports, timezone="Asia/Shanghai", enable_utc=False)

    return celery_app
