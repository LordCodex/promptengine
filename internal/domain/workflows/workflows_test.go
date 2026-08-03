package workflows

import (
	"context"
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestWorkflowEngine_PipelineAndEvents(t *testing.T) {
	fs := filesystem.NewMockFileSystem()

	// Write mock project files
	_ = fs.WriteFile("go.mod", []byte("module myapp"), 0644)
	_ = fs.WriteFile("docs/Database.md", []byte("database schema specs content"), 0644)
	_ = fs.WriteFile("AGENTS.md", []byte("Constitutional rules: reference to 05-universal-coding-standards.md"), 0644)

	// Create registries
	reg := NewRegistry()
	pub := NewPublisher()

	// Setup a mock workflow: precondition -> discovery -> postcondition
	w := Workflow{
		Name: "test-feature-flow",
		Stages: []PipelineStage{
			&PreconditionStage{RequiredFiles: []string{"go.mod"}},
			&DiscoveryStage{},
			&PostconditionStage{RequiredDocs: []string{"docs/Database.md"}},
		},
	}
	reg.Register(w)

	// Track lifecycle event dispatches
	var eventsReceived []EventType
	pub.Subscribe(func(e Event) {
		eventsReceived = append(eventsReceived, e.Type)
	})

	engine := NewEngine(fs, reg, pub)
	flowCtx := NewFlowContext("test-feature-flow")

	state, err := engine.RunWorkflow(context.Background(), "test-feature-flow", flowCtx)
	if err != nil {
		t.Fatalf("Expected workflow execution to complete successfully, got error: %v", err)
	}

	if state != StateCompleted {
		t.Errorf("Expected final state Completed, got %s", state)
	}

	// Verify events
	hasStarted := false
	hasCompleted := false
	for _, et := range eventsReceived {
		if et == EventWorkflowStarted {
			hasStarted = true
		}
		if et == EventWorkflowCompleted {
			hasCompleted = true
		}
	}

	if !hasStarted || !hasCompleted {
		t.Errorf("Expected lifecycle events to be published, got %v", eventsReceived)
	}

	// Verify preconditions and postconditions parameters
	if !flowCtx.PreconditionsMet {
		t.Errorf("Expected preconditions to validate successfully")
	}
	if !flowCtx.PostconditionsMet {
		t.Errorf("Expected postconditions to validate successfully")
	}
}

func TestWorkflowEngine_PreconditionFailure(t *testing.T) {
	fs := filesystem.NewMockFileSystem()

	// Missing required file
	reg := NewRegistry()
	pub := NewPublisher()

	w := Workflow{
		Name: "failing-flow",
		Stages: []PipelineStage{
			&PreconditionStage{RequiredFiles: []string{"missing.json"}},
		},
	}
	reg.Register(w)

	engine := NewEngine(fs, reg, pub)
	flowCtx := NewFlowContext("failing-flow")

	_, err := engine.RunWorkflow(context.Background(), "failing-flow", flowCtx)
	if err == nil {
		t.Errorf("Expected workflow to fail due to missing required files, got nil error")
	}

	if flowCtx.PreconditionsMet {
		t.Errorf("Expected preconditionsMet to remain false")
	}
}
