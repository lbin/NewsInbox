import time

import app

from newsinbox.common.config import conf, load_config
from newsinbox.common.logger import logger
from newsinbox.servers.feishu_helper import send_to_feishu_4_venture
from newsinbox.servers.llm_common_post import get_url_content, sum4all
from newsinbox.servers.feed_helper import feed_parser


@app.celery.task
def crawl4venture():
    load_config("schedule/config_venture.json")
    config = conf()

    rss_group = config.get("rss_group", {})
    for rss in rss_group:
        logger.info("RSS Source: {}".format(rss))

        entries = feed_parser(config, "schedule/rss_configs/" + rss)
        if entries is not None:
            for entry in entries:
                raw_content = get_url_content(entry.link)

                content = sum4all(config, raw_content)
                # TODO Source
                send_to_feishu_4_venture(config, content, raw_content)

                time.sleep(config.get("rss_interval"))

        time.sleep(config.get("rss_interval"))
