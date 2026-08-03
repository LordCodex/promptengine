package manifest

import "fmt"

func (e *Engine) RegisterMemoryManifest(sourceName string, m *DeclarativeManifest) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.legacyManifests[sourceName] = m
	e.legacyDirty = true
}

func (e *Engine) GetMergedManifest() *DeclarativeManifest {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.legacyDirty && e.legacyMerged != nil {
		return e.legacyMerged
	}
	merged := &DeclarativeManifest{
		SchemaVersion: "1",
		Workflows:     map[string]WorkflowDef{},
		Standards:     map[string]StandardDef{},
		Technologies:  map[string]TechDef{},
		Prompts:       map[string]PromptDef{},
		HealthRules:   map[string]HealthRuleDef{},
	}
	for _, source := range []string{"core", "organization", "plugin", "project"} {
		if m, ok := e.legacyManifests[source]; ok {
			mergeLegacy(merged, m)
		}
	}
	for source, m := range e.legacyManifests {
		if source != "core" && source != "organization" && source != "plugin" && source != "project" {
			mergeLegacy(merged, m)
		}
	}
	e.legacyMerged = merged
	e.legacyDirty = false
	return merged
}

func mergeLegacy(dest, src *DeclarativeManifest) {
	if src.SchemaVersion != "" {
		dest.SchemaVersion = src.SchemaVersion
	}
	if src.Compatibility.MinCLIVersion != "" || src.Compatibility.ManifestSchemaVer != "" {
		dest.Compatibility = src.Compatibility
	}
	for k, v := range src.Workflows {
		dest.Workflows[k] = v
	}
	for k, v := range src.Standards {
		dest.Standards[k] = v
	}
	for k, v := range src.Technologies {
		dest.Technologies[k] = v
	}
	for k, v := range src.Prompts {
		dest.Prompts[k] = v
	}
	for k, v := range src.HealthRules {
		dest.HealthRules[k] = v
	}
}

func (q *QueryEngine) FindWorkflow(id string) (WorkflowDef, error) {
	if w, ok := q.engine.GetMergedManifest().Workflows[id]; ok {
		return w, nil
	}
	if w, ok := findWorkflow(q.engine.ActiveManifest(), id); ok {
		return w, nil
	}
	return WorkflowDef{}, fmt.Errorf("workflow %q not found", id)
}

func (q *QueryEngine) FindStandard(id string) (StandardDef, error) {
	m := q.engine.GetMergedManifest()
	if s, ok := m.Standards[id]; ok {
		return s, nil
	}
	return StandardDef{}, fmt.Errorf("standard %q not found", id)
}

func (q *QueryEngine) FindTech(id string) (TechDef, error) {
	m := q.engine.GetMergedManifest()
	if t, ok := m.Technologies[id]; ok {
		return t, nil
	}
	return TechDef{}, fmt.Errorf("technology %q not found", id)
}

func (q *QueryEngine) FindPromptsForWorkflow(workflowID string) []PromptDef {
	m := q.engine.GetMergedManifest()
	var out []PromptDef
	for _, p := range m.Prompts {
		if p.WorkflowID == workflowID {
			out = append(out, p)
		}
	}
	return out
}
