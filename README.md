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

- **中文说明**: 测试上传、磁盘与同步场景时快速生成指定体积的占位文件；并可借助 ffmpeg 生成黑屏 MP4 与静音 MP3 样本。三个独立程序（`file-raw` / `file-mp4` / `file-mp3`）各自改常量、`go build` 后双击即可；输出文件按运行时刻**时分秒**命名（如 `143052.bin`），位于 exe 同目录的 `output/`（已 git 忽略）；看日志请在终端中运行 exe。
- **English**: Three tools: edit constants, build, double-click; output files use local **time HHmmss** names (e.g. `143052.mp4`) under `output/` next to each binary (gitignored); run from a terminal for logs.

### 4. [Mall Admin CRUD](./Mall-Admin-CRUD)
**Stack: React / TypeScript / Ant Design Pro / Go / Gin / GORM / TiDB Cloud / Docker**

- **中文说明**: 一个完整的商品管理 CRUD 后台案例，覆盖 `ProTable` 表格查询、新建编辑弹窗、删除确认、Gin + GORM 后端接口、TiDB Cloud 数据持久化、本地 Docker 与 Hugging Face Spaces 公网部署。重点记录从本地可运行到线上可访问之间遇到的真实工程问题。
- **English**: A full-stack product management CRUD case built with Ant Design Pro and Go. It covers table querying, modal forms, REST APIs, TiDB Cloud persistence, Docker local runtime, and Hugging Face Spaces deployment.

---

## 🏗️ 正在建设中 / In Progress
- **[Dev-Log](https://github.com/592837154/dev-log)**: 记录这些工具背后的开发心得与技术避坑指南。
- **[Life-Notes](https://github.com/592837154/life-notes)**: 记录通往 37 岁目标的学习与生活点滴。

---

## 📈 未来计划 / Future Roadmap
- [ ] 更多基于 Go 的系统自动化脚本 (More Go-based system scripts)
- [ ] AI 辅助生成的效能小工具 (AI-powered productivity tools)
- [ ] 网络环境一键切换方案 (One-click network environment switcher)
