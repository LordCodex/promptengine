package eventbus

import (
	"sync"
	"testing"
	"time"
)

func TestEventBus_Publish_Sync(t *testing.T) {
	eb := NewEventBus()
	received := false

	eb.Subscribe(ProjectInitialized, func(e Event) {
		received = true
		if e.Payload.(string) != "test-payload" {
			t.Errorf("expected test-payload, got: %v", e.Payload)
		}
	})

	eb.Publish(Event{
		Type:    ProjectInitialized,
		Payload: "test-payload",
	})

	if !received {
		t.Error("expected handler to receive published event")
	}
}

func TestEventBus_Publish_Async(t *testing.T) {
	eb := NewEventBus()
	var wg sync.WaitGroup
	wg.Add(1)

	eb.Subscribe(WorkflowCompleted, func(e Event) {
		defer wg.Done()
		if e.Payload.(string) != "async-payload" {
			t.Errorf("expected async-payload, got: %v", e.Payload)
		}
	})

	eb.PublishAsync(Event{
		Type:    WorkflowCompleted,
		Payload: "async-payload",
	})

	// Wait with timeout to prevent blocking indefinitely if test fails
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for async subscriber handler to execute")
	}
}

func TestEventBus_NoSubscribers(t *testing.T) {
	eb := NewEventBus()
	// Should not block or panic when publishing an event with no subscribers
	eb.Publish(Event{
		Type:    ProjectDetected,
		Payload: "no-subscribers",
	})
	eb.PublishAsync(Event{
		Type:    ProjectDetected,
		Payload: "no-subscribers",
	})
}
