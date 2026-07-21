package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"FileFlow/service"
)

func TestFindFile_findsFileInDir(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "hello.mp4")
	if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := service.FindFile(dir, "hello.mp4")
	if err != nil {
		t.Fatalf("FindFile err: %v", err)
	}
	if got != src {
		t.Errorf("expected %q, got %q", src, got)
	}
}

func TestFindFile_returnsErrorWhenNotFound(t *testing.T) {
	dir := t.TempDir()

	_, err := service.FindFile(dir, "nope.mkv")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFindFile_findsMediaInNamedDir(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "Some.Movie.2024.1080p.x265-GROUP")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	// Create a smaller file and a larger one — should pick the largest
	os.WriteFile(filepath.Join(subDir, "Some.Movie.2024.1080p.x265-GROUP.nfo"), []byte("info"), 0644)
	os.WriteFile(filepath.Join(subDir, "Some.Movie.2024.1080p.x265-GROUP.mkv"), []byte("larger content here"), 0644)

	got, err := service.FindFile(dir, "Some.Movie.2024.1080p.x265-GROUP")
	if err != nil {
		t.Fatalf("FindFile err: %v", err)
	}
	if filepath.Base(got) != "Some.Movie.2024.1080p.x265-GROUP.mkv" {
		t.Errorf("expected the .mkv file, got %q", got)
	}
}
