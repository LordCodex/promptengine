package validation

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestRegistry_MissingManifest_Critical(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	reg := NewRegistry()

	findings, err := reg.Run(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range findings {
		if f.Rule == "project-config" && f.Severity == quality.SeverityCritical {
			found = true
		}
	}
	if !found {
		t.Error("expected critical finding for missing manifest")
	}
}

func TestRegistry_EmptyManifest_Error(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("playbook-manifest.json", []byte(""), 0644)

	reg := NewRegistry()
	findings, err := reg.Run(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range findings {
		if f.Rule == "manifest-schema" && f.Severity == quality.SeverityError {
			found = true
		}
	}
	if !found {
		t.Error("expected error finding for empty manifest")
	}
}

func TestRegistry_MissingCoreDocs_Warnings(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("playbook-manifest.json", []byte(`{}`), 0644)
	_ = fs.WriteFile(".promptengine/.keep", []byte{}, 0644)

	reg := NewRegistry()
	findings, err := reg.Run(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	warningCount := 0
	for _, f := range findings {
		if f.Rule == "documentation-completeness" && f.Severity == quality.SeverityWarning {
			warningCount++
		}
	}
	if warningCount == 0 {
		t.Error("expected warning findings for missing core docs")
	}
}

func TestRegistry_GitMissing_Info(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	reg := NewRegistry()

	findings, err := reg.Run(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range findings {
		if f.Rule == "hooks-configured" && f.Severity == quality.SeverityInfo {
			found = true
		}
	}
	if !found {
		t.Error("expected info finding for missing .git directory")
	}
}

// pluginValidator simulates a third-party plugin validator
type pluginValidator struct{}

func (v *pluginValidator) ID() string       { return "custom-validator" }
func (v *pluginValidator) Category() string { return "custom" }
func (v *pluginValidator) Validate(_ filesystem.FileSystem) ([]quality.Finding, error) {
	return []quality.Finding{{
		Engine:   "validation",
		Rule:     v.ID(),
		Category: v.Category(),
		Severity: quality.SeverityInfo,
		Title:    "Plugin validator ran",
	}}, nil
}

func TestRegistry_PluginValidator(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	reg := NewRegistry()
	reg.Register(&pluginValidator{})

	findings, err := reg.Run(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range findings {
		if f.Rule == "custom-validator" {
			found = true
		}
	}
	if !found {
		t.Error("expected plugin validator finding")
	}
}
