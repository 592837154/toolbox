package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if err := cleanAll(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func cleanAll() error {
	baseDir, err := executableDir()
	if err != nil {
		return err
	}

	convertDir := resolveConvertDir(baseDir)
	inputDir := filepath.Join(convertDir, "input")
	outputDir := filepath.Join(convertDir, "output")

	for _, dir := range []string{inputDir, outputDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
		if err := cleanDir(dir); err != nil {
			return err
		}
		fmt.Printf("[CLEAN] %s\n", dir)
	}

	fmt.Println("Done. convert/input and convert/output are empty.")
	return nil
}

func resolveConvertDir(baseDir string) string {
	if filepath.Base(baseDir) == "convert" {
		return baseDir
	}
	return filepath.Clean(filepath.Join(baseDir, "..", "convert"))
}

func cleanDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove %s: %w", path, err)
		}
	}
	return nil
}

func executableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locate executable: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		resolved = exePath
	}
	return filepath.Dir(resolved), nil
}
