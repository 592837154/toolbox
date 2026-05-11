package main

import (
	"fmt"
	"os"
	"strings"

	"file-generator/internal/gen"
)

// ========== MP4：修改后保存，执行 go build -o file-mp4.exe ./cmd/file-mp4 ==========
// 需要系统 PATH 中有 ffmpeg。

const (
	outFileName = "sample.mp4"
	seconds     = 10.0
	width       = 1280
	height      = 720
)

func main() {
	fmt.Println("file-mp4 — 黑屏 H.264 样本 → exe 同目录/output/")
	if strings.TrimSpace(outFileName) == "" {
		fmt.Fprintln(os.Stderr, "请在 cmd/file-mp4/main.go 中填写 outFileName")
		os.Exit(1)
	}

	baseDir, err := gen.ExeOutputBaseDir()
	if err != nil {
		gen.Fatal(err)
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		gen.Fatal(err)
	}
	fmt.Printf("输出目录: %s\n\n", baseDir)

	path, err := gen.SafeOutputPath(baseDir, outFileName)
	if err != nil {
		gen.Fatal(err)
	}
	if err := gen.GenMP4(path, seconds, width, height); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n完成。")
}
