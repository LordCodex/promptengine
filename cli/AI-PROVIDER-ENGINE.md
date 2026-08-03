# AI Provider Engine Architecture Specification

The AI Provider Engine serves as the exclusive communication gateway between PromptEngine and LLM providers. No other subsystem should connect directly to an AI API. It isolates provider-specific schemas behind stable execution interfaces.

---

## 1. Capability System

Instead of hardcoding provider mappings, each provider defines a map of active `Capability` flags:
- **`CapStreaming`**: Supports real-time text segments rendering.
- **`CapToolCalling`**: Supports external function routing integrations (git, filesystem audits).
- **`CapStructuredOutput`**: Supports returning strict JSON mapping schemas.
- **`CapLongContext`**: Handles token windows greater than 100k.

PromptEngine automatically routes tasks based on capabilities: e.g., if a workflow requests JSON outputs, the selector chooses a model offering `CapStructuredOutput` or `CapJSONMode`.

---

## 2. Dynamic Provider Registry

Providers must compile to the common `Provider` interface and register themselves upon boot:

```go
type Provider interface {
    ID() string
    Metadata() ProviderMetadata
    Capabilities() map[Capability]bool
    Execute(ctx context.Context, req ExecutionRequest) (ExecutionResponse, error)
    Stream(ctx context.Context, req ExecutionRequest) (<-chan string, <-chan error, error)
}
```

---

## 3. Intelligent Selector Strategy

The selector picks the best match based on:
1. **Manual User Overrides**: Explicit `preferredID` parameters configuration.
2. **Task Capability Checks**: Validating that required capability flags are supported.
3. **Offline Constraints**: Forcing local LLM providers (e.g. Ollama local models) when offline mode is configured.

---

## 4. Conversation & Token Budgets

- **Conversation Tracker**: Models multi-turn memory buffers. `Conversation.EstimateTokens()` computes the estimated input token volume dynamically.
- **Provider Security**: Securely parses API parameters from configuration loaders. Integrations enforce that credentials are never logged, and files are sanitized before transmission.
