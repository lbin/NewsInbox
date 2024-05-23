import json

import requests
from common.logger import logger

from .file_io import get_timestamp, save_to_json


def _send_to_feishu(config, data):
    feishu_get_token_url = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"
    feishu_add_record_url = (
        "https://open.feishu.cn/open-apis/bitable/v1/apps/DM8Ib7DNeah45XsA8kOcxvYenHd/tables/tblELLHLYuKXsVf9/records"
    )
    
    headers = {"Content-Type": "application/json; charset=utf-8"}
    token_json_data = {"app_id": config.get("feishu_app_id"), "app_secret": config.get("feishu_app_secret")}
    
    feishu_token_response = requests.post(feishu_get_token_url, headers=headers, json=token_json_data)
    # TODO response error handling
    feishu_headers = {
        "Content-Type": "application/json; charset=utf-8",
        "Authorization": f"Bearer {feishu_token_response.json()['app_access_token']}",
    }
    # TODO response error handling
    status = requests.post(feishu_add_record_url, headers=feishu_headers, json=data)
    return status

def send_to_feishu(config, content, key_words, black_words, url, sender, title=None, published=None):
    logger.info("content: {}".format(content))
    if content is None:
        return None
    if content[0] != "{":
        data_index = content.find("json")
        json_data = content[data_index + 3 : -3]
        json_data = json_data.replace("\n", "")
        # json_data = json_data.replace(" ", "")
        json_data = json_data[1:]
    else:
        json_data = content

    try:
        json_data = json.loads(json_data)
    except Exception as e:
        logger.error("Error: {}\n Content: {}".format(e, json_data))
        return None


    key_points_str = ""
    for key in json_data["key_points"]:
        key_points_str += json_data["key_points"][key] + "\n"

    send_title = json_data["title"]

    tags = json_data["tags"]
    new_tags = []
    for tag in tags:
        tag = tag.replace("#", "")
        new_tags.append(tag)
    tags = new_tags
    
    logger.info("tags: {}".format(tags))
    logger.info("published: {}".format(published))

    new_json_data = {
        "分类": json_data["class"],
        "标签": tags,
        "项目名称": send_title,
        "来源": url,
        "总结": json_data["summary"],
        "关键要点": key_points_str,
        "项目来源": sender,
        "添加人": "Bot",
        "产品": json_data["products"],
        "团队": json_data["teams"],
        "发布时间": get_timestamp(published),
    }

    for black_word in black_words:
        for tag in tags:
            if black_word in tag:
                logger.info("黑名单匹配成功: black_word-{} tag-{}".format(black_word, tag))
                save_to_json(new_json_data, sender, False)
                return None

    for key_word in key_words:
        for tag in tags:
            if key_word in tag:
                logger.info("关键词匹配成功: key_word-{} tag-{}".format(key_word, tag))
                save_to_json(new_json_data, sender)
                new_data = {"fields": new_json_data}
                status = _send_to_feishu(config, new_data)
                return status

    save_to_json(new_json_data, sender, False)
    return None
