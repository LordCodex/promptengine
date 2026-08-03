package doctor

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestDoctorEngine_ManifestMissing(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewDoctorEngine()

	report, err := engine.Diagnose(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range report.Findings {
		if f.Rule == "manifest-integrity" && f.Severity == quality.SeverityCritical {
			found = true
		}
	}
	if !found {
		t.Error("expected critical finding for missing manifest")
	}
}

func TestDoctorEngine_ManifestPresent_NoFinding(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("playbook-manifest.json", []byte(`{}`), 0644)

	engine := NewDoctorEngine()
	report, err := engine.Diagnose(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, f := range report.Findings {
		if f.Rule == "manifest-integrity" {
			t.Errorf("expected no manifest finding when file exists, got: %s", f.Title)
		}
	}
}

func TestDoctorEngine_ConfigMissing(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewDoctorEngine()

	report, _ := engine.Diagnose(fs)
	found := false
	for _, f := range report.Findings {
		if f.Rule == "configuration-exists" {
			found = true
		}
	}
	if !found {
		t.Error("expected finding for missing .promptengine directory")
	}
}

func TestDoctorEngine_DocsMissing_Warning(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("playbook-manifest.json", []byte(`{}`), 0644)
	_ = fs.WriteFile(".promptengine/.keep", []byte{}, 0644)

	engine := NewDoctorEngine()
	report, _ := engine.Diagnose(fs)
	found := false
	for _, f := range report.Findings {
		if f.Rule == "docs-directory" && f.Severity == quality.SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected warning finding for missing docs/ directory")
	}
}

// pluginCheck simulates a plugin-contributed doctor check
type pluginCheck struct{ fired bool }

func (c *pluginCheck) ID() string          { return "plugin-custom-check" }
func (c *pluginCheck) Category() string    { return "integrity" }
func (c *pluginCheck) Description() string { return "Custom plugin check" }
func (c *pluginCheck) Run(_ filesystem.FileSystem) ([]quality.Finding, error) {
	c.fired = true
	return nil, nil
}

func TestDoctorEngine_PluginCheck(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewDoctorEngine()
	pc := &pluginCheck{}
	engine.Register(pc)

	_, err := engine.Diagnose(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !pc.fired {
		t.Error("expected plugin check to run during diagnosis")
	}
}

func TestDoctorEngine_Fix_DryRun(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewDoctorEngine()
	applied := false
	engine.RegisterRepair(RepairAction{
		FindingRuleID: "manifest-integrity",
		Description:   "Create manifest",
		Apply: func(fs filesystem.FileSystem) error {
			applied = true
			return fs.WriteFile("playbook-manifest.json", []byte(`{}`), 0644)
		},
	})

	if err := engine.Fix("manifest-integrity", fs, true); err != nil {
		t.Fatalf("dry-run fix should not error: %v", err)
	}
	if applied {
		t.Error("expected dry-run to NOT apply the fix")
	}
}

func TestDoctorEngine_Fix_Apply(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewDoctorEngine()
	engine.RegisterRepair(RepairAction{
		FindingRuleID: "manifest-integrity",
		Apply: func(fs filesystem.FileSystem) error {
			return fs.WriteFile("playbook-manifest.json", []byte(`{}`), 0644)
		},
	})

	if err := engine.Fix("manifest-integrity", fs, false); err != nil {
		t.Fatalf("expected fix to apply without error: %v", err)
	}
	if !fs.Exists("playbook-manifest.json") {
		t.Error("expected fix to create manifest file")
	}
}

func TestDoctorEngine_Score(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("playbook-manifest.json", []byte(`{}`), 0644)
	_ = fs.WriteFile(".promptengine/.keep", []byte{}, 0644)
	_ = fs.WriteFile("docs/.keep", []byte{}, 0644)

	engine := NewDoctorEngine()
	report, _ := engine.Diagnose(fs)
	if report.Score.Overall < 0 || report.Score.Overall > 100 {
		t.Errorf("expected score in 0-100 range, got %d", report.Score.Overall)
	}
	if report.Score.Rating == "" {
		t.Error("expected non-empty rating")
	}
}
