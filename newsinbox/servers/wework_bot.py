import requests


def send_to_wework_bot(webhook, data):
    header = {
                "Content-Type": "application/json",
                "Charset": "UTF-8"
                }
    info = requests.post(url=webhook, json=data, headers=header)
    print(info)