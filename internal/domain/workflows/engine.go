package workflows

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

type PipelineStage interface {
	Name() string
	Run(ctx context.Context, fs filesystem.FileSystem, flow *FlowContext) error
}

type StepHandler interface {
	Execute(ctx context.Context, step WorkflowStep, flow *FlowContext) (interface{}, error)
}

type StepHandlerFunc func(ctx context.Context, step WorkflowStep, flow *FlowContext) (interface{}, error)

func (f StepHandlerFunc) Execute(ctx context.Context, step WorkflowStep, flow *FlowContext) (interface{}, error) {
	return f(ctx, step, flow)
}

type Registry struct {
	mu        sync.RWMutex
	workflows map[string]Workflow
	sources   map[string]string
}

func NewRegistry() *Registry {
	r := &Registry{workflows: map[string]Workflow{}, sources: map[string]string{}}
	r.RegisterBuiltIns()
	return r
}

func (r *Registry) Register(w Workflow) {
	r.RegisterFromSource("built-in", w)
}

func (r *Registry) RegisterFromSource(source string, w Workflow) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if w.ID == "" {
		w.ID = w.Name
	}
	if w.Name == "" {
		w.Name = w.ID
	}
	r.workflows[w.ID] = w
	r.sources[w.ID] = source
}

func (r *Registry) Get(id string) (Workflow, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.workflows[id]
	if ok {
		return w, true
	}
	for _, candidate := range r.workflows {
		if candidate.Name == id {
			return candidate, true
		}
	}
	return Workflow{}, false
}

func (r *Registry) All() []Workflow {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Workflow, 0, len(r.workflows))
	for _, w := range r.workflows {
		out = append(out, w)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (r *Registry) RegisterBuiltIns() {
	r.RegisterFromSource("built-in", Workflow{
		ID: "new-project", Name: "New Project", Description: "Initialize a new project workflow.", Version: "1.0.0",
		RequiredEngines: []string{"discovery", "context", "documentation", "validation"},
		Steps: []WorkflowStep{
			{ID: "initialize", Name: "Initialize", Order: 1, Action: "initialize"},
			{ID: "discovery", Name: "Discovery", Order: 2, Action: "discovery", Dependencies: []string{"initialize"}, RequiredEngine: "discovery"},
			{ID: "context", Name: "Context", Order: 3, Action: "context", Dependencies: []string{"discovery"}, RequiredEngine: "context"},
			{ID: "documentation_setup", Name: "Documentation Setup", Order: 4, Action: "documentation_setup", Dependencies: []string{"context"}, RequiredEngine: "documentation"},
			{ID: "validation", Name: "Validation", Order: 5, Action: "validation", Dependencies: []string{"documentation_setup"}, RequiredEngine: "validation"},
		},
	})
	r.RegisterFromSource("built-in", Workflow{
		ID: "feature-implementation", Name: "Feature Implementation", Description: "Prepare context for feature implementation.", Version: "1.0.0",
		RequiredEngines: []string{"discovery", "context", "prompt", "validation"},
		Steps: []WorkflowStep{
			{ID: "discovery", Name: "Discovery", Order: 1, Action: "discovery", RequiredEngine: "discovery"},
			{ID: "context", Name: "Context", Order: 2, Action: "context", Dependencies: []string{"discovery"}, RequiredEngine: "context"},
			{ID: "prompt_preparation", Name: "Prompt Preparation", Order: 3, Action: "prompt_preparation", Dependencies: []string{"context"}, RequiredEngine: "prompt"},
			{ID: "validation", Name: "Validation", Order: 4, Action: "validation", Dependencies: []string{"prompt_preparation"}, RequiredEngine: "validation"},
		},
	})
	r.RegisterFromSource("built-in", Workflow{
		ID: "bug-fix", Name: "Bug Fix", Description: "Prepare context for bug fixing.", Version: "1.0.0",
		RequiredEngines: []string{"discovery", "context", "review", "validation"},
		Steps: []WorkflowStep{
			{ID: "discovery", Name: "Discovery", Order: 1, Action: "discovery", RequiredEngine: "discovery"},
			{ID: "context", Name: "Context", Order: 2, Action: "context", Dependencies: []string{"discovery"}, RequiredEngine: "context"},
			{ID: "review", Name: "Review", Order: 3, Action: "review", Dependencies: []string{"context"}, RequiredEngine: "review"},
			{ID: "validation", Name: "Validation", Order: 4, Action: "validation", Dependencies: []string{"review"}, RequiredEngine: "validation"},
		},
	})
}

type Engine struct {
	registry *Registry
	eventBus *eventbus.EventBus
	fs       filesystem.FileSystem
	handlers map[string]StepHandler
}

func NewEngine(fs filesystem.FileSystem, reg *Registry, eb *eventbus.EventBus) *Engine {
	if reg == nil {
		reg = NewRegistry()
	}
	if eb == nil {
		eb = eventbus.NewEventBus()
	}
	e := &Engine{fs: fs, registry: reg, eventBus: eb, handlers: map[string]StepHandler{}}
	e.registerNoopHandlers("initialize", "documentation_setup", "prompt_preparation", "review", "validation")
	return e
}

func (e *Engine) RegisterHandler(action string, handler StepHandler) {
	e.handlers[action] = handler
}

func (e *Engine) LoadFromManifest(m *manifest.Manifest) {
	if m == nil {
		return
	}
	for _, mw := range m.Workflows {
		steps := make([]WorkflowStep, 0, len(mw.Steps))
		for i, stepName := range mw.Steps {
			steps = append(steps, WorkflowStep{ID: stepName, Name: stepName, Order: i + 1, Action: stepName})
		}
		e.registry.RegisterFromSource("manifest", Workflow{
			ID: mw.ID, Name: mw.ID, Version: m.Metadata.SchemaVersion, Steps: steps,
			Inputs: mw.RequiredContext, RequiredEngines: inferRequiredEngines(steps),
		})
	}
}

func (e *Engine) RunWorkflow(ctx context.Context, flowName string, flowCtx *FlowContext) (State, error) {
	exec, err := e.Execute(ctx, flowName, flowCtx)
	if err != nil {
		return exec.Status, err
	}
	return exec.Status, nil
}

func (e *Engine) Execute(ctx context.Context, workflowID string, flow *FlowContext) (*WorkflowExecution, error) {
	w, ok := e.registry.Get(workflowID)
	if !ok {
		exec := newExecution(workflowID)
		exec.Status = StateFailed
		err := fmt.Errorf("workflow %q not registered", workflowID)
		exec.Errors = append(exec.Errors, WorkflowError{Message: err.Error(), RecommendedAction: "Register the workflow before executing it."})
		return exec, err
	}
	if flow == nil {
		flow = NewFlowContext(w.ID)
	}
	exec := newExecution(w.ID)
	flow.Execution = exec
	if err := ValidateWorkflow(w); err != nil {
		return e.fail(exec, "", err, "Fix the workflow definition."), err
	}
	exec.Status = StateRunning
	e.publish(EventWorkflowStarted, "workflow started", exec)

	results := map[string]bool{}
	steps := orderedSteps(w)
	for _, step := range steps {
		if err := ctx.Err(); err != nil {
			exec.Status = StateCancelled
			exec.EndTime = time.Now().UTC()
			e.publish(EventWorkflowCancelled, "workflow cancelled", exec)
			return exec, err
		}
		for _, dep := range step.Dependencies {
			if !results[dep] {
				err := fmt.Errorf("step %s dependency %s has not completed", step.ID, dep)
				return e.fail(exec, step.ID, err, "Verify step dependency order."), err
			}
		}
		handler, ok := e.handlers[step.Action]
		if !ok && len(w.Stages) > 0 {
			continue
		}
		if !ok {
			err := fmt.Errorf("no handler registered for action %q", step.Action)
			return e.fail(exec, step.ID, err, "Register a step handler for this action."), err
		}
		exec.CurrentStep = step.ID
		e.publish(EventWorkflowStepStarted, "workflow step started", step)
		result, err := handler.Execute(ctx, step, flow)
		if err != nil {
			return e.fail(exec, step.ID, err, step.FailureHandling.RecommendedAction), err
		}
		flow.Outputs[step.ID] = result
		exec.Results[step.ID] = result
		results[step.ID] = true
		e.publish(EventWorkflowStepCompleted, "workflow step completed", step)
	}

	for _, stage := range w.Stages {
		e.publish(EventWorkflowStepStarted, "workflow legacy stage started", stage.Name())
		if err := stage.Run(ctx, e.fs, flow); err != nil {
			return e.fail(exec, stage.Name(), err, "Inspect legacy stage failure."), err
		}
		e.publish(EventWorkflowStepCompleted, "workflow legacy stage completed", stage.Name())
	}

	exec.Status = StateCompleted
	exec.EndTime = time.Now().UTC()
	exec.CurrentStep = ""
	e.publish(EventWorkflowCompleted, "workflow completed", exec)
	return exec, nil
}

func ValidateWorkflow(w Workflow) error {
	if w.ID == "" && w.Name == "" {
		return fmt.Errorf("workflow id is required")
	}
	ids := map[string]bool{}
	for _, step := range w.Steps {
		if step.ID == "" {
			return fmt.Errorf("workflow step id is required")
		}
		if step.Action == "" {
			return fmt.Errorf("workflow step %s action is required", step.ID)
		}
		if ids[step.ID] {
			return fmt.Errorf("duplicate workflow step %q", step.ID)
		}
		ids[step.ID] = true
	}
	for _, step := range w.Steps {
		for _, dep := range step.Dependencies {
			if !ids[dep] {
				return fmt.Errorf("workflow step %s references missing dependency %s", step.ID, dep)
			}
		}
	}
	return nil
}

func (e *Engine) fail(exec *WorkflowExecution, stepID string, err error, recommended string) *WorkflowExecution {
	exec.Status = StateFailed
	exec.CurrentStep = stepID
	exec.EndTime = time.Now().UTC()
	exec.Errors = append(exec.Errors, WorkflowError{StepID: stepID, Message: err.Error(), RecommendedAction: recommended})
	e.publish(EventWorkflowFailed, "workflow failed", exec)
	return exec
}

func (e *Engine) publish(t eventbus.EventType, msg string, payload interface{}) {
	if e.eventBus != nil {
		e.eventBus.Publish(eventbus.Event{Type: t, Message: msg, Payload: payload})
	}
}

func (e *Engine) registerNoopHandlers(actions ...string) {
	for _, action := range actions {
		a := action
		e.RegisterHandler(a, StepHandlerFunc(func(ctx context.Context, step WorkflowStep, flow *FlowContext) (interface{}, error) {
			return map[string]string{"status": "prepared", "action": a}, nil
		}))
	}
}

func newExecution(workflowID string) *WorkflowExecution {
	now := time.Now().UTC()
	return &WorkflowExecution{
		WorkflowID:  workflowID,
		ExecutionID: fmt.Sprintf("%s-%d", workflowID, now.UnixNano()),
		StartTime:   now,
		Status:      StatePending,
		Results:     map[string]interface{}{},
	}
}

func orderedSteps(w Workflow) []WorkflowStep {
	steps := append([]WorkflowStep(nil), w.Steps...)
	sort.SliceStable(steps, func(i, j int) bool { return steps[i].Order < steps[j].Order })
	return steps
}

func inferRequiredEngines(steps []WorkflowStep) []string {
	seen := map[string]bool{}
	var out []string
	for _, step := range steps {
		engine := step.RequiredEngine
		if engine == "" {
			engine = step.Action
		}
		if !seen[engine] {
			seen[engine] = true
			out = append(out, engine)
		}
	}
	return out
}
