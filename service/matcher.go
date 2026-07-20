package service

import (
	"path/filepath"
	"strings"

	"FileFlow/model"
)

// expandBraces 展开大括号表达式，如 "*.{mp4,mkv}" → ["*.mp4", "*.mkv"]。
// 没有大括号时返回原 pattern。
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

// MatchRuleAgainst 在已按 priority DESC, id ASC 排序的 rules 中，
// 返回第一个匹配 filePath 的规则。没有匹配时返回 nil。
// 支持 glob 通配符和大括号展开（如 *.{mp4,mkv}）。
func MatchRuleAgainst(filePath string, rules []model.Rule) *model.Rule {
	for _, r := range rules {
		if !r.Enabled {
			continue
		}
		for _, pattern := range expandBraces(r.Pattern) {
			matched, err := filepath.Match(pattern, filePath)
			if err == nil && matched {
				return &r
			}
		}
	}
	return nil
}

// MatchRule 从数据库查询所有已启用的规则（按 priority DESC, id ASC 排序），
// 返回第一个匹配 filePath 的规则。没有匹配时返回 nil。
func MatchRule(filePath string) (*model.Rule, error) {
	rules, err := ListRules()
	if err != nil {
		return nil, err
	}
	return MatchRuleAgainst(filePath, rules), nil
}
