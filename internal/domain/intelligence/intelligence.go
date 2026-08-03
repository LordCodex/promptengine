package intelligence

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/domain/personal"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"gopkg.in/yaml.v3"
)

const DefaultDecisionPath = ".promptengine/decisions.yaml"

type Pattern struct {
	Name       string   `json:"name" yaml:"name"`
	Category   string   `json:"category" yaml:"category"`
	Evidence   []string `json:"evidence" yaml:"evidence"`
	Confidence float64  `json:"confidence" yaml:"confidence"`
}

type Decision struct {
	Title     string    `json:"title" yaml:"title"`
	Reason    string    `json:"reason" yaml:"reason"`
	Affected  []string  `json:"affected" yaml:"affected"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}

type Recommendation struct {
	Suggestion string   `json:"suggestion" yaml:"suggestion"`
	Reason     string   `json:"reason" yaml:"reason"`
	Confidence float64  `json:"confidence" yaml:"confidence"`
	Affected   []string `json:"affected_areas" yaml:"affected_areas"`
}

type SimilarReference struct {
	Path   string  `json:"path" yaml:"path"`
	Reason string  `json:"reason" yaml:"reason"`
	Score  float64 `json:"score" yaml:"score"`
}

type ImpactReport struct {
	ChangedFiles    []string         `json:"changed_files" yaml:"changed_files"`
	AffectedAreas   []string         `json:"affected_areas" yaml:"affected_areas"`
	Recommendations []Recommendation `json:"recommendations" yaml:"recommendations"`
}

type Insights struct {
	Patterns        []Pattern          `json:"patterns" yaml:"patterns"`
	Decisions       []Decision         `json:"decisions" yaml:"decisions"`
	Recommendations []Recommendation   `json:"recommendations" yaml:"recommendations"`
	References      []SimilarReference `json:"similar_references,omitempty" yaml:"similar_references,omitempty"`
}

type Platform struct {
	fs     filesystem.FileSystem
	events *eventbus.EventBus
}

func NewPlatform(fs filesystem.FileSystem, events *eventbus.EventBus) *Platform {
	return &Platform{fs: fs, events: events}
}

func (p *Platform) Analyze(project *discovery.ProjectModel) Insights {
	patterns := p.DetectPatterns(project)
	decisions, _ := p.ListDecisions()
	recs := p.Recommend(patterns, decisions, project)
	return Insights{Patterns: patterns, Decisions: decisions, Recommendations: recs}
}

func (p *Platform) DetectPatterns(project *discovery.ProjectModel) []Pattern {
	if project == nil {
		return nil
	}
	var patterns []Pattern
	hasDir := func(suffix string) []string {
		var hits []string
		for _, dir := range project.Repository.Directories {
			if strings.HasSuffix(strings.ToLower(dir), strings.ToLower(suffix)) {
				hits = append(hits, dir)
			}
		}
		return hits
	}
	add := func(name, category string, evidence []string, confidence float64) {
		if len(evidence) == 0 {
			return
		}
		pattern := Pattern{Name: name, Category: category, Evidence: uniqueSorted(evidence), Confidence: confidence}
		patterns = append(patterns, pattern)
		p.publish(eventbus.PatternDetected, "pattern detected", pattern)
	}
	add("Service Layer", "architecture", hasDir("Services"), 0.85)
	add("Action Classes", "architecture", hasDir("Actions"), 0.8)
	add("DTO/Data Objects", "architecture", append(hasDir("DTOs"), hasDir("Data")...), 0.75)
	add("Policy Authorization", "security", hasDir("Policies"), 0.75)
	add("Form Request Validation", "validation", hasDir("Requests"), 0.8)
	add("Flutter Feature Folders", "architecture", hasDir("features"), 0.75)
	add("Repository Pattern", "data", hasDir("repositories"), 0.75)
	add("Use Case Layer", "architecture", hasDir("usecases"), 0.7)
	add("Component-Based Frontend", "frontend", append(hasDir("components"), hasDir("hooks")...), 0.7)
	add("Documentation-Driven Development", "documentation", project.Repository.DocumentationFiles, 0.65)
	return patterns
}

func (p *Platform) StoreDecision(decision Decision) error {
	if strings.TrimSpace(decision.Title) == "" {
		return fmt.Errorf("decision title is required")
	}
	if decision.CreatedAt.IsZero() {
		decision.CreatedAt = time.Now().UTC()
	}
	decisions, _ := p.ListDecisions()
	decisions = append(decisions, decision)
	data, err := yaml.Marshal(map[string][]Decision{"decisions": decisions})
	if err != nil {
		return err
	}
	if err := p.fs.WriteFile(DefaultDecisionPath, data, 0644); err != nil {
		return err
	}
	p.publish(eventbus.DecisionStored, "decision stored", decision)
	return nil
}

func (p *Platform) ListDecisions() ([]Decision, error) {
	if !p.fs.Exists(DefaultDecisionPath) {
		return nil, nil
	}
	data, err := p.fs.ReadFile(DefaultDecisionPath)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Decisions []Decision `yaml:"decisions" json:"decisions"`
	}
	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	return wrapper.Decisions, nil
}

func (p *Platform) FindSimilar(task string, project *discovery.ProjectModel) []SimilarReference {
	if project == nil {
		return nil
	}
	terms := keywords(task)
	var refs []SimilarReference
	for _, path := range project.Repository.Files {
		lower := strings.ToLower(path)
		score := 0.0
		for _, term := range terms {
			if strings.Contains(lower, term) {
				score += 0.25
			}
		}
		if score > 0 {
			refs = append(refs, SimilarReference{Path: path, Score: min(score, 1), Reason: "Path matches task keyword(s)."})
		}
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].Score > refs[j].Score })
	if len(refs) > 10 {
		refs = refs[:10]
	}
	return refs
}

func (p *Platform) Recommend(patterns []Pattern, decisions []Decision, project *discovery.ProjectModel) []Recommendation {
	var recs []Recommendation
	for _, pattern := range patterns {
		switch pattern.Name {
		case "Service Layer":
			recs = append(recs, Recommendation{Suggestion: "Place business logic in existing service classes.", Reason: "Project already has a service-layer pattern.", Confidence: pattern.Confidence, Affected: []string{"services", "tests"}})
		case "Action Classes":
			recs = append(recs, Recommendation{Suggestion: "Use an action class for single-purpose application operations.", Reason: "Existing code uses action classes for workflow steps.", Confidence: pattern.Confidence, Affected: []string{"actions"}})
		case "Form Request Validation":
			recs = append(recs, Recommendation{Suggestion: "Keep request validation in Form Request objects.", Reason: "Form request folders are present.", Confidence: pattern.Confidence, Affected: []string{"controllers", "requests"}})
		case "Repository Pattern":
			recs = append(recs, Recommendation{Suggestion: "Follow the existing repository abstraction for data access.", Reason: "Repository folders were detected.", Confidence: pattern.Confidence, Affected: []string{"data", "tests"}})
		}
	}
	for _, decision := range decisions {
		recs = append(recs, Recommendation{Suggestion: "Respect decision: " + decision.Title, Reason: decision.Reason, Confidence: 0.9, Affected: decision.Affected})
	}
	for _, rec := range recs {
		p.publish(eventbus.RecommendationGenerated, "recommendation generated", rec)
	}
	return recs
}

func (p *Platform) AnalyzeImpact(git personal.GitContext, project *discovery.ProjectModel) ImpactReport {
	areas := map[string]bool{}
	for _, file := range git.ChangedFiles {
		lower := strings.ToLower(file)
		switch {
		case strings.Contains(lower, "user"):
			for _, area := range []string{"authentication", "api responses", "database relations", "tests", "documentation"} {
				areas[area] = true
			}
		case strings.Contains(lower, "model") || strings.Contains(lower, "migration") || strings.Contains(lower, "database"):
			for _, area := range []string{"database", "models", "api", "tests", "documentation"} {
				areas[area] = true
			}
		case strings.Contains(lower, "route") || strings.Contains(lower, "controller"):
			for _, area := range []string{"api", "authorization", "requests", "tests", "documentation"} {
				areas[area] = true
			}
		case strings.Contains(lower, "service"):
			for _, area := range []string{"business logic", "tests", "documentation"} {
				areas[area] = true
			}
		}
	}
	var affected []string
	for area := range areas {
		affected = append(affected, area)
	}
	sort.Strings(affected)
	recs := []Recommendation{}
	if len(affected) > 0 {
		recs = append(recs, Recommendation{Suggestion: "Review affected areas before finalizing changes.", Reason: "Changed file names imply related contracts may be affected.", Confidence: 0.7, Affected: affected})
	}
	report := ImpactReport{ChangedFiles: git.ChangedFiles, AffectedAreas: affected, Recommendations: recs}
	p.publish(eventbus.ImpactAnalysisCompleted, "impact analysis completed", report)
	return report
}

func FormatPromptEnhancement(insights Insights, impact ImpactReport) string {
	var b strings.Builder
	b.WriteString("\n## Existing Patterns\n\n")
	for _, pattern := range insights.Patterns {
		b.WriteString(fmt.Sprintf("- %s (%s, %.0f%%): %s\n", pattern.Name, pattern.Category, pattern.Confidence*100, strings.Join(pattern.Evidence, ", ")))
	}
	b.WriteString("\n## Previous Decisions\n\n")
	for _, decision := range insights.Decisions {
		b.WriteString(fmt.Sprintf("- %s: %s\n", decision.Title, decision.Reason))
	}
	b.WriteString("\n## Recommendations\n\n")
	for _, rec := range insights.Recommendations {
		b.WriteString(fmt.Sprintf("- %s Reason: %s Confidence: %.0f%%\n", rec.Suggestion, rec.Reason, rec.Confidence*100))
	}
	if len(impact.AffectedAreas) > 0 {
		b.WriteString("\n## Change Impact\n\n")
		for _, area := range impact.AffectedAreas {
			b.WriteString("- ")
			b.WriteString(area)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func keywords(task string) []string {
	words := strings.FieldsFunc(strings.ToLower(task), func(r rune) bool { return r < 'a' || r > 'z' })
	var out []string
	for _, word := range words {
		if len(word) >= 4 {
			out = append(out, word)
		}
	}
	return uniqueSorted(out)
}

func uniqueSorted(items []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, item := range items {
		item = filepath.ToSlash(strings.TrimSpace(item))
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	sort.Strings(out)
	return out
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (p *Platform) publish(t eventbus.EventType, msg string, payload interface{}) {
	if p.events != nil {
		p.events.Publish(eventbus.Event{Type: t, Message: msg, Payload: payload})
	}
}
