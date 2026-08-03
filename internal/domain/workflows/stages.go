package workflows

import (
	"context"
	"fmt"
	"path/filepath"

	ctxEngine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// DiscoveryStage executes standard repository audits
type DiscoveryStage struct{}

func (s *DiscoveryStage) Name() string { return "discovery_stage" }
func (s *DiscoveryStage) Run(ctx context.Context, fs filesystem.FileSystem, flow *FlowContext) error {
	pipeline := discovery.NewPipeline()
	// Register default stages
	pipeline.Register(&discovery.BaseStage{}, &discovery.PromptEngineStage{}, &discovery.TechStage{}, &discovery.ArchStage{}, &discovery.DocsStage{})

	pm, err := pipeline.Execute(ctx, fs, "")
	if err != nil {
		return err
	}
	flow.Project = pm
	return nil
}

// PreconditionStage asserts prerequisite file mappings
type PreconditionStage struct {
	RequiredFiles []string
}

func (s *PreconditionStage) Name() string { return "precondition_stage" }
func (s *PreconditionStage) Run(ctx context.Context, fs filesystem.FileSystem, flow *FlowContext) error {
	for _, f := range s.RequiredFiles {
		if !fs.Exists(f) {
			flow.ValidationErrors = append(flow.ValidationErrors, "Precondition failed: Missing file target "+f)
		}
	}

	if len(flow.ValidationErrors) > 0 {
		return fmt.Errorf("preconditions failed with %d validation error(s)", len(flow.ValidationErrors))
	}

	flow.PreconditionsMet = true
	return nil
}

// ContextStage executes playbooks context loading
type ContextStage struct {
	TaskType   ctxEngine.TaskType
	BudgetType ctxEngine.BudgetType
}

func (s *ContextStage) Name() string { return "context_stage" }
func (s *ContextStage) Run(ctx context.Context, fs filesystem.FileSystem, flow *FlowContext) error {
	if flow.Project == nil {
		return fmt.Errorf("context generation requires discovered project model metadata")
	}

	engine := ctxEngine.NewEngine(fs)
	pkg, err := engine.GenerateContext(ctx, s.TaskType, flow.Project, s.BudgetType)
	if err != nil {
		return err
	}

	flow.SelectedContext = pkg
	return nil
}

// PostconditionStage verifies final output compliance criteria
type PostconditionStage struct {
	RequiredDocs []string
}

func (s *PostconditionStage) Name() string { return "postcondition_stage" }
func (s *PostconditionStage) Run(ctx context.Context, fs filesystem.FileSystem, flow *FlowContext) error {
	// 1. Verify docs updates exist
	for _, rd := range s.RequiredDocs {
		if !fs.Exists(rd) {
			return fmt.Errorf("postcondition check failed: documentation file %s was not updated", rd)
		}
		
		// If exists, verify it isn't empty
		data, err := fs.ReadFile(rd)
		if err != nil || len(data) == 0 {
			return fmt.Errorf("postcondition check failed: documentation file %s is empty", rd)
		}
	}

	// 2. Validate linked path references in AGENTS.md
	if fs.Exists("AGENTS.md") {
		data, err := fs.ReadFile("AGENTS.md")
		if err == nil && len(data) > 0 {
			// Check if AGENTS.md inherits universal guidelines as required by AGENTS.md rules!
			content := string(data)
			if !filepath.IsAbs("core/05-universal-coding-standards.md") && !stringsContains(content, "05-universal-coding-standards.md") {
				return fmt.Errorf("postcondition check failed: AGENTS.md constitution must explicitly inherit from universal standards playbook")
			}
		}
	}

	flow.PostconditionsMet = true
	return nil
}

func stringsContains(s, sub string) bool {
	// Simple lookup helper to avoid string conversion overheads
	lenSub := len(sub)
	if lenSub == 0 {
		return true
	}
	lenS := len(s)
	if lenS < lenSub {
		return false
	}
	for i := 0; i <= lenS-lenSub; i++ {
		if s[i:i+lenSub] == sub {
			return true
		}
	}
	return false
}
