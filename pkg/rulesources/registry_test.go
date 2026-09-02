package rulesources

import "testing"

func TestResolveProfileLaravelInertiaVue(t *testing.T) {
	registry, err := LoadRegistry([]byte(`
version: 1
sources:
  universal:
    repository: example/universal
    ref: abc
    inherits: []
  php:
    repository: example/php
    ref: def
    inherits: [universal]
  laravel:
    repository: example/laravel
    ref: ghi
    inherits: [php]
  vue:
    repository: example/vue
    ref: jkl
    inherits: [universal]
promptengine:
  role: orchestrator
preservation_policy:
  classifications: [KEEP, ADD, MERGE, UPDATE, REFERENCE, CONFLICT]
  invariant: preserve before replacing
`))
	if err != nil {
		t.Fatalf("LoadRegistry() error = %v", err)
	}

	profile, err := LoadProfile([]byte(`
id: laravel-inertia-vue
version: 1
match:
  required_technologies: [php, laravel, inertia, vue]
inheritance: [universal, php, laravel, vue]
resolution_policy:
  concatenate_entire_repositories: false
`))
	if err != nil {
		t.Fatalf("LoadProfile() error = %v", err)
	}

	got, err := registry.ResolveProfile(profile)
	if err != nil {
		t.Fatalf("ResolveProfile() error = %v", err)
	}

	want := []string{"universal", "php", "laravel", "vue"}
	if len(got) != len(want) {
		t.Fatalf("ResolveProfile() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ResolveProfile() = %v, want %v", got, want)
		}
	}

	if !profile.Matches([]string{"Laravel", "PHP", "Vue", "Inertia"}) {
		t.Fatal("profile should match detected Laravel + Inertia + Vue technologies")
	}
	if profile.Matches([]string{"Laravel", "PHP", "React", "Inertia"}) {
		t.Fatal("profile should not match Laravel + Inertia + React")
	}
}

func TestRegistryRejectsInheritanceCycle(t *testing.T) {
	_, err := LoadRegistry([]byte(`
version: 1
sources:
  a:
    repository: example/a
    ref: a1
    inherits: [b]
  b:
    repository: example/b
    ref: b1
    inherits: [a]
`))
	if err == nil {
		t.Fatal("LoadRegistry() expected inheritance cycle error")
	}
}

func TestRegistryRejectsUnpinnedSource(t *testing.T) {
	_, err := LoadRegistry([]byte(`
version: 1
sources:
  universal:
    repository: example/universal
    ref: ""
    inherits: []
`))
	if err == nil {
		t.Fatal("LoadRegistry() expected missing ref error")
	}
}
