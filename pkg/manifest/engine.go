package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Engine is the central configuration coordinator
type Engine struct {
	mu           sync.RWMutex
	manifests    map[string]*DeclarativeManifest // sourceName -> parsed manifest
	mergedCache  *DeclarativeManifest
	cacheInvalid bool
}

func NewEngine() *Engine {
	return &Engine{
		manifests:    make(map[string]*DeclarativeManifest),
		cacheInvalid: true,
	}
}

func (e *Engine) LoadManifest(sourceName, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var m DeclarativeManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("invalid manifest schema format for source %s: %w", sourceName, err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.manifests[sourceName] = &m
	e.cacheInvalid = true
	return nil
}

// RegisterMemoryManifest registers a manifest in memory directly (useful for tests and plugins)
func (e *Engine) RegisterMemoryManifest(sourceName string, m *DeclarativeManifest) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.manifests[sourceName] = m
	e.cacheInvalid = true
}

func (e *Engine) GetMergedManifest() *DeclarativeManifest {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.cacheInvalid && e.mergedCache != nil {
		return e.mergedCache
	}

	// Merge all loaded manifests with priority: Project > Plugin > Organization > Core
	merged := &DeclarativeManifest{
		SchemaVersion: "1",
		Workflows:     make(map[string]WorkflowDef),
		Standards:     make(map[string]StandardDef),
		Technologies:  make(map[string]TechDef),
		Prompts:       make(map[string]PromptDef),
		HealthRules:   make(map[string]HealthRuleDef),
	}

	// Order of precedence: compile from lowest priority to highest priority
	sources := []string{"core", "organization", "plugin", "project"}
	for _, source := range sources {
		if m, ok := e.manifests[source]; ok {
			e.merge(merged, m)
		}
	}

	// Also catch any unstructured plugin sources
	for source, m := range e.manifests {
		isStructured := false
		for _, s := range sources {
			if s == source {
				isStructured = true
				break
			}
		}
		if !isStructured {
			e.merge(merged, m)
		}
	}

	e.mergedCache = merged
	e.cacheInvalid = false
	return merged
}

func (e *Engine) merge(dest, src *DeclarativeManifest) {
	if src.SchemaVersion != "" {
		dest.SchemaVersion = src.SchemaVersion
	}
	if src.Compatibility.MinCLIVersion != "" || src.Compatibility.ManifestSchemaVer != "" {
		dest.Compatibility = src.Compatibility
	}

	for k, v := range src.Workflows {
		dest.Workflows[k] = v
	}
	for k, v := range src.Standards {
		dest.Standards[k] = v
	}
	for k, v := range src.Technologies {
		dest.Technologies[k] = v
	}
	for k, v := range src.Prompts {
		dest.Prompts[k] = v
	}
	for k, v := range src.HealthRules {
		dest.HealthRules[k] = v
	}
}
