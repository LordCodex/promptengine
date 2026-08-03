package perf

import (
	"log/slog"
	"time"
)

// ExecutionTimer measures elapsed execution duration
type ExecutionTimer struct {
	name  string
	start time.Time
}

// Start initiates a named execution timer
func Start(name string) *ExecutionTimer {
	return &ExecutionTimer{
		name:  name,
		start: time.Now(),
	}
}

// Log prints details of elapsed duration to slog
func (t *ExecutionTimer) Log(logger *slog.Logger) {
	elapsed := time.Since(t.start)
	logger.Debug("Performance benchmark metric", "component", t.name, "duration_ms", elapsed.Milliseconds())
}

// Duration returns elapsed time duration
func (t *ExecutionTimer) Duration() time.Duration {
	return time.Since(t.start)
}
