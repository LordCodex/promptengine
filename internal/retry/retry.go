package retry

import (
	"context"
	"time"
)

// Policy defines retry parameters
type Policy struct {
	MaxAttempts int
	Backoff     time.Duration
	MaxBackoff  time.Duration
	Factor      float64
}

// DefaultPolicy returns a typical retry policy
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Backoff:     100 * time.Millisecond,
		MaxBackoff:  2 * time.Second,
		Factor:      2.0,
	}
}

// Execute runs the operation according to the retry policy
func (p Policy) Execute(ctx context.Context, op func() error) error {
	var err error
	backoff := p.Backoff

	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		if err = op(); err == nil {
			return nil
		}

		if attempt == p.MaxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}

		backoff = time.Duration(float64(backoff) * p.Factor)
		if backoff > p.MaxBackoff {
			backoff = p.MaxBackoff
		}
	}
	return err
}
