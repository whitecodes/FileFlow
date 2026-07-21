package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"FileFlow/service"
)

func TestLocate_findsFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "hello.mp4")
	if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	got, isDir, err := service.Locate(dir, "hello.mp4")
	if err != nil {
		t.Fatalf("Locate err: %v", err)
	}
	if isDir {
		t.Error("expected isDir=false")
	}
	if got != src {
		t.Errorf("expected %q, got %q", src, got)
	}
}

func TestLocate_findsDir(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "Some.Movie.2024")
	os.MkdirAll(subDir, 0755)

	got, isDir, err := service.Locate(dir, "Some.Movie.2024")
	if err != nil {
		t.Fatalf("Locate err: %v", err)
	}
	if !isDir {
		t.Error("expected isDir=true for a directory")
	}
	if got != subDir {
		t.Errorf("expected %q, got %q", subDir, got)
	}
}

func TestLocate_returnsErrorWhenNotFound(t *testing.T) {
	dir := t.TempDir()

	_, _, err := service.Locate(dir, "nope.mkv")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestListFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "movie.mkv"), []byte("larger"), 0644)
	os.WriteFile(filepath.Join(dir, "movie.nfo"), []byte("info"), 0644)
	os.WriteFile(filepath.Join(dir, "screenshot.png"), []byte("img"), 0644)

	// Subdirectory with more files
	sub := filepath.Join(dir, "sub")
	os.MkdirAll(sub, 0755)
	os.WriteFile(filepath.Join(sub, "clip.mp4"), []byte("small"), 0644)
	os.WriteFile(filepath.Join(sub, "readme.txt"), []byte("hello"), 0644)

	files, err := service.ListFiles(dir)
	if err != nil {
		t.Fatalf("ListFiles err: %v", err)
	}
	if len(files) != 5 {
		t.Errorf("expected 5 files total, got %d: %v", len(files), files)
	}
	// Should be sorted by size descending
	if filepath.Base(files[0]) != "movie.mkv" {
		t.Errorf("expected movie.mkv first (largest), got %q", files[0])
	}
}
