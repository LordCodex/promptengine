package fix

import (
	"fmt"
	"sync"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// FixSafety classifies how safe a repair action is
type FixSafety string

const (
	SafeAuto   FixSafety = "auto"   // safe to apply without confirmation
	SafeReview FixSafety = "review" // requires human review before apply
	SafeManual FixSafety = "manual" // cannot be automated; guidance only
)

// FixResult describes the outcome of applying a repair action
type FixResult struct {
	ActionID    string
	Applied     bool
	DryRun      bool
	Rolled      bool
	Description string
	Error       error
}

// RepairAction is a reversible, safe automated fix
type RepairAction interface {
	ID() string
	Description() string
	Safety() FixSafety
	Preview(fs filesystem.FileSystem) (string, error) // human-readable description of what will change
	Apply(fs filesystem.FileSystem) error
	Rollback(fs filesystem.FileSystem) error
}

// FixEngine orchestrates repair actions
type FixEngine struct {
	mu      sync.RWMutex
	actions map[string]RepairAction
	enabled bool // orgs can disable auto-fixes via manifest
}

func NewFixEngine() *FixEngine {
	e := &FixEngine{
		actions: make(map[string]RepairAction),
		enabled: true,
	}
	e.RegisterDefaults()
	return e
}

func (e *FixEngine) SetEnabled(enabled bool) {
	e.enabled = enabled
}

func (e *FixEngine) Register(action RepairAction) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.actions[action.ID()] = action
}

// EngineName implements quality.EngineRegistrar
func (e *FixEngine) EngineName() string { return "fix" }

// Preview returns a human-readable description of what a fix would do
func (e *FixEngine) Preview(id string, fs filesystem.FileSystem) (string, error) {
	action, err := e.get(id)
	if err != nil {
		return "", err
	}
	return action.Preview(fs)
}

// Apply executes a repair action; dryRun=true simulates without changes
func (e *FixEngine) Apply(id string, fs filesystem.FileSystem, dryRun bool) FixResult {
	if !e.enabled {
		return FixResult{ActionID: id, Applied: false, Description: "auto-fixes are disabled by organisation policy"}
	}

	action, err := e.get(id)
	if err != nil {
		return FixResult{ActionID: id, Error: err}
	}

	if action.Safety() == SafeManual {
		preview, _ := action.Preview(fs)
		return FixResult{ActionID: id, Applied: false, Description: "manual fix required: " + preview}
	}

	if dryRun {
		preview, _ := action.Preview(fs)
		return FixResult{ActionID: id, Applied: false, DryRun: true, Description: preview}
	}

	if err := action.Apply(fs); err != nil {
		return FixResult{ActionID: id, Applied: false, Error: err}
	}
	return FixResult{ActionID: id, Applied: true, Description: action.Description()}
}

// Rollback undoes a previously applied repair
func (e *FixEngine) Rollback(id string, fs filesystem.FileSystem) FixResult {
	action, err := e.get(id)
	if err != nil {
		return FixResult{ActionID: id, Error: err}
	}
	if err := action.Rollback(fs); err != nil {
		return FixResult{ActionID: id, Error: err}
	}
	return FixResult{ActionID: id, Rolled: true, Description: "rolled back: " + action.Description()}
}

func (e *FixEngine) get(id string) (RepairAction, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	action, ok := e.actions[id]
	if !ok {
		return nil, fmt.Errorf("no repair action registered for '%s'", id)
	}
	return action, nil
}

// ─── Default Repair Actions ────────────────────────────────────────────────

func (e *FixEngine) RegisterDefaults() {
	e.Register(&createDocsDirAction{})
	e.Register(&createManifestAction{})
	e.Register(&createPromptEngineDirAction{})
}

// createDocsDirAction creates the docs/ directory
type createDocsDirAction struct{}

func (a *createDocsDirAction) ID() string          { return "create-docs-dir" }
func (a *createDocsDirAction) Description() string { return "Create docs/ directory" }
func (a *createDocsDirAction) Safety() FixSafety   { return SafeAuto }
func (a *createDocsDirAction) Preview(_ filesystem.FileSystem) (string, error) {
	return "Will create docs/.keep to initialise the documentation directory.", nil
}
func (a *createDocsDirAction) Apply(fs filesystem.FileSystem) error {
	return fs.WriteFile("docs/.keep", []byte("# PromptEngine documentation directory\n"), 0644)
}
func (a *createDocsDirAction) Rollback(fs filesystem.FileSystem) error {
	// Marker removal — non-destructive
	return fs.WriteFile("docs/.keep", []byte(""), 0644)
}

// createManifestAction creates an empty playbook-manifest.json
type createManifestAction struct{}

func (a *createManifestAction) ID() string          { return "create-manifest" }
func (a *createManifestAction) Description() string { return "Create playbook-manifest.json" }
func (a *createManifestAction) Safety() FixSafety   { return SafeReview }
func (a *createManifestAction) Preview(_ filesystem.FileSystem) (string, error) {
	return "Will create playbook-manifest.json with default structure.", nil
}
func (a *createManifestAction) Apply(fs filesystem.FileSystem) error {
	skeleton := []byte(`{
  "version": "1.0",
  "project": {},
  "cli_foundation": [],
  "workflows": []
}
`)
	return fs.WriteFile("playbook-manifest.json", skeleton, 0644)
}
func (a *createManifestAction) Rollback(fs filesystem.FileSystem) error {
	return fs.WriteFile("playbook-manifest.json", []byte("{}"), 0644)
}

// createPromptEngineDirAction initialises the .promptengine directory
type createPromptEngineDirAction struct{}

func (a *createPromptEngineDirAction) ID() string          { return "create-promptengine-dir" }
func (a *createPromptEngineDirAction) Description() string { return "Create .promptengine directory" }
func (a *createPromptEngineDirAction) Safety() FixSafety   { return SafeAuto }
func (a *createPromptEngineDirAction) Preview(_ filesystem.FileSystem) (string, error) {
	return "Will create .promptengine/ directory.", nil
}
func (a *createPromptEngineDirAction) Apply(fs filesystem.FileSystem) error {
	return fs.WriteFile(".promptengine/.keep", []byte{}, 0644)
}
func (a *createPromptEngineDirAction) Rollback(fs filesystem.FileSystem) error {
	return fs.WriteFile(".promptengine/.keep", []byte("removed"), 0644)
}
