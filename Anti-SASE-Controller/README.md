
# Anti-SASE Controller (Go)

### 📖 项目背景 / Background
- **中文**: 公司运维强制要求安装 SASE 软件以访问内网。但该软件存在过度监控隐私、占用系统资源、干扰代理软件且关闭后强制自启动的痛点。
- **English**: Corporate IT requires SASE software for intranet access. However, it raises concerns due to invasive monitoring, high resource usage, interference with proxy tools, and forced auto-restarts after being closed.

---

### ✨ 核心功能 / Features
- **Power Kill**: 
    - **CN**: 一键强制杀掉所有相关的 SASE 进程。
    - **EN**: One-click termination of all related SASE processes.
- **Persistent Guard**: 
    - **CN**: 持续监听系统进程。一旦 SASE 尝试静默自启动，程序将立即检测并再次拦截。
    - **EN**: Continuous system process monitoring. Immediately detects and intercepts any unauthorized SASE auto-restart attempts.
- **User Control**: 
    - **CN**: 除非手动关闭此 Go 工具，否则 SASE 将无法在后台运行，确保用户对个人电脑的完全控制权。
    - **EN**: SASE remains inactive as long as this Go tool is running, granting users full control over their personal environment.

---

### 🛠️ 技术实现 / Technical Implementation
- **Language**: Go (Golang)
- **Mechanism**: 
    - 使用 `os/exec` 调用系统指令查询进程。
    - 循环监听（Polling）机制实现准实时拦截。
    - 编译建议 (Build Suggestion): `go build -ldflags "-H windowsgui"` (Hide CMD window on Windows).

---

### 📝 记录感悟 / Reflection
> “这是我转型全栈后的第一个实战小工具。通过技术手段解决工作中的不便，不仅提高了效率，更让我感受到了编码的自由。”
> "This is my first practical tool after pivoting to Full-stack. Solving workplace pain points through technology not only boosts efficiency but also reinforces the sense of freedom in coding."
