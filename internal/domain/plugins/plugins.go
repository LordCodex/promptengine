package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/internal/version"
	"github.com/LordCodex/promptengine/pkg/manifest"
	"gopkg.in/yaml.v3"
)

type Status string

const (
	StatusInstalled Status = "installed"
	StatusEnabled   Status = "enabled"
	StatusDisabled  Status = "disabled"
	StatusFailed    Status = "failed"
)

type Permission string

const (
	PermissionReadProject  Permission = "read_project"
	PermissionWriteProject Permission = "write_project"
	PermissionNetwork      Permission = "network"
	PermissionRunHooks     Permission = "run_hooks"
)

type PluginMetadata struct {
	ID                   string       `json:"id" yaml:"id"`
	Name                 string       `json:"name" yaml:"name"`
	Author               string       `json:"author,omitempty" yaml:"author,omitempty"`
	Version              string       `json:"version" yaml:"version"`
	Description          string       `json:"description,omitempty" yaml:"description,omitempty"`
	CompatibilityVersion string       `json:"compatibility_version,omitempty" yaml:"compatibility_version,omitempty"`
	Dependencies         []string     `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Permissions          []Permission `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	Status               Status       `json:"status" yaml:"status"`

	Enabled    bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	MinCoreVer string   `json:"min_core_ver,omitempty" yaml:"min_core_ver,omitempty"`
	Files      []string `json:"files,omitempty" yaml:"files,omitempty"`
}

type ContributionPoints struct {
	Commands        []string               `json:"commands,omitempty" yaml:"commands,omitempty"`
	Workflows       []string               `json:"workflows,omitempty" yaml:"workflows,omitempty"`
	Standards       []string               `json:"standards,omitempty" yaml:"standards,omitempty"`
	Detectors       []string               `json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Validators      []string               `json:"validators,omitempty" yaml:"validators,omitempty"`
	QualityRules    []quality.Rule         `json:"-" yaml:"-"`
	AIProviders     []string               `json:"ai_providers,omitempty" yaml:"ai_providers,omitempty"`
	DiscoveryStages []string               `json:"discovery_stages,omitempty" yaml:"discovery_stages,omitempty"`
	HookTypes       []string               `json:"hook_types,omitempty" yaml:"hook_types,omitempty"`
	Templates       []string               `json:"templates,omitempty" yaml:"templates,omitempty"`
	Manifest        *manifest.Manifest     `json:"manifest,omitempty" yaml:"manifest,omitempty"`
	CustomEngines   map[string]interface{} `json:"-" yaml:"-"`
}

type HealthFinding struct {
	PluginID          string `json:"plugin_id" yaml:"plugin_id"`
	Severity          string `json:"severity" yaml:"severity"`
	Message           string `json:"message" yaml:"message"`
	RecommendedAction string `json:"recommended_action,omitempty" yaml:"recommended_action,omitempty"`
}

type Plugin interface {
	Metadata() PluginMetadata
	Contributions() ContributionPoints
	Install(ctx context.Context, fs filesystem.FileSystem) error
	Load(ctx context.Context, fs filesystem.FileSystem) error
	Enable(ctx context.Context) error
	Disable(ctx context.Context) error
	Upgrade(ctx context.Context, fromVersion string) error
	Validate(ctx context.Context, fs filesystem.FileSystem) error
	HealthCheck(ctx context.Context, fs filesystem.FileSystem) []HealthFinding
	Unload(ctx context.Context) error
	Remove(ctx context.Context, fs filesystem.FileSystem) error
}

type ManifestPlugin struct {
	meta PluginMetadata
}

func NewManifestPlugin(meta PluginMetadata) *ManifestPlugin {
	return &ManifestPlugin{meta: meta}
}

func (p *ManifestPlugin) Metadata() PluginMetadata          { return p.meta }
func (p *ManifestPlugin) Contributions() ContributionPoints { return ContributionPoints{} }
func (p *ManifestPlugin) Install(ctx context.Context, fs filesystem.FileSystem) error {
	return nil
}
func (p *ManifestPlugin) Load(ctx context.Context, fs filesystem.FileSystem) error { return nil }
func (p *ManifestPlugin) Enable(ctx context.Context) error {
	p.meta.Enabled = true
	p.meta.Status = StatusEnabled
	return nil
}
func (p *ManifestPlugin) Disable(ctx context.Context) error {
	p.meta.Enabled = false
	p.meta.Status = StatusDisabled
	return nil
}
func (p *ManifestPlugin) Upgrade(ctx context.Context, fromVersion string) error { return nil }
func (p *ManifestPlugin) Validate(ctx context.Context, fs filesystem.FileSystem) error {
	return nil
}
func (p *ManifestPlugin) HealthCheck(ctx context.Context, fs filesystem.FileSystem) []HealthFinding {
	return nil
}
func (p *ManifestPlugin) Unload(ctx context.Context) error { return nil }
func (p *ManifestPlugin) Remove(ctx context.Context, fs filesystem.FileSystem) error {
	p.meta.Enabled = false
	p.meta.Status = StatusDisabled
	return nil
}

type CompatibilityResult struct {
	Compatible bool
	Reason     string
}

func CheckCompatibility(meta PluginMetadata, currentCoreVer string) CompatibilityResult {
	required := meta.CompatibilityVersion
	if required == "" {
		required = meta.MinCoreVer
	}
	if required == "" {
		return CompatibilityResult{Compatible: true, Reason: "no minimum version declared"}
	}
	if currentCoreVer < required {
		return CompatibilityResult{Compatible: false, Reason: fmt.Sprintf("plugin '%s' requires core >= %s, got %s", meta.ID, required, currentCoreVer)}
	}
	return CompatibilityResult{Compatible: true, Reason: "compatible"}
}

type ExtensionRegistry struct {
	Workflows       map[string]interface{}
	QualityRules    map[string]quality.Rule
	DocTemplates    map[string]string
	AIProviders     map[string]interface{}
	DiscoveryStages map[string]interface{}
	Commands        map[string]interface{}
	Hooks           map[string]interface{}
	CustomEngines   map[string]interface{}
}

func newExtensionRegistry() ExtensionRegistry {
	return ExtensionRegistry{
		Workflows:       map[string]interface{}{},
		QualityRules:    map[string]quality.Rule{},
		DocTemplates:    map[string]string{},
		AIProviders:     map[string]interface{}{},
		DiscoveryStages: map[string]interface{}{},
		Commands:        map[string]interface{}{},
		Hooks:           map[string]interface{}{},
		CustomEngines:   map[string]interface{}{},
	}
}

type Registry struct {
	mu                 sync.RWMutex
	plugins            map[string]Plugin
	status             map[string]Status
	config             map[string]map[string]interface{}
	extensions         ExtensionRegistry
	events             *eventbus.EventBus
	coreVersion        string
	allowedPermissions map[Permission]bool
}

func NewRegistry() *Registry {
	return NewRegistryWithEvents(nil)
}

func NewRegistryWithEvents(events *eventbus.EventBus) *Registry {
	return &Registry{
		plugins:            map[string]Plugin{},
		status:             map[string]Status{},
		config:             map[string]map[string]interface{}{},
		extensions:         newExtensionRegistry(),
		events:             events,
		coreVersion:        version.Version,
		allowedPermissions: map[Permission]bool{PermissionReadProject: true, PermissionWriteProject: true, PermissionRunHooks: true},
	}
}

func (r *Registry) Register(p Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	meta := p.Metadata()
	if meta.ID == "" {
		return fmt.Errorf("plugin id is required")
	}
	if _, exists := r.plugins[meta.ID]; exists {
		return fmt.Errorf("plugin '%s' already registered", meta.ID)
	}
	if err := r.validateMetadataLocked(meta); err != nil {
		r.status[meta.ID] = StatusFailed
		return err
	}
	r.plugins[meta.ID] = p
	r.status[meta.ID] = firstStatus(meta)
	r.registerContributionsLocked(meta.ID, p.Contributions())
	return nil
}

func (r *Registry) Get(id string) (Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plugins[id]
	return p, ok
}

func (r *Registry) Status(id string) Status {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.status[id]
}

func (r *Registry) Extensions() ExtensionRegistry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.extensions
}

func (r *Registry) Configure(id string, cfg map[string]interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.config[id] = cfg
}

func (r *Registry) Config(id string) map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := map[string]interface{}{}
	for k, v := range r.config[id] {
		out[k] = v
	}
	return out
}

func (r *Registry) Install(ctx context.Context, id string, fs filesystem.FileSystem) error {
	p, err := r.plugin(id)
	if err != nil {
		return err
	}
	if err := p.Install(ctx, fs); err != nil {
		r.setFailed(id)
		return err
	}
	r.setStatus(id, StatusInstalled)
	r.publish(eventbus.PluginInstalled, "plugin installed", safePluginMeta(p.Metadata(), StatusInstalled))
	return nil
}

func (r *Registry) Load(ctx context.Context, id string, fs filesystem.FileSystem) error {
	p, err := r.plugin(id)
	if err != nil {
		return err
	}
	if err := p.Load(ctx, fs); err != nil {
		r.setFailed(id)
		return err
	}
	return nil
}

func (r *Registry) Enable(id string) error {
	return r.EnableContext(context.Background(), id)
}

func (r *Registry) EnableContext(ctx context.Context, id string) error {
	p, err := r.plugin(id)
	if err != nil {
		return err
	}
	for _, dep := range p.Metadata().Dependencies {
		if _, depOK := r.Get(dep); !depOK {
			return fmt.Errorf("dependency '%s' required by plugin '%s' is not registered", dep, id)
		}
		if r.Status(dep) != StatusEnabled {
			if err := r.EnableContext(ctx, dep); err != nil {
				return err
			}
		}
	}
	if err := p.Enable(ctx); err != nil {
		r.setFailed(id)
		return err
	}
	r.setStatus(id, StatusEnabled)
	r.publish(eventbus.PluginEnabled, "plugin enabled", safePluginMeta(p.Metadata(), StatusEnabled))
	return nil
}

func (r *Registry) Disable(id string) error {
	p, err := r.plugin(id)
	if err != nil {
		return err
	}
	if err := p.Disable(context.Background()); err != nil {
		r.setFailed(id)
		return err
	}
	r.setStatus(id, StatusDisabled)
	r.publish(eventbus.PluginDisabled, "plugin disabled", safePluginMeta(p.Metadata(), StatusDisabled))
	return nil
}

func (r *Registry) Upgrade(ctx context.Context, id, fromVersion string) error {
	p, err := r.plugin(id)
	if err != nil {
		return err
	}
	if err := p.Upgrade(ctx, fromVersion); err != nil {
		r.setFailed(id)
		return err
	}
	r.publish(eventbus.PluginUpdated, "plugin updated", safePluginMeta(p.Metadata(), r.Status(id)))
	return nil
}

func (r *Registry) Remove(ctx context.Context, id string, fs filesystem.FileSystem) error {
	p, err := r.plugin(id)
	if err != nil {
		return err
	}
	if err := p.Remove(ctx, fs); err != nil {
		r.setFailed(id)
		return err
	}
	r.mu.Lock()
	delete(r.plugins, id)
	delete(r.status, id)
	delete(r.config, id)
	r.mu.Unlock()
	r.publish(eventbus.PluginRemoved, "plugin removed", map[string]string{"id": id})
	return nil
}

func (r *Registry) Health(ctx context.Context, fs filesystem.FileSystem) []HealthFinding {
	r.mu.RLock()
	plugins := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}
	r.mu.RUnlock()
	var findings []HealthFinding
	for _, p := range plugins {
		meta := p.Metadata()
		if err := p.Validate(ctx, fs); err != nil {
			findings = append(findings, HealthFinding{PluginID: meta.ID, Severity: "error", Message: err.Error(), RecommendedAction: "Fix plugin configuration or disable the plugin."})
		}
		findings = append(findings, p.HealthCheck(ctx, fs)...)
	}
	if len(findings) > 0 {
		r.publish(eventbus.PluginHealthFailed, "plugin health failed", findings)
	}
	return findings
}

func (r *Registry) ExtendManifest(base *manifest.Manifest) *manifest.Manifest {
	if base == nil {
		base = &manifest.Manifest{}
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.plugins {
		contrib := p.Contributions()
		if contrib.Manifest != nil {
			base.Playbooks = append(base.Playbooks, contrib.Manifest.Playbooks...)
			base.Workflows = append(base.Workflows, contrib.Manifest.Workflows...)
			base.Templates = append(base.Templates, contrib.Manifest.Templates...)
			base.Prompts = append(base.Prompts, contrib.Manifest.Prompts...)
			base.Technologies = append(base.Technologies, contrib.Manifest.Technologies...)
		}
	}
	return base
}

func (r *Registry) List() []PluginMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []PluginMetadata
	for id, p := range r.plugins {
		meta := p.Metadata()
		meta.Status = r.status[id]
		meta.Enabled = meta.Status == StatusEnabled
		list = append(list, meta)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

func (r *Registry) ResolveLoadOrder() ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	visited := map[string]bool{}
	inStack := map[string]bool{}
	var order []string
	var visit func(string) error
	visit = func(id string) error {
		if inStack[id] {
			return fmt.Errorf("cyclic dependency detected involving plugin '%s'", id)
		}
		if visited[id] {
			return nil
		}
		p, ok := r.plugins[id]
		if !ok {
			return fmt.Errorf("plugin '%s' not registered but required as dependency", id)
		}
		inStack[id] = true
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
	var ids []string
	for id := range r.plugins {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if err := visit(id); err != nil {
			return nil, err
		}
	}
	return order, nil
}

func (r *Registry) validateMetadataLocked(meta PluginMetadata) error {
	if meta.Version == "" {
		return fmt.Errorf("plugin '%s' version is required", meta.ID)
	}
	if result := CheckCompatibility(meta, r.coreVersion); !result.Compatible {
		return fmt.Errorf("%s", result.Reason)
	}
	for _, permission := range meta.Permissions {
		if !r.allowedPermissions[permission] {
			return fmt.Errorf("plugin '%s' requests disallowed permission '%s'", meta.ID, permission)
		}
	}
	return nil
}

func (r *Registry) registerContributionsLocked(pluginID string, c ContributionPoints) {
	for _, id := range c.Workflows {
		r.extensions.Workflows[id] = pluginID
	}
	for _, rule := range c.QualityRules {
		r.extensions.QualityRules[rule.ID()] = rule
	}
	for _, id := range c.Templates {
		r.extensions.DocTemplates[id] = pluginID
	}
	for _, id := range c.AIProviders {
		r.extensions.AIProviders[id] = pluginID
	}
	for _, id := range c.DiscoveryStages {
		r.extensions.DiscoveryStages[id] = pluginID
	}
	for _, id := range c.Commands {
		r.extensions.Commands[id] = pluginID
	}
	for _, id := range c.HookTypes {
		r.extensions.Hooks[id] = pluginID
	}
	for id, engine := range c.CustomEngines {
		r.extensions.CustomEngines[id] = engine
	}
}

func (r *Registry) plugin(id string) (Plugin, error) {
	r.mu.RLock()
	p, ok := r.plugins[id]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", id)
	}
	return p, nil
}

func (r *Registry) setStatus(id string, status Status) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status[id] = status
}

func (r *Registry) setFailed(id string) { r.setStatus(id, StatusFailed) }

func (r *Registry) publish(t eventbus.EventType, msg string, payload interface{}) {
	if r.events != nil {
		r.events.Publish(eventbus.Event{Type: t, Message: msg, Payload: payload})
	}
}

func firstStatus(meta PluginMetadata) Status {
	if meta.Status != "" {
		return meta.Status
	}
	if meta.Enabled {
		return StatusEnabled
	}
	return StatusInstalled
}

func safePluginMeta(meta PluginMetadata, status Status) PluginMetadata {
	meta.Status = status
	meta.Enabled = status == StatusEnabled
	return meta
}

type Loader struct {
	registry *Registry
	fs       filesystem.FileSystem
}

func NewLoader(fs filesystem.FileSystem) *Loader {
	return &Loader{fs: fs, registry: NewRegistry()}
}

func (l *Loader) Register(p Plugin)      { _ = l.registry.Register(p) }
func (l *Loader) List() []PluginMetadata { return l.registry.List() }

func (l *Loader) LoadManifest(path string) (PluginMetadata, error) {
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return PluginMetadata{}, err
	}
	var meta PluginMetadata
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		err = yaml.Unmarshal(data, &meta)
	} else {
		err = json.Unmarshal(data, &meta)
	}
	if err != nil {
		return PluginMetadata{}, err
	}
	if meta.ID == "" || meta.Version == "" {
		return PluginMetadata{}, fmt.Errorf("invalid plugin manifest: id and version are required")
	}
	return meta, nil
}

func (l *Loader) LoadManifestsFrom(dirs ...string) ([]PluginMetadata, error) {
	var manifests []PluginMetadata
	for _, dir := range dirs {
		entries, err := l.fs.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			base := filepath.Join(dir, entry.Name())
			for _, name := range []string{"plugin.yaml", "plugin.yml", "plugin.json"} {
				path := filepath.Join(base, name)
				if !l.fs.Exists(path) {
					continue
				}
				meta, err := l.LoadManifest(path)
				if err != nil {
					return nil, err
				}
				manifests = append(manifests, meta)
				break
			}
		}
	}
	sort.Slice(manifests, func(i, j int) bool { return manifests[i].ID < manifests[j].ID })
	return manifests, nil
}
