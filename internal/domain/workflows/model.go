package workflows

import (
	"github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
)

type State string

const (
	StatePending   State = "pending"
	StateRunning   State = "running"
	StatePaused    State = "paused"
	StateCompleted State = "completed"
	StateFailed    State = "failed"
	StateCancelled State = "cancelled"
)

type EventType string

const (
	EventWorkflowStarted    EventType = "WorkflowStarted"
	EventContextGenerated   EventType = "ContextGenerated"
	EventValidationPassed   EventType = "ValidationPassed"
	EventWorkflowCompleted EventType = "WorkflowCompleted"
	EventWorkflowFailed    EventType = "WorkflowFailed"
)

// Event tracks lifecycle transitions inside the execution pipeline
type Event struct {
	Type    EventType
	Message string
	Payload interface{}
}

// EventListener defines subscribers interfaces
type EventListener func(event Event)

// FlowContext holds variable state payloads during pipeline stages execution
type FlowContext struct {
	TaskName          string
	Variables         map[string]string
	Project           *discovery.ProjectModel
	SelectedContext   *context.ContextPackage
	PreconditionsMet  bool
	PostconditionsMet bool
	ValidationErrors  []string
}

func NewFlowContext(taskName string) *FlowContext {
	return &FlowContext{
		TaskName:         taskName,
		Variables:        make(map[string]string),
		ValidationErrors: make([]string, 0),
	}
}
