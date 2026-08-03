package manifest

import (
	"testing"
)

func TestManifestEngine_MergeAndQuery(t *testing.T) {
	engine := NewEngine()

	// 1. Setup Core Manifest
	core := &DeclarativeManifest{
		SchemaVersion: "1.0",
		Workflows: map[string]WorkflowDef{
			"bug-fix": {ID: "bug-fix", Name: "Core Bug Fix", Description: "Core fixes"},
		},
		Standards: map[string]StandardDef{
			"standard-lint": {ID: "standard-lint", Title: "Linting Rules", Priority: 50},
		},
	}
	engine.RegisterMemoryManifest("core", core)

	// 2. Setup Project Manifest overriding Core settings
	project := &DeclarativeManifest{
		Workflows: map[string]WorkflowDef{
			"bug-fix": {ID: "bug-fix", Name: "Project Bug Fix Override", Description: "Project customized workflow rules"},
		},
	}
	engine.RegisterMemoryManifest("project", project)

	// Verify merged state
	merged := engine.GetMergedManifest()

	// Expect the Project override to take precedence
	w := merged.Workflows["bug-fix"]
	if w.Name != "Project Bug Fix Override" {
		t.Errorf("Expected project workflow override name 'Project Bug Fix Override', got '%s'", w.Name)
	}

	// Expect non-overridden core standard to remain intact
	std := merged.Standards["standard-lint"]
	if std.Priority != 50 {
		t.Errorf("Expected standard priority 50, got %d", std.Priority)
	}

	// Test 3: Querying API
	query := NewQueryEngine(engine)
	retrievedW, err := query.FindWorkflow("bug-fix")
	if err != nil {
		t.Fatalf("Failed to query workflow: %v", err)
	}
	if retrievedW.Name != "Project Bug Fix Override" {
		t.Errorf("Expected query to return merged project workflow, got '%s'", retrievedW.Name)
	}
}

func TestManifestEngine_CompatibilityCheck(t *testing.T) {
	engine := NewEngine()

	core := &DeclarativeManifest{
		Compatibility: VersionCompatibilityDef{
			MinCLIVersion: "0.2.0",
		},
	}
	engine.RegisterMemoryManifest("core", core)

	query := NewQueryEngine(engine)

	// Test compatible CLI
	report := query.VerifyCompatibility("0.3.0")
	if !report.IsCompatible {
		t.Errorf("Expected CLI version 0.3.0 to be compatible, got incompatible: %s", report.Reason)
	}

	// Test incompatible CLI
	report = query.VerifyCompatibility("0.1.0")
	if report.IsCompatible {
		t.Errorf("Expected CLI version 0.1.0 to be incompatible, got compatible")
	}
}
