package generator

import (
	"strings"
	"testing"
)

func TestGeneratorRegistry_RegisterAndGenerate(t *testing.T) {
	reg := NewGeneratorRegistry()
	RegisterDefaults(reg)

	// Verify architecture generator
	g, ok := reg.Get(DocArchitecture)
	if !ok {
		t.Fatal("expected architecture generator to be registered")
	}
	out, err := g.Generate(GeneratorInput{
		ProjectName: "TestProject",
		Stack:       []string{"Go", "PostgreSQL"},
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out.Content, "TestProject") {
		t.Error("expected generated content to contain project name")
	}
	if out.Filename == "" {
		t.Error("expected non-empty filename in output")
	}
}

func TestGeneratorRegistry_AllDefaultsRegistered(t *testing.T) {
	reg := NewGeneratorRegistry()
	RegisterDefaults(reg)

	expected := []DocType{
		DocArchitecture, DocBusinessRules, DocDatabase, DocAPI,
		DocPRD, DocProgress, DocRoadmap, DocDeployment,
		DocTroubleshooting, DocDecisions, DocSecurity, DocTesting,
	}
	for _, dt := range expected {
		if _, ok := reg.Get(dt); !ok {
			t.Errorf("expected generator for '%s' to be registered", dt)
		}
	}
}

func TestGeneratorRegistry_DuplicateBlocked(t *testing.T) {
	reg := NewGeneratorRegistry()
	_ = reg.Register(&architectureGenerator{})
	if err := reg.Register(&architectureGenerator{}); err == nil {
		t.Error("expected error on duplicate generator registration")
	}
}

// pluginGenerator simulates a third-party plugin contributing a custom generator
type pluginGenerator struct{}

func (p *pluginGenerator) DocType() DocType { return DocType("custom-doc") }
func (p *pluginGenerator) Generate(in GeneratorInput) (GeneratorOutput, error) {
	return GeneratorOutput{DocType: DocType("custom-doc"), Filename: "docs/Custom.md", Content: "# Custom\n\nPlugin-generated content."}, nil
}

func TestGeneratorRegistry_PluginContribution(t *testing.T) {
	reg := NewGeneratorRegistry()
	if err := reg.Register(&pluginGenerator{}); err != nil {
		t.Fatalf("expected plugin generator to register without error, got: %v", err)
	}
	g, ok := reg.Get("custom-doc")
	if !ok {
		t.Fatal("expected plugin generator to be discoverable")
	}
	out, _ := g.Generate(GeneratorInput{ProjectName: "Proj"})
	if !strings.Contains(out.Content, "Plugin-generated") {
		t.Error("expected plugin generator content")
	}
}
