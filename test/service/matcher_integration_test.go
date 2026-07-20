package service_test

import (
	"testing"

	"FileFlow/db"
	"FileFlow/model"
	"FileFlow/service"
)

func TestMatchRule_integration(t *testing.T) {
	// Use in-memory SQLite
	db.Init("file:test_match?mode=memory&cache=shared")
	defer db.Close()

	// Insert test rules
	_, _ = db.DB.Exec("DELETE FROM rules")
	for _, r := range []model.Rule{
		{Name: "videos", Pattern: "*.mp4", TargetDir: "/media/videos", RenameTemplate: "{title}.{ext}", Priority: 5, Enabled: true},
		{Name: "high priority", Pattern: "*.mp4", TargetDir: "/media/priority", RenameTemplate: "{title}.{ext}", Priority: 10, Enabled: true},
		{Name: "disabled", Pattern: "*", TargetDir: "/media/disabled", RenameTemplate: "{title}.{ext}", Priority: 100, Enabled: false},
	} {
		_, err := db.DB.Exec(
			`INSERT INTO rules (name, pattern, target_dir, rename_template, priority, enabled)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			r.Name, r.Pattern, r.TargetDir, r.RenameTemplate, r.Priority, r.Enabled,
		)
		if err != nil {
			t.Fatalf("insert rule: %v", err)
		}
	}

	got, err := service.MatchRule("test.mp4")
	if err != nil {
		t.Fatalf("MatchRule err: %v", err)
	}
	if got == nil {
		t.Fatal("expected a match, got nil")
	}
	if got.Name != "high priority" {
		t.Errorf("expected highest priority rule 'high priority', got %q", got.Name)
	}

	got, err = service.MatchRule("test.txt")
	if err != nil {
		t.Fatalf("MatchRule err: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for no match, got %+v", got)
	}
}
