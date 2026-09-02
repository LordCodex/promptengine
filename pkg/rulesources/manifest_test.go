package rulesources

import (
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
	"gopkg.in/yaml.v3"
)

func TestResolverBuildManifestUsesPinnedCachedRules(t *testing.T) {
	registry := &Registry{
		Version: 1,
		Sources: map[string]Source{
			"universal": {Repository: "LordCodex/engineering-ai-rules", Ref: "u1", Owns: []string{"universal"}},
			"vue":       {Repository: "LordCodex/vue-ai-rules", Ref: "v1", Owns: []string{"vue"}, Inherits: []string{"universal"}},
		},
	}
	if err := registry.Validate(); err != nil {
		t.Fatal(err)
	}
	profiles := []Profile{{
		ID:          "vue",
		Version:     1,
		Match:       ProfileMatch{RequiredTechnologies: []string{"vue"}},
		Inheritance: []string{"universal", "vue"},
		RequiredRuleEntrypoints: map[string][]string{
			"universal": {"RULES.md"},
			"vue":       {"RULES.md"},
		},
	}}
	fsys := filesystem.NewMockFileSystem()
	writeSnapshotFixture(t, fsys, "universal", registry.Sources["universal"], []string{"RULES.md", "docs/SECURITY.md"})
	writeSnapshotFixture(t, fsys, "vue", registry.Sources["vue"], []string{"RULES.md", "docs/ARCHITECTURE.md"})

	resolver := NewResolver(registry, profiles, fsys)
	overlay, resolution, err := resolver.BuildManifest([]string{"Vue", "JavaScript"})
	if err != nil {
		t.Fatal(err)
	}
	if overlay == nil || resolution.ProfileID != "vue" {
		t.Fatalf("expected Vue authoritative overlay, overlay=%#v resolution=%#v", overlay, resolution)
	}
	if len(overlay.Technologies) != 1 || overlay.Technologies[0].ID != "vue" {
		t.Fatalf("unexpected technology overlay: %#v", overlay.Technologies)
	}
	if len(overlay.Playbooks) != 4 {
		t.Fatalf("expected all cached rule text files as playbooks, got %d", len(overlay.Playbooks))
	}
	foundRequired := false
	for _, playbook := range overlay.Playbooks {
		if playbook.Location == ".promptengine/rules/vue/v1/RULES.md" {
			foundRequired = true
			if playbook.Priority != 100 {
				t.Fatalf("expected authoritative priority, got %d", playbook.Priority)
			}
		}
	}
	if !foundRequired {
		t.Fatal("expected Vue RULES.md in authoritative overlay")
	}
}

func TestResolverBuildManifestFallsBackWhenSourceMissing(t *testing.T) {
	registry := &Registry{
		Version: 1,
		Sources: map[string]Source{
			"universal": {Repository: "LordCodex/engineering-ai-rules", Ref: "u1"},
			"vue":       {Repository: "LordCodex/vue-ai-rules", Ref: "v1", Inherits: []string{"universal"}},
		},
	}
	if err := registry.Validate(); err != nil {
		t.Fatal(err)
	}
	profile := Profile{
		ID:          "vue",
		Version:     1,
		Match:       ProfileMatch{RequiredTechnologies: []string{"vue"}},
		Inheritance: []string{"vue"},
		RequiredRuleEntrypoints: map[string][]string{
			"universal": {"RULES.md"},
			"vue":       {"RULES.md"},
		},
	}
	fsys := filesystem.NewMockFileSystem()
	writeSnapshotFixture(t, fsys, "universal", registry.Sources["universal"], []string{"RULES.md"})

	resolver := NewResolver(registry, []Profile{profile}, fsys)
	overlay, resolution, err := resolver.BuildManifest([]string{"Vue"})
	if err != nil {
		t.Fatal(err)
	}
	if overlay != nil {
		t.Fatal("expected bundled fallback instead of a partial authoritative overlay")
	}
	if len(resolution.MissingSources) != 1 || resolution.MissingSources[0] != "vue" {
		t.Fatalf("expected missing Vue source, got %#v", resolution.MissingSources)
	}
}

func writeSnapshotFixture(t *testing.T, fsys *filesystem.MockFileSystem, sourceID string, source Source, files []string) {
	t.Helper()
	for _, file := range files {
		fsys.WriteFile(DefaultCacheRoot+"/"+sourceID+"/"+source.Ref+"/"+file, []byte(file), 0644)
	}
	snapshot := Snapshot{SourceID: sourceID, Repository: source.Repository, Ref: source.Ref, SyncedAt: "2026-09-02T07:00:00Z", Files: files}
	data, err := yaml.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	fsys.WriteFile(DefaultCacheRoot+"/"+sourceID+"/"+source.Ref+"/.source.yaml", data, 0644)
}
