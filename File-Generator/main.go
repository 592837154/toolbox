package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	switch cmd {
	case "raw":
		if err := runRaw(args); err != nil {
			fmt.Fprintf(os.Stderr, "raw: %v\n", err)
			os.Exit(1)
		}
	case "mp4":
		if err := runMP4(args); err != nil {
			fmt.Fprintf(os.Stderr, "mp4: %v\n", err)
			os.Exit(1)
		}
	case "mp3":
		if err := runMP3(args); err != nil {
			fmt.Fprintf(os.Stderr, "mp3: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "未知子命令 %q\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `filegen — 生成指定大小占位文件，或通过 ffmpeg 生成 MP4 / MP3

用法:
  filegen raw  -out <相对或绝对路径> -size <大小> [-dir <目录>] [-sparse]
  filegen mp4  -out <相对或绝对路径> [-dir <目录>] [-duration <时长>] [-width WxH]
  filegen mp3  -out <相对或绝对路径> [-dir <目录>] [-duration <时长>]

输出目录:
  相对路径的 -out 会写入 -dir 之下（默认 -dir 为 output，即当前工作目录下的 output/）。
  绝对路径的 -out 不受 -dir 限制。仓库中已将 output/ 加入 .gitignore。

raw 说明:
  -size  支持 512、1kb、500mb、2gb、1.5g 等形式（默认按 1024 进制）
  -sparse  仅预分配长度（可能为稀疏文件，速度快）；省略则逐块写入 0（占满真实空间）

mp4 / mp3 说明:
  需要已在 PATH 中安装 ffmpeg。mp4 为黑屏 H.264；mp3 为静音 CBR 128k。

示例:
  filegen raw  -out big.bin -size 2gb
  filegen raw  -out full.bin -size 1gb        # 非稀疏，写入全 0 → output/full.bin
  filegen mp4  -out clip.mp4 -duration 30s
  filegen mp3  -out silence.mp3 -duration 5m
`)
}

func runRaw(args []string) error {
	fs := flag.NewFlagSet("raw", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	out := fs.String("out", "", "输出文件路径（相对路径写入 -dir）")
	outDir := fs.String("dir", "output", "相对 -out 的根目录，默认 output")
	sizeStr := fs.String("size", "", "大小，如 1gb、500mb")
	sparse := fs.Bool("sparse", false, "稀疏分配（Truncate），不逐字节写盘")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *out == "" {
		return errors.New("需要 -out")
	}
	if *sizeStr == "" {
		return errors.New("需要 -size")
	}
	n, err := parseHumanSize(*sizeStr)
	if err != nil {
		return err
	}
	if n <= 0 {
		return errors.New("size 必须大于 0")
	}
	path, err := resolveOutputPath(*outDir, *out)
	if err != nil {
		return err
	}
	if *sparse {
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

func runMP4(args []string) error {
	fs := flag.NewFlagSet("mp4", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	out := fs.String("out", "", "输出 .mp4 路径（相对路径写入 -dir）")
	outDir := fs.String("dir", "output", "相对 -out 的根目录，默认 output")
	dur := fs.String("duration", "10s", "时长，如 30s、2m、1h")
	res := fs.String("width", "1280x720", "分辨率 WxH")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *out == "" {
		return errors.New("需要 -out")
	}
	sec, err := parseDurationToSeconds(*dur)
	if err != nil {
		return err
	}
	if err := requireFFmpeg(); err != nil {
		return err
	}
	w, h, err := parseWxH(*res)
	if err != nil {
		return err
	}
	path, err := resolveOutputPath(*outDir, *out)
	if err != nil {
		return err
	}
	// 黑屏 + yuv420p，兼容常见播放器
	filter := fmt.Sprintf("color=c=black:s=%dx%d:r=30", w, h)
	cmd := exec.Command("ffmpeg",
		"-y",
		"-f", "lavfi", "-i", filter,
		"-t", fmt.Sprintf("%f", sec),
		"-pix_fmt", "yuv420p",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-movflags", "+faststart",
		path,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg: %w", err)
	}
	fmt.Printf("已生成 MP4: %s\n", path)
	return nil
}

func runMP3(args []string) error {
	fs := flag.NewFlagSet("mp3", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	out := fs.String("out", "", "输出 .mp3 路径（相对路径写入 -dir）")
	outDir := fs.String("dir", "output", "相对 -out 的根目录，默认 output")
	dur := fs.String("duration", "60s", "时长")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *out == "" {
		return errors.New("需要 -out")
	}
	sec, err := parseDurationToSeconds(*dur)
	if err != nil {
		return err
	}
	if err := requireFFmpeg(); err != nil {
		return err
	}
	path, err := resolveOutputPath(*outDir, *out)
	if err != nil {
		return err
	}
	cmd := exec.Command("ffmpeg",
		"-y",
		"-f", "lavfi", "-i", "anullsrc=r=44100:cl=stereo",
		"-t", fmt.Sprintf("%f", sec),
		"-c:a", "libmp3lame", "-b:a", "128k",
		path,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg: %w", err)
	}
	fmt.Printf("已生成 MP3: %s\n", path)
	return nil
}

func parseWxH(s string) (w, h int, err error) {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(s)), "x")
	if len(parts) != 2 {
		return 0, 0, errors.New("-width 格式应为 WxH，如 1280x720")
	}
	w, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	h, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}
	if w < 16 || h < 16 {
		return 0, 0, errors.New("宽高过小")
	}
	return w, h, nil
}

// parseDurationToSeconds 支持 30s、5m、1.5h、纯数字视为秒
func parseDurationToSeconds(s string) (float64, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, errors.New("空时长")
	}
	if !strings.ContainsAny(s, "smh") {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		if v <= 0 {
			return 0, errors.New("时长必须为正")
		}
		return v, nil
	}
	var num strings.Builder
	var unit byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' || c == '.' {
			num.WriteByte(c)
			continue
		}
		if c == 's' || c == 'm' || c == 'h' {
			if unit != 0 {
				return 0, errors.New("只能有一个时间单位")
			}
			unit = c
			continue
		}
		return 0, fmt.Errorf("非法字符 %q", c)
	}
	if num.Len() == 0 || unit == 0 {
		return 0, errors.New("时长格式应为 数字+单位，如 90s、2m")
	}
	v, err := strconv.ParseFloat(num.String(), 64)
	if err != nil {
		return 0, err
	}
	if v <= 0 {
		return 0, errors.New("时长必须为正")
	}
	switch unit {
	case 's':
		return v, nil
	case 'm':
		return v * 60, nil
	case 'h':
		return v * 3600, nil
	default:
		return 0, errors.New("单位仅支持 s/m/h")
	}
}

func requireFFmpeg() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return errors.New("未找到 ffmpeg，请先安装并加入 PATH（https://ffmpeg.org）")
	}
	return nil
}

// resolveOutputPath 将相对路径的 out 放到 outDir 下（默认 output/），并创建父目录。
// out 为绝对路径时直接使用 outDir 不参与拼接。
func resolveOutputPath(outDir, out string) (string, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return "", errors.New("空输出路径")
	}
	if filepath.IsAbs(out) {
		p := filepath.Clean(out)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return "", err
		}
		return p, nil
	}
	baseAbs, err := filepath.Abs(filepath.Clean(outDir))
	if err != nil {
		return "", err
	}
	joined := filepath.Join(baseAbs, filepath.Clean(out))
	finalAbs, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(baseAbs, finalAbs)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New("-out 相对路径必须在 -dir 目录内（不能使用 .. 跳出）")
	}
	if err := os.MkdirAll(filepath.Dir(finalAbs), 0o755); err != nil {
		return "", err
	}
	return finalAbs, nil
}
