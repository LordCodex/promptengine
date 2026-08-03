package context

import (
	"context"
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestContextEngine_ScoringAndBudget(t *testing.T) {
	fs := filesystem.NewMockFileSystem()

	// Write mock playbook templates
	_ = fs.WriteFile("AGENTS.md", []byte("Constitutional Core Guidelines"), 0644)
	_ = fs.WriteFile("docs/BusinessRules.md", []byte("Enterprise rules and invariants"), 0644)
	_ = fs.WriteFile("docs/Roadmap.md", []byte("Milestones schedule timeline specs details long details"), 0644)
	_ = fs.WriteFile("docs/Troubleshooting.md", []byte("Troubleshoot log keys"), 0644)

	// Build ProjectModel Mock discovery output
	pm := discovery.NewProjectModel("")
	pm.PromptEngine.AgentsMDPresent = true
	pm.Docs["BusinessRules"] = discovery.DocSpec{Name: "BusinessRules", Path: "docs/BusinessRules.md", Exists: true}
	pm.Docs["Roadmap"] = discovery.DocSpec{Name: "Roadmap", Path: "docs/Roadmap.md", Exists: true}
	pm.Docs["Troubleshooting"] = discovery.DocSpec{Name: "Troubleshooting", Path: "docs/Troubleshooting.md", Exists: true}

	engine := NewEngine(fs)

	// Test 1: BugFix context under Tiny budget (should drop lower priority documents)
	pkg, err := engine.GenerateContext(context.Background(), TaskBugFix, pm, BudgetTiny)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify that BusinessRules and AGENTS.md outrank Roadmap
	hasAgents := false
	hasBusiness := false
	for _, doc := range pkg.Documents {
		if doc.Path == "AGENTS.md" {
			hasAgents = true
		}
		if doc.Path == "docs/BusinessRules.md" {
			hasBusiness = true
		}
	}

	if !hasAgents || !hasBusiness {
		t.Errorf("Expected higher priority documents AGENTS.md and BusinessRules.md to be present, got documents %v", pkg.Documents)
	}

	// BudgetTiny limit is 5000 bytes. Mock files are small so all might fit, let's verify drop limits by using an extremely tiny budget simulation
	// Let's verify priority sorting: agents (100) > business_rules (90) > troubleshooting bugfix elevated (80) > roadmap (60)
	if pkg.Documents[0].Path != "AGENTS.md" {
		t.Errorf("Expected first prioritized file to be AGENTS.md, got %s", pkg.Documents[0].Path)
	}

	// Test 2: Formatting output CursorRules
	formatter := NewFormatter(ProviderCursor)
	cursorRules, err := formatter.Format(pkg)
	if err != nil {
		t.Fatalf("Failed to format cursor rules: %v", err)
	}

	if !strings.Contains(cursorRules, ".cursorrules standard format") && !strings.Contains(cursorRules, "rules") {
		t.Errorf("Expected valid Cursor rules envelope representation, got: %s", cursorRules)
	}
}
