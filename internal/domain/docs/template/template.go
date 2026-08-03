package template

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/LordCodex/promptengine/internal/filesystem"
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
	ID        string
	Name      string
	Source    TemplateSource
	Version   string
	ParentID  string // inheritance: parent template ID; empty = root
	Variables []string
	Sections  []Section
	Partials  map[string]string // name -> template ID to inline
}

func (t *Template) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("template id is required")
	}
	if t.Version == "" {
		return fmt.Errorf("template %q version is required", t.ID)
	}
	if len(t.Sections) == 0 {
		return fmt.Errorf("template %q must define at least one section", t.ID)
	}
	return nil
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
		result = strings.ReplaceAll(result, "{{"+k+"}}", v)
	}
	return result
}

type Loader struct {
	fs      filesystem.FileSystem
	rootDir string
}

func NewLoader(fs filesystem.FileSystem, rootDir string) *Loader {
	return &Loader{fs: fs, rootDir: rootDir}
}

func (l *Loader) Discover() ([]string, error) {
	var paths []string
	if !l.fs.Exists(l.rootDir) {
		return paths, nil
	}
	entries, err := l.fs.ReadDir(l.rootDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".template.md") {
			paths = append(paths, filepath.Join(l.rootDir, entry.Name()))
		}
	}
	return paths, nil
}

func (l *Loader) LoadAll() ([]*Template, error) {
	paths, err := l.Discover()
	if err != nil {
		return nil, err
	}
	var templates []*Template
	for _, path := range paths {
		tmpl, err := l.Load(path)
		if err != nil {
			return nil, err
		}
		templates = append(templates, tmpl)
	}
	return templates, nil
}

func (l *Loader) Load(path string) (*Template, error) {
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSuffix(filepath.Base(path), ".template.md")
	id := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	return &Template{
		ID:        id,
		Name:      name,
		Source:    SourceCore,
		Version:   "1.0.0",
		Variables: extractVariables(string(data)),
		Sections:  []Section{{Heading: name, Body: string(data)}},
		Partials:  map[string]string{},
	}, nil
}

func extractVariables(content string) []string {
	seen := map[string]bool{}
	var out []string
	collect := func(open, close string) {
		for _, part := range strings.Split(content, open) {
			if !strings.Contains(part, close) {
				continue
			}
			key := strings.TrimSpace(strings.SplitN(part, close, 2)[0])
			if key == "" || strings.HasPrefix(key, "partial:") || strings.ContainsAny(key, "\n\t") || seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, key)
		}
	}
	collect("{{", "}}")
	collect("{", "}")
	return out
}
