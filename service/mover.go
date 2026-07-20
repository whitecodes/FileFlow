package service

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ApplyTemplate 将命名模板中的占位符替换为文件信息。
// 支持: {title} {ext} {file}
// {index} 不被替换，由 MoveFile/ResolveIndex 处理。
func ApplyTemplate(fileName, template string) string {
	ext := filepath.Ext(fileName)
	title := strings.TrimSuffix(fileName, ext)
	ext = strings.TrimPrefix(ext, ".")

	result := template
	result = strings.ReplaceAll(result, "{file}", fileName)
	result = strings.ReplaceAll(result, "{title}", title)
	result = strings.ReplaceAll(result, "{ext}", ext)
	return result
}

// ResolveIndex 扫描 targetDir 中匹配 template（{index} 替换为 \d+ 做 regex）的文件，
// 返回当前最大序号 + 1。目录为空或无匹配时返回 1。
func ResolveIndex(targetDir, template string) int {
	pattern := strings.ReplaceAll(regexp.QuoteMeta(template), `\{index\}`, `(\d+)`)
	re := regexp.MustCompile(pattern)

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return 1
	}

	maxIdx := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		m := re.FindStringSubmatch(e.Name())
		if len(m) > 1 {
			if idx, err := strconv.Atoi(m[1]); err == nil && idx > maxIdx {
				maxIdx = idx
			}
		}
	}
	return maxIdx + 1
}

// RenameFile 按 template 重命名文件，返回新路径。
func RenameFile(src, template string) (string, error) {
	dir := filepath.Dir(src)
	oldName := filepath.Base(src)
	newName := ApplyTemplate(oldName, template)
	dst := filepath.Join(dir, newName)

	if err := os.Rename(src, dst); err != nil {
		return "", fmt.Errorf("rename %s -> %s: %w", src, dst, err)
	}
	return dst, nil
}

// MoveFile 先将文件按 template 重命名，再移动到 targetDir 目录下。
// {index} 占位符会被替换为根据目标目录已有文件计算出的自增序号。
func MoveFile(src, targetDir, template string) (string, error) {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", targetDir, err)
	}

	idx := ResolveIndex(targetDir, template)
	indexStr := fmt.Sprintf("%02d", idx)

	newName := ApplyTemplate(filepath.Base(src), template)
	newName = strings.ReplaceAll(newName, "{index}", indexStr)
	dst := filepath.Join(targetDir, newName)

	if err := os.Rename(src, dst); err != nil {
		return "", fmt.Errorf("move %s -> %s: %w", src, dst, err)
	}
	return dst, nil
}
