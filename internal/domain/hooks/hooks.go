package hooks

import (
	"context"
	"fmt"
	"strings"

	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

// HookType identifies the target automation system
type HookType string

const (
	GitPreCommit   HookType = "git-pre-commit"
	GitPrePush     HookType = "git-pre-push"
	GitHubWorkflow HookType = "github-actions"
	GitLabPipeline HookType = "gitlab-ci"
	AzurePipeline  HookType = "azure-pipelines"
	LocalScript    HookType = "local"
)

// HookPolicy controls whether a hook failure blocks or warns
type HookPolicy string

const (
	PolicyWarn    HookPolicy = "warn"    // default — allow commit but emit warning
	PolicyEnforce HookPolicy = "enforce" // block commit on failure
)

// HookConfig configures the behaviour of an individual hook
type HookConfig struct {
	Policy  HookPolicy
	Command string // command to run, e.g. "promptengine doctor"
}

// Hook handles automated validation scripting checks
type Hook interface {
	Name() string
	Type() HookType
	Config() HookConfig
	Install(fs filesystem.FileSystem) error
	Uninstall(fs filesystem.FileSystem) error
}

type EventHook interface {
	ID() string
	Event() eventbus.EventType
	Handle(ctx context.Context, event eventbus.Event) error
}

// GitHook is a concrete Git hook implementation that writes shell scripts
type GitHook struct {
	name   string
	htype  HookType
	config HookConfig
}

func NewGitHook(name string, htype HookType, cfg HookConfig) *GitHook {
	return &GitHook{name: name, htype: htype, config: cfg}
}

func (h *GitHook) Name() string       { return h.name }
func (h *GitHook) Type() HookType     { return h.htype }
func (h *GitHook) Config() HookConfig { return h.config }

func (h *GitHook) Install(fs filesystem.FileSystem) error {
	script := h.renderScript()
	hookPath := fmt.Sprintf(".git/hooks/%s", hookFileName(h.htype))
	return fs.WriteFile(hookPath, []byte(script), 0755)
}

func (h *GitHook) Uninstall(fs filesystem.FileSystem) error {
	hookPath := fmt.Sprintf(".git/hooks/%s", hookFileName(h.htype))
	if fs.Exists(hookPath) {
		return fs.WriteFile(hookPath, []byte("#!/bin/sh\n# removed by PromptEngine\n"), 0755)
	}
	return nil
}

func (h *GitHook) renderScript() string {
	var sb strings.Builder
	sb.WriteString("#!/bin/sh\n")
	sb.WriteString("# Managed by PromptEngine — do not edit manually\n\n")
	sb.WriteString(fmt.Sprintf("%s\n", h.config.Command))
	sb.WriteString("EXIT_CODE=$?\n")
	if h.config.Policy == PolicyEnforce {
		sb.WriteString("if [ $EXIT_CODE -ne 0 ]; then\n")
		sb.WriteString("  echo \"[PromptEngine] Hook failed. Commit blocked.\"\n")
		sb.WriteString("  exit $EXIT_CODE\n")
		sb.WriteString("fi\n")
	} else {
		sb.WriteString("if [ $EXIT_CODE -ne 0 ]; then\n")
		sb.WriteString("  echo \"[PromptEngine] Warning: hook check failed. Proceeding anyway.\"\n")
		sb.WriteString("fi\n")
	}
	return sb.String()
}

func hookFileName(t HookType) string {
	switch t {
	case GitPreCommit:
		return "pre-commit"
	case GitPrePush:
		return "pre-push"
	default:
		return string(t)
	}
}

// CITemplate generates CI pipeline configuration files
type CITemplate struct {
	htype   HookType
	command string
}

func NewCITemplate(htype HookType, command string) *CITemplate {
	return &CITemplate{htype: htype, command: command}
}

func (t *CITemplate) Generate() string {
	switch t.htype {
	case GitHubWorkflow:
		return fmt.Sprintf(`name: PromptEngine Check
on: [push, pull_request]
jobs:
  promptengine:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run PromptEngine
        run: %s
`, t.command)
	case GitLabPipeline:
		return fmt.Sprintf(`promptengine:
  stage: test
  script:
    - %s
`, t.command)
	case AzurePipeline:
		return fmt.Sprintf(`steps:
- script: %s
  displayName: 'PromptEngine Check'
`, t.command)
	default:
		return fmt.Sprintf("# Run: %s\n", t.command)
	}
}

// Registry handles optional compliance hooks management
type Registry struct {
	fs         filesystem.FileSystem
	hooks      []Hook
	eventHooks []EventHook
	bus        *eventbus.EventBus
	policy     HookPolicy // org-level default policy
}

func NewRegistry(fs filesystem.FileSystem) *Registry {
	return &Registry{
		fs:     fs,
		hooks:  make([]Hook, 0),
		policy: PolicyWarn,
	}
}

func (r *Registry) SetOrgPolicy(p HookPolicy) {
	r.policy = p
}

func (r *Registry) Register(h Hook) {
	r.hooks = append(r.hooks, h)
}

func (r *Registry) RegisterEventHook(h EventHook) {
	r.eventHooks = append(r.eventHooks, h)
	if r.bus != nil {
		r.subscribe(h)
	}
}

func (r *Registry) Attach(bus *eventbus.EventBus) {
	r.bus = bus
	for _, hook := range r.eventHooks {
		r.subscribe(hook)
	}
}

func (r *Registry) subscribe(hook EventHook) {
	h := hook
	r.bus.Subscribe(h.Event(), func(e eventbus.Event) {
		_ = h.Handle(context.Background(), e)
	})
}

func (r *Registry) Dispatch(ctx context.Context, e eventbus.Event) error {
	for _, hook := range r.eventHooks {
		if hook.Event() != e.Type {
			continue
		}
		if err := hook.Handle(ctx, e); err != nil {
			return fmt.Errorf("event hook '%s' failed: %w", hook.ID(), err)
		}
	}
	return nil
}

func (r *Registry) InstallAll() error {
	for _, h := range r.hooks {
		if err := h.Install(r.fs); err != nil {
			return fmt.Errorf("failed to install hook '%s': %w", h.Name(), err)
		}
	}
	return nil
}

func (r *Registry) UninstallAll() error {
	for _, h := range r.hooks {
		if err := h.Uninstall(r.fs); err != nil {
			return fmt.Errorf("failed to uninstall hook '%s': %w", h.Name(), err)
		}
	}
	return nil
}

func (r *Registry) ListHooks() []Hook {
	return r.hooks
}

func (r *Registry) ListEventHooks() []EventHook {
	return r.eventHooks
}
