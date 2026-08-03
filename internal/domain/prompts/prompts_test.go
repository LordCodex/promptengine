package prompts

import (
	"strings"
	"testing"
)

func TestPromptRegistry_RegisterAndGet(t *testing.T) {
	reg := NewPromptRegistry()
	def := &PromptDef{
		ID:       "new-project-bootstrap",
		Workflow: WorkflowNewProject,
		Source:   PromptSourceCore,
		Purpose:  "Bootstrap a new project AI session",
		Template: "You are assisting with a new project: {project_name}. Stack: {stack}.",
	}
	if err := reg.Register(def); err != nil {
		t.Fatalf("expected no error on register, got: %v", err)
	}
	// Duplicate should fail
	if err := reg.Register(def); err == nil {
		t.Error("expected error on duplicate prompt registration")
	}
	got, ok := reg.Get("new-project-bootstrap")
	if !ok || got.Purpose == "" {
		t.Error("expected to retrieve registered prompt")
	}
}

func TestPromptRegistry_ByWorkflow(t *testing.T) {
	reg := NewPromptRegistry()
	_ = reg.Register(&PromptDef{ID: "new-1", Workflow: WorkflowNewProject, Source: PromptSourceCore, Template: ""})
	_ = reg.Register(&PromptDef{ID: "new-2", Workflow: WorkflowNewProject, Source: PromptSourceOrg, Template: ""})
	_ = reg.Register(&PromptDef{ID: "fix-1", Workflow: WorkflowBugFix, Source: PromptSourceCore, Template: ""})

	results := reg.ByWorkflow(WorkflowNewProject)
	if len(results) != 2 {
		t.Errorf("expected 2 prompts for new-project workflow, got %d", len(results))
	}
}

func TestPromptBuilder_Build_WithContext(t *testing.T) {
	reg := NewPromptRegistry()
	_ = reg.Register(&PromptDef{
		ID:             "feature-dev",
		Workflow:       WorkflowFeature,
		Source:         PromptSourceCore,
		Purpose:        "Feature development",
		RequiredContext: []string{"project_name"},
		Template:       "Build a feature for {project_name} using {stack}.",
	})

	builder := NewPromptBuilder(reg)
	prompt, err := builder.Build("feature-dev", ContextPackage{
		"project_name": "PromptEngine",
		"stack":        "Go",
	}, "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(prompt.CopyAndPastePrompt, "PromptEngine") {
		t.Error("expected built prompt to contain project name from context")
	}
}

func TestPromptBuilder_Build_MissingRequiredContext(t *testing.T) {
	reg := NewPromptRegistry()
	_ = reg.Register(&PromptDef{
		ID:             "guarded-prompt",
		Workflow:       WorkflowFeature,
		Source:         PromptSourceCore,
		RequiredContext: []string{"project_name", "stack"},
		Template:       "...",
	})
	builder := NewPromptBuilder(reg)
	if _, err := builder.Build("guarded-prompt", ContextPackage{}, ""); err == nil {
		t.Error("expected error when required context is missing")
	}
}

func TestPromptBuilder_ProviderHint(t *testing.T) {
	reg := NewPromptRegistry()
	_ = reg.Register(&PromptDef{
		ID:       "claude-prompt",
		Workflow: WorkflowBootstrap,
		Source:   PromptSourceProvider,
		Template: "Base prompt body.",
		ProviderHints: map[string]string{
			"claude": "<claude-hint>",
		},
	})
	builder := NewPromptBuilder(reg)
	prompt, _ := builder.Build("claude-prompt", ContextPackage{}, "claude")
	if !strings.HasPrefix(prompt.CopyAndPastePrompt, "<claude-hint>") {
		t.Error("expected provider hint to be prepended for claude provider")
	}
}

func TestInject(t *testing.T) {
	result := Inject("Hello {name}, you are on {stack}.", map[string]string{"name": "World", "stack": "Go"})
	if result != "Hello World, you are on Go." {
		t.Errorf("unexpected inject result: %s", result)
	}
}
