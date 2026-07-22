package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// MoveFile 移动文件到目标路径。优先 os.Rename，跨设备时用 copy+delete。
func MoveFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// Cross-device: copy then delete
	fmt.Printf("[move] cross-device, fallback to copy: %s -> %s\n", src, dst)

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dest: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	if err := os.Remove(src); err != nil {
		return fmt.Errorf("remove source after copy: %w", err)
	}

	return nil
}
