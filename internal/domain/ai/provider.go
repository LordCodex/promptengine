package ai

import (
	"context"
	"fmt"
	"sync"
)

// ExecutionRequest holds prompts inputs
type ExecutionRequest struct {
	SystemPrompt string
	UserPrompt   string
	Temperature  float32
	JSONMode     bool
}

// ExecutionResponse represents raw or structured outcome
type ExecutionResponse struct {
	Text             string
	EstimatedTokens  int
	CapabilityOutput interface{}
}

// Provider defines the execution contract for any LLM API
type Provider interface {
	ID() string
	Metadata() ProviderMetadata
	Capabilities() map[Capability]bool
	Execute(ctx context.Context, req ExecutionRequest) (ExecutionResponse, error)
	Stream(ctx context.Context, req ExecutionRequest) (<-chan string, <-chan error, error)
}

// Registry manages concurrency-safe provider stubs lookup
type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

func (r *Registry) Register(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.ID()] = p
}

func (r *Registry) Get(id string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[id]
	return p, ok
}

func (r *Registry) List() []ProviderMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []ProviderMetadata
	for _, p := range r.providers {
		list = append(list, p.Metadata())
	}
	return list
}

// Selector manages automatic matches based on task constraints
type Selector struct {
	registry *Registry
}

func NewSelector(r *Registry) *Selector {
	return &Selector{registry: r}
}

func (s *Selector) SelectBest(preferredID string, requiredCaps []Capability, offlineOnly bool) (Provider, error) {
	// 1. Explicit override matches first
	if preferredID != "" {
		if p, ok := s.registry.Get(preferredID); ok {
			return p, nil
		}
	}

	// 2. Iterate and match capabilities
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()

	for _, p := range s.registry.providers {
		if offlineOnly && !p.Metadata().IsOffline {
			continue
		}

		matched := true
		for _, cap := range requiredCaps {
			if !p.Capabilities()[cap] {
				matched = false
				break
			}
		}

		if matched {
			return p, nil
		}
	}

	return nil, fmt.Errorf("no provider found satisfying required capabilities")
}
