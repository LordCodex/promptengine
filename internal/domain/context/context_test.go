package context

import (
	stdcontext "context"
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

func TestContextEngine_FeatureSelectionUsesManifestAndDiscovery(t *testing.T) {
	fs, pm := fixtureProject()
	engine := NewEngine(fs, WithManifestQuery(fixtureManifestQuery(fs)))

	pkg, err := engine.Build(stdcontext.Background(), ContextRequest{
		TaskType:           TaskAddFeature,
		WorkflowType:       "add-feature",
		Project:            pm,
		UserIntent:         "Add payment feature",
		RequestedOperation: "feature implementation",
		AffectedFiles:      []string{"app/Services/PaymentService.php"},
		Budget:             BudgetSmall,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertSelected(t, pkg, "AGENTS.md")
	assertSelected(t, pkg, "docs/BusinessRules.md")
	assertSelected(t, pkg, "app/Services/PaymentService.php")
	assertSelected(t, pkg, "standards/feature.md")
	assertNotSelected(t, pkg, "docs/Architecture.md")
	assertNotSelected(t, pkg, "docs/Database.md")
	assertNotSelected(t, pkg, "docs/API.md")
	if len(pkg.Items) == 0 || pkg.Items[0].Path != "app/Services/PaymentService.php" {
		t.Fatalf("expected affected file to rank first, got %#v", pkg.Items)
	}
}

func TestContextEngine_FeatureIncludesDatabaseAndAPIOnlyWhenTaskRequiresThem(t *testing.T) {
	fs, pm := fixtureProject()
	engine := NewEngine(fs)

	pkg, err := engine.Build(stdcontext.Background(), ContextRequest{
		TaskType:           TaskAddFeature,
		Project:            pm,
		UserIntent:         "Add payment webhook endpoint and persist provider transaction data with a migration",
		RequestedOperation: "API and database change",
		Budget:             BudgetSmall,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertSelected(t, pkg, "docs/Database.md")
	assertSelected(t, pkg, "docs/API.md")
	assertNotSelected(t, pkg, "docs/Architecture.md")
}

func TestContextEngine_FeatureDoesNotFloodUnrelatedTechnologyFiles(t *testing.T) {
	fs, pm := fixtureProject()
	fs.WriteFile("app/Services/EmailService.php", []byte("email service"), 0644)
	fs.WriteFile("app/Services/InventoryService.php", []byte("inventory service"), 0644)
	pm.Repository.Files = append(pm.Repository.Files, "app/Services/EmailService.php", "app/Services/InventoryService.php")
	engine := NewEngine(fs)

	pkg, err := engine.Build(stdcontext.Background(), ContextRequest{
		TaskType:      TaskAddFeature,
		Project:       pm,
		UserIntent:    "Improve payment processing",
		AffectedFiles: []string{"app/Services/PaymentService.php"},
		Budget:        BudgetLarge,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertSelected(t, pkg, "app/Services/PaymentService.php")
	assertNotSelected(t, pkg, "app/Services/EmailService.php")
	assertNotSelected(t, pkg, "app/Services/InventoryService.php")
}

func TestContextEngine_BugFixSelection(t *testing.T) {
	fs, pm := fixtureProject()
	fs.WriteFile("tests/Services/PaymentServiceTest.php", []byte("payment tests"), 0644)
	engine := NewEngine(fs)

	pkg, err := engine.Build(stdcontext.Background(), ContextRequest{
		TaskType:      TaskBugFix,
		Project:       pm,
		AffectedFiles: []string{"app/Services/PaymentService.php"},
		Budget:        BudgetSmall,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertSelected(t, pkg, "docs/Troubleshooting.md")
	assertSelected(t, pkg, "app/Services/PaymentService.php")
	assertSelected(t, pkg, "tests/Services/PaymentService.php")
}

func TestContextEngine_RefactorSelection(t *testing.T) {
	fs, pm := fixtureProject()
	engine := NewEngine(fs)

	pkg, err := engine.Build(stdcontext.Background(), ContextRequest{
		TaskType:      TaskRefactor,
		Project:       pm,
		AffectedFiles: []string{"app/Services/PaymentService.php"},
		Budget:        BudgetSmall,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assertSelected(t, pkg, "docs/Architecture.md")
	assertSelected(t, pkg, "app/Services/PaymentService.php")
}

func TestContextEngine_LargeFileSummarizedAndBudgeted(t *testing.T) {
	fs, pm := fixtureProject()
	large := strings.Repeat("large\n", 2000)
	fs.WriteFile("app/Services/HugeService.php", []byte(large), 0644)
	pm.Repository.Files = append(pm.Repository.Files, "app/Services/HugeService.php")
	engine := NewEngine(fs)

	pkg, err := engine.Build(stdcontext.Background(), ContextRequest{
		TaskType:      TaskAddFeature,
		Project:       pm,
		AffectedFiles: []string{"app/Services/HugeService.php"},
		MaxBytes:      1200,
		Budget:        BudgetTiny,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pkg.Summary.FinalSize > 1200 {
		t.Fatalf("expected budget respected, got %d", pkg.Summary.FinalSize)
	}
	foundTruncated := false
	for _, item := range pkg.Items {
		if item.Path == "app/Services/HugeService.php" && item.Truncated {
			foundTruncated = true
		}
	}
	if !foundTruncated {
		t.Fatalf("expected large affected file to be summarized, got %#v", pkg.Items)
	}
}

func TestContextEngine_MissingDocumentation(t *testing.T) {
	fs, pm := fixtureProject()
	delete(fs.Files, "docs/API.md")
	engine := NewEngine(fs)

	pkg, err := engine.Build(stdcontext.Background(), ContextRequest{TaskType: TaskAddFeature, Project: pm, Budget: BudgetSmall})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if packageHas(pkg, "docs/API.md") {
		t.Fatal("missing documentation should not be selected")
	}
}

func TestContextEngine_MonorepoSelection(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("apps/api/go.mod", []byte("module api"), 0644)
	fs.WriteFile("apps/api/internal/payment/service.go", []byte("service"), 0644)
	pm := discovery.NewProjectModel(".")
	pm.Repository.Files = []string{"apps/api/go.mod", "apps/api/internal/payment/service.go"}
	pm.Repository.IsMonorepo = true
	pm.Languages = []string{"Go"}
	pm.SyncLegacyFields()
	engine := NewEngine(fs)

	pkg, err := engine.Build(stdcontext.Background(), ContextRequest{TaskType: TaskBugFix, Project: pm, AffectedFiles: []string{"apps/api/internal/payment/service.go"}, Budget: BudgetSmall})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assertSelected(t, pkg, "apps/api/internal/payment/service.go")
}

func TestContextEngine_CacheInvalidation(t *testing.T) {
	fs, pm := fixtureProject()
	cache := NewCache()
	engine := NewEngine(fs, WithCache(cache))
	req := ContextRequest{TaskType: TaskAddFeature, Project: pm, AffectedFiles: []string{"app/Services/PaymentService.php"}, Budget: BudgetSmall}

	first, err := engine.Build(stdcontext.Background(), req)
	if err != nil {
		t.Fatalf("first build failed: %v", err)
	}
	second, err := engine.Build(stdcontext.Background(), req)
	if err != nil {
		t.Fatalf("second build failed: %v", err)
	}
	if !second.Summary.CacheHit {
		t.Fatal("expected second build to hit cache")
	}
	fs.WriteFile("app/Services/PaymentService.php", []byte("changed"), 0644)
	third, err := engine.Build(stdcontext.Background(), req)
	if err != nil {
		t.Fatalf("third build failed: %v", err)
	}
	if third.Summary.CacheHit {
		t.Fatal("expected cache invalidation after file change")
	}
	if first.Summary.FinalSize == third.Summary.FinalSize {
		t.Fatal("expected changed file to alter context fingerprint/output")
	}
}

func TestContextEngine_PublishesEvents(t *testing.T) {
	fs, pm := fixtureProject()
	events := eventbus.NewEventBus()
	var seen []eventbus.EventType
	for _, eventType := range []eventbus.EventType{eventbus.ContextBuildStarted, eventbus.ContextItemSelected, eventbus.ContextBuilt} {
		tp := eventType
		events.Subscribe(tp, func(e eventbus.Event) {
			seen = append(seen, e.Type)
		})
	}
	engine := NewEngine(fs, WithEventBus(events))
	if _, err := engine.Build(stdcontext.Background(), ContextRequest{TaskType: TaskAddFeature, Project: pm, Budget: BudgetSmall}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	for _, eventType := range []eventbus.EventType{eventbus.ContextBuildStarted, eventbus.ContextItemSelected, eventbus.ContextBuilt} {
		if !eventSeen(seen, eventType) {
			t.Fatalf("expected event %s, got %#v", eventType, seen)
		}
	}
}

func TestContextEngine_RanksTaskIntentKeywords(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("internal/domain/context/engine.go", []byte("context engine"), 0644)
	fs.WriteFile("internal/domain/quality/platform.go", []byte("quality platform"), 0644)
	fs.WriteFile("internal/cache/cache.go", []byte("cache"), 0644)
	pm := discovery.NewProjectModel(".")
	pm.Languages = []string{"Go"}
	pm.Repository.Files = []string{"internal/domain/context/engine.go", "internal/domain/quality/platform.go", "internal/cache/cache.go"}
	pm.SyncLegacyFields()

	pkg, err := NewEngine(fs).Build(stdcontext.Background(), ContextRequest{TaskType: TaskAddFeature, Project: pm, UserIntent: "Dogfood PromptEngine production hardening", Budget: BudgetSmall})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !packageHas(pkg, "internal/domain/context/engine.go") || !packageHas(pkg, "internal/domain/quality/platform.go") {
		t.Fatalf("expected task intent to prioritize context and quality files, got %#v", pkg.SelectedFiles)
	}
}

func TestContextEngine_RedactsSecretsAndSkipsSensitivePaths(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("internal/config/config.go", []byte("package config\nconst token = \"abc\"\nAPI_KEY=sk_live_secret\n"), 0644)
	fs.WriteFile(".env", []byte("PASSWORD=secret"), 0644)
	pm := discovery.NewProjectModel(".")
	pm.Languages = []string{"Go"}
	pm.Repository.Files = []string{"internal/config/config.go", ".env"}
	pm.SyncLegacyFields()

	pkg, err := NewEngine(fs).Build(stdcontext.Background(), ContextRequest{
		TaskType:      TaskBugFix,
		Project:       pm,
		AffectedFiles: []string{"internal/config/config.go", ".env"},
		Budget:        BudgetSmall,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if packageHas(pkg, ".env") {
		t.Fatal("sensitive environment files should not be selected")
	}
	assertSelected(t, pkg, "internal/config/config.go")
	for _, item := range pkg.Items {
		if item.Path == "internal/config/config.go" {
			if strings.Contains(item.Content, "sk_live_secret") {
				t.Fatalf("context leaked API key: %q", item.Content)
			}
			if !strings.Contains(item.Content, "[REDACTED]") {
				t.Fatalf("expected redacted marker, got %q", item.Content)
			}
			return
		}
	}
	t.Fatal("expected config file to be present")
}

func TestContextFormatter_JSONEnvelope(t *testing.T) {
	fs, pm := fixtureProject()
	pkg, err := NewEngine(fs).GenerateContext(stdcontext.Background(), TaskBugFix, pm, BudgetTiny)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	formatted, err := NewFormatter(ProviderCursor).Format(pkg)
	if err != nil {
		t.Fatalf("format failed: %v", err)
	}
	if !strings.Contains(formatted, `"rules"`) {
		t.Fatalf("expected JSON rules envelope, got %s", formatted)
	}
}

func fixtureProject() (*filesystem.MockFileSystem, *discovery.ProjectModel) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("AGENTS.md", []byte("project rules"), 0644)
	fs.WriteFile("docs/Architecture.md", []byte("payment architecture"), 0644)
	fs.WriteFile("docs/Database.md", []byte("payment tables"), 0644)
	fs.WriteFile("docs/API.md", []byte("payment endpoints"), 0644)
	fs.WriteFile("docs/BusinessRules.md", []byte("payment business rules"), 0644)
	fs.WriteFile("docs/Troubleshooting.md", []byte("payment troubleshooting"), 0644)
	fs.WriteFile("app/Services/PaymentService.php", []byte("payment service"), 0644)
	fs.WriteFile("tests/Services/PaymentService.php", []byte("payment tests"), 0644)
	fs.WriteFile("standards/feature.md", []byte("feature standard"), 0644)
	pm := discovery.NewProjectModel(".")
	pm.PromptEngine.AgentsMDPresent = true
	pm.Frameworks = []string{"Laravel"}
	pm.Languages = []string{"PHP"}
	pm.Repository.Files = []string{
		"AGENTS.md", "docs/Architecture.md", "docs/Database.md", "docs/API.md",
		"docs/BusinessRules.md", "docs/Troubleshooting.md", "app/Services/PaymentService.php",
		"tests/Services/PaymentService.php", "standards/feature.md",
	}
	pm.Docs["Architecture"] = discovery.DocSpec{Name: "Architecture", Path: "docs/Architecture.md", Exists: true}
	pm.Docs["Database"] = discovery.DocSpec{Name: "Database", Path: "docs/Database.md", Exists: true}
	pm.Docs["API"] = discovery.DocSpec{Name: "API", Path: "docs/API.md", Exists: true}
	pm.Docs["BusinessRules"] = discovery.DocSpec{Name: "BusinessRules", Path: "docs/BusinessRules.md", Exists: true}
	pm.Docs["Troubleshooting"] = discovery.DocSpec{Name: "Troubleshooting", Path: "docs/Troubleshooting.md", Exists: true}
	pm.SyncLegacyFields()
	return fs, pm
}

func fixtureManifestQuery(fs filesystem.FileSystem) *manifest.QueryEngine {
	_ = fs.WriteFile("prompts/feature.md", []byte("feature prompt"), 0644)
	engine := manifest.NewEngineWithFS(fs)
	if err := engine.Register("project", manifest.SourceProject, &manifest.Manifest{
		Metadata: manifest.ProjectMetadata{Name: "Project", Version: "1", SchemaVersion: manifest.SupportedSchemaVersion},
		Playbooks: []manifest.PlaybookDefinition{
			{ID: "feature-standard", Name: "Feature Standard", Category: manifest.CategoryWorkflows, Location: "standards/feature.md", Priority: 100},
		},
		Workflows: []manifest.WorkflowDefinition{
			{ID: "add-feature", Steps: []string{"build"}, RequiredPlaybooks: []string{"feature-standard"}, Prompts: []string{"add-feature"}},
		},
		Prompts: []manifest.PromptMapping{{TaskType: "add-feature", PromptTemplate: "prompts/feature.md"}},
	}); err != nil {
		panic(err)
	}
	return manifest.NewQueryEngine(engine)
}

func assertSelected(t *testing.T, pkg *ContextPackage, path string) {
	t.Helper()
	if !packageHas(pkg, path) {
		t.Fatalf("expected %s to be selected; items=%#v", path, pkg.Items)
	}
}

func assertNotSelected(t *testing.T, pkg *ContextPackage, path string) {
	t.Helper()
	if packageHas(pkg, path) {
		t.Fatalf("expected %s to be excluded; items=%#v", path, pkg.Items)
	}
}

func packageHas(pkg *ContextPackage, path string) bool {
	for _, item := range pkg.Items {
		if item.Path == path {
			return true
		}
	}
	return false
}

func eventSeen(events []eventbus.EventType, expected eventbus.EventType) bool {
	for _, eventType := range events {
		if eventType == expected {
			return true
		}
	}
	return false
}
