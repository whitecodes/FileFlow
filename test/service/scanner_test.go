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
