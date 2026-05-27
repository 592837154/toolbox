package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	inputDirName     = "input"
	outputDirName    = "output"
	referenceDocName = "reference.docx"
	pandocExeName    = "pandoc.exe"
)

type markdownFile struct {
	inputPath  string
	outputPath string
}

func main() {
	if err := convertAll(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func convertAll() error {
	baseDir, err := executableDir()
	if err != nil {
		return err
	}

	inputDir := filepath.Join(baseDir, inputDirName)
	outputDir := filepath.Join(baseDir, outputDirName)

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		return fmt.Errorf("create input directory: %w", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	pandocPath, err := findPandoc(baseDir)
	if err != nil {
		return err
	}
	fmt.Printf("Pandoc: %s\n", pandocPath)

	files, err := collectMarkdownFiles(inputDir, outputDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Printf("No .md files found in %s\n", inputDir)
		fmt.Println("Put Markdown files into the input folder, then run this program again.")
		return nil
	}

	referenceDoc := filepath.Join(baseDir, referenceDocName)
	useReferenceDoc := fileExists(referenceDoc)
	if useReferenceDoc {
		fmt.Printf("Reference doc: %s\n", referenceDoc)
	}

	var failed int
	for _, file := range files {
		if err := convertWithPandoc(pandocPath, file.inputPath, file.outputPath, referenceDoc, useReferenceDoc); err != nil {
			failed++
			fmt.Printf("[FAIL] %s: %v\n", displayPath(inputDir, file.inputPath), err)
			continue
		}
		fmt.Printf("[OK] %s -> %s\n", displayPath(inputDir, file.inputPath), displayPath(outputDir, file.outputPath))
	}

	if failed > 0 {
		return fmt.Errorf("%d file(s) failed", failed)
	}
	fmt.Printf("Done. Converted %d file(s). Output: %s\n", len(files), outputDir)
	return nil
}

func collectMarkdownFiles(inputDir, outputDir string) ([]markdownFile, error) {
	var files []markdownFile
	err := filepath.WalkDir(inputDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !isMarkdownFile(path) {
			return nil
		}

		relativePath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return err
		}
		outputRelativePath := strings.TrimSuffix(relativePath, filepath.Ext(relativePath)) + ".docx"
		files = append(files, markdownFile{
			inputPath:  path,
			outputPath: filepath.Join(outputDir, outputRelativePath),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan markdown files: %w", err)
	}
	return files, nil
}

func isMarkdownFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}

func convertWithPandoc(pandocPath, inputPath, outputPath, referenceDoc string, useReferenceDoc bool) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	args := []string{
		inputPath,
		"--from", "gfm+tex_math_dollars+emoji",
		"--to", "docx",
		"--output", outputPath,
		"--standalone",
		"--toc",
		"--toc-depth", "3",
	}
	if useReferenceDoc {
		args = append(args, "--reference-doc", referenceDoc)
	}

	cmd := exec.Command(pandocPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return fmt.Errorf("pandoc failed: %s", message)
	}
	return nil
}

func findPandoc(baseDir string) (string, error) {
	candidates := []string{
		filepath.Join(baseDir, pandocExeName),
		os.Getenv("PANDOC_PATH"),
	}

	if path, err := exec.LookPath("pandoc"); err == nil {
		candidates = append(candidates, path)
	}

	candidates = append(candidates, commonPandocPaths()...)

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if fileExists(candidate) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("pandoc not found. Install pandoc, add it to PATH, set PANDOC_PATH, or copy pandoc.exe next to this program")
}

func commonPandocPaths() []string {
	var paths []string
	for _, dir := range []string{
		os.Getenv("LOCALAPPDATA"),
		os.Getenv("ProgramFiles"),
		os.Getenv("ProgramFiles(x86)"),
	} {
		if dir == "" {
			continue
		}
		paths = append(paths, filepath.Join(dir, "Pandoc", pandocExeName))
	}

	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths,
			filepath.Join(home, "AppData", "Local", "Pandoc", pandocExeName),
			filepath.Join(home, "scoop", "apps", "pandoc", "current", pandocExeName),
		)
	}

	return paths
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

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func displayPath(baseDir, path string) string {
	relative, err := filepath.Rel(baseDir, path)
	if err != nil {
		return filepath.Base(path)
	}
	return relative
}
