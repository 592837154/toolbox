# 🛠️ Toolbox | 实用工具库

A collection of practical tools and scripts designed to solve real-world efficiency and privacy challenges.
这里记录了我为了解决实际工作中的痛点而编写的各类小工具与脚本。

---

## 🚀 核心项目 / Featured Projects

### 1. [Anti-SASE Controller](./Anti-SASE-Controller)
**Language: Go (Golang)**

- **中文说明**: 解决公司 SASE 软件过度监控与自启动问题，通过 Go 脚本实现进程的一键终止与持续拦截，找回开发环境的控制权。
- **English**: A system-level controller built with Go to terminate and intercept invasive SASE processes, reclaiming control over the development environment.

### 2. [Weather Push Bot](./Weather-Push-Bot)
**Language: Python**

- **中文说明**: 0 成本自动化天气提醒方案。结合 HuggingFace 托管、高德地图 API 及 UptimeRobot，通过企业微信机器人实现每日准时推送。
- **English**: A zero-cost automated weather notification bot. Utilizing HuggingFace, Amap API, and UptimeRobot to deliver daily alerts via WeChat Work.

### 3. [File Generator](./File-Generator)
**Language: Go (Golang)**

- **中文说明**: 测试上传、磁盘与同步场景时快速生成指定体积的占位文件；并可借助 ffmpeg 生成黑屏 MP4 与静音 MP3 样本。默认写入 `output/` 目录并已配置 git 忽略，避免生成物进入版本库。
- **English**: A Go CLI for fixed-size placeholder files and minimal MP4/MP3 samples (via ffmpeg) for upload, disk, and pipeline testing. Outputs default to `output/` with gitignore rules so artifacts stay out of the repo.

---

## 🏗️ 正在建设中 / In Progress
- **[Dev-Log](https://github.com/592837154/dev-log)**: 记录这些工具背后的开发心得与技术避坑指南。
- **[Life-Notes](https://github.com/592837154/life-notes)**: 记录通往 37 岁目标的学习与生活点滴。

---

## 📈 未来计划 / Future Roadmap
- [ ] 更多基于 Go 的系统自动化脚本 (More Go-based system scripts)
- [ ] AI 辅助生成的效能小工具 (AI-powered productivity tools)
- [ ] 网络环境一键切换方案 (One-click network environment switcher)
