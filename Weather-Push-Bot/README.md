# 🌦️ 0-Cost Weather Notification Bot (Python & HuggingFace)

### 📖 项目背景 / Background
- **中文**: 某天上班路上突降大雨，让我意识到需要一个出门前的提醒。为了实现零成本、高可用的方案，我开发了这个企业微信天气机器人。
- **English**: A sudden downpour on my way to work made me realize the need for a pre-commute reminder. I developed this WeChat Work bot to achieve a zero-cost, high-availability weather notification solution.

---

### ✨ 技术方案 / Technical Stack
- **Python**: 核心逻辑，通过高德地图 API 获取实时天气数据。
- **Hugging Face Spaces**: 作为免费的托管服务器运行 Python 脚本。
- **UptimeRobot**: 解决 Hugging Face 免费空间的休眠问题，通过定时心跳访问保持服务常驻。
- **WeChat Work Bot**: 利用企业微信群机器人 Webhook 实现消息推送。

---

### 🚀 核心逻辑 / Key Features
- **Smart Alert**: 每天早晨 7:30 准时推送，提醒是否需要带伞。
- **Zero Cost**: 结合多个免费云服务，实现 0 成本、全自动运行。
- **Anti-Sleep**: 巧妙利用 UptimeRobot 绕过服务器休眠限制。

---

### 📝 记录感悟 / Reflection
> “这个工具体现了我的架构思维：技术不一定要复杂，但一定要精准解决问题。通过巧妙组合免费资源，我拥有了一个私人定制的天气管家。”
> "This project reflects my architectural thinking: technology doesn't have to be complex to be effective. By combining free cloud resources, I've built a customized personal weather butler."
