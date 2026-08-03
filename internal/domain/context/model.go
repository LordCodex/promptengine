package context

import "github.com/LordCodex/promptengine/internal/domain/discovery"

type BudgetType string

const (
	BudgetTiny      BudgetType = "tiny"
	BudgetSmall     BudgetType = "small"
	BudgetMedium    BudgetType = "medium"
	BudgetLarge     BudgetType = "large"
	BudgetUnlimited BudgetType = "unlimited"
)

type ContextSourceType string

const (
	ContextFile          ContextSourceType = "file"
	ContextDocumentation ContextSourceType = "documentation"
	ContextManifestEntry ContextSourceType = "manifest_entry"
	ContextStandard      ContextSourceType = "standard"
	ContextWorkflow      ContextSourceType = "workflow"
	ContextTemplate      ContextSourceType = "template"
)

type InclusionLevel string

const (
	InclusionRequired    InclusionLevel = "required"
	InclusionConditional InclusionLevel = "conditional"
	InclusionExcluded    InclusionLevel = "excluded"
)

type ContextRequest struct {
	TaskType           TaskType                `json:"task_type"`
	WorkflowType       string                  `json:"workflow_type,omitempty"`
	Project            *discovery.ProjectModel `json:"project,omitempty"`
	TechnologyStack    []string                `json:"technology_stack,omitempty"`
	UserIntent         string                  `json:"user_intent,omitempty"`
	RequestedOperation string                  `json:"requested_operation,omitempty"`
	AffectedFiles      []string                `json:"affected_files,omitempty"`
	MaxBytes           int                     `json:"max_bytes,omitempty"`
	MinRelevanceScore  float64                 `json:"min_relevance_score,omitempty"`
	Explain            bool                    `json:"explain,omitempty"`
	Budget             BudgetType              `json:"budget,omitempty"`
	Metadata           map[string]string       `json:"metadata,omitempty"`
}

type ContextItem struct {
	Path            string            `json:"path"`
	Type            ContextSourceType `json:"type"`
	RelevanceScore  float64           `json:"relevance_score"`
	Reason          string            `json:"reason_selected"`
	Size            int               `json:"size"`
	Content         string            `json:"-"`
	Summary         string            `json:"summary,omitempty"`
	Truncated       bool              `json:"truncated"`
	InclusionLevel  InclusionLevel    `json:"inclusion_level"`
	ExclusionReason string            `json:"exclusion_reason,omitempty"`
}

type DocumentItem struct {
	Path        string  `json:"path"`
	Category    string  `json:"category"`
	Content     string  `json:"-"`
	Size        int     `json:"size"`
	Score       float64 `json:"score"`
	Explanation string  `json:"explanation"`
}

type OptimizationSummary struct {
	InitialCount int      `json:"initial_count"`
	FinalCount   int      `json:"final_count"`
	InitialSize  int      `json:"initial_size"`
	FinalSize    int      `json:"final_size"`
	DroppedFiles []string `json:"dropped_files"`
	BudgetLimit  int      `json:"budget_limit"`
	Deduplicated bool     `json:"deduplicated"`
	CacheHit     bool     `json:"cache_hit"`
}

type ContextPackage struct {
	TaskType          string              `json:"task_type"`
	WorkflowType      string              `json:"workflow_type,omitempty"`
	BudgetType        BudgetType          `json:"budget_type"`
	SystemPrompt      string              `json:"system_prompt"`
	Items             []ContextItem       `json:"items"`
	ExcludedItems     []ContextItem       `json:"excluded_items,omitempty"`
	Documents         []DocumentItem      `json:"documents"`
	SelectedFiles     []string            `json:"selected_files"`
	SelectedDocs      []string            `json:"selected_documents"`
	RelevantStandards []string            `json:"relevant_standards"`
	RelatedPlaybooks  []string            `json:"related_playbooks"`
	ProjectMetadata   map[string]string   `json:"project_metadata,omitempty"`
	Reasoning         []string            `json:"reasoning"`
	Explanations      map[string]string   `json:"explanations,omitempty"`
	Summary           OptimizationSummary `json:"optimization_summary"`
	EstimatedTokens   int                 `json:"estimated_tokens"`
}

func NewContextPackage(task string, budget BudgetType) *ContextPackage {
	return &ContextPackage{
		TaskType:        task,
		BudgetType:      budget,
		Items:           []ContextItem{},
		Documents:       []DocumentItem{},
		ProjectMetadata: map[string]string{},
		Reasoning:       []string{},
		Explanations:    map[string]string{},
	}
}

func GetBudgetLimit(b BudgetType) int {
	switch b {
	case BudgetTiny:
		return 5000
	case BudgetSmall:
		return 20000
	case BudgetMedium:
		return 100000
	case BudgetLarge:
		return 500000
	default:
		return 999999999
	}
}
