package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// 在aria2c下载完成的时候调用
// 参数: GID file_count first_file_path
// 流程: SSD→HDD → 加载YAML规则 → 匹配文件 → 移动到最终目录
func main() {
	flag.Parse()
	if len(flag.Args()) <= 2 {
		fmt.Println("args less than 2, not file download")
		return
	}

	fileCount := flag.Args()[1]
	firstFilePath := flag.Args()[2]

	fmt.Println("file_count:", fileCount)
	fmt.Println("first_file_path:", firstFilePath)

	if firstFilePath == "" || fileCount == "" {
		fmt.Println("empty args, return")
		return
	}

	fileCountNum, err := strconv.Atoi(fileCount)
	if err != nil || fileCountNum <= 0 {
		fmt.Println("invalid file_count:", fileCount)
		return
	}
	fmt.Printf("file_count_num: %d\n", fileCountNum)

	// --- 路径配置 ---
	sourceFolder := "/aria2/ssd/"
	targetFolder := "/aria2/hdd/"

	// --- 解析源文件/目录名 ---
	parts := strings.Split(firstFilePath, "/")
	if len(parts) < 3 {
		fmt.Println("first_file_path not right, return")
		return
	}
	entryName := parts[3] // 文件名或目录名
	sourcePath := sourceFolder + entryName
	targetPath := targetFolder + entryName

	fmt.Printf("source: %s\n", sourcePath)
	fmt.Printf("target: %s\n", targetPath)

	// --- SSD→HDD ---
	if _, err := os.Stat(sourcePath); err != nil {
		fmt.Printf("source not exist: %v\n", err)
		return
	}

	if err := os.Rename(sourcePath, targetPath); err != nil {
		fmt.Printf("move error, try copy: %v\n", err)
		// copy fallback: src(sourcePath) → dst(targetPath)
		if err := copyRecursive(targetPath, sourcePath); err != nil {
			fmt.Printf("copy error: %v\n", err)
			return
		}
		os.RemoveAll(sourcePath)
	}

	fmt.Println("SSD→HDD done")

	// chmod
	exec.Command("chmod", "-R", "o+w", targetPath).Run()

	// --- 加载配置并匹配 ---
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)
	cfgPath := os.Getenv("FILEFLOW_CONFIG")
	if cfgPath == "" {
		cfgPath = filepath.Join(execDir, "fileflow.yaml")
	}
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		fmt.Printf("load config error: %v\n", err)
		return
	}
	fmt.Printf("loaded %d rules\n", len(cfg.Rules))

	// --- 判断是文件还是目录 ---
	info, err := os.Stat(targetPath)
	if err != nil {
		fmt.Printf("stat error: %v\n", err)
		return
	}

	MatchAndMove(targetPath, info.IsDir(), cfg.Rules)
}

// copyRecursive 递归复制目录（跨设备 fallback）
func copyRecursive(dst, src string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		dest := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}
		return copyFile(dest, path)
	})
}

func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
