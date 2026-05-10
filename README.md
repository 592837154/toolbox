# 🛠️ Toolbox | 实用工具库

A collection of practical tools and scripts designed to solve real-world efficiency and privacy challenges.
这里记录了我为了解决实际工作中的痛点而编写的各类小工具与脚本。

---

## 🚀 核心项目 / Featured Project

### [Anti-SASE Controller](./Anti-SASE-Controller)
**Language: Go (Golang)**

#### **中文说明**
- **痛点**：公司强制安装的 SASE 软件存在过度监听上网记录、限制代理软件使用且关闭后频繁自启动的问题，严重干扰开发流程。
- **解决方案**：使用 Go 开发的系统级控制工具。
    - **一键终止**：快速杀掉所有 SASE 相关进程。
    - **持续拦截**：后台监听系统进程，一旦检测到 SASE 尝试自动重启，立即将其拦截，确保环境纯净。
- **价值**：实现了办公合规与个人隐私/开发自由的完美平衡。

#### **English Description**
- **Pain Point**: Mandatory corporate SASE software monitors web history, blocks proxy tools, and repeatedly auto-restarts, disrupting the development workflow.
- **Solution**: A system-level controller built with Go.
    - **One-Click Kill**: Instantly terminates all SASE-related processes.
    - **Active Interception**: Background monitoring that intercepts and kills SASE immediately upon any unauthorized auto-restart attempts.
- **Value**: Achieves a perfect balance between corporate compliance and personal privacy/developer freedom.

---

## 📈 未来计划 / Future Roadmap
- [ ] 更多基于 Go 的系统自动化脚本
- [ ] AI 辅助生成的效能小工具
- [ ] 网络环境一键切换方案
