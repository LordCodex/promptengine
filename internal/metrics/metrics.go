package metrics

import "time"

// Metrics defines the interface for recording application telemetry and metrics.
type Metrics interface {
	Counter(name string, value int64, tags map[string]string)
	Gauge(name string, value float64, tags map[string]string)
	Timer(name string, duration time.Duration, tags map[string]string)
}

// MockMetrics is a simple implementation of Metrics for testing
type MockMetrics struct {
	Counters map[string]int64
	Gauges   map[string]float64
	Timers   map[string]time.Duration
}

func NewMockMetrics() *MockMetrics {
	return &MockMetrics{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
		Timers:   make(map[string]time.Duration),
	}
}

func (m *MockMetrics) Counter(name string, value int64, tags map[string]string) {
	m.Counters[name] += value
}

func (m *MockMetrics) Gauge(name string, value float64, tags map[string]string) {
	m.Gauges[name] = value
}

func (m *MockMetrics) Timer(name string, duration time.Duration, tags map[string]string) {
	m.Timers[name] = duration
}
