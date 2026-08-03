package ai

import "context"

// StreamingCoordinator manages console feeds from streaming models
type StreamingCoordinator struct{}

func NewStreamingCoordinator() *StreamingCoordinator {
	return &StreamingCoordinator{}
}

func (s *StreamingCoordinator) StreamToWriter(ctx context.Context, p Provider, req ExecutionRequest, chunkCallback func(string)) (string, error) {
	textChan, errChan, err := p.Stream(ctx, req)
	if err != nil {
		return "", err
	}

	fullResponse := ""
	for {
		select {
		case <-ctx.Done():
			return fullResponse, ctx.Err()
		case chunk, ok := <-textChan:
			if !ok {
				return fullResponse, nil
			}
			fullResponse += chunk
			if chunkCallback != nil {
				chunkCallback(chunk)
			}
		case err, ok := <-errChan:
			if ok && err != nil {
				return fullResponse, err
			}
		}
	}
}
