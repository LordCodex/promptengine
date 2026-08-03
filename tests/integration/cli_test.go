package integration

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/app"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestIntegration_InitCommand(t *testing.T) {
	var buf bytes.Buffer
	cliApp, err := app.Bootstrap(&buf, false)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	// Substitute filesystem with memory mock to keep tests isolated
	cliApp.FS = filesystem.NewMockFileSystem()

	ctx := context.Background()
	args := []string{"init", "--force"}

	if err := cliApp.Execute(ctx, args); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Initialising PromptEngine") {
		t.Errorf("expected initialization message, got: %s", output)
	}

	if !cliApp.FS.Exists("playbook-manifest.json") {
		t.Error("expected playbook-manifest.json to be created")
	}
}

func TestIntegration_VersionCommand(t *testing.T) {
	var buf bytes.Buffer
	cliApp, err := app.Bootstrap(&buf, false)
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
	cliApp, err := app.Bootstrap(&buf, false)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	ctx := context.Background()
	args := []string{"--help"}

	if err := cliApp.Execute(ctx, args); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "PromptEngine is the AI Coding Standards Engine") {
		t.Errorf("expected help output, got: %s", output)
	}
}
