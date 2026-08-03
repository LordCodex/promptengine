package ai

import "context"

type StreamingCoordinator struct{}

func NewStreamingCoordinator() *StreamingCoordinator { return &StreamingCoordinator{} }

func (s *StreamingCoordinator) StreamToWriter(ctx context.Context, p Provider, req ExecutionRequest, chunkCallback func(string)) (string, error) {
	stream, err := p.Stream(ctx, Request{SystemInstructions: req.SystemPrompt, Prompt: req.UserPrompt, Temperature: req.Temperature})
	if err != nil {
		return "", err
	}
	full := ""
	for chunk := range stream {
		if chunk.Err != nil {
			return full, chunk.Err
		}
		full += chunk.Content
		if chunk.Content != "" && chunkCallback != nil {
			chunkCallback(chunk.Content)
		}
		if chunk.Done {
			return full, nil
		}
	}
	return full, nil
}
