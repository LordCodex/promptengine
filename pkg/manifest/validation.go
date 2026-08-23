package manifest

import (
	"fmt"
	"strings"

	"github.com/LordCodex/promptengine/internal/filesystem"
)

type ValidationIssue struct {
	Field   string `json:"field" yaml:"field"`
	Message string `json:"message" yaml:"message"`
}
type ValidationError struct {
	Issues []ValidationIssue `json:"issues" yaml:"issues"`
}

func (e *ValidationError) Error() string {
	if len(e.Issues) == 0 {
		return "manifest validation failed"
	}
	return fmt.Sprintf("manifest validation failed: %s", e.Issues[0].Message)
}

func Validate(m *Manifest, fs filesystem.FileSystem) error {
	var issues []ValidationIssue
	add := func(field, msg string) { issues = append(issues, ValidationIssue{Field: field, Message: msg}) }
	if m == nil {
		add("manifest", "manifest is nil")
		return &ValidationError{Issues: issues}
	}
	if m.Metadata.Name == "" {
		add("metadata.name", "project name is required")
	}
	if m.Metadata.Version == "" {
		add("metadata.version", "project version is required")
	}
	if m.Metadata.SchemaVersion == "" {
		add("metadata.schema_version", "schema version is required")
	} else if m.Metadata.SchemaVersion != SupportedSchemaVersion {
		add("metadata.schema_version", fmt.Sprintf("unsupported schema version %q", m.Metadata.SchemaVersion))
	}

	playbooks := map[string]PlaybookDefinition{}
	for i, p := range m.Playbooks {
		field := fmt.Sprintf("playbooks[%d]", i)
		if p.ID == "" {
			add(field+".id", "playbook identifier is required")
			continue
		}
		if _, exists := playbooks[p.ID]; exists {
			add(field+".id", fmt.Sprintf("duplicate playbook identifier %q", p.ID))
		}
		playbooks[p.ID] = p
		if p.Name == "" {
			add(field+".name", "playbook name is required")
		}
		if !validCategory(p.Category) {
			add(field+".category", fmt.Sprintf("invalid playbook category %q", p.Category))
		}
		if p.Location == "" {
			add(field+".location", "playbook location is required")
		} else if fs != nil && !fs.Exists(p.Location) {
			add(field+".location", fmt.Sprintf("playbook path %q does not exist", p.Location))
		}
	}

	workflows := map[string]WorkflowDefinition{}
	for i, w := range m.Workflows {
		field := fmt.Sprintf("workflows[%d]", i)
		if w.ID == "" {
			add(field+".id", "workflow identifier is required")
			continue
		}
		if _, exists := workflows[w.ID]; exists {
			add(field+".id", fmt.Sprintf("duplicate workflow identifier %q", w.ID))
		}
		workflows[w.ID] = w
		if len(w.Steps) == 0 {
			add(field+".steps", "workflow steps are required")
		}
		for _, id := range w.RequiredPlaybooks {
			if _, ok := playbooks[id]; !ok {
				add(field+".required_playbooks", fmt.Sprintf("missing playbook reference %q", id))
			}
		}
		for _, id := range w.OptionalPlaybooks {
			if _, ok := playbooks[id]; !ok {
				add(field+".optional_playbooks", fmt.Sprintf("missing playbook reference %q", id))
			}
		}
	}

	prompts := map[string]PromptMapping{}
	for i, p := range m.Prompts {
		field := fmt.Sprintf("prompts[%d]", i)
		if p.TaskType == "" {
			add(field+".task_type", "prompt task type is required")
			continue
		}
		if _, exists := prompts[p.TaskType]; exists {
			add(field+".task_type", fmt.Sprintf("duplicate prompt task type %q", p.TaskType))
		}
		prompts[p.TaskType] = p
		if p.PromptTemplate == "" {
			add(field+".prompt_template", "prompt template is required")
		} else if _, ok := findTemplate(m, p.PromptTemplate); !ok && fs != nil && !fs.Exists(p.PromptTemplate) {
			add(field+".prompt_template", fmt.Sprintf("prompt template %q is not registered and path does not exist", p.PromptTemplate))
		}
	}

	templates := map[string]TemplateDefinition{}
	for i, tmpl := range m.Templates {
		field := fmt.Sprintf("templates[%d]", i)
		if tmpl.Name == "" {
			add(field+".name", "template name is required")
			continue
		}
		key := strings.ToLower(tmpl.Name)
		if _, exists := templates[key]; exists {
			add(field+".name", fmt.Sprintf("duplicate template name %q", tmpl.Name))
		}
		templates[key] = tmpl
		if tmpl.Location == "" {
			add(field+".location", "template location is required")
		} else if fs != nil && !fs.Exists(tmpl.Location) {
			add(field+".location", fmt.Sprintf("template path %q does not exist", tmpl.Location))
		}
		for _, workflowID := range tmpl.SupportedWorkflows {
			if _, ok := workflows[workflowID]; !ok {
				add(field+".supported_workflows", fmt.Sprintf("missing workflow reference %q", workflowID))
			}
		}
	}

	techs := map[string]TechnologyDefinition{}
	for i, tech := range m.Technologies {
		field := fmt.Sprintf("technologies[%d]", i)
		key := strings.ToLower(firstNonEmpty(tech.ID, tech.Framework, tech.Language, tech.Stack))
		if key == "" {
			add(field+".id", "technology identifier, language, framework, or stack is required")
			continue
		}
		if _, exists := techs[key]; exists {
			add(field+".id", fmt.Sprintf("duplicate technology identifier %q", key))
		}
		techs[key] = tech
		for _, id := range tech.RelatedPlaybooks {
			if _, ok := playbooks[id]; !ok {
				add(field+".related_playbooks", fmt.Sprintf("missing playbook reference %q", id))
			}
		}
	}
	for i, mapping := range m.CommandMappings {
		if mapping.Command == "" {
			add(fmt.Sprintf("command_mappings[%d].command", i), "command is required")
		}
		if mapping.Workflow == "" {
			add(fmt.Sprintf("command_mappings[%d].workflow", i), "workflow is required")
		} else if _, ok := workflows[mapping.Workflow]; !ok {
			add(fmt.Sprintf("command_mappings[%d].workflow", i), fmt.Sprintf("missing workflow reference %q", mapping.Workflow))
		}
	}
	for i, rel := range m.TaskRelationships {
		if rel.TaskType == "" {
			add(fmt.Sprintf("task_relationships[%d].task_type", i), "task type is required")
		}
		if rel.RequiredWorkflow != "" {
			if _, ok := workflows[rel.RequiredWorkflow]; !ok {
				add(fmt.Sprintf("task_relationships[%d].required_workflow", i), fmt.Sprintf("missing workflow reference %q", rel.RequiredWorkflow))
			}
		}
	}
	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}

func validCategory(c PlaybookCategory) bool {
	switch c {
	case CategoryCore, CategoryStacks, CategorySecurity, CategoryPerformance, CategoryDesign, CategoryWorkflows, CategoryProject, CategoryBridge, CategoryChecklist, CategoryDecisionGuide, CategoryGuide, CategoryPrompt, CategoryAI, CategoryCLI:
		return true
	default:
		return false
	}
}
func findTemplate(m *Manifest, name string) (TemplateDefinition, bool) {
	for _, tmpl := range m.Templates {
		if strings.EqualFold(tmpl.Name, name) {
			return tmpl, true
		}
	}
	return TemplateDefinition{}, false
}
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
