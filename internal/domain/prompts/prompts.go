package prompts

import (
	"fmt"
	"strings"
	"sync"
)

// PromptWorkflow identifies the engineering activity a prompt serves
type PromptWorkflow string

const (
	WorkflowNewProject    PromptWorkflow = "new-project"
	WorkflowExisting      PromptWorkflow = "existing-project"
	WorkflowMigration     PromptWorkflow = "migration"
	WorkflowFeature       PromptWorkflow = "feature-development"
	WorkflowBugFix        PromptWorkflow = "bug-fix"
	WorkflowRefactor      PromptWorkflow = "refactoring"
	WorkflowArchReview    PromptWorkflow = "architecture-review"
	WorkflowSecurityReview PromptWorkflow = "security-review"
	WorkflowPerfReview    PromptWorkflow = "performance-review"
	WorkflowDeployment    PromptWorkflow = "deployment"
	WorkflowDocReview     PromptWorkflow = "documentation-review"
	WorkflowAudit         PromptWorkflow = "project-audit"
	WorkflowBootstrap     PromptWorkflow = "ai-session-bootstrap"
	WorkflowContextRefresh PromptWorkflow = "context-refresh"
	WorkflowImprovement   PromptWorkflow = "prompt-improvement"
)

// PromptSource identifies where a prompt originates
type PromptSource string

const (
	PromptSourceCore       PromptSource = "core"
	PromptSourceOrg        PromptSource = "org"
	PromptSourcePlugin     PromptSource = "plugin"
	PromptSourceTechnology PromptSource = "technology"
	PromptSourceWorkflow   PromptSource = "workflow"
	PromptSourceProvider   PromptSource = "provider"
	PromptSourceCustom     PromptSource = "custom"
)

// PromptDef is the full definition of a reusable prompt
type PromptDef struct {
	ID             string
	Workflow       PromptWorkflow
	Source         PromptSource
	Purpose        string
	RequiredContext []string // context item keys that must be populated
	Variables      []string // variable names used in Template
	Template       string   // prompt body with {variable} placeholders
	ProviderHints  map[string]string // provider-id -> formatting hint
}

// Prompt is the original struct kept for backward compatibility
type Prompt struct {
	ID                 string
	Purpose            string
	CopyAndPastePrompt string
}

// ContextPackage carries gathered project context (supplied by Context Engine)
type ContextPackage map[string]string

// PromptRegistry is the multi-source registry of all prompt definitions
type PromptRegistry struct {
	mu    sync.RWMutex
	defs  map[string]*PromptDef
}

func NewPromptRegistry() *PromptRegistry {
	return &PromptRegistry{defs: make(map[string]*PromptDef)}
}

func (r *PromptRegistry) Register(def *PromptDef) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.defs[def.ID]; exists {
		return fmt.Errorf("prompt '%s' already registered", def.ID)
	}
	r.defs[def.ID] = def
	return nil
}

func (r *PromptRegistry) Get(id string) (*PromptDef, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.defs[id]
	return d, ok
}

func (r *PromptRegistry) ByWorkflow(wf PromptWorkflow) []*PromptDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*PromptDef
	for _, def := range r.defs {
		if def.Workflow == wf {
			result = append(result, def)
		}
	}
	return result
}

func (r *PromptRegistry) BySource(src PromptSource) []*PromptDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*PromptDef
	for _, def := range r.defs {
		if def.Source == src {
			result = append(result, def)
		}
	}
	return result
}

// PromptBuilder constructs a fully populated prompt from a PromptDef + ContextPackage
type PromptBuilder struct {
	registry *PromptRegistry
}

func NewPromptBuilder(reg *PromptRegistry) *PromptBuilder {
	return &PromptBuilder{registry: reg}
}

// Build returns a fully resolved Prompt. It merges the ContextPackage into vars,
// then applies provider-specific formatting if providerID is known.
func (b *PromptBuilder) Build(id string, ctx ContextPackage, providerID string) (*Prompt, error) {
	def, ok := b.registry.Get(id)
	if !ok {
		return nil, fmt.Errorf("prompt '%s' not registered", id)
	}

	// Merge context into variables (context takes precedence)
	vars := make(map[string]string)
	for k, v := range ctx {
		vars[k] = v
	}

	// Validate required context
	var missing []string
	for _, req := range def.RequiredContext {
		if vars[req] == "" {
			missing = append(missing, req)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("prompt '%s' is missing required context: %s", id, strings.Join(missing, ", "))
	}

	body := Inject(def.Template, vars)

	// Apply provider hint prefix/suffix if declared
	if hint, ok := def.ProviderHints[providerID]; ok {
		body = hint + "\n\n" + body
	}

	return &Prompt{
		ID:                 def.ID,
		Purpose:            def.Purpose,
		CopyAndPastePrompt: body,
	}, nil
}

// --- Legacy API (preserved for backward compatibility) ---

// Registry is the original simple registry struct
type Registry struct {
	promptsPath string
}

func NewRegistry(_ interface{}, promptsPath string) *Registry {
	return &Registry{promptsPath: promptsPath}
}

func (r *Registry) LoadPrompt(id string, variables map[string]string) (*Prompt, error) {
	return &Prompt{
		ID:                 id,
		Purpose:            "Scaffolds workflow actions",
		CopyAndPastePrompt: "Generated workflow prompt content details",
	}, nil
}

// Inject replaces standard bracketed variable parameters inside prompt text
func Inject(text string, vars map[string]string) string {
	result := text
	for k, v := range vars {
		result = strings.ReplaceAll(result, "{"+k+"}", v)
	}
	return result
}
