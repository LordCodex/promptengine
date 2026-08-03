package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type baseAdapter struct {
	id       string
	name     string
	envKey   string
	offline  bool
	models   []string
	caps     []Capability
	costRank int
}

func NewOpenAIAdapter() Provider {
	return &baseAdapter{id: "openai", name: "OpenAI", envKey: "OPENAI_API_KEY", models: []string{"gpt-4.1", "gpt-4.1-mini"}, caps: []Capability{CapStructuredOutput, CapJSONMode, CapLongContext}, costRank: 4}
}

func NewAnthropicAdapter() Provider {
	return &baseAdapter{id: "anthropic", name: "Anthropic", envKey: "ANTHROPIC_API_KEY", models: []string{"claude-sonnet-4", "claude-haiku-3.5"}, caps: []Capability{CapLongContext}, costRank: 4}
}

func NewGeminiAdapter() Provider {
	return &baseAdapter{id: "gemini", name: "Gemini", envKey: "GEMINI_API_KEY", models: []string{"gemini-2.5-pro", "gemini-2.5-flash"}, caps: []Capability{CapStructuredOutput, CapLongContext}, costRank: 3}
}

func NewLocalAdapter(id, name string) Provider {
	return &baseAdapter{id: id, name: name, offline: true, models: []string{"local-default"}, caps: []Capability{}, costRank: 1}
}

func (a *baseAdapter) ID() string { return a.id }

func (a *baseAdapter) Metadata() ProviderMetadata {
	return ProviderMetadata{ID: a.id, Name: a.name, ContextSize: a.Capabilities().MaxContextTokens, CostRank: a.costRank, SpeedRank: 3, IsOffline: a.offline}
}

func (a *baseAdapter) Capabilities() CapabilitySet {
	return CapabilitySet{Provider: a.id, Models: append([]string(nil), a.models...), Capabilities: append([]Capability(nil), a.caps...), MaxContextTokens: 128000, SupportsStreaming: hasCapabilities(CapabilitySet{Capabilities: a.caps}, []Capability{CapStreaming}), Offline: a.offline}
}

func (a *baseAdapter) ValidateConnection(ctx context.Context) error {
	if a.offline {
		return nil
	}
	if strings.TrimSpace(os.Getenv(a.envKey)) == "" {
		return AIError{Category: ErrorAuthentication, Provider: a.id, Message: "provider credentials are not configured", RecommendedAction: "Set the provider API key in the environment."}
	}
	return nil
}

func (a *baseAdapter) Generate(ctx context.Context, req Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, AIError{Category: ErrorTimeout, Provider: a.id, Model: req.Model, Message: err.Error(), RecommendedAction: "Retry with a longer timeout."}
	}
	if req.Model != "" && !containsString(a.models, req.Model) {
		return Response{}, AIError{Category: ErrorModelUnavailable, Provider: a.id, Model: req.Model, Message: "requested model is not available for provider", RecommendedAction: "Choose one of the provider capabilities models."}
	}
	if err := a.ValidateConnection(ctx); err != nil {
		return Response{}, err
	}
	model := req.Model
	if model == "" && len(a.models) > 0 {
		model = a.models[0]
	}
	content, err := a.generateHTTP(ctx, model, req)
	if err != nil {
		return Response{}, err
	}
	usage := EstimateUsage(req, content, a.costRank)
	return Response{Content: content, Model: model, Provider: a.id, Usage: usage, FinishReason: "stop", Metadata: map[string]string{"adapter": a.id}}, nil
}

func (a *baseAdapter) Stream(ctx context.Context, req Request) (Stream, error) {
	if !hasCapabilities(a.Capabilities(), []Capability{CapStreaming}) {
		return nil, AIError{Category: ErrorProviderUnavailable, Provider: a.id, Message: "provider does not support streaming"}
	}
	out := make(chan StreamChunk, 2)
	go func() {
		defer close(out)
		resp, err := a.Generate(ctx, req)
		if err != nil {
			out <- StreamChunk{Err: err, Done: true}
			return
		}
		if resp.Content != "" {
			out <- StreamChunk{Content: resp.Content}
		}
		out <- StreamChunk{Done: true, Usage: resp.Usage}
	}()
	return out, nil
}

func (a *baseAdapter) generateHTTP(ctx context.Context, model string, req Request) (string, error) {
	switch a.id {
	case "openai":
		return postOpenAI(ctx, os.Getenv(a.envKey), model, req)
	case "anthropic":
		return postAnthropic(ctx, os.Getenv(a.envKey), model, req)
	case "gemini":
		return postGemini(ctx, os.Getenv(a.envKey), model, req)
	default:
		return postOllama(ctx, model, req)
	}
}

func postOpenAI(ctx context.Context, key, model string, req Request) (string, error) {
	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": req.SystemInstructions},
			{"role": "user", "content": strings.TrimSpace(req.Context + "\n\n" + req.Prompt)},
		},
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := postJSON(ctx, "https://api.openai.com/v1/chat/completions", key, body, &out, nil); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", AIError{Category: ErrorInvalidResponse, Provider: "openai", Model: model, Message: "provider returned no choices"}
	}
	return out.Choices[0].Message.Content, nil
}

func postAnthropic(ctx context.Context, key, model string, req Request) (string, error) {
	body := map[string]interface{}{
		"model":       model,
		"system":      req.SystemInstructions,
		"max_tokens":  firstPositive(req.MaxTokens, 1024),
		"temperature": req.Temperature,
		"messages": []map[string]string{
			{"role": "user", "content": strings.TrimSpace(req.Context + "\n\n" + req.Prompt)},
		},
	}
	var out struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	headers := map[string]string{"anthropic-version": "2023-06-01"}
	if err := postJSON(ctx, "https://api.anthropic.com/v1/messages", key, body, &out, headers); err != nil {
		return "", err
	}
	if len(out.Content) == 0 {
		return "", AIError{Category: ErrorInvalidResponse, Provider: "anthropic", Model: model, Message: "provider returned no content"}
	}
	return out.Content[0].Text, nil
}

func postGemini(ctx context.Context, key, model string, req Request) (string, error) {
	body := map[string]interface{}{
		"contents": []map[string]interface{}{{"parts": []map[string]string{{"text": strings.TrimSpace(req.SystemInstructions + "\n\n" + req.Context + "\n\n" + req.Prompt)}}}},
	}
	var out struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, key)
	if err := postJSON(ctx, url, "", body, &out, nil); err != nil {
		return "", err
	}
	if len(out.Candidates) == 0 || len(out.Candidates[0].Content.Parts) == 0 {
		return "", AIError{Category: ErrorInvalidResponse, Provider: "gemini", Model: model, Message: "provider returned no candidates"}
	}
	return out.Candidates[0].Content.Parts[0].Text, nil
}

func postOllama(ctx context.Context, model string, req Request) (string, error) {
	host := strings.TrimRight(firstNonEmptyString(os.Getenv("OLLAMA_HOST"), "http://127.0.0.1:11434"), "/")
	body := map[string]interface{}{
		"model":  model,
		"prompt": strings.TrimSpace(req.SystemInstructions + "\n\n" + req.Context + "\n\n" + req.Prompt),
		"stream": false,
	}
	var out struct {
		Response string `json:"response"`
	}
	if err := postJSON(ctx, host+"/api/generate", "", body, &out, nil); err != nil {
		return "", AIError{Category: ErrorProviderUnavailable, Provider: "ollama", Model: model, Message: err.Error(), RecommendedAction: "Start a compatible local model server or choose prompt export instead."}
	}
	if out.Response == "" {
		return "", AIError{Category: ErrorInvalidResponse, Provider: "ollama", Model: model, Message: "local model returned an empty response"}
	}
	return out.Response, nil
}

func postJSON(ctx context.Context, url, bearer string, body interface{}, target interface{}, headers map[string]string) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	httpReq.Header.Set("content-type", "application/json")
	if bearer != "" {
		httpReq.Header.Set("authorization", "Bearer "+bearer)
	}
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 400 {
		return AIError{Category: ErrorProviderUnavailable, Message: fmt.Sprintf("provider returned HTTP %d", resp.StatusCode), RecommendedAction: strings.TrimSpace(string(raw))}
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return AIError{Category: ErrorInvalidResponse, Message: err.Error()}
	}
	return nil
}

func firstPositive(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func (a *baseAdapter) Execute(ctx context.Context, req ExecutionRequest) (ExecutionResponse, error) {
	resp, err := a.Generate(ctx, Request{SystemInstructions: req.SystemPrompt, Prompt: req.UserPrompt, Temperature: req.Temperature})
	return ExecutionResponse{Text: resp.Content, EstimatedTokens: resp.Usage.TotalTokens, CapabilityOutput: resp.Metadata}, err
}

func containsString(items []string, expected string) bool {
	for _, item := range items {
		if item == expected {
			return true
		}
	}
	return false
}
