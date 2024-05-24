import html
from urllib.parse import urlparse

import requests
from openai import OpenAI

from ..common.logger import logger


def _get_openai_payload(config, content, max_tokens=0):
    prompt = config.get("prompt")
    logger.info(f"target_url_content_length: {len(content)}")
    max_words = config.get("max_words")
    if max_words > 0:
        content = content[:max_words]  # 通过字符串长度简单进行截断
    sum_prompt = f"{prompt}\n\n'''{content}'''"
    messages = [{"role": "user", "content": sum_prompt}]
    if max_tokens == 0:
        payload = {"model": config.get("open_ai_model"), "messages": messages, "temperature": 0.3}
    else:
        payload = {
            "model": config.get("open_ai_model"),
            "messages": messages,
            "temperature": 0.3,
            "max_tokens": max_tokens,
        }
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


def summary_url_content(config, content, max_tokens=0):
    open_ai_api_base = config.get("open_ai_api_base")
    open_ai_api_key = config.get("open_ai_api_key")
    openai_chat_url = config.get("openai_chat_url")
    model = config.get("open_ai_model")

    prompt = config.get("prompt")
    logger.info(f"target_url_content_length: {len(content)}")

    sum_prompt = f"{prompt}\n\n'''{content}'''"
    
    client = OpenAI(
        api_key=open_ai_api_key,
        base_url=open_ai_api_base,
    )

    completion = client.chat.completions.create(
        model=model,
        messages=[
            {
                "role": "system",
                "content": "请从产品描述中提取名称、尺寸、价格和颜色, 并在一个 JSON 对象中输出。",
            },
            {
                "role": "user",
                "content": sum_prompt,
            },
            {"role": "assistant", "content": "{", "partial": True},
        ],
        temperature=0.3,
    )

    print("{" + completion.choices[0].message.content)


def sum4all(config, content, max_tokens=0):
    try:
        open_ai_api_base = config.get("open_ai_api_base")
        open_ai_api_key = config.get("open_ai_api_key")
        openai_chat_url = config.get("openai_chat_url")

        openai_headers = {"Authorization": f"Bearer {open_ai_api_key}", "Host": urlparse(open_ai_api_base).netloc}
        openai_payload = _get_openai_payload(config, content, max_tokens)
        response = requests.post(openai_chat_url, headers=openai_headers, json=openai_payload, timeout=240)
        response.raise_for_status()
        result = response.json()["choices"][0]["message"]["content"]
        logger.info(response.json()["usage"])
        logger.info(response.json()["choices"][0]["finish_reason"])
        return result
    except Exception as e:
        logger.error("Error: {}".format(e))
        return None


def summary_stream(config, content, max_tokens=0):
    open_ai_api_base = config.get("open_ai_api_base")
    open_ai_api_key = config.get("open_ai_api_key")
    openai_chat_url = config.get("openai_chat_url")
    model = config.get("open_ai_model")

    prompt = config.get("prompt")
    logger.info(f"target_url_content_length: {len(content)}")

    sum_prompt = f"{prompt}\n\n'''{content}'''"

    client = OpenAI(
        api_key=open_ai_api_key,
        base_url=open_ai_api_base,
    )

    response = client.chat.completions.create(
        model=model,
        messages=[
            {
                "role": "system",
                "content": "你是 Kimi, 由 Moonshot AI 提供的人工智能助手, 你更擅长中文和英文的对话。你会为用户提供安全, 有帮助, 准确的回答。同时, 你会拒绝一切涉及恐怖主义, 种族歧视, 黄色暴力等问题的回答。Moonshot AI 为专有名词, 不可翻译成其他语言。",
            },
            {"role": "user", "content": sum_prompt},
        ],
        temperature=0.3,
        stream=True,
        max_tokens=max_tokens,
    )

    collected_messages = []
    for idx, chunk in enumerate(response):
        # print("Chunk received, value: ", chunk)
        chunk_message = chunk.choices[0].delta
        if not chunk_message.content:
            continue
        collected_messages.append(chunk_message)  # save the message
        # print(f"#{idx}: {chunk_message}")
        # print(f"#{idx}: {''.join([m.content for m in collected_messages])}")
    # print(f"Full conversation received: {''.join([m.content for m in collected_messages])}")
    # logger.info(response.json()["usage"])
    # logger.info(response.json()["choices"][0]["finish_reason"])
    messages = "".join([m.content for m in collected_messages])
    return messages
