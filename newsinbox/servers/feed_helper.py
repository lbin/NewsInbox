import json
import feedparser
from newsinbox.common.logger import logger


def feed_parser(config, rss_json_file):
    with open(rss_json_file) as json_file:
        rss_config = json.load(json_file)

    if rss_config["rss_url"].startswith("http") is False:
        url = config.get("rss_base_url") + rss_config["rss_url"]
    else:
        url = rss_config["rss_url"]

    feed = feedparser.parse(url)

    if feed.bozo:
        logger.warn("解析失败")
        return None
    else:
        if rss_config["rss_url"].startswith("http") is True:
            feed_updated = feed.feed.updated
        else:
            feed_updated = feed.updated

        if feed_updated == rss_config["rss_last_updated"]:
            logger.info("{} No New Feed".format(rss_config["rss_name"]))
            return None
        else:
            rss_config["rss_last_updated"] = feed_updated

            update_title_flag = False
            old_title = rss_config["rss_last_updated_title"]
            entries = []

            for entry in feed.entries:
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
                            
                    if hasattr(entry, "published") == False:
                        entry.published = feed_updated
                    entry.rss_name = rss_config["rss_name"] 
                    entries.append(entry)  
            return entries
        return None