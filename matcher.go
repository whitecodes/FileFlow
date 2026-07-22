package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// MatchResult 匹配结果，包含命中的规则和捕获的 source.N 值
type MatchResult struct {
	Rule    *Rule
	Sources []string // [0] = 完整文件名, [1..] = 通配符捕获
}

// globToRegex 将 glob pattern 转换为正则，同时处理 brace expansion
// 返回展开后的 pattern 列表
func globToRegex(pattern string) []string {
	// Handle brace expansion {a,b,c}
	patterns := expandBraces(pattern)
	return patterns
}

func expandBraces(pattern string) []string {
	start := strings.IndexByte(pattern, '{')
	end := strings.IndexByte(pattern, '}')
	if start == -1 || end == -1 || end <= start {
		return []string{pattern}
	}

	prefix := pattern[:start]
	body := pattern[start+1 : end]
	suffix := pattern[end+1:]

	parts := strings.Split(body, ",")
	result := make([]string, len(parts))
	for i, p := range parts {
		result[i] = prefix + p + suffix
	}
	return result
}

// patternToRegex 将单个 glob pattern 转换为带捕获组的正则
// * → (.*)
// ? → (.)
// 其他字符原样转义
func patternToRegex(pattern string) *regexp.Regexp {
	var buf strings.Builder
	buf.WriteString("^")
	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]
		switch ch {
		case '*':
			buf.WriteString("(.*)")
		case '?':
			buf.WriteString("(.)")
		case '.', '+', '^', '$', '[', ']', '(', ')', '|', '\\':
			buf.WriteByte('\\')
			buf.WriteByte(ch)
		default:
			buf.WriteByte(ch)
		}
	}
	buf.WriteString("$")
	return regexp.MustCompile(buf.String())
}

// MatchFile 将文件名与所有规则逐一匹配，返回第一个匹配的结果
func MatchFile(fileName string, rules []Rule) *MatchResult {
	for _, rule := range rules {
		patterns := globToRegex(rule.Pattern)
		for _, p := range patterns {
			re := patternToRegex(p)
			m := re.FindStringSubmatch(fileName)
			if m != nil {
				return &MatchResult{
					Rule:    &rule,
					Sources: m[1:], // m[0] 是整个匹配，m[1..] 是捕获组
				}
			}
		}
	}
	return nil
}

// ResolveTargetName 根据匹配结果和目标文件名模板，解析最终文件名
// 占位符:
//   {source.0}  - 完整源文件名
//   {source.N}  - 第 N 个通配符匹配内容 (N>=1)
//   {target.index} - 目标目录中已存在文件的递进序号
func ResolveTargetName(sourceFileName, targetName string, sources []string, targetDir string) string {
	result := targetName
	result = strings.ReplaceAll(result, "{source.0}", sourceFileName)
	for i, s := range sources {
		result = strings.ReplaceAll(result, fmt.Sprintf("{source.%d}", i+1), s)
	}
	if strings.Contains(result, "{target.index}") {
		idx := resolveIndex(targetDir, targetName)
		result = strings.ReplaceAll(result, "{target.index}", fmt.Sprintf("%02d", idx))
	}
	return result
}

// resolveIndex 扫描目标目录中匹配 targetName 的文件，返回最大序号+1
func resolveIndex(dir, targetName string) int {
	// Build regex from target_name template
	reStr := strings.ReplaceAll(regexp.QuoteMeta(targetName), `\{source\.\d+\}`, `.*`)
	reStr = strings.ReplaceAll(reStr, `\{target\.index\}`, `(\d+)`)
	re := regexp.MustCompile(reStr)

	entries, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return 1
	}

	maxIdx := 0
	for _, e := range entries {
		name := filepath.Base(e)
		m := re.FindStringSubmatch(name)
		if len(m) > 1 {
			idx := 0
			fmt.Sscanf(m[1], "%d", &idx)
			if idx > maxIdx {
				maxIdx = idx
			}
		}
	}
	return maxIdx + 1
}

// MatchAndMove 在 HDD 的目标上执行匹配和移动
// path 是 HDD 上已经存在的文件或目录路径
// fileName 是 aria2 传过来的原始文件名
func MatchAndMove(path string, isDir bool, rules []Rule) {
	if !isDir {
		matchAndMoveOne(path, rules)
		return
	}

	// 目录：递归所有文件逐个匹配
	filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		matchAndMoveOne(p, rules)
		return nil
	})
}

func matchAndMoveOne(filePath string, rules []Rule) {
	fileName := filepath.Base(filePath)
	result := MatchFile(fileName, rules)
	if result == nil {
		fmt.Printf("[skip] no matching rule for %s\n", fileName)
		return
	}

	rule := result.Rule
	newName := ResolveTargetName(fileName, rule.TargetName, result.Sources, rule.TargetDir)
	dst := filepath.Join(rule.TargetDir, newName)

	fmt.Printf("[match] %s -> %s (rule=%q)\n", fileName, dst, rule.Name)

	if err := MoveFile(filePath, dst); err != nil {
		fmt.Printf("[error] move %s -> %s: %v\n", filePath, dst, err)
	} else {
		fmt.Printf("[done] %s\n", dst)
	}
}
