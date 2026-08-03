// Package quality provides the shared types, registries, scoring framework,
// and reporting primitives consumed by all Quality Platform engines.
package quality

import (
	"encoding/json"
	"fmt"
)

// ─── Severity ──────────────────────────────────────────────────────────────

type Severity string

const (
	SeverityCritical   Severity = "critical"  // blocks CI/deploy
	SeverityError      Severity = "error"      // blocks commit
	SeverityWarning    Severity = "warning"    // advisory
	SeverityInfo       Severity = "info"       // informational
	SeveritySuggestion Severity = "suggestion" // optional improvement
)

// ─── Finding ───────────────────────────────────────────────────────────────

// Finding is the canonical quality observation emitted by every engine.
type Finding struct {
	Engine            string   // e.g. "doctor", "audit", "compliance"
	Rule              string   // rule/check ID
	Category          string   // e.g. "documentation", "security"
	Severity          Severity
	Title             string
	Explanation       string
	Impact            string
	Recommendation    string
	RelatedStandards  []string // PromptEngine doc IDs
	RelatedDocs       []string // project doc paths
	SuggestedWorkflow string
	AutoFixID         string // empty = no auto-fix available
	FilePath          string // optional - file the finding relates to
}

// ─── Score ─────────────────────────────────────────────────────────────────

// CategoryScore is a weighted score for a single quality dimension
type CategoryScore struct {
	Category       string
	Weight         float64 // 0.0–1.0; all weights in a registry must sum to 1.0
	Raw            int     // 0–100
	Weighted       float64 // Raw * Weight
	CriticalFail   bool    // if true, this category forces overall to 0
	FindingCount   int
}

// Score is the overall result of a quality evaluation
type Score struct {
	Overall    int
	Rating     string // A, B, C, D, F
	Categories []CategoryScore
	Passed     bool // true if Overall >= Threshold
	Threshold  int
}

func ComputeScore(categories []CategoryScore, threshold int) Score {
	critFail := false
	var totalWeight, weightedSum float64
	for _, c := range categories {
		if c.CriticalFail {
			critFail = true
		}
		totalWeight += c.Weight
		weightedSum += c.Weighted
	}

	overall := 0
	if !critFail && totalWeight > 0 {
		overall = int(weightedSum / totalWeight * 100)
		if overall > 100 {
			overall = 100
		}
	}

	rating := ratingFor(overall)
	return Score{
		Overall:    overall,
		Rating:     rating,
		Categories: categories,
		Passed:     overall >= threshold,
		Threshold:  threshold,
	}
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

// ─── Threshold ─────────────────────────────────────────────────────────────

// Threshold controls CI pass/fail behaviour
type Threshold struct {
	MinScore        int      // minimum overall score to pass
	BlockOnSeverity Severity // findings at or above this severity cause a block
}

var DefaultThreshold = Threshold{MinScore: 70, BlockOnSeverity: SeverityError}

// ─── Report ────────────────────────────────────────────────────────────────

// Report is the full output of a quality platform run
type Report struct {
	Title    string
	Score    Score
	Findings []Finding
	Meta     map[string]string
}

func (r *Report) CountBySeverity(s Severity) int {
	n := 0
	for _, f := range r.Findings {
		if f.Severity == s {
			n++
		}
	}
	return n
}

func (r *Report) ExceedsBLock(t Threshold) bool {
	for _, f := range r.Findings {
		if severityRank(f.Severity) >= severityRank(t.BlockOnSeverity) {
			return true
		}
	}
	return !r.Score.Passed
}

func (r *Report) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func severityRank(s Severity) int {
	switch s {
	case SeverityCritical:
		return 4
	case SeverityError:
		return 3
	case SeverityWarning:
		return 2
	case SeverityInfo:
		return 1
	default:
		return 0
	}
}

// ─── QualityRegistry ───────────────────────────────────────────────────────

// EngineRegistrar is implemented by each engine's plugin-extensible registry
type EngineRegistrar interface {
	EngineName() string
}

// QualityRegistry is the top-level plugin registration hub. All six engine
// registries are wired here so plugins call a single entry point.
type QualityRegistry struct {
	registrars map[string]EngineRegistrar
}

func NewQualityRegistry() *QualityRegistry {
	return &QualityRegistry{registrars: make(map[string]EngineRegistrar)}
}

func (q *QualityRegistry) Mount(r EngineRegistrar) {
	q.registrars[r.EngineName()] = r
}

func (q *QualityRegistry) Get(engine string) (EngineRegistrar, error) {
	r, ok := q.registrars[engine]
	if !ok {
		return nil, fmt.Errorf("quality registry: engine '%s' not mounted", engine)
	}
	return r, nil
}
