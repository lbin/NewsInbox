import datetime
import json
import os

import pytz

import dateutil.parser

def get_timestamp(time_string):
    # TODO: 优化时间戳获取方法
    if '-' in time_string:
        timestamp = dateutil.parser.isoparse(time_string)
    else:
        timestamp = datetime.datetime.strptime(time_string, "%a, %d %b %Y %H:%M:%S %Z")
        timestamp = timestamp.astimezone(pytz.timezone("Asia/Shanghai"))
    timestamp = int(timestamp.timestamp()) * 1000
    return timestamp


def save_to_json(json_data, sender=None, focus=True):
    # Generate file name based on current time and sender
    current_date = datetime.datetime.now().strftime("%Y-%m-%d")
    if focus:
        folder_path = f"./schedule/data/{current_date}/focus"
    else:
        folder_path = f"./schedule/data/{current_date}/non-focus"
    if not os.path.exists(folder_path):
        os.makedirs(folder_path)
    if sender is None:
        sender = "Unknown"
    file_name = f"{folder_path}/{datetime.datetime.now().strftime('%H-%M-%S')}_{sender}.json"
    # Save new_json_data to json file
    with open(file_name, "w") as json_file:
        json.dump(json_data, json_file, ensure_ascii=False, indent=4)