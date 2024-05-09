import html
import json
import time
from urllib.parse import urlparse

import feedparser
import requests

import app

from .config import conf, load_config, save_config
from .logger import logger


def load_key_words():
    # 读取key_words.txt, 返回关键词列表
    with open("schedule/key_words.txt") as f:
        key_words = f.readlines()
        key_words = [word.strip() for word in key_words]
    return key_words


def send_to_feishu(content, key_words, url, sender, title=None):
    feishu_get_token_url = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"
    feishu_add_record_url = (
        "https://open.feishu.cn/open-apis/bitable/v1/apps/DM8Ib7DNeah45XsA8kOcxvYenHd/tables/tblELLHLYuKXsVf9/records"
    )
    try:
        logger.info("content: {}".format(content))
        data_index = content.find("json")

        json_data = content[data_index + 3 : -3]
        json_data = json_data.replace("\n", "")
        json_data = json_data.replace(" ", "")
        json_data = json_data[1:]
        

        json_data = json.loads(json_data)
        config = conf()

        headers = {"Content-Type": "application/json; charset=utf-8"}
        token_json_data = {"app_id": config.get('feishu_app_id'), "app_secret": config.get('feishu_app_secret')}
        feishu_token_response = requests.post(feishu_get_token_url, headers=headers, json=token_json_data)

        feishu_headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {feishu_token_response.json()['app_access_token']}",
        }

        key_points_str = ""
        for key in json_data["key_points"]:
            key_points_str += json_data["key_points"][key] + "\n"

        send_title = title if title else json_data["summary"]

        new_json_data = {
            "分类": "Technology",
            "标签": json_data["tags"],
            "项目名称": send_title,
            "来源": url,
            "总结": json_data["summary"],
            "关键要点": key_points_str,
            "项目来源": sender,
            "添加人": "Bot",
        }

        tags = json_data["tags"]

        for key_word in key_words:
            for tag in tags:
                if key_word in tag:
                    logger.info("关键词匹配成功: key_word-{} tag-{}".format(key_word, tag))

                    new_data = {"fields": new_json_data}
                    status = requests.post(feishu_add_record_url, headers=feishu_headers, json=new_data)
                    return status
    except Exception as e:
        logger.error("Error: {}".format(e))
        return None

    return None


def _get_jina_url(target_url):
    jina_reader_base = "https://r.jina.ai"
    return jina_reader_base + "/" + target_url


def _get_openai_payload(target_url_content):
    config = conf()
    prompt = config.get('prompt')
    target_url_content = target_url_content[:8000]  # 通过字符串长度简单进行截断
    sum_prompt = f"{prompt}\n\n'''{target_url_content}'''"
    messages = [{"role": "user", "content": sum_prompt}]
    # payload = {"model": "moonshot-v1-8k", "messages": messages}
    payload = {"model": config.get('open_ai_model'), "messages": messages}
    return payload


def sum4all(url):
    try:
        target_url = html.unescape(url)
        jina_url = _get_jina_url(target_url)
        response = requests.get(jina_url, timeout=60)
        response.raise_for_status()
        target_url_content = response.text
        
        config = conf()
        
        # open_ai_api_base = "https://api.moonshot.cn/v1"
        # open_ai_api_key = "sk-I10kI0VKjld1LDx5AwFOHP4YSHF5rGkkDBaBTa5IEiNEnbiE"
        # openai_chat_url = "https://api.moonshot.cn/v1/chat/completions"
        
        open_ai_api_base = config.get("open_ai_api_base")
        open_ai_api_key = config.get("open_ai_api_key")
        openai_chat_url = config.get("openai_chat_url")
        
        openai_headers = {"Authorization": f"Bearer {open_ai_api_key}", "Host": urlparse(open_ai_api_base).netloc}
        openai_payload = _get_openai_payload(target_url_content)
        response = requests.post(openai_chat_url, headers=openai_headers, json=openai_payload, timeout=60)
        response.raise_for_status()
        result = response.json()["choices"][0]["message"]["content"]
        return result
    except Exception as e:
        logger.error("Error: {}".format(e))
        return None


def feed_parser(rss, key_words, sender):
    feed = feedparser.parse(rss["rss_url"])

    new_rss = rss

    if feed.bozo:
        logger.warn("解析失败")
        return rss
    else:
        logger.info("{}".format(feed.updated))

        if feed.updated == rss["rss_last_updated"]:
            logger.info("No new feed")
            return rss
        else:
            new_rss["rss_last_updated"] = feed.updated

            update_title_flag = False
            old_title = rss["rss_last_updated_title"]

            for entry in feed.entries:
                logger.info("标题: {}".format(entry.title))
                logger.info("链接: {}".format(entry.link))
                if entry.title == old_title:
                    logger.info("已经解析过")
                    return new_rss
                else:
                    if update_title_flag is False:
                        new_rss["rss_last_updated_title"] = entry.title
                        update_title_flag = True
                    content = sum4all(entry.link)
                    send_to_feishu(content, key_words, entry.link, sender, entry.title)
                    time.sleep(60 * 2)
            return new_rss

@app.celery.task
def crawl_from_rsshub_task():
    key_words = load_key_words()
    load_config()
    config = conf()
    rss_group = config.get("rss_group", {})
    new_rss_group = []
    for rss in rss_group:
        new_rss = feed_parser(rss, key_words, rss["rss_name"])
        new_rss_group.append(new_rss)
        time.sleep(60 * 2)
    config["rss_group"] = new_rss_group
    save_config()
    


if __name__ == "__main__":
    while True:
        crawl_from_rsshub_task()
        time.sleep(60 * 10 * 3)
