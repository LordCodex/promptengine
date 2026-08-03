package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryPolicy_ExecuteSuccess(t *testing.T) {
	p := Policy{
		MaxAttempts: 3,
		Backoff:     1 * time.Millisecond,
		MaxBackoff:  10 * time.Millisecond,
		Factor:      2.0,
	}

	attempts := 0
	err := p.Execute(context.Background(), func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected success, got: %v", err)
	}
	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got: %d", attempts)
	}
}

func TestRetryPolicy_ExecuteFailure(t *testing.T) {
	p := Policy{
		MaxAttempts: 3,
		Backoff:     1 * time.Millisecond,
		MaxBackoff:  10 * time.Millisecond,
		Factor:      2.0,
	}

	attempts := 0
	err := p.Execute(context.Background(), func() error {
		attempts++
		return errors.New("permanent error")
	})

	if err == nil || err.Error() != "permanent error" {
		t.Errorf("Expected permanent error, got: %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got: %d", attempts)
	}
}
