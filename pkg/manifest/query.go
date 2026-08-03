package manifest

import (
	"fmt"
	"strings"
)

// QueryEngine isolates manifest querying logic
type QueryEngine struct {
	engine *Engine
}

func NewQueryEngine(e *Engine) *QueryEngine {
	return &QueryEngine{engine: e}
}

func (q *QueryEngine) FindWorkflow(id string) (WorkflowDef, error) {
	m := q.engine.GetMergedManifest()
	w, ok := m.Workflows[id]
	if !ok {
		return WorkflowDef{}, fmt.Errorf("workflow '%s' not declared in active manifests", id)
	}
	return w, nil
}

func (q *QueryEngine) FindStandard(id string) (StandardDef, error) {
	m := q.engine.GetMergedManifest()
	s, ok := m.Standards[id]
	if !ok {
		return StandardDef{}, fmt.Errorf("standard '%s' not declared in active manifests", id)
	}
	return s, nil
}

func (q *QueryEngine) FindTech(id string) (TechDef, error) {
	m := q.engine.GetMergedManifest()
	t, ok := m.Technologies[id]
	if !ok {
		return TechDef{}, fmt.Errorf("technology stack profile '%s' not declared in active manifests", id)
	}
	return t, nil
}

// FindPromptsForWorkflow returns prompts mapped to the task type
func (q *QueryEngine) FindPromptsForWorkflow(workflowID string) []PromptDef {
	m := q.engine.GetMergedManifest()
	var list []PromptDef
	for _, p := range m.Prompts {
		if p.WorkflowID == workflowID {
			list = append(list, p)
		}
	}
	return list
}

type CompatibilityReport struct {
	IsCompatible bool
	Reason       string
}

// VerifyCompatibility checks if the target CLI version matches manifest restrictions
func (q *QueryEngine) VerifyCompatibility(cliVersion string) CompatibilityReport {
	m := q.engine.GetMergedManifest()
	comp := m.Compatibility
	if comp.MinCLIVersion != "" {
		// Basic prefix check, standard semver checks done in future versions
		if strings.Compare(cliVersion, comp.MinCLIVersion) < 0 {
			return CompatibilityReport{
				IsCompatible: false,
				Reason:       fmt.Sprintf("CLI version %s is below required minimum %s", cliVersion, comp.MinCLIVersion),
			}
		}
	}
	return CompatibilityReport{IsCompatible: true, Reason: "Compatible"}
}
