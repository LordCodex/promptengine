package discovery

import (
	"context"
	"testing"

	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

func TestDiscoveryPipeline_EmptyProject(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.MkdirAll(".", 0755)

	pm := runDiscovery(t, fs)
	if !pm.HasClassification(ClassGreenfield) {
		t.Fatalf("expected greenfield project, got %#v", pm.Classifications)
	}
}

func TestDiscoveryPipeline_LaravelVue(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("composer.json", []byte(`{"require": {"php": "^8.0", "laravel/framework": "^10.0"}}`), 0644)
	fs.WriteFile("package.json", []byte(`{"dependencies": {"vue": "^3.0"}}`), 0644)
	fs.WriteFile("artisan", []byte("cli entry"), 0644)
	fs.WriteFile(".env", []byte("DB_CONNECTION=mysql"), 0644)
	fs.WriteFile("app/Http/Controllers/Controller.php", []byte("controller"), 0644)
	fs.WriteFile("app/Models/User.php", []byte("model"), 0644)
	fs.WriteFile("resources/views/welcome.blade.php", []byte("view"), 0644)

	pm := runDiscovery(t, fs)
	assertContains(t, pm.Languages, "PHP")
	assertContains(t, pm.Languages, "JavaScript")
	assertContains(t, pm.Frameworks, "Laravel")
	assertContains(t, pm.Frameworks, "Vue")
	assertContains(t, pm.Databases, "MySQL")
	if !pm.Architecture.Backend || !pm.Architecture.Frontend {
		t.Fatalf("expected backend and frontend architecture flags, got %#v", pm.Architecture)
	}
}

func TestDiscoveryPipeline_GoProject(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("go.mod", []byte("module myapp"), 0644)
	fs.WriteFile("internal/domain/user.go", []byte("domain code"), 0644)
	fs.WriteFile("internal/usecase/auth.go", []byte("usecase code"), 0644)
	fs.WriteFile("internal/filesystem/storage.go", []byte("adapter code"), 0644)

	pm := runDiscovery(t, fs)
	assertContains(t, pm.Languages, "Go")
	assertContains(t, pm.PackageManagers, "go mod")
	assertContains(t, pm.Technology.Runtimes, "Go")
	if contains(pm.Technology.Runtimes, "") {
		t.Fatalf("runtime list should not contain blank entries: %#v", pm.Technology.Runtimes)
	}
	assertArchitecture(t, pm, "Clean Architecture")
}

func TestDiscoveryPipeline_FlutterProject(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("pubspec.yaml", []byte("dependencies:\n  flutter:\n    sdk: flutter\n"), 0644)
	fs.WriteFile("lib/main.dart", []byte("void main() {}"), 0644)

	pm := runDiscovery(t, fs)
	assertContains(t, pm.Languages, "Dart")
	assertContains(t, pm.Frameworks, "Flutter")
	if !pm.Architecture.Mobile {
		t.Fatal("expected mobile architecture flag")
	}
}

func TestDiscoveryPipeline_Monorepo(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("apps/api/go.mod", []byte("module api"), 0644)
	fs.WriteFile("apps/web/package.json", []byte(`{"dependencies":{"react":"latest","next":"latest"}}`), 0644)

	pm := runDiscovery(t, fs)
	if !pm.Repository.IsMonorepo || !pm.HasClassification(ClassMonorepo) {
		t.Fatalf("expected monorepo, got repo=%#v classes=%#v", pm.Repository, pm.Classifications)
	}
}

func TestDiscoveryPipeline_MissingConfiguration(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("README.md", []byte("hello"), 0644)

	pm := runDiscovery(t, fs)
	if len(pm.Repository.DocumentationFiles) != 1 {
		t.Fatalf("expected documentation file, got %#v", pm.Repository.DocumentationFiles)
	}
	if len(pm.Languages) != 0 {
		t.Fatalf("expected no languages without config, got %#v", pm.Languages)
	}
}

func TestRepositoryScanner_InvalidProject(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	scanner := NewRepositoryScanner()
	info, err := scanner.Scan(context.Background(), fs, "missing")
	if err != nil {
		t.Fatalf("expected missing project to be handled, got %v", err)
	}
	if info.RootPath != "missing" {
		t.Fatalf("expected root path to be preserved")
	}
}

func TestRepositoryScanner_IgnoresToolCacheDirectories(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("go.mod", []byte("module myapp"), 0644)
	fs.WriteFile(".gocache/00/cache-entry", []byte("cache"), 0644)

	scanner := NewRepositoryScanner()
	info, err := scanner.Scan(context.Background(), fs, ".")
	if err != nil {
		t.Fatalf("expected scan to run, got %v", err)
	}
	if contains(info.Directories, ".gocache") {
		t.Fatalf("expected .gocache to be ignored, got directories %#v", info.Directories)
	}
	if contains(info.Files, ".gocache/00/cache-entry") {
		t.Fatalf("expected .gocache contents to be ignored, got files %#v", info.Files)
	}
	assertContains(t, info.IgnoredFiles, ".gocache")
}

func TestDiscoveryPipeline_UsesManifestRules(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.MkdirAll(".", 0755)
	engine := manifest.NewEngineWithFS(fs)
	engine.Register("project", manifest.SourceProject, &manifest.Manifest{
		Metadata: manifest.ProjectMetadata{Name: "M", Version: "1", SchemaVersion: manifest.SupportedSchemaVersion},
		Technologies: []manifest.TechnologyDefinition{
			{ID: "custom", Language: "Elixir", Framework: "Phoenix", RelatedPlaybooks: []string{}},
		},
	})

	pipeline := NewDefaultPipeline(nil, engine)
	pm, err := pipeline.Execute(context.Background(), fs, ".")
	if err != nil {
		t.Fatalf("expected discovery to run, got %v", err)
	}
	// Manifest rules are registered as knowledge hints, but do not falsely detect without project evidence.
	if contains(pm.Frameworks, "Phoenix") {
		t.Fatal("manifest technology should not be detected without repository evidence")
	}
}

func TestDiscoveryPipeline_PublishesEvents(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("go.mod", []byte("module myapp"), 0644)
	events := eventbus.NewEventBus()
	var seen []eventbus.EventType
	for _, eventType := range []eventbus.EventType{
		eventbus.ProjectDiscoveryStarted,
		eventbus.TechnologyDetected,
		eventbus.ProjectDetected,
		eventbus.ProjectDiscoveryCompleted,
	} {
		tp := eventType
		events.Subscribe(tp, func(e eventbus.Event) {
			seen = append(seen, e.Type)
		})
	}

	pipeline := NewDefaultPipeline(events, nil)
	if _, err := pipeline.Execute(context.Background(), fs, "."); err != nil {
		t.Fatalf("expected discovery to run, got %v", err)
	}
	for _, eventType := range []eventbus.EventType{eventbus.ProjectDiscoveryStarted, eventbus.TechnologyDetected, eventbus.ProjectDetected, eventbus.ProjectDiscoveryCompleted} {
		if !containsEvent(seen, eventType) {
			t.Fatalf("expected event %s, got %#v", eventType, seen)
		}
	}
}

func runDiscovery(t *testing.T, fs filesystem.FileSystem) *ProjectModel {
	t.Helper()
	pipeline := NewDefaultPipeline(nil, nil)
	pm, err := pipeline.Execute(context.Background(), fs, ".")
	if err != nil {
		t.Fatalf("expected discovery to run, got %v", err)
	}
	return pm
}

func assertContains(t *testing.T, items []string, expected string) {
	t.Helper()
	if !contains(items, expected) {
		t.Fatalf("expected %q in %#v", expected, items)
	}
}

func contains(items []string, expected string) bool {
	for _, item := range items {
		if item == expected {
			return true
		}
	}
	return false
}

func containsEvent(items []eventbus.EventType, expected eventbus.EventType) bool {
	for _, item := range items {
		if item == expected {
			return true
		}
	}
	return false
}

func assertArchitecture(t *testing.T, pm *ProjectModel, expected string) {
	t.Helper()
	for _, arch := range pm.Architectures {
		if arch.Style == expected {
			return
		}
	}
	t.Fatalf("expected architecture %q, got %#v", expected, pm.Architectures)
}
