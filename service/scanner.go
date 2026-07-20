package service

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindFile 在 searchDir 及其子目录中查找 fileName，返回完整路径。
// 先匹配文件名，匹配到第一个就返回。
func FindFile(searchDir, fileName string) (string, error) {
	var found string
	err := filepath.WalkDir(searchDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip inaccessible
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == fileName {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("walk error: %w", err)
	}
	if found == "" {
		return "", fmt.Errorf("file not found: %s in %s", fileName, searchDir)
	}
	return found, nil
}
