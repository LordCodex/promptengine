package plugins

import (
	"context"
	"fmt"
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/quality"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

type mockPlugin struct {
	meta       PluginMetadata
	contrib    ContributionPoints
	enabled    bool
	installed  bool
	failHealth bool
}

func (m *mockPlugin) Metadata() PluginMetadata          { return m.meta }
func (m *mockPlugin) Contributions() ContributionPoints { return m.contrib }
func (m *mockPlugin) Install(ctx context.Context, fs filesystem.FileSystem) error {
	m.installed = true
	return nil
}
func (m *mockPlugin) Load(ctx context.Context, fs filesystem.FileSystem) error { return nil }
func (m *mockPlugin) Enable(ctx context.Context) error                         { m.enabled = true; return nil }
func (m *mockPlugin) Disable(ctx context.Context) error                        { m.enabled = false; return nil }
func (m *mockPlugin) Upgrade(ctx context.Context, fromVersion string) error {
	m.meta.Version = "2.0.0"
	return nil
}
func (m *mockPlugin) Validate(ctx context.Context, fs filesystem.FileSystem) error {
	if m.failHealth {
		return fmt.Errorf("invalid configuration")
	}
	return nil
}
func (m *mockPlugin) HealthCheck(ctx context.Context, fs filesystem.FileSystem) []HealthFinding {
	if m.failHealth {
		return []HealthFinding{{PluginID: m.meta.ID, Severity: "error", Message: "health failed"}}
	}
	return nil
}
func (m *mockPlugin) Unload(ctx context.Context) error                           { return nil }
func (m *mockPlugin) Remove(ctx context.Context, fs filesystem.FileSystem) error { return nil }

func TestPluginRegistry_RegisterLifecycleAndEvents(t *testing.T) {
	events := eventbus.NewEventBus()
	var seen []eventbus.EventType
	for _, tp := range []eventbus.EventType{eventbus.PluginInstalled, eventbus.PluginEnabled, eventbus.PluginDisabled, eventbus.PluginUpdated, eventbus.PluginRemoved} {
		eventType := tp
		events.Subscribe(eventType, func(e eventbus.Event) { seen = append(seen, e.Type) })
	}
	reg := NewRegistryWithEvents(events)
	fs := filesystem.NewMockFileSystem()
	p := &mockPlugin{meta: PluginMetadata{ID: "company-standard", Name: "Company Standard", Version: "1.0.0", Permissions: []Permission{PermissionReadProject}}}
	if err := reg.Register(p); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if err := reg.Install(context.Background(), p.meta.ID, fs); err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if err := reg.Enable(p.meta.ID); err != nil {
		t.Fatalf("enable failed: %v", err)
	}
	if reg.Status(p.meta.ID) != StatusEnabled || !p.enabled {
		t.Fatalf("expected enabled")
	}
	if err := reg.Disable(p.meta.ID); err != nil {
		t.Fatalf("disable failed: %v", err)
	}
	if err := reg.Upgrade(context.Background(), p.meta.ID, "1.0.0"); err != nil {
		t.Fatalf("upgrade failed: %v", err)
	}
	if err := reg.Remove(context.Background(), p.meta.ID, fs); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	for _, tp := range []eventbus.EventType{eventbus.PluginInstalled, eventbus.PluginEnabled, eventbus.PluginDisabled, eventbus.PluginUpdated, eventbus.PluginRemoved} {
		if !hasEvent(seen, tp) {
			t.Fatalf("expected event %s, got %v", tp, seen)
		}
	}
}

func TestPluginRegistry_RegistrationAndContributions(t *testing.T) {
	reg := NewRegistry()
	rule := quality.RuleFunc{RuleID: "company-rule", RuleEngine: "validation", RuleCategory: "security", Fn: func(ctx context.Context, fs filesystem.FileSystem) ([]quality.Finding, error) { return nil, nil }}
	p := &mockPlugin{meta: PluginMetadata{ID: "ext", Version: "1.0.0"}, contrib: ContributionPoints{
		Workflows: []string{"wf"}, Templates: []string{"tmpl"}, AIProviders: []string{"ai"}, DiscoveryStages: []string{"stage"}, Commands: []string{"cmd"}, HookTypes: []string{"hook"}, QualityRules: []quality.Rule{rule}, CustomEngines: map[string]interface{}{"engine": "custom"},
	}}
	if err := reg.Register(p); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	ext := reg.Extensions()
	if ext.Workflows["wf"] == nil || ext.QualityRules["company-rule"] == nil || ext.CustomEngines["engine"] == nil {
		t.Fatalf("expected contributions registered, got %#v", ext)
	}
	if err := reg.Register(p); err == nil {
		t.Fatal("expected duplicate register error")
	}
}

func TestPluginRegistry_EnableWithDependency(t *testing.T) {
	reg := NewRegistry()
	dep := &mockPlugin{meta: PluginMetadata{ID: "base-pack", Version: "1.0.0"}}
	child := &mockPlugin{meta: PluginMetadata{ID: "laravel-pack", Version: "1.0.0", Dependencies: []string{"base-pack"}}}
	_ = reg.Register(dep)
	_ = reg.Register(child)
	if err := reg.Enable("laravel-pack"); err != nil {
		t.Fatalf("expected enable with dependency, got %v", err)
	}
	if !dep.enabled || !child.enabled {
		t.Fatal("expected dependency and child enabled")
	}
}

func TestPluginRegistry_EnableMissingDependency(t *testing.T) {
	reg := NewRegistry()
	child := &mockPlugin{meta: PluginMetadata{ID: "laravel-pack", Version: "1.0.0", Dependencies: []string{"missing-dep"}}}
	_ = reg.Register(child)
	if err := reg.Enable("laravel-pack"); err == nil {
		t.Fatal("expected missing dependency error")
	}
}

func TestPluginRegistry_TopologicalOrder(t *testing.T) {
	reg := NewRegistry()
	_ = reg.Register(&mockPlugin{meta: PluginMetadata{ID: "top", Version: "1.0.0", Dependencies: []string{"mid"}}})
	_ = reg.Register(&mockPlugin{meta: PluginMetadata{ID: "mid", Version: "1.0.0", Dependencies: []string{"base"}}})
	_ = reg.Register(&mockPlugin{meta: PluginMetadata{ID: "base", Version: "1.0.0"}})
	order, err := reg.ResolveLoadOrder()
	if err != nil {
		t.Fatalf("order failed: %v", err)
	}
	if idx(order, "base") >= idx(order, "mid") || idx(order, "mid") >= idx(order, "top") {
		t.Fatalf("bad order: %v", order)
	}
}

func TestCompatibilityPermissionsAndConfiguration(t *testing.T) {
	if r := CheckCompatibility(PluginMetadata{ID: "test", CompatibilityVersion: "0.5.0"}, "0.4.0"); r.Compatible {
		t.Fatal("expected incompatible")
	}
	reg := NewRegistry()
	if err := reg.Register(&mockPlugin{meta: PluginMetadata{ID: "bad", Version: "1.0.0", Permissions: []Permission{PermissionNetwork}}}); err == nil {
		t.Fatal("expected disallowed permission")
	}
	reg.Configure("company-standard", map[string]interface{}{"security_level": "strict", "enable_checks": true})
	cfg := reg.Config("company-standard")
	if cfg["security_level"] != "strict" || cfg["enable_checks"] != true {
		t.Fatalf("bad config: %#v", cfg)
	}
}

func TestPluginHealthAndManifestExtension(t *testing.T) {
	events := eventbus.NewEventBus()
	healthEvent := false
	events.Subscribe(eventbus.PluginHealthFailed, func(e eventbus.Event) { healthEvent = true })
	reg := NewRegistryWithEvents(events)
	p := &mockPlugin{meta: PluginMetadata{ID: "manifest-ext", Version: "1.0.0"}, failHealth: true, contrib: ContributionPoints{Manifest: &manifest.Manifest{
		Playbooks: []manifest.PlaybookDefinition{{ID: "std", Name: "Std", Category: manifest.CategoryCore, Location: "std.md"}},
	}}}
	_ = reg.Register(p)
	health := reg.Health(context.Background(), filesystem.NewMockFileSystem())
	if len(health) == 0 || !healthEvent {
		t.Fatal("expected health findings and event")
	}
	extended := reg.ExtendManifest(&manifest.Manifest{})
	if len(extended.Playbooks) != 1 || extended.Playbooks[0].ID != "std" {
		t.Fatalf("expected manifest contribution, got %#v", extended.Playbooks)
	}
}

func TestLoader_LoadManifest(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("plugins/company/plugin.yaml", []byte("id: company\nname: Company\nversion: 1.0.0\n"), 0644)
	meta, err := NewLoader(fs).LoadManifest("plugins/company/plugin.yaml")
	if err != nil {
		t.Fatalf("expected load manifest, got %v", err)
	}
	if meta.ID != "company" {
		t.Fatalf("expected company manifest, got %#v", meta)
	}
}

func TestLoader_LoadManifestsFrom(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("local/company/plugin.yaml", []byte("id: company\nversion: 1.0.0\n"), 0644)
	fs.WriteFile("project/stack/plugin.json", []byte(`{"id":"stack","version":"1.0.0"}`), 0644)
	metas, err := NewLoader(fs).LoadManifestsFrom("local", "project", "org")
	if err != nil {
		t.Fatalf("expected load manifests, got %v", err)
	}
	if len(metas) != 2 || metas[0].ID != "company" || metas[1].ID != "stack" {
		t.Fatalf("unexpected manifests: %#v", metas)
	}
}

func idx(items []string, item string) int {
	for i, v := range items {
		if v == item {
			return i
		}
	}
	return -1
}

func hasEvent(events []eventbus.EventType, expected eventbus.EventType) bool {
	for _, eventType := range events {
		if eventType == expected {
			return true
		}
	}
	return false
}
