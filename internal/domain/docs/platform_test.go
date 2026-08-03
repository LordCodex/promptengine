package docs

import (
	"context"
	"strings"
	"testing"

	"github.com/LordCodex/promptengine/internal/domain/discovery"
	docsync "github.com/LordCodex/promptengine/internal/domain/docs/sync"
	"github.com/LordCodex/promptengine/internal/domain/workflows"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
	"github.com/LordCodex/promptengine/pkg/manifest"
)

func TestPlatform_TemplateLoadingAndGeneration(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("project/templates/Architecture.template.md", []byte("# Architecture\n\n## Overview\n\n{{project_name}} uses {{stack}}.\n"), 0644)
	pm := discovery.NewProjectModel(".")
	pm.Project.Name = "Billing"
	pm.Languages = []string{"Go"}
	pm.Frameworks = []string{"React"}
	pm.SyncLegacyFields()

	platform := NewPlatform(fs, nil, nil)
	doc, err := platform.Generate(context.Background(), GenerationRequest{DocumentID: "architecture", Project: pm, Overwrite: true})
	if err != nil {
		t.Fatalf("expected generation to succeed, got %v", err)
	}
	if doc.Path != "docs/Architecture.md" || !fs.Exists("docs/Architecture.md") {
		t.Fatalf("expected generated architecture doc, got %#v", doc)
	}
	data, _ := fs.ReadFile("docs/Architecture.md")
	if string(data) == "" || !contains(string(data), "Billing") {
		t.Fatalf("expected variables to be replaced, got %s", string(data))
	}
}

func TestPlatform_DiscoversExistingDocumentation(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("README.md", []byte("# Overview\n\nHello"), 0644)
	fs.WriteFile("docs/API.md", []byte("# API\n\nEndpoints"), 0644)
	pm := discovery.NewProjectModel(".")
	pm.Repository.DocumentationFiles = []string{"README.md", "docs/API.md", "docs/Extra.md"}

	docs := NewPlatform(fs, nil, nil).Discover(pm)
	if len(docs) == 0 {
		t.Fatal("expected discovered docs")
	}
	if !hasDocument(docs, "README.md") || !hasDocument(docs, "docs/API.md") {
		t.Fatalf("expected README and API docs, got %#v", docs)
	}
}

func TestPlatform_ValidationFindsMissingAndInvalidDocs(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("docs/Architecture.md", []byte("No headings here"), 0644)
	platform := NewPlatform(fs, nil, nil)
	report, err := platform.Validate(nil)
	if err != nil {
		t.Fatalf("expected validation to run, got %v", err)
	}
	if len(report.Findings) == 0 {
		t.Fatal("expected validation findings")
	}
}

func TestPlatform_SyncDetectsDrift(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	platform := NewPlatform(fs, nil, nil)
	result := platform.Sync([]string{"database/migrations/2026_create_users.php", "routes/api.php"}, true)
	if len(result.Recommendations) == 0 {
		t.Fatal("expected sync recommendations")
	}
	if !hasRecommendation(result.Recommendations, "database") || !hasRecommendation(result.Recommendations, "api") {
		t.Fatalf("expected database and api recommendations, got %#v", result.Recommendations)
	}
}

func TestPlatform_ManifestControlsDocumentationRules(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	platform := NewPlatform(fs, nil, nil)
	platform.ApplyManifestRules(&manifest.Manifest{
		Templates: []manifest.TemplateDefinition{{Name: "Runbook", Location: "project/templates/Runbook.template.md", Version: "1.0.0"}},
	})
	if _, ok := platform.Registry().Get("runbook"); !ok {
		t.Fatal("expected manifest template to register document rule")
	}
}

func TestPlatform_WorkflowIntegration(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("project/templates/Architecture.template.md", []byte("# Architecture\n\n## Overview\n\n{project_name}\n"), 0644)
	platform := NewPlatform(fs, nil, nil)
	reg := workflows.NewRegistry()
	reg.RegisterFromSource("test", workflows.Workflow{
		ID:    "docs-flow",
		Steps: []workflows.WorkflowStep{{ID: "generate", Order: 1, Action: "generate_documentation", InputMapping: map[string]string{"document_id": "architecture"}}},
	})
	engine := workflows.NewEngine(fs, reg, eventbus.NewEventBus())
	for action, handler := range platform.WorkflowHandlers() {
		engine.RegisterHandler(action, handler)
	}
	if _, err := engine.Execute(context.Background(), "docs-flow", workflows.NewFlowContext("docs-flow")); err != nil {
		t.Fatalf("expected workflow docs action to run, got %v", err)
	}
	if !fs.Exists("docs/Architecture.md") {
		t.Fatal("expected workflow handler to generate doc")
	}
}

func TestPlatform_PublishesEvents(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("project/templates/Architecture.template.md", []byte("# Architecture\n\n## Overview\n\n{project_name}\n"), 0644)
	events := eventbus.NewEventBus()
	var generated bool
	var syncRequired bool
	events.Subscribe(eventbus.DocumentationGenerated, func(e eventbus.Event) { generated = true })
	events.Subscribe(eventbus.DocumentationSyncRequired, func(e eventbus.Event) { syncRequired = true })
	platform := NewPlatform(fs, events, nil)
	if _, err := platform.Generate(context.Background(), GenerationRequest{DocumentID: "architecture", Overwrite: true}); err != nil {
		t.Fatalf("expected generation, got %v", err)
	}
	platform.Sync([]string{"routes/api.php"}, true)
	if !generated || !syncRequired {
		t.Fatalf("expected generated and sync events, generated=%v sync=%v", generated, syncRequired)
	}
}

func hasDocument(docs []Document, path string) bool {
	for _, doc := range docs {
		if doc.Path == path {
			return true
		}
	}
	return false
}

func hasRecommendation(recs []docsync.SyncRecommendation, docID string) bool {
	for _, rec := range recs {
		if rec.DocID == docID {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
