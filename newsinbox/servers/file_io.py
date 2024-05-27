import codecs
import datetime
import json
import os

import dateutil.parser
import markdown
import pdfkit
import pytz


def get_timestamp(time_string):
    # TODO: 优化时间戳获取方法
    if "-" in time_string:
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


def save_to_pdf(markdown_str):

    html_content = markdown.markdown(markdown_str)
    # html_file_name = f"templates/{datetime.datetime.now().strftime('%Y-%m-%d')}.html"
    html_file_name = "templates/today.html"
    
    with codecs.open(html_file_name, "w", encoding="utf-8") as f:
        f.write('<meta content="text/html; charset=utf-8" http-equiv="Content-Type"/>')
        f.write(html_content)
        
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d")
    pdf_file_name = f"schedule/data/{timestamp}.pdf"
 
    pdfkit.from_file(html_file_name, pdf_file_name)
