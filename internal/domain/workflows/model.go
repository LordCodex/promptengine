package workflows

import (
	"time"

	ctxengine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/eventbus"
)

type State string

const (
	StatePending   State = "pending"
	StateRunning   State = "running"
	StateCompleted State = "completed"
	StateFailed    State = "failed"
	StateCancelled State = "cancelled"
)

type EventType = eventbus.EventType

const (
	EventWorkflowStarted       EventType = eventbus.WorkflowStarted
	EventWorkflowStepStarted   EventType = eventbus.WorkflowStepStarted
	EventWorkflowStepCompleted EventType = eventbus.WorkflowStepCompleted
	EventWorkflowCompleted     EventType = eventbus.WorkflowCompleted
	EventWorkflowFailed        EventType = eventbus.WorkflowFailed
	EventWorkflowCancelled     EventType = eventbus.WorkflowCancelled
	EventContextGenerated      EventType = eventbus.ContextGenerated
	EventValidationPassed      EventType = eventbus.ValidationPassed
)

type Event = eventbus.Event
type EventListener func(event Event)

type Workflow struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	Version         string         `json:"version"`
	Steps           []WorkflowStep `json:"steps"`
	Inputs          []string       `json:"inputs,omitempty"`
	Outputs         []string       `json:"outputs,omitempty"`
	RequiredEngines []string       `json:"required_engines,omitempty"`

	// Legacy field retained for older tests and callers.
	Stages []PipelineStage `json:"-"`
}

type WorkflowStep struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Order           int               `json:"order"`
	Action          string            `json:"action"`
	Dependencies    []string          `json:"dependencies,omitempty"`
	InputMapping    map[string]string `json:"input_mapping,omitempty"`
	OutputMapping   map[string]string `json:"output_mapping,omitempty"`
	FailureHandling FailureHandling   `json:"failure_handling,omitempty"`
	RequiredEngine  string            `json:"required_engine,omitempty"`
}

type FailureHandling struct {
	Strategy          string `json:"strategy,omitempty"`
	RecommendedAction string `json:"recommended_action,omitempty"`
}

type WorkflowExecution struct {
	WorkflowID  string                 `json:"workflow_id"`
	ExecutionID string                 `json:"execution_id"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     time.Time              `json:"end_time,omitempty"`
	Status      State                  `json:"status"`
	CurrentStep string                 `json:"current_step,omitempty"`
	Results     map[string]interface{} `json:"results"`
	Errors      []WorkflowError        `json:"errors,omitempty"`
}

type WorkflowError struct {
	StepID            string `json:"step_id"`
	Message           string `json:"message"`
	RecommendedAction string `json:"recommended_action,omitempty"`
}

type FlowContext struct {
	TaskName          string
	Variables         map[string]string
	Project           *discovery.ProjectModel
	SelectedContext   *ctxengine.ContextPackage
	PreconditionsMet  bool
	PostconditionsMet bool
	ValidationErrors  []string
	Execution         *WorkflowExecution
	Inputs            map[string]interface{}
	Outputs           map[string]interface{}
}

func NewFlowContext(taskName string) *FlowContext {
	return &FlowContext{
		TaskName:         taskName,
		Variables:        map[string]string{},
		ValidationErrors: []string{},
		Inputs:           map[string]interface{}{},
		Outputs:          map[string]interface{}{},
	}
}
