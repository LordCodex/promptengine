package hooks

import (
	"context"
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestGitHook_Install_WarnPolicy(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	// Git hooks directory must exist
	_ = fs.WriteFile(".git/hooks/.keep", []byte{}, 0644)

	hook := NewGitHook("promptengine-precommit", GitPreCommit, HookConfig{
		Policy:  PolicyWarn,
		Command: "promptengine doctor",
	})

	if err := hook.Install(fs); err != nil {
		t.Fatalf("expected no error installing hook, got %v", err)
	}

	data, err := fs.ReadFile(".git/hooks/pre-commit")
	if err != nil {
		t.Fatalf("expected hook script to be written, got error: %v", err)
	}
	script := string(data)
	if !strings.Contains(script, "promptengine doctor") {
		t.Error("expected hook script to contain the configured command")
	}
	if !strings.Contains(script, "Proceeding anyway") {
		t.Error("expected warn-policy hook to contain non-blocking message")
	}
}

func TestGitHook_Install_EnforcePolicy(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile(".git/hooks/.keep", []byte{}, 0644)

	hook := NewGitHook("pe-enforce", GitPreCommit, HookConfig{
		Policy:  PolicyEnforce,
		Command: "promptengine doctor",
	})
	_ = hook.Install(fs)

	data, _ := fs.ReadFile(".git/hooks/pre-commit")
	if !strings.Contains(string(data), "Commit blocked") {
		t.Error("expected enforce-policy hook to contain blocking message")
	}
}

func TestCITemplate_Generate_GitHub(t *testing.T) {
	tmpl := NewCITemplate(GitHubWorkflow, "promptengine doctor")
	out := tmpl.Generate()
	if !strings.Contains(out, "actions/checkout") {
		t.Error("expected GitHub Actions template to reference checkout action")
	}
	if !strings.Contains(out, "promptengine doctor") {
		t.Error("expected GitHub Actions template to include the command")
	}
}

func TestRegistry_InstallAll(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	_ = fs.WriteFile(".git/hooks/.keep", []byte{}, 0644)

	reg := NewRegistry(fs)
	reg.Register(NewGitHook("h1", GitPreCommit, HookConfig{Policy: PolicyWarn, Command: "echo ok"}))
	reg.Register(NewGitHook("h2", GitPrePush, HookConfig{Policy: PolicyWarn, Command: "echo ok"}))

	if err := reg.InstallAll(); err != nil {
		t.Errorf("expected all hooks to install cleanly, got: %v", err)
	}
	if len(reg.ListHooks()) != 2 {
		t.Errorf("expected 2 hooks registered, got %d", len(reg.ListHooks()))
	}
}

type mockEventHook struct {
	id     string
	event  eventbus.EventType
	called bool
}

func (m *mockEventHook) ID() string                { return m.id }
func (m *mockEventHook) Event() eventbus.EventType { return m.event }
func (m *mockEventHook) Handle(ctx context.Context, e eventbus.Event) error {
	m.called = true
	return nil
}

func TestRegistry_EventHookDispatchAndAttach(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	reg := NewRegistry(fs)
	hook := &mockEventHook{id: "context-hook", event: eventbus.ContextBuilt}
	reg.RegisterEventHook(hook)
	if err := reg.Dispatch(context.Background(), eventbus.Event{Type: eventbus.ContextBuilt}); err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}
	if !hook.called {
		t.Fatal("expected hook dispatch")
	}
	hook.called = false
	bus := eventbus.NewEventBus()
	reg.Attach(bus)
	bus.Publish(eventbus.Event{Type: eventbus.ContextBuilt})
	if !hook.called {
		t.Fatal("expected attached hook to run")
	}
	lateHook := &mockEventHook{id: "late-hook", event: eventbus.PluginInstalled}
	reg.RegisterEventHook(lateHook)
	bus.Publish(eventbus.Event{Type: eventbus.PluginInstalled})
	if !lateHook.called {
		t.Fatal("expected hook registered after attach to run")
	}
}
