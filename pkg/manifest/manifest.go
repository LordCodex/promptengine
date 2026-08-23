package manifest

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/LordCodex/promptengine/internal/filesystem"
	"gopkg.in/yaml.v3"
)

const DefaultFilename = "playbook-manifest.json"

type Playbook struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type TaskMapping struct {
	RequiredPlaybookIDs []string `json:"required_playbook_ids"`
	OptionalPlaybookIDs []string `json:"optional_playbook_ids"`
}

type PlaybookManifest struct {
	RepositoryNavigation map[string]string      `json:"repository_navigation"`
	CorePlaybooks        []Playbook             `json:"core_playbooks"`
	TechnologyStacks     map[string][]Playbook  `json:"technology_stacks"`
	DomainPlaybooks      map[string][]Playbook  `json:"domain_playbooks"`
	ProjectPlaybooks     []Playbook             `json:"project_playbooks"`
	BridgePlaybooks      []Playbook             `json:"bridge_playbooks"`
	Checklists           []Playbook             `json:"checklists"`
	Workflows            []Playbook             `json:"workflows"`
	DecisionGuides       []Playbook             `json:"decision_guides"`
	AIBootstrap          []Playbook             `json:"ai_bootstrap"`
	Guides               []Playbook             `json:"guides"`
	PromptsLibrary       []Playbook             `json:"prompts_library"`
	CliFoundation        []Playbook             `json:"cli_foundation"`
	CliCommandSpecs      []Playbook             `json:"cli_command_specifications"`
	TaskMappings         map[string]TaskMapping `json:"task_mappings"`
}

type Loader struct{ fs filesystem.FileSystem }

func NewLoader(fs filesystem.FileSystem) *Loader { return &Loader{fs: fs} }

func (l *Loader) Discover(startDir string) (string, bool) {
	if startDir == "" {
		startDir = "."
	}
	dir := filepath.Clean(startDir)
	for {
		candidate := filepath.Join(dir, DefaultFilename)
		if l.fs.Exists(candidate) {
			return filepath.ToSlash(candidate), true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func (l *Loader) Load(path string) (*Manifest, error) {
	if path == "" {
		path = DefaultFilename
	}
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load manifest %q: %w", path, err)
	}
	var m Manifest
	if err := decodeManifest(path, data, &m); err != nil {
		return nil, err
	}
	if m.Metadata.SchemaVersion == "" {
		var legacy PlaybookManifest
		if err := json.Unmarshal(data, &legacy); err == nil {
			m = convertLegacyManifest(legacy)
		}
	}
	if m.Metadata.GeneratedAt.IsZero() && m.Metadata.GeneratedAtRaw != "" {
		parsed, err := time.Parse(time.RFC3339, m.Metadata.GeneratedAtRaw)
		if err != nil {
			return nil, fmt.Errorf("manifest metadata generated_at is invalid: %w", err)
		}
		m.Metadata.GeneratedAt = parsed
	}
	return &m, nil
}

func convertLegacyManifest(old PlaybookManifest) Manifest {
	m := Manifest{Metadata: ProjectMetadata{Name: "PromptEngine", Version: "1.0.0", SchemaVersion: SupportedSchemaVersion, GeneratedAt: time.Now().UTC()}}
	seen := map[string]bool{}
	add := func(category PlaybookCategory, books []Playbook) {
		for _, book := range books {
			if book.ID == "" || book.Path == "" || seen[book.ID] {
				continue
			}
			seen[book.ID] = true
			m.Playbooks = append(m.Playbooks, PlaybookDefinition{ID: book.ID, Name: strings.ReplaceAll(book.ID, "_", " "), Category: category, Location: book.Path, Priority: 50})
		}
	}
	add(CategoryCore, old.CorePlaybooks)
	for domain, books := range old.DomainPlaybooks {
		add(domainCategory(domain), books)
	}
	add(CategoryProject, old.ProjectPlaybooks)
	add(CategoryBridge, old.BridgePlaybooks)
	add(CategoryChecklist, old.Checklists)
	add(CategoryWorkflows, old.Workflows)
	add(CategoryDecisionGuide, old.DecisionGuides)
	add(CategoryGuide, old.Guides)
	add(CategoryAI, old.AIBootstrap)
	add(CategoryPrompt, old.PromptsLibrary)
	add(CategoryCLI, old.CliFoundation)
	add(CategoryCLI, old.CliCommandSpecs)
	for stack, books := range old.TechnologyStacks {
		var related []string
		for _, book := range books {
			related = append(related, book.ID)
		}
		add(CategoryStacks, books)
		m.Technologies = append(m.Technologies, TechnologyDefinition{ID: stack, Stack: stack, RelatedPlaybooks: related})
	}
	for task, mapping := range old.TaskMappings {
		workflowID := normalizeTask(task)
		m.Workflows = append(m.Workflows, WorkflowDefinition{ID: workflowID, Steps: []string{"prepare", "execute", "review"}, RequiredPlaybooks: mapping.RequiredPlaybookIDs, OptionalPlaybooks: mapping.OptionalPlaybookIDs})
		m.TaskRelationships = append(m.TaskRelationships, TaskRelationship{TaskType: task, RequiredWorkflow: workflowID})
	}
	return m
}

func domainCategory(domain string) PlaybookCategory {
	switch strings.ToLower(domain) {
	case "security":
		return CategorySecurity
	case "performance":
		return CategoryPerformance
	case "ui_ux", "accessibility", "seo":
		return CategoryDesign
	default:
		return CategoryCore
	}
}

func (l *Loader) LoadDiscovered(startDir string) (string, *Manifest, error) {
	path, ok := l.Discover(startDir)
	if !ok {
		return "", nil, fmt.Errorf("manifest %q was not found from %q", DefaultFilename, startDir)
	}
	m, err := l.Load(path)
	return path, m, err
}

func decodeManifest(path string, data []byte, m *Manifest) error {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, m); err != nil {
			return fmt.Errorf("parse manifest %q as YAML: %w", path, err)
		}
	default:
		if err := json.Unmarshal(data, m); err != nil {
			return fmt.Errorf("parse manifest %q as JSON: %w", path, err)
		}
	}
	return nil
}

// Load preserves the old disk-loading API for packages that have not moved to Loader.
func Load(path string) (*PlaybookManifest, error) {
	fs := &filesystem.OSFileSystem{}
	data, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m PlaybookManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
