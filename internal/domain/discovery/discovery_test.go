package discovery

import (
	"context"
	"testing"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

func TestDiscoveryPipeline_LaravelVue(t *testing.T) {
	fs := filesystem.NewMockFileSystem()

	// Write mock project files
	_ = fs.WriteFile("composer.json", []byte(`{"require": {"php": "^8.0"}}`), 0644)
	_ = fs.WriteFile("package.json", []byte(`{"dependencies": {"vue": "^3.0"}}`), 0644)
	_ = fs.WriteFile("artisan", []byte("cli entry"), 0644)
	_ = fs.WriteFile(".env", []byte("DB_CONNECTION=mysql"), 0644)
	_ = fs.WriteFile("app/Http/Controllers/Controller.php", []byte("controller"), 0644)
	_ = fs.WriteFile("app/Models/User.php", []byte("model"), 0644)
	_ = fs.WriteFile("resources/views/welcome.blade.php", []byte("view"), 0644)

	pipeline := NewPipeline()
	pipeline.Register(&BaseStage{}, &PromptEngineStage{}, &TechStage{}, &ArchStage{}, &DocsStage{})

	pm, err := pipeline.Execute(context.Background(), fs, "")
	if err != nil {
		t.Fatalf("Expected no error executing pipeline, got: %v", err)
	}

	// Verify Languages
	hasPHP := false
	hasJS := false
	for _, l := range pm.Languages {
		if l == "PHP" {
			hasPHP = true
		}
		if l == "JavaScript" {
			hasJS = true
		}
	}
	if !hasPHP || !hasJS {
		t.Errorf("Expected PHP and JavaScript, got %v", pm.Languages)
	}

	// Verify Frameworks
	hasLaravel := false
	hasVue := false
	for _, f := range pm.Frameworks {
		if f == "Laravel" {
			hasLaravel = true
		}
		if f == "Vue" {
			hasVue = true
		}
	}
	if !hasLaravel || !hasVue {
		t.Errorf("Expected Laravel and Vue frameworks, got %v", pm.Frameworks)
	}

	// Verify Databases
	if len(pm.Databases) == 0 || pm.Databases[0] != "MySQL" {
		t.Errorf("Expected MySQL connection, got %v", pm.Databases)
	}

	// Verify Architecture inference
	hasMVC := false
	for _, arch := range pm.Architectures {
		if arch.Style == "MVC" {
			hasMVC = true
			if arch.Confidence < 0.8 {
				t.Errorf("Expected confidence >= 0.8 for MVC, got %f", arch.Confidence)
			}
		}
	}
	if !hasMVC {
		t.Errorf("Expected MVC style to be inferred, got %v", pm.Architectures)
	}
}

func TestDiscoveryPipeline_GoCleanArch(t *testing.T) {
	fs := filesystem.NewMockFileSystem()

	_ = fs.WriteFile("go.mod", []byte("module myapp"), 0644)
	_ = fs.WriteFile("internal/domain/user.go", []byte("domain code"), 0644)
	_ = fs.WriteFile("internal/usecase/auth.go", []byte("usecase code"), 0644)
	_ = fs.WriteFile("internal/filesystem/storage.go", []byte("adapter code"), 0644)

	pipeline := NewPipeline()
	pipeline.Register(&BaseStage{}, &PromptEngineStage{}, &TechStage{}, &ArchStage{}, &DocsStage{})

	pm, err := pipeline.Execute(context.Background(), fs, "")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	hasClean := false
	for _, arch := range pm.Architectures {
		if arch.Style == "Clean Architecture" {
			hasClean = true
		}
	}
	if !hasClean {
		t.Errorf("Expected Clean Architecture styles, got %v", pm.Architectures)
	}
}
