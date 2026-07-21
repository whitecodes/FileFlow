package service_test

import (
	"testing"

	"FileFlow/db"
	"FileFlow/service"
)

func TestRecordHistory_insertsRecord(t *testing.T) {
	db.Init("file:test_history?mode=memory&cache=shared")
	defer db.Close()
	db.DB.Exec("DELETE FROM history")

	err := service.RecordHistory("test.mp4", "file_moved", "movies", "/src/test.mp4", "/dst/test.mp4", "matched", "")
	if err != nil {
		t.Fatalf("RecordHistory err: %v", err)
	}

	records, err := service.ListHistory(10)
	if err != nil {
		t.Fatalf("ListHistory err: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].FileName != "test.mp4" {
		t.Errorf("FileName = %q, want %q", records[0].FileName, "test.mp4")
	}
	if records[0].Status != "matched" {
		t.Errorf("Status = %q, want %q", records[0].Status, "matched")
	}
}

func TestListHistory_returnsLatestFirst(t *testing.T) {
	db.Init("file:test_history?mode=memory&cache=shared")
	defer db.Close()
	db.DB.Exec("DELETE FROM history")

	for i := 0; i < 5; i++ {
		service.RecordHistory("file.mp4", "", "", "", "", "matched", "")
	}

	records, err := service.ListHistory(3)
	if err != nil {
		t.Fatalf("ListHistory err: %v", err)
	}
	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}
	// IDs should be descending (5, 4, 3)
	if records[0].ID < records[1].ID {
		t.Error("expected latest record first (descending order)")
	}
}

func TestRecordHistory_trimsToMaxRecords(t *testing.T) {
	db.Init("file:test_history?mode=memory&cache=shared")
	defer db.Close()
	db.DB.Exec("DELETE FROM history")

	for i := 0; i < 105; i++ {
		service.RecordHistory("file.mp4", "", "", "", "", "matched", "")
	}

	records, err := service.ListHistory(200)
	if err != nil {
		t.Fatalf("ListHistory err: %v", err)
	}
	if len(records) > 100 {
		t.Errorf("expected at most 100 records, got %d", len(records))
	}
}
