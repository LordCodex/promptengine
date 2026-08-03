package filesystem

import (
	"os"
	"testing"
)

func TestMockFileSystem_Exists(t *testing.T) {
	fs := NewMockFileSystem()
	testPath := "docs/PRD.md"

	if fs.Exists(testPath) {
		t.Fatalf("Expected file '%s' not to exist initially", testPath)
	}

	err := fs.WriteFile(testPath, []byte("PRD Content"), 0644)
	if err != nil {
		t.Fatalf("Expected no error writing file, got: %v", err)
	}

	if !fs.Exists(testPath) {
		t.Errorf("Expected file '%s' to exist after writing", testPath)
	}
}

func TestMockFileSystem_IsSafePath(t *testing.T) {
	fs := NewMockFileSystem()

	if !fs.IsSafePath("base", "base/docs/PRD.md") {
		t.Errorf("Expected path base/docs/PRD.md to be safe")
	}

	if fs.IsSafePath("base", "base/../outside.md") {
		t.Errorf("Expected path base/../outside.md to be unsafe due to relative traversal")
	}
}

func TestMockFileSystem_IsSafePathAllowsDotDotPrefixName(t *testing.T) {
	fs := NewMockFileSystem()

	if !fs.IsSafePath("base", "base/..not-parent/file.md") {
		t.Fatal("path segment that merely starts with dots should be safe inside base")
	}
}

func TestOSFileSystem_IsSafePathRejectsSymlinkEscape(t *testing.T) {
	dir := t.TempDir()
	outside := t.TempDir()
	link := dir + "/link"
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	fs := &OSFileSystem{}
	if fs.IsSafePath(dir, link+"/secret.txt") {
		t.Fatal("symlink escaping the base directory should be unsafe")
	}
}
