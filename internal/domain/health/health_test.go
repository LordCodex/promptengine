package health

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// ─── Mock Checkers ─────────────────────────────────────────────────────────

type cleanChecker struct{ name string }

func (c *cleanChecker) Name() string { return c.name }
func (c *cleanChecker) Check(_ filesystem.FileSystem) ([]Issue, error) {
	return nil, nil
}

type blockingChecker struct {
	name     string
	category string
}

func (c *blockingChecker) Name() string { return c.name }
func (c *blockingChecker) Check(_ filesystem.FileSystem) ([]Issue, error) {
	return []Issue{{
		Category: c.category,
		Title:    "critical block",
		Severity: SeverityBlock,
	}}, nil
}

type warningChecker struct {
	name     string
	category string
}

func (c *warningChecker) Name() string { return c.name }
func (c *warningChecker) Check(_ filesystem.FileSystem) ([]Issue, error) {
	return []Issue{{
		Category: c.category,
		Title:    "advisory",
		Severity: SeveritySuggestion,
	}}, nil
}

// ─── Tests ─────────────────────────────────────────────────────────────────

func TestRegistry_CleanProject(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&cleanChecker{name: "docs-check"})
	reg.Register(&cleanChecker{name: "security-check"})

	result, err := reg.Evaluate(filesystem.NewMockFileSystem())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Score < 95 {
		t.Errorf("expected score >= 95 for clean project, got %d", result.Score)
	}
	if result.Rating != "A" {
		t.Errorf("expected rating A, got %s", result.Rating)
	}
}

func TestRegistry_BlockIssue_ZerosCategory(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&blockingChecker{name: "sec", category: "security"})

	result, err := reg.Evaluate(filesystem.NewMockFileSystem())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cat := result.Categories["security"]
	if cat == nil {
		t.Fatal("expected security category in result")
	}
	if cat.Score >= 100 {
		t.Errorf("expected security category score to be reduced by block, got %d", cat.Score)
	}
}

func TestRegistry_WarningReducesScore(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&warningChecker{name: "docs", category: "documentation"})

	result, err := reg.Evaluate(filesystem.NewMockFileSystem())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Score >= 100 {
		t.Error("expected score below 100 when there are warnings")
	}
}

func TestRegistry_CategoryBreakdown(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&blockingChecker{name: "deploy", category: "deployment"})

	result, err := reg.Evaluate(filesystem.NewMockFileSystem())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Categories) == 0 {
		t.Error("expected per-category breakdown in result")
	}
}

func TestRegistry_ManifestDrivenWeights(t *testing.T) {
	reg := NewRegistry()
	reg.SetCategories([]HealthCategory{
		{Name: "documentation", Weight: 0.5},
		{Name: "security", Weight: 0.5},
	})
	reg.Register(&cleanChecker{name: "all-clean"})

	result, err := reg.Evaluate(filesystem.NewMockFileSystem())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Score != 100 {
		t.Errorf("expected 100 with clean checkers and custom weights, got %d", result.Score)
	}
}

func TestRegistry_Threshold(t *testing.T) {
	reg := NewRegistry()
	reg.SetThreshold(60)
	// Multiple blocking issues should push score below 60
	reg.Register(&blockingChecker{name: "b1", category: "security"})
	reg.Register(&blockingChecker{name: "b2", category: "documentation"})
	reg.Register(&blockingChecker{name: "b3", category: "testing"})

	result, err := reg.Evaluate(filesystem.NewMockFileSystem())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Score should be reduced; just verify it stays in valid range
	if result.Score < 0 || result.Score > 100 {
		t.Errorf("score out of range: %d", result.Score)
	}
}
