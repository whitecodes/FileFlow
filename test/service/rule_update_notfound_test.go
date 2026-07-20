package service_test

import (
	"database/sql"
	"testing"

	"FileFlow/db"
	"FileFlow/model"
	"FileFlow/service"
)

func TestUpdateRule_notFound(t *testing.T) {
	db.Init("file:test_update_rule?mode=memory&cache=shared")
	defer db.Close()
	db.DB.Exec("DELETE FROM rules")

	err := service.UpdateRule(999, &model.Rule{
		Name:           "nope",
		Pattern:        "*",
		TargetDir:      "/dev/null",
		RenameTemplate: "{title}.{ext}",
		Priority:       0,
		Enabled:        true,
	})
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}
