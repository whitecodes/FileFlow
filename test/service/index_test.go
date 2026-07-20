package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"FileFlow/service"
)

func TestMoveFile_withIndex(t *testing.T) {
	targetDir := t.TempDir()

	// 准备源文件
	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "Rick.and.Morty.S09E03.1080p.x265-ELiTE.mkv")
	if err := os.WriteFile(src, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	// 目标目录已有 S01, S02
	os.WriteFile(filepath.Join(targetDir, "Rick.and.Morty.S09E01.1080p.x265-ELiTE.mkv"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(targetDir, "Rick.and.Morty.S09E02.1080p.x265-ELiTE.mkv"), []byte("x"), 0644)

	dst, err := service.MoveFile(src, targetDir, "Rick.and.Morty.S09E{index}.1080p.x265-ELiTE.mkv")
	if err != nil {
		t.Fatalf("MoveFile err: %v", err)
	}

	want := filepath.Join(targetDir, "Rick.and.Morty.S09E03.1080p.x265-ELiTE.mkv")
	if dst != want {
		t.Errorf("expected %q, got %q", want, dst)
	}
}

func TestMoveFile_withIndex_emptyDir(t *testing.T) {
	targetDir := t.TempDir()

	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "Rick.and.Morty.S09E05.1080p.x265-ELiTE.mkv")
	if err := os.WriteFile(src, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	// 目标目录为空 → 第一条为 01
	dst, err := service.MoveFile(src, targetDir, "Rick.and.Morty.S09E{index}.1080p.x265-ELiTE.mkv")
	if err != nil {
		t.Fatalf("MoveFile err: %v", err)
	}

	want := filepath.Join(targetDir, "Rick.and.Morty.S09E01.1080p.x265-ELiTE.mkv")
	if dst != want {
		t.Errorf("expected %q, got %q", want, dst)
	}
}

func TestResolveIndex(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "show.S01E01.1080p.mkv"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(dir, "show.S01E02.1080p.mkv"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(dir, "show.S01E05.1080p.mkv"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(dir, "other.txt"), []byte("x"), 0644) // 不匹配模板

	got := service.ResolveIndex(dir, "show.S01E{index}.1080p.mkv")
	// max existing = 05 → next = 06
	if got != 6 {
		t.Errorf("expected 6, got %d", got)
	}
}

func TestResolveIndex_emptyDir(t *testing.T) {
	dir := t.TempDir()
	got := service.ResolveIndex(dir, "show.S01E{index}.1080p.mkv")
	if got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestApplyTemplate_withIndex(t *testing.T) {
	got := service.ApplyTemplate("Rick.and.Morty.S09E03.1080p.x265-ELiTE.mkv", "Rick.and.Morty.S09E{index}.1080p.x265-ELiTE.mkv")
	// {index} 在源文件名解析阶段还不知道值，但模板本身应保持不变（由 MoveFile 注入）
	want := "Rick.and.Morty.S09E{index}.1080p.x265-ELiTE.mkv"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
