package manifest

import (
	"fmt"
	"strings"
)

type QueryEngine struct {
	engine *Engine
}

func NewQueryEngine(e *Engine) *QueryEngine {
	return &QueryEngine{engine: e}
}

func (q *QueryEngine) StandardsByTechnology(name string) []PlaybookDefinition {
	m := q.engine.ActiveManifest()
	needle := strings.ToLower(name)
	playbookIDs := map[string]bool{}
	for _, tech := range m.Technologies {
		if techMatches(tech, needle) {
			for _, id := range tech.RelatedPlaybooks {
				playbookIDs[id] = true
			}
		}
	}
	var out []PlaybookDefinition
	for _, playbook := range m.Playbooks {
		if playbookIDs[playbook.ID] {
			out = append(out, playbook)
		}
	}
	return out
}

func (q *QueryEngine) PlaybooksByTask(taskType string) []PlaybookDefinition {
	m := q.engine.ActiveManifest()
	workflowID := ""
	for _, rel := range m.TaskRelationships {
		if strings.EqualFold(rel.TaskType, taskType) {
			workflowID = rel.RequiredWorkflow
			break
		}
	}
	if workflowID == "" {
		workflowID = normalizeTask(taskType)
	}

	workflow, ok := findWorkflow(m, workflowID)
	if !ok {
		return nil
	}
	byID := map[string]PlaybookDefinition{}
	for _, playbook := range m.Playbooks {
		byID[playbook.ID] = playbook
	}
	out := make([]PlaybookDefinition, 0, len(workflow.RequiredPlaybooks))
	for _, id := range workflow.RequiredPlaybooks {
		if playbook, ok := byID[id]; ok {
			out = append(out, playbook)
		}
	}
	return out
}

func (q *QueryEngine) PromptsByWorkflow(workflowID string) []PromptMapping {
	m := q.engine.ActiveManifest()
	workflow, ok := findWorkflow(m, workflowID)
	if !ok {
		return nil
	}
	allowed := map[string]bool{}
	for _, prompt := range workflow.Prompts {
		allowed[prompt] = true
	}
	var out []PromptMapping
	for _, prompt := range m.Prompts {
		if allowed[prompt.TaskType] || strings.EqualFold(prompt.TaskType, workflowID) {
			out = append(out, prompt)
		}
	}
	return out
}

func (q *QueryEngine) TemplatesByType(templateType string) []TemplateDefinition {
	m := q.engine.ActiveManifest()
	needle := strings.ToLower(templateType)
	var out []TemplateDefinition
	for _, tmpl := range m.Templates {
		if strings.EqualFold(tmpl.Type, templateType) || strings.Contains(strings.ToLower(tmpl.Name), needle) {
			out = append(out, tmpl)
		}
	}
	return out
}

func (q *QueryEngine) WorkflowRequirements(workflowID string) (WorkflowDefinition, []PlaybookDefinition, []PromptMapping, error) {
	workflow, ok := findWorkflow(q.engine.ActiveManifest(), workflowID)
	if !ok {
		return WorkflowDefinition{}, nil, nil, fmt.Errorf("workflow %q not found", workflowID)
	}
	return workflow, q.PlaybooksByTask(workflow.ID), q.PromptsByWorkflow(workflow.ID), nil
}

func findWorkflow(m *Manifest, id string) (WorkflowDefinition, bool) {
	for _, workflow := range m.Workflows {
		if strings.EqualFold(workflow.ID, id) || normalizeTask(workflow.ID) == normalizeTask(id) {
			return workflow, true
		}
	}
	return WorkflowDefinition{}, false
}

func techMatches(tech TechnologyDefinition, needle string) bool {
	values := []string{tech.ID, tech.Language, tech.Framework, tech.Stack}
	for _, value := range values {
		if strings.EqualFold(value, needle) || strings.Contains(strings.ToLower(value), needle) {
			return true
		}
	}
	return false
}

func normalizeTask(task string) string {
	normalized := strings.ReplaceAll(strings.ToLower(task), " ", "_")
	return strings.ReplaceAll(normalized, "-", "_")
}

type CompatibilityReport struct {
	IsCompatible bool
	Reason       string
}

func (q *QueryEngine) VerifyCompatibility(cliVersion string) CompatibilityReport {
	q.engine.mu.RLock()
	hasActiveSources := len(q.engine.sources) > 0
	q.engine.mu.RUnlock()

	m := q.engine.ActiveManifest()
	if hasActiveSources && m.Metadata.SchemaVersion != "" {
		if m.Metadata.SchemaVersion != SupportedSchemaVersion {
			return CompatibilityReport{IsCompatible: false, Reason: fmt.Sprintf("unsupported manifest schema version %s", m.Metadata.SchemaVersion)}
		}
		return CompatibilityReport{IsCompatible: true, Reason: "Compatible"}
	}

	legacy := q.engine.GetMergedManifest()
	if legacy.Compatibility.MinCLIVersion != "" && strings.Compare(cliVersion, legacy.Compatibility.MinCLIVersion) < 0 {
		return CompatibilityReport{
			IsCompatible: false,
			Reason:       fmt.Sprintf("CLI version %s is below required minimum %s", cliVersion, legacy.Compatibility.MinCLIVersion),
		}
	}
	if legacy.SchemaVersion == "" {
		return CompatibilityReport{IsCompatible: false, Reason: "manifest schema version is missing"}
	}
	return CompatibilityReport{IsCompatible: true, Reason: "Compatible"}
}
