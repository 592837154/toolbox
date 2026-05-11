package gen

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GenRaw 生成占位文件：sparse 为 true 时仅 Truncate（可能稀疏）。
func GenRaw(path, sizeStr string, sparse bool) error {
	n, err := parseHumanSize(sizeStr)
	if err != nil {
		return err
	}
	if n <= 0 {
		return errors.New("size 必须大于 0")
	}
	if sparse {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := f.Truncate(n); err != nil {
			return err
		}
		fmt.Printf("已创建（稀疏/预分配） %s，逻辑大小 %d 字节\n", path, n)
		return nil
	}
	return writeZeroFilledFile(path, n)
}

func writeZeroFilledFile(path string, total int64) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	const chunk = 4 << 20 // 4 MiB
	buf := make([]byte, chunk)
	var written int64
	start := time.Now()
	for written < total {
		n := int64(len(buf))
		if remain := total - written; remain < n {
			n = remain
		}
		if _, err := f.Write(buf[:n]); err != nil {
			return err
		}
		written += n
		if written%(64<<20) == 0 || written == total {
			elapsed := time.Since(start).Seconds()
			mb := float64(written) / (1024 * 1024)
			fmt.Printf("\r已写入 %.1f MiB", mb)
			if elapsed > 0 {
				fmt.Printf("  (%.1f MiB/s)", mb/elapsed)
			}
		}
	}
	fmt.Println()
	fmt.Printf("已创建 %s，大小 %d 字节\n", path, total)
	return nil
}

func parseHumanSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, " ", "")
	if s == "" {
		return 0, errors.New("空的大小字符串")
	}
	i := len(s)
	for i > 0 && (s[i-1] >= 'a' && s[i-1] <= 'z') {
		i--
	}
	if i == 0 {
		return 0, errors.New("无效的大小格式")
	}
	numStr := s[:i]
	suffix := s[i:]
	v, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("解析数字: %w", err)
	}
	if v <= 0 {
		return 0, errors.New("数值必须为正")
	}
	mult := float64(1)
	switch suffix {
	case "", "b", "byte", "bytes":
		mult = 1
	case "k", "kb", "kib":
		mult = 1 << 10
	case "m", "mb", "mib":
		mult = 1 << 20
	case "g", "gb", "gib":
		mult = 1 << 30
	case "t", "tb", "tib":
		mult = 1 << 40
	default:
		return 0, fmt.Errorf("未知单位 %q（可用 kb/mb/gb/tb）", suffix)
	}
	out := int64(v * mult)
	if out < 1 {
		return 0, errors.New("结果过小")
	}
	return out, nil
}
