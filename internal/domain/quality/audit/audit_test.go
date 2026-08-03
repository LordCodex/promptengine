package audit

import (
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestAuditEngine_EmptyProject(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewAuditEngine()

	report, err := engine.Run(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Findings) == 0 {
		t.Error("expected findings for empty project")
	}
}

func TestAuditEngine_PromptEngineNotAdopted_Critical(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewAuditEngine()

	report, _ := engine.Run(fs)
	found := false
	for _, f := range report.Findings {
		if f.Rule == "promptengine-adoption" && f.Severity == quality.SeverityCritical {
			found = true
		}
	}
	if !found {
		t.Error("expected critical finding for missing PromptEngine adoption")
	}
}

func TestAuditEngine_MissingCoreDocs_Errors(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("playbook-manifest.json", []byte(`{}`), 0644)
	_ = fs.WriteFile(".promptengine/.keep", []byte{}, 0644)
	_ = fs.WriteFile("docs/.keep", []byte{}, 0644)

	engine := NewAuditEngine()
	report, _ := engine.Run(fs)

	errorCount := 0
	for _, f := range report.Findings {
		if f.Rule == "missing-documentation" && f.Severity == quality.SeverityError {
			errorCount++
		}
	}
	if errorCount == 0 {
		t.Error("expected error findings for missing core documents")
	}
}

func TestAuditEngine_CleanProject_Score(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("playbook-manifest.json", []byte(`{}`), 0644)
	_ = fs.WriteFile(".promptengine/.keep", []byte{}, 0644)
	_ = fs.WriteFile("docs/.keep", []byte{}, 0644)
	_ = fs.WriteFile("docs/Architecture.md", []byte("# Architecture\n"), 0644)
	_ = fs.WriteFile("docs/Database.md", []byte("# Database\n"), 0644)
	_ = fs.WriteFile("docs/Decisions.md", []byte("# Decisions\n"), 0644)

	engine := NewAuditEngine()
	report, _ := engine.Run(fs)

	if report.Score.Overall < 0 || report.Score.Overall > 100 {
		t.Errorf("score out of range: %d", report.Score.Overall)
	}
}

func TestAuditReport_Export_JSON(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewAuditEngine()
	report, _ := engine.Run(fs)

	data, err := report.Export("json")
	if err != nil {
		t.Fatalf("expected JSON export to succeed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty JSON export")
	}
	if !strings.Contains(string(data), "PromptEngine") {
		t.Error("expected JSON to contain report title")
	}
}

func TestAuditReport_Export_Markdown(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewAuditEngine()
	report, _ := engine.Run(fs)

	data, err := report.Export("markdown")
	if err != nil {
		t.Fatalf("expected markdown export to succeed: %v", err)
	}
	if !strings.HasPrefix(string(data), "# ") {
		t.Error("expected markdown to start with a heading")
	}
}

// customAuditRule simulates a plugin-contributed audit rule
type customAuditRule struct{}

func (r *customAuditRule) ID() string          { return "custom-audit" }
func (r *customAuditRule) Area() AuditArea     { return AuditArea("custom") }
func (r *customAuditRule) Description() string { return "Custom plugin audit rule" }
func (r *customAuditRule) Run(_ filesystem.FileSystem) ([]quality.Finding, error) {
	return []quality.Finding{{
		Engine:   "audit",
		Rule:     r.ID(),
		Severity: quality.SeverityInfo,
		Title:    "Plugin audit ran",
	}}, nil
}

func TestAuditEngine_PluginRule(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewAuditEngine()
	engine.Register(&customAuditRule{})

	report, err := engine.Run(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range report.Findings {
		if f.Rule == "custom-audit" {
			found = true
		}
	}
	if !found {
		t.Error("expected custom plugin audit rule to produce a finding")
	}
}
