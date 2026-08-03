package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/LordCodex/promptengine/internal/domain/ai"
	contextengine "github.com/LordCodex/promptengine/internal/domain/context"
)

type promptPackage struct {
	Client               string   `json:"client" yaml:"client"`
	TaskType             string   `json:"task_type" yaml:"task_type"`
	Instructions         string   `json:"instructions" yaml:"instructions"`
	Context              string   `json:"context" yaml:"context"`
	Task                 string   `json:"task" yaml:"task"`
	Constraints          []string `json:"constraints" yaml:"constraints"`
	DeveloperPreferences []string `json:"developer_preferences,omitempty" yaml:"developer_preferences,omitempty"`
	ExpectedOutput       []string `json:"expected_output" yaml:"expected_output"`
	SelectedFiles        []string `json:"selected_files" yaml:"selected_files"`
	SelectedDocuments    []string `json:"selected_documents" yaml:"selected_documents"`
	EstimatedContextSize int      `json:"estimated_context_size" yaml:"estimated_context_size"`
}

func buildPromptPackage(client, task string, compiled ai.Request, pkg *contextengine.ContextPackage, preferences []string) promptPackage {
	if client == "" {
		client = "generic"
	}
	out := promptPackage{
		Client:               client,
		TaskType:             task,
		Instructions:         clientInstructions(client),
		Context:              compiled.Context,
		Task:                 compiled.Prompt,
		Constraints:          clientConstraints(client, pkg),
		DeveloperPreferences: dedupeStrings(preferences),
		ExpectedOutput:       expectedOutput(client),
	}
	if pkg != nil {
		out.SelectedFiles = append([]string(nil), pkg.SelectedFiles...)
		out.SelectedDocuments = append([]string(nil), pkg.SelectedDocs...)
		sort.Strings(out.SelectedFiles)
		sort.Strings(out.SelectedDocuments)
		out.EstimatedContextSize = pkg.Summary.FinalSize
		if out.EstimatedContextSize == 0 {
			out.EstimatedContextSize = len(compiled.Context)
		}
	}
	return out
}

func renderPromptPackage(pkg promptPackage, format string) ([]byte, string, error) {
	switch strings.ToLower(format) {
	case "json":
		data, err := json.MarshalIndent(pkg, "", "  ")
		return data, "json", err
	case "text", "plain", "plain-text":
		return []byte(renderPromptText(pkg)), "txt", nil
	case "markdown", "md", "":
		return []byte(renderPromptMarkdown(pkg)), "md", nil
	default:
		return nil, "", fmt.Errorf("unsupported prompt format %q", format)
	}
}

func renderPromptMarkdown(pkg promptPackage) string {
	var b strings.Builder
	b.WriteString("# AI Prompt Package\n\n")
	b.WriteString("## Instructions\n\n")
	b.WriteString(pkg.Instructions)
	b.WriteString("\n\n## Context\n\n")
	b.WriteString(emptyFallback(pkg.Context, "No project context selected."))
	b.WriteString("\n\n## Task\n\n")
	b.WriteString(emptyFallback(pkg.Task, "No task supplied."))
	b.WriteString("\n\n## Constraints\n\n")
	writeList(&b, pkg.Constraints)
	if len(pkg.DeveloperPreferences) > 0 {
		b.WriteString("\n## Developer Preferences\n\n")
		writeList(&b, pkg.DeveloperPreferences)
	}
	b.WriteString("\n## Expected Output\n\n")
	writeList(&b, pkg.ExpectedOutput)
	b.WriteString("\n## Context Size\n\n")
	b.WriteString(fmt.Sprintf("Estimated context size: %d bytes\n", pkg.EstimatedContextSize))
	return b.String()
}

func renderPromptText(pkg promptPackage) string {
	var b strings.Builder
	b.WriteString("Instructions\n")
	b.WriteString(pkg.Instructions)
	b.WriteString("\n\nContext\n")
	b.WriteString(emptyFallback(pkg.Context, "No project context selected."))
	b.WriteString("\n\nTask\n")
	b.WriteString(emptyFallback(pkg.Task, "No task supplied."))
	b.WriteString("\n\nConstraints\n")
	for _, item := range pkg.Constraints {
		b.WriteString("- ")
		b.WriteString(item)
		b.WriteByte('\n')
	}
	if len(pkg.DeveloperPreferences) > 0 {
		b.WriteString("\nDeveloper Preferences\n")
		for _, item := range pkg.DeveloperPreferences {
			b.WriteString("- ")
			b.WriteString(item)
			b.WriteByte('\n')
		}
	}
	b.WriteString("\nExpected Output\n")
	for _, item := range pkg.ExpectedOutput {
		b.WriteString("- ")
		b.WriteString(item)
		b.WriteByte('\n')
	}
	b.WriteString(fmt.Sprintf("\nEstimated context size: %d bytes\n", pkg.EstimatedContextSize))
	return b.String()
}

func writeList(b *strings.Builder, items []string) {
	if len(items) == 0 {
		b.WriteString("- None\n")
		return
	}
	for _, item := range items {
		b.WriteString("- ")
		b.WriteString(item)
		b.WriteByte('\n')
	}
}

func clientInstructions(client string) string {
	switch strings.ToLower(client) {
	case "claude":
		return "Act as a careful senior engineer. Use the provided context first, ask only if required, and produce concise implementation guidance with file-aware reasoning."
	case "codex":
		return "Act as a coding agent in the repository. Follow the selected context, preserve existing architecture, make scoped changes, and verify with tests where applicable."
	case "chatgpt":
		return "Act as a software engineering assistant. Use the context and constraints to produce practical, stepwise guidance or code-ready instructions."
	default:
		return "Act as a senior software engineering assistant. Use only relevant context, respect project standards, and keep the response actionable."
	}
}

func clientConstraints(client string, pkg *contextengine.ContextPackage) []string {
	constraints := []string{
		"Do not redesign the existing architecture unless explicitly requested.",
		"Prefer existing project conventions and domain boundaries.",
		"Keep changes focused on the requested task.",
	}
	if pkg != nil {
		constraints = append(constraints, pkg.RelevantStandards...)
		constraints = append(constraints, pkg.RelatedPlaybooks...)
		if pkg.Summary.BudgetLimit > 0 {
			constraints = append(constraints, fmt.Sprintf("Context was optimized to a %d byte budget.", pkg.Summary.BudgetLimit))
		}
		if len(pkg.Summary.DroppedFiles) > 0 {
			constraints = append(constraints, "Some low-relevance files were omitted during token optimization.")
		}
	}
	switch strings.ToLower(client) {
	case "codex":
		constraints = append(constraints, "When editing files, run relevant tests and report verification.")
	case "claude", "chatgpt":
		constraints = append(constraints, "When code changes are needed, identify exact files and provide implementation-ready snippets.")
	}
	return dedupeStrings(constraints)
}

func expectedOutput(client string) []string {
	out := []string{
		"A short summary of the intended approach.",
		"Implementation details grounded in the selected context.",
		"Verification steps or tests to run.",
	}
	if strings.EqualFold(client, "codex") {
		out = append(out, "A concise final change summary with files touched.")
	}
	return out
}

func defaultPromptExportPath(task, ext string) string {
	task = strings.TrimSpace(strings.ToLower(task))
	if task == "" {
		task = "task"
	}
	task = strings.NewReplacer(" ", "-", "_", "-", "/", "-").Replace(task)
	return task + "-prompt." + ext
}

var copyToClipboard = systemClipboardCopy

func systemClipboardCopy(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	default:
		if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd = exec.Command("wl-copy")
		} else {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		}
	}
	cmd.Stdin = bytes.NewBufferString(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("copy prompt to clipboard: %w", err)
	}
	return nil
}

func emptyFallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func dedupeStrings(items []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}
