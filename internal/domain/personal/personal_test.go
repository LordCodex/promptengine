package personal

import (
	"strings"
	"testing"

	contextengine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestProfileLoadSave(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	p := NewPlatform(fs, nil)
	profile := DefaultProfile()
	profile.Languages = []string{"Go"}
	if err := p.SaveProfile("", profile); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := p.LoadProfile("")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(loaded.Languages) != 1 || loaded.Languages[0] != "Go" {
		t.Fatalf("unexpected profile %#v", loaded)
	}
}

func TestMemoryRejectsSensitiveValues(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	p := NewPlatform(fs, nil)
	if _, err := p.AddMemory("token", "api_key=secret"); err == nil {
		t.Fatal("expected sensitive memory rejection")
	}
}

func TestBuildTaskPackageUsesGitAndProfile(t *testing.T) {
	pkg := BuildTaskPackage(TaskRequest{
		Description: "Add payment refund support",
		Template:    "feature",
		Profile:     DefaultProfile(),
		Git:         GitContext{Branch: "feature/refunds", ChangedFiles: []string{"PaymentService.php"}},
		Context:     &contextengine.ContextPackage{SelectedDocs: []string{"docs/API.md"}, Summary: contextengine.OptimizationSummary{FinalSize: 512}},
	})
	if pkg.RecommendedWorkflow != "feature-implementation" {
		t.Fatalf("unexpected workflow %s", pkg.RecommendedWorkflow)
	}
	if !strings.Contains(pkg.Prompt, "Developer Preferences") || !strings.Contains(strings.Join(pkg.SelectedFiles, ","), "PaymentService.php") {
		t.Fatalf("unexpected task package %#v", pkg)
	}
}
