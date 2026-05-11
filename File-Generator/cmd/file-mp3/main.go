package main

import (
	"fmt"
	"os"

	"file-generator/internal/gen"
)

// ========== MP3：修改后保存，执行 go build -o file-mp3.exe ./cmd/file-mp3 ==========
// 需要系统 PATH 中有 ffmpeg。生成文件名：运行时刻的时分秒（HHmmss）+ .mp3

const (
	seconds = 60.0
)

func main() {
	fmt.Println("file-mp3 — 静音 MP3 样本 → exe 同目录/output/")
	outFileName := gen.FileNameByHMS(".mp3")

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
	if err := gen.GenMP3(path, seconds); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n完成。")
}
