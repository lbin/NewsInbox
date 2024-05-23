import datetime
import glob
import json
import os

import app

from newsinbox.common.logger import logger
from newsinbox.servers.llm_common_post import sum4all
from newsinbox.common.config import conf, load_config

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
    # logger.info(f"{json_files}")
    
    load_config("schedule/config.json")
    config = conf()
    config['max_words'] = 0
    config['open_ai_model'] = "moonshot-v1-128k"
    config['prompt'] = '我需要对下面引号内的字符串进行归纳整理, 合并相同的内容, 并用一句话说明合并的内容要点'
    content = sum4all(config, json_files)
    logger.info(f"{content}")
