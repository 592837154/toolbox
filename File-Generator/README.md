
# File Generator (Go)

### 📖 项目背景 / Background
- **中文**: 测试上传限速、磁盘配额、备份与同步工具时，经常需要「指定体积」的占位文件；音视频流水线也需要可预期的最小 MP4 / MP3 样本。本工具用 Go 实现本地生成，并借助 ffmpeg 输出标准容器格式。
- **English**: When testing upload throttling, disk quotas, or backup/sync tools, fixed-size placeholder files are often needed; media pipelines also benefit from predictable minimal MP4/MP3 samples. This tool generates raw files in Go and produces standard containers via ffmpeg.

---

### ✨ 核心功能 / Features
- **Raw / 占位文件**:
    - **CN**: 按人类可读大小（如 `1gb`、`500mb`、`2gb`）生成文件；可选稀疏预分配（`-sparse`）或全量写 0（占满真实磁盘）。
    - **EN**: Create files from human-readable sizes (e.g. `1gb`, `500mb`); optional sparse preallocation (`-sparse`) or fully zero-filled output.
- **MP4 / 视频样本**:
    - **CN**: 调用 ffmpeg 生成黑屏 H.264 MP4，可调时长与分辨率。
    - **EN**: Invokes ffmpeg to produce black-screen H.264 MP4 with configurable duration and resolution.
- **MP3 / 音频样本**:
    - **CN**: 调用 ffmpeg 生成静音立体声 MP3（固定码率示例）。
    - **EN**: Invokes ffmpeg to produce silent stereo MP3 (CBR example).

---

### 🛠️ 技术实现 / Technical Implementation
- **Language**: Go (Golang)
- **Mechanism**:
    - `raw`：`os.Create`、`Truncate` 或分块写入零字节。
    - `mp4` / `mp3`：`os/exec` 调用 `ffmpeg`（需已安装并在 PATH 中）。
    - 相对路径的 `-out` 写入 `-dir` 目录（默认当前工作目录下的 `output/`，可用 `-dir` 修改）；`-out` 为绝对路径时不受 `-dir` 影响。`output/` 已在 `.gitignore` 中忽略，避免生成物进入版本库。
- **Build / 编译**:
    - `go build -o file-generator.exe .`（Windows 示例，可自定义输出名）
- **Usage / 用法示例**:
    - `file-generator.exe raw -out big.bin -size 2gb` → 实际为 `output/big.bin`
    - `file-generator.exe raw -out sparse.bin -size 1gb -sparse`
    - `file-generator.exe mp4 -out clip.mp4 -duration 30s -width 1280x720`
    - `file-generator.exe mp3 -out silence.mp3 -duration 5m`
    - `file-generator.exe raw -dir D:\tmp -out x.bin -size 10mb` → 写入 `D:\tmp\x.bin`

---

### 📝 记录感悟 / Reflection
> “小工具的价值在于把重复、易错的手工步骤固化成一条命令。占位文件与样本媒体不需要华丽，但要可靠、可脚本化。”
> "The value of a small tool is turning repetitive, error-prone steps into one command. Placeholders and sample media need not be fancy—only dependable and scriptable."
