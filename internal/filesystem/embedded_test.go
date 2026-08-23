package filesystem

import (
	"testing"
	"testing/fstest"
)

func TestOverlayReadsBundledLibraryWithoutExposingItToProjectWalks(t *testing.T) {
	project := NewMockFileSystem()
	if err := project.WriteFile("app/main.go", []byte("package app"), 0644); err != nil {
		t.Fatal(err)
	}
	library := NewEmbeddedFileSystem(fstest.MapFS{
		"core/05-universal-coding-standards.md": &fstest.MapFile{Data: []byte("standards")},
	})
	overlay := NewOverlayFileSystem(project, library)

	data, err := overlay.ReadFile("core/05-universal-coding-standards.md")
	if err != nil {
		t.Fatalf("read bundled standard: %v", err)
	}
	if string(data) != "standards" {
		t.Fatalf("unexpected bundled content: %q", data)
	}

	entries, err := overlay.ReadDir(".")
	if err != nil {
		t.Fatalf("read project root: %v", err)
	}
	for _, entry := range entries {
		if entry.Name() == "core" {
			t.Fatal("bundled standards must not appear as project files during discovery")
		}
	}
}

func TestOverlayWritesOnlyToProjectFilesystem(t *testing.T) {
	project := NewMockFileSystem()
	library := NewEmbeddedFileSystem(fstest.MapFS{
		"core/rule.md": &fstest.MapFile{Data: []byte("rule")},
	})
	overlay := NewOverlayFileSystem(project, library)

	if err := overlay.WriteFile("docs/Architecture.md", []byte("project architecture"), 0644); err != nil {
		t.Fatalf("write project file: %v", err)
	}
	if !project.Exists("docs/Architecture.md") {
		t.Fatal("expected writes to go to the project filesystem")
	}
}
