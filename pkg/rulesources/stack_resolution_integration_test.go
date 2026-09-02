package rulesources

import (
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestLaravelInertiaVueResolutionExcludesUnrelatedFrameworkSources(t *testing.T) {
	registry := &Registry{
		Version: 1,
		Sources: map[string]Source{
			"universal": {Repository: "LordCodex/engineering-ai-rules", Ref: "u1", Owns: []string{"universal"}},
			"php":       {Repository: "LordCodex/php-broilerplate", Ref: "p1", Owns: []string{"php"}, Inherits: []string{"universal"}},
			"laravel":   {Repository: "LordCodex/laravel-ai-rules", Ref: "l1", Owns: []string{"laravel", "inertia-laravel"}, Inherits: []string{"php"}},
			"vue":       {Repository: "LordCodex/vue-ai-rules", Ref: "v1", Owns: []string{"vue"}, Inherits: []string{"universal"}},
			"react":     {Repository: "LordCodex/react-ai-rules", Ref: "r1", Owns: []string{"react"}, Inherits: []string{"universal"}},
		},
	}
	if err := registry.Validate(); err != nil {
		t.Fatal(err)
	}

	profiles := []Profile{
		{
			ID:          "laravel-inertia-vue", Version: 1,
			Match:       ProfileMatch{RequiredTechnologies: []string{"php", "laravel", "inertia", "vue"}},
			Inheritance: []string{"universal", "php", "laravel", "vue"},
			RequiredRuleEntrypoints: map[string][]string{
				"universal": {"RULES.md"},
				"php":       {"AGENTS.md"},
				"laravel":   {"RULES.md", "docs/INERTIA.md"},
				"vue":       {"RULES.md"},
			},
		},
		{ID: "laravel", Version: 1, Match: ProfileMatch{RequiredTechnologies: []string{"php", "laravel"}}, Inheritance: []string{"universal", "php", "laravel"}},
		{ID: "vue", Version: 1, Match: ProfileMatch{RequiredTechnologies: []string{"vue"}}, Inheritance: []string{"universal", "vue"}},
		{ID: "react", Version: 1, Match: ProfileMatch{RequiredTechnologies: []string{"react"}}, Inheritance: []string{"universal", "react"}},
	}

	fsys := filesystem.NewMockFileSystem()
	writeSnapshotFixture(t, fsys, "universal", registry.Sources["universal"], []string{"RULES.md", "docs/SECURITY.md"})
	writeSnapshotFixture(t, fsys, "php", registry.Sources["php"], []string{"AGENTS.md", "docs/PHP-LANGUAGE.md"})
	writeSnapshotFixture(t, fsys, "laravel", registry.Sources["laravel"], []string{"RULES.md", "docs/INERTIA.md", "docs/LIVEWIRE.md"})
	writeSnapshotFixture(t, fsys, "vue", registry.Sources["vue"], []string{"RULES.md", "docs/COMPONENTS.md", "docs/I18N.md"})
	writeSnapshotFixture(t, fsys, "react", registry.Sources["react"], []string{"RULES.md", "docs/COMPONENTS.md"})

	resolver := NewResolver(registry, profiles, fsys)
	overlay, resolution, err := resolver.BuildManifest([]string{"PHP", "Laravel", "Inertia", "Vue"})
	if err != nil {
		t.Fatal(err)
	}
	if overlay == nil {
		t.Fatalf("expected authoritative overlay, resolution=%#v", resolution)
	}
	if resolution.ProfileID != "laravel-inertia-vue" {
		t.Fatalf("expected most-specific profile, got %q", resolution.ProfileID)
	}
	if strings.Join(resolution.SourceIDs, ",") != "universal,php,laravel,vue" {
		t.Fatalf("unexpected source inheritance: %#v", resolution.SourceIDs)
	}
	for _, playbook := range overlay.Playbooks {
		if strings.Contains(playbook.Location, "/react/") {
			t.Fatalf("unrelated React source leaked into Vue stack: %s", playbook.Location)
		}
	}
}

func TestResolveLocalizationIntentSelectsI18NWithoutLoadingUnrelatedOptionalRules(t *testing.T) {
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
		ID:          "vue", Version: 1,
		Match:       ProfileMatch{RequiredTechnologies: []string{"vue"}},
		Inheritance: []string{"universal", "vue"},
		RequiredRuleEntrypoints: map[string][]string{
			"universal": {"RULES.md"},
			"vue":       {"RULES.md"},
		},
	}}
	fsys := filesystem.NewMockFileSystem()
	writeSnapshotFixture(t, fsys, "universal", registry.Sources["universal"], []string{"RULES.md"})
	writeSnapshotFixture(t, fsys, "vue", registry.Sources["vue"], []string{"RULES.md", "docs/I18N.md", "docs/STATE-MANAGEMENT.md"})

	resolution, err := NewResolver(registry, profiles, fsys).Resolve([]string{"Vue"}, "add localization and timezone support")
	if err != nil {
		t.Fatal(err)
	}
	foundI18N := false
	for _, file := range resolution.Files {
		if file.RulePath == "docs/I18N.md" {
			foundI18N = true
		}
		if file.RulePath == "docs/STATE-MANAGEMENT.md" {
			t.Fatalf("unrelated optional rule selected for localization intent: %#v", file)
		}
	}
	if !foundI18N {
		t.Fatalf("expected localization intent to select docs/I18N.md, got %#v", resolution.Files)
	}
}

func TestAuthoritativeI18NIsDiscoverableAsGuide(t *testing.T) {
	source := Source{Repository: "LordCodex/vue-ai-rules", Ref: "v1"}
	if got := authoritativeCategory("vue", "docs/I18N.md"); got != "guide" {
		t.Fatalf("expected i18n rule to be categorized as guide, got %q", got)
	}
	description := authoritativeDescription(source, "docs/I18N.md")
	for _, term := range []string{"internationalization", "localization", "timezone", "currency"} {
		if !strings.Contains(strings.ToLower(description), term) {
			t.Fatalf("expected i18n description to contain %q: %s", term, description)
		}
	}
}
