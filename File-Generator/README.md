
# File Generator (Go)

### 📖 项目背景 / Background
- **中文**: 测试上传限速、磁盘配额、备份与同步工具时，经常需要「指定体积」的占位文件；音视频流水线也需要可预期的最小 MP4 / MP3 样本。本工具用 Go 实现本地生成，并借助 ffmpeg 输出标准容器格式。
- **English**: When testing upload throttling, disk quotas, or backup/sync tools, fixed-size placeholder files are often needed; media pipelines also benefit from predictable minimal MP4/MP3 samples. This tool generates raw files in Go and produces standard containers via ffmpeg.

---

### ✨ 核心功能 / Features
- **file-raw.exe / 占位文件**:
    - **CN**: 修改 `cmd/file-raw/main.go` 顶部常量（文件名、大小、是否稀疏），编译后双击即生成。
    - **EN**: Edit constants in `cmd/file-raw/main.go`, build, double-click.
- **file-mp4.exe / 视频样本**:
    - **CN**: 修改 `cmd/file-mp4/main.go`，生成黑屏 H.264 MP4（依赖 ffmpeg）。
    - **EN**: Edit `cmd/file-mp4/main.go`; black-screen H.264 (requires ffmpeg).
- **file-mp3.exe / 音频样本**:
    - **CN**: 修改 `cmd/file-mp3/main.go`，生成静音 MP3（依赖 ffmpeg）。
    - **EN**: Edit `cmd/file-mp3/main.go`; silent MP3 (requires ffmpeg).

---

### 🛠️ 技术实现 / Technical Implementation
- **Language**: Go (Golang)
- **Mechanism**:
    - 三个独立 `main`：`cmd/file-raw`、`cmd/file-mp4`、`cmd/file-mp3`；共享逻辑在 `internal/gen`。
    - 输出目录均为 **exe 同目录下的 `output/`**，`output/` 已在 `.gitignore` 中忽略。
    - 任务结束后进程立即退出，双击运行时控制台窗口会自行关闭（若需看日志可在终端里手动运行 exe）。
- **Build / 编译**（在 `File-Generator` 目录下）:
    - `go build -o file-raw.exe ./cmd/file-raw`
    - `go build -o file-mp4.exe ./cmd/file-mp4`
    - `go build -o file-mp3.exe ./cmd/file-mp3`

---

### 📝 记录感悟 / Reflection
> “拆成三个 exe 后，改一种需求只动一个 main、只编一个程序，双击目标更单一。”
> "Three executables mean one main to edit and one binary to build per need—each double-click does one job."
