package config

import (
	stderrors "errors"
	"os"
	"path/filepath"
	"testing"

	apperrors "github.com/LordCodex/promptengine/internal/errors"
)

func TestNewConfigLoader_Defaults(t *testing.T) {
	loader := NewConfigLoader()
	cfg, err := loader.Load("", "")
	if err != nil {
		t.Fatalf("Expected no error loading default configurations, got: %v", err)
	}

	if cfg.Project.Name != "PromptEngineProject" {
		t.Errorf("Expected default project name 'PromptEngineProject', got '%s'", cfg.Project.Name)
	}

	if cfg.Docs.RootDir != "docs" {
		t.Errorf("Expected default docs root dir 'docs', got '%s'", cfg.Docs.RootDir)
	}
}

func TestConfigLoader_Precedence(t *testing.T) {
	t.Setenv("PROMPTENGINE_PROJECT_NAME", "from-env")

	dir := t.TempDir()
	globalPath := filepath.Join(dir, "global.yaml")
	projectPath := filepath.Join(dir, "project.yaml")

	if err := os.WriteFile(globalPath, []byte("project:\n  name: from-global\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(projectPath, []byte("project:\n  name: from-project\ncli:\n  verbose: false\n"), 0644); err != nil {
		t.Fatal(err)
	}

	verbose := true
	loader := NewConfigLoader()
	cfg, err := loader.LoadWithFlags(projectPath, globalPath, Flags{Verbose: &verbose})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Project.Name != "from-env" {
		t.Fatalf("expected env to override project/global config, got %q", cfg.Project.Name)
	}
	if !cfg.CLI.Verbose {
		t.Fatal("expected CLI flag to override file config")
	}
}

func TestConfigLoader_EnvOverrides(t *testing.T) {
	os.Setenv("PROMPTENGINE_PROJECT_NAME", "EnvOverriddenName")
	defer os.Unsetenv("PROMPTENGINE_PROJECT_NAME")

	loader := NewConfigLoader()
	cfg, err := loader.Load("", "")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if cfg.Project.Name != "EnvOverriddenName" {
		t.Errorf("Expected env-overridden project name 'EnvOverriddenName', got '%s'", cfg.Project.Name)
	}
}

func TestConfigLoader_InvalidYAMLReturnsStructuredError(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, ".promptengine.yaml")
	if err := os.WriteFile(projectPath, []byte("profile: [broken\n"), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewConfigLoader()
	_, err := loader.LoadWithFlags(projectPath, "", Flags{})
	if err == nil {
		t.Fatal("expected invalid YAML to fail")
	}
	var appErr *apperrors.AppError
	if !stderrors.As(err, &appErr) {
		t.Fatalf("expected structured app error, got %T: %v", err, err)
	}
	if appErr.Category != apperrors.CategoryConfiguration {
		t.Fatalf("expected configuration category, got %s", appErr.Category)
	}
	if appErr.ExitCode() != apperrors.ExitConfiguration {
		t.Fatalf("expected configuration exit code, got %d", appErr.ExitCode())
	}
	if appErr.Recommendation == "" {
		t.Fatal("expected remediation recommendation")
	}
}
