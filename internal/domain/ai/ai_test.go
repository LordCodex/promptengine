package ai

import (
	"context"
	"testing"
)

// MockProvider stub
type MockProvider struct {
	id     string
	meta   ProviderMetadata
	caps   map[Capability]bool
	result string
}

func (m *MockProvider) ID() string                       { return m.id }
func (m *MockProvider) Metadata() ProviderMetadata       { return m.meta }
func (m *MockProvider) Capabilities() map[Capability]bool { return m.caps }
func (m *MockProvider) Execute(ctx context.Context, req ExecutionRequest) (ExecutionResponse, error) {
	return ExecutionResponse{Text: m.result}, nil
}
func (m *MockProvider) Stream(ctx context.Context, req ExecutionRequest) (<-chan string, <-chan error, error) {
	tChan := make(chan string, 2)
	eChan := make(chan error, 1)
	tChan <- "chunk1"
	tChan <- "chunk2"
	close(tChan)
	close(eChan)
	return tChan, eChan, nil
}

func TestAIProvider_Selection(t *testing.T) {
	reg := NewRegistry()

	onlineModel := &MockProvider{
		id:   "gemini-ultra",
		meta: ProviderMetadata{ID: "gemini-ultra", IsOffline: false},
		caps: map[Capability]bool{CapStreaming: true, CapStructuredOutput: true},
	}
	offlineModel := &MockProvider{
		id:   "ollama-llama3",
		meta: ProviderMetadata{ID: "ollama-llama3", IsOffline: true},
		caps: map[Capability]bool{CapStreaming: true},
	}

	reg.Register(onlineModel)
	reg.Register(offlineModel)

	selector := NewSelector(reg)

	// Test 1: Choose offline only model
	prov, err := selector.SelectBest("", []Capability{CapStreaming}, true)
	if err != nil {
		t.Fatalf("Expected model match, got error: %v", err)
	}
	if prov.ID() != "ollama-llama3" {
		t.Errorf("Expected ollama-llama3, got %s", prov.ID())
	}

	// Test 2: Choose matching capabilities model
	prov, err = selector.SelectBest("", []Capability{CapStructuredOutput}, false)
	if err != nil {
		t.Fatalf("Expected model match, got error: %v", err)
	}
	if prov.ID() != "gemini-ultra" {
		t.Errorf("Expected gemini-ultra, got %s", prov.ID())
	}
}

func TestAIProvider_ConversationTracking(t *testing.T) {
	c := NewConversation("You are an assistant.")
	c.Append(RoleUser, "Hello")
	c.Append(RoleAssistant, "Hi there")

	tokens := c.EstimateTokens()
	if tokens <= 0 {
		t.Errorf("Expected positive token count estimation, got %d", tokens)
	}

	c.Reset("New Prompt")
	if len(c.Messages) != 1 || c.Messages[0].Content != "New Prompt" {
		t.Errorf("Expected reset message structure to reset messages list")
	}
}
