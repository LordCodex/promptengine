package scheduler

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// Job defines the function signature of scheduled CLI tasks.
type Job func(ctx context.Context) error

// Scheduler maintains a catalog of execution tasks that can be run on demand.
type Scheduler struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs: make(map[string]Job),
	}
}

// Register adds a new job callback.
func (s *Scheduler) Register(name string, job Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[name] = job
}

// Run executes a registered scheduled job by name.
func (s *Scheduler) Run(ctx context.Context, name string) error {
	s.mu.RLock()
	job, ok := s.jobs[name]
	s.mu.RUnlock()

	if !ok {
		return fmt.Errorf("scheduler: job '%s' not registered", name)
	}
	return job(ctx)
}

// List returns a sorted slice of registered scheduled job names.
func (s *Scheduler) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []string
	for k := range s.jobs {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}
