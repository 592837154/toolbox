package gen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// OutputDirName 为 exe 同目录下用于存放生成物的文件夹名。
const OutputDirName = "output"

// FileNameByHMS 返回以当前本地时间「时分秒」命名的文件名，如 143052.mp4。
// ext 须带点，例如 ".bin"、".mp4"。
func FileNameByHMS(ext string) string {
	if ext != "" && ext[0] != '.' {
		ext = "." + ext
	}
	return time.Now().Format("150405") + ext
}

// ExeOutputBaseDir 返回「可执行文件所在目录/output」的绝对路径。
func ExeOutputBaseDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exe), OutputDirName), nil
}

// SafeOutputPath 将 name 解析为 baseDir 下的安全绝对路径，并创建父目录。
func SafeOutputPath(baseDir, name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("文件名为空")
	}
	cleanName := filepath.Clean(name)
	if cleanName == "." || cleanName == ".." || strings.HasPrefix(cleanName, ".."+string(filepath.Separator)) {
		return "", errors.New("文件名不能包含 ..")
	}
	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	final := filepath.Join(baseAbs, cleanName)
	finalAbs, err := filepath.Abs(final)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(baseAbs, finalAbs)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New("路径必须位于 output 目录内")
	}
	if err := os.MkdirAll(filepath.Dir(finalAbs), 0o755); err != nil {
		return "", err
	}
	return finalAbs, nil
}

// Fatal 打印错误后以退出码 1 结束（双击运行时控制台会随之关闭）。
func Fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
