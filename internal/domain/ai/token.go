package ai

import "strings"

func EstimateTokens(text string) int {
	if text == "" {
		return 0
	}
	tokens := len([]rune(text)) / 4
	if tokens == 0 {
		return 1
	}
	return tokens
}

func EstimateUsage(req Request, output string, costRank int) TokenUsage {
	input := EstimateTokens(strings.Join([]string{req.SystemInstructions, req.Context, req.Prompt}, "\n"))
	out := EstimateTokens(output)
	total := input + out
	return TokenUsage{
		InputTokens:   input,
		OutputTokens:  out,
		TotalTokens:   total,
		EstimatedCost: float64(total*costRank) / 1000000,
	}
}

func toAIError(err error, req Request) AIError {
	if aiErr, ok := err.(AIError); ok {
		return aiErr
	}
	return AIError{Category: ErrorProviderUnavailable, Provider: req.Provider, Model: req.Model, Message: err.Error(), RecommendedAction: "Check provider configuration and retry."}
}

func safeRequest(req Request) map[string]interface{} {
	return map[string]interface{}{
		"provider":    req.Provider,
		"model":       req.Model,
		"max_tokens":  req.MaxTokens,
		"temperature": req.Temperature,
		"metadata":    req.Metadata,
	}
}

func safeResponse(resp Response) map[string]interface{} {
	return map[string]interface{}{
		"provider":      resp.Provider,
		"model":         resp.Model,
		"token_usage":   resp.Usage,
		"finish_reason": resp.FinishReason,
		"metadata":      resp.Metadata,
	}
}

func safeError(err error) map[string]string {
	if aiErr, ok := err.(AIError); ok {
		return map[string]string{"category": string(aiErr.Category), "provider": aiErr.Provider, "model": aiErr.Model, "message": aiErr.Message}
	}
	return map[string]string{"message": err.Error()}
}
