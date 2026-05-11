import requests
import json
import os
from flask import Flask
from datetime import datetime, timedelta
from apscheduler.schedulers.background import BackgroundScheduler

# ==================== 配置区 ====================
WEBHOOK_URL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=973582d2-a263-4e4f-ad8c-14b9c53a6fca"
AMAP_KEY = "a990a2ff2ada0e5fda4f04697f93a5a8"
CITY_CODE = "610100"  # 西安
CITY_NAME = "西安"
# ===============================================

app = Flask(__name__)

@app.route('/')
def home():
    # 增加当前服务器时间显示，方便通过网页确认服务是否活着
    now = datetime.utcnow() + timedelta(hours=8)
    return f"Service is Running! 🚀<br>Current Beijing Time: {now.strftime('%Y-%m-%d %H:%M:%S')}<br>Schedule: 07:30 & 18:35"

def get_full_weather():
    """获取高德全量天气数据"""
    try:
        url = f"https://restapi.amap.com/v3/weather/weatherInfo?city={CITY_CODE}&key={AMAP_KEY}&extensions=all"
        res = requests.get(url, timeout=10)
        data = res.json()
        if data["status"] == "1":
            return data["forecasts"][0]
        return None
    except Exception as e:
        print(f"天气数据获取失败: {e}")
        return None

def send_weather_msg():
    """任务1：发送详细天气消息 (每天 07:30)"""
    forecast = get_full_weather()
    if not forecast: return
    
    beijing_now = datetime.utcnow() + timedelta(hours=8)
    casts = forecast['casts']
    today = casts[0]
    
    # 计算今日温差
    diff = int(today['daytemp']) - int(today['nighttemp'])
    
    msg = [
        f"🌞 【西安天气全预报】",
        f"🕒 推送时间：{beijing_now.strftime('%H:%M:%S')}",
        f"📍 城市：{CITY_NAME} ({forecast['province']})",
        f"---------------------------",
        f"今日 ({today['date']})：",
        f"● 天气：白天 {today['dayweather']} | 夜间 {today['nightweather']}",
        f"● 气温：{today['nighttemp']}℃ ~ {today['daytemp']}℃",
        f"● 昼夜温差：{diff}℃",
        f"● 风力：{today['daywind']}风 {today['daypower']}级",
        f"---------------------------",
        f"🔮 未来三日趋势：",
        f"● 明日：{casts[1]['dayweather']} | {casts[1]['nighttemp']}~{casts[1]['daytemp']}℃",
        f"● 后日：{casts[2]['dayweather']} | {casts[2]['nighttemp']}~{casts[2]['daytemp']}℃",
        f"● 大后日：{casts[3]['dayweather']} | {casts[3]['nighttemp']}~{casts[3]['daytemp']}℃",
        f"---------------------------",
        f"💡 生活建议：西安气候干燥，注意补水；{'早晚温差大，记得披件外套' if diff > 10 else '气温较平稳，体感舒适'}"
    ]
    
    requests.post(WEBHOOK_URL, json={"msgtype": "text", "text": {"content": "\n".join(msg)}})
    print(f"[{datetime.now()}] 已发送：早上天气消息")

def send_remind_msg():
    """任务2：发送下班提醒消息 (每天 18:35)"""
    beijing_now = datetime.utcnow() + timedelta(hours=8)
    remind_text = [
        f"⏰ 兄弟，该下班了！",
        f"🕒 当前时间：{beijing_now.strftime('%H:%M:%S')}",
        f"━━━━━━━━━━━━",
        f"✅ 别忘了打卡（早晚都要记牢）",
        f"📊 别忘了甘特图（进度同步了吗？）",
        f"📝 别忘了写周会（小心发红包）",
        f"━━━━━━━━━━━━",
        f"🏃 辛苦了，早点回家休息！"
    ]
    requests.post(WEBHOOK_URL, json={"msgtype": "text", "text": {"content": "\n".join(remind_text)}})
    print(f"[{datetime.now()}] 已发送：傍晚下班提醒")

# 配置定时任务
scheduler = BackgroundScheduler(timezone='Asia/Shanghai')
scheduler.add_job(send_weather_msg, 'cron', hour=7, minute=30)
scheduler.add_job(send_remind_msg, 'cron', hour=18, minute=35)
scheduler.start()

if __name__ == "__main__":
    # 启动时立刻各发送一次，用于确认 Webhook 和 API 是否配置成功
    print("正在执行启动首次验证推送...")
    send_weather_msg()
    send_remind_msg()
    
    # 启动 Web 服务，监听 Hugging Face 的 7860 端口
    # debug=False 是为了防止定时任务被初始化两次
    app.run(host='0.0.0.0', port=7860, debug=False)