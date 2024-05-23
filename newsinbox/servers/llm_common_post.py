import html
from urllib.parse import urlparse

import requests

from ..common.logger import logger


def _get_openai_payload(config, content):
    prompt = config.get('prompt')
    logger.info(f"target_url_content_length: {len(content)}")
    max_words = config.get('max_words')
    if max_words > 0:
        content = content[:max_words]  # 通过字符串长度简单进行截断
    sum_prompt = f"{prompt}\n\n'''{content}'''"
    messages = [{"role": "user", "content": sum_prompt}]
    payload = {"model": config.get('open_ai_model'), "messages": messages}
    return payload

def _get_jina_url(target_url):
    jina_reader_base = "https://r.jina.ai"
    return jina_reader_base + "/" + target_url

def get_url_content(url):
    try:
        target_url = html.unescape(url)
        jina_url = _get_jina_url(target_url)
        response = requests.get(jina_url, timeout=120)
        response.raise_for_status()
        return response.text
    except Exception as e:
        logger.error("Error: {}".format(e))
        return None

def sum4all(config, content):
    try:
        open_ai_api_base = config.get("open_ai_api_base")
        open_ai_api_key = config.get("open_ai_api_key")
        openai_chat_url = config.get("openai_chat_url")
        
        openai_headers = {"Authorization": f"Bearer {open_ai_api_key}", "Host": urlparse(open_ai_api_base).netloc}
        openai_payload = _get_openai_payload(config, content)
        response = requests.post(openai_chat_url, headers=openai_headers, json=openai_payload, timeout=120)
        response.raise_for_status()
        result = response.json()["choices"][0]["message"]["content"]
        return result
    except Exception as e:
        logger.error("Error: {}".format(e))
        return None