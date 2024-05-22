import datetime
import glob
import json
import os

import app

from .logger import logger


def get_current_data():
    current_date = datetime.datetime.now().strftime("%Y-%m-%d")
    folder_path = f"./schedule/data/{current_date}/focus"

    json_files = []
    for file in glob.glob(os.path.join(folder_path, "*.json")):
        with open(file) as f:
            json_data = json.load(f)
            json_files.append(json_data)

    merged_json = json.dumps(json_files, ensure_ascii=False, indent=4)
    return merged_json


@app.celery.task
def daily_summary_task():
    logger.info("daily_summary_task")
    json_files = get_current_data()
    logger.info(f"{json_files}")
