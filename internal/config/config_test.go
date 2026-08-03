package config

import (
	"os"
	"testing"
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
