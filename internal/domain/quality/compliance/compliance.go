package compliance

import (
	"fmt"
	"sync"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// ProfileType classifies a compliance profile
type ProfileType string

const (
	ProfilePromptEngine  ProfileType = "promptengine"
	ProfileOrganization  ProfileType = "organization"
	ProfileTechnology    ProfileType = "technology"
	ProfileSecurity      ProfileType = "security"
	ProfileDocumentation ProfileType = "documentation"
	ProfileNaming        ProfileType = "naming"
	ProfileArchitecture  ProfileType = "architecture"
)

// ComplianceRule is a single check within a compliance profile
type ComplianceRule interface {
	ID() string
	Description() string
	Check(fs filesystem.FileSystem) ([]quality.Finding, error)
}

// ComplianceProfile groups rules under a named, typed profile
type ComplianceProfile struct {
	ID          string
	Name        string
	Type        ProfileType
	Description string
	Rules       []ComplianceRule
}

// ComplianceReport is the full result of a compliance run
type ComplianceReport struct {
	ProfileResults []ProfileResult
	Overall        quality.Score
}

// ProfileResult is the result for a single compliance profile
type ProfileResult struct {
	Profile  ComplianceProfile
	Findings []quality.Finding
	Passed   bool
}

// ComplianceEngine runs all registered compliance profiles
type ComplianceEngine struct {
	mu       sync.RWMutex
	profiles []ComplianceProfile
}

func NewComplianceEngine() *ComplianceEngine {
	e := &ComplianceEngine{}
	e.RegisterDefaults()
	return e
}

func (e *ComplianceEngine) RegisterProfile(profile ComplianceProfile) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.profiles = append(e.profiles, profile)
}

// EngineName implements quality.EngineRegistrar
func (e *ComplianceEngine) EngineName() string { return "compliance" }

// Run evaluates all profiles and returns an aggregated compliance report
func (e *ComplianceEngine) Run(fs filesystem.FileSystem) (*ComplianceReport, error) {
	e.mu.RLock()
	profiles := e.profiles
	e.mu.RUnlock()

	var allFindings []quality.Finding
	var profileResults []ProfileResult

	for _, profile := range profiles {
		var profileFindings []quality.Finding
		for _, rule := range profile.Rules {
			f, err := rule.Check(fs)
			if err != nil {
				return nil, fmt.Errorf("compliance rule '%s' in profile '%s' failed: %w", rule.ID(), profile.ID, err)
			}
			profileFindings = append(profileFindings, f...)
		}
		allFindings = append(allFindings, profileFindings...)
		passed := true
		for _, f := range profileFindings {
			if f.Severity == quality.SeverityCritical || f.Severity == quality.SeverityError {
				passed = false
				break
			}
		}
		profileResults = append(profileResults, ProfileResult{
			Profile:  profile,
			Findings: profileFindings,
			Passed:   passed,
		})
	}

	score := computeComplianceScore(allFindings)
	return &ComplianceReport{
		ProfileResults: profileResults,
		Overall:        score,
	}, nil
}

func computeComplianceScore(findings []quality.Finding) quality.Score {
	cats := []quality.CategoryScore{
		{Category: "promptengine", Weight: 0.30, Raw: 100},
		{Category: "organization", Weight: 0.25, Raw: 100},
		{Category: "security", Weight: 0.25, Raw: 100},
		{Category: "documentation", Weight: 0.20, Raw: 100},
	}
	deduct := map[string]int{}
	for _, f := range findings {
		switch f.Severity {
		case quality.SeverityCritical:
			deduct[f.Category] += 50
		case quality.SeverityError:
			deduct[f.Category] += 25
		case quality.SeverityWarning:
			deduct[f.Category] += 10
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

// ─── Default Profiles ──────────────────────────────────────────────────────

func (e *ComplianceEngine) RegisterDefaults() {
	e.profiles = []ComplianceProfile{
		{
			ID:          "promptengine-core",
			Name:        "PromptEngine Core Standards",
			Type:        ProfilePromptEngine,
			Description: "Core PromptEngine adoption and configuration requirements",
			Rules: []ComplianceRule{
				&manifestRequiredRule{},
				&docsDirectoryRule{},
				&decisionsDocRule{},
			},
		},
		{
			ID:          "security-baseline",
			Name:        "Security Baseline",
			Type:        ProfileSecurity,
			Description: "Minimum security documentation requirements",
			Rules: []ComplianceRule{
				&securityDocRequiredRule{},
			},
		},
	}
}

// ─── Built-in compliance rules ─────────────────────────────────────────────

type manifestRequiredRule struct{}

func (r *manifestRequiredRule) ID() string          { return "manifest-required" }
func (r *manifestRequiredRule) Description() string { return "playbook-manifest.json must exist" }
func (r *manifestRequiredRule) Check(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists("playbook-manifest.json") {
		return []quality.Finding{{
			Engine:         "compliance",
			Rule:           r.ID(),
			Category:       "promptengine",
			Severity:       quality.SeverityCritical,
			Title:          "PromptEngine manifest is missing",
			Recommendation: "Run 'promptengine init'.",
		}}, nil
	}
	return nil, nil
}

type docsDirectoryRule struct{}

func (r *docsDirectoryRule) ID() string          { return "docs-directory-required" }
func (r *docsDirectoryRule) Description() string { return "docs/ directory must exist" }
func (r *docsDirectoryRule) Check(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists("docs") {
		return []quality.Finding{{
			Engine:         "compliance",
			Rule:           r.ID(),
			Category:       "documentation",
			Severity:       quality.SeverityError,
			Title:          "docs/ directory is missing",
			Recommendation: "Run 'promptengine generate'.",
		}}, nil
	}
	return nil, nil
}

type decisionsDocRule struct{}

func (r *decisionsDocRule) ID() string          { return "decisions-doc-required" }
func (r *decisionsDocRule) Description() string { return "Decisions.md must exist" }
func (r *decisionsDocRule) Check(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists("docs/Decisions.md") {
		return []quality.Finding{{
			Engine:         "compliance",
			Rule:           r.ID(),
			Category:       "documentation",
			Severity:       quality.SeverityWarning,
			Title:          "docs/Decisions.md is missing",
			Recommendation: "Run 'promptengine generate decisions'.",
		}}, nil
	}
	return nil, nil
}

type securityDocRequiredRule struct{}

func (r *securityDocRequiredRule) ID() string          { return "security-doc-required" }
func (r *securityDocRequiredRule) Description() string { return "Security.md must exist" }
func (r *securityDocRequiredRule) Check(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists("docs/Security.md") {
		return []quality.Finding{{
			Engine:         "compliance",
			Rule:           r.ID(),
			Category:       "security",
			Severity:       quality.SeverityError,
			Title:          "docs/Security.md is missing",
			Recommendation: "Run 'promptengine generate security'.",
		}}, nil
	}
	return nil, nil
}
