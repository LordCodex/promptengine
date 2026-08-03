package sdk

import (
	"context"
	"strings"
	"testing"
)

type fakeRunner struct {
	workingDir string
	name       string
	args       []string
	stdout     []byte
	stderr     []byte
	err        error
}

func (r *fakeRunner) Run(ctx context.Context, workingDir string, name string, args ...string) ([]byte, []byte, error) {
	r.workingDir = workingDir
	r.name = name
	r.args = append([]string(nil), args...)
	return r.stdout, r.stderr, r.err
}

func TestClientGeneratePromptUsesCLI(t *testing.T) {
	runner := &fakeRunner{stdout: []byte(`{"status":"exported"}`)}
	client := NewClientWithRunner(Options{BinaryPath: "/bin/promptengine", WorkingDir: "/repo", ConfigPath: ".promptengine.yaml"}, runner)
	out, err := client.GeneratePrompt(context.Background(), PromptRequest{Task: "bug_fix", Request: "Fix retry", Agent: "codex", Format: "markdown", OutputPath: "retry.md"})
	if err != nil {
		t.Fatalf("prompt failed: %v", err)
	}
	if string(out) != `{"status":"exported"}` {
		t.Fatalf("unexpected output %s", out)
	}
	joined := strings.Join(runner.args, " ")
	for _, expected := range []string{"--config .promptengine.yaml", "prompt", "--json", "--task bug_fix", "--client codex", "--format markdown", "--out retry.md"} {
		if !strings.Contains(joined, expected) {
			t.Fatalf("expected args to contain %q, got %q", expected, joined)
		}
	}
}

func TestClientAnalyzeProjectUsesWorkspace(t *testing.T) {
	runner := &fakeRunner{stdout: []byte(`{"project":{}}`)}
	client := NewClientWithRunner(Options{WorkingDir: "/repo"}, runner)
	if _, err := client.AnalyzeProject(context.Background()); err != nil {
		t.Fatalf("analyze failed: %v", err)
	}
	if runner.workingDir != "/repo" || runner.name != "promptengine" || strings.Join(runner.args, " ") != "scan --json" {
		t.Fatalf("unexpected invocation: dir=%s name=%s args=%v", runner.workingDir, runner.name, runner.args)
	}
}
