package health

import (
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// ─── Severity constants (kept for back-compat) ─────────────────────────────

const (
	SeverityBlock      = "block"
	SeverityImportant  = "important"
	SeveritySuggestion = "suggestion"
)

// Issue details a specific compliance drift finding (original type preserved)
type Issue struct {
	Category string
	Title    string
	Warning  string
	Severity string // "block", "important", "suggestion"
}

// Result logs computed metrics from workspace audits (original type preserved)
type Result struct {
	Score      int
	Rating     string // "A", "B", "C", "D", "F"
	Issues     []Issue
	Categories map[string]*CategoryResult // NEW: per-category breakdown
}

// CategoryResult holds per-dimension health data
type CategoryResult struct {
	Name         string
	Score        int
	Weight       float64
	CriticalFail bool
	Issues       []Issue
}

// HealthCategory declares a weighted quality dimension
type HealthCategory struct {
	Name         string
	Weight       float64 // 0.0–1.0; must sum to 1.0 across all categories
	CriticalRule string  // checker ID whose "block" finding zeros this category
}

// DefaultHealthCategories are the standard PromptEngine health dimensions
var DefaultHealthCategories = []HealthCategory{
	{Name: "documentation", Weight: 0.20},
	{Name: "architecture", Weight: 0.15},
	{Name: "security", Weight: 0.15},
	{Name: "testing", Weight: 0.10},
	{Name: "performance", Weight: 0.05},
	{Name: "maintainability", Weight: 0.10},
	{Name: "deployment", Weight: 0.05},
	{Name: "observability", Weight: 0.05},
	{Name: "project-knowledge", Weight: 0.05},
	{Name: "promptengine-adoption", Weight: 0.05},
	{Name: "dependency-health", Weight: 0.05},
}

// Checker validates a specific subsystem health (original interface preserved)
type Checker interface {
	Name() string
	Check(fs filesystem.FileSystem) ([]Issue, error)
}

// ─── Registry ──────────────────────────────────────────────────────────────

// Registry manages the set of registered checks (original preserved + extended)
type Registry struct {
	checkers   []Checker
	categories []HealthCategory
	threshold  int // minimum overall score to pass (default 70)
}

func NewRegistry() *Registry {
	return &Registry{
		checkers:   make([]Checker, 0),
		categories: DefaultHealthCategories,
		threshold:  70,
	}
}

// SetCategories allows manifest-driven weight overrides
func (r *Registry) SetCategories(cats []HealthCategory) {
	r.categories = cats
}

// SetThreshold configures the CI pass/fail score floor
func (r *Registry) SetThreshold(t int) {
	r.threshold = t
}

func (r *Registry) Register(c Checker) {
	r.checkers = append(r.checkers, c)
}

// EngineName implements quality.EngineRegistrar
func (r *Registry) EngineName() string { return "health" }

// Evaluate runs all checks and returns a weighted, categorised Result.
func (r *Registry) Evaluate(fs filesystem.FileSystem) (*Result, error) {
	// Group issues by category
	catIssues := make(map[string][]Issue)
	for _, cat := range r.categories {
		catIssues[cat.Name] = nil
	}

	for _, c := range r.checkers {
		issues, err := c.Check(fs)
		if err != nil {
			return nil, err
		}
		for _, issue := range issues {
			catIssues[issue.Category] = append(catIssues[issue.Category], issue)
		}
	}

	// Compute per-category scores
	catResults := make(map[string]*CategoryResult)
	var totalWeight, weightedSum float64
	critFail := false

	for _, cat := range r.categories {
		issues := catIssues[cat.Name]
		raw := 100
		for _, issue := range issues {
			switch issue.Severity {
			case SeverityBlock:
				raw -= 40
			case SeverityImportant:
				raw -= 15
			default:
				raw -= 5
			}
		}
		if raw < 0 {
			raw = 0
		}

		isCritical := false
		for _, issue := range issues {
			if issue.Severity == SeverityBlock {
				isCritical = true
				critFail = true
			}
		}

		catResults[cat.Name] = &CategoryResult{
			Name:         cat.Name,
			Score:        raw,
			Weight:       cat.Weight,
			CriticalFail: isCritical,
			Issues:       issues,
		}
		totalWeight += cat.Weight
		weightedSum += float64(raw) * cat.Weight
	}

	// Aggregate all issues for backward compatibility
	var allIssues []Issue
	for _, issues := range catIssues {
		allIssues = append(allIssues, issues...)
	}

	overall := 0
	if !critFail && totalWeight > 0 {
		overall = int(weightedSum / totalWeight)
	}
	if overall > 100 {
		overall = 100
	}

	rating := ratingFor(overall)
	return &Result{
		Score:      overall,
		Rating:     rating,
		Issues:     allIssues,
		Categories: catResults,
	}, nil
}

func ratingFor(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 65:
		return "C"
	case score >= 50:
		return "D"
	default:
		return "F"
	}
}
