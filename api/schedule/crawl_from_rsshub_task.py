import json
import time

import app
import feedparser

from newsinbox.common.config import conf, load_config
from newsinbox.common.logger import logger
from newsinbox.servers.feishu_helper import send_to_feishu
from newsinbox.servers.file_io import load_key_words
from newsinbox.servers.llm_common_post import get_url_content, sum4all


def feed_parser(rss, key_words, black_words, key_web_list):
    config = conf()

    rss_json_file = "schedule/rss_configs/" + rss
    with open(rss_json_file) as json_file:
        rss_config = json.load(json_file)

    if rss_config["rss_url"].startswith("http") is False:
        url = config.get("rss_base_url") + rss_config["rss_url"]
    else:
        url = rss_config["rss_url"]
    
    if rss_config["rss_name"] in key_web_list:
        key_words = None
        black_words = None

    feed = feedparser.parse(url)

    if feed.bozo:
        logger.warn("解析失败")
        return rss
    else:
        if rss_config["rss_url"].startswith("http") is True:
            feed_updated = feed.feed.updated
        else:
            feed_updated = feed.updated

        # logger.info("{}".format(feed_updated))

        if feed_updated == rss_config["rss_last_updated"]:
            logger.info("{} No New Feed".format(rss_config["rss_name"]))
            return rss
        else:
            rss_config["rss_last_updated"] = feed_updated

            update_title_flag = False
            old_title = rss_config["rss_last_updated_title"]

            for entry in feed.entries:
                logger.info("标题: {}".format(entry.title))
                logger.info("链接: {}".format(entry.link))
                if entry.title == old_title:
                    logger.info("old-标题: {}".format(old_title))
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
                    raw_content = get_url_content(entry.link)
                    content = sum4all(config, raw_content)
                    published = entry.published if hasattr(entry, "published") else feed_updated
                    send_to_feishu(config, content, key_words, black_words, entry.link, rss_config["rss_name"], published, raw_content)
                    time.sleep(config.get("rss_interval"))

            return None


@app.celery.task
def crawl_from_rsshub_task():
    key_words, black_words, key_web_list = load_key_words()

    load_config("schedule/config.json")
    config = conf()
    rss_group = config.get("rss_group", {})

    for rss in rss_group:
        logger.info("RSS Source: {}".format(rss))
        feed_parser(rss, key_words, black_words, key_web_list)
        time.sleep(config.get("rss_interval"))
