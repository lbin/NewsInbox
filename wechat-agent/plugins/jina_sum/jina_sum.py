import html
import json
import os
from urllib.parse import urlparse

import requests
from bridge.context import ContextType
from bridge.reply import Reply, ReplyType
from common.log import logger

import plugins
from plugins import *


@plugins.register(
    name="JinaSum",
    desire_priority=10,
    hidden=False,
    desc="Sum url link content with jina reader and llm",
    version="0.0.1",
    author="hanfangyuan",
)
class JinaSum(Plugin):
    jina_reader_base = "https://r.jina.ai"
    open_ai_api_base = "https://api.moonshot.cn/v1"
    open_ai_model = "moonshot-v1-8k"

    max_words = 8000
    prompt = "我需要对下面引号内文档进行总结并对涉及的产品和团队信息进行抽取, 根据以下标准将输入的文本数据进行分类：软件创新、硬件创新、商业应用、教育与学习、医疗健康、心理健康与情绪陪伴、社会影响、人机交互、老龄以及其他。每个类别的定义和重点如下, [软件创新：涉及新的编程模型、软件架构、算法优化、AI应用开发框架等],[硬件创新：包括新的芯片设计、传感器技术、可穿戴设备和任何增强计算或数据处理能力的物理设备。],[商业应用：涉及AI技术的市场推广，包括产品落地、市场策略、商业模式的创新等。],[教育与学习：包括利用AI技术进行教育创新，如智能教学系统、在线教育平台、学习管理系统等。],[医疗健康：侧重于AI在传统医疗和健康监测中的应用。],[心理健康与情绪陪伴：专注于AI在情绪支持、心理健康干预和治疗等方面的应用。],[社会影响：探讨AI技术如何影响社会结构、伦理、隐私和安全等问题。],[人机交互：关注改进的用户交互体验，如多模态交互和自然语言处理。],[老龄：专门为老年人设计的技术，如健康管理和生活辅助。],[其他：用于分类不符合以上任何一个类别的信息。],请确保分类结果的准确性，并注意某些文本可能需要归类到多个标签。总结内容全部用中文表述, 总结输出包括以下七个部分:\n标题\n一句话总结\n关键要点,用数字序号列出3-5个文章的核心内容\n所属类别:xx, xx\n标签: xx, xx, xx\n涉及的产品: xx, xx, xx\n涉及的公司或团队:xx, xx. 标签的提取不要受到所属类别的影响，类别限定在所提出的类别中，标签要根据全文内容自由提取。结果以JSON格式返回: {title: 标题, summary: 总结的内容, class:[所属类别1, 所属类别2], key_points: {'1': 关键要点1, '2': 关键要点, '3': 关键要点}, tags: [xx, xx, xx], products:[xx, xx, xx], teams:[xx, xx]}, 严格按照json的格式返回回答, json中只有key、value的形式,不要带有其他的字符串和信息"
    white_url_list = []
    black_url_list = [
        "https://support.weixin.qq.com",  # 视频号视频
        "https://channels-aladin.wxqcloud.qq.com",  # 视频号音乐
    ]

    def __init__(self):
        super().__init__()
        try:
            self.config = super().load_config()
            if not self.config:
                self.config = self._load_config_template()
            self.jina_reader_base = self.config.get("jina_reader_base", self.jina_reader_base)
            self.open_ai_api_base = self.config.get("open_ai_api_base", self.open_ai_api_base)
            self.open_ai_api_key = self.config.get("open_ai_api_key", "")
            self.open_ai_model = self.config.get("open_ai_model", self.open_ai_model)
            self.max_words = self.config.get("max_words", self.max_words)
            self.prompt = self.config.get("prompt", self.prompt)
            self.white_url_list = self.config.get("white_url_list", self.white_url_list)
            self.black_url_list = self.config.get("black_url_list", self.black_url_list)
            logger.info("[JinaSum] inited")
            self.handlers[Event.ON_HANDLE_CONTEXT] = self.on_handle_context
        except Exception as e:
            logger.error(f"[JinaSum] 初始化异常：{e}")
            raise "[JinaSum] init failed, ignore "

    def on_handle_context(self, e_context: EventContext, retry_count: int = 0):
        try:
            context = e_context["context"]
            content = context.content
            logger.info(f"[JinaSum] on_handle_context. context: {context}")
            pdf_flag = False
            if context.type == ContextType.FILE and content.endswith(".pdf"):
                pdf_flag = True
            else:
                if context.type != ContextType.SHARING and context.type != ContextType.TEXT:
                    return
                else:
                    if not self._check_url(content):
                        logger.debug(f"[JinaSum] {content} is not a valid url, skip")
                        return

            if retry_count == 0:
                logger.debug("[JinaSum] on_handle_context. content: {}".format(content))
                reply = Reply(ReplyType.TEXT, "🎉正在为您生成总结，请稍候...")
                channel = e_context["channel"]
                channel.send(reply, context)

            if pdf_flag:
                paths = content.split("/")
                target_url = "https://halfjourney.xyz/pdf/" + paths[-2] + "/" + paths[-1]
                target_url = html.unescape(target_url)
            else:
                target_url = html.unescape(content)  # 解决公众号卡片链接校验问题，参考 https://github.com/fatwang2/sum4all/commit/b983c49473fc55f13ba2c44e4d8b226db3517c45
            logger.info(f"[JinaSum] target_url: {target_url}")
            jina_url = self._get_jina_url(target_url)
            # response = requests.get(jina_url, timeout=120)
            # response.raise_for_status()
            
            import httpx
            from fake_useragent import UserAgent
            ua = UserAgent()
            ua.random
            
            headers = {
                "User-Agent": ua.chrome,
            }

            response = httpx.get(url=jina_url, headers=headers, timeout=120, verify=False)
            response.raise_for_status()
            target_url_content = response.text

            openai_chat_url = self._get_openai_chat_url()
            openai_headers = self._get_openai_headers()
            openai_payload = self._get_openai_payload(target_url_content)
            logger.debug(f"[JinaSum] openai_chat_url: {openai_chat_url}, openai_headers: {openai_headers}, openai_payload: {openai_payload}")
            response = requests.post(openai_chat_url, headers=openai_headers, json=openai_payload, timeout=60)
            response.raise_for_status()
            result = response.json()["choices"][0]["message"]["content"]
            if result.startswith("```"):
                pass
            else:
                result = "```json" + result + "```"
            reply = Reply(ReplyType.TEXT, result, target_url)
            e_context["reply"] = reply
            e_context.action = EventAction.BREAK_PASS

        except Exception as e:
            if retry_count < 1:
                logger.warning(f"[JinaSum] {str(e)}, retry {retry_count + 1}")
                self.on_handle_context(e_context, retry_count + 1)
                return

            logger.exception(f"[JinaSum] {str(e)}")
            reply = Reply(ReplyType.ERROR, "我暂时无法总结链接，请稍后再试")
            e_context["reply"] = reply
            e_context.action = EventAction.BREAK_PASS

    def get_help_text(self, verbose, **kwargs):
        return "使用jina reader和ChatGPT总结网页链接内容"

    def _load_config_template(self):
        logger.debug("No Suno plugin config.json, use plugins/jina_sum/config.json.template")
        try:
            plugin_config_path = os.path.join(self.path, "config.json.template")
            if os.path.exists(plugin_config_path):
                with open(plugin_config_path, encoding="utf-8") as f:
                    plugin_conf = json.load(f)
                    return plugin_conf
        except Exception as e:
            logger.exception(e)

    def _get_jina_url(self, target_url):
        return self.jina_reader_base + "/" + target_url

    def _get_openai_chat_url(self):
        return self.open_ai_api_base + "/chat/completions"

    def _get_openai_headers(self):
        return {"Authorization": f"Bearer {self.open_ai_api_key}", "Host": urlparse(self.open_ai_api_base).netloc}

    def _get_openai_payload(self, target_url_content):
        target_url_content = target_url_content[: self.max_words]  # 通过字符串长度简单进行截断
        sum_prompt = f"{self.prompt}\n\n'''{target_url_content}'''"
        messages = [{"role": "user", "content": sum_prompt}]
        payload = {"model": self.open_ai_model, "messages": messages}
        return payload

    def _check_url(self, target_url: str):
        stripped_url = target_url.strip()
        # 简单校验是否是url
        if not stripped_url.startswith("http://") and not stripped_url.startswith("https://"):
            return False

        # 检查白名单
        if len(self.white_url_list):
            if not any(stripped_url.startswith(white_url) for white_url in self.white_url_list):
                return False

        # 排除黑名单，黑名单优先级>白名单
        for black_url in self.black_url_list:
            if stripped_url.startswith(black_url):
                return False

        return True
