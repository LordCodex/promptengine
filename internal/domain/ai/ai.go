package ai

type ProviderName string

const (
	ProviderChatGPT  ProviderName = "chatgpt"
	ProviderClaude   ProviderName = "claude"
	ProviderGemini   ProviderName = "gemini"
	ProviderCursor   ProviderName = "cursor"
	ProviderWindsurf ProviderName = "windsurf"
)

type PromptPayload struct {
	SystemPrompt string
	UserPrompt   string
	ContextFiles []string
}

// Adapter isolates provider-specific formats
type Adapter interface {
	Name() ProviderName
	Format(payload *PromptPayload) (string, error)
}

// StandardAIAdapter implements Adapter interface
type StandardAIAdapter struct {
	provider ProviderName
}

func NewStandardAIAdapter(p ProviderName) *StandardAIAdapter {
	return &StandardAIAdapter{provider: p}
}

func (a *StandardAIAdapter) Name() ProviderName {
	return a.provider
}

func (a *StandardAIAdapter) Format(payload *PromptPayload) (string, error) {
	// Abstract serializations logic
	switch a.provider {
	case ProviderCursor:
		return ".cursorrules standard format", nil
	case ProviderWindsurf:
		return ".windsurfrules format details", nil
	default:
		return payload.SystemPrompt + "\n" + payload.UserPrompt, nil
	}
}
