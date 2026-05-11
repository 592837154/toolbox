package main

import (
	"fmt"
	"os"

	"file-generator/internal/gen"
)

// ========== MP4：修改后保存，执行 go build -o file-mp4.exe ./cmd/file-mp4 ==========
// 需要系统 PATH 中有 ffmpeg。生成文件名：运行时刻的时分秒（HHmmss）+ .mp4

const (
	seconds = 10.0
	width   = 1280
	height  = 720
)

func main() {
	fmt.Println("file-mp4 — 黑屏 H.264 样本 → exe 同目录/output/")
	outFileName := gen.FileNameByHMS(".mp4")

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
