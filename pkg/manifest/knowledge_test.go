package manifest

import (
	"testing"

	promptengineassets "github.com/LordCodex/promptengine"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestBundledKnowledgeManifestPreservesTaskMappings(t *testing.T) {
	libraryFS := filesystem.NewEmbeddedFileSystem(promptengineassets.StandardsFS)
	loaded, err := NewLoader(libraryFS).Load(DefaultFilename)
	if err != nil {
		t.Fatalf("load bundled manifest: %v", err)
	}

	engine := NewEngineWithFS(libraryFS)
	if err := engine.Register("promptengine-core", SourceCore, loaded); err != nil {
		t.Fatalf("register bundled manifest: %v", err)
	}
	query := NewQueryEngine(engine)

	required := query.PlaybooksByTask("new_feature")
	assertPlaybookLocation(t, required, "core/05-universal-coding-standards.md")
	assertPlaybookLocation(t, required, "workflows/01-feature-implementation.md")
	assertPlaybookLocation(t, required, "checklists/01-feature-implementation-checklist.md")

	optional := query.OptionalPlaybooksByTask("new_feature")
	assertPlaybookLocation(t, optional, "core/02-architecture-and-simplicity.md")

	prompts := query.PlaybooksByCategory(CategoryPrompt)
	assertPlaybookLocation(t, prompts, "prompts/04-add-feature.md")
}

func TestBundledKnowledgeManifestExposesTechnologyStandards(t *testing.T) {
	libraryFS := filesystem.NewEmbeddedFileSystem(promptengineassets.StandardsFS)
	loaded, err := NewLoader(libraryFS).Load(DefaultFilename)
	if err != nil {
		t.Fatalf("load bundled manifest: %v", err)
	}
	engine := NewEngineWithFS(libraryFS)
	if err := engine.Register("promptengine-core", SourceCore, loaded); err != nil {
		t.Fatalf("register bundled manifest: %v", err)
	}

	standards := NewQueryEngine(engine).StandardsByTechnology("Laravel")
	assertPlaybookLocation(t, standards, "stacks/php-laravel/laravel-engineering-standard.md")
	assertPlaybookLocation(t, standards, "stacks/php-laravel/laravel-routing.md")
}

func assertPlaybookLocation(t *testing.T, playbooks []PlaybookDefinition, location string) {
	t.Helper()
	for _, playbook := range playbooks {
		if playbook.Location == location {
			return
		}
	}
	t.Fatalf("expected playbook %s; got %#v", location, playbooks)
}
