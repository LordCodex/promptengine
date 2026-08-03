package app

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/internal/history"
	"github.com/LordCodex/promptengine/pkg/manifest"
	"github.com/spf13/cobra"
)

func TestLifecycle_PublishesCommandEvents(t *testing.T) {
	var out bytes.Buffer
	a, err := Bootstrap(BootstrapOptions{Out: &out, Err: &out})
	if err != nil {
		t.Fatalf("expected bootstrap to succeed, got %v", err)
	}
	a.History = history.NewRecorder(filesystem.NewMockFileSystem())

	var events []eventbus.EventType
	a.EventBus.Subscribe(eventbus.CommandStarted, func(e eventbus.Event) {
		events = append(events, e.Type)
	})
	a.EventBus.Subscribe(eventbus.CommandCompleted, func(e eventbus.Event) {
		events = append(events, e.Type)
	})

	wrapped := a.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		if lc.Config == nil || lc.FS == nil || lc.Logger == nil {
			t.Fatal("expected lifecycle context dependencies")
		}
		return nil
	})

	cmd := &cobra.Command{Use: "test"}
	cmd.SetContext(context.Background())
	if err := wrapped(cmd, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(events) != 2 || events[0] != eventbus.CommandStarted || events[1] != eventbus.CommandCompleted {
		t.Fatalf("expected command started/completed events, got %#v", events)
	}
}

func TestLifecycle_RecordsAuditHistoryWithoutArgs(t *testing.T) {
	var out bytes.Buffer
	a, err := Bootstrap(BootstrapOptions{Out: &out, Err: &out})
	if err != nil {
		t.Fatalf("expected bootstrap to succeed, got %v", err)
	}
	fs := filesystem.NewMockFileSystem()
	a.History = history.NewRecorder(fs)

	wrapped := a.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		return nil
	})

	cmd := &cobra.Command{Use: "task"}
	cmd.SetContext(context.Background())
	if err := wrapped(cmd, []string{"api_key=sk_live_secret"}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	data, err := fs.ReadFile(history.DefaultPath)
	if err != nil {
		t.Fatalf("expected history file: %v", err)
	}
	if strings.Contains(string(data), "sk_live_secret") || strings.Contains(string(data), "api_key=") {
		t.Fatalf("history should not contain raw args or secrets: %s", data)
	}
	if !strings.Contains(string(data), `"command":"task"`) {
		t.Fatalf("expected command name in history: %s", data)
	}
}

func TestLifecycle_PublishesManifestLoaded(t *testing.T) {
	var out bytes.Buffer
	a, err := Bootstrap(BootstrapOptions{Out: &out, Err: &out})
	if err != nil {
		t.Fatalf("expected bootstrap to succeed, got %v", err)
	}
	a.History = history.NewRecorder(filesystem.NewMockFileSystem())
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("playbook-manifest.json", []byte(`{
		"metadata": {"name": "Example", "version": "1.0.0", "schema_version": "2.0.0"},
		"playbooks": [{"id": "p", "name": "P", "category": "core", "location": "p.md", "priority": 1}],
		"workflows": [{"id": "w", "steps": ["s"], "required_playbooks": ["p"]}]
	}`), 0644)
	fs.WriteFile("p.md", []byte("playbook"), 0644)
	a.FS = fs
	a.Manifest = manifest.NewEngineWithFS(fs)
	a.ManifestQ = manifest.NewQueryEngine(a.Manifest)

	loaded := false
	a.EventBus.Subscribe(eventbus.ManifestLoaded, func(e eventbus.Event) {
		loaded = true
	})

	cmd := &cobra.Command{Use: "test"}
	cmd.SetContext(context.Background())
	wrapped := a.EnforceLifecycle(func(lc *LifecycleContext, args []string) error {
		if lc.Manifest == nil {
			t.Fatal("expected manifest in lifecycle context")
		}
		return nil
	})
	if err := wrapped(cmd, nil); err != nil {
		t.Fatalf("expected lifecycle to succeed, got %v", err)
	}
	if !loaded {
		t.Fatal("expected ManifestLoaded event")
	}
}
