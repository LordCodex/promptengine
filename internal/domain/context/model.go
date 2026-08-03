package context

// BudgetType represents available token capacity constraints
type BudgetType string

const (
	BudgetTiny      BudgetType = "tiny"      // ~1,000 tokens (5,000 bytes limit)
	BudgetSmall     BudgetType = "small"     // ~4,000 tokens (20,000 bytes limit)
	BudgetMedium    BudgetType = "medium"    // ~20,000 tokens (100,000 bytes limit)
	BudgetLarge     BudgetType = "large"     // ~100,000 tokens (500,000 bytes limit)
	BudgetUnlimited BudgetType = "unlimited" // no limits
)

// DocumentItem represents a spec or playbook candidate
type DocumentItem struct {
	Path        string  `json:"path"`
	Category    string  `json:"category"` // "business_rules", "architecture", "workflow", "stack", "generic"
	Content     string  `json:"-"`
	Size        int     `json:"size"` // in bytes
	Score       float64 `json:"score"`
	Explanation string  `json:"explanation"`
}

// OptimizationSummary provides metadata diagnostics
type OptimizationSummary struct {
	InitialCount   int     `json:"initial_count"`
	FinalCount     int     `json:"final_count"`
	InitialSize    int     `json:"initial_size"`
	FinalSize      int     `json:"final_size"`
	DroppedFiles   []string `json:"dropped_files"`
	BudgetLimit    int     `json:"budget_limit"`
	Deduplicated   bool    `json:"deduplicated"`
}

// ContextPackage is the complete payload delivered to prompt formatters
type ContextPackage struct {
	TaskType     string              `json:"task_type"`
	BudgetType   BudgetType          `json:"budget_type"`
	SystemPrompt string              `json:"system_prompt"`
	Documents    []DocumentItem      `json:"documents"`
	Explanations map[string]string   `json:"explanations"` // filepath -> explanation logs
	Summary      OptimizationSummary `json:"optimization_summary"`
}

func NewContextPackage(task string, budget BudgetType) *ContextPackage {
	return &ContextPackage{
		TaskType:     task,
		BudgetType:   budget,
		Documents:    make([]DocumentItem, 0),
		Explanations: make(map[string]string),
	}
}

// GetBudgetLimit returns the byte limit mapping the token capacity
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
