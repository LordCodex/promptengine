package intelligence

import (
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/domain/personal"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func testProject() *discovery.ProjectModel {
	pm := discovery.NewProjectModel(".")
	pm.Repository.Directories = []string{"app/Services", "app/Actions", "app/Http/Requests", "tests"}
	pm.Repository.Files = []string{"app/Services/PaymentService.php", "tests/PaymentServiceTest.php", "docs/API.md"}
	pm.Repository.DocumentationFiles = []string{"docs/API.md"}
	return pm
}

func TestDetectPatterns(t *testing.T) {
	p := NewPlatform(filesystem.NewMockFileSystem(), nil)
	patterns := p.DetectPatterns(testProject())
	names := map[string]bool{}
	for _, pattern := range patterns {
		names[pattern.Name] = true
	}
	if !names["Service Layer"] || !names["Action Classes"] || !names["Form Request Validation"] {
		t.Fatalf("expected Laravel-style patterns, got %#v", patterns)
	}
}

func TestDecisionStorageAndRecommendation(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	p := NewPlatform(fs, nil)
	if err := p.StoreDecision(Decision{Title: "Use UUID primary keys", Reason: "Separate public identifiers from database IDs", Affected: []string{"models", "apis"}}); err != nil {
		t.Fatalf("store failed: %v", err)
	}
	decisions, err := p.ListDecisions()
	if err != nil || len(decisions) != 1 {
		t.Fatalf("expected decision, got %#v err=%v", decisions, err)
	}
	recs := p.Recommend(nil, decisions, testProject())
	if len(recs) == 0 || !strings.Contains(recs[0].Suggestion, "Use UUID") {
		t.Fatalf("expected decision recommendation, got %#v", recs)
	}
}

func TestSimilarReferencesAndImpact(t *testing.T) {
	p := NewPlatform(filesystem.NewMockFileSystem(), nil)
	refs := p.FindSimilar("Add subscription payment support", testProject())
	if len(refs) == 0 || refs[0].Path != "app/Services/PaymentService.php" {
		t.Fatalf("expected payment reference, got %#v", refs)
	}
	impact := p.AnalyzeImpact(personal.GitContext{ChangedFiles: []string{"app/Models/User.php"}}, testProject())
	if !contains(impact.AffectedAreas, "authentication") || !contains(impact.AffectedAreas, "tests") {
		t.Fatalf("unexpected impact %#v", impact)
	}
}

func TestPromptEnhancement(t *testing.T) {
	out := FormatPromptEnhancement(Insights{Patterns: []Pattern{{Name: "Service Layer", Category: "architecture", Confidence: 0.8, Evidence: []string{"app/Services"}}}, Decisions: []Decision{{Title: "Use UUID primary keys", Reason: "Public IDs"}}, Recommendations: []Recommendation{{Suggestion: "Use services", Reason: "Existing pattern", Confidence: 0.8}}}, ImpactReport{AffectedAreas: []string{"tests"}})
	if !strings.Contains(out, "Existing Patterns") || !strings.Contains(out, "Previous Decisions") || !strings.Contains(out, "Recommendations") {
		t.Fatalf("unexpected enhancement:\n%s", out)
	}
}

func contains(items []string, expected string) bool {
	for _, item := range items {
		if item == expected {
			return true
		}
	}
	return false
}
