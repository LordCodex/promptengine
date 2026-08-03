package quality

import (
	"context"
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

func TestFindingCreationAndNormalization(t *testing.T) {
	f := NewFinding("validation", "config", "configuration", SeverityError, "Missing config", "Create config").WithFile(".promptengine")
	f.Line = 12
	normalized := NormalizeFindings([]Finding{f})[0]
	if normalized.ID == "" {
		t.Fatal("expected finding id")
	}
	if normalized.StartLine != 12 || normalized.EndLine != 12 {
		t.Fatalf("expected line range from line, got %#v", normalized)
	}
}

func TestRuleExecutionAndPluginRules(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	platform := NewPlatform(fs, nil)
	platform.Rules().Register("plugin:test", RuleFunc{
		RuleID: "plugin-rule", RuleEngine: "validation", RuleCategory: "maintainability",
		Fn: func(ctx context.Context, fs filesystem.FileSystem) ([]Finding, error) {
			return []Finding{NewFinding("validation", "plugin-rule", "maintainability", SeverityInfo, "Plugin rule ran", "No action")}, nil
		},
	})
	report, err := platform.Validate(context.Background())
	if err != nil {
		t.Fatalf("expected validate to run, got %v", err)
	}
	if !hasFinding(report.Findings, "plugin-rule") {
		t.Fatalf("expected plugin rule finding, got %#v", report.Findings)
	}
}

func TestHealthScoringSeverityCalculation(t *testing.T) {
	findings := []Finding{
		NewFinding("review", "security-doc", "security", SeverityCritical, "Critical", "Fix"),
		NewFinding("review", "testing-doc", "testing", SeverityWarning, "Warning", "Fix"),
	}
	score := ScoreFindings(findings)
	if score.Overall != 0 {
		t.Fatalf("critical issue should force overall score to 0, got %d", score.Overall)
	}
	if score.Rating != "F" {
		t.Fatalf("expected F rating, got %s", score.Rating)
	}
}

func TestValidationFailuresAndAuditGeneration(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	platform := NewPlatform(fs, nil)
	validation, err := platform.Validate(context.Background())
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}
	if len(validation.Findings) == 0 {
		t.Fatal("expected validation findings")
	}
	audit, err := platform.Audit(context.Background())
	if err != nil {
		t.Fatalf("audit failed: %v", err)
	}
	if audit.Title != "Quality Audit Report" || len(audit.Findings) == 0 {
		t.Fatalf("expected audit report with findings, got %#v", audit)
	}
}

func TestReviewAcceptsRootSecurityDocument(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("SECURITY.md", []byte("# Security\n"), 0644)
	platform := NewPlatform(fs, nil)
	report, err := platform.Review(context.Background())
	if err != nil {
		t.Fatalf("review failed: %v", err)
	}
	if hasFinding(report.Findings, "security-doc") {
		t.Fatalf("did not expect security-doc finding with root SECURITY.md, got %#v", report.Findings)
	}
}

func TestManifestControlsQualityRules(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	platform := NewPlatform(fs, nil)
	platform.Rules().ApplyManifest(&manifest.Manifest{
		Playbooks: []manifest.PlaybookDefinition{{ID: "security-baseline", Category: manifest.CategorySecurity, Location: "security/baseline.md"}},
	})
	report, err := platform.Validate(context.Background())
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if !hasFinding(report.Findings, "manifest-security-baseline") {
		t.Fatalf("expected manifest-controlled rule finding, got %#v", report.Findings)
	}
}

func TestQualityEventsPublished(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	events := eventbus.NewEventBus()
	var seen []eventbus.EventType
	for _, eventType := range []eventbus.EventType{eventbus.ValidationStarted, eventbus.ValidationCompleted, eventbus.ReviewStarted, eventbus.ReviewCompleted, eventbus.AuditCompleted, eventbus.HealthCalculated} {
		tp := eventType
		events.Subscribe(tp, func(e eventbus.Event) { seen = append(seen, e.Type) })
	}
	platform := NewPlatform(fs, events)
	_, _ = platform.Validate(context.Background())
	_, _ = platform.Review(context.Background())
	_, _ = platform.Audit(context.Background())
	_, _ = platform.Health(context.Background())
	for _, eventType := range []eventbus.EventType{eventbus.ValidationStarted, eventbus.ValidationCompleted, eventbus.ReviewStarted, eventbus.ReviewCompleted, eventbus.AuditCompleted, eventbus.HealthCalculated} {
		if !hasEvent(seen, eventType) {
			t.Fatalf("expected event %s, got %#v", eventType, seen)
		}
	}
}

func TestQualityReportOutputs(t *testing.T) {
	report := NewReport("Quality Report", []Finding{NewFinding("review", "security", "security", SeverityWarning, "Security warning", "Fix it")}, nil)
	jsonData, err := report.ToJSON()
	if err != nil || !strings.Contains(string(jsonData), "Security warning") {
		t.Fatalf("expected JSON report, got %s err=%v", string(jsonData), err)
	}
	yamlData, err := report.ToYAML()
	if err != nil || !strings.Contains(string(yamlData), "Security warning") {
		t.Fatalf("expected YAML report, got %s err=%v", string(yamlData), err)
	}
	human := report.Human()
	if !strings.Contains(human, "Quality Report") || !strings.Contains(human, "Health Score") {
		t.Fatalf("expected human report, got %s", human)
	}
}

func hasFinding(findings []Finding, rule string) bool {
	for _, finding := range findings {
		if finding.Rule == rule {
			return true
		}
	}
	return false
}

func hasEvent(events []eventbus.EventType, expected eventbus.EventType) bool {
	for _, eventType := range events {
		if eventType == expected {
			return true
		}
	}
	return false
}
