package app

import (
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/ai"
	contextengine "github.com/LordCodex/promptengine/internal/domain/context"
)

func TestRenderPromptPackageMarkdown(t *testing.T) {
	pkg := buildPromptPackage("claude", "feature", ai.Request{Prompt: "Add billing", Context: "BillingService.php"}, &contextengine.ContextPackage{
		SelectedFiles: []string{"app/BillingService.php"},
		Summary:       contextengine.OptimizationSummary{FinalSize: 128, BudgetLimit: 500},
	}, []string{"prefer simple solutions"})
	data, ext, err := renderPromptPackage(pkg, "markdown")
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	out := string(data)
	if ext != "md" || !strings.Contains(out, "## Instructions") || !strings.Contains(out, "## Developer Preferences") || !strings.Contains(out, "Estimated context size: 128 bytes") {
		t.Fatalf("unexpected markdown output:\n%s", out)
	}
}

func TestRenderPromptPackageJSON(t *testing.T) {
	pkg := buildPromptPackage("codex", "bug_fix", ai.Request{Prompt: "Fix retry", Context: "RetryService.go"}, nil, nil)
	data, ext, err := renderPromptPackage(pkg, "json")
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if ext != "json" || !strings.Contains(string(data), `"client": "codex"`) || !strings.Contains(string(data), `"task": "Fix retry"`) {
		t.Fatalf("unexpected json output:\n%s", string(data))
	}
}

func TestPromptPackageTemplateAndLargeContextHints(t *testing.T) {
	ctxPkg := &contextengine.ContextPackage{
		RelatedPlaybooks: []string{"security-standard"},
		Summary: contextengine.OptimizationSummary{
			FinalSize:   2048,
			BudgetLimit: 4096,
			DroppedFiles: []string{
				"logs/noisy.log",
			},
		},
	}
	pkg := buildPromptPackage("chatgpt", "review", ai.Request{Prompt: "Review auth", Context: "auth context"}, ctxPkg, nil)
	joined := strings.Join(pkg.Constraints, "\n")
	if !strings.Contains(pkg.Instructions, "software engineering assistant") {
		t.Fatalf("expected chatgpt instructions, got %q", pkg.Instructions)
	}
	if !strings.Contains(joined, "security-standard") || !strings.Contains(joined, "omitted") {
		t.Fatalf("expected optimization constraints, got %#v", pkg.Constraints)
	}
}

func TestClipboardCopyCanBeInjected(t *testing.T) {
	old := copyToClipboard
	defer func() { copyToClipboard = old }()
	called := false
	copyToClipboard = func(text string) error {
		called = strings.Contains(text, "hello")
		return nil
	}
	if err := copyToClipboard("hello"); err != nil {
		t.Fatalf("copy failed: %v", err)
	}
	if !called {
		t.Fatal("expected injected clipboard copier to be called")
	}
}
