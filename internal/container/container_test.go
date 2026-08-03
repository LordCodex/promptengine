package container

import (
	"bytes"
	"testing"

	"github.com/LordCodex/promptengine/internal/config"
)

func TestNewContainer_InitializesPhaseOneDependencies(t *testing.T) {
	cfg := config.DefaultConfig()
	var out bytes.Buffer
	var errOut bytes.Buffer

	c, err := NewContainer(Options{Config: cfg, Out: &out, Err: &errOut})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if c.Config == nil || c.Logger == nil || c.FS == nil || c.Cache == nil || c.Renderer == nil || c.EventBus == nil {
		t.Fatal("expected all phase one dependencies to be initialized")
	}
}
