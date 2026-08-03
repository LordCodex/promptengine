package manifest

import (
	"fmt"
	"sync"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

type SourceKind string

const (
	SourceCore         SourceKind = "core"
	SourceOrganization SourceKind = "organization"
	SourceMarketplace  SourceKind = "marketplace"
	SourcePlugin       SourceKind = "plugin"
	SourceProject      SourceKind = "project"
)

type Source struct {
	Name     string
	Kind     SourceKind
	Priority int
	Manifest *Manifest
}

type Engine struct {
	mu      sync.RWMutex
	fs      filesystem.FileSystem
	sources map[string]Source
	merged  *Manifest
	dirty   bool

	// Legacy registry retained for existing tests/packages during rollout.
	legacyManifests map[string]*DeclarativeManifest
	legacyMerged    *DeclarativeManifest
	legacyDirty     bool
}

func NewEngine() *Engine {
	return NewEngineWithFS(&filesystem.OSFileSystem{})
}

func NewEngineWithFS(fs filesystem.FileSystem) *Engine {
	if fs == nil {
		fs = &filesystem.OSFileSystem{}
	}
	return &Engine{
		fs:              fs,
		sources:         map[string]Source{},
		dirty:           true,
		legacyManifests: map[string]*DeclarativeManifest{},
		legacyDirty:     true,
	}
}

func (e *Engine) LoadManifest(sourceName, path string) error {
	loader := NewLoader(e.fs)
	m, err := loader.Load(path)
	if err != nil {
		return err
	}
	return e.Register(sourceName, SourceProject, m)
}

func (e *Engine) Register(sourceName string, kind SourceKind, m *Manifest) error {
	if sourceName == "" {
		return fmt.Errorf("manifest source name is required")
	}
	if m == nil {
		return fmt.Errorf("manifest source %q is nil", sourceName)
	}
	if err := Validate(m, e.fs); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.sources[sourceName] = Source{Name: sourceName, Kind: kind, Priority: sourcePriority(kind), Manifest: m}
	e.dirty = true
	return nil
}

func (e *Engine) RegisterPluginManifest(pluginID string, m *Manifest) error {
	return e.Register(pluginID, SourcePlugin, m)
}

func (e *Engine) Sources() []Source {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]Source, 0, len(e.sources))
	for _, source := range e.sources {
		out = append(out, source)
	}
	return out
}

func (e *Engine) ActiveManifest() *Manifest {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.dirty && e.merged != nil {
		return cloneManifest(e.merged)
	}

	merged := &Manifest{
		Metadata:   ProjectMetadata{SchemaVersion: SupportedSchemaVersion},
		PluginData: map[string]map[string]any{},
		Extensions: map[string][]ExtensionResource{},
	}
	for _, source := range orderedSources(e.sources) {
		mergeManifest(merged, source.Manifest)
	}
	e.merged = merged
	e.dirty = false
	return cloneManifest(merged)
}

func sourcePriority(kind SourceKind) int {
	switch kind {
	case SourceCore:
		return 10
	case SourceOrganization:
		return 20
	case SourceMarketplace:
		return 30
	case SourcePlugin:
		return 40
	case SourceProject:
		return 50
	default:
		return 0
	}
}

func orderedSources(sources map[string]Source) []Source {
	out := make([]Source, 0, len(sources))
	for _, source := range sources {
		out = append(out, source)
	}
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Priority < out[i].Priority || (out[j].Priority == out[i].Priority && out[j].Name < out[i].Name) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

func mergeManifest(dest, src *Manifest) {
	if src.Metadata.Name != "" {
		dest.Metadata = src.Metadata
	}
	dest.Technologies = upsertTech(dest.Technologies, src.Technologies)
	dest.Playbooks = upsertPlaybooks(dest.Playbooks, src.Playbooks)
	dest.Workflows = upsertWorkflows(dest.Workflows, src.Workflows)
	dest.Prompts = upsertPrompts(dest.Prompts, src.Prompts)
	dest.Templates = upsertTemplates(dest.Templates, src.Templates)
	dest.CommandMappings = append(dest.CommandMappings, src.CommandMappings...)
	dest.TaskRelationships = append(dest.TaskRelationships, src.TaskRelationships...)
	if dest.PluginData == nil {
		dest.PluginData = map[string]map[string]any{}
	}
	for k, v := range src.PluginData {
		dest.PluginData[k] = v
	}
	if dest.Extensions == nil {
		dest.Extensions = map[string][]ExtensionResource{}
	}
	for k, v := range src.Extensions {
		dest.Extensions[k] = append(dest.Extensions[k], v...)
	}
}

func cloneManifest(m *Manifest) *Manifest {
	if m == nil {
		return nil
	}
	c := *m
	c.Technologies = append([]TechnologyDefinition(nil), m.Technologies...)
	c.Playbooks = append([]PlaybookDefinition(nil), m.Playbooks...)
	c.Workflows = append([]WorkflowDefinition(nil), m.Workflows...)
	c.Prompts = append([]PromptMapping(nil), m.Prompts...)
	c.Templates = append([]TemplateDefinition(nil), m.Templates...)
	c.CommandMappings = append([]CommandMapping(nil), m.CommandMappings...)
	c.TaskRelationships = append([]TaskRelationship(nil), m.TaskRelationships...)
	return &c
}
