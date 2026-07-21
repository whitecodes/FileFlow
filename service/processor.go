package service

import (
	"fmt"
	"log"
	"path/filepath"

	"FileFlow/model"
)

// ProcessResult 单次处理结果。
type ProcessResult struct {
	Rule    *model.Rule
	SrcPath string
	DstPath string
}

// ProcessFile 查找文件并匹配规则：
// - 传的是文件 → 匹配规则 → 移动
// - 传的是目录 → 递归查找所有媒体文件 → 逐个匹配规则 → 逐个移动
// 返回处理结果列表。无匹配时不返回错误，result 列表为空。
func ProcessFile(searchDir, fileName string) ([]*ProcessResult, error) {
	path, isDir, err := Locate(searchDir, fileName)
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}

	if !isDir {
		// 单个文件
		result, err := processSingleFile(path)
		if err != nil {
			return nil, err
		}
		if result == nil {
			return nil, nil
		}
		return []*ProcessResult{result}, nil
	}

	// 目录：列出所有文件逐个处理
	files, err := ListFiles(path)
	if err != nil {
		log.Printf("[process] no media files in directory %s: %v", path, err)
		return nil, nil
	}

	var results []*ProcessResult
	for _, f := range files {
		result, err := processSingleFile(f)
		if err != nil {
			log.Printf("[process] skip %s: %v", f, err)
			continue
		}
		if result != nil {
			results = append(results, result)
		}
	}
	return results, nil
}

// processSingleFile 处理单个文件：匹配规则，移动到目标目录。
func processSingleFile(filePath string) (*ProcessResult, error) {
	fileName := filepath.Base(filePath)
	rule, err := MatchRule(fileName)
	if err != nil {
		return nil, fmt.Errorf("match rule: %w", err)
	}
	if rule == nil {
		log.Printf("[process] no matching rule for %s", fileName)
		return nil, nil
	}

	log.Printf("[process] matched rule=%q for %s, moving -> %s", rule.Name, fileName, rule.TargetDir)

	dst, err := MoveFile(filePath, rule.TargetDir, rule.RenameTemplate)
	if err != nil {
		return nil, fmt.Errorf("move file: %w", err)
	}

	return &ProcessResult{
		Rule:    rule,
		SrcPath: filePath,
		DstPath: dst,
	}, nil
}
