package main

import (
	"fmt"
	"os"
	"strings"

	"file-generator/internal/gen"
)

// ========== 占位文件：修改后保存，执行 go build -o file-raw.exe ./cmd/file-raw ==========
// 生成文件名：运行时刻的时分秒（HHmmss）+ .bin

const (
	outSize = "300mb" // 如 512、1kb、500mb、2gb、1.5g（1024 进制）
	sparse  = false   // true=仅预分配长度；false=逐块写 0
)

func main() {
	fmt.Println("file-raw — 生成指定大小的占位文件 → exe 同目录/output/")
	if strings.TrimSpace(outSize) == "" {
		fmt.Fprintln(os.Stderr, "请在 cmd/file-raw/main.go 中填写 outSize")
		os.Exit(1)
	}
	outFileName := gen.FileNameByHMS(".bin")

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
	if err := gen.GenRaw(path, outSize, sparse); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n完成。")
}
