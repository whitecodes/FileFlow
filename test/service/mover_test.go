package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"FileFlow/service"
)

func TestRenameFile_basicTemplate(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "My Movie.mkv")
	if err := os.WriteFile(src, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	dst, err := service.RenameFile(src, "{title}.{ext}")
	if err != nil {
		t.Fatalf("RenameFile err: %v", err)
	}

	want := filepath.Join(dir, "My Movie.mkv")
	if dst != want {
		t.Errorf("expected %q, got %q", want, dst)
	}
}

func TestRenameFile_changesName(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "My Movie.mkv")
	if err := os.WriteFile(src, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	dst, err := service.RenameFile(src, "newname.{ext}")
	if err != nil {
		t.Fatalf("RenameFile err: %v", err)
	}

	want := filepath.Join(dir, "newname.mkv")
	if dst != want {
		t.Errorf("expected %q, got %q", want, dst)
	}

	// Verify rename happened on disk
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Error("original file should not exist after rename")
	}
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Error("renamed file should exist on disk")
	}
}

func TestApplyTemplate_replacesPlaceholders(t *testing.T) {
	tests := []struct {
		file     string
		template string
		want     string
	}{
		{"hello.mp4", "{title}.{ext}", "hello.mp4"},
		{"hello.mp4", "{file}", "hello.mp4"},
		{"hello.mp4", "prefix-{title}.{ext}", "prefix-hello.mp4"},
		{"hello.world.tar.gz", "{title}.{ext}", "hello.world.tar.gz"},
	}
	for _, tt := range tests {
		got := service.ApplyTemplate(tt.file, tt.template)
		if got != tt.want {
			t.Errorf("ApplyTemplate(%q, %q) = %q, want %q", tt.file, tt.template, got, tt.want)
		}
	}
}
