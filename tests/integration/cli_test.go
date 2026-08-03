package integration

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/app"
)

func TestIntegration_RootStarts(t *testing.T) {
	var buf bytes.Buffer
	cliApp, err := app.Bootstrap(app.BootstrapOptions{Out: &buf, Err: &buf})
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	ctx := context.Background()
	if err := cliApp.Execute(ctx, []string{}); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
}

func TestIntegration_VersionCommand(t *testing.T) {
	var buf bytes.Buffer
	cliApp, err := app.Bootstrap(app.BootstrapOptions{Out: &buf, Err: &buf})
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	ctx := context.Background()
	args := []string{"version"}

	if err := cliApp.Execute(ctx, args); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "PromptEngine v") {
		t.Errorf("expected version output, got: %s", output)
	}
}

func TestIntegration_HelpCommand(t *testing.T) {
	var buf bytes.Buffer
	cliApp, err := app.Bootstrap(app.BootstrapOptions{Out: &buf, Err: &buf})
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	ctx := context.Background()
	args := []string{"--help"}

	if err := cliApp.Execute(ctx, args); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "PromptEngine CLI foundation") {
		t.Errorf("expected help output, got: %s", output)
	}
}

func TestIntegration_ProductionCommands(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{name: "scan", args: []string{"scan"}, want: "Project:"},
		{name: "doctor", args: []string{"doctor"}, want: "Quality Report"},
		{name: "docs-generate", args: []string{"docs", "generate", "--doc", "architecture", "--overwrite"}, want: "Generated Architecture"},
		{name: "prompt-markdown", args: []string{"prompt", "--task", "feature", "--request", "Add billing", "--format", "markdown", "--out", "billing-prompt.md"}, want: "Prompt package exported"},
		{name: "prompt-json", args: []string{"prompt", "--task", "bug_fix", "--request", "Fix retry", "--format", "json", "--out", "retry-prompt.json"}, want: "Prompt package exported"},
		{name: "agents-sync", args: []string{"agents", "sync", "--agent", "claude"}, want: "Agent instructions synchronized"},
		{name: "context-export", args: []string{"context", "export", "--task", "feature", "--agent", "codex", "--format", "markdown"}, want: "Context exported"},
		{name: "profile-init", args: []string{"profile", "init"}, want: "Created .promptengine/profile.yaml"},
		{name: "task", args: []string{"task", "--template", "feature", "--out", "subscription-task.md", "Add subscription billing"}, want: "Workflow: feature-implementation"},
		{name: "verify", args: []string{"verify"}, want: "Verification"},
		{name: "memory-add", args: []string{"memory", "add", "--key", "test-command", "--value", "go test ./..."}, want: "Stored memory"},
		{name: "decision-store", args: []string{"decisions", "store", "--title", "Use UUID primary keys", "--reason", "Separate public IDs", "--affected", "models,apis"}, want: "Stored decision"},
		{name: "decisions-list", args: []string{"decisions", "list"}, want: "Use UUID primary keys"},
		{name: "insights", args: []string{"insights"}, want: "Insights"},
		{name: "impact", args: []string{"impact"}, want: "Impact Analysis"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cliApp, err := app.Bootstrap(app.BootstrapOptions{Out: &buf, Err: &buf})
			if err != nil {
				t.Fatalf("bootstrap failed: %v", err)
			}
			if err := cliApp.Execute(context.Background(), tc.args); err != nil {
				t.Fatalf("execute failed: %v\noutput:\n%s", err, buf.String())
			}
			if !strings.Contains(buf.String(), tc.want) {
				t.Fatalf("expected output containing %q, got:\n%s", tc.want, buf.String())
			}
		})
	}
	for _, file := range []string{"billing-prompt.md", "retry-prompt.json"} {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("expected prompt export %s: %v", file, err)
		}
	}
	for _, file := range []string{"CLAUDE.md", "codex-context.md"} {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("expected agent artifact %s: %v", file, err)
		}
	}
	for _, file := range []string{".promptengine/profile.yaml", ".promptengine/memory.yaml", ".promptengine/decisions.yaml", "subscription-task.md"} {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("expected personal workflow artifact %s: %v", file, err)
		}
	}
}

func TestIntegration_InitWithAgentInstructions(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	var buf bytes.Buffer
	cliApp, err := app.Bootstrap(app.BootstrapOptions{Out: &buf, Err: &buf})
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}
	if err := cliApp.Execute(context.Background(), []string{"init", "--agents", "codex,claude,cursor,windsurf"}); err != nil {
		t.Fatalf("execute failed: %v\noutput:\n%s", err, buf.String())
	}
	for _, file := range []string{"AGENTS.md", "CLAUDE.md", ".cursor/rules/promptengine.md", ".windsurf/rules/promptengine.md", ".promptengine"} {
		if _, err := os.Stat(file); err != nil {
			t.Fatalf("expected initialized agent artifact %s: %v", file, err)
		}
	}
}
