package plugins

import (
	"fmt"
	"sync"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// PluginMetadata describes an installed plugin
type PluginMetadata struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Author       string   `json:"author"`
	Version      string   `json:"version"`
	Enabled      bool     `json:"enabled"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies"` // IDs of required plugins
	MinCoreVer   string   `json:"min_core_ver"` // minimum PromptEngine version
}

// ContributionPoints declares what a plugin contributes to PromptEngine
type ContributionPoints struct {
	Commands    []string // CLI command IDs
	Workflows   []string // workflow IDs contributed
	Standards   []string // standard IDs contributed
	Detectors   []string // technology detector IDs
	Validators  []string // validator IDs
	AIProviders []string // AI provider IDs
	HookTypes   []string // hook type IDs
	Templates   []string // template IDs
}

// Plugin is the full extension boundary interface
type Plugin interface {
	Metadata() PluginMetadata
	Contributions() ContributionPoints
	OnInstall(fs filesystem.FileSystem) error
	OnUninstall(fs filesystem.FileSystem) error
	OnEnable() error
	OnDisable() error
	OnUpdate(fromVersion string) error
}

// CompatibilityResult holds version check outcomes
type CompatibilityResult struct {
	Compatible bool
	Reason     string
}

// CheckCompatibility verifies a plugin's core version requirement
func CheckCompatibility(meta PluginMetadata, currentCoreVer string) CompatibilityResult {
	if meta.MinCoreVer == "" {
		return CompatibilityResult{Compatible: true, Reason: "no minimum version declared"}
	}
	if currentCoreVer < meta.MinCoreVer {
		return CompatibilityResult{
			Compatible: false,
			Reason:     fmt.Sprintf("plugin '%s' requires core >= %s, got %s", meta.ID, meta.MinCoreVer, currentCoreVer),
		}
	}
	return CompatibilityResult{Compatible: true, Reason: "compatible"}
}

// Registry manages plugin discovery, activation, and dependency ordering
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
}

func NewRegistry() *Registry {
	return &Registry{plugins: make(map[string]Plugin)}
}

func (r *Registry) Register(p Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := p.Metadata().ID
	if _, exists := r.plugins[id]; exists {
		return fmt.Errorf("plugin '%s' already registered", id)
	}
	r.plugins[id] = p
	return nil
}

func (r *Registry) Get(id string) (Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plugins[id]
	return p, ok
}

func (r *Registry) Enable(id string) error {
	r.mu.RLock()
	p, ok := r.plugins[id]
	r.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin '%s' not found", id)
	}
	// Resolve dependencies first
	for _, dep := range p.Metadata().Dependencies {
		if _, depOK := r.Get(dep); !depOK {
			return fmt.Errorf("dependency '%s' required by plugin '%s' is not registered", dep, id)
		}
	}
	return p.OnEnable()
}

func (r *Registry) Disable(id string) error {
	r.mu.RLock()
	p, ok := r.plugins[id]
	r.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin '%s' not found", id)
	}
	return p.OnDisable()
}

func (r *Registry) List() []PluginMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []PluginMetadata
	for _, p := range r.plugins {
		list = append(list, p.Metadata())
	}
	return list
}

// ResolveLoadOrder returns plugins in dependency-safe topological order
func (r *Registry) ResolveLoadOrder() ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	var order []string

	var visit func(id string) error
	visit = func(id string) error {
		if inStack[id] {
			return fmt.Errorf("cyclic dependency detected involving plugin '%s'", id)
		}
		if visited[id] {
			return nil
		}
		inStack[id] = true
		p, ok := r.plugins[id]
		if !ok {
			return fmt.Errorf("plugin '%s' not registered but required as dependency", id)
		}
		for _, dep := range p.Metadata().Dependencies {
			if err := visit(dep); err != nil {
				return err
			}
		}
		inStack[id] = false
		visited[id] = true
		order = append(order, id)
		return nil
	}

	for id := range r.plugins {
		if err := visit(id); err != nil {
			return nil, err
		}
	}
	return order, nil
}

// Loader is the legacy adapter kept for backward compatibility
type Loader struct {
	registry *Registry
	fs       filesystem.FileSystem
}

func NewLoader(fs filesystem.FileSystem) *Loader {
	return &Loader{fs: fs, registry: NewRegistry()}
}

func (l *Loader) Register(p Plugin) {
	_ = l.registry.Register(p) // swallow duplicate errors in legacy path
}

func (l *Loader) List() []PluginMetadata {
	return l.registry.List()
}
