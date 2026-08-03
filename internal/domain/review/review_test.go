package review

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestRegistry_ReviewFindings(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	reg := NewRegistry()
	RegisterDefaultReviewers(reg)

	// No manifest, no docs — should produce findings
	session, err := reg.RunSession(fs, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(session.Findings) == 0 {
		t.Error("expected findings for empty project")
	}
}

func TestRegistry_CleanProject_NoBlockFindings(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("playbook-manifest.json", []byte(`{}`), 0644)
	_ = fs.WriteFile("docs/Security.md", []byte("# Security\n"), 0644)
	_ = fs.WriteFile("docs/Testing.md", []byte("# Testing\n"), 0644)
	_ = fs.WriteFile("docs/.keep", []byte{}, 0644)

	reg := NewRegistry()
	RegisterDefaultReviewers(reg)

	session, err := reg.RunSession(fs, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, f := range session.Findings {
		if f.Severity == "block" {
			t.Errorf("unexpected block finding on clean project: %s", f.Message)
		}
	}
}

func TestRegistry_SessionSummary(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	reg := NewRegistry()
	RegisterDefaultReviewers(reg)

	session, err := reg.RunSession(fs, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Summary map should have at least one entry
	if len(session.Summary) == 0 {
		t.Error("expected session summary to contain review type counts")
	}
}

// pluginReviewer simulates a plugin-contributed reviewer
type pluginReviewer struct{ fired bool }

type alwaysFindingRule struct{ parent *pluginReviewer }

func (r *alwaysFindingRule) Name() string { return "plugin-rule" }
func (r *alwaysFindingRule) Evaluate(_ filesystem.FileSystem, _ string) ([]Finding, error) {
	r.parent.fired = true
	return []Finding{{Category: "custom", Severity: "suggestion", Message: "plugin finding"}}, nil
}

func (p *pluginReviewer) Type() ReviewType    { return ReviewType("custom") }
func (p *pluginReviewer) Description() string { return "Plugin reviewer" }
func (p *pluginReviewer) Rules() []Rule       { return []Rule{&alwaysFindingRule{parent: p}} }

func TestRegistry_PluginReviewer(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	reg := NewRegistry()
	pr := &pluginReviewer{}
	reg.RegisterReviewer(pr)

	session, err := reg.RunSession(fs, ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !pr.fired {
		t.Error("expected plugin reviewer rules to run")
	}
	if len(session.Findings) == 0 {
		t.Error("expected plugin finding in session")
	}
}

func TestSession_CountBySeverity(t *testing.T) {
	session := &ReviewSession{Summary: make(map[ReviewType]int)}
	session.Add(Finding{Severity: "block"})
	session.Add(Finding{Severity: "block"})
	session.Add(Finding{Severity: "suggestion"})

	if session.CountBySeverity("block") != 2 {
		t.Errorf("expected 2 block findings, got %d", session.CountBySeverity("block"))
	}
	if session.CountBySeverity("suggestion") != 1 {
		t.Errorf("expected 1 suggestion finding, got %d", session.CountBySeverity("suggestion"))
	}
}
