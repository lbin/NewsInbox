import html
from urllib.parse import urlparse

import requests
from common.logger import logger


def _get_openai_payload(config, target_url_content):
    prompt = config.get('prompt')
    logger.info(f"target_url_content_length: {len(target_url_content)}")
    target_url_content = target_url_content[:config.get('max_words')]  # 通过字符串长度简单进行截断
    sum_prompt = f"{prompt}\n\n'''{target_url_content}'''"
    messages = [{"role": "user", "content": sum_prompt}]
    # payload = {"model": "moonshot-v1-8k", "messages": messages}
    payload = {"model": config.get('open_ai_model'), "messages": messages}
    return payload

def _get_jina_url(target_url):
    jina_reader_base = "https://r.jina.ai"
    return jina_reader_base + "/" + target_url

def sum4all(config, url):
    try:
        target_url = html.unescape(url)
        jina_url = _get_jina_url(target_url)
        # TODO response error handling
        response = requests.get(jina_url, timeout=60)
        response.raise_for_status()
        target_url_content = response.text
        
        
        open_ai_api_base = config.get("open_ai_api_base")
        open_ai_api_key = config.get("open_ai_api_key")
        openai_chat_url = config.get("openai_chat_url")
        
        openai_headers = {"Authorization": f"Bearer {open_ai_api_key}", "Host": urlparse(open_ai_api_base).netloc}
        openai_payload = _get_openai_payload(config, target_url_content)
        response = requests.post(openai_chat_url, headers=openai_headers, json=openai_payload, timeout=60)
        response.raise_for_status()
        result = response.json()["choices"][0]["message"]["content"]
        return result
    except Exception as e:
        logger.error("Error: {}".format(e))
        return None