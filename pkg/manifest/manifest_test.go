package manifest

import (
	"errors"
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestLoader_LoadsAndValidatesManifest(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	writeManifestFixtures(fs)
	fs.WriteFile("playbook-manifest.json", []byte(validManifestJSON()), 0644)

	loader := NewLoader(fs)
	m, err := loader.Load("playbook-manifest.json")
	if err != nil {
		t.Fatalf("expected load to succeed, got %v", err)
	}
	if err := Validate(m, fs); err != nil {
		t.Fatalf("expected validation to succeed, got %v", err)
	}
	if m.Metadata.Name != "Example" {
		t.Fatalf("expected metadata name, got %q", m.Metadata.Name)
	}
}

func TestLoader_ConvertsLegacyRepositoryManifest(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("core/a.md", []byte("a"), 0644)
	fs.WriteFile("playbook-manifest.json", []byte(`{
		"repository_metadata": {"promptengine_version": "1.0.0", "manifest_version": "1.1.0"},
		"core_playbooks": [{"id": "a", "path": "core/a.md"}],
		"task_mappings": {"feature": {"required_playbook_ids": ["a"]}}
	}`), 0644)

	m, err := NewLoader(fs).Load("playbook-manifest.json")
	if err != nil {
		t.Fatalf("expected legacy load to succeed, got %v", err)
	}
	if m.Metadata.Name != "PromptEngine" || m.Metadata.SchemaVersion != SupportedSchemaVersion {
		t.Fatalf("expected converted metadata, got %#v", m.Metadata)
	}
	if err := Validate(m, fs); err != nil {
		t.Fatalf("expected converted manifest to validate, got %v", err)
	}
}

func TestLoader_Discover(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	fs.WriteFile("project/playbook-manifest.json", []byte(validManifestJSON()), 0644)
	loader := NewLoader(fs)

	path, ok := loader.Discover("project/internal")
	if !ok {
		t.Fatal("expected manifest discovery")
	}
	if path != "project/playbook-manifest.json" {
		t.Fatalf("unexpected manifest path %q", path)
	}
}

func TestValidate_DetectsDuplicates(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	writeManifestFixtures(fs)
	m := sampleManifest()
	m.Playbooks = append(m.Playbooks, m.Playbooks[0])

	var validationErr *ValidationError
	if err := Validate(m, fs); !errors.As(err, &validationErr) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestValidate_DetectsInvalidReferences(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	writeManifestFixtures(fs)
	m := sampleManifest()
	m.Workflows[0].RequiredPlaybooks = append(m.Workflows[0].RequiredPlaybooks, "missing")

	var validationErr *ValidationError
	if err := Validate(m, fs); !errors.As(err, &validationErr) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestValidate_DetectsUnsupportedVersion(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	writeManifestFixtures(fs)
	m := sampleManifest()
	m.Metadata.SchemaVersion = "9.9.9"

	if err := Validate(m, fs); err == nil {
		t.Fatal("expected unsupported version error")
	}
}

func TestQueryEngine_ResolvesKnowledge(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	writeManifestFixtures(fs)
	engine := NewEngineWithFS(fs)
	if err := engine.Register("project", SourceProject, sampleManifest()); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	query := NewQueryEngine(engine)

	standards := query.StandardsByTechnology("Laravel")
	if len(standards) != 1 || standards[0].ID != "laravel-standard" {
		t.Fatalf("expected Laravel standard playbook, got %#v", standards)
	}
	playbooks := query.PlaybooksByTask("feature")
	if len(playbooks) != 1 || playbooks[0].ID != "feature-playbook" {
		t.Fatalf("expected feature playbook, got %#v", playbooks)
	}
	prompts := query.PromptsByWorkflow("feature")
	if len(prompts) != 1 || prompts[0].TaskType != "feature" {
		t.Fatalf("expected feature prompt, got %#v", prompts)
	}
	templates := query.TemplatesByType("architecture")
	if len(templates) != 1 || templates[0].Name != "Architecture" {
		t.Fatalf("expected architecture template, got %#v", templates)
	}
}

func TestEngine_PluginExtensionOverrides(t *testing.T) {
	fs := filesystem.NewMockFileSystem()
	writeManifestFixtures(fs)
	engine := NewEngineWithFS(fs)
	if err := engine.Register("core", SourceCore, sampleManifest()); err != nil {
		t.Fatalf("register core failed: %v", err)
	}
	plugin := sampleManifest()
	plugin.Metadata.Name = "Plugin"
	plugin.Playbooks[0].Name = "Plugin Laravel Standard"
	plugin.Extensions = map[string][]ExtensionResource{
		"plugin.example": {{Kind: "metadata", ID: "plugin.example.capability"}},
	}
	if err := engine.RegisterPluginManifest("plugin.example", plugin); err != nil {
		t.Fatalf("register plugin failed: %v", err)
	}

	active := engine.ActiveManifest()
	if active.Playbooks[0].Name != "Plugin Laravel Standard" {
		t.Fatalf("expected plugin override, got %q", active.Playbooks[0].Name)
	}
	if len(active.Extensions["plugin.example"]) != 1 {
		t.Fatal("expected plugin extension resource")
	}
}

func TestManifestEngine_LegacyMergeAndCompatibility(t *testing.T) {
	engine := NewEngine()
	engine.RegisterMemoryManifest("core", &DeclarativeManifest{
		SchemaVersion: "1.0",
		Workflows: map[string]WorkflowDef{
			"bug-fix": {ID: "bug-fix", Steps: []string{"fix"}},
		},
		Standards: map[string]StandardDef{
			"standard-lint": {ID: "standard-lint", Title: "Linting Rules", Priority: 50},
		},
		Compatibility: VersionCompatibilityDef{MinCLIVersion: "0.2.0"},
	})
	engine.RegisterMemoryManifest("project", &DeclarativeManifest{
		Workflows: map[string]WorkflowDef{
			"bug-fix": {ID: "bug-fix", Steps: []string{"project-fix"}},
		},
	})

	query := NewQueryEngine(engine)
	workflow, err := query.FindWorkflow("bug-fix")
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if workflow.Steps[0] != "project-fix" {
		t.Fatalf("expected project override, got %#v", workflow)
	}
	if query.VerifyCompatibility("0.1.0").IsCompatible {
		t.Fatal("expected old CLI version to be incompatible")
	}
}

func sampleManifest() *Manifest {
	return &Manifest{
		Metadata: ProjectMetadata{
			Name:          "Example",
			Version:       "1.0.0",
			SchemaVersion: SupportedSchemaVersion,
		},
		Playbooks: []PlaybookDefinition{
			{ID: "laravel-standard", Name: "Laravel Standard", Category: CategoryStacks, Location: "standards/laravel.md", Priority: 100},
			{ID: "feature-playbook", Name: "Feature Implementation", Category: CategoryWorkflows, Location: "workflows/feature.md", Priority: 90},
		},
		Technologies: []TechnologyDefinition{
			{ID: "laravel", Language: "php", Framework: "Laravel", Stack: "php-laravel", RelatedPlaybooks: []string{"laravel-standard"}},
		},
		Workflows: []WorkflowDefinition{
			{ID: "feature", Steps: []string{"plan", "implement", "verify"}, RequiredContext: []string{"requirements"}, RequiredPlaybooks: []string{"feature-playbook"}, Prompts: []string{"feature"}},
		},
		Prompts: []PromptMapping{
			{TaskType: "feature", PromptTemplate: "Feature", RequiredContext: []string{"requirements"}, AIFormattingRules: []string{"concise"}},
		},
		Templates: []TemplateDefinition{
			{Name: "Feature", Type: "prompt", Location: "prompts/feature.md", Version: "1.0.0", SupportedWorkflows: []string{"feature"}},
			{Name: "Architecture", Type: "architecture", Location: "project/templates/Architecture.template.md", Version: "1.0.0", Variables: []string{"name"}},
		},
		TaskRelationships: []TaskRelationship{
			{TaskType: "feature", RequiredWorkflow: "feature"},
		},
	}
}

func writeManifestFixtures(fs *filesystem.MockFileSystem) {
	fs.WriteFile("standards/laravel.md", []byte("laravel"), 0644)
	fs.WriteFile("workflows/feature.md", []byte("feature"), 0644)
	fs.WriteFile("prompts/feature.md", []byte("prompt"), 0644)
	fs.WriteFile("project/templates/Architecture.template.md", []byte("architecture"), 0644)
}

func validManifestJSON() string {
	return `{
		"metadata": {"name": "Example", "version": "1.0.0", "schema_version": "2.0.0"},
		"playbooks": [
			{"id": "feature-playbook", "name": "Feature", "category": "workflows", "location": "workflows/feature.md", "priority": 90}
		],
		"workflows": [
			{"id": "feature", "steps": ["plan"], "required_playbooks": ["feature-playbook"], "prompts": ["feature"]}
		],
		"prompts": [
			{"task_type": "feature", "prompt_template": "prompts/feature.md"}
		],
		"templates": [
			{"name": "Architecture", "type": "architecture", "location": "project/templates/Architecture.template.md", "version": "1.0.0"}
		]
	}`
}
