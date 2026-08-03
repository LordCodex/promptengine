package validation

import (
	"fmt"
	"sync"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// Validator is the general-purpose interface for all quality validators.
// This is distinct from the docs-layer validator — it operates at project level.
type Validator interface {
	ID() string
	Category() string
	Validate(fs filesystem.FileSystem) ([]quality.Finding, error)
}

// Registry is the plugin-extensible catalogue of all quality validators
type Registry struct {
	mu         sync.RWMutex
	validators []Validator
}

func NewRegistry() *Registry {
	r := &Registry{}
	r.RegisterDefaults()
	return r
}

func (r *Registry) Register(v Validator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.validators = append(r.validators, v)
}

// EngineName implements quality.EngineRegistrar
func (r *Registry) EngineName() string { return "validation" }

// Run executes all validators and returns aggregated findings
func (r *Registry) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	r.mu.RLock()
	validators := r.validators
	r.mu.RUnlock()

	var findings []quality.Finding
	for _, v := range validators {
		f, err := v.Validate(fs)
		if err != nil {
			return nil, fmt.Errorf("validator '%s' failed: %w", v.ID(), err)
		}
		findings = append(findings, f...)
	}
	return findings, nil
}

func (r *Registry) RegisterDefaults() {
	r.validators = []Validator{
		&projectConfigValidator{},
		&promptEngineConfigValidator{},
		&manifestValidator{},
		&documentationValidator{},
		&hooksValidator{},
		&technologyCompatibilityValidator{},
	}
}

// ─── Built-in Validators ───────────────────────────────────────────────────

type projectConfigValidator struct{}

func (v *projectConfigValidator) ID() string       { return "project-config" }
func (v *projectConfigValidator) Category() string { return "configuration" }
func (v *projectConfigValidator) Validate(fs filesystem.FileSystem) ([]quality.Finding, error) {
	var findings []quality.Finding
	if !fs.Exists("playbook-manifest.json") {
		findings = append(findings, quality.Finding{
			Engine:         "validation",
			Rule:           v.ID(),
			Category:       v.Category(),
			Severity:       quality.SeverityCritical,
			Title:          "Missing playbook-manifest.json",
			Explanation:    "The PromptEngine manifest is required for all quality features.",
			Recommendation: "Run 'promptengine init' to generate the manifest.",
		})
	}
	return findings, nil
}

type promptEngineConfigValidator struct{}

func (v *promptEngineConfigValidator) ID() string       { return "promptengine-config" }
func (v *promptEngineConfigValidator) Category() string { return "configuration" }
func (v *promptEngineConfigValidator) Validate(fs filesystem.FileSystem) ([]quality.Finding, error) {
	var findings []quality.Finding
	if !fs.Exists(".promptengine") {
		findings = append(findings, quality.Finding{
			Engine:         "validation",
			Rule:           v.ID(),
			Category:       v.Category(),
			Severity:       quality.SeverityError,
			Title:          "PromptEngine not initialised",
			Explanation:    ".promptengine directory is missing.",
			Recommendation: "Run 'promptengine init'.",
		})
	}
	return findings, nil
}

type manifestValidator struct{}

func (v *manifestValidator) ID() string       { return "manifest-schema" }
func (v *manifestValidator) Category() string { return "integrity" }
func (v *manifestValidator) Validate(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists("playbook-manifest.json") {
		return nil, nil // covered by project-config
	}
	data, err := fs.ReadFile("playbook-manifest.json")
	if err != nil {
		return nil, err
	}
	if len(data) < 2 { // must be at least "{}"
		return []quality.Finding{{
			Engine:         "validation",
			Rule:           v.ID(),
			Category:       v.Category(),
			Severity:       quality.SeverityError,
			Title:          "playbook-manifest.json is empty",
			Recommendation: "Restore or regenerate the manifest.",
		}}, nil
	}
	return nil, nil
}

type documentationValidator struct{}

func (v *documentationValidator) ID() string       { return "documentation-completeness" }
func (v *documentationValidator) Category() string { return "documentation" }
func (v *documentationValidator) Validate(fs filesystem.FileSystem) ([]quality.Finding, error) {
	var findings []quality.Finding
	required := []string{
		"docs/Architecture.md",
		"docs/Database.md",
		"docs/API.md",
	}
	for _, path := range required {
		if !fs.Exists(path) {
			findings = append(findings, quality.Finding{
				Engine:         "validation",
				Rule:           v.ID(),
				Category:       v.Category(),
				Severity:       quality.SeverityWarning,
				Title:          fmt.Sprintf("Missing core document: %s", path),
				Recommendation: fmt.Sprintf("Run 'promptengine generate' to create %s.", path),
				FilePath:       path,
			})
		}
	}
	return findings, nil
}

type hooksValidator struct{}

func (v *hooksValidator) ID() string       { return "hooks-configured" }
func (v *hooksValidator) Category() string { return "integrity" }
func (v *hooksValidator) Validate(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists(".git") {
		return []quality.Finding{{
			Engine:         "validation",
			Rule:           v.ID(),
			Category:       v.Category(),
			Severity:       quality.SeverityInfo,
			Title:          "Not a Git repository",
			Explanation:    "Git hooks cannot be installed without a .git directory.",
			Recommendation: "Initialise Git with 'git init' before running promptengine sync.",
		}}, nil
	}
	return nil, nil
}

type technologyCompatibilityValidator struct{}

func (v *technologyCompatibilityValidator) ID() string       { return "technology-compatibility" }
func (v *technologyCompatibilityValidator) Category() string { return "compatibility" }
func (v *technologyCompatibilityValidator) Validate(fs filesystem.FileSystem) ([]quality.Finding, error) {
	// Compatibility analysis is performed by the Discovery Engine at runtime.
	// This validator acts as the quality-layer entry point.
	return nil, nil
}
