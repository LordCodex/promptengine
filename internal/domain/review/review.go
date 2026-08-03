package review

import (
	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// Finding is now an alias for the consolidated quality.Finding struct
type Finding = quality.Finding

// ─── ReviewType ────────────────────────────────────────────────────────────

type ReviewType string

const (
	ReviewArchitecture    ReviewType = "architecture"
	ReviewCodeOrganization ReviewType = "code-organization"
	ReviewDocumentation   ReviewType = "documentation"
	ReviewSecurity        ReviewType = "security"
	ReviewPerformance     ReviewType = "performance"
	ReviewTesting         ReviewType = "testing"
	ReviewMaintainability ReviewType = "maintainability"
	ReviewScalability     ReviewType = "scalability"
	ReviewDeployment      ReviewType = "deployment"
	ReviewObservability   ReviewType = "observability"
	ReviewCompliance      ReviewType = "promptengine-compliance"
	ReviewOrgStandards    ReviewType = "org-standards"
	ReviewTechBestPractice ReviewType = "tech-best-practice"
)

// ─── Rule (original preserved) ─────────────────────────────────────────────

// Rule checks a target file/code path for standard violations
type Rule interface {
	Name() string
	Evaluate(fs filesystem.FileSystem, path string) ([]Finding, error)
}

// ─── Reviewer ──────────────────────────────────────────────────────────────

// Reviewer groups a set of rules under a specific ReviewType
type Reviewer interface {
	Type() ReviewType
	Description() string
	Rules() []Rule
}

// ─── ReviewSession aggregates findings across all reviewers ────────────────

type ReviewSession struct {
	TargetPath string
	Findings   []Finding
	Summary    map[ReviewType]int // count of findings per review type
}

func (s *ReviewSession) Add(findings ...Finding) {
	s.Findings = append(s.Findings, findings...)
}

func (s *ReviewSession) CountBySeverity(severity string) int {
	n := 0
	for _, f := range s.Findings {
		if string(f.Severity) == severity {
			n++
		}
	}
	return n
}

// ─── Registry (original preserved + extended) ──────────────────────────────

// Registry manages code review rules list
type Registry struct {
	rules     []Rule
	reviewers []Reviewer
}

func NewRegistry() *Registry {
	return &Registry{
		rules:     make([]Rule, 0),
		reviewers: make([]Reviewer, 0),
	}
}

func (r *Registry) Register(rule Rule) {
	r.rules = append(r.rules, rule)
}

func (r *Registry) RegisterReviewer(reviewer Reviewer) {
	r.reviewers = append(r.reviewers, reviewer)
}

// EngineName implements quality.EngineRegistrar
func (r *Registry) EngineName() string { return "review" }

// Review runs all registered rules against a target path (original API preserved)
func (r *Registry) Review(fs filesystem.FileSystem, targetPath string) ([]Finding, error) {
	var findings []Finding
	for _, rule := range r.rules {
		res, err := rule.Evaluate(fs, targetPath)
		if err != nil {
			return nil, err
		}
		findings = append(findings, res...)
	}
	return findings, nil
}

// RunSession runs all Reviewers (each contributing grouped rules) and returns a ReviewSession
func (r *Registry) RunSession(fs filesystem.FileSystem, targetPath string) (*ReviewSession, error) {
	session := &ReviewSession{
		TargetPath: targetPath,
		Summary:    make(map[ReviewType]int),
	}

	for _, reviewer := range r.reviewers {
		for _, rule := range reviewer.Rules() {
			findings, err := rule.Evaluate(fs, targetPath)
			if err != nil {
				return nil, err
			}
			for _, f := range findings {
				session.Add(f)
				session.Summary[reviewer.Type()]++
			}
		}
	}

	// Also run any directly registered rules
	for _, rule := range r.rules {
		findings, err := rule.Evaluate(fs, targetPath)
		if err != nil {
			return nil, err
		}
		session.Add(findings...)
	}

	return session, nil
}

// ─── Built-in Reviewers ────────────────────────────────────────────────────

// RegisterDefaultReviewers adds the standard PromptEngine review set
func RegisterDefaultReviewers(reg *Registry) {
	reg.RegisterReviewer(&documentationReviewer{})
	reg.RegisterReviewer(&securityReviewer{})
	reg.RegisterReviewer(&testingReviewer{})
	reg.RegisterReviewer(&complianceReviewer{})
}

// documentationReviewer checks documentation completeness
type documentationReviewer struct{}

func (r *documentationReviewer) Type() ReviewType    { return ReviewDocumentation }
func (r *documentationReviewer) Description() string { return "Checks project documentation completeness" }
func (r *documentationReviewer) Rules() []Rule       { return []Rule{&docsExistRule{}} }

type docsExistRule struct{}

func (r *docsExistRule) Name() string { return "docs-exist" }
func (r *docsExistRule) Evaluate(fs filesystem.FileSystem, path string) ([]Finding, error) {
	if !fs.Exists("docs") {
		return []Finding{{
			Engine:         "review",
			Rule:           r.Name(),
			Category:       string(ReviewDocumentation),
			Severity:       quality.SeverityError,
			Title:          "docs/ directory is missing",
			Recommendation: "run 'promptengine generate' to scaffold documentation",
			Explanation:    "Create docs/ directory with all standard PromptEngine documents",
		}}, nil
	}
	return nil, nil
}

// securityReviewer checks for common security red flags
type securityReviewer struct{}

func (r *securityReviewer) Type() ReviewType    { return ReviewSecurity }
func (r *securityReviewer) Description() string { return "Checks for common security issues" }
func (r *securityReviewer) Rules() []Rule       { return []Rule{&securityDocRule{}} }

type securityDocRule struct{}

func (r *securityDocRule) Name() string { return "security-doc-exists" }
func (r *securityDocRule) Evaluate(fs filesystem.FileSystem, path string) ([]Finding, error) {
	if !fs.Exists("docs/Security.md") {
		return []Finding{{
			Engine:         "review",
			Rule:           r.Name(),
			Category:       string(ReviewSecurity),
			Severity:       quality.SeverityError,
			Title:          "docs/Security.md is missing",
			Recommendation: "run 'promptengine generate security' to scaffold security documentation",
		}}, nil
	}
	return nil, nil
}

// testingReviewer checks testing standards
type testingReviewer struct{}

func (r *testingReviewer) Type() ReviewType    { return ReviewTesting }
func (r *testingReviewer) Description() string { return "Checks testing standards" }
func (r *testingReviewer) Rules() []Rule       { return []Rule{&testingDocRule{}} }

type testingDocRule struct{}

func (r *testingDocRule) Name() string { return "testing-doc-exists" }
func (r *testingDocRule) Evaluate(fs filesystem.FileSystem, path string) ([]Finding, error) {
	if !fs.Exists("docs/Testing.md") {
		return []Finding{{
			Engine:         "review",
			Rule:           r.Name(),
			Category:       string(ReviewTesting),
			Severity:       quality.SeveritySuggestion,
			Title:          "docs/Testing.md is missing",
			Recommendation: "run 'promptengine generate testing' to scaffold testing documentation",
		}}, nil
	}
	return nil, nil
}

// complianceReviewer checks PromptEngine compliance
type complianceReviewer struct{}

func (r *complianceReviewer) Type() ReviewType    { return ReviewCompliance }
func (r *complianceReviewer) Description() string { return "Checks PromptEngine compliance" }
func (r *complianceReviewer) Rules() []Rule       { return []Rule{&manifestRule{}} }

type manifestRule struct{}

func (r *manifestRule) Name() string { return "manifest-exists" }
func (r *manifestRule) Evaluate(fs filesystem.FileSystem, path string) ([]Finding, error) {
	if !fs.Exists("playbook-manifest.json") {
		return []Finding{{
			Engine:         "review",
			Rule:           r.Name(),
			Category:       string(ReviewCompliance),
			Severity:       quality.SeverityCritical,
			Title:          "playbook-manifest.json is missing",
			Recommendation: "run 'promptengine init' to initialise PromptEngine",
		}}, nil
	}
	return nil, nil
}
