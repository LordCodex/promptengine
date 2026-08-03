package template

import (
	"strings"
	"testing"
)

func TestRenderer_BasicRender(t *testing.T) {
	reg := NewTemplateRegistry()
	_ = reg.Register(&Template{
		ID:      "architecture",
		Name:    "Architecture",
		Source:  SourceCore,
		Version: "1.0",
		Sections: []Section{
			{Heading: "Overview", Body: "System overview for {project_name}."},
			{Heading: "Stack", Body: "Built with {stack}."},
		},
	})

	renderer := NewRenderer(reg)
	out, err := renderer.Render("architecture", map[string]string{
		"project_name": "MyApp",
		"stack":        "Go + PostgreSQL",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "MyApp") {
		t.Error("expected rendered output to contain project name")
	}
	if !strings.Contains(out, "Go + PostgreSQL") {
		t.Error("expected rendered output to contain stack")
	}
}

func TestRenderer_ConditionalSection(t *testing.T) {
	reg := NewTemplateRegistry()
	_ = reg.Register(&Template{
		ID: "prd",
		Sections: []Section{
			{Heading: "Always", Body: "always present"},
			{Heading: "Optional", Body: "only when feature_flag set", Condition: "feature_flag"},
		},
	})

	renderer := NewRenderer(reg)

	// Without condition variable — optional section should be absent
	out, _ := renderer.Render("prd", map[string]string{})
	if strings.Contains(out, "Optional") {
		t.Error("expected conditional section to be omitted when variable is absent")
	}

	// With condition variable — optional section should appear
	out, _ = renderer.Render("prd", map[string]string{"feature_flag": "true"})
	if !strings.Contains(out, "Optional") {
		t.Error("expected conditional section to appear when variable is set")
	}
}

func TestRenderer_Inheritance(t *testing.T) {
	reg := NewTemplateRegistry()
	_ = reg.Register(&Template{
		ID: "base",
		Sections: []Section{
			{Heading: "Introduction", Body: "Base intro."},
		},
	})
	_ = reg.Register(&Template{
		ID:       "extended",
		ParentID: "base",
		Sections: []Section{
			{Heading: "Extension", Body: "Child-specific content."},
		},
	})

	renderer := NewRenderer(reg)
	out, err := renderer.Render("extended", map[string]string{})
	if err != nil {
		t.Fatalf("expected no error on inherited template, got: %v", err)
	}
	if !strings.Contains(out, "Introduction") {
		t.Error("expected inherited section 'Introduction' from parent")
	}
	if !strings.Contains(out, "Extension") {
		t.Error("expected child section 'Extension'")
	}
}

func TestTemplateRegistry_DuplicateBlocked(t *testing.T) {
	reg := NewTemplateRegistry()
	tmpl := &Template{ID: "arch"}
	_ = reg.Register(tmpl)
	if err := reg.Register(tmpl); err == nil {
		t.Error("expected error on duplicate template registration")
	}
}
