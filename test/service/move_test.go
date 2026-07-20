package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"FileFlow/service"
)

func TestMoveFile_movesToTargetDir(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := filepath.Join(t.TempDir(), "movies")
	src := filepath.Join(srcDir, "test.mp4")
	if err := os.WriteFile(src, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	dst, err := service.MoveFile(src, dstDir, "{title}.{ext}")
	if err != nil {
		t.Fatalf("MoveFile err: %v", err)
	}

	want := filepath.Join(dstDir, "test.mp4")
	if dst != want {
		t.Errorf("expected %q, got %q", want, dst)
	}

	// Original gone, target exists
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Error("original should not exist after move")
	}
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Error("moved file should exist")
	}
}

func TestMoveFile_createsTargetDir(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := filepath.Join(t.TempDir(), "deep", "nested", "movies")
	src := filepath.Join(srcDir, "test.mp4")
	if err := os.WriteFile(src, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	dst, err := service.MoveFile(src, dstDir, "{title}.{ext}")
	if err != nil {
		t.Fatalf("MoveFile err: %v", err)
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Errorf("moved file should exist at %q", dst)
	}
}
