package template

import (
	"fmt"
	"strings"
	"sync"
)

// TemplateSource identifies where a template originates
type TemplateSource string

const (
	SourceCore        TemplateSource = "core"
	SourceOrg         TemplateSource = "org"
	SourcePlugin      TemplateSource = "plugin"
	SourceTechnology  TemplateSource = "technology"
	SourceMarketplace TemplateSource = "marketplace"
)

// Section is a conditional or unconditional block inside a template
type Section struct {
	Heading   string
	Body      string
	Condition string // variable name that must be non-empty to include; empty = always include
}

// Template is the full definition of a renderable document template
type Template struct {
	ID         string
	Name       string
	Source     TemplateSource
	Version    string
	ParentID   string // inheritance: parent template ID; empty = root
	Variables  []string
	Sections   []Section
	Partials   map[string]string // name -> template ID to inline
}

// TemplateRegistry manages all available templates
type TemplateRegistry struct {
	mu        sync.RWMutex
	templates map[string]*Template
}

func NewTemplateRegistry() *TemplateRegistry {
	return &TemplateRegistry{templates: make(map[string]*Template)}
}

func (r *TemplateRegistry) Register(t *Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.templates[t.ID]; exists {
		return fmt.Errorf("template '%s' already registered", t.ID)
	}
	r.templates[t.ID] = t
	return nil
}

func (r *TemplateRegistry) Get(id string) (*Template, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[id]
	return t, ok
}

// Renderer resolves variables, partials, and inheritance before returning filled content
type Renderer struct {
	registry *TemplateRegistry
}

func NewRenderer(reg *TemplateRegistry) *Renderer {
	return &Renderer{registry: reg}
}

func (r *Renderer) Render(id string, vars map[string]string) (string, error) {
	tmpl, ok := r.registry.Get(id)
	if !ok {
		return "", fmt.Errorf("template '%s' not found", id)
	}

	// Resolve inherited parent first
	base := ""
	if tmpl.ParentID != "" {
		parentContent, err := r.Render(tmpl.ParentID, vars)
		if err != nil {
			return "", fmt.Errorf("failed to render parent template '%s': %w", tmpl.ParentID, err)
		}
		base = parentContent + "\n\n"
	}

	var sb strings.Builder
	sb.WriteString(base)

	for _, section := range tmpl.Sections {
		// Conditional section inclusion
		if section.Condition != "" {
			if val, ok := vars[section.Condition]; !ok || val == "" {
				continue
			}
		}
		sb.WriteString(fmt.Sprintf("# %s\n\n", section.Heading))
		sb.WriteString(injectVars(section.Body, vars))
		sb.WriteString("\n\n")
	}

	// Inline partials
	result := sb.String()
	for partialName, partialID := range tmpl.Partials {
		partialContent, err := r.Render(partialID, vars)
		if err != nil {
			return "", fmt.Errorf("failed to render partial '%s': %w", partialName, err)
		}
		result = strings.ReplaceAll(result, "{{partial:"+partialName+"}}", partialContent)
	}

	return injectVars(result, vars), nil
}

func injectVars(text string, vars map[string]string) string {
	result := text
	for k, v := range vars {
		result = strings.ReplaceAll(result, "{"+k+"}", v)
	}
	return result
}
