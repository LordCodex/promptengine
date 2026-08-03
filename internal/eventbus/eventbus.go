package eventbus

import (
	"sync"
)

// EventType categorises the event payload
type EventType string

const (
	ProjectDetected        EventType = "ProjectDetected"
	ProjectInitialized     EventType = "ProjectInitialized"
	ProjectMigrated        EventType = "ProjectMigrated"
	WorkflowStarted        EventType = "WorkflowStarted"
	WorkflowCompleted      EventType = "WorkflowCompleted"
	ContextBuilt           EventType = "ContextBuilt"
	PromptGenerated        EventType = "PromptGenerated"
	PromptExecuted         EventType = "PromptExecuted"
	DocumentationGenerated EventType = "DocumentationGenerated"
	DocumentationUpdated   EventType = "DocumentationUpdated"
	HealthCalculated       EventType = "HealthCalculated"
	ReviewCompleted        EventType = "ReviewCompleted"
	AuditCompleted         EventType = "AuditCompleted"
	ValidationCompleted    EventType = "ValidationCompleted"
	PluginInstalled        EventType = "PluginInstalled"
	PluginUpdated          EventType = "PluginUpdated"
	PluginRemoved          EventType = "PluginRemoved"
	ConfigurationChanged   EventType = "ConfigurationChanged"
	ManifestLoaded         EventType = "ManifestLoaded"
	CommandStarted         EventType = "CommandStarted"
	CommandCompleted       EventType = "CommandCompleted"
	HooksExecuted          EventType = "HooksExecuted"
)

// Event carries type and optional data payloads
type Event struct {
	Type    EventType
	Payload interface{}
}

// Handler functions process published events
type Handler func(e Event)

// EventBus coordinates event publishing and subscriptions synchronously by default
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[EventType][]Handler
}

// NewEventBus instantiates a clean EventBus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]Handler),
	}
}

// Subscribe adds a handler function for an EventType
func (b *EventBus) Subscribe(t EventType, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[t] = append(b.subscribers[t], h)
}

// Publish notifies all registered subscribers of the event synchronously
func (b *EventBus) Publish(e Event) {
	b.mu.RLock()
	handlers, ok := b.subscribers[e.Type]
	b.mu.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		handler(e)
	}
}

// PublishAsync notifies subscribers in non-blocking goroutines
func (b *EventBus) PublishAsync(e Event) {
	b.mu.RLock()
	handlers, ok := b.subscribers[e.Type]
	b.mu.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		go handler(e)
	}
}
