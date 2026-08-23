package audit

import (
	"fmt"
	"sync"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// AuditArea classifies what an audit rule examines
type AuditArea string

const (
	AreaProjectStructure     AuditArea = "project-structure"
	AreaEngineeringStandards AuditArea = "engineering-standards"
	AreaPromptEngineAdoption AuditArea = "promptengine-adoption"
	AreaMissingDocumentation AuditArea = "missing-documentation"
	AreaMissingDecisions     AuditArea = "missing-decisions"
	AreaDeprecatedTech       AuditArea = "deprecated-technology"
	AreaDependencyRisk       AuditArea = "dependency-risk"
	AreaArchitectureDrift    AuditArea = "architecture-drift"
	AreaConfigDrift          AuditArea = "config-drift"
	AreaDocumentationDrift   AuditArea = "documentation-drift"
)

// AuditRule defines a single audit check
type AuditRule interface {
	ID() string
	Area() AuditArea
	Description() string
	Run(fs filesystem.FileSystem) ([]quality.Finding, error)
}

// AuditReport is the full exportable output of an audit run
type AuditReport struct {
	Title    string
	Findings []quality.Finding
	Summary  map[AuditArea]int // finding count per area
	Score    quality.Score
}

func (r *AuditReport) Export(format string) ([]byte, error) {
	switch format {
	case "json":
		rep := quality.Report{
			Title:    r.Title,
			Score:    r.Score,
			Findings: r.Findings,
		}
		return rep.ToJSON()
	default:
		return []byte(r.toMarkdown()), nil
	}
}

func (r *AuditReport) toMarkdown() string {
	out := fmt.Sprintf("# %s\n\n", r.Title)
	out += fmt.Sprintf("**Overall Score**: %d (%s)\n\n", r.Score.Overall, r.Score.Rating)
	out += "## Findings\n\n"
	for _, f := range r.Findings {
		out += fmt.Sprintf("- **[%s]** %s — %s\n", f.Severity, f.Title, f.Recommendation)
	}
	return out
}

// AuditEngine orchestrates all registered audit rules
type AuditEngine struct {
	mu    sync.RWMutex
	rules []AuditRule
}

func NewAuditEngine() *AuditEngine {
	e := &AuditEngine{}
	e.RegisterDefaults()
	return e
}

func (e *AuditEngine) Register(rule AuditRule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = append(e.rules, rule)
}

// EngineName implements quality.EngineRegistrar
func (e *AuditEngine) EngineName() string { return "audit" }

// Run executes all audit rules and produces a structured report
func (e *AuditEngine) Run(fs filesystem.FileSystem) (*AuditReport, error) {
	e.mu.RLock()
	rules := e.rules
	e.mu.RUnlock()

	var findings []quality.Finding
	summary := make(map[AuditArea]int)

	for _, rule := range rules {
		f, err := rule.Run(fs)
		if err != nil {
			return nil, fmt.Errorf("audit rule '%s' failed: %w", rule.ID(), err)
		}
		findings = append(findings, f...)
		if len(f) > 0 {
			summary[rule.Area()] += len(f)
		}
	}

	score := computeAuditScore(findings)
	return &AuditReport{
		Title:    "PromptEngine Audit Report",
		Findings: findings,
		Summary:  summary,
		Score:    score,
	}, nil
}

func computeAuditScore(findings []quality.Finding) quality.Score {
	cats := []quality.CategoryScore{
		{Category: "project-structure", Weight: 0.15, Raw: 100},
		{Category: "adoption", Weight: 0.20, Raw: 100},
		{Category: "documentation", Weight: 0.25, Raw: 100},
		{Category: "architecture", Weight: 0.20, Raw: 100},
		{Category: "dependencies", Weight: 0.20, Raw: 100},
	}
	deduct := map[string]int{}
	for _, f := range findings {
		switch f.Severity {
		case quality.SeverityCritical:
			deduct[f.Category] += 40
		case quality.SeverityError:
			deduct[f.Category] += 20
		case quality.SeverityWarning:
			deduct[f.Category] += 8
		default:
			deduct[f.Category] += 2
		}
	}
	for i := range cats {
		raw := 100 - deduct[cats[i].Category]
		if raw < 0 {
			raw = 0
		}
		cats[i].Raw = raw
		cats[i].Weighted = float64(raw) * cats[i].Weight
	}
	return quality.ComputeScore(cats, 70)
}

// ─── Default Audit Rules ───────────────────────────────────────────────────

func (e *AuditEngine) RegisterDefaults() {
	e.rules = []AuditRule{
		&projectStructureRule{},
		&promptEngineAdoptionRule{},
		&missingDocumentationRule{},
		&missingDecisionsRule{},
		&architectureDriftRule{},
		&configDriftRule{},
	}
}

type projectStructureRule struct{}

func (r *projectStructureRule) ID() string      { return "project-structure" }
func (r *projectStructureRule) Area() AuditArea { return AreaProjectStructure }
func (r *projectStructureRule) Description() string {
	return "Checks standard project directories exist"
}
func (r *projectStructureRule) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	var findings []quality.Finding
	dirs := []string{"docs", ".promptengine"}
	for _, dir := range dirs {
		if !fs.Exists(dir) {
			findings = append(findings, quality.Finding{
				Engine:         "audit",
				Rule:           r.ID(),
				Category:       "project-structure",
				Severity:       quality.SeverityWarning,
				Title:          fmt.Sprintf("Missing directory: %s/", dir),
				Recommendation: fmt.Sprintf("Create the %s/ directory as part of project initialisation.", dir),
			})
		}
	}
	return findings, nil
}

type promptEngineAdoptionRule struct{}

func (r *promptEngineAdoptionRule) ID() string      { return "promptengine-adoption" }
func (r *promptEngineAdoptionRule) Area() AuditArea { return AreaPromptEngineAdoption }
func (r *promptEngineAdoptionRule) Description() string {
	return "Measures PromptEngine adoption level"
}
func (r *promptEngineAdoptionRule) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists("playbook-manifest.json") {
		return []quality.Finding{{
			Engine:         "audit",
			Rule:           r.ID(),
			Category:       "adoption",
			Severity:       quality.SeverityCritical,
			Title:          "PromptEngine not adopted",
			Explanation:    "No playbook-manifest.json found — PromptEngine has not been initialised.",
			Recommendation: "Run 'promptengine init' to adopt PromptEngine.",
		}}, nil
	}
	return nil, nil
}

type missingDocumentationRule struct{}

func (r *missingDocumentationRule) ID() string          { return "missing-documentation" }
func (r *missingDocumentationRule) Area() AuditArea     { return AreaMissingDocumentation }
func (r *missingDocumentationRule) Description() string { return "Identifies missing core documents" }
func (r *missingDocumentationRule) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	var findings []quality.Finding
	critical := []string{"docs/Architecture.md", "docs/Database.md"}
	for _, p := range critical {
		if !fs.Exists(p) {
			findings = append(findings, quality.Finding{
				Engine:         "audit",
				Rule:           r.ID(),
				Category:       "documentation",
				Severity:       quality.SeverityError,
				Title:          fmt.Sprintf("Missing: %s", p),
				Recommendation: "Run 'promptengine generate' to create this document.",
				FilePath:       p,
			})
		}
	}
	return findings, nil
}

type missingDecisionsRule struct{}

func (r *missingDecisionsRule) ID() string          { return "missing-decisions" }
func (r *missingDecisionsRule) Area() AuditArea     { return AreaMissingDecisions }
func (r *missingDecisionsRule) Description() string { return "Checks for Decisions.md" }
func (r *missingDecisionsRule) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists("docs/Decisions.md") {
		return []quality.Finding{{
			Engine:         "audit",
			Rule:           r.ID(),
			Category:       "documentation",
			Severity:       quality.SeverityWarning,
			Title:          "docs/Decisions.md is missing",
			Explanation:    "Architecture decisions should be formally recorded.",
			Recommendation: "Run 'promptengine generate decisions' to scaffold the document.",
		}}, nil
	}
	return nil, nil
}

type architectureDriftRule struct{}

func (r *architectureDriftRule) ID() string      { return "architecture-drift" }
func (r *architectureDriftRule) Area() AuditArea { return AreaArchitectureDrift }
func (r *architectureDriftRule) Description() string {
	return "Detects architecture documentation drift"
}
func (r *architectureDriftRule) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	// Full implementation integrates with Discovery Engine diff.
	return nil, nil
}

type configDriftRule struct{}

func (r *configDriftRule) ID() string          { return "config-drift" }
func (r *configDriftRule) Area() AuditArea     { return AreaConfigDrift }
func (r *configDriftRule) Description() string { return "Detects configuration drift" }
func (r *configDriftRule) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	// Full implementation integrates with Manifest Engine diff.
	return nil, nil
}
