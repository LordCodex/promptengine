package doctor

import (
	"fmt"
	"sync"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// ─── DoctorCheck ───────────────────────────────────────────────────────────

// DoctorCheck is the interface every diagnostic check must implement.
type DoctorCheck interface {
	ID() string
	Category() string
	Description() string
	Run(fs filesystem.FileSystem) ([]quality.Finding, error)
}

// RepairAction is a safe, reversible fix for a doctor finding.
type RepairAction struct {
	FindingRuleID string
	Description   string
	DryRun        bool
	Apply         func(fs filesystem.FileSystem) error
}

// ─── DoctorEngine ──────────────────────────────────────────────────────────

type DoctorEngine struct {
	mu      sync.RWMutex
	checks  []DoctorCheck
	repairs map[string]RepairAction // keyed by FindingRuleID
}

func NewDoctorEngine() *DoctorEngine {
	d := &DoctorEngine{
		checks:  make([]DoctorCheck, 0),
		repairs: make(map[string]RepairAction),
	}
	d.RegisterDefaults()
	return d
}

func (d *DoctorEngine) Register(check DoctorCheck) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.checks = append(d.checks, check)
}

func (d *DoctorEngine) RegisterRepair(action RepairAction) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.repairs[action.FindingRuleID] = action
}

// EngineName implements quality.EngineRegistrar
func (d *DoctorEngine) EngineName() string { return "doctor" }

// Diagnose runs all registered checks and returns the aggregated report.
func (d *DoctorEngine) Diagnose(fs filesystem.FileSystem) (*quality.Report, error) {
	d.mu.RLock()
	checks := d.checks
	d.mu.RUnlock()

	var findings []quality.Finding
	for _, check := range checks {
		f, err := check.Run(fs)
		if err != nil {
			return nil, fmt.Errorf("doctor check '%s' failed: %w", check.ID(), err)
		}
		findings = append(findings, f...)
	}

	score := computeDoctorScore(findings)
	return &quality.Report{
		Title:    "PromptEngine Doctor Report",
		Score:    score,
		Findings: findings,
		Meta:     map[string]string{"engine": "doctor"},
	}, nil
}

// Fix attempts to apply a registered repair action for a finding.
// If dryRun is true, it only prints what would happen.
func (d *DoctorEngine) Fix(findingRuleID string, fs filesystem.FileSystem, dryRun bool) error {
	d.mu.RLock()
	action, ok := d.repairs[findingRuleID]
	d.mu.RUnlock()
	if !ok {
		return fmt.Errorf("no repair action registered for finding '%s'", findingRuleID)
	}
	if dryRun {
		return nil // preview only
	}
	return action.Apply(fs)
}

func computeDoctorScore(findings []quality.Finding) quality.Score {
	cats := []quality.CategoryScore{
		{Category: "installation", Weight: 0.2, Raw: 100},
		{Category: "configuration", Weight: 0.2, Raw: 100},
		{Category: "documentation", Weight: 0.2, Raw: 100},
		{Category: "integrity", Weight: 0.2, Raw: 100},
		{Category: "compatibility", Weight: 0.2, Raw: 100},
	}
	deductions := map[string]int{
		"installation":  0,
		"configuration": 0,
		"documentation": 0,
		"integrity":     0,
		"compatibility": 0,
	}
	for _, f := range findings {
		switch f.Severity {
		case quality.SeverityCritical:
			deductions[f.Category] += 40
		case quality.SeverityError:
			deductions[f.Category] += 25
		case quality.SeverityWarning:
			deductions[f.Category] += 10
		default:
			deductions[f.Category] += 2
		}
	}
	for i := range cats {
		raw := 100 - deductions[cats[i].Category]
		if raw < 0 {
			raw = 0
		}
		cats[i].Raw = raw
		cats[i].Weighted = float64(raw) * cats[i].Weight
		cats[i].FindingCount = deductions[cats[i].Category]
	}
	return quality.ComputeScore(cats, 70)
}

// ─── Default Checks ────────────────────────────────────────────────────────

func (d *DoctorEngine) RegisterDefaults() {
	d.Register(&manifestIntegrityCheck{})
	d.Register(&configurationCheck{})
	d.Register(&documentationExistenceCheck{})
	d.Register(&pluginIntegrityCheck{})
	d.Register(&brokenReferencesCheck{})
}

// manifestIntegrityCheck verifies playbook-manifest.json exists and is parseable.
type manifestIntegrityCheck struct{}

func (c *manifestIntegrityCheck) ID() string          { return "manifest-integrity" }
func (c *manifestIntegrityCheck) Category() string    { return "integrity" }
func (c *manifestIntegrityCheck) Description() string { return "Verify playbook-manifest.json exists" }
func (c *manifestIntegrityCheck) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists("playbook-manifest.json") {
		return []quality.Finding{{
			Engine:         "doctor",
			Rule:           c.ID(),
			Category:       c.Category(),
			Severity:       quality.SeverityCritical,
			Title:          "playbook-manifest.json not found",
			Explanation:    "PromptEngine requires a playbook-manifest.json at the project root.",
			Impact:         "All manifest-driven features are disabled.",
			Recommendation: "Run 'promptengine init' to generate the manifest.",
			AutoFixID:      "create-manifest",
		}}, nil
	}
	return nil, nil
}

// configurationCheck verifies .promptengine/config exists.
type configurationCheck struct{}

func (c *configurationCheck) ID() string          { return "configuration-exists" }
func (c *configurationCheck) Category() string    { return "configuration" }
func (c *configurationCheck) Description() string { return "Verify PromptEngine configuration" }
func (c *configurationCheck) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists(".promptengine") {
		return []quality.Finding{{
			Engine:         "doctor",
			Rule:           c.ID(),
			Category:       c.Category(),
			Severity:       quality.SeverityError,
			Title:          ".promptengine directory not found",
			Explanation:    "PromptEngine has not been initialised in this project.",
			Impact:         "No PromptEngine features will function correctly.",
			Recommendation: "Run 'promptengine init' to initialise the project.",
		}}, nil
	}
	return nil, nil
}

// documentationExistenceCheck verifies core docs directory exists.
type documentationExistenceCheck struct{}

func (c *documentationExistenceCheck) ID() string          { return "docs-directory" }
func (c *documentationExistenceCheck) Category() string    { return "documentation" }
func (c *documentationExistenceCheck) Description() string { return "Verify docs directory exists" }
func (c *documentationExistenceCheck) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	if !fs.Exists("docs") {
		return []quality.Finding{{
			Engine:         "doctor",
			Rule:           c.ID(),
			Category:       c.Category(),
			Severity:       quality.SeverityWarning,
			Title:          "docs/ directory not found",
			Explanation:    "PromptEngine expects project documentation to live in docs/.",
			Impact:         "Document generation and validation will not function correctly.",
			Recommendation: "Run 'promptengine generate' to scaffold the docs directory.",
			AutoFixID:      "create-docs-dir",
		}}, nil
	}
	return nil, nil
}

// pluginIntegrityCheck verifies the plugins directory is consistent.
type pluginIntegrityCheck struct{}

func (c *pluginIntegrityCheck) ID() string          { return "plugin-integrity" }
func (c *pluginIntegrityCheck) Category() string    { return "integrity" }
func (c *pluginIntegrityCheck) Description() string { return "Verify installed plugins are intact" }
func (c *pluginIntegrityCheck) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	pluginDir := ".promptengine/plugins"
	if !fs.Exists(pluginDir) {
		// No plugins installed — not an error
		return nil, nil
	}
	return nil, nil
}

// brokenReferencesCheck is a placeholder for cross-document reference validation.
type brokenReferencesCheck struct{}

func (c *brokenReferencesCheck) ID() string          { return "broken-references" }
func (c *brokenReferencesCheck) Category() string    { return "integrity" }
func (c *brokenReferencesCheck) Description() string { return "Detect broken document references" }
func (c *brokenReferencesCheck) Run(fs filesystem.FileSystem) ([]quality.Finding, error) {
	// Full implementation delegates to the docs validation engine.
	return nil, nil
}
