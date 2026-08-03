package report

import (
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/quality"
)

func makeReport(score int, passed bool) *quality.Report {
	return &quality.Report{
		Title: "Test Report",
		Score: quality.Score{
			Overall: score,
			Rating:  "B",
			Passed:  passed,
		},
		Findings: []quality.Finding{
			{
				Engine:         "test",
				Rule:           "test-rule",
				Category:       "documentation",
				Severity:       quality.SeverityWarning,
				Title:          "Missing document",
				Recommendation: "Create the document.",
			},
		},
	}
}

func TestTextRenderer(t *testing.T) {
	reg := NewRendererRegistry()
	data, err := reg.Render("text", makeReport(85, true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := string(data)
	if !strings.Contains(out, "85") {
		t.Error("expected score in text output")
	}
	if !strings.Contains(out, "PASSED") {
		t.Error("expected PASSED label in text output")
	}
}

func TestJSONRenderer(t *testing.T) {
	reg := NewRendererRegistry()
	data, err := reg.Render("json", makeReport(85, true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(data), `"Title"`) {
		t.Error("expected JSON to contain Title field")
	}
}

func TestYAMLRenderer(t *testing.T) {
	reg := NewRendererRegistry()
	data, err := reg.Render("yaml", makeReport(85, true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty YAML output")
	}
}

func TestMarkdownRenderer(t *testing.T) {
	reg := NewRendererRegistry()
	data, err := reg.Render("markdown", makeReport(72, false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := string(data)
	if !strings.HasPrefix(out, "# ") {
		t.Error("expected markdown to start with heading")
	}
	if !strings.Contains(out, "FAILED") {
		t.Error("expected FAILED label in markdown output")
	}
}

func TestSARIFRenderer(t *testing.T) {
	reg := NewRendererRegistry()
	data, err := reg.Render("sarif", makeReport(85, true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(data), "PromptEngine") {
		t.Error("expected SARIF to reference PromptEngine tool name")
	}
}

func TestUnknownFormat(t *testing.T) {
	reg := NewRendererRegistry()
	if _, err := reg.Render("html", makeReport(85, true)); err == nil {
		t.Error("expected error for unknown report format")
	}
}

func TestCIThreshold_Pass(t *testing.T) {
	report := makeReport(85, true)
	result := EvaluateCIThreshold(report, quality.DefaultThreshold)
	if !result.Passed {
		t.Errorf("expected CI pass for score 85, got: %s", result.Message)
	}
}

func TestCIThreshold_Fail_LowScore(t *testing.T) {
	report := &quality.Report{
		Title: "Low Score",
		Score: quality.Score{Overall: 50, Passed: false},
	}
	result := EvaluateCIThreshold(report, quality.DefaultThreshold)
	if result.Passed {
		t.Error("expected CI fail for score below threshold")
	}
}
