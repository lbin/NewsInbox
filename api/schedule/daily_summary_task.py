import datetime
import glob
import json
import os

import app

from newsinbox.common.config import conf, load_config
from newsinbox.common.logger import logger
from newsinbox.servers.file_io import save_to_pdf
from newsinbox.servers.llm_common_post import summary_stream

# from newsinbox.servers.wework_bot import send_to_wework_bot


def get_current_data():
    current_date = datetime.datetime.now().strftime("%Y-%m-%d")
    folder_path = f"./schedule/data/{current_date}/focus"

    json_files = []
    for file in glob.glob(os.path.join(folder_path, "*.json")):
        with open(file) as f:
            json_data = json.load(f)
            del json_data["分类"]
            del json_data["标签"]
            del json_data["添加人"]
            del json_data["产品"]
            del json_data["团队"]
            json_files.append(json_data)

    merged_json = json.dumps(json_files, ensure_ascii=False, indent=4)
    # Save merged_json as a JSON file
    output_file = f"./schedule/data/{current_date}/merged.json"
    with open(output_file, "w") as f:
        f.write(merged_json)
    return merged_json


@app.celery.task
def daily_summary_task():
    logger.info("daily_summary_task")
    json_files = get_current_data()
    # logger.info(f"{json_files}")

    load_config("schedule/config.json")
    config = conf()
    config["max_words"] = 0
    config["open_ai_model"] = "moonshot-v1-128k"
    config["prompt"] = (
        "下面引号内的内容是一个json数组, 每一条json代表一个新闻, 包含标题、总结, 关键要点和来源链接等信息。\
        这些新闻来自不同的媒体, 可能会报道相同的内容, 先找出并合并这些新闻。\
        针对合并后的每一条新闻给出一个更好的标题, 用一句话总结全文, 然后给出三个关键的内容点, 并在最后给出文章来源的链接。\
        最终得到的内容需要完全用中文表达, 并组织成markdown的格式  \
        一级标题类似于daily news about sports或todays news about finance之类的 \
        每一条新闻的标题作为二级标题 \
        总结内容高亮, 紧接着列出三个关键要点, 最后给出这些新闻的来源以及链接。\
        请不要遗漏任何一条新闻。内容全部输出, 不需要因为内容多而询问我是否需要全部输出"
    )
    content = summary_stream(config, json_files, len(json_files))
    # content = content.encode('utf-8')
    logger.info(f"{content}")

    # message = {"msgtype": "markdown", "markdown": {"content": content}}
    # send_to_wework_bot(config["we_work_webhook"], message)
    save_to_pdf(content)
