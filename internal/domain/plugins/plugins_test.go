package plugins

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

// --- Mock Plugin ---

type mockPlugin struct {
	meta    PluginMetadata
	contrib ContributionPoints
	enabled bool
}

func (m *mockPlugin) Metadata() PluginMetadata       { return m.meta }
func (m *mockPlugin) Contributions() ContributionPoints { return m.contrib }
func (m *mockPlugin) OnInstall(_ filesystem.FileSystem) error { return nil }
func (m *mockPlugin) OnUninstall(_ filesystem.FileSystem) error { return nil }
func (m *mockPlugin) OnEnable() error  { m.enabled = true; return nil }
func (m *mockPlugin) OnDisable() error { m.enabled = false; return nil }
func (m *mockPlugin) OnUpdate(_ string) error { return nil }

func TestPluginRegistry_RegisterAndList(t *testing.T) {
	reg := NewRegistry()
	p := &mockPlugin{meta: PluginMetadata{ID: "laravel-pack", Name: "Laravel Pack", Version: "1.0.0"}}
	if err := reg.Register(p); err != nil {
		t.Fatalf("expected no error on first register, got %v", err)
	}
	// Duplicate registration must fail
	if err := reg.Register(p); err == nil {
		t.Error("expected error on duplicate registration, got nil")
	}
	if len(reg.List()) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(reg.List()))
	}
}

func TestPluginRegistry_EnableWithDependency(t *testing.T) {
	reg := NewRegistry()
	dep := &mockPlugin{meta: PluginMetadata{ID: "base-pack"}}
	child := &mockPlugin{meta: PluginMetadata{ID: "laravel-pack", Dependencies: []string{"base-pack"}}}

	_ = reg.Register(dep)
	_ = reg.Register(child)

	// Enable child — dependency exists, should succeed
	if err := reg.Enable("laravel-pack"); err != nil {
		t.Errorf("expected enable to succeed when dependency is present, got %v", err)
	}
}

func TestPluginRegistry_EnableMissingDependency(t *testing.T) {
	reg := NewRegistry()
	child := &mockPlugin{meta: PluginMetadata{ID: "laravel-pack", Dependencies: []string{"missing-dep"}}}
	_ = reg.Register(child)

	if err := reg.Enable("laravel-pack"); err == nil {
		t.Error("expected error when dependency is missing, got nil")
	}
}

func TestPluginRegistry_TopologicalOrder(t *testing.T) {
	reg := NewRegistry()
	base := &mockPlugin{meta: PluginMetadata{ID: "base"}}
	mid := &mockPlugin{meta: PluginMetadata{ID: "mid", Dependencies: []string{"base"}}}
	top := &mockPlugin{meta: PluginMetadata{ID: "top", Dependencies: []string{"mid"}}}

	_ = reg.Register(top)
	_ = reg.Register(mid)
	_ = reg.Register(base)

	order, err := reg.ResolveLoadOrder()
	if err != nil {
		t.Fatalf("expected no error resolving load order, got %v", err)
	}
	// base must appear before mid, mid before top
	idx := func(s string) int {
		for i, v := range order { if v == s { return i } }
		return -1
	}
	if idx("base") >= idx("mid") || idx("mid") >= idx("top") {
		t.Errorf("unexpected load order: %v", order)
	}
}

func TestCompatibilityCheck(t *testing.T) {
	meta := PluginMetadata{ID: "test", MinCoreVer: "0.5.0"}
	if r := CheckCompatibility(meta, "0.4.0"); r.Compatible {
		t.Error("expected incompatible for older core version")
	}
	if r := CheckCompatibility(meta, "0.6.0"); !r.Compatible {
		t.Errorf("expected compatible for newer core version, got: %s", r.Reason)
	}
}
