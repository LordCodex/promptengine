package ai

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	ctxengine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/eventbus"
)

type MockProvider struct {
	id     string
	meta   ProviderMetadata
	caps   CapabilitySet
	result string
	err    error
}

func (m *MockProvider) ID() string                                   { return m.id }
func (m *MockProvider) Metadata() ProviderMetadata                   { return m.meta }
func (m *MockProvider) Capabilities() CapabilitySet                  { return m.caps }
func (m *MockProvider) ValidateConnection(ctx context.Context) error { return m.err }
func (m *MockProvider) Generate(ctx context.Context, req Request) (Response, error) {
	if m.err != nil {
		return Response{}, m.err
	}
	return Response{Content: m.result, Model: req.Model, Provider: m.id, Usage: EstimateUsage(req, m.result, m.meta.CostRank), FinishReason: "stop"}, nil
}
func (m *MockProvider) Stream(ctx context.Context, req Request) (Stream, error) {
	if m.err != nil {
		return nil, m.err
	}
	ch := make(chan StreamChunk, 3)
	ch <- StreamChunk{Content: "chunk1"}
	ch <- StreamChunk{Content: "chunk2"}
	ch <- StreamChunk{Done: true}
	close(ch)
	return ch, nil
}
func (m *MockProvider) Execute(ctx context.Context, req ExecutionRequest) (ExecutionResponse, error) {
	resp, err := m.Generate(ctx, Request{SystemInstructions: req.SystemPrompt, Prompt: req.UserPrompt})
	return ExecutionResponse{Text: resp.Content, EstimatedTokens: resp.Usage.TotalTokens}, err
}

func TestAIProvider_Selection(t *testing.T) {
	reg := &Registry{providers: map[string]Provider{}}
	reg.Register(&MockProvider{id: "gemini-ultra", meta: ProviderMetadata{ID: "gemini-ultra", IsOffline: false}, caps: CapabilitySet{Provider: "gemini-ultra", Capabilities: []Capability{CapStreaming, CapStructuredOutput}}})
	reg.Register(&MockProvider{id: "ollama-llama3", meta: ProviderMetadata{ID: "ollama-llama3", IsOffline: true}, caps: CapabilitySet{Provider: "ollama-llama3", Capabilities: []Capability{CapStreaming}, Offline: true}})

	selector := NewSelector(reg)
	prov, err := selector.SelectBest("", []Capability{CapStreaming}, true)
	if err != nil {
		t.Fatalf("Expected model match, got error: %v", err)
	}
	if prov.ID() != "ollama-llama3" {
		t.Errorf("Expected ollama-llama3, got %s", prov.ID())
	}

	prov, err = selector.SelectBest("", []Capability{CapStructuredOutput}, false)
	if err != nil {
		t.Fatalf("Expected model match, got error: %v", err)
	}
	if prov.ID() != "gemini-ultra" {
		t.Errorf("Expected gemini-ultra, got %s", prov.ID())
	}
}

func TestPromptCompiler_CompilesFromContext(t *testing.T) {
	pkg := ctxengine.NewContextPackage("feature", ctxengine.BudgetSmall)
	pkg.Items = append(pkg.Items, ctxengine.ContextItem{Path: "docs/API.md", Type: ctxengine.ContextDocumentation, Content: "API docs"})
	req, err := NewPromptCompiler().Compile(CompileInput{
		Provider: "mock", Model: "mock-1", UserRequest: "Build feature", SystemInstructions: "Follow standards", WorkflowRequirements: []string{"validate"}, Standards: []string{"security"}, ContextPackage: pkg,
	})
	if err != nil {
		t.Fatalf("expected compile success, got %v", err)
	}
	if !strings.Contains(req.Context, "docs/API.md") || !strings.Contains(req.Context, "security") {
		t.Fatalf("expected context and standards in prompt, got %s", req.Context)
	}
}

func TestPlatform_GenerateTracksUsageAndEvents(t *testing.T) {
	reg := &Registry{providers: map[string]Provider{}}
	reg.Register(&MockProvider{id: "mock", meta: ProviderMetadata{ID: "mock", CostRank: 2}, caps: CapabilitySet{Provider: "mock"}, result: "hello"})
	events := eventbus.NewEventBus()
	var started, completed bool
	events.Subscribe(eventbus.AIRequestStarted, func(e eventbus.Event) { started = true })
	events.Subscribe(eventbus.AIRequestCompleted, func(e eventbus.Event) { completed = true })

	resp, err := NewPlatform(reg, events).Generate(context.Background(), Request{Provider: "mock", Model: "mock-1", Prompt: "hello"})
	if err != nil {
		t.Fatalf("expected generate success, got %v", err)
	}
	if resp.Usage.TotalTokens == 0 || resp.ExecutionDuration < 0 {
		t.Fatalf("expected usage and duration, got %#v", resp)
	}
	if !started || !completed {
		t.Fatalf("expected started/completed events")
	}
}

func TestProviderFailureStructured(t *testing.T) {
	reg := &Registry{providers: map[string]Provider{}}
	reg.Register(&MockProvider{id: "mock", meta: ProviderMetadata{ID: "mock"}, caps: CapabilitySet{Provider: "mock"}, err: AIError{Category: ErrorRateLimit, Provider: "mock", Message: "limited"}})
	events := eventbus.NewEventBus()
	var failed bool
	events.Subscribe(eventbus.AIRequestFailed, func(e eventbus.Event) { failed = true })

	resp, err := NewPlatform(reg, events).Generate(context.Background(), Request{Provider: "mock", Prompt: "hello"})
	if err == nil {
		t.Fatal("expected provider failure")
	}
	if len(resp.Errors) != 1 || resp.Errors[0].Category != ErrorRateLimit {
		t.Fatalf("expected structured rate-limit error, got %#v", resp.Errors)
	}
	if !failed {
		t.Fatal("expected failure event")
	}
}

func TestStreamingBehavior(t *testing.T) {
	prov := &MockProvider{id: "mock", meta: ProviderMetadata{ID: "mock"}, caps: CapabilitySet{Provider: "mock", Capabilities: []Capability{CapStreaming}}, result: "ok"}
	got, err := NewStreamingCoordinator().StreamToWriter(context.Background(), prov, ExecutionRequest{UserPrompt: "hello"}, nil)
	if err != nil {
		t.Fatalf("expected stream success, got %v", err)
	}
	if got != "chunk1chunk2" {
		t.Fatalf("unexpected stream output %q", got)
	}
}

func TestAdaptersValidateWithoutLeakingSecrets(t *testing.T) {
	err := NewOpenAIAdapter().ValidateConnection(context.Background())
	if err == nil {
		t.Fatal("expected missing key error")
	}
	if strings.Contains(err.Error(), "OPENAI_API_KEY") {
		t.Fatal("error should not leak secret variable details")
	}
}

func TestConversationTracking(t *testing.T) {
	c := NewConversation("You are an assistant.")
	c.Append(RoleUser, "Hello")
	c.Append(RoleAssistant, "Hi there")
	if c.EstimateTokens() <= 0 {
		t.Fatal("expected positive token estimate")
	}
	c.Reset("New Prompt")
	if len(c.Messages) != 1 || c.Messages[0].Content != "New Prompt" {
		t.Fatal("expected reset message")
	}
}

func TestPlatform_CompilePublishesPromptGenerated(t *testing.T) {
	events := eventbus.NewEventBus()
	var generated bool
	events.Subscribe(eventbus.PromptGenerated, func(e eventbus.Event) { generated = true })
	_, err := NewPlatform(&Registry{providers: map[string]Provider{}}, events).Compile(CompileInput{UserRequest: "Review code"})
	if err != nil {
		t.Fatalf("expected compile success, got %v", err)
	}
	if !generated {
		t.Fatal("expected PromptGenerated event")
	}
}

func TestTimeoutErrorMapping(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	_, err := NewLocalAdapter("local", "Local").Generate(ctx, Request{Prompt: "hello"})
	if err == nil {
		t.Fatal("expected timeout error")
	}
	var aiErr AIError
	if !errors.As(err, &aiErr) || aiErr.Category != ErrorTimeout {
		t.Fatalf("expected timeout AIError, got %v", err)
	}
}
