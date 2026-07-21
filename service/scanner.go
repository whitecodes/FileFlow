package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// Locate 在 searchDir 中查找 fileName。返回完整路径，以及是否是目录。
func Locate(searchDir, fileName string) (string, bool, error) {
	log.Printf("[scanner] Locate searchDir=%q fileName=%q", searchDir, fileName)

	// Try exact file match
	var found string
	err := filepath.WalkDir(searchDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() || d.Name() != fileName {
			return nil
		}
		found = path
		return filepath.SkipAll
	})
	if err != nil {
		return "", false, fmt.Errorf("walk error: %w", err)
	}
	if found != "" {
		log.Printf("[scanner] found file: %q", found)
		return found, false, nil
	}

	// Try directory match
	dirPath := filepath.Join(searchDir, fileName)
	if stat, err := os.Stat(dirPath); err == nil && stat.IsDir() {
		log.Printf("[scanner] found directory: %q", dirPath)
		return dirPath, true, nil
	}

	return "", false, fmt.Errorf("not found: %s in %s", fileName, searchDir)
}

// ListFiles 递归扫描 dir 及其子目录，返回所有文件的完整路径（按大小降序）。
func ListFiles(dir string) ([]string, error) {
	var files []string
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if len(files) == 0 {
		return nil, fmt.Errorf("no files in %s", dir)
	}
	sort.Slice(files, func(i, j int) bool {
		fi, _ := os.Stat(files[i])
		fj, _ := os.Stat(files[j])
		return fi.Size() > fj.Size()
	})
	return files, nil
}
