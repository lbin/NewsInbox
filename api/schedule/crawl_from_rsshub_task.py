import datetime
import html
import json
import os
import time
from urllib.parse import urlparse

import feedparser
import requests

import app

from .config import conf, load_config
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
        
        tags = json_data["tags"]
        logger.info("#tags: {}".format(tags))
        
        new_tags = []
        
        for tag in tags:
            tag = tag.replace("#", "")
            new_tags.append(tag)
        tags = new_tags
        logger.info("tags: {}".format(tags))
        
        new_json_data = {
            "分类": "Technology",
            "标签": tags,
            "项目名称": send_title,
            "来源": url,
            "总结": json_data["summary"],
            "关键要点": key_points_str,
            "项目来源": sender,
            "添加人": "Bot",
        }
        
        # Generate file name based on current time and sender
        current_date = datetime.datetime.now().strftime("%Y-%m-%d")
        folder_path = f"./schedule/data/{current_date}"
        if not os.path.exists(folder_path):
            os.makedirs(folder_path)
        file_name = f"{folder_path}/{datetime.datetime.now().strftime('%H-%M-%S')}_{sender}.json"
        # Save new_json_data to json file
        with open(file_name, "w") as json_file:
            json.dump(new_json_data, json_file, ensure_ascii=False, indent=4)

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
    target_url_content = target_url_content[:config.get('max_words')]  # 通过字符串长度简单进行截断
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


def feed_parser(rss, key_words):
    config = conf()
    
    rss_json_file = 'schedule/rss_configs/'+ rss
    with open(rss_json_file) as json_file:
        rss_config = json.load(json_file)
        
    if rss_config["rss_url"].startswith("http") is False:
        url = config.get("rss_base_url") + rss_config["rss_url"]
    else:
        url = rss_config["rss_url"]
        
    feed = feedparser.parse(url)

    if feed.bozo:
        logger.warn("解析失败")
        return rss
    else:
        if rss_config["rss_url"].startswith("http") is True:
            feed_updated = feed.feed.updated
        else:
            feed_updated = feed.updated
 
        logger.info("{}".format(feed_updated))
        

        if feed_updated == rss_config["rss_last_updated"]:
            logger.info("{} No New Feed".format(rss_config['rss_name']))
            return rss
        else:
            rss_config["rss_last_updated"] = feed_updated

            update_title_flag = False
            old_title = rss_config["rss_last_updated_title"]

            for entry in feed.entries:
                entry.title = entry.title.replace(" ", "")
                logger.info("标题: {}".format(entry.title))
                logger.info("链接: {}".format(entry.link))
                if entry.title == old_title:
                    logger.info("已经解析过")
                    with open(rss_json_file, "w") as json_file:
                        json.dump(rss_config, json_file, ensure_ascii=False, indent=4)
                    return None
                else:
                    if update_title_flag is False:
                        rss_config["rss_last_updated_title"] = entry.title
                        update_title_flag = True
                        with open(rss_json_file, "w") as json_file:
                            json.dump(rss_config, json_file, ensure_ascii=False, indent=4)
                    content = sum4all(entry.link)
                    send_to_feishu(content, key_words, entry.link, rss_config['rss_name'], entry.title)
                    time.sleep(config.get("rss_interval"))

            return None

@app.celery.task
def crawl_from_rsshub_task():
    key_words = load_key_words()
    
    load_config()
    config = conf()
    rss_group = config.get("rss_group", {})
    
    for rss in rss_group:
        logger.info("RSS Source: {}".format(rss))
        feed_parser(rss, key_words)
        time.sleep(config.get("rss_interval"))
