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
    config['prompt'] = "下面引号内的内容json数据, 每一条json代表今天的一个新闻, 包含标题、内容、链接等信息。这些新闻可能报道了同一件事情, 但是观察的出发点可能不一样。请先找出并合并这些新闻报道的核心内容。然后请返回给我一个markdown文档, 一级标题类似于daily news about sports 或todays news about finance, 然后每一条新闻作为二级标题, 总结内容是这条新闻的核心内容, 在核心内容之后可以给一个点评, 最后列出这些新闻的出处, 出处列出最多三个链接"
    content = sum4all(config, json_files)
    logger.info(f"{content}")
