package service_test

import (
	"testing"

	"FileFlow/model"
	"FileFlow/service"
)

func TestMatchRuleAgainst_basicGlob(t *testing.T) {
	rules := []model.Rule{
		{ID: 1, Pattern: "*.mp4", Priority: 10, Enabled: true},
		{ID: 2, Pattern: "*.txt", Priority: 10, Enabled: true},
	}

	got := service.MatchRuleAgainst("test.mp4", rules)
	if got == nil || got.ID != 1 {
		t.Errorf("expected rule 1 (pattern=*.mp4) to match test.mp4, got %+v", got)
	}

	got = service.MatchRuleAgainst("test.txt", rules)
	if got == nil || got.ID != 2 {
		t.Errorf("expected rule 2 (pattern=*.txt) to match test.txt, got %+v", got)
	}

	got = service.MatchRuleAgainst("test.png", rules)
	if got != nil {
		t.Errorf("expected no match for test.png, got rule %+v", got)
	}
}

func TestMatchRuleAgainst_braceExpansion(t *testing.T) {
	rules := []model.Rule{
		{ID: 1, Pattern: "*.{mp4,mkv,avi}", Priority: 10, Enabled: true},
	}

	got := service.MatchRuleAgainst("movie.mp4", rules)
	if got == nil || got.ID != 1 {
		t.Errorf("expected rule to match movie.mp4, got %+v", got)
	}

	got = service.MatchRuleAgainst("movie.mkv", rules)
	if got == nil || got.ID != 1 {
		t.Errorf("expected rule to match movie.mkv, got %+v", got)
	}

	got = service.MatchRuleAgainst("movie.avi", rules)
	if got == nil || got.ID != 1 {
		t.Errorf("expected rule to match movie.avi, got %+v", got)
	}

	got = service.MatchRuleAgainst("movie.png", rules)
	if got != nil {
		t.Errorf("expected no match for movie.png, got rule %+v", got)
	}
}

func TestMatchRuleAgainst_pathWithWildcard(t *testing.T) {
	rules := []model.Rule{
		{ID: 1, Pattern: "*/movies/*.mkv", Priority: 10, Enabled: true},
	}

	got := service.MatchRuleAgainst("downloads/movies/test.mkv", rules)
	if got == nil || got.ID != 1 {
		t.Errorf("expected rule to match downloads/movies/test.mkv, got %+v", got)
	}

	got = service.MatchRuleAgainst("downloads/test.mkv", rules)
	if got != nil {
		t.Errorf("expected no match for downloads/test.mkv, got rule %+v", got)
	}
}

func TestMatchRuleAgainst_priorityOrder(t *testing.T) {
	rules := []model.Rule{
		{ID: 2, Pattern: "*.mp4", Priority: 10, Enabled: true},
		{ID: 1, Pattern: "*.mp4", Priority: 5, Enabled: true},
	}

	// priority 更高（10 > 5）的优先
	got := service.MatchRuleAgainst("test.mp4", rules)
	if got == nil || got.ID != 2 {
		t.Errorf("expected rule 2 (priority=10) to match first, got %+v", got)
	}

	// 同优先级时 id 大的后匹配，所以 id=1 的先返回
	rules = []model.Rule{
		{ID: 1, Pattern: "*.txt", Priority: 10, Enabled: true},
		{ID: 3, Pattern: "*.txt", Priority: 10, Enabled: true},
	}
	got = service.MatchRuleAgainst("test.txt", rules)
	if got == nil || got.ID != 1 {
		t.Errorf("expected rule 1 to match first when same priority, got %+v", got)
	}
}

func TestMatchRuleAgainst_disabledRule(t *testing.T) {
	rules := []model.Rule{
		{ID: 1, Pattern: "*.mp4", Priority: 10, Enabled: false},
		{ID: 2, Pattern: "*", Priority: 5, Enabled: true},
	}

	got := service.MatchRuleAgainst("test.mp4", rules)
	if got == nil || got.ID != 2 {
		t.Errorf("expected enabled rule 2 to match, got %+v", got)
	}
}

func TestMatchRuleAgainst_noMatch(t *testing.T) {
	rules := []model.Rule{
		{ID: 1, Pattern: "*.mp4", Priority: 10, Enabled: true},
	}

	got := service.MatchRuleAgainst("test.txt", rules)
	if got != nil {
		t.Errorf("expected nil for no match, got %+v", got)
	}
}
