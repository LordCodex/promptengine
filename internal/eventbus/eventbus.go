package eventbus

import (
	"sync"
)

// EventType categorises the event payload
type EventType string

const (
	ApplicationStarted            EventType = "ApplicationStarted"
	ApplicationReady              EventType = "ApplicationReady"
	ApplicationShutdown           EventType = "ApplicationShutdown"
	ProjectDiscoveryStarted       EventType = "ProjectDiscoveryStarted"
	ProjectDetected               EventType = "ProjectDetected"
	TechnologyDetected            EventType = "TechnologyDetected"
	ProjectDiscoveryCompleted     EventType = "ProjectDiscoveryCompleted"
	ProjectDiscoveryFailed        EventType = "ProjectDiscoveryFailed"
	ProjectInitialized            EventType = "ProjectInitialized"
	ProjectMigrated               EventType = "ProjectMigrated"
	WorkflowStarted               EventType = "WorkflowStarted"
	WorkflowStepStarted           EventType = "WorkflowStepStarted"
	WorkflowStepCompleted         EventType = "WorkflowStepCompleted"
	WorkflowCompleted             EventType = "WorkflowCompleted"
	WorkflowFailed                EventType = "WorkflowFailed"
	WorkflowCancelled             EventType = "WorkflowCancelled"
	ContextGenerated              EventType = "ContextGenerated"
	ContextBuildStarted           EventType = "ContextBuildStarted"
	ContextItemSelected           EventType = "ContextItemSelected"
	ContextBuildFailed            EventType = "ContextBuildFailed"
	ValidationStarted             EventType = "ValidationStarted"
	ValidationPassed              EventType = "ValidationPassed"
	ContextBuilt                  EventType = "ContextBuilt"
	PromptGenerated               EventType = "PromptGenerated"
	AIRequestStarted              EventType = "AIRequestStarted"
	AIRequestCompleted            EventType = "AIRequestCompleted"
	AIRequestFailed               EventType = "AIRequestFailed"
	AIStreamingStarted            EventType = "AIStreamingStarted"
	PromptExecuted                EventType = "PromptExecuted"
	DocumentationGenerated        EventType = "DocumentationGenerated"
	DocumentationUpdated          EventType = "DocumentationUpdated"
	DocumentationValidationFailed EventType = "DocumentationValidationFailed"
	DocumentationSyncRequired     EventType = "DocumentationSyncRequired"
	HealthCalculated              EventType = "HealthCalculated"
	CriticalIssueDetected         EventType = "CriticalIssueDetected"
	ReviewStarted                 EventType = "ReviewStarted"
	ReviewCompleted               EventType = "ReviewCompleted"
	AuditCompleted                EventType = "AuditCompleted"
	ValidationCompleted           EventType = "ValidationCompleted"
	PluginInstalled               EventType = "PluginInstalled"
	PluginEnabled                 EventType = "PluginEnabled"
	PluginDisabled                EventType = "PluginDisabled"
	PluginUpdated                 EventType = "PluginUpdated"
	PluginRemoved                 EventType = "PluginRemoved"
	PluginHealthFailed            EventType = "PluginHealthFailed"
	AgentConfigGenerated          EventType = "AgentConfigGenerated"
	AgentConfigUpdated            EventType = "AgentConfigUpdated"
	ContextExported               EventType = "ContextExported"
	PromptPackageCreated          EventType = "PromptPackageCreated"
	PatternDetected               EventType = "PatternDetected"
	DecisionStored                EventType = "DecisionStored"
	RecommendationGenerated       EventType = "RecommendationGenerated"
	ImpactAnalysisCompleted       EventType = "ImpactAnalysisCompleted"
	ConfigurationChanged          EventType = "ConfigurationChanged"
	ManifestLoaded                EventType = "ManifestLoaded"
	ManifestValidationFailed      EventType = "ManifestValidationFailed"
	CommandStarted                EventType = "CommandStarted"
	CommandCompleted              EventType = "CommandCompleted"
	HooksExecuted                 EventType = "HooksExecuted"
)

// Event carries type and optional data payloads
type Event struct {
	Type    EventType
	Message string
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
