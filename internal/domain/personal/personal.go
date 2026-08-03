package personal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	contextengine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"gopkg.in/yaml.v3"
)

const (
	DefaultProfilePath = ".promptengine/profile.yaml"
	DefaultMemoryPath  = ".promptengine/memory.yaml"
)

type Profile struct {
	Developer     DeveloperInfo `json:"developer" yaml:"developer"`
	Languages     []string      `json:"languages" yaml:"languages"`
	Frameworks    []string      `json:"frameworks" yaml:"frameworks"`
	Tools         []string      `json:"tools" yaml:"tools"`
	Architecture  []string      `json:"architecture_preferences" yaml:"architecture_preferences"`
	Principles    []string      `json:"principles" yaml:"principles"`
	Testing       []string      `json:"testing_preferences" yaml:"testing_preferences"`
	Documentation []string      `json:"documentation_preferences" yaml:"documentation_preferences"`
	AI            []string      `json:"ai_prompting_preferences" yaml:"ai_prompting_preferences"`
}

type DeveloperInfo struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

type Memory struct {
	ProjectPreferences []string          `json:"project_preferences" yaml:"project_preferences"`
	Decisions          []string          `json:"decisions" yaml:"decisions"`
	CommonCommands     []string          `json:"common_commands" yaml:"common_commands"`
	Workflows          []string          `json:"frequently_used_workflows" yaml:"frequently_used_workflows"`
	Notes              map[string]string `json:"notes,omitempty" yaml:"notes,omitempty"`
	UpdatedAt          time.Time         `json:"updated_at" yaml:"updated_at"`
}

type GitContext struct {
	Branch        string   `json:"branch" yaml:"branch"`
	ChangedFiles  []string `json:"changed_files" yaml:"changed_files"`
	RecentCommits []string `json:"recent_commits" yaml:"recent_commits"`
	DiffSummary   string   `json:"diff_summary,omitempty" yaml:"diff_summary,omitempty"`
}

type TaskTemplate struct {
	ID              string   `json:"id" yaml:"id"`
	TaskType        string   `json:"task_type" yaml:"task_type"`
	Workflow        string   `json:"workflow" yaml:"workflow"`
	RequiredContext []string `json:"required_context" yaml:"required_context"`
	ValidationSteps []string `json:"validation_steps" yaml:"validation_steps"`
	ExpectedOutput  []string `json:"expected_output" yaml:"expected_output"`
}

type TaskRequest struct {
	Description string
	Template    string
	Profile     Profile
	Memory      Memory
	Git         GitContext
	Project     *discovery.ProjectModel
	Context     *contextengine.ContextPackage
}

type TaskPackage struct {
	Summary             string                            `json:"task_summary" yaml:"task_summary"`
	Intent              string                            `json:"intent" yaml:"intent"`
	RecommendedWorkflow string                            `json:"recommended_workflow" yaml:"recommended_workflow"`
	SelectedFiles       []string                          `json:"selected_files" yaml:"selected_files"`
	SelectedDocuments   []string                          `json:"selected_documents" yaml:"selected_documents"`
	Prompt              string                            `json:"generated_prompt" yaml:"generated_prompt"`
	Template            TaskTemplate                      `json:"template" yaml:"template"`
	Profile             Profile                           `json:"profile" yaml:"profile"`
	Git                 GitContext                        `json:"git_context" yaml:"git_context"`
	ContextSummary      contextengine.OptimizationSummary `json:"context_summary" yaml:"context_summary"`
}

type Platform struct {
	fs     filesystem.FileSystem
	events *eventbus.EventBus
}

func NewPlatform(fs filesystem.FileSystem, events *eventbus.EventBus) *Platform {
	return &Platform{fs: fs, events: events}
}

func DefaultProfile() Profile {
	return Profile{
		Languages:  []string{"PHP", "Dart", "TypeScript"},
		Frameworks: []string{"Laravel", "Vue", "Flutter"},
		Principles: []string{
			"prefer simple solutions",
			"avoid unnecessary abstractions",
			"follow existing architecture",
			"explain important technical decisions",
			"prioritize maintainability",
		},
		Testing:       []string{"write behavior-focused tests", "run the narrowest useful verification first"},
		Documentation: []string{"keep architecture and API docs synchronized with meaningful changes"},
		AI:            []string{"be concise", "state assumptions", "produce actionable implementation steps"},
	}
}

func (p *Platform) LoadProfile(path string) (Profile, error) {
	if path == "" {
		path = DefaultProfilePath
	}
	if !p.fs.Exists(path) {
		return DefaultProfile(), nil
	}
	data, err := p.fs.ReadFile(path)
	if err != nil {
		return Profile{}, err
	}
	var wrapper struct {
		Profile Profile `json:"profile" yaml:"profile"`
	}
	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return Profile{}, err
	}
	if len(wrapper.Profile.Principles) == 0 {
		return DefaultProfile(), nil
	}
	return wrapper.Profile, nil
}

func (p *Platform) SaveProfile(path string, profile Profile) error {
	if path == "" {
		path = DefaultProfilePath
	}
	data, err := yaml.Marshal(map[string]Profile{"profile": profile})
	if err != nil {
		return err
	}
	return p.fs.WriteFile(path, data, 0644)
}

func (p *Platform) LoadMemory(path string) (Memory, error) {
	if path == "" {
		path = DefaultMemoryPath
	}
	if !p.fs.Exists(path) {
		return Memory{Notes: map[string]string{}}, nil
	}
	data, err := p.fs.ReadFile(path)
	if err != nil {
		return Memory{}, err
	}
	var mem Memory
	if err := yaml.Unmarshal(data, &mem); err != nil {
		return Memory{}, err
	}
	if mem.Notes == nil {
		mem.Notes = map[string]string{}
	}
	return mem, nil
}

func (p *Platform) SaveMemory(path string, mem Memory) error {
	if path == "" {
		path = DefaultMemoryPath
	}
	mem.UpdatedAt = time.Now().UTC()
	data, err := yaml.Marshal(mem)
	if err != nil {
		return err
	}
	return p.fs.WriteFile(path, data, 0644)
}

func (p *Platform) AddMemory(key, value string) (Memory, error) {
	mem, err := p.LoadMemory("")
	if err != nil {
		return mem, err
	}
	if mem.Notes == nil {
		mem.Notes = map[string]string{}
	}
	if looksSensitive(value) {
		return mem, fmt.Errorf("memory value appears sensitive and was not stored")
	}
	mem.Notes[key] = value
	return mem, p.SaveMemory("", mem)
}

func DetectGitContext(ctx context.Context, root string) GitContext {
	return GitContext{
		Branch:        runGit(ctx, root, "rev-parse", "--abbrev-ref", "HEAD"),
		ChangedFiles:  splitLines(runGit(ctx, root, "diff", "--name-only", "HEAD")),
		RecentCommits: splitLines(runGit(ctx, root, "log", "--oneline", "-5")),
		DiffSummary:   runGit(ctx, root, "diff", "--stat", "HEAD"),
	}
}

func ResolveTemplate(id string) TaskTemplate {
	if id == "" {
		id = "feature"
	}
	templates := map[string]TaskTemplate{
		"feature":      {ID: "feature", TaskType: "feature", Workflow: "feature-implementation", RequiredContext: []string{"architecture docs", "related source files", "tests"}, ValidationSteps: []string{"run relevant tests", "check documentation impact"}, ExpectedOutput: []string{"implementation plan", "files to change", "tests to run"}},
		"bug":          {ID: "bug", TaskType: "bug_fix", Workflow: "bug-fix", RequiredContext: []string{"affected module", "tests", "recent changes"}, ValidationSteps: []string{"reproduce or explain bug", "run regression tests"}, ExpectedOutput: []string{"root cause", "fix plan", "verification"}},
		"bug_fix":      {ID: "bug_fix", TaskType: "bug_fix", Workflow: "bug-fix", RequiredContext: []string{"affected module", "tests", "recent changes"}, ValidationSteps: []string{"reproduce or explain bug", "run regression tests"}, ExpectedOutput: []string{"root cause", "fix plan", "verification"}},
		"refactor":     {ID: "refactor", TaskType: "refactor", Workflow: "feature-implementation", RequiredContext: []string{"target module", "tests", "architecture rules"}, ValidationSteps: []string{"run existing tests", "verify behavior unchanged"}, ExpectedOutput: []string{"safe refactor plan", "risk notes"}},
		"review":       {ID: "review", TaskType: "review", Workflow: "bug-fix", RequiredContext: []string{"changed files", "standards", "tests"}, ValidationSteps: []string{"quality review", "security review"}, ExpectedOutput: []string{"findings", "recommendations"}},
		"architecture": {ID: "architecture", TaskType: "architecture_review", Workflow: "feature-implementation", RequiredContext: []string{"architecture docs", "project structure", "dependencies"}, ValidationSteps: []string{"document decision", "review tradeoffs"}, ExpectedOutput: []string{"architecture recommendation"}},
		"docs":         {ID: "docs", TaskType: "documentation_update", Workflow: "feature-implementation", RequiredContext: []string{"existing docs", "changed files"}, ValidationSteps: []string{"validate docs"}, ExpectedOutput: []string{"documentation changes"}},
		"security":     {ID: "security", TaskType: "review", Workflow: "bug-fix", RequiredContext: []string{"security standards", "auth flows", "changed files"}, ValidationSteps: []string{"threat review", "security tests"}, ExpectedOutput: []string{"security findings"}},
	}
	if tmpl, ok := templates[id]; ok {
		return tmpl
	}
	return templates["feature"]
}

func BuildTaskPackage(req TaskRequest) TaskPackage {
	tmpl := ResolveTemplate(req.Template)
	if req.Template == "" {
		tmpl = inferTemplate(req.Description)
	}
	selectedFiles := []string{}
	selectedDocs := []string{}
	var summary contextengine.OptimizationSummary
	if req.Context != nil {
		selectedFiles = append(selectedFiles, req.Context.SelectedFiles...)
		selectedDocs = append(selectedDocs, req.Context.SelectedDocs...)
		summary = req.Context.Summary
	}
	selectedFiles = uniqueSorted(append(selectedFiles, req.Git.ChangedFiles...))
	prompt := RenderTaskPrompt(req.Description, tmpl, req.Profile, req.Memory, req.Git, selectedFiles, selectedDocs)
	return TaskPackage{
		Summary:             summarizeTask(req.Description),
		Intent:              tmpl.TaskType,
		RecommendedWorkflow: tmpl.Workflow,
		SelectedFiles:       selectedFiles,
		SelectedDocuments:   uniqueSorted(selectedDocs),
		Prompt:              prompt,
		Template:            tmpl,
		Profile:             req.Profile,
		Git:                 req.Git,
		ContextSummary:      summary,
	}
}

func RenderTaskPrompt(task string, tmpl TaskTemplate, profile Profile, mem Memory, git GitContext, files, docs []string) string {
	var b strings.Builder
	b.WriteString("# AI Task Package\n\n## Context\n\n")
	b.WriteString("- Current branch: ")
	b.WriteString(empty(git.Branch, "unknown"))
	b.WriteString("\n- Changed files:\n")
	writeBullets(&b, files)
	b.WriteString("\n- Relevant documentation:\n")
	writeBullets(&b, docs)
	b.WriteString("\n## Task\n\n")
	b.WriteString(task)
	b.WriteString("\n\n## Constraints\n\n")
	writeBullets(&b, append(tmpl.RequiredContext, tmpl.ValidationSteps...))
	b.WriteString("\n## Developer Preferences\n\n")
	writeBullets(&b, profile.PreferenceLines())
	if len(mem.Decisions) > 0 || len(mem.ProjectPreferences) > 0 {
		b.WriteString("\n## Personal Memory\n\n")
		writeBullets(&b, append(mem.ProjectPreferences, mem.Decisions...))
	}
	b.WriteString("\n## Expected Output\n\n")
	writeBullets(&b, tmpl.ExpectedOutput)
	return b.String()
}

func (p Profile) PreferenceLines() []string {
	var out []string
	out = append(out, prefixed("Language", p.Languages)...)
	out = append(out, prefixed("Framework", p.Frameworks)...)
	out = append(out, p.Principles...)
	out = append(out, p.Architecture...)
	out = append(out, p.Testing...)
	out = append(out, p.Documentation...)
	out = append(out, p.AI...)
	return uniqueSorted(out)
}

func (p Profile) JSON() string {
	data, _ := json.Marshal(p)
	return string(data)
}

func runGit(ctx context.Context, root string, args ...string) string {
	cmd := exec.CommandContext(ctx, "git", args...)
	if root != "" {
		cmd.Dir = root
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}

func inferTemplate(task string) TaskTemplate {
	lower := strings.ToLower(task)
	switch {
	case strings.Contains(lower, "bug") || strings.Contains(lower, "fix"):
		return ResolveTemplate("bug_fix")
	case strings.Contains(lower, "refactor"):
		return ResolveTemplate("refactor")
	case strings.Contains(lower, "review"):
		return ResolveTemplate("review")
	case strings.Contains(lower, "security"):
		return ResolveTemplate("security")
	case strings.Contains(lower, "doc"):
		return ResolveTemplate("docs")
	default:
		return ResolveTemplate("feature")
	}
}

func summarizeTask(task string) string {
	task = strings.TrimSpace(task)
	if len(task) <= 120 {
		return task
	}
	return task[:117] + "..."
}

func looksSensitive(value string) bool {
	lower := strings.ToLower(value)
	for _, marker := range []string{"api_key", "apikey", "secret", "password", "token=", "private_key"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func splitLines(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func prefixed(prefix string, items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, prefix+": "+item)
	}
	return out
}

func writeBullets(b *strings.Builder, items []string) {
	if len(items) == 0 {
		b.WriteString("- None\n")
		return
	}
	for _, item := range uniqueSorted(items) {
		b.WriteString("- ")
		b.WriteString(item)
		b.WriteByte('\n')
	}
}

func uniqueSorted(items []string) []string {
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
	sort.Strings(out)
	return out
}

func empty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
