package validation

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestValidator_MissingDocument(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	v := NewValidator()
	findings, err := v.Validate(fs, "docs/Architecture.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) == 0 {
		t.Error("expected findings for missing document")
	}
	if findings[0].Severity != SeverityError {
		t.Errorf("expected error severity for missing document, got %s", findings[0].Severity)
	}
}

func TestValidator_MissingSections(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("docs/Empty.md", []byte("Just some text with no headings.\n"), 0644)

	v := NewValidator()
	findings, _ := v.Validate(fs, "docs/Empty.md")
	found := false
	for _, f := range findings {
		if f.Rule == "missing-sections" {
			found = true
		}
	}
	if !found {
		t.Error("expected missing-sections finding for document with no headings")
	}
}

func TestValidator_StaleDocument(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("docs/Stale.md", []byte("# Architecture\n\n_Define components here._\n"), 0644)

	v := NewValidator()
	findings, _ := v.Validate(fs, "docs/Stale.md")
	found := false
	for _, f := range findings {
		if f.Rule == "stale-document" {
			found = true
		}
	}
	if !found {
		t.Error("expected stale-document finding for document with placeholder text")
	}
}

func TestValidator_DuplicateHeadings(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("docs/Duplicate.md", []byte("# Overview\n\nContent.\n\n# Overview\n\nMore content.\n"), 0644)

	v := NewValidator()
	findings, _ := v.Validate(fs, "docs/Duplicate.md")
	found := false
	for _, f := range findings {
		if f.Rule == "duplicate-content" {
			found = true
		}
	}
	if !found {
		t.Error("expected duplicate-content finding for document with repeated headings")
	}
}

func TestValidator_CleanDocument(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	content := "# Architecture\n\n## Overview\n\nThis is the architecture of the system. It uses Go and PostgreSQL.\n\n## Components\n\nThe API server, the database, and the worker.\n"
	_ = fs.WriteFile("docs/Clean.md", []byte(content), 0644)

	v := NewValidator()
	findings, _ := v.Validate(fs, "docs/Clean.md")
	for _, f := range findings {
		if f.Severity == SeverityError {
			t.Errorf("did not expect error-level finding on clean document: %s", f.Title)
		}
	}
}

func TestValidator_MarkdownLinksOnlyAndRelativePaths(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("docs/Guide.md", []byte("# Guide\n\nThis prose has [brackets] and (parentheses) but is not a link.\n\nSee [Architecture](Architecture.md).\n"), 0644)
	_ = fs.WriteFile("docs/Architecture.md", []byte("# Architecture\n"), 0644)

	v := NewValidator()
	findings, err := v.Validate(fs, "docs/Guide.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, finding := range findings {
		if finding.Rule == "broken-references" {
			t.Fatalf("did not expect broken reference false positive: %#v", finding)
		}
	}
}

func TestValidator_BrokenMarkdownLink(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("docs/Guide.md", []byte("# Guide\n\nSee [Missing](Missing.md).\n"), 0644)
	v := NewValidator()
	findings, err := v.Validate(fs, "docs/Guide.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, finding := range findings {
		if finding.Rule == "broken-references" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected broken markdown link finding")
	}
}

func TestValidator_LocalFileURILinks(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("docs/Guide.md", []byte("# Guide\n\nSee [Source](file:///repo/internal/app/root.go).\n"), 0644)
	_ = fs.WriteFile("/repo/internal/app/root.go", []byte("package app"), 0644)
	v := NewValidator()
	findings, err := v.Validate(fs, "docs/Guide.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, finding := range findings {
		if finding.Rule == "broken-references" {
			t.Fatalf("did not expect local file URI to be broken: %#v", finding)
		}
	}
}

// customRule simulates a plugin-provided validation rule
type customRule struct{}

func (r *customRule) Name() string { return "custom-rule" }
func (r *customRule) Run(_ filesystem.FileSystem, path string, content string) []ValidationFinding {
	return []ValidationFinding{{Rule: r.Name(), FilePath: path, Severity: SeverityInfo, Title: "custom check passed"}}
}

func TestValidator_PluginRule(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("docs/Any.md", []byte("# Title\n\nBody.\n"), 0644)

	v := NewValidator()
	v.Register(&customRule{})
	findings, _ := v.Validate(fs, "docs/Any.md")
	found := false
	for _, f := range findings {
		if f.Rule == "custom-rule" {
			found = true
		}
	}
	if !found {
		t.Error("expected custom plugin rule to fire")
	}
}
