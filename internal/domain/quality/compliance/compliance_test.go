package compliance

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestComplianceEngine_EmptyProject(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewComplianceEngine()

	report, err := engine.Run(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.ProfileResults) == 0 {
		t.Error("expected profile results from default profiles")
	}
}

func TestComplianceEngine_MissingManifest_Critical(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewComplianceEngine()

	report, _ := engine.Run(fs)
	found := false
	for _, pr := range report.ProfileResults {
		for _, f := range pr.Findings {
			if f.Rule == "manifest-required" && f.Severity == quality.SeverityCritical {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected critical finding for missing manifest")
	}
}

func TestComplianceEngine_CoreProfile_Fails_Without_Manifest(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewComplianceEngine()

	report, _ := engine.Run(fs)
	for _, pr := range report.ProfileResults {
		if pr.Profile.ID == "promptengine-core" && pr.Passed {
			t.Error("expected promptengine-core profile to fail without manifest")
		}
	}
}

func TestComplianceEngine_SecurityProfile(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile("docs/Security.md", []byte("# Security\n"), 0644)

	engine := NewComplianceEngine()
	report, _ := engine.Run(fs)
	for _, pr := range report.ProfileResults {
		if pr.Profile.ID == "security-baseline" {
			for _, f := range pr.Findings {
				if f.Rule == "security-doc-required" {
					t.Error("expected no finding when Security.md exists")
				}
			}
		}
	}
}

// customComplianceRule simulates a plugin-contributed compliance rule
type customComplianceRule struct{}

func (r *customComplianceRule) ID() string          { return "custom-compliance" }
func (r *customComplianceRule) Description() string { return "Custom compliance rule" }
func (r *customComplianceRule) Check(_ filesystem.FileSystem) ([]quality.Finding, error) {
	return []quality.Finding{{
		Engine:   "compliance",
		Rule:     r.ID(),
		Category: "organization",
		Severity: quality.SeverityInfo,
		Title:    "Custom compliance check passed",
	}}, nil
}

func TestComplianceEngine_CustomProfile(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewComplianceEngine()
	engine.RegisterProfile(ComplianceProfile{
		ID:    "org-custom",
		Name:  "Org Custom",
		Type:  ProfileOrganization,
		Rules: []ComplianceRule{&customComplianceRule{}},
	})

	report, err := engine.Run(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, pr := range report.ProfileResults {
		if pr.Profile.ID == "org-custom" {
			for _, f := range pr.Findings {
				if f.Rule == "custom-compliance" {
					found = true
				}
			}
		}
	}
	if !found {
		t.Error("expected custom compliance rule to produce a finding")
	}
}

func TestComplianceEngine_OverallScore(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	engine := NewComplianceEngine()
	report, _ := engine.Run(fs)

	if report.Overall.Overall < 0 || report.Overall.Overall > 100 {
		t.Errorf("overall score out of range: %d", report.Overall.Overall)
	}
}
