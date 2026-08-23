package manifest

import (
	"fmt"
	"strings"
)

type QueryEngine struct { engine *Engine }
func NewQueryEngine(e *Engine) *QueryEngine { return &QueryEngine{engine: e} }

func (q *QueryEngine) StandardsByTechnology(name string) []PlaybookDefinition {
	m := q.engine.ActiveManifest()
	needle := strings.ToLower(name)
	ids := map[string]bool{}
	for _, tech := range m.Technologies {
		if techMatches(tech, needle) {
			for _, id := range tech.RelatedPlaybooks { ids[id] = true }
			for _, id := range tech.RelatedStandards { ids[id] = true }
		}
	}
	return playbooksByIDs(m, ids)
}

func (q *QueryEngine) PlaybooksByTask(taskType string) []PlaybookDefinition {
	workflow, ok := q.workflowForTask(taskType)
	if !ok { return nil }
	ids := map[string]bool{}
	for _, id := range workflow.RequiredPlaybooks { ids[id] = true }
	return playbooksByIDs(q.engine.ActiveManifest(), ids)
}

func (q *QueryEngine) OptionalPlaybooksByTask(taskType string) []PlaybookDefinition {
	workflow, ok := q.workflowForTask(taskType)
	if !ok { return nil }
	ids := map[string]bool{}
	for _, id := range workflow.OptionalPlaybooks { ids[id] = true }
	return playbooksByIDs(q.engine.ActiveManifest(), ids)
}

func (q *QueryEngine) PlaybooksByCategory(category PlaybookCategory) []PlaybookDefinition {
	var out []PlaybookDefinition
	for _, playbook := range q.engine.ActiveManifest().Playbooks {
		if playbook.Category == category { out = append(out, playbook) }
	}
	return out
}

// MatchingKnowledge returns non-required library resources whose identity or
// path matches concrete task terms. It is deliberately conservative: these are
// candidates for context ranking, never unconditional context.
func (q *QueryEngine) MatchingKnowledge(intent string, categories ...PlaybookCategory) []PlaybookDefinition {
	allowed := map[PlaybookCategory]bool{}
	for _, category := range categories { allowed[category] = true }
	terms := knowledgeTerms(intent)
	if len(terms) == 0 { return nil }
	var out []PlaybookDefinition
	for _, playbook := range q.engine.ActiveManifest().Playbooks {
		if len(allowed) > 0 && !allowed[playbook.Category] { continue }
		haystack := strings.ToLower(playbook.ID + " " + playbook.Name + " " + playbook.Location + " " + playbook.Description)
		for _, term := range terms {
			if strings.Contains(haystack, term) { out = append(out, playbook); break }
		}
	}
	return out
}

func (q *QueryEngine) PromptsByWorkflow(workflowID string) []PromptMapping {
	m := q.engine.ActiveManifest()
	workflow, ok := findWorkflow(m, workflowID)
	if !ok { return nil }
	allowed := map[string]bool{}
	for _, prompt := range workflow.Prompts { allowed[prompt] = true }
	var out []PromptMapping
	for _, prompt := range m.Prompts {
		if allowed[prompt.TaskType] || strings.EqualFold(prompt.TaskType, workflowID) { out = append(out, prompt) }
	}
	return out
}

func (q *QueryEngine) TemplatesByType(templateType string) []TemplateDefinition {
	m := q.engine.ActiveManifest()
	needle := strings.ToLower(templateType)
	var out []TemplateDefinition
	for _, tmpl := range m.Templates {
		if strings.EqualFold(tmpl.Type, templateType) || strings.Contains(strings.ToLower(tmpl.Name), needle) { out = append(out, tmpl) }
	}
	return out
}

func (q *QueryEngine) WorkflowRequirements(workflowID string) (WorkflowDefinition, []PlaybookDefinition, []PromptMapping, error) {
	workflow, ok := findWorkflow(q.engine.ActiveManifest(), workflowID)
	if !ok { return WorkflowDefinition{}, nil, nil, fmt.Errorf("workflow %q not found", workflowID) }
	return workflow, q.PlaybooksByTask(workflow.ID), q.PromptsByWorkflow(workflow.ID), nil
}

func (q *QueryEngine) workflowForTask(taskType string) (WorkflowDefinition, bool) {
	m := q.engine.ActiveManifest()
	workflowID := ""
	for _, rel := range m.TaskRelationships {
		if strings.EqualFold(normalizeTask(rel.TaskType), normalizeTask(taskType)) {
			workflowID = rel.RequiredWorkflow
			break
		}
	}
	if workflowID == "" { workflowID = normalizeTask(taskType) }
	return findWorkflow(m, workflowID)
}

func playbooksByIDs(m *Manifest, ids map[string]bool) []PlaybookDefinition {
	var out []PlaybookDefinition
	for _, playbook := range m.Playbooks {
		if ids[playbook.ID] { out = append(out, playbook) }
	}
	return out
}

func findWorkflow(m *Manifest, id string) (WorkflowDefinition, bool) {
	for _, workflow := range m.Workflows {
		if strings.EqualFold(workflow.ID, id) || normalizeTask(workflow.ID) == normalizeTask(id) { return workflow, true }
	}
	return WorkflowDefinition{}, false
}

func techMatches(tech TechnologyDefinition, needle string) bool {
	for _, value := range []string{tech.ID, tech.Language, tech.Framework, tech.Stack} {
		if strings.EqualFold(value, needle) || strings.Contains(strings.ToLower(value), needle) { return true }
	}
	return false
}

func normalizeTask(task string) string {
	normalized := strings.ReplaceAll(strings.ToLower(task), " ", "_")
	return strings.ReplaceAll(normalized, "-", "_")
}

func knowledgeTerms(value string) []string {
	stop := map[string]bool{"with":true,"from":true,"into":true,"this":true,"that":true,"have":true,"using":true,"feature":true,"project":true,"change":true,"update":true,"implement":true}
	seen := map[string]bool{}
	var out []string
	for _, term := range strings.FieldsFunc(strings.ToLower(value), func(r rune) bool { return (r < 'a' || r > 'z') && (r < '0' || r > '9') }) {
		if len(term) < 4 || stop[term] || seen[term] { continue }
		seen[term] = true
		out = append(out, term)
	}
	return out
}

type CompatibilityReport struct { IsCompatible bool; Reason string }
func (q *QueryEngine) VerifyCompatibility(cliVersion string) CompatibilityReport {
	q.engine.mu.RLock(); hasActiveSources := len(q.engine.sources) > 0; q.engine.mu.RUnlock()
	m := q.engine.ActiveManifest()
	if hasActiveSources && m.Metadata.SchemaVersion != "" {
		if m.Metadata.SchemaVersion != SupportedSchemaVersion { return CompatibilityReport{false, fmt.Sprintf("unsupported manifest schema version %s", m.Metadata.SchemaVersion)} }
		return CompatibilityReport{true, "Compatible"}
	}
	legacy := q.engine.GetMergedManifest()
	if legacy.Compatibility.MinCLIVersion != "" && strings.Compare(cliVersion, legacy.Compatibility.MinCLIVersion) < 0 { return CompatibilityReport{false, fmt.Sprintf("CLI version %s is below required minimum %s", cliVersion, legacy.Compatibility.MinCLIVersion)} }
	if legacy.SchemaVersion == "" { return CompatibilityReport{false, "manifest schema version is missing"} }
	return CompatibilityReport{true, "Compatible"}
}
