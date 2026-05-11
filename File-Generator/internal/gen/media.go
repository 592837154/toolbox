package gen

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// GenMP4 调用 ffmpeg 生成黑屏 H.264 MP4。
func GenMP4(path string, sec float64, w, h int) error {
	if sec <= 0 {
		return errors.New("mp4 时长必须大于 0")
	}
	if w < 16 || h < 16 {
		return errors.New("mp4 宽高过小")
	}
	if err := requireFFmpeg(); err != nil {
		return err
	}
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

// GenMP3 调用 ffmpeg 生成静音立体声 MP3。
func GenMP3(path string, sec float64) error {
	if sec <= 0 {
		return errors.New("mp3 时长必须大于 0")
	}
	if err := requireFFmpeg(); err != nil {
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

func requireFFmpeg() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return errors.New("未找到 ffmpeg，请先安装并加入 PATH（https://ffmpeg.org）")
	}
	return nil
}
