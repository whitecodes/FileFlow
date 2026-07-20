package service

import (
	"fmt"
	"log"

	"FileFlow/model"
)

// ProcessFile 执行完整流程：在 searchDir 中定位文件 → 匹配规则 → 移动到目标目录。
// 如果没有匹配的规则，返回 (nil, nil)。
// 如果文件不存在或移动失败，返回错误。
func ProcessFile(searchDir, fileName string) (*ProcessResult, error) {
	filePath, err := FindFile(searchDir, fileName)
	if err != nil {
		return nil, fmt.Errorf("find file: %w", err)
	}

	rule, err := MatchRule(fileName)
	if err != nil {
		return nil, fmt.Errorf("match rule: %w", err)
	}
	if rule == nil {
		log.Printf("[process] no matching rule for %s", fileName)
		return nil, nil
	}

	log.Printf("[process] matched rule=%q, moving %s -> %s", rule.Name, filePath, rule.TargetDir)

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

type ProcessResult struct {
	Rule    *model.Rule
	SrcPath string
	DstPath string
}
