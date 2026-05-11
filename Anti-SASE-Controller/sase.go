package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// go build -o killSASE.exe sase.go
func main() {
	fmt.Println("======================================")
	fmt.Println("🛡️  SASE 终极猎杀器 (EXE 版)")
	fmt.Println("🚀 即使系统报错，本程序也绝不退出")
	fmt.Println("======================================")

	for {
		// 1. 定义命令
		cmd := exec.Command("taskkill", "/F", "/IM", "SASE*", "/T")

		// 2. 强行执行，完全不接收任何返回值和错误
		// 我们用 CombinedOutput 但直接把所有结果丢进黑洞
		output, _ := cmd.CombinedOutput()

		res := string(output)
		if strings.Contains(res, "成功") || strings.Contains(res, "SUCCESS") {
			fmt.Printf("[%s] 🎯 已捕获并处决 SASE 进程！\n", time.Now().Format("15:04:05"))
		}

		// 3. 每一秒扫描一次
		time.Sleep(1 * time.Second)
	}
}
