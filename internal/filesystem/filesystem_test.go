package filesystem

import (
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
