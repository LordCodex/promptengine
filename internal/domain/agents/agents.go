package agents

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	contextengine "github.com/LordCodex/promptengine/internal/domain/context"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
	"gopkg.in/yaml.v3"
)

type AgentProfile struct {
	ID              string `json:"id" yaml:"id"`
	Name            string `json:"name" yaml:"name"`
	InstructionFile string `json:"instruction_file" yaml:"instruction_file"`
	Format          string `json:"format" yaml:"format"`
}

type InstructionRequest struct {
	Profile     string
	Project     *discovery.ProjectModel
	Manifest    *manifest.Manifest
	Preferences []string
	Config      map[string]string
}

type GeneratedConfig struct {
	Profile string `json:"profile" yaml:"profile"`
	Path    string `json:"path" yaml:"path"`
	Status  string `json:"status" yaml:"status"`
	Content string `json:"-" yaml:"-"`
}

type SyncReport struct {
	Generated []GeneratedConfig `json:"generated" yaml:"generated"`
	Updated   []GeneratedConfig `json:"updated" yaml:"updated"`
	Current   []GeneratedConfig `json:"current" yaml:"current"`
}

type ContextExportRequest struct {
	Task    string
	Agent   string
	Format  string
	Package *contextengine.ContextPackage
}

type ContextExport struct {
	File                  string   `json:"file" yaml:"file"`
	Format                string   `json:"format" yaml:"format"`
	Agent                 string   `json:"agent" yaml:"agent"`
	SelectedFiles         []string `json:"selected_files" yaml:"selected_files"`
	SelectedDocuments     []string `json:"selected_documents" yaml:"selected_documents"`
	EstimatedContextBytes int      `json:"estimated_context_bytes" yaml:"estimated_context_bytes"`
}

type Platform struct {
	fs       filesystem.FileSystem
	events   *eventbus.EventBus
	profiles map[string]AgentProfile
}

func NewPlatform(fs filesystem.FileSystem, events *eventbus.EventBus) *Platform {
	p := &Platform{fs: fs, events: events, profiles: map[string]AgentProfile{}}
	for _, profile := range defaultProfiles() {
		p.Register(profile)
	}
	return p
}

func (p *Platform) Register(profile AgentProfile) {
	if profile.ID == "" {
		return
	}
	p.profiles[profile.ID] = profile
}

func (p *Platform) Profiles() []AgentProfile {
	out := make([]AgentProfile, 0, len(p.profiles))
	for _, profile := range p.profiles {
		out = append(out, profile)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (p *Platform) Generate(req InstructionRequest) (GeneratedConfig, error) {
	profile, ok := p.profile(req.Profile)
	if !ok {
		return GeneratedConfig{}, fmt.Errorf("agent profile %q is not registered", req.Profile)
	}
	content := renderInstruction(profile, req.Project, req.Manifest, req.Preferences)
	if err := p.fs.WriteFile(profile.InstructionFile, []byte(content), 0644); err != nil {
		return GeneratedConfig{}, err
	}
	result := GeneratedConfig{Profile: profile.ID, Path: profile.InstructionFile, Status: "generated", Content: content}
	p.publish(eventbus.AgentConfigGenerated, "agent config generated", result)
	return result, nil
}

func (p *Platform) Sync(req InstructionRequest) (SyncReport, error) {
	var report SyncReport
	targets := p.Profiles()
	if req.Profile != "" && req.Profile != "all" {
		profile, ok := p.profile(req.Profile)
		if !ok {
			return report, fmt.Errorf("agent profile %q is not registered", req.Profile)
		}
		targets = []AgentProfile{profile}
	}
	for _, profile := range targets {
		content := renderInstruction(profile, req.Project, req.Manifest, req.Preferences)
		status := "current"
		if !p.fs.Exists(profile.InstructionFile) {
			status = "generated"
		} else if data, err := p.fs.ReadFile(profile.InstructionFile); err == nil && string(data) != content {
			status = "updated"
		}
		result := GeneratedConfig{Profile: profile.ID, Path: profile.InstructionFile, Status: status, Content: content}
		switch status {
		case "current":
			report.Current = append(report.Current, result)
		case "updated":
			if err := p.fs.WriteFile(profile.InstructionFile, []byte(content), 0644); err != nil {
				return report, err
			}
			report.Updated = append(report.Updated, result)
			p.publish(eventbus.AgentConfigUpdated, "agent config updated", result)
		default:
			if err := p.fs.WriteFile(profile.InstructionFile, []byte(content), 0644); err != nil {
				return report, err
			}
			report.Generated = append(report.Generated, result)
			p.publish(eventbus.AgentConfigGenerated, "agent config generated", result)
		}
	}
	return report, nil
}

func (p *Platform) ExportContext(req ContextExportRequest) (ContextExport, error) {
	if req.Package == nil {
		return ContextExport{}, fmt.Errorf("context package is required")
	}
	format, ext := normalizeFormat(req.Format)
	agent := req.Agent
	if agent == "" {
		agent = "generic"
	}
	path := defaultContextPath(req.Task, agent, ext)
	content, err := renderContext(req.Package, agent, format)
	if err != nil {
		return ContextExport{}, err
	}
	if err := p.fs.WriteFile(path, content, 0644); err != nil {
		return ContextExport{}, err
	}
	result := ContextExport{
		File:                  path,
		Format:                format,
		Agent:                 agent,
		SelectedFiles:         append([]string(nil), req.Package.SelectedFiles...),
		SelectedDocuments:     append([]string(nil), req.Package.SelectedDocs...),
		EstimatedContextBytes: req.Package.Summary.FinalSize,
	}
	if result.EstimatedContextBytes == 0 {
		result.EstimatedContextBytes = len(req.Package.SystemPrompt)
	}
	p.publish(eventbus.ContextExported, "context exported", result)
	return result, nil
}

func (p *Platform) profile(id string) (AgentProfile, bool) {
	if id == "" {
		id = "codex"
	}
	profile, ok := p.profiles[id]
	return profile, ok
}

func (p *Platform) publish(t eventbus.EventType, msg string, payload interface{}) {
	if p.events != nil {
		p.events.Publish(eventbus.Event{Type: t, Message: msg, Payload: payload})
	}
}

func defaultProfiles() []AgentProfile {
	return []AgentProfile{
		{ID: "codex", Name: "Codex", InstructionFile: "AGENTS.md", Format: "markdown"},
		{ID: "claude", Name: "Claude Code", InstructionFile: "CLAUDE.md", Format: "markdown"},
		{ID: "codex-md", Name: "Codex Instructions", InstructionFile: "CODEX.md", Format: "markdown"},
		{ID: "cursor", Name: "Cursor", InstructionFile: ".cursor/rules/promptengine.md", Format: "markdown"},
		{ID: "windsurf", Name: "Windsurf", InstructionFile: ".windsurf/rules/promptengine.md", Format: "markdown"},
	}
}

func renderInstruction(profile AgentProfile, project *discovery.ProjectModel, m *manifest.Manifest, preferences []string) string {
	var b strings.Builder
	b.WriteString("# ")
	b.WriteString(profile.Name)
	b.WriteString(" Instructions\n\n")
	b.WriteString("Generated by PromptEngine. Keep this file synchronized with `promptengine agents sync`.\n\n")
	b.WriteString("## PromptEngine Location\n\n")
	b.WriteString("- Manifest: `playbook-manifest.json`\n")
	b.WriteString("- Project docs: `docs/`\n")
	b.WriteString("- Standards root: repository PromptEngine playbooks referenced by the manifest\n\n")
	b.WriteString("## Project Context\n\n")
	if project != nil {
		b.WriteString("- Root: `")
		b.WriteString(project.RootDir)
		b.WriteString("`\n")
		b.WriteString("- Languages: ")
		b.WriteString(joinList(project.Languages))
		b.WriteString("\n- Frameworks: ")
		b.WriteString(joinList(project.Frameworks))
		b.WriteString("\n")
	} else {
		b.WriteString("- Project discovery has not been run yet.\n")
	}
	b.WriteString("\n## Workflow Rules\n\n")
	b.WriteString("- Use PromptEngine context before editing code.\n")
	b.WriteString("- Follow the active workflow for the task type.\n")
	b.WriteString("- Keep changes scoped and preserve existing architecture.\n")
	b.WriteString("- Report verification steps before finalizing work.\n\n")
	b.WriteString("## Engineering Standards\n\n")
	for _, ref := range standardRefs(m) {
		b.WriteString("- `")
		b.WriteString(ref)
		b.WriteString("`\n")
	}
	b.WriteString("\n## Documentation Rules\n\n")
	b.WriteString("- Update docs when behavior, architecture, APIs, schemas, deployment, or troubleshooting guidance changes.\n")
	b.WriteString("- Prefer references to PromptEngine standards over copying long rule text.\n\n")
	if len(preferences) > 0 {
		b.WriteString("## Developer Preferences\n\n")
		writeBullets(&b, preferences)
		b.WriteString("\n")
	}
	b.WriteString("## Security Rules\n\n")
	b.WriteString("- Do not expose secrets in code, logs, reports, generated prompts, or examples.\n")
	b.WriteString("- Treat authentication, authorization, data validation, and dependency changes as security-sensitive.\n\n")
	b.WriteString("## Testing Expectations\n\n")
	b.WriteString("- Run the narrowest useful verification first.\n")
	b.WriteString("- Broaden test coverage when shared contracts or user-facing flows change.\n\n")
	b.WriteString("## Agent Formatting\n\n")
	switch profile.ID {
	case "claude":
		b.WriteString("- Claude should explain assumptions and keep implementation plans concise.\n")
	case "cursor":
		b.WriteString("- Cursor should apply these rules to indexed workspace edits.\n")
	case "windsurf":
		b.WriteString("- Windsurf should use these instructions as workspace-level engineering rules.\n")
	default:
		b.WriteString("- Coding agents should use this file as their first project instruction source.\n")
	}
	return b.String()
}

func standardRefs(m *manifest.Manifest) []string {
	refs := []string{"core/05-universal-coding-standards.md", "core/11-testing-engineering-standard.md", "core/19-documentation-engineering-standard.md", "core/20-ai-agent-engineering-workflow-standard.md"}
	if m != nil {
		for _, playbook := range m.Playbooks {
			if playbook.Location != "" {
				refs = append(refs, playbook.Location)
			}
		}
	}
	return uniqueSorted(refs)
}

func renderContext(pkg *contextengine.ContextPackage, agent, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(pkg, "", "  ")
	case "markdown":
		var b strings.Builder
		b.WriteString("# ")
		b.WriteString(strings.Title(agent))
		b.WriteString(" Context Package\n\n")
		b.WriteString("Generated: ")
		b.WriteString(time.Now().UTC().Format(time.RFC3339))
		b.WriteString("\n\n## Project Context\n\n")
		b.WriteString(pkg.SystemPrompt)
		b.WriteString("\n\n## Relevant Files\n\n")
		writeBullets(&b, pkg.SelectedFiles)
		b.WriteString("\n## Relevant Documents\n\n")
		writeBullets(&b, pkg.SelectedDocs)
		b.WriteString("\n## Standards\n\n")
		writeBullets(&b, pkg.RelevantStandards)
		b.WriteString("\n## Reasoning\n\n")
		writeBullets(&b, pkg.Reasoning)
		b.WriteString(fmt.Sprintf("\nEstimated context size: %d bytes\n", pkg.Summary.FinalSize))
		return []byte(b.String()), nil
	case "yaml":
		return yaml.Marshal(pkg)
	default:
		return nil, fmt.Errorf("unsupported context export format %q", format)
	}
}

func normalizeFormat(format string) (string, string) {
	switch strings.ToLower(format) {
	case "json":
		return "json", "json"
	case "yaml", "yml":
		return "yaml", "yaml"
	default:
		return "markdown", "md"
	}
}

func defaultContextPath(task, agent, ext string) string {
	task = strings.TrimSpace(strings.ToLower(task))
	if task == "" {
		task = "task"
	}
	task = strings.NewReplacer(" ", "-", "_", "-", "/", "-").Replace(task)
	if agent == "" || agent == "generic" {
		return task + "-context." + ext
	}
	return filepath.Clean(agent + "-context." + ext)
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

func joinList(items []string) string {
	if len(items) == 0 {
		return "none"
	}
	return strings.Join(uniqueSorted(items), ", ")
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
