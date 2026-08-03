package workflows

import (
	"context"
	"fmt"
	"sync"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// PipelineStage defines an execution step inside a workflow
type PipelineStage interface {
	Name() string
	Run(ctx context.Context, fs filesystem.FileSystem, flow *FlowContext) error
}

// Workflow definition properties
type Workflow struct {
	Name   string
	Stages []PipelineStage
}

// Registry manages the set of registered built-in and plugin workflows
type Registry struct {
	mu        sync.RWMutex
	workflows map[string]Workflow
}

func NewRegistry() *Registry {
	return &Registry{
		workflows: make(map[string]Workflow),
	}
}

func (r *Registry) Register(w Workflow) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.workflows[w.Name] = w
}

func (r *Registry) Get(name string) (Workflow, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.workflows[name]
	return w, ok
}

// Publisher routes event updates to subscribers
type Publisher struct {
	mu        sync.Mutex
	listeners []EventListener
}

func NewPublisher() *Publisher {
	return &Publisher{
		listeners: make([]EventListener, 0),
	}
}

func (p *Publisher) Subscribe(l EventListener) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.listeners = append(p.listeners, l)
}

func (p *Publisher) Publish(e Event) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, listener := range p.listeners {
		listener(e)
	}
}

// Engine coordinates pipeline runs
type Engine struct {
	registry  *Registry
	publisher *Publisher
	fs        filesystem.FileSystem
}

func NewEngine(fs filesystem.FileSystem, reg *Registry, pub *Publisher) *Engine {
	return &Engine{
		fs:        fs,
		registry:  reg,
		publisher: pub,
	}
}

func (e *Engine) RunWorkflow(ctx context.Context, flowName string, flowCtx *FlowContext) (State, error) {
	w, ok := e.registry.Get(flowName)
	if !ok {
		return StateFailed, fmt.Errorf("workflow '%s' not registered in context", flowName)
	}

	e.publisher.Publish(Event{Type: EventWorkflowStarted, Message: "Starting workflow pipeline " + flowName, Payload: flowCtx})

	for _, stage := range w.Stages {
		if err := stage.Run(ctx, e.fs, flowCtx); err != nil {
			e.publisher.Publish(Event{Type: EventWorkflowFailed, Message: fmt.Sprintf("Stage %s failed: %v", stage.Name(), err), Payload: err})
			return StateFailed, err
		}
	}

	e.publisher.Publish(Event{Type: EventWorkflowCompleted, Message: "Completed workflow execution " + flowName, Payload: flowCtx})
	return StateCompleted, nil
}
