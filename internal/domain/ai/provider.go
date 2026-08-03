package ai

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/LordCodex/promptengine/internal/eventbus"
)

type Request struct {
	Provider           string            `json:"provider,omitempty"`
	Model              string            `json:"model,omitempty"`
	Prompt             string            `json:"prompt"`
	Context            string            `json:"context,omitempty"`
	SystemInstructions string            `json:"system_instructions,omitempty"`
	Temperature        float32           `json:"temperature,omitempty"`
	MaxTokens          int               `json:"max_tokens,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
}

type Response struct {
	Content           string            `json:"generated_content"`
	Model             string            `json:"model"`
	Provider          string            `json:"provider"`
	Usage             TokenUsage        `json:"token_usage"`
	ExecutionDuration time.Duration     `json:"execution_duration"`
	FinishReason      string            `json:"finish_reason"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	Errors            []AIError         `json:"errors,omitempty"`
}

type TokenUsage struct {
	InputTokens   int     `json:"input_tokens"`
	OutputTokens  int     `json:"output_tokens"`
	TotalTokens   int     `json:"total_tokens"`
	EstimatedCost float64 `json:"estimated_cost"`
}

type StreamChunk struct {
	Content string     `json:"content"`
	Done    bool       `json:"done"`
	Usage   TokenUsage `json:"usage,omitempty"`
	Err     error      `json:"-"`
}

type Stream <-chan StreamChunk

type CapabilitySet struct {
	Provider          string       `json:"provider"`
	Models            []string     `json:"models"`
	Capabilities      []Capability `json:"capabilities"`
	MaxContextTokens  int          `json:"max_context_tokens"`
	SupportsStreaming bool         `json:"supports_streaming"`
	Offline           bool         `json:"offline"`
}

type ErrorCategory string

const (
	ErrorProviderUnavailable ErrorCategory = "provider_unavailable"
	ErrorAuthentication      ErrorCategory = "authentication_failure"
	ErrorRateLimit           ErrorCategory = "rate_limit"
	ErrorTimeout             ErrorCategory = "timeout"
	ErrorInvalidResponse     ErrorCategory = "invalid_response"
	ErrorModelUnavailable    ErrorCategory = "model_unavailable"
)

type AIError struct {
	Category          ErrorCategory `json:"category"`
	Provider          string        `json:"provider,omitempty"`
	Model             string        `json:"model,omitempty"`
	Message           string        `json:"message"`
	RecommendedAction string        `json:"recommended_action,omitempty"`
}

func (e AIError) Error() string { return e.Message }

type Provider interface {
	ID() string
	Generate(ctx context.Context, request Request) (Response, error)
	Stream(ctx context.Context, request Request) (Stream, error)
	ValidateConnection(ctx context.Context) error
	Capabilities() CapabilitySet

	// Legacy compatibility.
	Metadata() ProviderMetadata
	Execute(ctx context.Context, req ExecutionRequest) (ExecutionResponse, error)
}

type ExecutionRequest struct {
	SystemPrompt string
	UserPrompt   string
	Temperature  float32
	JSONMode     bool
}

type ExecutionResponse struct {
	Text             string
	EstimatedTokens  int
	CapabilityOutput interface{}
}

type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

func NewRegistry() *Registry {
	r := &Registry{providers: map[string]Provider{}}
	r.Register(NewOpenAIAdapter())
	r.Register(NewAnthropicAdapter())
	r.Register(NewGeminiAdapter())
	r.Register(NewLocalAdapter("ollama", "Ollama"))
	return r
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
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

func (r *Registry) Discover() []CapabilitySet {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []CapabilitySet
	for _, p := range r.providers {
		out = append(out, p.Capabilities())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Provider < out[j].Provider })
	return out
}

type Selector struct{ registry *Registry }

func NewSelector(r *Registry) *Selector { return &Selector{registry: r} }

func (s *Selector) SelectBest(preferredID string, requiredCaps []Capability, offlineOnly bool) (Provider, error) {
	if preferredID != "" {
		if p, ok := s.registry.Get(preferredID); ok {
			return p, nil
		}
	}
	s.registry.mu.RLock()
	defer s.registry.mu.RUnlock()
	for _, p := range s.registry.providers {
		caps := p.Capabilities()
		if offlineOnly && !caps.Offline {
			continue
		}
		if hasCapabilities(caps, requiredCaps) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no provider found satisfying required capabilities")
}

type Platform struct {
	registry *Registry
	compiler *PromptCompiler
	events   *eventbus.EventBus
}

func NewPlatform(registry *Registry, events *eventbus.EventBus) *Platform {
	if registry == nil {
		registry = NewRegistry()
	}
	return &Platform{registry: registry, compiler: NewPromptCompiler(), events: events}
}

func (p *Platform) Registry() *Registry { return p.registry }

func (p *Platform) Generate(ctx context.Context, req Request) (Response, error) {
	provider, err := p.selectProvider(req)
	if err != nil {
		p.publish(eventbus.AIRequestFailed, "ai request failed", safeError(err))
		return Response{Errors: []AIError{toAIError(err, req)}}, err
	}
	req.Provider = provider.ID()
	p.publish(eventbus.AIRequestStarted, "ai request started", safeRequest(req))
	start := time.Now()
	resp, err := provider.Generate(ctx, req)
	resp.ExecutionDuration = time.Since(start)
	if resp.Usage.TotalTokens == 0 {
		resp.Usage = EstimateUsage(req, resp.Content, provider.Metadata().CostRank)
	}
	if err != nil {
		resp.Errors = append(resp.Errors, toAIError(err, req))
		p.publish(eventbus.AIRequestFailed, "ai request failed", safeError(err))
		return resp, err
	}
	p.publish(eventbus.AIRequestCompleted, "ai request completed", safeResponse(resp))
	return resp, nil
}

func (p *Platform) Stream(ctx context.Context, req Request) (Stream, error) {
	provider, err := p.selectProvider(req)
	if err != nil {
		p.publish(eventbus.AIRequestFailed, "ai request failed", safeError(err))
		return nil, err
	}
	req.Provider = provider.ID()
	p.publish(eventbus.AIStreamingStarted, "ai streaming started", safeRequest(req))
	return provider.Stream(ctx, req)
}

func (p *Platform) Compile(input CompileInput) (Request, error) {
	req, err := p.compiler.Compile(input)
	if err == nil {
		p.publish(eventbus.PromptGenerated, "prompt generated", safeRequest(req))
	}
	return req, err
}

func (p *Platform) selectProvider(req Request) (Provider, error) {
	if req.Provider != "" {
		provider, ok := p.registry.Get(req.Provider)
		if !ok {
			return nil, AIError{Category: ErrorProviderUnavailable, Provider: req.Provider, Message: "requested provider is not registered", RecommendedAction: "Register or configure the provider before use."}
		}
		return provider, nil
	}
	return NewSelector(p.registry).SelectBest("", nil, false)
}

func (p *Platform) publish(t eventbus.EventType, msg string, payload interface{}) {
	if p.events != nil {
		p.events.Publish(eventbus.Event{Type: t, Message: msg, Payload: payload})
	}
}

func hasCapabilities(set CapabilitySet, required []Capability) bool {
	available := map[Capability]bool{}
	for _, cap := range set.Capabilities {
		available[cap] = true
	}
	for _, cap := range required {
		if !available[cap] {
			return false
		}
	}
	return true
}
