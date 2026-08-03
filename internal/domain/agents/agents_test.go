package agents

import (
	"strings"
	"testing"

	contextengine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

func TestPlatform_GenerateInstruction(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	platform := NewPlatform(fs, nil)
	project := discovery.NewProjectModel(".")
	project.Languages = []string{"Go"}
	generated, err := platform.Generate(InstructionRequest{
		Profile: "claude",
		Project: project,
		Manifest: &manifest.Manifest{Playbooks: []manifest.PlaybookDefinition{{
			ID: "universal", Name: "Universal", Category: manifest.CategoryCore, Location: "core/05-universal-coding-standards.md",
		}}},
	})
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	data, err := fs.ReadFile(generated.Path)
	if err != nil {
		t.Fatalf("expected generated file: %v", err)
	}
	out := string(data)
	if generated.Path != "CLAUDE.md" || !strings.Contains(out, "Workflow Rules") || !strings.Contains(out, "core/05-universal-coding-standards.md") {
		t.Fatalf("unexpected instruction output:\n%s", out)
	}
}

func TestPlatform_SyncDetectsGeneratedUpdatedCurrent(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	platform := NewPlatform(fs, nil)
	first, err := platform.Sync(InstructionRequest{Profile: "codex", Project: discovery.NewProjectModel(".")})
	if err != nil {
		t.Fatalf("first sync failed: %v", err)
	}
	if len(first.Generated) != 1 {
		t.Fatalf("expected generated file, got %#v", first)
	}
	second, err := platform.Sync(InstructionRequest{Profile: "codex", Project: discovery.NewProjectModel(".")})
	if err != nil {
		t.Fatalf("second sync failed: %v", err)
	}
	if len(second.Current) != 1 {
		t.Fatalf("expected current file, got %#v", second)
	}
	if err := fs.WriteFile("AGENTS.md", []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	third, err := platform.Sync(InstructionRequest{Profile: "codex", Project: discovery.NewProjectModel(".")})
	if err != nil {
		t.Fatalf("third sync failed: %v", err)
	}
	if len(third.Updated) != 1 {
		t.Fatalf("expected updated file, got %#v", third)
	}
}

func TestPlatform_ExportContextAndEvent(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	events := eventbus.NewEventBus()
	seen := false
	events.Subscribe(eventbus.ContextExported, func(e eventbus.Event) { seen = true })
	platform := NewPlatform(fs, events)
	export, err := platform.ExportContext(ContextExportRequest{
		Task:   "feature",
		Agent:  "codex",
		Format: "markdown",
		Package: &contextengine.ContextPackage{
			SystemPrompt:  "Use these files",
			SelectedFiles: []string{"main.go"},
			Summary:       contextengine.OptimizationSummary{FinalSize: 128},
		},
	})
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	data, err := fs.ReadFile(export.File)
	if err != nil {
		t.Fatalf("expected exported context: %v", err)
	}
	if export.File != "codex-context.md" || !strings.Contains(string(data), "Relevant Files") || !seen {
		t.Fatalf("unexpected export: %#v\n%s", export, string(data))
	}
}

func TestPlatform_CustomProfileRegistration(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	platform := NewPlatform(fs, nil)
	platform.Register(AgentProfile{ID: "custom", Name: "Custom Agent", InstructionFile: ".custom/instructions.md", Format: "markdown"})
	report, err := platform.Sync(InstructionRequest{Profile: "custom", Project: discovery.NewProjectModel(".")})
	if err != nil {
		t.Fatalf("sync failed: %v", err)
	}
	if len(report.Generated) != 1 || report.Generated[0].Path != ".custom/instructions.md" {
		t.Fatalf("expected custom profile generation, got %#v", report)
	}
}
