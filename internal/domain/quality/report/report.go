package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"gopkg.in/yaml.v3"
)

// ReportRenderer is the interface all output formatters must implement
type ReportRenderer interface {
	Format() string
	Render(report *quality.Report) ([]byte, error)
}

// ─── Text Renderer ─────────────────────────────────────────────────────────

type TextRenderer struct{}

func (r *TextRenderer) Format() string { return "text" }
func (r *TextRenderer) Render(report *quality.Report) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== %s ===\n\n", report.Title))
	sb.WriteString(fmt.Sprintf("Overall Score : %d / 100 (%s)\n", report.Score.Overall, report.Score.Rating))
	sb.WriteString(fmt.Sprintf("Result        : %s\n\n", passFailLabel(report.Score.Passed)))

	if len(report.Findings) == 0 {
		sb.WriteString("No issues found. 🎉\n")
		return []byte(sb.String()), nil
	}

	sb.WriteString(fmt.Sprintf("Findings (%d):\n", len(report.Findings)))
	for _, f := range report.Findings {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n", strings.ToUpper(string(f.Severity)), f.Title))
		if f.Recommendation != "" {
			sb.WriteString(fmt.Sprintf("    → %s\n", f.Recommendation))
		}
	}
	return []byte(sb.String()), nil
}

// ─── JSON Renderer ─────────────────────────────────────────────────────────

type JSONRenderer struct{}

func (r *JSONRenderer) Format() string { return "json" }
func (r *JSONRenderer) Render(report *quality.Report) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

// ─── YAML Renderer ─────────────────────────────────────────────────────────

type YAMLRenderer struct{}

func (r *YAMLRenderer) Format() string { return "yaml" }
func (r *YAMLRenderer) Render(report *quality.Report) ([]byte, error) {
	return yaml.Marshal(report)
}

// ─── Markdown Renderer ─────────────────────────────────────────────────────

type MarkdownRenderer struct{}

func (r *MarkdownRenderer) Format() string { return "markdown" }
func (r *MarkdownRenderer) Render(report *quality.Report) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", report.Title))
	sb.WriteString(fmt.Sprintf("**Score**: %d/100 (%s) — %s\n\n", report.Score.Overall, report.Score.Rating, passFailLabel(report.Score.Passed)))

	// Category table
	if len(report.Score.Categories) > 0 {
		sb.WriteString("## Category Breakdown\n\n")
		sb.WriteString("| Category | Score | Weight |\n|---|---|---|\n")
		for _, cat := range report.Score.Categories {
			sb.WriteString(fmt.Sprintf("| %s | %d | %.0f%% |\n", cat.Category, cat.Raw, cat.Weight*100))
		}
		sb.WriteString("\n")
	}

	if len(report.Findings) == 0 {
		sb.WriteString("## Findings\n\nNo issues found.\n")
		return []byte(sb.String()), nil
	}

	sb.WriteString(fmt.Sprintf("## Findings (%d)\n\n", len(report.Findings)))
	for _, f := range report.Findings {
		sb.WriteString(fmt.Sprintf("### [%s] %s\n\n", strings.ToUpper(string(f.Severity)), f.Title))
		if f.Explanation != "" {
			sb.WriteString(fmt.Sprintf("**Why**: %s\n\n", f.Explanation))
		}
		if f.Impact != "" {
			sb.WriteString(fmt.Sprintf("**Impact**: %s\n\n", f.Impact))
		}
		if f.Recommendation != "" {
			sb.WriteString(fmt.Sprintf("**Action**: %s\n\n", f.Recommendation))
		}
		if f.AutoFixID != "" {
			sb.WriteString(fmt.Sprintf("**Auto-fix available**: `promptengine fix %s`\n\n", f.AutoFixID))
		}
	}
	return []byte(sb.String()), nil
}

// ─── SARIF Stub ─────────────────────────────────────────────────────────────

// SARIFRenderer produces SARIF 2.1 output (stub — structure only)
type SARIFRenderer struct{}

func (r *SARIFRenderer) Format() string { return "sarif" }
func (r *SARIFRenderer) Render(report *quality.Report) ([]byte, error) {
	type sarifResult struct {
		RuleID  string `json:"ruleId"`
		Message struct {
			Text string `json:"text"`
		} `json:"message"`
	}
	type sarifRun struct {
		Tool struct {
			Driver struct {
				Name string `json:"name"`
			} `json:"driver"`
		} `json:"tool"`
		Results []sarifResult `json:"results"`
	}
	type sarif struct {
		Version string     `json:"version"`
		Runs    []sarifRun `json:"runs"`
	}

	var results []sarifResult
	for _, f := range report.Findings {
		r := sarifResult{RuleID: f.Rule}
		r.Message.Text = f.Title
		results = append(results, r)
	}

	run := sarifRun{}
	run.Tool.Driver.Name = "PromptEngine"
	run.Results = results

	doc := sarif{Version: "2.1.0", Runs: []sarifRun{run}}
	return json.MarshalIndent(doc, "", "  ")
}

// ─── Registry ──────────────────────────────────────────────────────────────

// RendererRegistry allows plugins to register additional output formats
type RendererRegistry struct {
	renderers map[string]ReportRenderer
}

func NewRendererRegistry() *RendererRegistry {
	reg := &RendererRegistry{renderers: make(map[string]ReportRenderer)}
	reg.RegisterDefaults()
	return reg
}

func (r *RendererRegistry) Register(renderer ReportRenderer) {
	r.renderers[renderer.Format()] = renderer
}

func (r *RendererRegistry) RegisterDefaults() {
	r.Register(&TextRenderer{})
	r.Register(&JSONRenderer{})
	r.Register(&YAMLRenderer{})
	r.Register(&MarkdownRenderer{})
	r.Register(&SARIFRenderer{})
}

func (r *RendererRegistry) Render(format string, report *quality.Report) ([]byte, error) {
	renderer, ok := r.renderers[format]
	if !ok {
		return nil, fmt.Errorf("unknown report format '%s'", format)
	}
	return renderer.Render(report)
}

// CIResult determines CI pass/fail outcome
type CIResult struct {
	Passed  bool
	Message string
	Score   int
}

func EvaluateCIThreshold(report *quality.Report, threshold quality.Threshold) CIResult {
	if report.ExceedsBLock(threshold) {
		return CIResult{
			Passed:  false,
			Message: fmt.Sprintf("Quality check FAILED: score %d (threshold %d)", report.Score.Overall, threshold.MinScore),
			Score:   report.Score.Overall,
		}
	}
	return CIResult{
		Passed:  true,
		Message: fmt.Sprintf("Quality check PASSED: score %d", report.Score.Overall),
		Score:   report.Score.Overall,
	}
}

func passFailLabel(passed bool) string {
	if passed {
		return "PASSED ✓"
	}
	return "FAILED ✗"
}
