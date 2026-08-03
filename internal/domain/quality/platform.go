package quality

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
	"gopkg.in/yaml.v3"
)

type Rule interface {
	ID() string
	Engine() string
	Category() string
	Run(ctx context.Context, fs filesystem.FileSystem) ([]Finding, error)
}

type RuleFunc struct {
	RuleID       string
	RuleEngine   string
	RuleCategory string
	Fn           func(context.Context, filesystem.FileSystem) ([]Finding, error)
}

func (r RuleFunc) ID() string       { return r.RuleID }
func (r RuleFunc) Engine() string   { return r.RuleEngine }
func (r RuleFunc) Category() string { return r.RuleCategory }
func (r RuleFunc) Run(ctx context.Context, fs filesystem.FileSystem) ([]Finding, error) {
	return r.Fn(ctx, fs)
}

type RuleRegistry struct {
	rules   map[string]Rule
	sources map[string]string
}

func NewRuleRegistry() *RuleRegistry {
	return &RuleRegistry{rules: map[string]Rule{}, sources: map[string]string{}}
}

func (r *RuleRegistry) Register(source string, rule Rule) {
	r.rules[rule.ID()] = rule
	r.sources[rule.ID()] = source
}

func (r *RuleRegistry) All() []Rule {
	out := make([]Rule, 0, len(r.rules))
	for _, rule := range r.rules {
		out = append(out, rule)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

func (r *RuleRegistry) ApplyManifest(m *manifest.Manifest) {
	if m == nil {
		return
	}
	for _, playbook := range m.Playbooks {
		if strings.Contains(strings.ToLower(string(playbook.Category)), "security") {
			id := "manifest-" + playbook.ID
			location := playbook.Location
			r.Register("manifest", RuleFunc{
				RuleID: id, RuleEngine: "validation", RuleCategory: "security",
				Fn: func(ctx context.Context, fs filesystem.FileSystem) ([]Finding, error) {
					if location != "" && !fs.Exists(location) {
						return []Finding{NewFinding("validation", id, "security", SeverityWarning, "Manifest standard is missing", "Restore the standard referenced by the manifest.").WithFile(location)}, nil
					}
					return nil, nil
				},
			})
		}
	}
}

type Platform struct {
	fs     filesystem.FileSystem
	events *eventbus.EventBus
	rules  *RuleRegistry
}

func NewPlatform(fs filesystem.FileSystem, events *eventbus.EventBus) *Platform {
	p := &Platform{fs: fs, events: events, rules: NewRuleRegistry()}
	p.RegisterBuiltIns()
	return p
}

func (p *Platform) Rules() *RuleRegistry { return p.rules }

func (p *Platform) RegisterBuiltIns() {
	p.rules.Register("built-in", RuleFunc{RuleID: "required-structure", RuleEngine: "validation", RuleCategory: "architecture", Fn: func(ctx context.Context, fs filesystem.FileSystem) ([]Finding, error) {
		var findings []Finding
		for _, dir := range []string{"docs"} {
			if !fs.Exists(dir) {
				findings = append(findings, NewFinding("validation", "required-structure", "architecture", SeverityWarning, "Required directory is missing", "Create the missing project directory.").WithFile(dir))
			}
		}
		return findings, nil
	}})
	p.rules.Register("built-in", RuleFunc{RuleID: "security-doc", RuleEngine: "review", RuleCategory: "security", Fn: func(ctx context.Context, fs filesystem.FileSystem) ([]Finding, error) {
		if !fs.Exists("docs/Security.md") && !fs.Exists("SECURITY.md") {
			return []Finding{NewFinding("review", "security-doc", "security", SeverityError, "Security documentation is missing", "Generate docs/Security.md and document baseline controls.").WithFile("docs/Security.md")}, nil
		}
		return nil, nil
	}})
	p.rules.Register("built-in", RuleFunc{RuleID: "testing-doc", RuleEngine: "review", RuleCategory: "testing", Fn: func(ctx context.Context, fs filesystem.FileSystem) ([]Finding, error) {
		if !fs.Exists("docs/Testing.md") {
			return []Finding{NewFinding("review", "testing-doc", "testing", SeverityWarning, "Testing documentation is missing", "Generate docs/Testing.md and define testing practices.").WithFile("docs/Testing.md")}, nil
		}
		return nil, nil
	}})
}

func (p *Platform) Validate(ctx context.Context) (*Report, error) {
	p.publish(eventbus.ValidationStarted, "validation started", nil)
	report, err := p.run(ctx, "validation", "Validation Report")
	if err != nil {
		return nil, err
	}
	p.publish(eventbus.ValidationCompleted, "validation completed", report)
	return report, nil
}

func (p *Platform) Review(ctx context.Context) (*Report, error) {
	p.publish(eventbus.ReviewStarted, "review started", nil)
	report, err := p.run(ctx, "review", "Review Report")
	if err != nil {
		return nil, err
	}
	p.publish(eventbus.ReviewCompleted, "review completed", report)
	return report, nil
}

func (p *Platform) Audit(ctx context.Context) (*Report, error) {
	var all []Finding
	for _, engine := range []string{"validation", "review", "doctor", "compliance"} {
		report, err := p.run(ctx, engine, engine+" report")
		if err != nil {
			return nil, err
		}
		all = append(all, report.Findings...)
	}
	report := NewReport("Quality Audit Report", all, map[string]string{"engine": "audit"})
	p.publish(eventbus.AuditCompleted, "audit completed", report)
	return report, nil
}

func (p *Platform) Health(ctx context.Context) (*Report, error) {
	report, err := p.Audit(ctx)
	if err != nil {
		return nil, err
	}
	report.Title = "Quality Health Report"
	p.publish(eventbus.HealthCalculated, "health calculated", report)
	return report, nil
}

func (p *Platform) run(ctx context.Context, engineName, title string) (*Report, error) {
	var findings []Finding
	for _, rule := range p.rules.All() {
		if rule.Engine() != engineName {
			continue
		}
		res, err := rule.Run(ctx, p.fs)
		if err != nil {
			return nil, fmt.Errorf("quality rule %s failed: %w", rule.ID(), err)
		}
		for _, finding := range NormalizeFindings(res) {
			findings = append(findings, finding)
			if finding.Severity == SeverityCritical {
				p.publish(eventbus.CriticalIssueDetected, "critical issue detected", finding)
			}
		}
	}
	return NewReport(title, findings, map[string]string{"engine": engineName}), nil
}

func NewFinding(engine, rule, category string, severity Severity, title, recommendation string) Finding {
	f := Finding{Engine: engine, Rule: rule, Category: category, Severity: severity, Title: title, Recommendation: recommendation}
	f.ID = FindingID(f)
	return f
}

func (f Finding) WithFile(path string) Finding {
	f.FilePath = path
	f.ID = FindingID(f)
	return f
}

func FindingID(f Finding) string {
	parts := []string{f.Engine, f.Rule, f.Category, f.FilePath, f.Title}
	return strings.ToLower(strings.ReplaceAll(strings.Join(parts, ":"), " ", "-"))
}

func NormalizeFindings(findings []Finding) []Finding {
	out := make([]Finding, len(findings))
	for i, f := range findings {
		if f.ID == "" {
			f.ID = FindingID(f)
		}
		if f.StartLine == 0 && f.Line > 0 {
			f.StartLine = f.Line
		}
		if f.EndLine == 0 && f.Line > 0 {
			f.EndLine = f.Line
		}
		out[i] = f
	}
	return out
}

func NewReport(title string, findings []Finding, meta map[string]string) *Report {
	findings = NormalizeFindings(findings)
	return &Report{Title: title, Findings: findings, Score: ScoreFindings(findings), Meta: meta}
}

func ScoreFindings(findings []Finding) Score {
	cats := []string{"security", "documentation", "architecture", "testing", "performance", "maintainability"}
	weights := map[string]float64{"security": 0.2, "documentation": 0.2, "architecture": 0.2, "testing": 0.15, "performance": 0.1, "maintainability": 0.15}
	var scores []CategoryScore
	for _, cat := range cats {
		raw := 100
		critical := false
		count := 0
		for _, f := range findings {
			if f.Category != cat {
				continue
			}
			count++
			switch f.Severity {
			case SeverityCritical:
				raw -= 60
				critical = true
			case SeverityError:
				raw -= 25
			case SeverityWarning:
				raw -= 10
			default:
				raw -= 3
			}
		}
		if raw < 0 {
			raw = 0
		}
		scores = append(scores, CategoryScore{Category: cat, Weight: weights[cat], Raw: raw, Weighted: float64(raw) * weights[cat], CriticalFail: critical, FindingCount: count})
	}
	return ComputeScore(scores, 70)
}

func (r *Report) ToYAML() ([]byte, error) {
	return yaml.Marshal(r)
}

func (r *Report) Human() string {
	var b strings.Builder
	b.WriteString("Quality Report\n\n")
	counts := map[string]map[Severity]int{}
	for _, f := range r.Findings {
		if counts[f.Category] == nil {
			counts[f.Category] = map[Severity]int{}
		}
		counts[f.Category][f.Severity]++
	}
	var cats []string
	for cat := range counts {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	for _, cat := range cats {
		b.WriteString(strings.Title(cat))
		b.WriteString(":\n")
		for _, sev := range []Severity{SeverityCritical, SeverityError, SeverityWarning, SeverityInfo} {
			if counts[cat][sev] > 0 {
				b.WriteString(fmt.Sprintf("  %s: %d\n", strings.Title(string(sev)), counts[cat][sev]))
			}
		}
		b.WriteString("\n")
	}
	b.WriteString(fmt.Sprintf("Health Score:\n  %d/100\n", r.Score.Overall))
	return b.String()
}

func (p *Platform) publish(t eventbus.EventType, msg string, payload interface{}) {
	if p.events != nil {
		p.events.Publish(eventbus.Event{Type: t, Message: msg, Payload: payload})
	}
}
