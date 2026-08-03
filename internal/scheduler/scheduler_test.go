package scheduler

import (
	"context"
	"errors"
	"testing"
)

func TestScheduler(t *testing.T) {
	s := NewScheduler()

	var called bool
	s.Register("test-job", func(ctx context.Context) error {
		called = true
		return nil
	})

	s.Register("error-job", func(ctx context.Context) error {
		return errors.New("job failed")
	})

	list := s.List()
	if len(list) != 2 || list[0] != "error-job" || list[1] != "test-job" {
		t.Errorf("Unexpected job list: %v", list)
	}

	err := s.Run(context.Background(), "test-job")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !called {
		t.Error("Expected test-job to be executed")
	}

	err = s.Run(context.Background(), "error-job")
	if err == nil || err.Error() != "job failed" {
		t.Errorf("Expected job failure error, got: %v", err)
	}

	err = s.Run(context.Background(), "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent job, got nil")
	}
}
