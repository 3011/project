import os
import time

import requests


tgbottoken = "botxxx"
chatid = "xxx"


def check():
    headers = {
        'User-Agent': "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36"
    }
    result = requests.get(
        "https://www.netflix.com/title/81215567", headers=headers)
    return result.status_code == 200


def send_msg(text):
    data = {"chat_id": chatid, "text": text}
    requests.post("https://api.telegram.org/" +
                  tgbottoken+"/sendMessage", json=data)


def get_ip():
    headers = {
        'User-Agent': "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36"
    }
    result = requests.get("http://ipinfo.io/ip", headers=headers)
    return result.text


def main():
    while 1:
        if check():
            time.sleep(30 * 60)
        else:
            os.system('./warp.sh rewg')
            send_msg("\nNetflix: " + str(check()) + "\nNew IP: " + get_ip())


if __name__ == '__main__':
    while 1:
        try:
            main()
        except:
            os.system('./warp.sh rewg')
