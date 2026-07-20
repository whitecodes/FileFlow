package service_test

import (
	"testing"
	"time"

	"FileFlow/db"
	"FileFlow/model"
	"FileFlow/service"
)

func TestUpdateRule_success(t *testing.T) {
	db.Init("file:test_update_rule?mode=memory&cache=shared")
	defer db.Close()
	db.DB.Exec("DELETE FROM rules")

	// Create a rule first
	id, err := service.CreateRule(&model.Rule{
		Name:           "old name",
		Pattern:        "*.txt",
		TargetDir:      "/old/path",
		RenameTemplate: "{title}.{ext}",
		Priority:       5,
		Enabled:        true,
	})
	if err != nil {
		t.Fatalf("create rule: %v", err)
	}

	// Track time before update
	before := time.Now().Truncate(time.Second)

	// Update the rule
	err = service.UpdateRule(id, &model.Rule{
		Name:           "new name",
		Pattern:        "*.go",
		TargetDir:      "/new/path",
		RenameTemplate: "{title}.{ext}",
		Priority:       10,
		Enabled:        false,
	})
	if err != nil {
		t.Fatalf("update rule: %v", err)
	}

	// Fetch and verify
	updated, err := service.GetRule(id)
	if err != nil {
		t.Fatalf("get updated rule: %v", err)
	}
	if updated.Name != "new name" {
		t.Errorf("name = %q, want %q", updated.Name, "new name")
	}
	if updated.Pattern != "*.go" {
		t.Errorf("pattern = %q, want %q", updated.Pattern, "*.go")
	}
	if updated.TargetDir != "/new/path" {
		t.Errorf("target_dir = %q, want %q", updated.TargetDir, "/new/path")
	}
	if updated.Priority != 10 {
		t.Errorf("priority = %d, want 10", updated.Priority)
	}
	if updated.Enabled != false {
		t.Errorf("enabled = %v, want false", updated.Enabled)
	}
	// SQLite CURRENT_TIMESTAMP has second granularity
	if updated.UpdatedAt.Before(before) {
		t.Errorf("updated_at should not be before update call: was %v, before was %v",
			updated.UpdatedAt, before)
	}
}
