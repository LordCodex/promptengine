package context

import (
	"context"
	"testing"

	promptengineassets "github.com/LordCodex/promptengine"
	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

func TestContextBuildCombinesProjectFilesWithRelevantBundledStandards(t *testing.T) {
	projectFS := filesystem.NewMockFileSystem()
	if err := projectFS.WriteFile("AGENTS.md", []byte("project-specific rules"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := projectFS.WriteFile("app/Services/PaymentService.php", []byte("payment service"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := projectFS.WriteFile("routes/api.php", []byte("api routes"), 0644); err != nil {
		t.Fatal(err)
	}

	libraryFS := filesystem.NewEmbeddedFileSystem(promptengineassets.StandardsFS)
	resourceFS := filesystem.NewOverlayFileSystem(projectFS, libraryFS)
	manifestEngine := manifest.NewEngineWithFS(resourceFS)
	coreManifest, err := manifest.NewLoader(libraryFS).Load(manifest.DefaultFilename)
	if err != nil {
		t.Fatalf("load bundled manifest: %v", err)
	}
	if err := manifestEngine.Register("promptengine-core", manifest.SourceCore, coreManifest); err != nil {
		t.Fatalf("register bundled manifest: %v", err)
	}

	project := discovery.NewProjectModel(".")
	project.Frameworks = []string{"Laravel"}
	project.Languages = []string{"PHP"}
	project.Repository.Files = []string{"AGENTS.md", "app/Services/PaymentService.php", "routes/api.php"}
	project.SyncLegacyFields()

	engine := NewEngine(resourceFS, WithManifestQuery(manifest.NewQueryEngine(manifestEngine)))
	pkg, err := engine.Build(context.Background(), ContextRequest{
		TaskType:   TaskType("new_feature"),
		Project:    project,
		UserIntent: "Add an API endpoint for author withdrawal payments",
		MaxBytes:   1_000_000,
	})
	if err != nil {
		t.Fatalf("build context: %v", err)
	}

	assertContextPath(t, pkg.SelectedDocs, "AGENTS.md")
	assertContextPath(t, pkg.RelevantStandards, "core/05-universal-coding-standards.md")
	assertContextPath(t, pkg.RelevantStandards, "workflows/01-feature-implementation.md")
	assertContextPath(t, pkg.RelevantStandards, "checklists/01-feature-implementation-checklist.md")
	assertContextPath(t, pkg.RelevantStandards, "stacks/php-laravel/laravel-engineering-standard.md")
	assertContextPath(t, pkg.SelectedFiles, "app/Services/PaymentService.php")

	assertContextPathMissing(t, pkg.RelevantStandards, "stacks/dart-flutter/flutter-dart-engineering-standard.md")
	assertContextPathMissing(t, pkg.RelevantStandards, "core/27-seo-engineering-standard.md")
}

func assertContextPath(t *testing.T, paths []string, expected string) {
	t.Helper()
	for _, path := range paths {
		if path == expected {
			return
		}
	}
	t.Fatalf("expected %s in %#v", expected, paths)
}

func assertContextPathMissing(t *testing.T, paths []string, unexpected string) {
	t.Helper()
	for _, path := range paths {
		if path == unexpected {
			t.Fatalf("did not expect %s in %#v", unexpected, paths)
		}
	}
}
