package workflows

import (
	"context"
	"fmt"
	"testing"

	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

func TestWorkflowRegistry_LoadingBuiltIns(t *testing.T) {
	reg := NewRegistry()
	if _, ok := reg.Get("new-project"); !ok {
		t.Fatal("expected new-project built-in workflow")
	}
	if _, ok := reg.Get("feature-implementation"); !ok {
		t.Fatal("expected feature-implementation built-in workflow")
	}
	if _, ok := reg.Get("bug-fix"); !ok {
		t.Fatal("expected bug-fix built-in workflow")
	}
}

func TestWorkflowEngine_StepExecutionOrderAndOutputs(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	reg := NewRegistry()
	reg.RegisterFromSource("test", Workflow{
		ID: "ordered",
		Steps: []WorkflowStep{
			{ID: "second", Name: "Second", Order: 2, Action: "second", Dependencies: []string{"first"}},
			{ID: "first", Name: "First", Order: 1, Action: "first"},
		},
	})
	engine := NewEngine(fs, reg, eventbus.NewEventBus())
	var order []string
	engine.RegisterHandler("first", StepHandlerFunc(func(ctx context.Context, step WorkflowStep, flow *FlowContext) (interface{}, error) {
		order = append(order, step.ID)
		return "one", nil
	}))
	engine.RegisterHandler("second", StepHandlerFunc(func(ctx context.Context, step WorkflowStep, flow *FlowContext) (interface{}, error) {
		order = append(order, step.ID)
		if flow.Outputs["first"] != "one" {
			t.Fatalf("expected first output to be available")
		}
		return "two", nil
	}))

	exec, err := engine.Execute(context.Background(), "ordered", NewFlowContext("ordered"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exec.Status != StateCompleted {
		t.Fatalf("expected completed, got %s", exec.Status)
	}
	if fmt.Sprint(order) != "[first second]" {
		t.Fatalf("unexpected order %v", order)
	}
}

func TestWorkflowEngine_FailedStepIsMachineReadable(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterFromSource("test", Workflow{
		ID:    "failing",
		Steps: []WorkflowStep{{ID: "fail", Name: "Fail", Order: 1, Action: "fail", FailureHandling: FailureHandling{RecommendedAction: "Fix the mock."}}},
	})
	events := eventbus.NewEventBus()
	failedEvent := false
	events.Subscribe(EventWorkflowFailed, func(e Event) { failedEvent = true })
	engine := NewEngine(filesystem.NewMockFileSystem(), reg, events)
	engine.RegisterHandler("fail", StepHandlerFunc(func(ctx context.Context, step WorkflowStep, flow *FlowContext) (interface{}, error) {
		return nil, fmt.Errorf("boom")
	}))

	exec, err := engine.Execute(context.Background(), "failing", NewFlowContext("failing"))
	if err == nil {
		t.Fatal("expected failure")
	}
	if exec.Status != StateFailed || exec.CurrentStep != "fail" {
		t.Fatalf("expected failed execution at fail step, got %#v", exec)
	}
	if len(exec.Errors) != 1 || exec.Errors[0].RecommendedAction != "Fix the mock." {
		t.Fatalf("expected machine-readable error with recommendation, got %#v", exec.Errors)
	}
	if !failedEvent {
		t.Fatal("expected WorkflowFailed event")
	}
}

func TestWorkflowEngine_EventPublishing(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterFromSource("test", Workflow{ID: "events", Steps: []WorkflowStep{{ID: "a", Order: 1, Action: "a"}}})
	events := eventbus.NewEventBus()
	var seen []EventType
	for _, eventType := range []EventType{EventWorkflowStarted, EventWorkflowStepStarted, EventWorkflowStepCompleted, EventWorkflowCompleted} {
		tp := eventType
		events.Subscribe(tp, func(e Event) { seen = append(seen, e.Type) })
	}
	engine := NewEngine(filesystem.NewMockFileSystem(), reg, events)
	engine.RegisterHandler("a", StepHandlerFunc(func(ctx context.Context, step WorkflowStep, flow *FlowContext) (interface{}, error) {
		return "ok", nil
	}))
	if _, err := engine.Execute(context.Background(), "events", NewFlowContext("events")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	for _, eventType := range []EventType{EventWorkflowStarted, EventWorkflowStepStarted, EventWorkflowStepCompleted, EventWorkflowCompleted} {
		if !hasEvent(seen, eventType) {
			t.Fatalf("expected event %s, got %v", eventType, seen)
		}
	}
}

func TestWorkflowEngine_DependencyInjection(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterFromSource("test", Workflow{ID: "di", Steps: []WorkflowStep{{ID: "discovery", Order: 1, Action: "discovery"}}})
	mock := &mockDiscoveryHandler{called: false}
	engine := NewEngine(filesystem.NewMockFileSystem(), reg, eventbus.NewEventBus())
	engine.RegisterHandler("discovery", mock)
	if _, err := engine.Execute(context.Background(), "di", NewFlowContext("di")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !mock.called {
		t.Fatal("expected injected discovery handler to be called")
	}
}

func TestWorkflowEngine_LoadFromManifest(t *testing.T) {
	reg := NewRegistry()
	engine := NewEngine(filesystem.NewMockFileSystem(), reg, eventbus.NewEventBus())
	engine.LoadFromManifest(&manifest.Manifest{
		Metadata:  manifest.ProjectMetadata{Name: "Project", SchemaVersion: manifest.SupportedSchemaVersion},
		Workflows: []manifest.WorkflowDefinition{{ID: "manifest-flow", Steps: []string{"discovery", "context"}, RequiredContext: []string{"project"}}},
	})
	w, ok := reg.Get("manifest-flow")
	if !ok {
		t.Fatal("expected workflow loaded from manifest")
	}
	if len(w.Steps) != 2 || w.Steps[1].Action != "context" {
		t.Fatalf("unexpected manifest workflow %#v", w)
	}
}

func TestWorkflowRegistry_PluginWorkflowRegistration(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterFromSource("plugin:example", Workflow{ID: "plugin-flow", Name: "Plugin Flow", Steps: []WorkflowStep{{ID: "a", Order: 1, Action: "a"}}})
	if _, ok := reg.Get("plugin-flow"); !ok {
		t.Fatal("expected plugin workflow registration")
	}
}

func TestValidateWorkflow_InvalidDependency(t *testing.T) {
	err := ValidateWorkflow(Workflow{ID: "bad", Steps: []WorkflowStep{{ID: "a", Action: "a", Dependencies: []string{"missing"}}}})
	if err == nil {
		t.Fatal("expected invalid dependency error")
	}
}

type mockDiscoveryHandler struct {
	called bool
}

func (m *mockDiscoveryHandler) Execute(ctx context.Context, step WorkflowStep, flow *FlowContext) (interface{}, error) {
	m.called = true
	return "discovered", nil
}

func hasEvent(events []EventType, expected EventType) bool {
	for _, eventType := range events {
		if eventType == expected {
			return true
		}
	}
	return false
}
