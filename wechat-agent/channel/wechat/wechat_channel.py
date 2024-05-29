# encoding:utf-8

"""
wechat channel
"""
import base64
import hashlib
import io
import json
import os
import threading
import time

import requests

from bridge.context import *
from bridge.reply import *
from channel.chat_channel import ChatChannel
from channel import chat_channel
from channel.wechat.wechat_message import *
from common.expired_dict import ExpiredDict
from common.log import logger
from common.singleton import singleton
from common.time_check import time_checker
from newsinbox.common.config import conf, get_appdata_dir
from lib import itchat
from lib.itchat.content import *


@itchat.msg_register([TEXT, VOICE, PICTURE, NOTE, ATTACHMENT, SHARING])
def handler_single_msg(msg):
    try:
        cmsg = WechatMessage(msg, False)
    except NotImplementedError as e:
        logger.debug("[WX]single message {} skipped: {}".format(msg["MsgId"], e))
        return None
    WechatChannel().handle_single(cmsg)
    return None


@itchat.msg_register([TEXT, VOICE, PICTURE, NOTE, ATTACHMENT, SHARING], isGroupChat=True)
def handler_group_msg(msg):
    try:
        cmsg = WechatMessage(msg, True)
    except NotImplementedError as e:
        logger.debug("[WX]group message {} skipped: {}".format(msg["MsgId"], e))
        return None
    WechatChannel().handle_group(cmsg)
    return None

# 自动接受加好友
@itchat.msg_register(FRIENDS)
def deal_with_friend(msg):
    try:
        cmsg = WechatMessage(msg, False)
    except NotImplementedError as e:
        logger.debug("[WX]friend request {} skipped: {}".format(msg["MsgId"], e))
        return None
    WechatChannel().handle_friend_request(cmsg)
    return None

def _check(func):
    def wrapper(self, cmsg: ChatMessage):
        msgId = cmsg.msg_id
        if msgId in self.receivedMsgs:
            logger.info("Wechat message {} already received, ignore".format(msgId))
            return
        self.receivedMsgs[msgId] = True
        create_time = cmsg.create_time  # 消息时间戳
        if conf().get("hot_reload") == True and int(create_time) < int(time.time()) - 60:  # 跳过1分钟前的历史消息
            logger.debug("[WX]history message {} skipped".format(msgId))
            return
        if cmsg.my_msg and not cmsg.is_group:
            logger.debug("[WX]my message {} skipped".format(msgId))
            return
        return func(self, cmsg)

    return wrapper

def send_wechat(webhook, data):
    header = {
                "Content-Type": "application/json",
                "Charset": "UTF-8"
                }
    info = requests.post(url=webhook, json=data, headers=header)
    print(info)
    
    
# 可用的二维码生成接口
# https://api.qrserver.com/v1/create-qr-code/?size=400×400&data=https://www.abc.com
# https://api.isoyu.com/qr/?m=1&e=L&p=20&url=https://www.abc.com
def qrCallback(uuid, status, qrcode):
    # logger.debug("qrCallback: {} {}".format(uuid,status))
    if status == "0":
        try:
            from PIL import Image

            img = Image.open(io.BytesIO(qrcode))
            _thread = threading.Thread(target=img.show, args=("QRCode",))
            _thread.setDaemon(True)
            _thread.start()
        except Exception as e:
            pass

        import qrcode

        url = f"https://login.weixin.qq.com/l/{uuid}"

        qr_api1 = "https://api.isoyu.com/qr/?m=1&e=L&p=20&url={}".format(url)
        qr_api2 = "https://api.qrserver.com/v1/create-qr-code/?size=400×400&data={}".format(url)
        qr_api3 = "https://api.pwmqr.com/qrcode/create/?url={}".format(url)
        qr_api4 = "https://my.tv.sohu.com/user/a/wvideo/getQRCode.do?text={}".format(url)
        print("You can also scan QRCode in any website below:")
        print(qr_api3)
        print(qr_api4)
        print(qr_api2)
        print(qr_api1)
        _send_qr_code([qr_api3, qr_api4, qr_api2, qr_api1])
        qr = qrcode.QRCode(border=1)
        qr.add_data(url)
        qr.make(fit=True)
        qr.print_ascii(invert=True)
        
        im = qr.make_image()
        im.save("qrcode.png")
        
        with open(r'qrcode.png', 'rb') as file:
            encodestr = base64.b64encode(file.read())
            image_data = str(encodestr, 'utf-8')

        with open(r'qrcode.png', 'rb') as file:
            md = hashlib.md5()
            md.update(file.read())
            image_md5 = md.hexdigest()
        
        
        data = {
            "msgtype": "image",
            "image": {
                "base64": image_data,
                "md5": image_md5
            }
        }
        send_wechat(conf().get("we_work_webhook"), data)
        

@singleton
class WechatChannel(ChatChannel):
    NOT_SUPPORT_REPLYTYPE = []

    def __init__(self):
        super().__init__()
        self.receivedMsgs = ExpiredDict(60 * 60)
        self.auto_login_times = 0

    def startup(self):
        try:
            itchat.instance.receivingRetryCount = 600  # 修改断线超时时间
            # login by scan QRCode
            hotReload = conf().get("hot_reload", False)
            status_path = os.path.join(get_appdata_dir(), "itchat.pkl")
            itchat.auto_login(
                enableCmdQR=2,
                hotReload=hotReload,
                statusStorageDir=status_path,
                qrCallback=qrCallback,
                exitCallback=self.exitCallback,
                loginCallback=self.loginCallback
            )
            self.user_id = itchat.instance.storageClass.userName
            self.name = itchat.instance.storageClass.nickName
            logger.info("Wechat login success, user_id: {}, nickname: {}".format(self.user_id, self.name))
            # start message listener
            itchat.run()
        except Exception as e:
            logger.error(e)

    def exitCallback(self):
        try:
            _send_logout()
            time.sleep(2)
            self.auto_login_times += 1
            if self.auto_login_times < 100:
                chat_channel.handler_pool._shutdown = False
                self.startup()
                
        except Exception as e:
            pass

    def loginCallback(self):
        logger.debug("Login success")
        _send_login_success()

    # handle_* 系列函数处理收到的消息后构造Context，然后传入produce函数中处理Context和发送回复
    # Context包含了消息的所有信息，包括以下属性
    #   type 消息类型, 包括TEXT、VOICE、IMAGE_CREATE
    #   content 消息内容，如果是TEXT类型，content就是文本内容，如果是VOICE类型，content就是语音文件名，如果是IMAGE_CREATE类型，content就是图片生成命令
    #   kwargs 附加参数字典，包含以下的key：
    #        session_id: 会话id
    #        isgroup: 是否是群聊
    #        receiver: 需要回复的对象
    #        msg: ChatMessage消息对象
    #        origin_ctype: 原始消息类型，语音转文字后，私聊时如果匹配前缀失败，会根据初始消息是否是语音来放宽触发规则
    #        desire_rtype: 希望回复类型，默认是文本回复，设置为ReplyType.VOICE是语音回复
    @time_checker
    @_check
    def handle_single(self, cmsg: ChatMessage):
        # filter system message
        if cmsg.other_user_id in ["weixin"]:
            return
        if cmsg.ctype == ContextType.VOICE:
            if conf().get("speech_recognition") != True:
                return
            logger.debug("[WX]receive voice msg: {}".format(cmsg.content))
        elif cmsg.ctype == ContextType.IMAGE:
            logger.debug("[WX]receive image msg: {}".format(cmsg.content))
        elif cmsg.ctype == ContextType.PATPAT:
            logger.debug("[WX]receive patpat msg: {}".format(cmsg.content))
        elif cmsg.ctype == ContextType.TEXT:
            logger.debug("[WX]receive text msg: {}, cmsg={}".format(json.dumps(cmsg._rawmsg, ensure_ascii=False), cmsg))
        else:
            logger.debug("[WX]receive msg: {}, cmsg={}".format(cmsg.content, cmsg))
        context = self._compose_context(cmsg.ctype, cmsg.content, isgroup=False, msg=cmsg)
        if context:
            self.produce(context)

    @time_checker
    @_check
    def handle_group(self, cmsg: ChatMessage):
        logger.info(f"{cmsg.ctype}")
        if cmsg.ctype == ContextType.VOICE:
            if conf().get("group_speech_recognition") != True:
                return
            logger.debug("[WX]receive voice for group msg: {}".format(cmsg.content))
        elif cmsg.ctype == ContextType.IMAGE:
            logger.debug("[WX]receive image for group msg: {}".format(cmsg.content))
        elif cmsg.ctype in [ContextType.JOIN_GROUP, ContextType.PATPAT, ContextType.ACCEPT_FRIEND, ContextType.EXIT_GROUP]:
            logger.debug("[WX]receive note msg: {}".format(cmsg.content))
        elif cmsg.ctype == ContextType.TEXT:
            # logger.debug("[WX]receive group msg: {}, cmsg={}".format(json.dumps(cmsg._rawmsg, ensure_ascii=False), cmsg))
            pass
        elif cmsg.ctype == ContextType.FILE:
            logger.debug(f"[WX]receive attachment msg, file_name={cmsg.content}")
        else:
            logger.debug("[WX]receive group msg: {}".format(cmsg.content))
        context = self._compose_context(cmsg.ctype, cmsg.content, isgroup=True, msg=cmsg)
        if context:
            self.produce(context)
            
    def send_to_feishu(self, reply, resceiver, group_name=None, nick_name=None):

        feishu_get_token_url = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal"
        feishu_add_record_url = "https://open.feishu.cn/open-apis/bitable/v1/apps/DM8Ib7DNeah45XsA8kOcxvYenHd/tables/tblELLHLYuKXsVf9/records"

        content = reply.content
        if content is None:
            return None
        
        url = reply.url
        data_index = content.find("json")
        user = content[1:data_index-3]
        user = user.replace("\n", "")
        
        user_group = conf().get("user_group")
        for tmp_user in user_group:
            if nick_name == tmp_user["user_name"]:
                user = nick_name
                feishu_add_record_url = tmp_user["feishu_add_record_url"]
                break
        
        user_group_group = conf().get("user_group_group")
        for tmp_user in user_group_group:
            if group_name == tmp_user["user_name"]:
                # user = nick_name
                feishu_add_record_url = tmp_user["feishu_add_record_url"]
                break
        

        json_data = content[data_index+3:-3]
        json_data = json_data.replace("\n", "")
        json_data = json_data.replace(" ", "")
        json_data = json_data[1:]

        json_data = json.loads(json_data)

        headers = {
            "Content-Type": "application/json; charset=utf-8"
            }
        token_json_data = {
            "app_id": "cli_a6bafbc7acbdd013",
            "app_secret": "8U3AAH3AOPtwVOFdWusswctp5EOcqFM5"
            }
        feishu_token_response = requests.post(feishu_get_token_url, headers=headers, json=token_json_data)

        feishu_headers = {
            "Content-Type": "application/json; charset=utf-8",
            'Authorization': f"Bearer {feishu_token_response.json()['app_access_token']}"
        }

        key_points_str = ""
        for key in json_data['key_points']:
            key_points_str += json_data['key_points'][key] + "\n"
        
        tags = json_data["tags"]
        new_tags = []
        
        for tag in tags:
            tag = tag.replace("#", "")
            new_tags.append(tag)
        tags = new_tags

        new_json_data = {'分类': ['Technology'], '标签': tags, '项目来源': '个人分享', '项目名称': json_data['title'], '来源': url, '总结': json_data['summary'], '关键要点': key_points_str, '添加人': user}

        new_data={'fields': new_json_data}
        status = requests.post(feishu_add_record_url, headers=feishu_headers, json=new_data)
        return status

    @time_checker
    @_check
    def handle_friend_request(self, cmsg: ChatMessage):
        if cmsg.ctype == ContextType.ACCEPT_FRIEND:
            logger.debug("[WX]receive friend request: {}".format(cmsg.content["NickName"]))
        else:
            logger.debug("[WX]receive friend request: {}, cmsg={}".format(cmsg.content["NickName"], cmsg))
        context = self._compose_context(cmsg.ctype, cmsg.content, msg=cmsg)
        if context:
            self.produce(context)


    # 统一的发送函数，每个Channel自行实现，根据reply的type字段发送不同类型的消息
    def send(self, reply: Reply, context: Context):
        receiver = context.get("receiver")
        if reply.type == ReplyType.TEXT:
            itchat.send(reply.content, toUserName=receiver)
            if reply.url is not None:
                nick_name = context["msg"].from_user_nickname
                group_name = context["msg"].other_user_nickname
                is_group = context["isgroup"]
                logger.info("[WX] sendMsg={}, group_name={}, nick_name={}, is_group={}".format(reply.content, group_name, nick_name, context["isgroup"]))
                if is_group is False:
                    self.send_to_feishu(reply, receiver, None, nick_name)
                else:
                    self.send_to_feishu(reply, receiver, group_name, None)
            logger.info("[WX] sendMsg={}, receiver={}".format(reply, receiver))
        elif reply.type == ReplyType.ERROR or reply.type == ReplyType.INFO:
            itchat.send(reply.content, toUserName=receiver)
            logger.info("[WX] sendMsg={}, receiver={}".format(reply, receiver))
        elif reply.type == ReplyType.VOICE:
            itchat.send_file(reply.content, toUserName=receiver)
            logger.info("[WX] sendFile={}, receiver={}".format(reply.content, receiver))
        elif reply.type == ReplyType.IMAGE_URL:  # 从网络下载图片
            img_url = reply.content
            logger.debug(f"[WX] start download image, img_url={img_url}")
            pic_res = requests.get(img_url, stream=True)
            image_storage = io.BytesIO()
            size = 0
            for block in pic_res.iter_content(1024):
                size += len(block)
                image_storage.write(block)
            logger.info(f"[WX] download image success, size={size}, img_url={img_url}")
            image_storage.seek(0)
            itchat.send_image(image_storage, toUserName=receiver)
            logger.info("[WX] sendImage url={}, receiver={}".format(img_url, receiver))
        elif reply.type == ReplyType.IMAGE:  # 从文件读取图片
            image_storage = reply.content
            image_storage.seek(0)
            itchat.send_image(image_storage, toUserName=receiver)
            logger.info("[WX] sendImage, receiver={}".format(receiver))
        elif reply.type == ReplyType.FILE:  # 新增文件回复类型
            file_storage = reply.content
            itchat.send_file(file_storage, toUserName=receiver)
            logger.info("[WX] sendFile, receiver={}".format(receiver))
        elif reply.type == ReplyType.VIDEO:  # 新增视频回复类型
            video_storage = reply.content
            itchat.send_video(video_storage, toUserName=receiver)
            logger.info("[WX] sendFile, receiver={}".format(receiver))
        elif reply.type == ReplyType.VIDEO_URL:  # 新增视频URL回复类型
            video_url = reply.content
            logger.debug(f"[WX] start download video, video_url={video_url}")
            video_res = requests.get(video_url, stream=True)
            video_storage = io.BytesIO()
            size = 0
            for block in video_res.iter_content(1024):
                size += len(block)
                video_storage.write(block)
            logger.info(f"[WX] download video success, size={size}, video_url={video_url}")
            video_storage.seek(0)
            itchat.send_video(video_storage, toUserName=receiver)
            logger.info("[WX] sendVideo url={}, receiver={}".format(video_url, receiver))

        elif reply.type == ReplyType.ACCEPT_FRIEND:  # 新增接受好友申请回复类型
            # 假设 reply.content 包含了新好友的用户名
            is_accept = reply.content
            if is_accept:
                try:
                    # 自动接受好友申请
                    debug_msg = itchat.accept_friend(userName=context.content["UserName"], v4=context.content["Ticket"])
                    logger.debug("[WX] accept_friend return: {}".format(debug_msg))
                    logger.info("[WX] Accepted new friend, UserName={}, NickName={}".format(context.content["UserName"], context.content["NickName"]))
                except Exception as e:
                    logger.error("[WX] Failed to add friend. Error: {}".format(e))
            else:
                logger.info("[WX] Ignored new friend, username={}".format(context.content["NickName"]))
        elif reply.type == ReplyType.INVITE_ROOM:  # 新增邀请好友进群回复类型
            # 假设 reply.content 包含了群聊的名字

            def get_group_id(group_name):
                """
                根据群聊名称获取群聊ID。
                :param group_name: 群聊的名称。
                :return: 群聊的ID (UserName)。
                """
                group_list = itchat.search_chatrooms(name=group_name)
                if group_list:
                    return group_list[0]["UserName"]
                else:
                    return None

            try:
                chatroomUserName = reply.content
                group_id = get_group_id(chatroomUserName)
                logger.debug("[WX] find group_id={}, where chatroom={}".format(group_id, chatroomUserName))
                if group_id is None:
                    raise ValueError("The specified group chat was not found: {}".format(chatroomUserName))
                # 调用 itchat 的 add_member_into_chatroom 方法来添加成员
                debug_msg = itchat.add_member_into_chatroom(group_id, receiver)
                logger.debug("[WX] add_member_into_chatroom return: {}".format(debug_msg))
                logger.info("[WX] invite members={}, to chatroom={}".format(receiver, chatroomUserName))
            except ValueError as ve:
                # 记录查找群聊失败的错误信息
                logger.error("[WX] {}".format(ve))
            except Exception as e:
                # 记录添加成员失败的错误信息
                logger.error("[WX] Failed to invite members to chatroom. Error: {}".format(e))

def _send_login_success():
    try:
        from common.linkai_client import chat_client
        if chat_client.client_id:
            chat_client.send_login_success()
    except Exception as e:
        pass

def _send_logout():
    try:
        from common.linkai_client import chat_client
        if chat_client.client_id:
            chat_client.send_logout()
    except Exception as e:
        pass

def _send_qr_code(qrcode_list: list):
    try:
        from common.linkai_client import chat_client

        if chat_client.client_id:
            chat_client.send_qrcode(qrcode_list)
    except Exception as e:
        pass
