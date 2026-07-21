package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// FindFile 在 searchDir 中查找 fileName。
// 1. 先按文件名精确匹配文件
// 2. 如果没找到，查找同名目录，在目录中递归搜索媒体文件
// 返回第一个匹配的完整路径。
func FindFile(searchDir, fileName string) (string, error) {
	log.Printf("[scanner] FindFile searchDir=%q fileName=%q", searchDir, fileName)

	// Phase 1: exact file match
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
		return "", fmt.Errorf("walk error: %w", err)
	}
	if found != "" {
		log.Printf("[scanner] phase1: found file %q", found)
		return found, nil
	}
	log.Printf("[scanner] phase1: no exact file match for %q, trying directory", fileName)

	// Phase 2: look for a directory named fileName, find the largest media file inside
	dirPath := filepath.Join(searchDir, fileName)
	log.Printf("[scanner] phase2: checking directory path %q", dirPath)
	if stat, err := os.Stat(dirPath); err == nil && stat.IsDir() {
		log.Printf("[scanner] phase2: directory exists, scanning for media files")
		file, err := findMediaInDir(dirPath)
		if err != nil {
			log.Printf("[scanner] phase2: scan error: %v", err)
			return "", fmt.Errorf("find media in dir %s: %w", dirPath, err)
		}
		log.Printf("[scanner] phase2: found media file %q", file)
		return file, nil
	} else if err != nil {
		log.Printf("[scanner] phase2: stat error: %v", err)
	} else {
		log.Printf("[scanner] phase2: %q is not a directory", dirPath)
	}

	return "", fmt.Errorf("file not found: %s in %s", fileName, searchDir)
}

// findMediaInDir 在目录中递归查找视频文件，返回最大的那个。
// 支持的扩展名：.mkv .mp4 .avi .mov .ts .m2ts
func findMediaInDir(dir string) (string, error) {
	var candidates []string
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("[scanner] walk error in %s: %v", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(d.Name())
		switch ext {
		case ".mkv", ".mp4", ".avi", ".mov", ".ts", ".m2ts":
			log.Printf("[scanner] candidate: %s (ext=%s)", path, ext)
			candidates = append(candidates, path)
		}
		return nil
	})
	if len(candidates) == 0 {
		log.Printf("[scanner] no media files found in %s", dir)
		return "", fmt.Errorf("no media file found")
	}
	// Return the largest file
	sort.Slice(candidates, func(i, j int) bool {
		fi, _ := os.Stat(candidates[i])
		fj, _ := os.Stat(candidates[j])
		return fi.Size() > fj.Size()
	})
	log.Printf("[scanner] selected largest: %s", candidates[0])
	return candidates[0], nil
}
